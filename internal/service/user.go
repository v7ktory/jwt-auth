package service

import (
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
	return s.repos.GetUserByID(userID)
}
