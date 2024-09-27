package controller

import "github.com/labstack/echo/v4"

type PostController interface {
	FindPostByDistance(pctx echo.Context) error
	GetPostByUserID(pctx echo.Context) error
	CreateWorkerPost(pctx echo.Context) error
	EditWorkerPost(pctx echo.Context) error
	getPostID(pctx echo.Context) (uint64, error)
	DeleteWorkerPost(pctx echo.Context) error

	ActivatePost(pctx echo.Context) error
	UnActivatePost(pctx echo.Context) error
}
