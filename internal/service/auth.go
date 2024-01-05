package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/internal/repository"
	"github.com/v7ktory/fullstack/pkg/hasher"
)

type AuthService struct {
	repos  repository.Authorization
	hasher hasher.PasswordHasher
}

func NewAuthService(repos repository.Authorization, hasher hasher.PasswordHasher) *AuthService {
	return &AuthService{
		repos:  repos,
		hasher: hasher,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user model.User) error {

	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.repos.Create(user)
}

func (s *AuthService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repos.GetByCredentials(email)
	if err != nil {
		return "", err
	}

	// Check if the provided password matches the stored hashed password
	if !s.hasher.CheckPasswordHash(password, user.Password) {
		return "", errors.New("authentication failed: incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
