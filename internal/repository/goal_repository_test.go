package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestCreateAndGetGoal tests creating a goal and retrieving it.
func TestCreateAndGetGoal(t *testing.T) {
	// Set up a context with timeout for the test.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to a test MongoDB instance (you can use a local instance or a dedicated test database).
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Use a dedicated test database
	db := client.Database("achievement_manager_test")

	// Clean up the goals collection before testing.
	err = db.Collection("goals").Drop(ctx)
	if err != nil {
		t.Fatalf("Failed to drop collection: %v", err)
	}

	// Initialize the repository.
	repo := NewGoalRepository(db)

	// Define a new goal.
	goal := &models.Goal{
		Name:        "Test Goal",
		Description: "Testing create and get functionality",
		Steps:       []string{"Step1", "Step2"},
		Status:      "pending",
	}

	// Test creating the goal.
	createdGoal, err := repo.CreateGoal(ctx, goal)
	assert.NoError(t, err)
	assert.NotNil(t, createdGoal)
	assert.NotEqual(t, primitive.NilObjectID, createdGoal.ID)

	// Test retrieving the created goal by its ID.
	fetchedGoal, err := repo.GetGoalByID(ctx, createdGoal.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdGoal.Name, fetchedGoal.Name)
	assert.Equal(t, createdGoal.Description, fetchedGoal.Description)
}
