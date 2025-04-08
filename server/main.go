package main

import (
	"context"
	"log"
	"net"

	pb "mubashir-crud/proto/userpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func validateToken(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	// check if the token is present in the metadata
	tokens := md["authorization"]
	if len(tokens) == 0 {
		return false
	}

	return tokens[0] == "abcd12356##"
}

// Middleware: Logging interceptor for unary RPCs
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Log request details
	log.Printf("Received request: %+v, Method: %s", req, info.FullMethod) // info.FullMethod is a string, so no parentheses needed

	// Call the handler to continue the RPC execution
	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error occurred while processing request: %v", err)
	} else {
		log.Printf("Successfully processed request: %s", info.FullMethod) // info.FullMethod is a string
	}

	return resp, err
}

type server struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
}

func (s *server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {

	// Token Validation
	if !validateToken(ctx) {
		return nil, grpc.Errorf(grpc.Code(grpc.ErrClientConnClosing), "Invalid token")
	}

	user.Id = uuid.NewString()
	s.users[user.Id] = user
	return user, nil
}

func (s *server) GetUser(ctx context.Context, id *pb.UserId) (*pb.User, error) {
	user, exists := s.users[id.Id]
	if !exists {
		return nil, grpc.Errorf(grpc.Code(grpc.ErrClientConnClosing), "User not found")
	}
	return user, nil
}

func (s *server) UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error) {

	// Token Validation
	if !validateToken(ctx) {
		return nil, grpc.Errorf(grpc.Code(grpc.ErrClientConnClosing), "Invalid token")
	}

	_, exits := s.users[user.Id]
	if !exits {
		return nil, grpc.Errorf(grpc.Code(grpc.ErrClientConnClosing), "User not found")
	}
	s.users[user.Id] = user
	return user, nil
}

func (s *server) DeleteUser(ctx context.Context, id *pb.UserId) (*pb.Empty, error) {
	delete(s.users, id.Id)
	return &pb.Empty{}, nil
}

func (s *server) ListUsers(ctx context.Context, empty *pb.Empty) (*pb.UserList, error) {
	users := []*pb.User{}
	for _, user := range s.users {
		users = append(users, user)
	}
	return &pb.UserList{Users: users}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// s := grpc.NewServer()

	// Create a new gRPC server with middleware (logging interceptor)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor), // Add the interceptor
	)

	pb.RegisterUserServiceServer(s, &server{users: make(map[string]*pb.User)})

	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
