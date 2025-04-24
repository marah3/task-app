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

type TaskHandler struct {
	repo  repository.TaskRepository
	cache *cache.RedisCache // Redis cache injected
}

func NewTaskHandler(repo repository.TaskRepository, cache *cache.RedisCache) *TaskHandler {
	return &TaskHandler{repo: repo, cache: cache}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	log.Println("Incoming request body:", r.Body)

	// Decode the request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	log.Printf("Creating task: %+v\n", task)

	// Call the repository to save the task to the database
	err := h.repo.CreateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		log.Println("Error creating task:", err)
		return
	}

	// Log the created task ID for debugging purposes
	log.Printf("Created task with ID: %d\n", task.ID)

	// Respond with the created task
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Println("Error parsing user ID:", err)
		return
	}

	cacheKey := "user_tasks_"
	cachedTasks, err := h.cache.Get(cacheKey)
	if err == nil && cachedTasks != "" {
		log.Println("Returning tasks from cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedTasks))
		return
	}

	tasks, err := h.repo.GetTasksByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		log.Println("Error fetching tasks for user:", userID, err)
		return
	}

	tasksJSON, _ := json.Marshal(tasks)
	h.cache.Set(cacheKey, string(tasksJSON)) // Save to cache

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.ListTasks(r.Context())
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		log.Println("Error fetching all tasks:", err)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

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

func (h *TaskHandler) AssignTaskToUser(w http.ResponseWriter, r *http.Request) {
	taskIDStr := mux.Vars(r)["id"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		log.Println("Error parsing task ID:", err)
		return
	}

	var input struct {
		UserIDs []int `json:"user_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if len(input.UserIDs) == 0 {
		http.Error(w, "At least one User ID is required", http.StatusBadRequest)
		log.Println("User IDs are missing or invalid")
		return
	}

	for _, userID := range input.UserIDs {
		err = h.repo.AssignTaskToUser(r.Context(), taskID, int64(userID))
		if err != nil {
			log.Printf("Error assigning task %d to user %d: %v", taskID, userID, err)
			http.Error(w, "Error assigning task", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Assigned task ID: %d to users: %+v\n", taskID, input.UserIDs)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id":  taskID,
		"user_ids": input.UserIDs,
		"message":  "Task assigned successfully",
	})
}
