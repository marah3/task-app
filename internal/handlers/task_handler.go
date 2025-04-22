package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"taskapp/internal/cache"
	"taskapp/internal/models"
	"taskapp/internal/repository"
	"time"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	repo  repository.TaskRepository
	cache *cache.RedisCache // Redis cache injected
}

// NewTaskHandler creates a new TaskHandler with repository and cache
func NewTaskHandler(repo repository.TaskRepository, cache *cache.RedisCache) *TaskHandler {
	return &TaskHandler{repo: repo, cache: cache}
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Log the incoming body for debugging
	log.Println("Incoming request body:", r.Body)

	// Decode the request body to a Task struct
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	// Log the decoded task for debugging
	log.Printf("Decoded task: %+v\n", task)

	// Validate the user ID (if necessary)
	if task.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		log.Println("User ID is missing or invalid")
		return
	}

	// Set timestamps
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Log the task before creating it
	log.Printf("Creating task for user %d: %+v\n", task.UserID, task)

	// Create the task
	err := h.repo.CreateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		log.Println("Error creating task:", err)
		return
	}

	// Log the created task ID for debugging
	log.Printf("Created task with ID: %d\n", task.ID)

	// Send the response with the task data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetUserTasks
func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Println("Error parsing user ID:", err)
		return
	}

	// Check if the tasks are cached in Redis
	cacheKey := "user_tasks_" + strconv.FormatInt(userID, 10)
	cachedTasks, err := h.cache.Get(cacheKey)
	if err == nil && cachedTasks != "" {
		log.Println("Returning tasks from cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedTasks))
		return
	}

	// Fetch tasks from the database if not in cache
	tasks, err := h.repo.GetTasksByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		log.Println("Error fetching tasks for user:", userID, err)
		return
	}

	// Cache the tasks for future requests
	tasksJSON, _ := json.Marshal(tasks)
	h.cache.Set(cacheKey, string(tasksJSON)) // Save to cache with an appropriate expiration time

	// Return the tasks in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// ListTasks returns all tasks
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.ListTasks(r.Context())
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		log.Println("Error fetching all tasks:", err)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskIDStr := mux.Vars(r)["id"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		log.Println("Error parsing task ID:", err)
		return
	}

	task, err := h.repo.GetTaskByID(r.Context(), taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Println("Task not found with ID:", taskID)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// AssignTaskToUser
func (h *TaskHandler) AssignTaskToUser(w http.ResponseWriter, r *http.Request) {
	taskIDStr := mux.Vars(r)["id"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		log.Println("Error parsing task ID:", err)
		return
	}

	var input struct {
		UserID int64 `json:"user_id"`
	}

	// Decode the user ID from the request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	// Validate the user ID
	if input.UserID == 0 {
		http.Error(w, "User ID is required to assign task", http.StatusBadRequest)
		log.Println("User ID is missing or invalid")
		return
	}

	// Assign the task to the user
	err = h.repo.AssignTaskToUser(r.Context(), taskID, input.UserID)
	if err != nil {
		http.Error(w, "Error assigning task", http.StatusInternalServerError)
		log.Println("Error assigning task:", err)
		return
	}

	// Log the successful assignment
	log.Printf("Assigned task ID: %d to user %d\n", taskID, input.UserID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id": taskID,
		"user_id": input.UserID,
		"message": "Task assigned successfully",
	})
}
