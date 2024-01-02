package model

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with such email already exists")
	ErrBadCredentials    = errors.New("bad credentials")
	ErrInternalServer    = errors.New("internal server error")
	ErrUnauthorized      = errors.New("unauthorized")
)
