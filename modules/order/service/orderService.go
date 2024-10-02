package service

import (
	"github.com/Kamila3820/hoca-backend/helper"
	_paymentModel "github.com/Kamila3820/hoca-backend/helper/model"
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
)

type OrderService interface {
	CreatingOrder(orderCreatingReq *_orderModel.OrderReq, postID uint64) (*_orderModel.Order, error)
	UpdateOrderProgress(updaterID string, orderID uint64, newStatus string) (*_orderModel.Order, error)
	CancelOrder(orderID uint64, reason string, cancelledBy string) error
	StartConfirmationTimer(orderID uint64)

	GetUserOrder(orderID uint64, userID string) (*_orderModel.UserOrder, error)
	GetWorkerOrder(orderID uint64, userID string) (*_orderModel.WorkerOrder, error)
	GetPreparingOrder(orderID uint64, customerLat, customerLong string) (*_orderModel.Order, *helper.DirectionsResponse, error)
	GetQRpayment(userID string, orderID uint64) (*_paymentModel.CreateOrderQrResponse, error)
	InquiryQRpayment(transactionID string) (*_paymentModel.PaymentInquiryResponse, error)

	GetUserByID(userID string) (*_userModel.User, error)
}
