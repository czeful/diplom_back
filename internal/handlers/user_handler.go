package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dias221467/Achievemenet_Manager/internal/config"
	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"github.com/Dias221467/Achievemenet_Manager/internal/services"
	jwtutil "github.com/Dias221467/Achievemenet_Manager/pkg/jwt"
	"github.com/gorilla/mux"
)

// UserHandler handles HTTP requests related to user operations.
type UserHandler struct {
	Service *services.UserService
	Config  *config.Config
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(service *services.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		Service: service,
		Config:  cfg,
	}
}

// RegisterUserHandler handles user registration.
func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	createdUser, err := h.Service.RegisterUser(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}

// LoginUserHandler handles user login.
func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Define a simple struct to receive login credentials.
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.Service.AuthenticateUser(r.Context(), credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := jwtutil.GenerateToken(user.ID.Hex(), user.Email, h.Config.JWTSecret, 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token and user details
	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserHandler handles fetching a user by ID.
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.Service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
