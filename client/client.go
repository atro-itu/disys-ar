package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gRPC "github.com/duckth/disys-ar/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var servers []gRPC.PingerClient

func main() {
	ports := []int64{5000, 5001, 5002}

	for i := 0; i < len(ports); i++ {
		servers = append(servers, ConnectToServer(ports[i]))
	}

	readInput()
}

func readInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()
		if in == "exit" {
			break
		} else {
			Ping(servers)
		}
	}
}

func Ping(servers []gRPC.PingerClient) {
	hasReceivedPing := false
	for i := 0; i < len(servers); i++ {
		server := servers[i]

		response, _ := server.Ping(context.Background(), &gRPC.PingRequest{})

		if hasReceivedPing {
			continue
		}

		if response == nil {
			log.Printf("Received no response from server %d", i)
			continue
		}

		log.Printf(response.Message)
		hasReceivedPing = true
	}
}

func ConnectToServer(port int64) gRPC.PingerClient {
	var opts []grpc.DialOption
	var target = fmt.Sprintf("localhost:%d", port)

	opts = append(
		opts,
		grpc.WithBlock(),
		grpc.WithTimeout(1*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	fmt.Printf("Dialing on %s \n", target)
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("Fail to dial: %v\n", err)
	}

	return gRPC.NewPingerClient(conn)
}
