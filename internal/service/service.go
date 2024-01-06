package service

import (
	"context"

	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/internal/repository"
	"github.com/v7ktory/fullstack/pkg/hasher"
)

type Authorization interface {
	SignUp(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, hasher hasher.PasswordHasher) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, hasher),
	}
}
