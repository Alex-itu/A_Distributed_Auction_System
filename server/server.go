package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	// inspired by https://github.com/PatrickMatthiesen/DSYS-gRPC-template and https://articles.wesionary.team/grpc-console-chat-application-in-go-dd77a29bb5c3
	Auction "github.com/Alex-itu/A_Distributed_Auction_System/proto"

	"google.golang.org/grpc"
)

// Run server with:
// go run server/server.go -port 8080 -id 0

type RMserver struct {
	Auction.UnimplementedAuctionServiceServer        //need this if it's a server
	port                                      string // Not required but useful if your server needs to know what port it's listening to
	Id                                        int

	mutex sync.Mutex // used to lock the server to avoid race conditions.

}

var clientID = 0

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var port = flag.String("port", "8080", "Server port") // set with "-port <port>" in terminal
var serverId = flag.Int("id", 0, "Server id")
var server *RMserver

// Maps
var clientNames = make(map[int]string)
var CurrentBids = make(map[int32]float32)
var BackupAckRecieved = make(map[int]bool)

func main() {

	f := setLog() //uncomment this line to log to a log.txt file instead of the console
	defer f.Close()

	// This parses the flags and sets the correct/given corresponding values.
	flag.Parse()
	fmt.Println(".:server is starting:.")

	// launch the server
	launchServer()

	// code here is unreachable because launchServer occupies the current thread.
	for {

	}
}

func launchServer() {
	fmt.Printf("Server %d: Attempts to create listener on port %s\n", *serverId, *port)
	log.Printf("Server %d: Attempts to create listener on port %s\n", *serverId, *port)

	// Create listener for the RMserver connection
	listOnServerClient, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		fmt.Printf("Server %d: Failed to listen on port %s: %v \n", *serverId, *port, err)
		log.Printf("Server %d: Failed to listen on port %s: %v", *serverId, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// Create listener for all client that wants to bid
	// listToAllClients, err := net.Listen("tcp", "localhost:"+*port)
	// if err != nil {
	// 	fmt.Printf("Server %s: Failed to listen on port %s: %v \n", *serverName, *port, err)
	// 	log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
	// 	return
	// }

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// makes a new server instance using the name and port from the flags.
	server = &RMserver{
		port: *port,
		Id:   *serverId,
	}

	Auction.RegisterAuctionServiceServer(grpcServer, server) //Registers the server to the gRPC server.

	fmt.Printf("Server %d: Listening at %v \n", *serverId, listOnServerClient.Addr())
	log.Printf("Server %d: Listening at %v \n", *serverId, listOnServerClient.Addr())

	if err := grpcServer.Serve(listOnServerClient); err != nil {
		fmt.Printf("failed to serve %v", err)
		log.Fatalf("failed to serve %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}



func (s *RMserver) Bid(cxt context.Context, msg *Auction.BidAmount) (*Auction.Ack, error) {
	if msg.GetAmount() > HighestBid() { //This check may need to be for the current highest bid for that auction idk
		CurrentBids[int32(clientID)] = msg.Amount
		return &Auction.Ack{Message: "Nice job team", ClientID: msg.ClientID}, nil
	} else {
		return &Auction.Ack{Message: "Bid is lower than current highest bid: " + fmt.Sprint(HighestBid()), ClientID: msg.ClientID}, nil
	} 
}

func HighestBid() float32 {
	max := float32(-1.0)
	for _, CurrentBid := range CurrentBids {
		if CurrentBid > max {
			max = CurrentBid
		}
	}
	return max
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
