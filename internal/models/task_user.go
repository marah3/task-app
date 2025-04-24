package models

type TaskUser struct {
	ID     int `bun:"id,pk,autoincrement"` // Primary key for the join table
	TaskID int `bun:"task_id,pk"`          // Foreign key to Task
	UserID int `bun:"user_id,pk"`          // Foreign key to User

	// Relations
	Task *Task `bun:"rel:belongs-to,join:task_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
