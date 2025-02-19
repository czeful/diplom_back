package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dias221467/Achievemenet_Manager/internal/config"
	"github.com/Dias221467/Achievemenet_Manager/internal/models"
	"github.com/Dias221467/Achievemenet_Manager/internal/services"
	jwtutil "github.com/Dias221467/Achievemenet_Manager/pkg/jwt"
	"github.com/Dias221467/Achievemenet_Manager/pkg/middleware"
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
	requestedUserID := vars["id"]

	// Get the logged-in user from the request context
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ensure that the requested user ID matches the logged-in user’s ID
	if requestedUserID != claims.UserID {
		http.Error(w, "Forbidden: You can only access your own profile", http.StatusForbidden)
		return
	}

	// Fetch the user from the database
	user, err := h.Service.GetUser(r.Context(), requestedUserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedUserID := vars["id"]

	// Get logged-in user
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ensure only the logged-in user can update their own profile
	if requestedUserID != claims.UserID {
		http.Error(w, "Forbidden: You can only update your own profile", http.StatusForbidden)
		return
	}

	// Decode request body
	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update user in DB
	updatedUserData, err := h.Service.UpdateUser(r.Context(), requestedUserID, &updatedUser)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUserData)
}
