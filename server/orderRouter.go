package server

import (
	_orderController "github.com/Kamila3820/hoca-backend/modules/order/controller"
	_orderRepository "github.com/Kamila3820/hoca-backend/modules/order/repository"
	_orderService "github.com/Kamila3820/hoca-backend/modules/order/service"
)

func (s *echoServer) initOrderRouter(m *authorizingMiddleware) {
	router := s.app.Group("/v1/order")

	orderRepository := _orderRepository.NewOrderRepositoryImpl(s.db, s.app.Logger)
	orderService := _orderService.NewOrderServiceImpl(orderRepository)
	orderController := _orderController.NewOrderControllerImpl(orderService)

	router.POST("/create/:postID", orderController.PlaceOrder, m.UserAuthorizing)
	router.GET("/contact", orderController.GetUserContact, m.UserAuthorizing)
	router.PATCH("/update/:orderID", orderController.WorkerUpdateProgress, m.UserAuthorizing)
	router.PATCH("/cancel/:orderID", orderController.CancelOrder, m.UserAuthorizing)
	router.GET("/timer/:orderID", orderController.ConfirmationTimerOrder, m.UserAuthorizing)

	router.GET("/prepare/:orderID", orderController.GetPreparingOrder, m.UserAuthorizing)
}
