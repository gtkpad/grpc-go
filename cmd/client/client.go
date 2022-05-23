package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gtkpad/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer connection.Close()
	
	client := pb.NewUserServiceClient(connection)

	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Name: "John Doe",
		Email: "john@email.com",
		Id: "0",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("failed to run gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Name: "John Doe",
		Email: "john@email.com",
		Id: "0",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("failed to run gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("failed to receive stream: %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id: "1",
			Name: "John Doe 1",
			Email: "john1@email.com",
		},
		&pb.User{
			Id: "2",
			Name: "John Doe 2",
			Email: "john2@email.com",
		},
		&pb.User{
			Id: "3",
			Name: "John Doe 3",
			Email: "john3@email.com",
		},
		&pb.User{
			Id: "4",
			Name: "John Doe 4",
			Email: "john4@email.com",
		},
		&pb.User{
			Id: "5",
			Name: "John Doe 5",
			Email: "john5@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("failed to run gRPC request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()


	if err != nil {
		log.Fatalf("failed to receive stream: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
		stream, err := client.AddUserStreamBoth(context.Background())

		if err != nil {
			log.Fatalf("failed to run gRPC request: %v", err)
		}

		reqs := []*pb.User{
			&pb.User{
				Id: "1",
				Name: "John Doe 1",
				Email: "john1@email.com",
			},
			&pb.User{
				Id: "2",
				Name: "John Doe 2",
				Email: "john2@email.com",
			},
			&pb.User{
				Id: "3",
				Name: "John Doe 3",
				Email: "john3@email.com",
			},
			&pb.User{
				Id: "4",
				Name: "John Doe 4",
				Email: "john4@email.com",
			},
			&pb.User{
				Id: "5",
				Name: "John Doe 5",
				Email: "john5@email.com",
			},
		}

		wait := make(chan int)

		go func() {
			for _, req := range reqs {
				fmt.Println("Sending request: ", req)
				stream.Send(req)
				time.Sleep(time.Second * 2)
			}
			stream.CloseSend()
		}()

		go func() {
			for {
				res, err := stream.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatalf("failed to receive stream: %v", err)
				}

				fmt.Println("User %v Status: %v", res.GetUser().GetName(), res.GetStatus())
			}

			close(wait)
		}()

		<-wait
}