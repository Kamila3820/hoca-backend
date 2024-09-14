package server

import (
	_userRatingController "github.com/Kamila3820/hoca-backend/modules/user_rating/controller"
	_userRatingRepository "github.com/Kamila3820/hoca-backend/modules/user_rating/repository"
	_userRatingService "github.com/Kamila3820/hoca-backend/modules/user_rating/service"
)

func (s *echoServer) initUserRatingRouter(m *authorizingMiddleware) {
	router := s.app.Group("/v1/rating")

	userRatingRepository := _userRatingRepository.NewUserRatingRepositoryImpl(s.db, s.app.Logger)
	userRatingService := _userRatingService.NewUserRatingServiceImpl(userRatingRepository)
	userRatingController := _userRatingController.NewUserRatingControllerImpl(userRatingService)

	router.GET("/list/:postID", userRatingController.ListRatingByPostID, m.UserAuthorizing)
	router.POST("/create/:postID", userRatingController.RatingWorker, m.UserAuthorizing)
}
