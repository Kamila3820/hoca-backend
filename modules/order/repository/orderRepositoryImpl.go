package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type orderRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewOrderRepositoryImpl(db databases.Database, logger echo.Logger) OrderRepository {
	return &orderRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *orderRepositoryImpl) CreatingOrder(orderEntity *entities.Order) (*entities.Order, error) {
	order := new(entities.Order)

	if err := r.db.Connect().Create(orderEntity).Scan(order).Error; err != nil {
		r.logger.Errorf("Failed to create new order: %s", err.Error())
		return nil, err
	}

	return order, nil
}

func (r *orderRepositoryImpl) FindPostByID(postID uint64) (*entities.Post, error) {
	post := new(entities.Post)

	if err := r.db.Connect().Preload("PlaceTypes").First(post, postID).Error; err != nil {
		r.logger.Errorf("Failed to find post by ID: %s", err.Error())
		return nil, err
	}

	return post, nil
}
