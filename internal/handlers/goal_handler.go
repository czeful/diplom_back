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

	// Convert UserID from string to ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	goal.UserID = userID // Assign the logged-in user's ID
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	// Initialize `progress`: All steps = false
	goal.Progress = make(map[string]bool)
	for _, step := range goal.Steps {
		goal.Progress[step] = false
	}

	// Save the goal in DB
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
	id := vars["id"]

	goal, err := h.Service.GetGoal(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goal)
}

// UpdateGoalHandler handles updating an existing goal.
func (h *GoalHandler) UpdateGoalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Get the logged-in user from JWT token
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

	// Fetch the goal from DB
	goal, err := h.Service.GetGoal(r.Context(), goalID)
	if err != nil || goal == nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}

	// Check if the logged-in user is the owner of the goal
	if goal.UserID.Hex() != claims.UserID {
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

	// Perform update
	updatedGoal.ID = objID
	updatedGoal.UserID = goal.UserID // Keep the same owner
	updatedGoal.UpdatedAt = time.Now()

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

	// 🔹 Проверяем, что пользователь - владелец цели
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

	// 🔹 Проверяем, есть ли шаг в списке шагов
	if _, exists := goal.Progress[progressUpdate.Step]; !exists {
		http.Error(w, "Step not found in goal", http.StatusBadRequest)
		return
	}

	// 🔹 Обновляем статус шага
	goal.Progress[progressUpdate.Step] = progressUpdate.Done
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
