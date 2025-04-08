package main

import (
	"context"
	"log"
	"time"

	pb "mubashir-crud/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loggingInterceptor(
	ctx context.Context,
	method string, // The full method name (/package.Service/Method)
	req interface{},
	reply interface{},
	cc *grpc.ClientConn, // The ClientConn associated with this call
	invoker grpc.UnaryInvoker, // The actual RPC invoker
	opts ...grpc.CallOption,
) error {
	// Log request details
	log.Printf("Client: Sending request: %+v, Method: %s", req, method)

	// Call the invoker to continue the RPC execution
	err := invoker(ctx, method, req, reply, cc, opts...)

	// Log the response or error
	if err != nil {
		log.Printf("Client: Error occurred while processing request: %v", err)
	} else {
		log.Printf("Client: Successfully processed request: %s", method)
	}

	return err
}

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(loggingInterceptor), // Add the interceptor
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	log.Println("Client created: ", 5*time.Second)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// Create context with token metadata for secured methods
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "abcd12356##"))

	// Create
	user, err := client.CreateUser(ctx, &pb.User{Name: "Mubashir", Email: "mubashir@example.com"})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("User created: %s", user.Id)

	// Create 2
	user, err = client.CreateUser(ctx, &pb.User{Name: "Ali", Email: "ali@example.com"})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("User created: %s", user.Id)

	// Get
	gotUser, err := client.GetUser(ctx, &pb.UserId{Id: user.Id})
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	log.Printf("User retrieved: %s", gotUser)

	// Udpate
	user.Name = "Mubashir Updated"
	updatedUser, err := client.UpdateUser(ctx, user)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	log.Printf("User updated: %s", updatedUser)

	// List
	users, err := client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	log.Printf("Users listed: %v", users.Users)

	// Delete
	_, err = client.DeleteUser(ctx, &pb.UserId{Id: user.Id})
	if err != nil {
		log.Fatalf("DeleteUser error: %v", err)
	}
	log.Printf("User deleted: %s", user.Id)
}
