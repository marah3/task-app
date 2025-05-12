package models

import "time"

type Task struct {
	ID          int       `json:"id" bun:"id,pk,autoincrement"`
	Title       string    `json:"title" bun:"title"`
	Description string    `json:"description" bun:"description"`
	Status      string    `json:"status" bun:"status,default:'pending'"`
	AssignedAt  time.Time `json:"assigned_at" bun:"assigned_at,default:current_timestamp"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
}

type TaskUser struct {
	ID     int `bun:"id,pk,autoincrement"` // Primary key for the join table
	TaskID int `bun:"task_id,pk"`          // Foreign key to Task
	UserID int `bun:"user_id,pk"`          // Foreign key to User

	// Relations
	Task *Task `bun:"rel:belongs-to,join:task_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

type User struct {
	ID        int       `bun:"id,pk,autoincrement"`
	Username  string    `bun:"username,notnull"`
	Email     string    `bun:"email,unique,notnull"`
	Password  string    `bun:"password,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
