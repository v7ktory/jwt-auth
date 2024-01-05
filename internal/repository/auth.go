package repository

import (
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
	if gorm.ErrDuplicatedKey == err {
		return model.ErrUserAlreadyExists
	}
	return err
}

func (r *AuthRepository) GetByCredentials(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if gorm.ErrRecordNotFound == err {
		return model.User{}, model.ErrUserNotFound
	}
	return user, err
}
