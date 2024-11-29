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

	activeOrder, err := s.orderRepository.FindActiveOrder(orderCreatingReq.UserID)
	if err != nil {
		return nil, err
	}

	if activeOrder != nil {
		return nil, errors.New("Can not proceed order, if you still have the running order")
	}

	workActiveOrder, err := s.orderRepository.FindWorkerOrder(orderCreatingReq.UserID)
	if err != nil {
		return nil, err
	}

	if workActiveOrder != nil {
		return nil, errors.New("Your post still running the process, please finish it first")
	}

	orderEntity := &entities.Order{
		UserID:        orderCreatingReq.UserID,
		WorkerPostID:  postID,
		ContactName:   orderCreatingReq.ContactName,
		ContactPhone:  orderCreatingReq.ContactPhone,
		PaymentType:   orderCreatingReq.PaymentType,
		SpecificPlace: orderCreatingReq.SpecificPlace,
		Note:          orderCreatingReq.Note,
		Duration:      orderCreatingReq.Duration,
		Price:         orderCreatingReq.Price,
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

	// confirmDeadline := order.CreatedAt.Add(4 * time.Minute)

	// if order.OrderStatus == "confirmation" && time.Now().After(confirmDeadline) {
	// 	order.OrderStatus = "cancelled"
	// 	order.IsCancel = true
	// 	order.CancellationReason = "worker did not confirm within the time"
	// 	order.CancelledBy = "system"
	// 	s.orderRepository.UpdateOrder(order)
	// }

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
		time.Sleep(4 * time.Minute)

		order, err := s.orderRepository.FindOrderByID(orderID)
		if err != nil {
			fmt.Println("service: failed to find order", err)
			return
		}

		if order.OrderStatus == "confirmation" {
			triggerSystem0 := "afbf87a5-9114-42a9-bb58-82ab18809ecd"
			triggerSystem1 := "bfbf87a5-9114-42a9-bb58-82ab18809ecd"
			// Auto cancel if the order is still in the "confirmation" after 4 minutes
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
			cancelTypeSystem := _notiModel.NotificationSystemCancel
			notification := &entities.Notification{
				Trigger:          nil,
				TriggerID:        &triggerSystem0,
				Triggee:          nil,
				TriggeeID:        &order.UserID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &cancelTypeSystem,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notification); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}

			// Create noti to worker
			notificationWorker := &entities.Notification{
				Trigger:          nil,
				TriggerID:        &triggerSystem1,
				Triggee:          nil,
				TriggeeID:        &workerPost.OwnerID,
				Order:            nil,
				OrderID:          &order.ID,
				UserRating:       nil,
				UserRatingID:     nil,
				NotificationType: &cancelTypeSystem,
				CreatedAt:        nil,
			}

			if err := s.orderRepository.CreateNotification(notificationWorker); err != nil {
				fmt.Printf("service: unable to create notification %v", err.Error)
			}
		}

	}()
}

func (s *orderServiceImpl) GetPreparingOrder(orderID uint64) (*_orderModel.UserOrder, *helper.DirectionsResponse, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, nil, errors.New("service: cannot find order by id")
	}

	if order.OrderStatus != "preparing" {
		return order.ToUserOrder(), nil, nil
	}

	user, err := s.orderRepository.FindUserByID(order.UserID)
	if err != nil {
		return nil, nil, errors.New("service: cannot find user by id")
	}

	workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
	if err != nil {
		return nil, nil, errors.New("service: order cannot be update")
	}

	workerLat := workerPost.LocationLat
	workerLong := workerPost.LocationLong

	workerLocation := fmt.Sprintf("%s,%s", workerLat, workerLong)
	customerLocation := fmt.Sprintf("%s,%s", user.Latitude, user.Longtitude)

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

	return order.ToUserOrder(), directionsResponse, nil
}

func (s *orderServiceImpl) GetWorkerPrepare(orderID uint64) (*_orderModel.WorkerOrder, *helper.DirectionsResponse, error) {
	order, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, nil, errors.New("service: cannot find order by id")
	}

	if order.OrderStatus != "preparing" {
		return order.ToWorkerOrder(), nil, nil
	}

	user, err := s.orderRepository.FindUserByID(order.UserID)
	if err != nil {
		return nil, nil, errors.New("service: cannot find user by id")
	}

	workerPost, err := s.orderRepository.FindPostByID(order.WorkerPostID)
	if err != nil {
		return nil, nil, errors.New("service: order cannot be update")
	}

	workerLat := workerPost.LocationLat
	workerLong := workerPost.LocationLong

	workerLocation := fmt.Sprintf("%s,%s", workerLat, workerLong)
	customerLocation := fmt.Sprintf("%s,%s", user.Latitude, user.Longtitude)

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

	return order.ToWorkerOrder(), directionsResponse, nil
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

