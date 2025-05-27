package handler

import (
	"auth-service/pb"
	"auth-service/service"
	"context"
	"log"
)

// AuthHandler implements AuthService gRPC server
type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register request: %s", req.Email)

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Username, email and password are required",
		}, nil
	}

	// Register user
	user, err := service.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Register error: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  user.ID.Hex(),
	}, nil
}

// Login handles user authentication
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Login request: %s", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" {
		return &pb.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	// Login user
	token, user, err := service.LoginUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Success:     true,
		Message:     "Login successful",
		AccessToken: token,
		UserId:      user.ID.Hex(),
	}, nil
}

// Logout handles user logout (blacklist token)
func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	log.Printf("Logout request")

	// Validate input
	if req.AccessToken == "" {
		return &pb.LogoutResponse{
			Success: false,
			Message: "Access token is required",
		}, nil
	}

	// Logout user
	err := service.LogoutUser(req.AccessToken)
	if err != nil {
		log.Printf("Logout error: %v", err)
		return &pb.LogoutResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}
