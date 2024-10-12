package repository

import "github.com/Kamila3820/hoca-backend/entities"

type AccountRepository interface {
	CheckDuplicateEmail(email string) error
	Creating(userEntity *entities.User) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
}
