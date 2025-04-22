package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"taskapp/internal/auth"
	"taskapp/internal/models"
	"taskapp/internal/repository"
)

type UserHandler struct {
	UserRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body into a User struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		log.Println("Error hashing password:", err)
		return
	}
	user.Password = hashedPassword

	err = h.UserRepo.CreateUser(context.Background(), &user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println("Error creating user:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginUser
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Email == "" || input.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Invalid email or password input:", err)
		return
	}

	user, err := h.UserRepo.FindUserByEmail(context.Background(), input.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Println("Error fetching user:", err)
		return
	}

	if !auth.CheckPasswordHash(input.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Println("Error generating JWT:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
