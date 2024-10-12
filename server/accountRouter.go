package server

import (
	_accountController "github.com/Kamila3820/hoca-backend/modules/account/controller"
	_accountRepository "github.com/Kamila3820/hoca-backend/modules/account/repository"
	_accountService "github.com/Kamila3820/hoca-backend/modules/account/service"
)

func (s *echoServer) initAccountRouter() {
	router := s.app.Group("/v1/account")

	accountRepository := _accountRepository.NewAccountRepositoryImpl(s.db, s.app.Logger)
	accountService := _accountService.NewAccountServiceImpl(accountRepository)
	accountController := _accountController.NewAccountControllerImpl(accountService)

	router.POST("/register", accountController.Register)
	router.POST("/login", accountController.Login)
}
