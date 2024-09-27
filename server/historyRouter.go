package server

import (
	_historyController "github.com/Kamila3820/hoca-backend/modules/history/controller"
	_historyRepository "github.com/Kamila3820/hoca-backend/modules/history/repository"
	_historyService "github.com/Kamila3820/hoca-backend/modules/history/service"
)

func (s *echoServer) initHistoryRouter(m *authorizingMiddleware) {
	router := s.app.Group("/v1/history")

	historyRepository := _historyRepository.NewHistoryRepositoryImpl(s.db, s.app.Logger)
	historyService := _historyService.NewHistoryServiceImpl(historyRepository)
	historyController := _historyController.NewHistoryControllerImpl(historyService)

	router.GET("/list", historyController.GetHistoryByUserID, m.UserAuthorizing)
	router.GET("/work", historyController.GetWorkingHistory, m.UserAuthorizing)
}
