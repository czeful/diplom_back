package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/czeful/diplom_back/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB initializes a MongoDB connection
func ConnectDB(cfg *config.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)
	log.Println("Connected to MongoDB")
	return db, nil
}
