package main

import (
	"api/grpc/protobufStorage"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitiatorSend (send_to_URL string) {

	time.Sleep(time.Second)
	fmt.Printf("send_to_URL: %v\n", send_to_URL)

	conn, err := grpc.Dial(send_to_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	lecturesCreate := protobufStorage.Todo{Task: "predavanja", Completed: false, Commited: false}
	
	fmt.Print("1. Create: ")
	if _, err := grpcClient.Put(contextPGC, &lecturesCreate); err != nil {
		panic(err)
	}
	fmt.Println("done")
	
}

func InitiatorGet (send_to_URL string) {

	time.Sleep(time.Second)
	fmt.Printf("send_to_URL: %v\n", send_to_URL)

	conn, err := grpc.Dial(send_to_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	lecturesCreate := protobufStorage.Todo{Task: "predavanja", Completed: false, Commited: false}

	fmt.Print("Read 1: ")
	if response, err := grpcClient.Get(contextPGC, &lecturesCreate); err == nil {
		fmt.Println(response.Todos, ": done")
	} else {
		panic(err)
	}
}

func Server (read_from_URL string) {
	fmt.Printf("readURL: %v\n", read_from_URL)
	// pripravimo strežnik gRPC
	grpcServer := grpc.NewServer()

	// pripravimo strukuro za streženje metod CRUD na shrambi TodoStorage
	pgcServer := NewServerPGC()

	// streženje metod CRUD na shrambi TodoStorage povežemo s strežnikom gRPC
	protobufStorage.RegisterPGCServer(grpcServer, pgcServer)

	// izpišemo ime strežnika
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// odpremo vtičnico
	listener, err := net.Listen("tcp", read_from_URL)
	if err != nil {
		panic(err)
	}
	fmt.Printf("gRPC server listening at %v%v\n", hostName, read_from_URL)
	// začnemo s streženjem
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func main () {

	go Server("localhost:8099")

	InitiatorSend("localhost:8100")
	InitiatorGet("localhost:8109")

	time.Sleep(20 * time.Second)
}