package controller

import (
	//_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"

	"fmt"
	"net/http"
	"strconv"

	_paymentModel "github.com/Kamila3820/hoca-backend/helper/model"
	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_orderService "github.com/Kamila3820/hoca-backend/modules/order/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type orderControllerImpl struct {
	orderService _orderService.OrderService
}

func NewOrderControllerImpl(orderService _orderService.OrderService) OrderController {
	return &orderControllerImpl{
		orderService: orderService,
	}
}

func (c *orderControllerImpl) PlaceOrder(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	fmt.Println("1")
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	orderCreatingReq := new(_orderModel.OrderReq)
	fmt.Println("2")
	fmt.Println(orderCreatingReq.ContactName)
	fmt.Println(orderCreatingReq.ContactPhone)
	fmt.Println(orderCreatingReq.PaymentType)
	fmt.Println(orderCreatingReq.SpecificPlace)
	fmt.Println(orderCreatingReq.Note)
	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(orderCreatingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	fmt.Println("3")
	fmt.Println(orderCreatingReq.ContactName)
	fmt.Println(orderCreatingReq.ContactPhone)
	fmt.Println(orderCreatingReq.PaymentType)
	fmt.Println(orderCreatingReq.SpecificPlace)
	fmt.Println(orderCreatingReq.Note)
	orderCreatingReq.UserID = userID.ID

	newOrder, err := c.orderService.CreatingOrder(orderCreatingReq, postID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	fmt.Println("4")
	return pctx.JSON(http.StatusCreated, newOrder)
}

func (c *orderControllerImpl) getPostID(pctx echo.Context) (uint64, error) {
	postIDStr := pctx.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return postID, nil
}

func (c *orderControllerImpl) getOrderID(pctx echo.Context) (uint64, error) {
	orderIDStr := pctx.Param("orderID")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return orderID, nil
}

func (c *orderControllerImpl) GetUserContact(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	user, err := c.orderService.GetUserByID(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusNotFound, err)
	}

	return pctx.JSON(http.StatusOK, map[string]interface{}{
		"username":  user.UserName,
		"phone":     user.PhoneNumber,
		"location":  user.Location,
		"latitude":  user.Latitude,
		"longitude": user.Longtitude,
	})
}

func (c *orderControllerImpl) WorkerUpdateProgress(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	statusUpdate := new(struct {
		Status string `json:"status" validate:"required,oneof=confirmation preparing working complete"`
	})

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(statusUpdate); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	order, err := c.orderService.UpdateOrderProgress(userID.ID, orderID, statusUpdate.Status)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, order)
}

func (c *orderControllerImpl) CancelOrder(pctx echo.Context) error {
	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	cancelBy := pctx.QueryParam("cancelBy")

	cancelOrderReq := new(_orderModel.CancelOrderReq)
	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(cancelOrderReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.orderService.CancelOrder(orderID, cancelOrderReq.CancellationReason, cancelBy); err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, map[string]string{
		"message": "Order canceled successfully",
	})
}

func (c *orderControllerImpl) ConfirmationTimerOrder(pctx echo.Context) error {
	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	c.orderService.StartConfirmationTimer(orderID)

	return pctx.JSON(http.StatusOK, map[string]string{
		"message": "Timer Starting, please wait for the worker to confirm your order",
	})
}

func (c *orderControllerImpl) GetPreparingOrder(pctx echo.Context) error {
	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	order, distance, err := c.orderService.GetPreparingOrder(orderID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"worker_name":   order.WorkerName,
		"worker_phone":  order.WorkerPhone,
		"worker_avatar": order.WorkerAvatar,
		"payment_type":  order.PaymentType,
		"price":         order.Price,
		"order_status":  order.OrderStatus,
	}

	if distance != nil {
		response["distance"] = distance.Routes

	}

	return pctx.JSON(http.StatusOK, response)
}

func (c *orderControllerImpl) GetWorkerPrepare(pctx echo.Context) error {
	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	order, distance, err := c.orderService.GetWorkerPrepare(orderID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"contact_name":   order.ContactName,
		"contact_phone":  order.ContactPhone,
		"user_avatar":    order.UserAvatar,
		"payment_type":   order.PaymentType,
		"location":       order.Location,
		"specific_place": order.SpecificPlace,
		"note":           order.Note,
		"duration":       order.Duration,
		"price":          order.Price,
		"created_at":     order.CreatedAt,
		"order_status":   order.OrderStatus,
	}

	if distance != nil {
		response["distance"] = distance.Routes

	}

	return pctx.JSON(http.StatusOK, response)
}

func (c *orderControllerImpl) GetQRpayment(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	qrPayment, err := c.orderService.GetQRpayment(userID.ID, orderID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusCreated, qrPayment)
}

func (c *orderControllerImpl) GetWorkerFeePayment(pctx echo.Context) error {
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	qrPayment, err := c.orderService.GetWorkerFeePayment(postID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	if qrPayment == nil {
		return pctx.JSON(http.StatusCreated, nil)
	}

	return pctx.JSON(http.StatusCreated, qrPayment)
}

func (c *orderControllerImpl) InquiryFeePayment(pctx echo.Context) error {
	paymentOrderReq := new(_paymentModel.PaymentInquiryRequest)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(paymentOrderReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	paymentResponse, err := c.orderService.InquiryWorkerFeePayment(paymentOrderReq.TransactionId)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, paymentResponse)
}

func (c *orderControllerImpl) InquiryQRpayment(pctx echo.Context) error {
	paymentOrderReq := new(_paymentModel.PaymentInquiryRequest)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(paymentOrderReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	paymentResponse, err := c.orderService.InquiryQRpayment(paymentOrderReq.TransactionId)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, paymentResponse)
}

func (c *orderControllerImpl) GetUserOrder(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	userOrder, err := c.orderService.GetUserOrder(orderID, userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, userOrder)
}

func (c *orderControllerImpl) GetWorkerOrder(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	orderID, err := c.getOrderID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	workerOrder, err := c.orderService.GetWorkerOrder(orderID, userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, workerOrder)
}

func (c *orderControllerImpl) GetActiveOrder(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	userOrder, err := c.orderService.GetActiveOrder(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	if userOrder == nil {
		return pctx.JSON(http.StatusOK, nil)
	}

	return pctx.JSON(http.StatusOK, userOrder)
}

func (c *orderControllerImpl) GetActiveWorkerOrder(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	workerOrder, err := c.orderService.GetWorkerActiveOrder(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	if workerOrder == nil {
		return pctx.JSON(http.StatusOK, nil)
	}

	return pctx.JSON(http.StatusOK, workerOrder)

}
