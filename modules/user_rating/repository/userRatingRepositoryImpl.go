package repository

import (
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
