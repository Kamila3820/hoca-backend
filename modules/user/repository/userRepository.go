package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
)

type UserRepository interface {
	Creating(userEntity *entities.User) (*entities.User, error)
	FindByID(userID string) (*entities.User, error)

	FindUserByID(userID string) (*entities.User, error)
	EditingUser(userID string, userEditingReq *_userModel.UserEditingReq) (string, error)

	EditingUserLocation(userID string, userLocationReq *_userModel.UserLocationReq) (string, error)
}
