package main

import (
	"log"
	"net"
	"reset-password-service/config"
	"reset-password-service/db"
	"reset-password-service/handler"
	"reset-password-service/pb"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal("DB Connection error:", err)
	}

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	server := grpc.NewServer()
	pb.RegisterResetPasswordServiceServer(server, &handler.ResetPasswordService{JWTSecret: cfg.JWTSecret})

	log.Println("ResetPasswordService running on", cfg.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
