package models

import "time"

type User struct {
	ID        int       `bun:"id,pk,autoincrement"`
	Username  string    `bun:"username,notnull"`
	Email     string    `bun:"email,unique,notnull"`
	Password  string    `bun:"password,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
