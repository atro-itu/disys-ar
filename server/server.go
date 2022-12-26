package main

import (
	"context"
	"flag"
	"log"
	"net"

	gRPC "github.com/duckth/disys-ar/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedPingerServer
	port string
}

var port = flag.String("port", "5000", "Server port")

func main() {
	flag.Parse()
	listen, _ := net.Listen("tcp", "localhost:"+*port)
	grpcServer := grpc.NewServer()
	pinger := &Server{
		port: *port,
	}

	gRPC.RegisterPingerServer(grpcServer, pinger)
	log.Printf("Listening on port %s...", *port)
	grpcServer.Serve(listen)
}

func (s *Server) Ping(ctx context.Context, req *gRPC.PingRequest) (*gRPC.PongResponse, error) {
	return &gRPC.PongResponse{Message: "Pong!"}, nil
}
