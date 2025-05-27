package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	MongoURI  string
	Port      string
	JWTSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27019"),
		Port:      getEnv("PORT", ":50051"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
	}

	return config
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
