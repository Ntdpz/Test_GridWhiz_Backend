package main

import (
	"auth-service/config"
	"auth-service/db"
	"auth-service/handler"
	"auth-service/pb/auth-service/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	err := db.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register auth service
	authHandler := handler.NewAuthHandler()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Start listening
	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("gRPC Auth Service started on port %s", cfg.Port)
	log.Printf("MongoDB connected to: %s", cfg.MongoURI)

	// Start server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
