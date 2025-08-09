package domain

import "errors"

var (
	ErrNotFound     = errors.New("required item not found")
	ErrAlreadyExist = errors.New("item already exist")
)