func (s *orderServiceImpl) GetWorkerFeePayment(postID uint64) (*_paymentModel.CreateWorkerFeeQrResponse, error) {
	weekOrder, err := s.orderRepository.FindLastWeekOrderByPostID(postID)
	if err != nil {
		return nil, errors.New("service: cannot find last week order by postID")
	}

	if weekOrder == nil {
		return nil, nil
	}

	postCheck, err := s.orderRepository.FindPostByID(postID)
	if err != nil {
		return nil, errors.New("service: cannot find post by postID")
	}

	if postCheck.ID != postID {
		return nil, errors.New("service: invalid postID")
	}

	now := time.Now()
	startFrom := now.AddDate(0, 0, -7)
	endedAt := now.AddDate(0, 0, 5)
	orderCount := len(weekOrder)
	serviceFee := 70 * (len(weekOrder))
	transactionID := text.GenerateTransactionId(10)

	paymentEntity := &entities.WorkerFeepayment{
		Post:          nil,
		PostId:        postID,
		OrderCount:    uint64(orderCount),
		Amount:        uint64(serviceFee),
		TransactionID: transactionID,
		Paid:          false,
		StartFrom:     startFrom.Format("2006-01-02 15:04"),
		EndFrom:       now.Format("2006-01-02 15:04"),
	}

	if err := s.orderRepository.CreatingWorkerFeePayment(paymentEntity); err != nil {
		return nil, errors.New("service: cannot create worker fee payment")
	}

	// Create qr code
	qrData := helper.ScbCreateQrPayment(uint(serviceFee), transactionID)

	fmt.Println(serviceFee)
	fmt.Println(orderCount)
	fmt.Println(weekOrder)
	paymentResponse := &_paymentModel.CreateWorkerFeeQrResponse{
		QrRawData:     qrData.QrRawData,
		QrImage:       qrData.QrImage,
		TransactionId: transactionID,
		OrderCount:    orderCount,
		Amount:        uint64(serviceFee),
		StartFrom:     startFrom.Format("2006-01-02"),
		EndFrom:       now.Format("2006-01-02"),
		EndedAt:       endedAt.Format("2006-01-02"),
	}

	return paymentResponse, nil
}

func (s *orderServiceImpl) InquiryWorkerFeePayment(transactionID string) (*_paymentModel.PaymentInquiryResponse, error) {
	inquiryData, err := helper.ScbInquiryPayment(transactionID)
	if err != nil {
		return nil, errors.New("service: unable to inquiry payment")
	}

	paymentStatusResponse := new(_paymentModel.PaymentInquiryResponse)

	if inquiryData != nil && *inquiryData.PayeeName != "" && inquiryData.PayeeName != nil {
		paymentStatusResponse.PaymentSuccess = true

		paymentOrder, err := s.orderRepository.FindWorkerTransactionByID(transactionID)
		if err != nil {
			return nil, errors.New("service: cannot find payment order by ID")
		}

		// Update Paid in the order
		paymentOrder.Paid = true
		if err := s.orderRepository.UpdateTransactionFee(paymentOrder); err != nil {
			return nil, errors.New("service: cannot update payment order")
		}
	} else {
		paymentStatusResponse.PaymentSuccess = false
		paymentStatusResponse.Message = "Payment Not Found"
	}

	return paymentStatusResponse, nil
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

		order, err := s.orderRepository.FindOrderByID(paymentOrder.OrderID)
		if err != nil {
			return nil, errors.New("service: cannot find order by ID")
		}

		order.Paid = true
		s.orderRepository.UpdateOrder(order)

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

func (s *orderServiceImpl) GetActiveOrder(userID string) (*_orderModel.UserOrder, error) {
	order, err := s.orderRepository.FindActiveOrder(userID)
	if err != nil {
		return nil, errors.New("service: cannot query active worker order")
	}

	if order == nil {
		return nil, nil
	}

	return order.ToUserOrder(), nil
}

func (s *orderServiceImpl) GetWorkerActiveOrder(userID string) (*_orderModel.WorkerOrder, error) {
	order, err := s.orderRepository.FindWorkerOrder(userID)
	if err != nil {
		return nil, errors.New("service: cannot query active order")
	}

	if order == nil {
		return nil, nil
	}

	return order.ToWorkerOrder(), nil
}
