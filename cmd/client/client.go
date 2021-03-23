package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eduardodataa/fullCycle-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Não conseguiu se conectar no gRPC Server", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)

	// AddUserVerbose(client)

	// AddUsers(client)

	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:       "0",
		Name:     "Joao",
		Email:    "j@j.com",
		Telefone: "12456",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Não conseguiu fazer requisição gRPC", err)
	}

	fmt.Println(res)

}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:       "0",
		Name:     "Joao",
		Email:    "j@j.com",
		Telefone: "12456",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Não conseguiu fazer requisição gRPC %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Não conseguiu receber a mensagem %v", err)
		}
		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Edu1",
			Email: "edu1@edu.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Edu2",
			Email: "edu2@edu.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Edu3",
			Email: "edu3@edu.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Edu4",
			Email: "edu4@edu.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Erro ao receber resposta: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Edu1",
			Email: "edu1@edu.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Edu2",
			Email: "edu2@edu.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Edu3",
			Email: "edu3@edu.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Edu4",
			Email: "edu4@edu.com",
		},
	}

	wait := make(chan int)

	go func(){
		for _, req := range reqs {
			fmt.Println("Enviando usuário: ",req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func(){
		for  {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil{
				log.Fatalf("Erro receber dados: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status: %v", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

		<-wait

}