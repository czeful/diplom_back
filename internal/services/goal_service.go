package services

import (
	"context"
	"fmt"

	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"github.com/Dias221467/Achievemenet_Manager/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GoalService encapsulates the business logic for goals.
type GoalService struct {
	repo *repository.GoalRepository
}

// NewGoalService creates a new instance of GoalService.
func NewGoalService(repo *repository.GoalRepository) *GoalService {
	return &GoalService{
		repo: repo,
	}
}

// CreateGoal processes the goal creation logic and stores it in the database.
func (s *GoalService) CreateGoal(ctx context.Context, goal *models.Goal) (*models.Goal, error) {
	// Here you can add additional business logic,
	// such as validating the goal, generating steps automatically, etc.
	if goal.Name == "" {
		return nil, fmt.Errorf("goal name is required")
	}
	createdGoal, err := s.repo.CreateGoal(ctx, goal)
	if err != nil {
		return nil, fmt.Errorf("failed to create goal: %v", err)
	}
	return createdGoal, nil
}

// GetGoal retrieves a goal by its ID.
func (s *GoalService) GetGoal(ctx context.Context, id string) (*models.Goal, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %v", err)
	}
	goal, err := s.repo.GetGoalByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %v", err)
	}
	return goal, nil
}

// UpdateGoal updates an existing goal.
func (s *GoalService) UpdateGoal(ctx context.Context, id string, updatedGoal *models.Goal) (*models.Goal, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %v", err)
	}
	goal, err := s.repo.UpdateGoal(ctx, objID, updatedGoal)
	if err != nil {
		return nil, fmt.Errorf("failed to update goal: %v", err)
	}
	return goal, nil
}

// DeleteGoal removes a goal from the database.
func (s *GoalService) DeleteGoal(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid goal ID: %v", err)
	}
	if err := s.repo.DeleteGoal(ctx, objID); err != nil {
		return fmt.Errorf("failed to delete goal: %v", err)
	}
	return nil
}

// GetAllGoals retrieves a list of goals with an optional limit.
func (s *GoalService) GetAllGoals(ctx context.Context, limit int64) ([]models.Goal, error) {
	goals, err := s.repo.GetAllGoals(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goals: %v", err)
	}
	return goals, nil
}
