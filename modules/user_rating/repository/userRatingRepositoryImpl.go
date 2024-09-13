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
