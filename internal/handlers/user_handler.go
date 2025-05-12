package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"taskapp/internal/models"
	"taskapp/internal/services" // Make sure to import your service package
)

type UserHandler struct {
	UserService *services.UserService // Ensure UserService is a pointer to the service
}

// NewUserHandler creates a new UserHandler with the UserService injected as a pointer
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding user:", err)
		return
	}

	if err := h.UserService.RegisterUser(context.Background(), &user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println("Registration error:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Login input error:", err)
		return
	}

	token, err := h.UserService.LoginUser(context.Background(), input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
