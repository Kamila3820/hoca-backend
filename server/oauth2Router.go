package server

import (
	_oauth2Controller "github.com/Kamila3820/hoca-backend/modules/oauth2/controller"
	_oauth2Service "github.com/Kamila3820/hoca-backend/modules/oauth2/service"
	_userRepository "github.com/Kamila3820/hoca-backend/modules/user/repository"
)

func (s *echoServer) initOAuth2Router() {
	router := s.app.Group("/v1/oauth2/google")

	userRepository := _userRepository.NewUserRepositoryImpl(s.db, s.app.Logger)

	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(userRepository)
	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(oauth2Service, s.conf.OAuth2, s.app.Logger)

	router.GET("/user/login", oauth2Controller.UserLogin)
	router.GET("/user/login/callback", oauth2Controller.UserLoginCallback)
	router.DELETE("/logout", oauth2Controller.Logout)

}
