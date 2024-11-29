package controller

import "github.com/labstack/echo/v4"

type UserRatingController interface {
	ListRatingByPostID(pctx echo.Context) error
	RatingWorker(pctx echo.Context) error
	GetRatingMetrics(pctx echo.Context) error
}
