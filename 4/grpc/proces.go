// Komunikacija po protokolu gRPC
// strežnik

package main

import (
	"api/grpc/protobufStorage"
	"api/storage"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func ClientPut (send_to_URL string, in *protobufStorage.Todo) {

	conn, err := grpc.Dial(send_to_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	lecturesCreate := protobufStorage.Todo{Task: in.Task, Completed: in.Completed, Commited: false}

	if _, err := grpcClient.Put(contextPGC, &lecturesCreate); err != nil {
		panic(err)
	}
}

func ClientCommit (send_back_URL string, in *protobufStorage.Todo, commit bool) {

	conn, err := grpc.Dial(send_back_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	lecturesCreate := protobufStorage.Todo{Task: in.Task, Completed: in.Completed, Commited: commit}

	if _, err := grpcClient.Commit(contextPGC, &lecturesCreate); err != nil {
		panic(err)
	}
}

func ClientGet (send_to_URL string, in *protobufStorage.Todo) {

	conn, err := grpc.Dial(send_to_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	lecturesCreate := protobufStorage.Todo{Task: in.Task, Completed: in.Completed, Commited: in.Commited}

	if response, err := grpcClient.Get(contextPGC, &lecturesCreate); err == nil {
		fmt.Println(response.Todos, ": done")
	} else {
		panic(err)
	}

}

func Server (read_from_URL string) {
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
	fmt.Printf("Server listening at %v%v\n", hostName, read_from_URL)
	// začnemo s streženjem
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}


func InitiatorSend (send_to_URL string, in *protobufStorage.Todo) {

	time.Sleep(time.Second)

	conn, err := grpc.Dial(send_to_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	if _, err := grpcClient.Put(contextPGC, in); err != nil {
		panic(err)
	}
}


type Mssg struct {
	TaskName string
	Completed bool
	Commited bool
}

func InitiatorGet (send_to_URL string, in *protobufStorage.Todo) (Mssg) {

	time.Sleep(time.Second)

	conn, err := grpc.Dial(send_to_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	contextPGC, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	grpcClient := protobufStorage.NewPGCClient(conn)

	if response, err := grpcClient.Get(contextPGC, in); err == nil {
		return Mssg {
			TaskName: response.Todos[0].Task,
			Completed: response.Todos[0].Completed,
			Commited: response.Todos[0].Commited,
		}
	} else {
		panic(err)
	}

	return Mssg {
		TaskName: "",
		Completed: false,
		Commited: false,
	}
}


var previousURL string
var currentURL string
var forwardURL string

func Initiator(glava string, rep string, read_from_URL string, send_to_URL string, send_back_URL string) {

	previousURL = send_back_URL
	currentURL = read_from_URL
	forwardURL = send_to_URL

	go Server (read_from_URL)


	task1 := protobufStorage.Todo{Task: "apple", Completed: false}
	task2 := protobufStorage.Todo{Task: "banana", Completed: false}
	task3 := protobufStorage.Todo{Task: "apple", Completed: true}
	task4 := protobufStorage.Todo{Task: "notes", Completed: true}
	task5 := protobufStorage.Todo{Task: "doctor", Completed: false}

	todoArray := []protobufStorage.Todo{task1, task2, task3, task4, task5}

	
	for _, v := range todoArray {
		
		for {
			InitiatorSend(glava, &v)
			
			getter := InitiatorGet(rep, &v)
			if getter.Commited == true {
				break
			}

			time.Sleep(3 * time.Second)
		}
	}


	time.Sleep(20 * time.Second)
}

func Proces(read_from_URL string, send_to_URL string, send_back_URL string) {

	previousURL = send_back_URL
	currentURL = read_from_URL
	forwardURL = send_to_URL

	go Server(read_from_URL)

	time.Sleep(20 * time.Second)
}

// stuktura za strežnik CRUD za shrambo TodoStorage
type serverPGC struct {
	protobufStorage.UnimplementedPGCServer
	todoStore storage.TodoStorage
}

// pripravimo nov strežnik CRUD za shrambo TodoStorage
func NewServerPGC() *serverPGC {
	todoStorePtr := storage.NewTodoStorage()
	return &serverPGC{protobufStorage.UnimplementedPGCServer{}, *todoStorePtr}
}

func (s* serverPGC) Put(ctx context.Context, in *protobufStorage.Todo) (*emptypb.Empty, error) {
	var ret struct{}
	err := s.todoStore.Put(&storage.Todo{Task: in.Task, Completed: in.Completed, Commited: false}, &ret)
	if err == nil {
		fmt.Println("Added: ", in.Task, in.Completed, in.Commited)
	}

	if forwardURL != "unknown" {
		go ClientPut(forwardURL, in)
	} else {
		if err == nil {
			ClientCommit(currentURL, in, true)
		} else {
			ClientCommit(currentURL, in, false)
		}
	}

	return &emptypb.Empty{}, err
}

func (s* serverPGC) Get(ctx context.Context, in *protobufStorage.Todo) (*protobufStorage.TodoStorage, error) {
	dict := make(map[string]storage.Todo)
	err := s.todoStore.Get(&storage.Todo{Task: in.Task, Completed: in.Completed, Commited: in.Commited}, &dict)

	pbDict := protobufStorage.TodoStorage{}

	for k, v := range dict {
		fmt.Println(currentURL)
		fmt.Println(dict)
		pbDict.Todos = append(pbDict.Todos, &protobufStorage.Todo{Task: k, Completed: v.Completed, Commited: v.Commited})
	}
	
	return &pbDict, err
}

func (s* serverPGC) Commit(ctx context.Context, in *protobufStorage.Todo) (*emptypb.Empty, error) {
	if previousURL != "unknown" {
		var ret struct{}
		err := s.todoStore.Commit(&storage.Todo{Task: in.Task, Completed: in.Completed, Commited: in.Commited}, &ret)
		
		go ClientCommit(previousURL, in, in.Commited)

		if err == nil {
			fmt.Println("Committed: ", in.Task, in.Completed, in.Commited)
		}

		return &emptypb.Empty{}, err
	} else {
		if in.Commited {
			fmt.Println(in.Task, " SUCCESSFULLY added to database")
		} else {
			fmt.Println(in.Task, " FAILED to add to database")
		}

		return &emptypb.Empty{}, nil
	}
}