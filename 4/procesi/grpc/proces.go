package main

import (
	"fmt"
	"time"
	"api/grpc/protobufStorage"
	"api/storage"
	"context"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/protobuf/types/known/emptypb"
)

var listenURL string
var getPutURL string
var commitURL string

var pgcClientGetPut protobufStorage.PGCClient
var pgcClientCommit protobufStorage.PGCClient
var pgcSelfCommit protobufStorage.PGCClient

func Proces(ListenURL string, GetPutURL string, CommitURL string) {

	fmt.Println("ListenURL: ", ListenURL)
	fmt.Println("GetPutURL: ", GetPutURL)
	fmt.Println("CommitURL: ", CommitURL)

	fmt.Println()

	listenURL = ListenURL
	getPutURL = GetPutURL
	commitURL = CommitURL
	
	go Server(ListenURL)
	go ClientGetPut(GetPutURL)
	go ClientCommit(CommitURL)
	
	if getPutURL == "rep" {
		go SelfCommit(ListenURL)
	}
		
	// Čakamo, da se strežnik in odjemalca pripravijo
	time.Sleep(1 * time.Second)

	// Pripravimo odjemalca
	if listenURL == ":8100" {

		connGlava, err := grpc.Dial("localhost:8100", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		defer connGlava.Close()
		connRep, err := grpc.Dial("localhost:8109", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer connRep.Close()

		contextPGC, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		grpcClientGlava := protobufStorage.NewPGCClient(connGlava)
		grpcClientRep := protobufStorage.NewPGCClient(connRep)

		lecturesCreate := protobufStorage.Todo{Task: "test", Completed: false, Commited: false}
		
		if val, err := grpcClientGlava.Put(contextPGC, &lecturesCreate); err != nil {
			panic(err)
		} else {
			fmt.Println("Result: ", val)
		}

		if val, err := grpcClientRep.Get(contextPGC, &lecturesCreate); err != nil {
			panic(err)
		} else {
			fmt.Println("Result: ", val)
		}

	}
	time.Sleep(10 * time.Second)
}

/**
 * 	Strežnik
 *  Pripravimo strežnik, ki bo poslušal na vratih ListenPORT
 *  Strežnik bo uporabljal funkcije put, get, commit
 */
// stuktura za strežnik CRUD za shrambo TodoStorage

func (s* serverPGC) Put(ctx context.Context, in *protobufStorage.Todo) (*protobufStorage.StatusResponse, error) {

	var ret struct{}
	err := s.todoStore.Put(&storage.Todo{Task: in.Task, Completed: in.Completed, Commited: false}, &ret)

	val := &protobufStorage.StatusResponse{}
	
	if err == nil {
		fmt.Println("Added: [", in.Task, in.Completed, in.Commited, "]")

		contextPGC, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Če je proces rep naj pošlje commit
		if getPutURL == "rep" {
			val, err = pgcSelfCommit.Commit(contextPGC, in)
		// Če proces ni rep naj pošlje naprej
		} else {
			val, err = pgcClientGetPut.Put(contextPGC, in)
		}

	} else {
		val.Value = false
		fmt.Println("Error: Failed to add: [", in.Task, "]")
	}

	return val, err
}

// Funkcija, ki vrne vse podatke iz shrambe - iz repa
func (s* serverPGC) Get(ctx context.Context, in *protobufStorage.Todo) (*protobufStorage.TodoStorage, error) {

	// Naredimo mapo, v katero bomo prebrali vse podatke
	dict := make(map[string]storage.Todo)
	// Preberemo podatke
	err := s.todoStore.Get(&storage.Todo{Task: in.Task, Commited: true}, &dict)
	// Pripravimo mapo, ki jo bomo vrnili
	pbDict := protobufStorage.TodoStorage{}

	
	for k, v := range dict {
		pbDict.Todos = append(pbDict.Todos, &protobufStorage.Todo{Task: k, Completed: v.Completed, Commited: v.Commited})
	}
	
	// Vrnemo podatke
	return &pbDict, err
}

func (s* serverPGC) Commit(ctx context.Context, in *protobufStorage.Todo) (*protobufStorage.StatusResponse, error) {
	var ret struct{}
	err := s.todoStore.Commit(&storage.Todo{Task: in.Task, Completed: in.Completed, Commited: true}, &ret)

	val := &protobufStorage.StatusResponse{}

	if err == nil {

		// Če je glava naj ne pošlje naprej ampak naj vrne vrednost
		if commitURL == "glava" {
			val.Value = true
		} else {
			contextPGC, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
			defer cancel()

			val, err = pgcClientCommit.Commit(contextPGC, in)

			// Če je kateri od procesov vrnil false naj tudi ta nazaj vrne false
			// in nase podatke ne shrani - torej ne commita
			// if val.Value == false {
			// 	for {
			// 		err := s.todoStore.Commit(&storage.Todo{Task: in.Task, Completed: in.Completed, Commited: false}, &ret)
			// 		if err == nil {
			// 			break
			// 		}
			// 	}
			// }

		}
	} else {
		val.Value = false
		fmt.Println("Error: Failed to add: [", in.Task, "]")
	}

	return val, err
}

type serverPGC struct {
	protobufStorage.UnimplementedPGCServer
	todoStore storage.TodoStorage
}

// pripravimo nov strežnik CRUD za shrambo TodoStorage
func NewServerPGC() *serverPGC {
	todoStorePtr := storage.NewTodoStorage()
	return &serverPGC{protobufStorage.UnimplementedPGCServer{}, *todoStorePtr}
}

func Server (ListenURL string) {
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
	listener, err := net.Listen("tcp", ListenURL)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server listening at %v%v\n", hostName, ListenURL)
	// začnemo s streženjem
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func ClientGetPut (GetPutURL string) {

	conn, err := grpc.Dial(GetPutURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	pgcClientGetPut = protobufStorage.NewPGCClient(conn)
}

func ClientCommit (CommitURL string) {
	
	conn, err := grpc.Dial(CommitURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	
	pgcClientCommit = protobufStorage.NewPGCClient(conn)
}

func SelfCommit (SelfURL string) {
	
	conn, err := grpc.Dial(SelfURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	
	pgcSelfCommit = protobufStorage.NewPGCClient(conn)
}