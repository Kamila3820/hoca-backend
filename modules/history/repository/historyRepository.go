package repository

import "github.com/Kamila3820/hoca-backend/entities"

type HistoryRepository interface {
	GetOrderHistory(userID string) ([]*entities.Order, error)
	GetHistory(userID string) ([]*entities.History, error)
	GetWorkingHistory(userID string) ([]*entities.History, error)
}
