package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Predefined categories (optional, for validation)
var AllowedCategories = map[string]bool{
	"Health":        true,
	"Career":        true,
	"Education":     true,
	"Personal":      true,
	"Finance":       true,
	"Hobby":         true,
	"Relationships": true,
}

// Goal represents a user's goal.
type Goal struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"` // New Field
	Steps       []string           `bson:"steps" json:"steps"`
	Progress    map[string]bool    `bson:"progress" json:"progress"`
	Status      string             `bson:"status" json:"status"`
	DueDate     time.Time          `bson:"due_date,omitempty" json:"due_date,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
