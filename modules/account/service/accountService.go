package service

import (
	_accountModel "github.com/Kamila3820/hoca-backend/modules/account/model"
)

type AccountService interface {
	Register(registerReq *_accountModel.RegisterRequest) (*_accountModel.RegisterResponse, error)
	Login(loginReq *_accountModel.LoginRequest) (*_accountModel.LoginResponse, error)
}
