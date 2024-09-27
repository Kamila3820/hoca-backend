package service

import (
	"errors"
	"strconv"

	"github.com/Kamila3820/hoca-backend/entities"
	_historyModel "github.com/Kamila3820/hoca-backend/modules/history/model"
	_historyRepository "github.com/Kamila3820/hoca-backend/modules/history/repository"
)

type historyServiceImpl struct {
	historyRepository _historyRepository.HistoryRepository
}

func NewHistoryServiceImpl(historyRepository _historyRepository.HistoryRepository) HistoryService {
	return &historyServiceImpl{
		historyRepository: historyRepository,
	}
}

func (s *historyServiceImpl) GetOrderHistory(userID string) ([]*_historyModel.History, error) {
	// Get order history from repository
	orders, err := s.historyRepository.GetOrderHistory(userID)
	if err != nil {
		return nil, errors.New("service: cannot find order history")
	}

	// Initialize the history slice to avoid panic
	history := make([]*entities.History, len(orders))

	// Convert orders to history model
	for i, order := range orders {
		history[i] = &entities.History{
			UserID:             order.UserID,
			OrderID:            strconv.Itoa(int(order.ID)),
			Status:             order.OrderStatus,
			CancelledBy:        order.CancelledBy,
			CancellationReason: order.CancellationReason,
			IsRated:            false,
			CreatedAt:          order.UpdatedAt,
		}
	}

	// Convert to HistoryModel
	orderHistory := make([]*_historyModel.History, 0)
	for _, oh := range history {
		orderHistory = append(orderHistory, oh.ToHistoryModel())
	}

	return orderHistory, nil
}

func (s *historyServiceImpl) GetHistory(userID string) ([]*_historyModel.History, error) {
	histories, err := s.historyRepository.GetHistory(userID)
	if err != nil {
		return nil, errors.New("service: cannot find history")
	}

	userHistory := make([]*_historyModel.History, 0)
	for _, history := range histories {
		userHistory = append(userHistory, history.ToHistoryModel())
	}

	return userHistory, nil
}

func (s *historyServiceImpl) GetWorkingHistory(userID string) ([]*_historyModel.WorkingHistory, error) {
	histories, err := s.historyRepository.GetWorkingHistory(userID)
	if err != nil {
		return nil, errors.New("service: cannot find working history")
	}

	workingHistory := make([]*_historyModel.WorkingHistory, 0)
	for _, history := range histories {
		workingHistory = append(workingHistory, history.ToWorkingHistoryModel())
	}

	return workingHistory, nil
}
