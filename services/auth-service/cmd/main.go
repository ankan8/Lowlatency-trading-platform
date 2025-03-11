package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"

    pb "github.com/ankan8/swapsync/backend/services/auth-service/proto"
    "github.com/ankan8/swapsync/backend/services/auth-service/service"
    "github.com/ankan8/swapsync/backend/internal/config" 
	"github.com/ankan8/swapsync/backend/internal/middleware"
)



type server struct {
	pb.UnimplementedAuthServiceServer
}


func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Printf("Received Login request for email: %s\n", req.GetEmail())
	token, err := service.AuthenticateUser(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Printf("Received Register request for email: %s\n", req.GetEmail())
	userID, err := service.RegisterUser(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{UserId: userID}, nil
}

func main() {
	config.ConnectDB()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryJWTInterceptor),)
	pb.RegisterAuthServiceServer(grpcServer, &server{})

	log.Printf("Auth Service gRPC server is listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
