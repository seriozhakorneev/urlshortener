package main

import (
	"fmt"
	"log"
	"net"

	"github.com/urlshortener/shortener"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Server is up.")
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen on port 3000: %v", err)
	}

	s := shortener.Server{}

	grpcServer := grpc.NewServer()

	shortener.RegisterUrlServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 3000: %v", err)
	}

}
