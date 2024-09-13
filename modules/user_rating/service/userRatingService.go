package service

import (
	_userRatingModel "github.com/Kamila3820/hoca-backend/modules/user_rating/model"
)

type UserRatingService interface {
	ListRatingByPost(postID uint64) ([]*_userRatingModel.UserRating, error)
}
