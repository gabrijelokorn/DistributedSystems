package main

import (
	"api/grpc/protobufStorage"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main () {

	conn, err := grpc.Dial("localhost:8100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	lecturesCreate := protobufStorage.Todo{Task: "test", Completed: false, Commited: false}

	if _, err := grpcClient.Put(contextPGC, &lecturesCreate); err != nil {
		panic(err)
	}
}