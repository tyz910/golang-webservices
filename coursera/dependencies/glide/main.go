package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	fmt.Println("starting GRPC server at :8081")
	server.Serve(lis)
}
