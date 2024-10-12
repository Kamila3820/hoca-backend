package server

import (
	_postController "github.com/Kamila3820/hoca-backend/modules/post/controller"
	_postRepository "github.com/Kamila3820/hoca-backend/modules/post/repository"
	_postService "github.com/Kamila3820/hoca-backend/modules/post/service"
)

func (s *echoServer) initPostRouter() {
	router := s.app.Group("/v1/post", Jwt())

	postRepository := _postRepository.NewPostRepositoryImpl(s.db, s.app.Logger)
	postService := _postService.NewPostServiceImpl(postRepository)
	postController := _postController.NewPostControllerImpl(postService)

	router.GET("/list", postController.FindPostByDistance)
	router.GET("/me", postController.GetOwnPost)
	router.GET("/:postID", postController.GetPostByPostID)
	router.POST("/create", postController.CreateWorkerPost)
	router.PATCH("/edit/:postID", postController.EditWorkerPost)
	router.DELETE("/delete/:postID", postController.DeleteWorkerPost)

	router.DELETE("/open/:postID", postController.ActivatePost)
	router.DELETE("/close/:postID", postController.UnActivatePost)
}
