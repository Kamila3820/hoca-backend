package controller

import (
	"net/http"
	"strconv"

	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	_userRatingModel "github.com/Kamila3820/hoca-backend/modules/user_rating/model"
	_userRatingService "github.com/Kamila3820/hoca-backend/modules/user_rating/service"
	"github.com/golang-jwt/jwt/v5"
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
func (c *userRatingControllerImpl) getHistoryID(pctx echo.Context) (uint64, error) {
	historyIDStr := pctx.Param("historyID")
	historyID, err := strconv.ParseUint(historyIDStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return historyID, nil
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

func (c *userRatingControllerImpl) RatingWorker(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	historyID, err := c.getHistoryID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	userRatingCreateReq := new(_userRatingModel.UserRatingCreateReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(userRatingCreateReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	userRating, err := c.userRatingService.CreateRating(userID.ID, historyID, userRatingCreateReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusCreated, userRating)
}

func (c *userRatingControllerImpl) GetRatingMetrics(pctx echo.Context) error {
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	metrics, err := c.userRatingService.GetRatingMetrics(postID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, metrics)
}
