package domain

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("required item not found")
	ErrAlreadyExist      = errors.New("item already exist")
	ErrInvalidToken      = errors.New("invalid jwt token")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotFound      = errors.New("user not found")
	ErrBadRequest        = errors.New("bad request")
	ErrInvalidInterval   = errors.New("invalid interval")
	ErrIncorrectPassword = errors.New("incorrect password")
)
