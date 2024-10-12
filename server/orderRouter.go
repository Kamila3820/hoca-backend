package server

import (
	_orderController "github.com/Kamila3820/hoca-backend/modules/order/controller"
	_orderRepository "github.com/Kamila3820/hoca-backend/modules/order/repository"
	_orderService "github.com/Kamila3820/hoca-backend/modules/order/service"
)

func (s *echoServer) initOrderRouter() {
	router := s.app.Group("/v1/order", Jwt())

	orderRepository := _orderRepository.NewOrderRepositoryImpl(s.db, s.app.Logger)
	orderService := _orderService.NewOrderServiceImpl(orderRepository)
	orderController := _orderController.NewOrderControllerImpl(orderService)

	router.POST("/create/:postID", orderController.PlaceOrder)
	router.GET("/contact", orderController.GetUserContact)
	router.PATCH("/update/:orderID", orderController.WorkerUpdateProgress)
	router.PATCH("/cancel/:orderID", orderController.CancelOrder)

	router.GET("/timer/:orderID", orderController.ConfirmationTimerOrder)
	router.GET("/user/:orderID", orderController.GetUserOrder)
	router.GET("/worker/:orderID", orderController.GetWorkerOrder)

	router.GET("/prepare/:orderID", orderController.GetPreparingOrder)
	router.GET("/payment/qr/:orderID", orderController.GetQRpayment)
	router.GET("/payment/inquiry", orderController.InquiryQRpayment)
}
