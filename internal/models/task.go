package models

import (
	"time"
)

type Task struct {
	ID          int       `json:"id" bun:"id,pk,autoincrement"`
	Title       string    `json:"title" bun:"title"`
	Description string    `json:"description" bun:"description"`
	UserID      int       `json:"user_id" bun:"user_id"`
	Status      string    `json:"status" bun:"status,default:'pending'"`
	AssignedAt  time.Time `json:"assigned_at" bun:"assigned_at,default:current_timestamp"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
}
