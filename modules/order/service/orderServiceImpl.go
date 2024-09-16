package service

import (
	"github.com/Kamila3820/hoca-backend/entities"
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_orderRepository "github.com/Kamila3820/hoca-backend/modules/order/repository"
)

type orderServiceImpl struct {
	orderRepository _orderRepository.OrderRepository
}

func NewOrderServiceImpl(orderRepository _orderRepository.OrderRepository) OrderService {
	return &orderServiceImpl{
		orderRepository: orderRepository,
	}
}

func (s *orderServiceImpl) CreatingOrder(orderCreatingReq *_orderModel.OrderReq, postID uint64) (*_orderModel.Order, error) {
	workerPost, err := s.orderRepository.FindPostByID(postID)
	if err != nil {
		return nil, err
	}

	orderEntity := &entities.Order{
		UserID:        orderCreatingReq.UserID,
		WorkerPostID:  postID,
		PaymentType:   orderCreatingReq.PaymentType,
		SpecificPlace: orderCreatingReq.SpecificPlace,
		Note:          orderCreatingReq.Note,
		Price:         workerPost.Price,
		OrderStatus:   "confirmation",
	}

	if workerPost.OwnerID == orderCreatingReq.UserID {
		return nil, err
	}

	order, err := s.orderRepository.CreatingOrder(orderEntity)
	if err != nil {
		return nil, err
	}

	return order.ToOrderModel(), nil
}
