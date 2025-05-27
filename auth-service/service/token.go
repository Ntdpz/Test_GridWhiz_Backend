package service

import (
	"auth-service/db"
	"auth-service/models"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

// JWT Secret (ในการใช้งานจริงควรเก็บใน ENV)
var jwtSecret = []byte("your-secret-key")

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token
func GenerateJWT(userID, email string) (string, error) {
	// Set expiration time (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	// Check if token is in end_tokens (blacklisted)
	if IsTokenBlacklisted(tokenString) {
		return nil, fmt.Errorf("token is blacklisted")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// BlacklistToken adds token to blacklist (end_tokens collection)
func BlacklistToken(tokenString string) error {
	ctx := context.Background()
	collection := db.GetEndTokensCollection()

	endToken := models.EndToken{
		Token:     tokenString,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(ctx, endToken)
	return err
}

// IsTokenBlacklisted checks if token is in blacklist
func IsTokenBlacklisted(tokenString string) bool {
	ctx := context.Background()
	collection := db.GetEndTokensCollection()

	var endToken models.EndToken
	err := collection.FindOne(ctx, bson.M{"token": tokenString}).Decode(&endToken)

	return err == nil // if no error, token exists in blacklist
}
