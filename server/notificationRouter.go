package server

import (
	_notificationController "github.com/Kamila3820/hoca-backend/modules/notification/controller"
	_notificationRepository "github.com/Kamila3820/hoca-backend/modules/notification/repository"
	_notificationService "github.com/Kamila3820/hoca-backend/modules/notification/service"
)

func (s *echoServer) initNotificationRouter(m *authorizingMiddleware) {
	router := s.app.Group("/v1/notification")

	notificationRepository := _notificationRepository.NewNotificationRepositoryImpl(s.db, s.app.Logger)
	notificationService := _notificationService.NewNotificationServiceImpl(notificationRepository)
	notificationController := _notificationController.NewNotificationControllerImpl(notificationService)

	router.GET("/", notificationController.GetNotificationsByUser, m.UserAuthorizing)
}
