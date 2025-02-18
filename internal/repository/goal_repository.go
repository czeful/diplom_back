package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GoalRepository struct handles database operations related to goals
type GoalRepository struct {
	collection *mongo.Collection
}

// NewGoalRepository creates a new instance of GoalRepository
func NewGoalRepository(db *mongo.Database) *GoalRepository {
	return &GoalRepository{
		collection: db.Collection("goals"),
	}
}

// CreateGoal creates a new goal in the database
func (r *GoalRepository) CreateGoal(ctx context.Context, goal *models.Goal) (*models.Goal, error) {
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, goal)
	if err != nil {
		return nil, fmt.Errorf("failed to insert goal: %v", err)
	}

	// Cast the inserted ID and assign it to the goal object
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to cast inserted ID")
	}
	goal.ID = insertedID

	return goal, nil
}

// GetGoalByID fetches a goal by its ID
func (r *GoalRepository) GetGoalByID(ctx context.Context, id primitive.ObjectID) (*models.Goal, error) {
	var goal models.Goal

	// Find the goal by its ID
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&goal)
	if err != nil {
		return nil, fmt.Errorf("failed to find goal by id: %v", err)
	}

	return &goal, nil
}

// UpdateGoal updates an existing goal in the database
func (r *GoalRepository) UpdateGoal(ctx context.Context, id primitive.ObjectID, goal *models.Goal) (*models.Goal, error) {
	goal.UpdatedAt = time.Now()

	// Update the goal in the database
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": goal},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update goal: %v", err)
	}

	return goal, nil
}

// DeleteGoal deletes a goal from the database by its ID
func (r *GoalRepository) DeleteGoal(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete goal: %v", err)
	}

	return nil
}

// GetAllGoals fetches all goals from the database
func (r *GoalRepository) GetAllGoals(ctx context.Context, limit int64) ([]models.Goal, error) {
	var goals []models.Goal

	findOptions := options.Find().SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goals: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var goal models.Goal
		if err := cursor.Decode(&goal); err != nil {
			return nil, fmt.Errorf("failed to decode goal: %v", err)
		}
		goals = append(goals, goal)
	}

	return goals, nil
}
