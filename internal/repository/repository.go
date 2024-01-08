package repository

import (
	"github.com/v7ktory/fullstack/internal/model"
	"gorm.io/gorm"
)

type Authorization interface {
	Create(user model.User) error
	Get(email string) (model.User, error)
}

type User interface {
	GetUserByID(userID string) (*model.User, error)
}
type Repository struct {
	Authorization
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		User:          NewUserRepository(db),
	}
}
