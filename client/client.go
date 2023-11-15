package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	// inspired by https://github.com/PatrickMatthiesen/DSYS-gRPC-template and https://articles.wesionary.team/grpc-console-chat-application-in-go-dd77a29bb5c3
	gRPC "github.com/JonasSkjodt/chitty-chat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "5400", "Tcp server")

var ServerConn *grpc.ClientConn //the server connection
var chatServer gRPC.ChatClient  // new chat server client

var vectorClock = []int32{0, 0} // vector clock for the client
var clientID = -1               // clientID is set to 1 by default
var hasher = fnv.New32()

func main() {
	//parse flag/arguments
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	f := setLog()
	defer f.Close()

	//connect to server and close the connection when program closes
	fmt.Println("--- join Server ---")
	ConnectToServer()
	//defer SendMessage("exit", ChatStream)
	defer ServerConn.Close()

	ChatStream, err := chatServer.MessageStream(context.Background())
	if err != nil {
		fmt.Printf("Error on receive: %v \n", err)
		log.Fatalf("Error on receive: %v", err)
	}
	hasher.Write([]byte(*clientsName))

	// Client will connect with a timestamp of 0 because client has not received it's id.
	// This will be updated when the reponse from the server is received.
	SendMessage(fmt.Sprint(hasher.Sum32()), ChatStream)

	//start the biding

	go listenForMessages(ChatStream)
	parseInput(ChatStream)
}

// connect to server
func ConnectToServer() {

	//dial options
	//the server is not using TLS, so we use insecure credentials
	//(should be fine for local testing but not in the real world)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	//dial the server, with the flag "server", to get a connection to it
	fmt.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%s", *serverPort), opts...)
	if err != nil {
		fmt.Printf("Fail to Dial : %v \n", err)
		log.Printf("Fail to Dial : %v", err)
		return
	}

	//for the chat implementation
	// makes a client from the server connection and saves the connection
	// and prints rather or not the connection was is READY
	chatServer = gRPC.NewChatClient(conn)
	ServerConn = conn
	fmt.Println("the connection is: ", conn.GetState().String())
	log.Println("the connection is: ", conn.GetState().String())
}

func parseInput(stream gRPC.Chat_MessageStreamClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Chitty Chat!")
	fmt.Println("--------------------")

	//Infinite loop to listen for clients input.
	for {
		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%v \n", err)
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		// message must be under 128 characters long
		if len(input) > 128 {
			println("Message is too long. Your message must be under 128 characters long")
			continue
		}

		if !conReady(chatServer) {
			fmt.Printf("Client %s: something was wrong with the connection to the server :(", *clientsName)
			log.Printf("Client %s: something was wrong with the connection to the server :(", *clientsName)
			continue
		}

		if input == "exit" {
			SendMessage("Participant "+*clientsName+" left chitty-chat", stream)
			//chatServer.DisconnectFromServer(stream.Context(), &gRPC.ClientName{ClientName: *clientsName})
			time.Sleep(1 * time.Second)
			os.Exit(1)
		} else {
			SendMessage(input, stream)
		}
	}
}

// Function which returns a true boolean if the connection to the server is ready, and false if it's not.
func conReady(s gRPC.ChatClient) bool {
	return ServerConn.GetState().String() == "READY"
}

// sets the logger to use a log.txt file instead of the console
func setLog() *os.File {
	if err := os.Truncate("log_"+*clientsName+".txt", 0); err != nil {
		fmt.Printf("Failed to truncate: %v \n", err)
		log.Printf("Failed to truncate: %v", err)
	}

	f, err := os.OpenFile("log_"+*clientsName+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

func SendMessage(content string, stream gRPC.Chat_MessageStreamClient) {
	if clientID != -1 {
		vectorClock[clientID]++
	}
	message := &gRPC.ChatMessage{
		Content:     content,
		ClientName:  *clientsName,
		VectorClock: vectorClock,
	}

	stream.Send(message)
}

// watch the god
func listenForMessages(stream gRPC.Chat_MessageStreamClient) {
	for {
		time.Sleep(1 * time.Second)
		if stream != nil {
			msg, err := stream.Recv()
			if err == io.EOF {
				fmt.Printf("Error: io.EOF in listenForMessages in client.go \n")
				log.Printf("Error: io.EOF in listenForMessages in client.go")
				break
			}
			if err != nil {
				fmt.Printf("%v \n", err)
				log.Fatalf("%v", err)
			}
			if strings.Contains(msg.Content, *clientsName+" joined chitty-chat") {
				// Updates the clientID
				clientID = int(msg.ClientID)
			}

			//Updates the clients vector clock
			updateVectorClock(msg.VectorClock)
			if msg.ClientName != *clientsName {
				fmt.Printf("%s: \"%s\" at lamport timestamp: %d \n", msg.ClientName, msg.Content, vectorClock)
				log.Printf("%s: \"%s\" at lamport timestamp: %d", msg.ClientName, msg.Content, vectorClock)
			}

		}
	}
}

func updateVectorClock(msgVectorClock []int32) {
	for i := 0; i < len(msgVectorClock); i++ {
		// Add values to vectorClock if different lengths
		if len(msgVectorClock) >= len(vectorClock) {
			var lenDiff int = len(msgVectorClock) - len(vectorClock)
			for j := 0; j < lenDiff; j++ {
				vectorClock = append(vectorClock, 0)
			}
		}
		// Compare and update values
		if vectorClock[i] < msgVectorClock[i] {
			vectorClock[i] = msgVectorClock[i]
		}
	}
	vectorClock[clientID]++
}
