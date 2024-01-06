package repository

import (
	"fmt"
	"strings"

	"github.com/v7ktory/fullstack/internal/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(user model.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return model.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetByCredentials(email string) (model.User, error) {

	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if gorm.ErrRecordNotFound == err {
		return model.User{}, fmt.Errorf("user not found: %w", model.ErrUserNotFound)
	} else if err != nil {
		return user, fmt.Errorf("failed to get user by credentials: %w", err)
	}
	return user, nil
}
