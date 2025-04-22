package main

import (
	"fmt"
	"log"
	"net/http"
	"taskapp/config"
	"taskapp/internal/cache"
	"taskapp/internal/database"
	"taskapp/internal/dependencies" // new import
	"taskapp/internal/handlers"
	"taskapp/internal/repository"
	"taskapp/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Initialize the database connection
	database.Init(cfg)

	// Initialize Redis client
	redisClient := cache.NewRedisCache()

	// Create the dependencies struct
	deps := &dependencies.Dependencies{
		DB:    database.DB,
		Redis: redisClient,
	}

	// Initialize repositories and handlers using dependencies
	userRepo := repository.NewUserRepository(deps.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	taskRepo := repository.NewTaskRepository(deps.DB)
	taskHandler := handlers.NewTaskHandler(taskRepo, deps.Redis)

	// Set up the router
	router := routes.SetupRoutes(userHandler, taskHandler)

	fmt.Printf("Server is running on port %s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
