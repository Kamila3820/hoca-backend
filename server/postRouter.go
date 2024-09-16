package server

import (
	_postController "github.com/Kamila3820/hoca-backend/modules/post/controller"
	_postRepository "github.com/Kamila3820/hoca-backend/modules/post/repository"
	_postService "github.com/Kamila3820/hoca-backend/modules/post/service"
)

func (s *echoServer) initPostRouter(m *authorizingMiddleware) {
	router := s.app.Group("/v1/post")

	postRepository := _postRepository.NewPostRepositoryImpl(s.db, s.app.Logger)
	postService := _postService.NewPostServiceImpl(postRepository)
	postController := _postController.NewPostControllerImpl(postService)

	router.GET("/list", postController.FindPostByDistance, m.UserAuthorizing)
	router.POST("/create", postController.CreateWorkerPost, m.UserAuthorizing)
	router.PATCH("/edit/:postID", postController.EditWorkerPost, m.UserAuthorizing)
	router.DELETE("/delete/:postID", postController.DeleteWorkerPost, m.UserAuthorizing)

	router.DELETE("/open/:postID", postController.ActivatePost, m.UserAuthorizing)
	router.DELETE("/close/:postID", postController.UnActivatePost, m.UserAuthorizing)
}
