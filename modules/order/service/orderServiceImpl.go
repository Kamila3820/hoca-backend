package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/helper"
	_paymentModel "github.com/Kamila3820/hoca-backend/helper/model"
	_notiModel "github.com/Kamila3820/hoca-backend/modules/notification/model"
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_orderRepository "github.com/Kamila3820/hoca-backend/modules/order/repository"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
	"github.com/Kamila3820/hoca-backend/utils/text"
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

	order, err := s.orderRepository.CreatingOrder(orderEntity)
	if err != nil {
		return nil, err
	}

	go func() {
		workerPost.IsReserved = true
		s.orderRepository.UpdatePost(workerPost)

		// Create noti
		orderType := _notiModel.NotificationPlaceOrder
		notification := &entities.Notification{
			Trigger:          nil,
			TriggerID:        &order.UserID,
			Triggee:          nil,
			TriggeeID:        &workerPost.OwnerID,
			Order:            nil,
			OrderID:          &order.ID,
			UserRating:       nil,
			UserRatingID:     nil,
			NotificationType: &orderType,
			CreatedAt:        nil,
		}

		if err := s.orderRepository.CreateNotification(notification); err != nil {
			fmt.Printf("service: unable to create notification %v", err.Error)
		}
	}()

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
		order.CancelledBy = "system"
		s.orderRepository.UpdateOrder(order)
	}

	order.OrderStatus = newStatus
	s.orderRepository.UpdateOrder(order)

	if newStatus == "preparing" {
		go func() {
			// Create noti
			orderType := _notiModel.NotificationPreparing
			notification := &entities.Notification{
				Trigger:          nil,
				TriggerID:        &workerPost.OwnerID,
				Triggee:          nil,
				TriggeeID:        &order.UserID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &orderType,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notification); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}
		}()
	}

	if newStatus == "working" {
		go func() {
			// Create noti
			orderType := _notiModel.NotificationWorking
			notification := &entities.Notification{
				Trigger:          nil,
				TriggerID:        &workerPost.OwnerID,
				Triggee:          nil,
				TriggeeID:        &order.UserID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &orderType,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notification); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}
		}()
	}

	if newStatus == "complete" {
		// Update worker post
		workerPost.IsReserved = false
		s.orderRepository.UpdatePost(workerPost)

		// Update order
		order.Paid = true
		s.orderRepository.UpdateOrder(order)

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

		go func() {
			// Create noti
			orderType := _notiModel.NotificationComplete
			notification := &entities.Notification{
				Trigger:          nil,
				TriggerID:        &workerPost.OwnerID,
				Triggee:          nil,
				TriggeeID:        &order.UserID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &orderType,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notification); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}
		}()
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

	go func() {
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
			fmt.Printf("service: unable to create history %v", err.Error)
		}

		// Create noti
		var notification *entities.Notification

		if cancelledBy == "customer" {
			cancelType := _notiModel.NotificationUserCancel
			notification = &entities.Notification{
				Trigger:          nil,
				TriggerID:        &order.UserID,
				Triggee:          nil,
				TriggeeID:        &order.Post.OwnerID,
				Order:            nil,
				OrderID:          &orderID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &cancelType,
				CreatedAt:        nil,
			}
		}
		if cancelledBy == "worker" {
			cancelType := _notiModel.NotificationWorkerCancel
			notification = &entities.Notification{
				Trigger:          nil,
				TriggerID:        &order.Post.OwnerID,
				Triggee:          nil,
				TriggeeID:        &order.UserID,
				Order:            nil,
				OrderID:          &orderID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &cancelType,
				CreatedAt:        nil,
			}
		}

		if err := s.orderRepository.CreateNotification(notification); err != nil {
			fmt.Printf("service: unable to create notification %v", err.Error)
		}
	}()

	return nil
}

func (s *orderServiceImpl) StartConfirmationTimer(orderID uint64) {
	go func() {
		time.Sleep(7 * time.Minute)

		order, err := s.orderRepository.FindOrderByID(orderID)
		if err != nil {
			fmt.Println("service: failed to find order", err)
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
				fmt.Println("service: failed to find post", err)
				return
			}
			workerPost.IsReserved = false

			s.orderRepository.UpdatePost(workerPost)

			// Create noti to customer
			cancelTypeCustomer := _notiModel.NotificationPlaceOrder
			notification := &entities.Notification{
				Trigger:          nil,
				TriggerID:        nil,
				Triggee:          nil,
				TriggeeID:        &order.UserID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &cancelTypeCustomer,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notification); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}

			// Create noti to worker
			cancelTypeWorker := _notiModel.NotificationPlaceOrder
			notificationWorker := &entities.Notification{
				Trigger:          nil,
				TriggerID:        nil,
				Triggee:          nil,
				TriggeeID:        &workerPost.OwnerID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &cancelTypeWorker,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notificationWorker); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}
		}

	}()
}

