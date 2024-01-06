package service

import (
	"context"

	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/internal/repository"
	"github.com/v7ktory/fullstack/pkg/hasher"
	"github.com/v7ktory/fullstack/pkg/token"
)

type Authorization interface {
	SignUp(ctx context.Context, user model.User) (string, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type User interface {
	GetUserByID(userID string) (*model.User, error)
}

type Service struct {
	Authorization
	User
}

func NewService(repos *repository.Repository, hasher hasher.PasswordHasher, jwt token.JWTService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, hasher, jwt),
		User:          NewUserService(repos.User),
	}
}
