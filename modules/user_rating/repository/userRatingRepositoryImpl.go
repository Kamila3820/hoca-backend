package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type userRatingRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewUserRatingRepositoryImpl(db databases.Database, logger echo.Logger) UserRatingRepository {
	return &userRatingRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *userRatingRepositoryImpl) ListRatingByPost(postID uint64) ([]*entities.UserRating, error) {
	ratings := make([]*entities.UserRating, 0)
	if err := r.db.Connect().Preload("User").Where("worker_post_id = ?", postID).Find(&ratings).Error; err != nil {
		r.logger.Errorf("Failed to list ratings: %s", err.Error())
		return nil, err
	}

	return ratings, nil
}

func (r *userRatingRepositoryImpl) CreateRating(ratingEntity *entities.UserRating) (*entities.UserRating, error) {
	userRating := new(entities.UserRating)

	if err := r.db.Connect().Create(ratingEntity).
		Preload("User").
		Preload("Post").
		Scan(userRating).Error; err != nil {
		r.logger.Errorf("Failed to create ratings: %s", err.Error())
		return nil, err
	}

	return userRating, nil
}

// History
func (r *userRatingRepositoryImpl) GetHistoryByID(historyID uint64) (*entities.History, error) {
	history := new(entities.History)

	if err := r.db.Connect().First(&history, historyID).Error; err != nil {
		r.logger.Errorf("Failed to find history by ID: %s", err.Error())
		return nil, err
	}

	return history, nil
}

func (r *userRatingRepositoryImpl) UpdateHistoryByID(historyEntity *entities.History) error {
	if err := r.db.Connect().Save(historyEntity).Error; err != nil {
		r.logger.Errorf("Failed to update history: %s", err.Error())
		return err
	}

	return nil
}

// Order
func (r *userRatingRepositoryImpl) FindOrderByID(orderID uint64) (*entities.Order, error) {
	order := new(entities.Order)

	if err := r.db.Connect().Where("id = ?", orderID).Preload("User").Preload("Post").First(&order).Error; err != nil {
		r.logger.Errorf("Failed to find order by ID: %s", err.Error())
		return nil, err
	}

	return order, nil
}
