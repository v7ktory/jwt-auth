package repository

import (
	"github.com/v7ktory/fullstack/internal/model"
	"gorm.io/gorm"
)

type Authorization interface {
	Create(user model.User) error
	GetByCredentials(email, password string) (model.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
	}
}
