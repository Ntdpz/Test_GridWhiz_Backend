package handler

import (
	"context"
	"reset-password-service/db"
	"reset-password-service/models"
	"reset-password-service/pb"
	"reset-password-service/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordService struct {
	pb.UnimplementedResetPasswordServiceServer
	JWTSecret string
}

// RequestReset
func (s *ResetPasswordService) RequestReset(ctx context.Context, req *pb.RequestResetRequest) (*pb.RequestResetResponse, error) {
	userCol := db.GetUsersCollection()
	var user models.User
	err := userCol.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return &pb.RequestResetResponse{Message: "Email not found"}, nil
	}

	token, err := utils.GenerateResetToken(user.Email, s.JWTSecret, 15)
	if err != nil {
		return nil, err
	}

	_, err = db.GetResetPasswordTokensCollection().InsertOne(ctx, models.ResetToken{
		Email:     user.Email,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IsUsed:    false,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	// คุณสามารถส่ง email ได้ที่นี่ (mock หรือ integrate SMTP)
	return &pb.RequestResetResponse{Message: "Reset token generated"}, nil
}

// VerifyResetToken
func (s *ResetPasswordService) VerifyResetToken(ctx context.Context, req *pb.VerifyResetTokenRequest) (*pb.VerifyResetTokenResponse, error) {
	email, err := utils.ParseResetToken(req.Token, s.JWTSecret)
	if err != nil {
		return &pb.VerifyResetTokenResponse{Valid: false}, nil
	}

	var token models.ResetToken
	err = db.GetResetPasswordTokensCollection().FindOne(ctx, bson.M{"token": req.Token, "is_used": false}).Decode(&token)
	if err != nil || token.ExpiresAt.Before(time.Now()) {
		return &pb.VerifyResetTokenResponse{Valid: false}, nil
	}

	return &pb.VerifyResetTokenResponse{Valid: true, Email: email}, nil
}

// ResetPassword
func (s *ResetPasswordService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	email, err := utils.ParseResetToken(req.Token, s.JWTSecret)
	if err != nil {
		return nil, err
	}

	tokensCol := db.GetResetPasswordTokensCollection()
	usersCol := db.GetUsersCollection()

	// Check token again
	var token models.ResetToken
	err = tokensCol.FindOne(ctx, bson.M{"token": req.Token, "is_used": false}).Decode(&token)
	if err != nil || token.ExpiresAt.Before(time.Now()) {
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	_, err = usersCol.UpdateOne(ctx, bson.M{"email": email}, bson.M{
		"$set": bson.M{"password": string(hashedPassword)},
	})

	if err != nil {
		return nil, err
	}

	// Mark token used
	_, _ = tokensCol.UpdateOne(ctx, bson.M{"token": req.Token}, bson.M{"$set": bson.M{"is_used": true}})

	return &pb.ResetPasswordResponse{Message: "Password updated"}, nil
}
