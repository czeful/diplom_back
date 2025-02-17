package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"github.com/Dias221467/Achievemenet_Manager/internal/services"
	"github.com/gorilla/mux"
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
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	createdGoal, err := h.Service.CreateGoal(r.Context(), &goal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	id := vars["id"]

	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedGoal, err := h.Service.UpdateGoal(r.Context(), id, &goal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGoal)
}

// DeleteGoalHandler handles deleting a goal by its ID.
func (h *GoalHandler) DeleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.Service.DeleteGoal(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
