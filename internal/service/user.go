package service

import (
	"cryptoserver/domain"
	"cryptoserver/pkg/jwt"
)

type UserRepository interface {
	Login(*domain.User) (string, error)
	Register(*domain.User) (string, error)
}

type UserService struct {
	ur        UserRepository
	jwtConfig jwt.JWTConfig
}

func NewUserService(jwtConfig jwt.JWTConfig, ur UserRepository) *UserService {
	return &UserService{
		ur:        ur,
		jwtConfig: jwtConfig,
	}
}

func (us *UserService) Register(user *domain.User) (string, error) {
	userId, err := us.ur.Register(user)
	if err != nil {
		return "", err
	}

	return us.jwtConfig.GenerateJWT(userId)
}

func (us *UserService) Login(user *domain.User) (string, error) {
	userId, err := us.ur.Login(user)
	if err != nil {
		return "", err
	}

	return us.jwtConfig.GenerateJWT(userId)
}
