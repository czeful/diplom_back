package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the Achievement Manager system.
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username"`
	Email          string             `bson:"email"`
	HashedPassword string             `json:"hashed_password"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}
