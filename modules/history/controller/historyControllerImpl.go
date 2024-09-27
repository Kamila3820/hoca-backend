package controller

import (
	"net/http"

	"github.com/Kamila3820/hoca-backend/modules/custom"
	_historyService "github.com/Kamila3820/hoca-backend/modules/history/service"
	"github.com/labstack/echo/v4"
)

type historyControllerImpl struct {
	historyService _historyService.HistoryService
}

func NewHistoryControllerImpl(historyService _historyService.HistoryService) HistoryController {
	return &historyControllerImpl{
		historyService: historyService,
	}
}

func (c *historyControllerImpl) GetOrderHistoryByUserID(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	history, err := c.historyService.GetOrderHistory(userIDStr)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, history)
}

func (c *historyControllerImpl) GetHistoryByUserID(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	history, err := c.historyService.GetHistory(userIDStr)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, history)
}

func (c *historyControllerImpl) GetWorkingHistory(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	history, err := c.historyService.GetWorkingHistory(userIDStr)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, history)
}
