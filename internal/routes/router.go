package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"taskapp/internal/auth"
	"taskapp/internal/handlers"
)

func SetupRoutes(userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/taskapp/users/register", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/taskapp/users/login", userHandler.LoginUser).Methods("POST")

	// Protected Routes
	router.Handle("/taskapp/tasks/create", auth.AuthMiddleware(http.HandlerFunc(taskHandler.CreateTask))).Methods("POST")
	router.Handle("/taskapp/users/{user_id}/tasks", auth.AuthMiddleware(http.HandlerFunc(taskHandler.GetUserTasks))).Methods("GET")
	router.Handle("/taskapp/tasks/{id}", auth.AuthMiddleware(http.HandlerFunc(taskHandler.GetTaskByID))).Methods("GET")
	router.Handle("/taskapp/tasks", auth.AuthMiddleware(http.HandlerFunc(taskHandler.ListTasks))).Methods("GET")
	router.Handle("/taskapp/tasks/{id}/assign", auth.AuthMiddleware(http.HandlerFunc(taskHandler.AssignTaskToUser))).Methods("PUT")

	return router
}
