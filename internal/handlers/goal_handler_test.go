package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/czeful/diplom_back/internal/models"
	"github.com/czeful/diplom_back/internal/repository"
	"github.com/czeful/diplom_back/internal/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestRouter() (*mux.Router, *mongo.Database, context.Context, func()) {
	// Set up context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to the test MongoDB.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := client.Database("achievement_manager_test")
	// Drop collection for a clean test state.
	_ = db.Collection("goals").Drop(ctx)

	// Initialize repository, service, and handler.
	goalRepo := repository.NewGoalRepository(db)
	goalService := services.NewGoalService(goalRepo)
	goalHandler := NewGoalHandler(goalService)

	// Set up router and register routes.
	router := mux.NewRouter()
	router.HandleFunc("/goals", goalHandler.CreateGoalHandler).Methods("POST")
	router.HandleFunc("/goals/{id}", goalHandler.GetGoalHandler).Methods("GET")

	// Return a cleanup function.
	cleanup := func() {
		_ = client.Disconnect(ctx)
		cancel()
	}

	return router, db, ctx, cleanup
}

func TestCreateGoalHandler(t *testing.T) {
	router, _, ctx, cleanup := setupTestRouter()
	defer cleanup()

	// Define a new goal.
	goal := models.Goal{
		Name:        "Integration Test Goal",
		Description: "Testing the create goal endpoint",
		Steps:       []string{"Step A", "Step B"},
		Status:      "pending",
	}

	goalJSON, _ := json.Marshal(goal)
	req, err := http.NewRequest("POST", "/goals", bytes.NewBuffer(goalJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx) // Attach the context

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK (or 201 Created, if you change the status).
	assert.Equal(t, http.StatusOK, rr.Code)

	// Optionally, decode the response and check the values.
	var createdGoal models.Goal
	err = json.NewDecoder(rr.Body).Decode(&createdGoal)
	assert.NoError(t, err)
	assert.Equal(t, goal.Name, createdGoal.Name)
}

func TestGetGoalHandler(t *testing.T) {
	router, db, ctx, cleanup := setupTestRouter()
	defer cleanup()

	// Insert a goal directly using repository for testing the GET handler.
	goalRepo := repository.NewGoalRepository(db)
	goal := &models.Goal{
		Name:        "Test Goal for GET",
		Description: "A goal to test GET endpoint",
		Steps:       []string{"Step 1", "Step 2"},
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	createdGoal, err := goalRepo.CreateGoal(ctx, goal)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to GET the goal.
	req, err := http.NewRequest("GET", "/goals/"+createdGoal.ID.Hex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// Since our router uses mux, we need to set the route variables.
	router = mux.NewRouter()
	router.HandleFunc("/goals/{id}", NewGoalHandler(services.NewGoalService(goalRepo)).GetGoalHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	// Assert status code.
	assert.Equal(t, http.StatusOK, rr.Code)
	var fetchedGoal models.Goal
	err = json.NewDecoder(rr.Body).Decode(&fetchedGoal)
	assert.NoError(t, err)
	assert.Equal(t, createdGoal.Name, fetchedGoal.Name)
}
