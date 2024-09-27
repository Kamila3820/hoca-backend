package service

import (
	_historyModel "github.com/Kamila3820/hoca-backend/modules/history/model"
)

type HistoryService interface {
	GetOrderHistory(userID string) ([]*_historyModel.History, error)
	GetHistory(userID string) ([]*_historyModel.History, error)
	GetWorkingHistory(userID string) ([]*_historyModel.WorkingHistory, error)
}
