package service

import (
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
)

type OrderService interface {
	CreatingOrder(orderCreatingReq *_orderModel.OrderReq, postID uint64) (*_orderModel.Order, error)
}
