package services

import (
	"context"
	"errors"
	"log"
	"taskapp/internal/auth"
	"taskapp/internal/models"
	"taskapp/internal/repository"
)

type UserService struct {
	UserRepo repository.UserRepository
	TaskRepo repository.TaskRepository
}

func NewUserService(userRepo repository.UserRepository, taskRepo repository.TaskRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		TaskRepo: taskRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.UserRepo.CreateUser(ctx, user)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		log.Println("JWT generation failed:", err)
		return "", errors.New("could not generate token")
	}
	return token, nil
}
