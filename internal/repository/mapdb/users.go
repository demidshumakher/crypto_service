package mapdb

import (
	"cryptoserver/domain"
	"cryptoserver/internal/repository"
)

type UserRepository struct {
	db map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: map[string]string{},
	}
}
func (ur *UserRepository) Login(user *domain.User) (string, error) {
	if _, ok := ur.db[user.Username]; !ok {
		return "", domain.ErrUserNotFound
	}

	hashed := ur.db[user.Username]
	if repository.CheckPassword(user.Password, hashed) {
		return user.Username, nil
	}

	return "", domain.ErrIncorrectPassword
}

func (ur *UserRepository) Register(user *domain.User) (string, error) {
	if _, ok := ur.db[user.Username]; ok {
		return "", domain.ErrAlreadyExist
	}

	hashed, err := repository.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	ur.db[user.Username] = hashed
	return user.Username, nil
}
