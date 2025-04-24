package repository

import (
	"context"
	"log"
	"taskapp/internal/models"

	"github.com/uptrace/bun"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error)
	GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error)
	ListTasks(ctx context.Context) ([]models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	AssignTaskToUser(ctx context.Context, taskID int64, userID int64) error
}

type taskRepository struct {
	db *bun.DB
}

func NewTaskRepository(db *bun.DB) TaskRepository {
	return &taskRepository{db: db}
}

// CreateTask
func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	_, err := r.db.NewInsert().Model(task).Exec(ctx)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		return err
	}
	return nil
}

// GetTasksByUserID (updated to use task_user join table)
func (r *taskRepository) GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.NewSelect().
		Model(&tasks).
		Join("JOIN task_users tu ON tu.task_id = task.id").
		Where("tu.user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		log.Printf("Error fetching tasks for user %d: %v", userID, err)
		return nil, err
	}
	return tasks, nil
}

// GetTaskByID
func (r *taskRepository) GetTaskByID(ctx context.Context, taskID int64) (*models.Task, error) {
	var task models.Task
	err := r.db.NewSelect().Model(&task).Where("id = ?", taskID).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching task by ID %d: %v", taskID, err)
		return nil, err
	}
	return &task, nil
}

// UpdateTask
func (r *taskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	_, err := r.db.NewUpdate().
		Model(task).
		Column("title", "description", "status", "updated_at").
		Where("id = ?", task.ID).
		Exec(ctx)

	if err != nil {
		log.Printf("Error updating task ID %d: %v", task.ID, err)
		return err
	}
	return nil
}

func (r *taskRepository) AssignTaskToUser(ctx context.Context, taskID int64, userID int64) error {
	taskUser := &models.TaskUser{
		TaskID: int(taskID),
		UserID: int(userID),
	}

	_, err := r.db.NewInsert().
		Model(taskUser).
		On("CONFLICT (task_id, user_id) DO NOTHING"). // Prevent duplicate assignment
		Exec(ctx)

	if err != nil {
		log.Printf("Error assigning task ID %d to user %d: %v", taskID, userID, err)
		return err
	}
	return nil
}

// ListTasks
func (r *taskRepository) ListTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.NewSelect().Model(&tasks).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching all tasks: %v", err)
		return nil, err
	}
	return tasks, nil
}
