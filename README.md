# A_Distributed_Auction_System

# How to run
Server:

\- Boot up three terminal windows dedicated to being your Servers

\- Run these to boot up your servers:

go run server/server.go -port 8080 -id 0 -endTime {some_timestamp_in_the_future_given_as_HH:MM:SS} (e.g. 15:39:20)

go run server/server.go -port 8081 -id 1 -endTime {some_timestamp_in_the_future_given_as_HH:MM:SS} (e.g. 15:39:20)

go run server/server.go -port 8082 -id 2 -endTime {some_timestamp_in_the_future_given_as_HH:MM:SS} (e.g. 15:39:20)

# How to Run Client
\- Boot up a terminal window for you client

\- To boot up a client you can use this:

go run client/client.go -name Bames Nond -serverPorts :8080 :8081 :8082 -id 0 

# Some notes about the different paramaters for Server
The servers are hardcoded to only have 3 processes so please dont try to do it with more

\- The port is the port given to the server. These cant be changed but will have to be changed on the client side as well. Default value is 8080

\- The id is the id the server is known by. Default value is 0

\- The endTime is the value that sets when the auction ends. This time is given in the format HH:MM:SS: Default value is 00:00:00

# Some notes about the different paramters for Clients

\- The name is the name of the client that gets printed on the result call. Default value is Bames Nond

\- The server ports of the 3 servers. This is just given as a string seperated by spaces and the ports must contain a ":". Default value is :8080 :8081 :8082

\- The id is just the client id. This value has to be different from other clients otherwise it will add the bid to the same client. Default value is 0
