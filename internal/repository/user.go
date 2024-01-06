package repository

import (
	"fmt"

	"github.com/v7ktory/fullstack/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByID(userID string) (*model.User, error) {

	var user model.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if gorm.ErrRecordNotFound == err {
		return nil, fmt.Errorf("user not found: %w", model.ErrUserNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}