func (s *orderServiceImpl) GetPreparingOrder(orderID uint64, customerLat, customerLong string) (*_orderModel.Order, *helper.DirectionsResponse, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, nil, errors.New("service: cannot find order by id")
	}

	if order.OrderStatus != "preparing" {
		return order.ToOrderModel(), nil, nil
	}

	workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
	if err != nil {
		return nil, nil, errors.New("service: order cannot be update")
	}

	workerLat := workerPost.LocationLat
	workerLong := workerPost.LocationLong

	workerLocation := fmt.Sprintf("%s,%s", workerLat, workerLong)
	customerLocation := fmt.Sprintf("%s,%s", customerLat, customerLong)

	directionsRequest := &helper.DirectionsRequest{
		Origin:      workerLocation,
		Destination: customerLocation,
		Mode:        "driving",
	}

	client := helper.Client{APIKey: config.ConfigGetting().Google.ApiKey}
	directionsResponse, err := client.Directions(directionsRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to get directions: %v", err)
	}

	return order.ToOrderModel(), directionsResponse, nil
}

func (s *orderServiceImpl) GetQRpayment(userID string, orderID uint64) (*_paymentModel.CreateOrderQrResponse, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("service: cannot find order by id")
	}

	if order.OrderStatus != "working" {
		return nil, errors.New("service: cannot create QR payment in this stage")
	}

	if order.PaymentType != "qrcode" {
		return nil, errors.New("service: cannot proceed QR payment type in this order")
	}

	transactionID := text.GenerateTransactionId(10)

	paymentEntity := &entities.OrderQrpayment{
		User:          nil,
		UserId:        userID,
		Order:         nil,
		OrderID:       orderID,
		Amount:        uint64(order.Price),
		Paid:          false,
		TransactionID: transactionID,
	}

	if err := s.orderRepository.CreatingQRpayment(paymentEntity); err != nil {
		return nil, errors.New("service: cannot create order qr payment")
	}

	// Create qr code
	qrData := helper.ScbCreateQrPayment(uint(order.Price), transactionID)

	paymentResponse := &_paymentModel.CreateOrderQrResponse{
		QrRawData:     qrData.QrRawData,
		QrImage:       qrData.QrImage,
		TransactionId: transactionID,
	}

	return paymentResponse, nil
}

func (s *orderServiceImpl) InquiryQRpayment(transactionID string) (*_paymentModel.PaymentInquiryResponse, error) {
	inquiryData, err := helper.ScbInquiryPayment(transactionID)
	if err != nil {
		return nil, errors.New("service: unable to inquiry payment")
	}

	paymentStatusResponse := new(_paymentModel.PaymentInquiryResponse)

	if inquiryData != nil && *inquiryData.PayeeName != "" && inquiryData.PayeeName != nil {
		paymentStatusResponse.PaymentSuccess = true

		paymentOrder, err := s.orderRepository.FindTransactionByID(transactionID)
		if err != nil {
			return nil, errors.New("service: cannot find payment order by ID")
		}

		// Update Paid in the order
		paymentOrder.Paid = true
		if err := s.orderRepository.UpdateTransactionOrder(paymentOrder); err != nil {
			return nil, errors.New("service: cannot update payment order")
		}

	} else {
		paymentStatusResponse.PaymentSuccess = false
		paymentStatusResponse.Message = "Payment Not Found"
	}

	return paymentStatusResponse, nil
}

func (s *orderServiceImpl) GetUserOrder(orderID uint64, userID string) (*_orderModel.UserOrder, error) {
	userOrder, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("service: cannot query order by id")
	}

	if userID != userOrder.UserID {
		return nil, errors.New("service: cannot proceed the user order")
	}

	return userOrder.ToUserOrder(), nil
}

func (s *orderServiceImpl) GetWorkerOrder(orderID uint64, userID string) (*_orderModel.WorkerOrder, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("service: cannot query order by id")
	}

	workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
	if err != nil {
		return nil, errors.New("service: cannot query post by id")
	}

	if userID != workerPost.OwnerID {
		return nil, errors.New("service: cannot proceed the worker order")
	}

	return order.ToWorkerOrder(), nil
}
