package service

import (
	"context"

	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/internal/repository"
	"github.com/v7ktory/fullstack/pkg/hasher"
	"github.com/v7ktory/fullstack/pkg/token"
)

type Authorization interface {
	SignUp(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, hasher hasher.PasswordHasher, jwt token.JWTService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, hasher, jwt),
	}
}
