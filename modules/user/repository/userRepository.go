package repository

import "github.com/Kamila3820/hoca-backend/entities"

type UserRepository interface {
	Creating(userEntity *entities.User) (*entities.User, error)
	FindByID(userID string) (*entities.User, error)
}
