package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/helper"
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_orderRepository "github.com/Kamila3820/hoca-backend/modules/order/repository"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
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
		ContactName:   orderCreatingReq.ContactName,
		ContactPhone:  orderCreatingReq.ContactPhone,
		PaymentType:   orderCreatingReq.PaymentType,
		SpecificPlace: orderCreatingReq.SpecificPlace,
		Note:          orderCreatingReq.Note,
		Price:         workerPost.Price,
		OrderStatus:   "confirmation",
	}

	if workerPost.OwnerID == orderCreatingReq.UserID {
		return nil, errors.New("Can not proceed yourself post")
	}

	if workerPost.IsReserved {
		return nil, errors.New("Can not order a post that running the progress")
	}

	workerPost.IsReserved = true
	s.orderRepository.UpdatePost(workerPost)

	order, err := s.orderRepository.CreatingOrder(orderEntity)
	if err != nil {
		return nil, err
	}

	return order.ToOrderModel(), nil
}

func (s *orderServiceImpl) GetUserByID(userID string) (*_userModel.User, error) {
	user, err := s.orderRepository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Location == "" {
		return nil, errors.New("user location not set")
	}

	return user.ToUserModel(), nil
}

func (s *orderServiceImpl) UpdateOrderProgress(updaterID string, orderID uint64, newStatus string) (*_orderModel.Order, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("service: cannot find order by id")
	}

	workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
	if err != nil {
		return nil, errors.New("service: cannot find worker post by id")
	}

	if updaterID != workerPost.OwnerID {
		return nil, errors.New("service: you have no permission to update the order")
	}

	if order.IsCancel {
		return nil, fmt.Errorf("order has been cancelled")
	}

	confirmDeadline := order.CreatedAt.Add(7 * time.Minute)

	if order.OrderStatus == "confirmation" && time.Now().After(confirmDeadline) {
		order.OrderStatus = "cancelled"
		order.IsCancel = true
		order.CancellationReason = "worker did not confirm within the time"
		order.CancelledBy = "worker"
		s.orderRepository.UpdateOrder(order)
		return nil, fmt.Errorf("order automatically cancelled due to no confirmation")
	}

	order.OrderStatus = newStatus
	s.orderRepository.UpdateOrder(order)

	if newStatus == "complete" {
		// Update worker post
		workerPost.IsReserved = false
		s.orderRepository.UpdatePost(workerPost)

		// Create history
		history := &entities.History{
			UserID:    updaterID,
			OrderID:   strconv.Itoa(int(order.ID)),
			Status:    "complete",
			IsRated:   false,
			CreatedAt: time.Now(),
		}
		_, err = s.orderRepository.CreatingHistory(history)
		if err != nil {
			return nil, errors.New("cannot create order history")
		}
	}

	return order.ToOrderModel(), nil
}

func (s *orderServiceImpl) CancelOrder(orderID uint64, reason string, cancelledBy string) error {
	//Update order
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return errors.New("cannot find order by id")
	}

	if order.OrderStatus != "confirmation" || order.IsCancel {
		return errors.New("order cannot be cancelled")
	}

	order.OrderStatus = "cancelled"
	order.IsCancel = true
	order.CancellationReason = reason
	order.CancelledBy = cancelledBy

	if err := s.orderRepository.UpdateOrder(order); err != nil {
		return errors.New("order cannot be update")
	}

	// Update post
	workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
	if err != nil {
		return errors.New("cannot find post by postID")
	}

	workerPost.IsReserved = false
	s.orderRepository.UpdatePost(workerPost)

	// Create history
	history := &entities.History{
		UserID:             "",
		OrderID:            strconv.Itoa(int(order.ID)),
		Status:             "cancelled",
		CancellationReason: reason,
		CancelledBy:        cancelledBy,
		IsRated:            false,
		CreatedAt:          time.Now(),
	}
	_, err = s.orderRepository.CreatingHistory(history)
	if err != nil {
		return errors.New("cannot create order history")
	}

	return nil
}

func (s *orderServiceImpl) StartConfirmationTimer(orderID uint64) {
	go func() {
		time.Sleep(7 * time.Minute)

		order, err := s.orderRepository.FindOrderByID(orderID)
		if err != nil {
			fmt.Println("Failed to find order", err)
			return
		}

		if order.OrderStatus == "confirmation" {
			// Auto cancel if the order is still in the "confirmation" after 7 minutes
			order.OrderStatus = "cancelled"
			order.IsCancel = true
			order.CancelledBy = "system"
			order.CancellationReason = "worker did not confirm within the time"
			s.orderRepository.UpdateOrder(order)

			workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
			if err != nil {
				fmt.Println("Failed to find post", err)
				return
			}
			workerPost.IsReserved = false

			s.orderRepository.UpdatePost(workerPost)
		}
	}()
}

func (s *orderServiceImpl) GetPreparingOrder(orderID uint64) (*_orderModel.Order, *_orderModel.DistanceMatrixResponse, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, nil, errors.New("cannot find order by id")
	}

	if order.OrderStatus != "preparing" {
		return order.ToOrderModel(), nil, nil
	}

	distanceResponse, err := helper.DistanceMatrix(order.User.Latitude+","+order.User.Longtitude, order.Post.LocationLat+","+order.Post.LocationLong)
	if err != nil {
		return nil, nil, errors.New("failed to get distance information")
	}

	return order.ToOrderModel(), (*_orderModel.DistanceMatrixResponse)(distanceResponse), nil
}
