package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/czeful/diplom_back/internal/config"
	"github.com/czeful/diplom_back/internal/database"
	"github.com/czeful/diplom_back/internal/handlers"
	"github.com/czeful/diplom_back/internal/repository"
	"github.com/czeful/diplom_back/internal/services"
	"github.com/czeful/diplom_back/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration from .env file
	cfg := config.LoadConfig()


	// Connect to MongoDB Atlas
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Initialize repositories, services, and handlers for goals
	goalRepo := repository.NewGoalRepository(db)
	goalService := services.NewGoalService(goalRepo)
	goalHandler := handlers.NewGoalHandler(goalService)

	// Initialize repositories, services, and handlers for users
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, cfg)

	// Initialize Gorilla Mux router
	router := mux.NewRouter()

	// Apply authentication middleware to goal routes
	protectedRoutes := router.PathPrefix("/goals").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protectedRoutes.HandleFunc("", goalHandler.CreateGoalHandler).Methods("POST")
	protectedRoutes.HandleFunc("/{id}", goalHandler.GetGoalHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", goalHandler.UpdateGoalHandler).Methods("PUT")
	protectedRoutes.HandleFunc("/{id}", goalHandler.DeleteGoalHandler).Methods("DELETE")
	protectedRoutes.HandleFunc("/{id}/progress", goalHandler.UpdateGoalProgressHandler).Methods("PATCH")
	protectedRoutes.HandleFunc("/{id}/progress", goalHandler.GetGoalProgressHandler).Methods("GET")
	protectedRoutes.HandleFunc("", goalHandler.GetGoalsHandler).Methods("GET")

	// Register User routes
	router.HandleFunc("/users/register", userHandler.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/users/login", userHandler.LoginUserHandler).Methods("POST")

	// Protected user routes (only authenticated users can access)
	protectedUserRoutes := router.PathPrefix("/users").Subrouter()
	protectedUserRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protectedUserRoutes.HandleFunc("/{id}", userHandler.GetUserHandler).Methods("GET")
	protectedUserRoutes.HandleFunc("/{id}", userHandler.UpdateUserHandler).Methods("PUT")

	// Apply middleware for logging
	router.Use(middleware.LoggingMiddleware)

	// Start the HTTP server
	port := cfg.Port
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
