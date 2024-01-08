package repository

import (
	"errors"
	"fmt"
	"strconv"

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

	numericUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user ID to numeric format: %w", err)
	}

	err = r.db.Where("id = ?", numericUserID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found: %w", model.ErrUserNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}
