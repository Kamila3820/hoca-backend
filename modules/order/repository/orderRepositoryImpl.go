package repository

import (
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
