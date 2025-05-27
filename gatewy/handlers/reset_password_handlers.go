package handlers

import (
	"GridWhiz/proto/reset-password-service/pb"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var resetPasswordClient pb.ResetPasswordServiceClient

func InitResetPasswordClient(c pb.ResetPasswordServiceClient) {
	resetPasswordClient = c
}

// RequestResetHandler รับ email เพื่อขอรีเซ็ตรหัสผ่าน
func RequestResetHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := resetPasswordClient.RequestReset(ctx, &pb.RequestResetRequest{
		Email: req.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

// VerifyResetTokenHandler ตรวจสอบว่า token ถูกต้องไหม
func VerifyResetTokenHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := resetPasswordClient.VerifyResetToken(ctx, &pb.VerifyResetTokenRequest{
		Token: req.Token,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": resp.Valid,
		"email": resp.Email,
	})
}

// ResetPasswordHandler รีเซ็ตรหัสผ่านโดยใช้ token และ new_password
func ResetPasswordHandler(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := resetPasswordClient.ResetPassword(ctx, &pb.ResetPasswordRequest{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}
