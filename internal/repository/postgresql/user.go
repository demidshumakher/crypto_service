package postgresql

import (
	"cryptoserver/domain"
	"cryptoserver/internal/repository"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Login(user *domain.User) (string, error) {
	var hashedPassword string

	err := ur.db.QueryRow("select password from users where username = $1", user.Username).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}

	if repository.CheckPassword(user.Password, hashedPassword) {
		return user.Username, nil
	}

	return "", domain.ErrIncorrectPassword
}

func (ur *UserRepository) Register(user *domain.User) (string, error) {
	hashedPassword, err := repository.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	_, err = ur.db.Exec("insert into users (username, password) values ($1, $2)", user.Username, hashedPassword)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}
