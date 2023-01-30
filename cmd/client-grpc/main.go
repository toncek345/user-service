package main

import (
	"context"
	"log"

	pb "github.com/toncek345/userservice/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUsersClient(conn)
	client.AddUser(context.Background(), &pb.AddUserMessage{})
}
