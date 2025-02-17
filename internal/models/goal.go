package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Goal represents a user's goal in the system.
type Goal struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Steps       []string           `bson:"steps"`      // Steps to complete the goal
	Status      string             `bson:"status"`     // e.g., "pending", "completed", etc.
	CreatedAt   time.Time          `bson:"created_at"` // Timestamp of goal creation
	UpdatedAt   time.Time          `bson:"updated_at"` // Timestamp for the last update
}
