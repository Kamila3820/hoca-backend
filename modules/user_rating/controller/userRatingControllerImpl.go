package controller

import (
	"net/http"
	"strconv"

	"github.com/Kamila3820/hoca-backend/modules/custom"
	_userRatingService "github.com/Kamila3820/hoca-backend/modules/user_rating/service"
	"github.com/labstack/echo/v4"
)

type userRatingControllerImpl struct {
	userRatingService _userRatingService.UserRatingService
}

func NewUserRatingControllerImpl(userRatingService _userRatingService.UserRatingService) UserRatingController {
	return &userRatingControllerImpl{
		userRatingService: userRatingService,
	}
}

func (c *userRatingControllerImpl) getPostID(pctx echo.Context) (uint64, error) {
	postIDStr := pctx.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return postID, nil
}

func (c *userRatingControllerImpl) ListRatingByPostID(pctx echo.Context) error {
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	postRatings, err := c.userRatingService.ListRatingByPost(postID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, postRatings)
}
