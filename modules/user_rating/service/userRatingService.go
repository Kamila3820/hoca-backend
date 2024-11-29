package service

import (
	_userRatingModel "github.com/Kamila3820/hoca-backend/modules/user_rating/model"
)

type UserRatingService interface {
	ListRatingByPost(postID uint64) ([]*_userRatingModel.UserRating, error)
	CreateRating(raterID string, historyID uint64, ratingCreateReq *_userRatingModel.UserRatingCreateReq) (*_userRatingModel.UserRating, error)
	GetRatingMetrics(postID uint64) (*_userRatingModel.RatingMetricsResponse, error)
}
