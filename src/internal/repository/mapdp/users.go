package mapdp

type UserRepository struct {
	db map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: map[string]string{},
	}
}
