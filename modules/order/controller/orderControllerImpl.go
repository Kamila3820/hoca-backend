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
