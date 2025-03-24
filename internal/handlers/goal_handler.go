package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"github.com/Dias221467/Achievemenet_Manager/internal/services"
	"github.com/Dias221467/Achievemenet_Manager/pkg/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GoalHandler handles HTTP requests related to goals.
type GoalHandler struct {
	Service *services.GoalService
}

// NewGoalHandler creates a new instance of GoalHandler.
func NewGoalHandler(service *services.GoalService) *GoalHandler {
	return &GoalHandler{Service: service}
}

// CreateGoalHandler handles the creation of a new goal.
func (h *GoalHandler) CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	// Get the logged-in user from JWT token
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Convert UserID to ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	goal.UserID = userID
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	//  Validate & Parse Due Date (Optional)
	if !goal.DueDate.IsZero() && goal.DueDate.Before(time.Now()) {
		http.Error(w, "Due date cannot be in the past", http.StatusBadRequest)
		return
	}

	//  Validate & Set Category (Optional)
	if goal.Category != "" {
		if _, exists := models.AllowedCategories[goal.Category]; !exists {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}
	}

	// Initialize progress field
	goal.Progress = make(map[string]bool)
	for _, step := range goal.Steps {
		goal.Progress[step] = false
	}

	// Save to DB
	createdGoal, err := h.Service.CreateGoal(r.Context(), &goal)
	if err != nil {
		http.Error(w, "Failed to create goal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdGoal)
}

// GetGoalHandler handles fetching a single goal by its ID.
func (h *GoalHandler) GetGoalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Fetch the goal from DB
	goal, err := h.Service.GetGoal(r.Context(), goalID)
	if err != nil || goal == nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	// Check if the goal is overdue
	if !goal.DueDate.IsZero() && goal.DueDate.Before(time.Now()) {
		goal.Status = "expired"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goal)
}

// UpdateGoalHandler handles updating an existing goal.
func (h *GoalHandler) UpdateGoalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Get the logged-in user
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert goalID to ObjectID
	objID, err := primitive.ObjectIDFromHex(goalID)
	if err != nil {
		http.Error(w, "Invalid goal ID", http.StatusBadRequest)
		return
	}

	// Fetch the existing goal
	existingGoal, err := h.Service.GetGoal(r.Context(), goalID)
	if err != nil || existingGoal == nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	// Ensure the logged-in user is the owner of the goal
	if existingGoal.UserID.Hex() != claims.UserID {
		http.Error(w, "Forbidden: You can only update your own goals", http.StatusForbidden)
		return
	}

	// Decode request body
	var updatedGoal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&updatedGoal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//  Validate & Parse Due Date (Optional)
	if !updatedGoal.DueDate.IsZero() && updatedGoal.DueDate.Before(time.Now()) {
		http.Error(w, "Due date cannot be in the past", http.StatusBadRequest)
		return
	}

	// Sync Progress Field (Only Keep Steps That Exist in Updated Goal)
	newProgress := make(map[string]bool)

	//  Validate & Set Category (Optional)
	if updatedGoal.Category != "" {
		if _, exists := models.AllowedCategories[updatedGoal.Category]; !exists {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}
	}

	// Create a set of valid steps (to remove old progress)
	validSteps := make(map[string]bool)
	for _, step := range updatedGoal.Steps {
		validSteps[step] = true
	}

	// Keep only the progress of valid steps
	for step, done := range existingGoal.Progress {
		if validSteps[step] {
			newProgress[step] = done // Keep existing progress
		}
	}

	// Add new steps with default `false`
	for _, step := range updatedGoal.Steps {
		if _, exists := newProgress[step]; !exists {
			newProgress[step] = false
		}
	}

	//  Assign updated values
	updatedGoal.ID = objID
	updatedGoal.UserID = existingGoal.UserID
	updatedGoal.Progress = newProgress // Ensure old steps are removed
	updatedGoal.UpdatedAt = time.Now()

	// Save the updated goal
	updatedGoalData, err := h.Service.UpdateGoal(r.Context(), goalID, &updatedGoal)
	if err != nil {
		http.Error(w, "Failed to update goal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGoalData)
}

func (h *GoalHandler) UpdateGoalProgressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Get logged-in user
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch goal from DB
	goal, err := h.Service.GetGoal(r.Context(), goalID)
	if err != nil || goal == nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	// Ensure the logged-in user owns the goal
	if goal.UserID.Hex() != claims.UserID {
		http.Error(w, "Forbidden: You can only update your own goals", http.StatusForbidden)
		return
	}

	// Decode request body
	var progressUpdate struct {
		Step string `json:"step"`
		Done bool   `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&progressUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Ensure the step exists in the goal
	if _, exists := goal.Progress[progressUpdate.Step]; !exists {
		http.Error(w, "Step not found in goal", http.StatusBadRequest)
		return
	}

	// Update step progress
	goal.Progress[progressUpdate.Step] = progressUpdate.Done

	// Check if all steps are completed
	allCompleted := true
	for _, done := range goal.Progress {
		if !done {
			allCompleted = false
			break
		}
	}

	// Update goal status
	if allCompleted {
		goal.Status = "completed"
	} else {
		goal.Status = "in_progress"
	}

	goal.UpdatedAt = time.Now()

	// Save changes
	updatedGoal, err := h.Service.UpdateGoal(r.Context(), goalID, goal)
	if err != nil {
		http.Error(w, "Failed to update progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGoal)
}

// DeleteGoalHandler handles deleting a goal by its ID.
func (h *GoalHandler) DeleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Get the logged-in user from JWT token
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch the goal from DB
	goal, err := h.Service.GetGoal(r.Context(), goalID)
	if err != nil || goal == nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	// Check if the logged-in user is the owner
	if goal.UserID.Hex() != claims.UserID {
		http.Error(w, "Forbidden: You can only delete your own goals", http.StatusForbidden)
		return
	}

	// Perform delete
	err = h.Service.DeleteGoal(r.Context(), goalID)
	if err != nil {
		http.Error(w, "Failed to delete goal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllGoalsHandler handles fetching all goals, with an optional limit.

// Its not working right now, we will need it later when we will add admins and their rights with functions
func (h *GoalHandler) GetAllGoalsHandler(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	var limit int64 = 10 // default limit
	if limitParam != "" {
		parsed, err := strconv.ParseInt(limitParam, 10, 64)
		if err == nil {
			limit = parsed
		}
	}

	goals, err := h.Service.GetAllGoals(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goals)
}

func (h *GoalHandler) GetGoalProgressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Get the logged-in user
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch the goal from DB
	goal, err := h.Service.GetGoal(r.Context(), goalID)
	if err != nil || goal == nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	// Ensure the logged-in user is the owner of the goal
	if goal.UserID.Hex() != claims.UserID {
		http.Error(w, "Forbidden: You can only view your own goal progress", http.StatusForbidden)
		return
	}

	// Return only the progress field
	response := map[string]interface{}{
		"progress": goal.Progress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GoalHandler) GetGoalsHandler(w http.ResponseWriter, r *http.Request) {
	// Get logged-in user
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert UserID to ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Get category filter from query params (optional)
	category := r.URL.Query().Get("category")

	// Fetch goals from DB with optional category filter
	goals, err := h.Service.GetGoals(r.Context(), userID, category)
	if err != nil {
		http.Error(w, "Failed to retrieve goals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goals)
}
