package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/v7ktory/fullstack/internal/model"
	"github.com/v7ktory/fullstack/internal/repository"
	"github.com/v7ktory/fullstack/pkg/hasher"
	"github.com/v7ktory/fullstack/pkg/token"
)

type AuthService struct {
	repos  repository.Authorization
	hasher hasher.PasswordHasher
	jwt    token.JWTService
}

func NewAuthService(repos repository.Authorization, hasher hasher.PasswordHasher, jwt token.JWTService) *AuthService {
	return &AuthService{
		repos:  repos,
		hasher: hasher,
		jwt:    jwt,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user model.User) (string, error) {

	// Хешируем пароль
	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Заменяем оригинальный пароль на хешированный
	user.Password = string(hashedPassword)

	// Создаем пользователя
	err = s.repos.Create(user)
	if errors.Is(err, model.ErrUserAlreadyExists) {
		return "", model.ErrUserAlreadyExists
	} else if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// Генерируем токен с UserID, Email и Username
	token, err := s.jwt.GenerateJWT(user.Email, user.Username, strconv.Itoa(int(user.ID)))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {

	// Получаем пользователя по электронной почте
	user, err := s.repos.GetByCredentials(email)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	// Проверяем, найден ли пользователь
	if user == (model.User{}) {
		return "", errors.New("authentication failed: user not found")
	}

	// Проверяем корректность пароля
	if !s.hasher.CheckPasswordHash(password, user.Password) {
		return "", errors.New("authentication failed: incorrect password")
	}

	// Генерируем токен с UserID, Email и Username
	tokenString, err := s.jwt.GenerateJWT(email, user.Username, strconv.Itoa(int(user.ID)))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}
