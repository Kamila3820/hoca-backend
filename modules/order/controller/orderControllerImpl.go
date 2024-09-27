package controller

import (
	//_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	"net/http"
	"strconv"

	"github.com/Kamila3820/hoca-backend/modules/custom"
	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
	_orderService "github.com/Kamila3820/hoca-backend/modules/order/service"
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
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	orderCreatingReq := new(_orderModel.OrderReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(orderCreatingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	orderCreatingReq.UserID = userIDStr

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
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	user, err := c.orderService.GetUserByID(userIDStr)
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
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

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

	order, err := c.orderService.UpdateOrderProgress(userIDStr, orderID, statusUpdate.Status)
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
		"order_status": order.OrderStatus,
	}

	if distance != nil {
		response["distance"] = distance
	}

	return pctx.JSON(http.StatusOK, response)
}
