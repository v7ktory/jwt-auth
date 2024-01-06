package service

import (
	"context"
	"errors"
	"fmt"
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

func (s *AuthService) SignUp(ctx context.Context, user model.User) error {

	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	err = s.repos.Create(user)
	if errors.Is(err, model.ErrUserAlreadyExists) {
		return model.ErrUserAlreadyExists
	} else if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {

	user, err := s.repos.GetByCredentials(email)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	if user == (model.User{}) {
		return "", errors.New("authentication failed: user not found")
	}

	if !s.hasher.CheckPasswordHash(password, user.Password) {
		return "", errors.New("authentication failed: incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}
