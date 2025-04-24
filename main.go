package main

import (
	"fmt"
	"log"
	"net/http"
	"taskapp/config"
	"taskapp/internal/cache"
	"taskapp/internal/database"
	"taskapp/internal/dependencies"
	"taskapp/internal/handlers"
	"taskapp/internal/repository"
	"taskapp/internal/routes"
	"taskapp/internal/services" // Import the services package

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

	// Initialize repositories
	userRepo := repository.NewUserRepository(deps.DB)
	taskRepo := repository.NewTaskRepository(deps.DB)

	userService := services.NewUserService(userRepo, taskRepo)
	taskService := services.NewTaskService(taskRepo, redisClient)

	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Set up the router
	router := routes.SetupRoutes(userHandler, taskHandler)

	// Start the server
	fmt.Printf("Server is running on port %s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
