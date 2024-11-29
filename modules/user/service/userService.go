package service

import (
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
)

type UserService interface {
	FindUserByID(userID string) (*_userModel.ProfileUser, error)
	EditingUser(userID string, userEditingReq *_userModel.UserEditingReq) (*_userModel.ProfileUser, error)

	FindLocation(userID string) (*_userModel.UserLocation, error)
	EditingLocation(userID string, userLocationReq *_userModel.UserLocationReq) (*_userModel.UserLocation, error)
}
