package service

import _userModel "github.com/Kamila3820/hoca-backend/modules/user/model"

type OAuth2Service interface {
	UserAccountCreating(userCreatingReq *_userModel.UserCreatingReq) error
	IsUserExists(userID string) bool
}
