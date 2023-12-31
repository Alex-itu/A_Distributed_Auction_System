package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	// inspired by https://github.com/PatrickMatthiesen/DSYS-gRPC-template and https://articles.wesionary.team/grpc-console-chat-application-in-go-dd77a29bb5c3
	gRPC "github.com/Alex-itu/A_Distributed_Auction_System/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "Bames Nond", "Senders name")
var serverPorts = flag.String("serverPorts", ":8080 :8081 :8082", "TcP SeRvEr pOrTs UwU")
var clientId = flag.Int("id", 0, "Client id")

var ServerConn1 *grpc.ClientConn //the server connection
var ServerConn2 *grpc.ClientConn
var ServerConn3 *grpc.ClientConn
var auctionServer1 gRPC.AuctionServiceClient // new chat server client
var auctionServer2 gRPC.AuctionServiceClient
var auctionServer3 gRPC.AuctionServiceClient

var servers []string

var clientID int32 // clientID is set to 1 by default

func main() {
	//parse flag/arguments
	flag.Parse()
	servers = strings.Split(*serverPorts, " ")
	clientID = int32(*clientId)
	
	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	f := setLog()
	defer f.Close()

	//connect to server and close the connection when program closes
	fmt.Println("--- join Server ---")
	ConnectToServers()
	
	defer ServerConn1.Close()
	defer ServerConn2.Close()
	defer ServerConn3.Close()	

	//start the biding
	parseInput()
}

// connect to server
func ConnectToServers() {

	//dial options
	//the server is not using TLS, so we use insecure credentials
	//(should be fine for local testing but not in the real world)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	//dial the server, with the flag "server", to get a connection to it
	conn1, err := grpc.Dial(servers[0], opts...)
	if err != nil {
		fmt.Printf("Fail to Dial : %v \n", err)
		log.Printf("Fail to Dial : %v", err)
		return
	}
	auctionServer1 = gRPC.NewAuctionServiceClient(conn1)
	ServerConn1 = conn1
	fmt.Println("Connected to server 0")

	conn2, err := grpc.Dial(servers[1], opts...)
	if err != nil {
		fmt.Printf("Fail to Dial : %v \n", err)
		log.Printf("Fail to Dial : %v", err)
		return
	}
	auctionServer2 = gRPC.NewAuctionServiceClient(conn2)
	ServerConn2 = conn2
	fmt.Println("Connected to server 1")

	conn3, err := grpc.Dial(servers[2], opts...)
	if err != nil {
		fmt.Printf("Fail to Dial : %v \n", err)
		log.Printf("Fail to Dial : %v", err)
		return
	}
	auctionServer3 = gRPC.NewAuctionServiceClient(conn3)
	ServerConn3 = conn3
	fmt.Println("Connected to server 2")

	//for the chat implementation
	// makes a client from the server connection and saves the connection
	// and prints rather or not the connection was is READY
}

// watch the god
func parseInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the auction!")
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
		splitInput := strings.Split(input, " ")

		if splitInput[0] == "exit" { 
			time.Sleep(1 * time.Second)
			os.Exit(1)
		} else if splitInput[0] == "bid" {
			amount64, err := strconv.ParseFloat(splitInput[1], 32)
			amount32 := float32(amount64)
			if err != nil {
				fmt.Printf("%v \n", err)
				log.Fatalf("%v", err)
			}
			fmt.Println(amount32)
			ack1, err := auctionServer1.Bid(context.Background(), &gRPC.BidAmount{Amount: amount32, ClientID: clientID, ClientName: *clientsName})
			if err != nil {
				fmt.Printf("Server 0 is down \n")
				log.Printf("Server 0 is down")
			} else {
				fmt.Println(ack1.Message)
				log.Println(ack1.Message)
			}
			
			ack2, err := auctionServer2.Bid(context.Background(), &gRPC.BidAmount{Amount: amount32, ClientID: clientID, ClientName: *clientsName})
			if err != nil {
				fmt.Printf("Server 1 is down \n")
				log.Printf("Server 1 is down")
			} else {
				fmt.Println(ack2.Message)
				log.Println(ack2.Message)
			}
			
			ack3, err := auctionServer3.Bid(context.Background(), &gRPC.BidAmount{Amount: amount32, ClientID: clientID, ClientName: *clientsName})
			if err != nil {
				fmt.Printf("Server 2 is down \n")
				log.Printf("Server 2 is down")
			} else {
				fmt.Println(ack3.Message)
				log.Println(ack3.Message)
			}
			

		} else if splitInput[0] == "result" {
			result, err := auctionServer1.Result(context.Background(), &gRPC.Void{})
			if err != nil {
				fmt.Printf("Server 0 is down. Trying on connection 1 \n")
				log.Printf("Server 0 is down. Trying on connection 1")
				result, err = auctionServer2.Result(context.Background(), &gRPC.Void{})
				if err != nil {
					fmt.Printf("Server 1 is down. Trying on connection 2 \n")
					log.Printf("Server 1 is down. Trying on connection 2")
					result, err = auctionServer3.Result(context.Background(), &gRPC.Void{})
					if err != nil {
						fmt.Printf("you are offcially fucked. All servers are dead \n")
						log.Printf("you are offcially fucked. All servers are dead")
					}
				}
			}
			
			
			if result.BidDone {
				fmt.Printf("The bid is over and the winner is: %s \nWith a bid of: %f \n", result.ClientName, result.Amount)
				log.Printf("The bid is over and the winner is: %s \nWith a bid of: %f", result.ClientName, result.Amount)
			} else {
				fmt.Printf("The current highest bid is: %s \nWith a bid of: %f \n", result.ClientName, result.Amount)
				log.Printf("The current highest bid is: %s \nWith a bid of: %f", result.ClientName, result.Amount)
			}
		}
	}
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

