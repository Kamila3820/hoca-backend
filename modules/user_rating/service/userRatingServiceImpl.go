package service

import (
	_userRatingModel "github.com/Kamila3820/hoca-backend/modules/user_rating/model"
	_userRatingRepository "github.com/Kamila3820/hoca-backend/modules/user_rating/repository"
)

type userRatingServiceImpl struct {
	userRatingRepository _userRatingRepository.UserRatingRepository
}

func NewUserRatingServiceImpl(userRatingRepository _userRatingRepository.UserRatingRepository) UserRatingService {
	return &userRatingServiceImpl{
		userRatingRepository: userRatingRepository,
	}
}

func (s *userRatingServiceImpl) ListRatingByPost(postID uint64) ([]*_userRatingModel.UserRating, error) {
	ratings, err := s.userRatingRepository.ListRatingByPost(postID)
	if err != nil {
		return nil, err
	}

	postRating := make([]*_userRatingModel.UserRating, 0)
	for _, rating := range ratings {
		postRating = append(postRating, rating.ToUserRatingModel())
	}

	return postRating, nil
}
