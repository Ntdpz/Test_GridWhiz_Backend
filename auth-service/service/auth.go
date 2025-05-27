package service

import (
	"auth-service/db"
	"auth-service/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUser creates a new user
func RegisterUser(email, password string) (*models.User, error) {
	ctx := context.Background()
	collection := db.GetUsersCollection()

	// Check if user already exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create new user
	user := models.User{
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	// Insert user to database
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Set user ID
	user.ID = result.InsertedID.(primitive.ObjectID)

	return &user, nil
}

// LoginUser authenticates user and returns JWT token
func LoginUser(email, password string) (string, *models.User, error) {
	ctx := context.Background()
	collection := db.GetUsersCollection()

	// Find user by email
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil, fmt.Errorf("invalid email or password")
		}
		return "", nil, fmt.Errorf("database error: %v", err)
	}

	// Check password
	if !CheckPassword(password, user.Password) {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return token, &user, nil
}

// LogoutUser adds token to blacklist
func LogoutUser(tokenString string) error {
	// Validate token first
	_, err := ValidateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	// Add token to blacklist
	err = BlacklistToken(tokenString)
	if err != nil {
		return fmt.Errorf("failed to logout: %v", err)
	}

	return nil
}
