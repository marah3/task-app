package repository

import (
	"context"
	"github.com/uptrace/bun"
	"log"
	"taskapp/internal/models"
)

// UserRepository defines methods for interacting with the user table
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *bun.DB // Changed from *database.DB to *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

// FindUserByEmail finds a user by email
func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil, err
	}
	return &user, nil
}
