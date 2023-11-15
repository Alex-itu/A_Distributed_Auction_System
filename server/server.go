package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	// inspired by https://github.com/PatrickMatthiesen/DSYS-gRPC-template and https://articles.wesionary.team/grpc-console-chat-application-in-go-dd77a29bb5c3
	gRPC "github.com/JonasSkjodt/chitty-chat/proto"
	"google.golang.org/grpc"
)

type chatServer struct {
	gRPC.UnimplementedChatServer        // You need this line if you have a server
	name                         string // Not required but useful if you want to name your server
	port                         string // Not required but useful if your server needs to know what port it's listening to

	mutex sync.Mutex // used to lock the server to avoid race conditions.

}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var serverName = flag.String("name", "default", "Senders name") // set with "-name <name>" in terminal
var port = flag.String("port", "5400", "Server port")           // set with "-port <port>" in terminal
var vectorClock = []int32{0}
var clientID = 1
var nextNumClock = 0 // vector clock for the server

// Maps
var clientNames = make(map[string]gRPC.Chat_MessageStreamServer)
var clientIDs = make(map[string]int)

func main() {

	f := setLog() //uncomment this line to log to a log.txt file instead of the console
	defer f.Close()

	// This parses the flags and sets the correct/given corresponding values.
	flag.Parse()
	fmt.Println(".:server is starting:.")

	// launch the server
	launchServer()

	// code here is unreachable because launchServer occupies the current thread.
}

func launchServer() {
	fmt.Printf("Server %s: Attempts to create listener on port %s\n", *serverName, *port)
	log.Printf("Server %s: Attempts to create listener on port %s\n", *serverName, *port)

	// Create listener tcp on given port or default port 5400
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		fmt.Printf("Server %s: Failed to listen on port %s: %v \n", *serverName, *port, err)
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// makes a new server instance using the name and port from the flags.
	server := &chatServer{
		name: *serverName,
		port: *port,
	}

	gRPC.RegisterChatServer(grpcServer, server) //Registers the server to the gRPC server.

	fmt.Printf("Server %s: Listening at %v \n", *serverName, list.Addr())
	log.Printf("Server %s: Listening at %v \n", *serverName, list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		fmt.Printf("failed to serve %v", err)
		log.Fatalf("failed to serve %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}

func DeleteUser(clientName string) {
	if clientName != "" {
		//Deletes the client from the clientNames map
		delete(clientNames, clientName)
	}
}

func (s *chatServer) MessageStream(msgStream gRPC.Chat_MessageStreamServer) error {
	for {
		// get the next message from the stream
		msg, err := msgStream.Recv()
		if err == io.EOF {
			break
		}
		// some other error
		if err != nil {
			return err
		}
		hasher := fnv.New32()
		hasher.Write([]byte(msg.ClientName))
		if msg.Content == fmt.Sprint(hasher.Sum32()) {
			clientNames[msg.ClientName] = msgStream
			clientIDs[msg.ClientName] = clientID
			clientID++

			//Adds the client to the vector clock
			vectorClock = append(vectorClock, 1)
			UpdateVectorClock(msg.VectorClock)

			fmt.Printf("Participant %s joined chitty-chat at lamport timestamp: %d \n", msg.ClientName, vectorClock)
			log.Printf("Participant %s joined chitty-chat at lamport timestamp: %d", msg.ClientName, vectorClock)

			//Sends the message that a client has connected to the other clients
			SendMessages(&gRPC.ChatMessage{VectorClock: vectorClock, ClientID: int32(clientIDs[msg.ClientName]), ClientName: "Server", Content: fmt.Sprintf("Participant %s joined chitty-chat", msg.ClientName)})

			hasher = nil

		} else if strings.Contains(msg.Content, "Participant "+msg.ClientName+" left chitty-chat") {
			// Counts the clients vector clock up
			UpdateVectorClock(msg.VectorClock)

			fmt.Printf("Participant %s left chitty-chat at lamport timestamp: %d \n", msg.ClientName, vectorClock)
			log.Printf("Participant %s left chitty-chat At lamport timestamp: %d", msg.ClientName, vectorClock)

			//Adds the vector clock to the message
			msg.VectorClock = vectorClock

			DeleteUser(msg.ClientName)

			// send the message to all clients
			SendMessages(msg)

		} else {
			// Counts the clients vector clock up
			UpdateVectorClock(msg.VectorClock)

			// the stream is closed so we can exit the loop
			// log the message
			fmt.Printf("Received message: from %s: \"%s\" At lamport timestamp: %d \n", msg.ClientName, msg.Content, vectorClock)
			log.Printf("Received message: from %s: \"%s\" At lamport timestamp: %d", msg.ClientName, msg.Content, vectorClock)

			//Adds the vector clock to the message
			msg.VectorClock = vectorClock

			// send the message to all clients
			SendMessages(msg)

		}
	}

	return nil
}

func SendMessages(msg *gRPC.ChatMessage) {
	for name := range clientNames {
		if msg.ClientName != name {
			vectorClock[0]++
			msg.VectorClock = vectorClock
			clientNames[name].Send(msg)
		}
	}
}

func UpdateVectorClock(msgVectorClock []int32) {
	for i := 0; i < len(vectorClock); i++ {
		// Add dummy values to msgVectorClock so that values can be compared
		if len(msgVectorClock) <= len(vectorClock) {
			var lenDiff int = len(vectorClock) - len(msgVectorClock)
			for j := 0; j < lenDiff; j++ {
				msgVectorClock = append(msgVectorClock, 0)
			}
		}
		// Compare and update vectorclock values
		if vectorClock[i] < msgVectorClock[i] {
			vectorClock[i] = msgVectorClock[i]
		}
	}
	vectorClock[0]++
}

// Get preferred outbound ip of this machine
// Usefull if you have to know which ip you should dial, in a client running on an other computer
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// sets the logger to use a log.txt file instead of the console
func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	if err := os.Truncate("log_server.txt", 0); err != nil {
		fmt.Printf("Failed to truncate: %v \n", err)
		log.Printf("Failed to truncate: %v", err)
	}

	// This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log_server.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
