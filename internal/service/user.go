package service

import (
	"fmt"

	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/internal/repository"
)

type UserService struct {
	repos repository.User
}

func NewUserService(repos repository.User) *UserService {
	return &UserService{
		repos: repos,
	}
}

func (s *UserService) GetUserByID(userID string) (*model.User, error) {

	user, err := s.repos.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
