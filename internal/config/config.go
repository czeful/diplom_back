package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config struct holds application configuration
type Config struct {
	MongoURI    string
	Database    string
	Port        string
	JWTSecret   string
	TokenExpiry time.Duration
}

// LoadConfig reads from the .env file
func LoadConfig() *Config {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	expiryStr := os.Getenv("TOKEN_EXPIRY") // Get TOKEN_EXPIRY as string

	// Convert string to time.Duration
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		log.Printf("Invalid TOKEN_EXPIRY format, defaulting to 24h: %v", err)
		expiry = 24 * time.Hour // Default to 24 hours if parsing fails
	}

	return &Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		Database:    os.Getenv("DB_NAME"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		TokenExpiry: expiry,
	}
}
