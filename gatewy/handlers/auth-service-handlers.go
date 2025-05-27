package handlers

import (
	"GridWhiz/proto/auth-service/pb"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// User model (เหมือนกับ models.User แต่อาจไม่ต้องใช้ primitive.ObjectID ที่ Gateway)
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// EndToken model (ใช้สำหรับ logout)
type EndToken struct {
	Token string `json:"access_token"`
}

// Assume authClient is gRPC client injected somehow
var authClient pb.AuthServiceClient

func InitAuthClient(c pb.AuthServiceClient) {
	authClient = c
}

// Register Handler
func RegisterHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Register(ctx, &pb.RegisterRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": resp.Success,
		"message": resp.Message,
		"user_id": resp.UserId,
	})
}

// Login Handler
func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      resp.Success,
		"message":      resp.Message,
		"access_token": resp.AccessToken,
		"user_id":      resp.UserId,
	})
}

// Logout Handler
func LogoutHandler(c *gin.Context) {
	var token EndToken
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Logout(ctx, &pb.LogoutRequest{
		AccessToken: token.Token,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": resp.Success,
		"message": resp.Message,
	})
}
