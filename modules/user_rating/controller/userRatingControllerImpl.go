package controller

import (
	_userRatingService "github.com/Kamila3820/hoca-backend/modules/user_rating/service"
)

type userRatingControllerImpl struct {
	userRatingService _userRatingService.UserRatingService
}

func NewUserRatingControllerImpl(userRatingService _userRatingService.UserRatingService) UserRatingController {
	return &userRatingControllerImpl{
		userRatingService: userRatingService,
	}
}
