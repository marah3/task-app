package services

import (
	"context"
	"encoding/json"
	"taskapp/internal/cache"
	"taskapp/internal/models"
	"taskapp/internal/repository"
	"time"
)

type TaskService struct {
	repo  repository.TaskRepository
	cache *cache.RedisCache
}

func NewTaskService(repo repository.TaskRepository, cache *cache.RedisCache) *TaskService {
	return &TaskService{repo: repo, cache: cache}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*models.Task, error) {
	task := &models.Task{
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.repo.CreateTask(ctx, task)
	return task, err
}

func (s *TaskService) GetUserTasks(ctx context.Context, userID int64) ([]models.Task, error) {
	cacheKey := "user_tasks_"

	// Try from cache
	cached, err := s.cache.Get(cacheKey)
	if err == nil && cached != "" {
		var tasks []models.Task
		if jsonErr := json.Unmarshal([]byte(cached), &tasks); jsonErr == nil {
			return tasks, nil
		}
	}

	// Fetch from DB and cache it
	tasks, err := s.repo.GetTasksByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	tasksJSON, _ := json.Marshal(tasks)
	_ = s.cache.Set(cacheKey, string(tasksJSON))
	return tasks, nil
}

func (s *TaskService) ListTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.ListTasks(ctx)
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	return s.repo.GetTaskByID(ctx, taskID)
}

func (s *TaskService) AssignTaskToUsers(ctx context.Context, taskID int64, userIDs []int) error {
	for _, userID := range userIDs {
		if err := s.repo.AssignTaskToUser(ctx, taskID, int64(userID)); err != nil {
			return err
		}
	}
	return nil
}
