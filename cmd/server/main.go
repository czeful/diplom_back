package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dias221467/Achievemenet_Manager/internal/config"
	"github.com/Dias221467/Achievemenet_Manager/internal/database"
	"github.com/Dias221467/Achievemenet_Manager/internal/handlers"
	"github.com/Dias221467/Achievemenet_Manager/internal/repository"
	"github.com/Dias221467/Achievemenet_Manager/internal/services"
	"github.com/Dias221467/Achievemenet_Manager/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration from .env
	cfg := config.LoadConfig()

	// Connect to MongoDB Atlas
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Initialize repository, service, and handler for goals
	goalRepo := repository.NewGoalRepository(db)
	goalService := services.NewGoalService(goalRepo)
	goalHandler := handlers.NewGoalHandler(goalService)

	// Initialize Gorilla Mux router
	router := mux.NewRouter()

	// Register API routes for goal operations
	router.HandleFunc("/goals", goalHandler.CreateGoalHandler).Methods("POST")
	router.HandleFunc("/goals", goalHandler.GetAllGoalsHandler).Methods("GET")
	router.HandleFunc("/goals/{id}", goalHandler.GetGoalHandler).Methods("GET")
	router.HandleFunc("/goals/{id}", goalHandler.UpdateGoalHandler).Methods("PUT")
	router.HandleFunc("/goals/{id}", goalHandler.DeleteGoalHandler).Methods("DELETE")

	// Apply logging middleware to all routes
	router.Use(middleware.LoggingMiddleware)

	// Start the HTTP server
	port := cfg.Port
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
