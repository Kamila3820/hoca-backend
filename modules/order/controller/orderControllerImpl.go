package controller

import (
	//_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"

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

	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	orderCreatingReq := new(_orderModel.OrderReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(orderCreatingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	orderCreatingReq.UserID = userID.ID

	newOrder, err := c.orderService.CreatingOrder(orderCreatingReq, postID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

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

	customerLat := pctx.QueryParam("lat")
	customerLong := pctx.QueryParam("long")

	order, distance, err := c.orderService.GetPreparingOrder(orderID, customerLat, customerLong)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"order_status": order.OrderStatus,
	}

	if distance != nil {
		response["distance"] = distance
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

	return pctx.JSON(http.StatusOK, qrPayment)
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
