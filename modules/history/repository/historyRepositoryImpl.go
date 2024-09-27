package repository

import (
	"fmt"

	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type historyRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewHistoryRepositoryImpl(db databases.Database, logger echo.Logger) HistoryRepository {
	return &historyRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *historyRepositoryImpl) GetOrderHistory(userID string) ([]*entities.Order, error) {
	var orders []*entities.Order
	var posts []*entities.Post

	if err := r.db.Connect().Where("owner_id = ?", userID).Find(&posts).Error; err != nil {
		r.logger.Errorf("Failed to find post by userID: %s", err.Error())
		return nil, err
	}

	var workerPostID []uint64
	for _, post := range posts {
		workerPostID = append(workerPostID, post.ID)
	}
	fmt.Println(workerPostID)

	if err := r.db.Connect().
		Where("user_id = ? OR worker_post_id IN ?", userID, workerPostID).
		Where("order_status = ? OR is_cancel = ?", "complete", true).
		Order("created_at DESC").
		Find(&orders).
		Error; err != nil {
		r.logger.Errorf("Failed to find order history: %s", err.Error())
		return nil, err
	}
	fmt.Println(orders)

	return orders, nil
}

func (r *historyRepositoryImpl) GetHistory(userID string) ([]*entities.History, error) {
	var orders []*entities.Order
	var history []*entities.History

	if err := r.db.Connect().Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		r.logger.Errorf("Failed to find order by userID: %s", err.Error())
		return nil, err
	}

	var orderIDs []uint64
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	if err := r.db.Connect().Preload("Order.Post").Where("order_id IN ?", orderIDs).Order("created_at DESC").Find(&history).Error; err != nil {
		r.logger.Errorf("Failed to find history: %s", err.Error())
		return nil, err
	}

	return history, nil
}

func (r *historyRepositoryImpl) GetWorkingHistory(userID string) ([]*entities.History, error) {
	var history []*entities.History

	if err := r.db.Connect().Preload("Order.User").Where("user_id = ?", userID).Order("created_at DESC").Find(&history).Error; err != nil {
		r.logger.Errorf("Failed to find worker history by ID: %s", err.Error())
		return nil, err
	}

	return history, nil
}
