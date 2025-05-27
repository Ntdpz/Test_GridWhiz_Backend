package main

import (
	"log"
	"net"
	"user-service/config"
	"user-service/db"
	"user-service/pb"
	"user-service/service"

	"google.golang.org/grpc"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to DB
	err := db.ConnectMongoDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	// Prepare gRPC server
	grpcServer := grpc.NewServer()

	// Register UserService
	userService := service.NewUserService(db.GetUsersCollection())
	pb.RegisterUserServiceServer(grpcServer, userService)

	// Start listening
	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("gRPC server started on %s", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
