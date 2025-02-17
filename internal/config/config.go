package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds application configuration
type Config struct {
	MongoURI    string
	Database    string
	Port        string
	JWTSecret   string
	TokenExpiry string
}

// LoadConfig reads from the .env file
func LoadConfig() *Config {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	return &Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		Database:    os.Getenv("DB_NAME"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		TokenExpiry: os.Getenv("TOKEN_EXPIRY"),
	}
}
