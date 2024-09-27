package service

import (
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
)

type OrderService interface {
	CreatingOrder(orderCreatingReq *_orderModel.OrderReq, postID uint64) (*_orderModel.Order, error)
	UpdateOrderProgress(updaterID string, orderID uint64, newStatus string) (*_orderModel.Order, error)
	CancelOrder(orderID uint64, reason string, cancelledBy string) error
	StartConfirmationTimer(orderID uint64)
	GetPreparingOrder(orderID uint64) (*_orderModel.Order, *_orderModel.DistanceMatrixResponse, error)

	GetUserByID(userID string) (*_userModel.User, error)
}
