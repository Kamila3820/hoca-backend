package server

import (
	_userController "github.com/Kamila3820/hoca-backend/modules/user/controller"
	_userRepository "github.com/Kamila3820/hoca-backend/modules/user/repository"
	_userService "github.com/Kamila3820/hoca-backend/modules/user/service"
)

func (s *echoServer) initUserRouter() {
	router := s.app.Group("/v1/user", Jwt())

	userRepository := _userRepository.NewUserRepositoryImpl(s.db, s.app.Logger)
	userService := _userService.NewUserServiceImpl(userRepository)
	userController := _userController.NewUserControllerImpl(userService)

	router.GET("/profile", userController.GetUserByID)
	router.PATCH("/profile/edit", userController.EditUserProfile)
	router.GET("/location", userController.GetUserLocation)
	router.PATCH("/location/edit", userController.UpdateUserLocation)
}
