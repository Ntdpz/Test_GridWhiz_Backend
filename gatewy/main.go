package main

import (
	"GridWhiz/handlers"
	"GridWhiz/proto/auth-service/pb"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var mongoClient *mongo.Client
var authClient pb.AuthServiceClient

func connectMongoDB() error {
	// MongoDB connection string สำหรับ local MongoDB ที่ port 27019
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27019")

	// สร้าง context สำหรับ timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// เชื่อมต่อกับ MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// ทดสอบการเชื่อมต่อ
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	mongoClient = client
	log.Println("Connected to MongoDB successfully!")
	return nil
}

func disconnectMongoDB() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoClient.Disconnect(ctx)
		log.Println("Disconnected from MongoDB")
	}
}

func main() {
	PORT := ":8081"

	// สร้าง Gin router
	router := gin.Default()

	// Endpoint สำหรับทดสอบ API
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello GoGin",
			"status":  "API is running",
		})
	})

	// Endpoint สำหรับเชื่อมต่อ Database
	router.GET("/connect-database", func(c *gin.Context) {
		// ตรวจสอบว่าเชื่อมต่อแล้วหรือยัง
		if mongoClient != nil {
			// ทดสอบการเชื่อมต่อ
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := mongoClient.Ping(ctx, nil)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "Already connected to MongoDB",
					"status":  "success",
					"port":    "27019",
				})
				return
			}
		}

		// พยายามเชื่อมต่อใหม่
		err := connectMongoDB()
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to connect to MongoDB",
				"error":   err.Error(),
				"status":  "error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Connected to MongoDB successfully",
			"status":  "success",
			"port":    "27019",
		})
	})

	// Endpoint สำหรับตรวจสอบสถานะ Database
	router.GET("/database-status", func(c *gin.Context) {
		if mongoClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "Not connected to MongoDB",
				"status":  "disconnected",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := mongoClient.Ping(ctx, nil)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "MongoDB connection lost",
				"status":  "error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "MongoDB is connected and healthy",
			"status":  "connected",
			"port":    "27019",
		})
	})

	log.Printf("Starting server on port %s", PORT)

	// เชื่อมต่อกับ MongoDB เมื่อเริ่มต้นแอพ
	if err := connectMongoDB(); err != nil {
		log.Printf("Warning: Failed to connect to MongoDB on startup: %v", err)
	}

	// ปิดการเชื่อมต่อเมื่อแอพปิด
	defer disconnectMongoDB()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect grpc: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)
	handlers.InitAuthClient(authClient)

	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)
	router.POST("/logout", handlers.LogoutHandler)

	// เริ่มต้น server
	router.Run(PORT)
}
