package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

// ConnectMongoDB connects to MongoDB
func ConnectMongoDB() error {
	// MongoDB connection string
	uri := "mongodb://localhost:27019"

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	Client = client
	Database = client.Database("authdb") // database name

	log.Println("Connected to MongoDB successfully!")
	return nil
}

// GetUsersCollection returns users collection
func GetUsersCollection() *mongo.Collection {
	return Database.Collection("users")
}

// GetEndTokensCollection returns end_tokens collection
func GetEndTokensCollection() *mongo.Collection {
	return Database.Collection("end_tokens")
}
