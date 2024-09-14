package service

import (
	"fmt"

	"github.com/Kamila3820/hoca-backend/entities"
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

func (s *userRatingServiceImpl) CreateRating(ratingCreateReq *_userRatingModel.UserRatingCreateReq) (*_userRatingModel.UserRating, error) {
	userRatingEntity := &entities.UserRating{
		UserID:        ratingCreateReq.UserID,
		WorkerPostID:  ratingCreateReq.WorkerPostID,
		WorkScore:     ratingCreateReq.WorkScore,
		SecurityScore: ratingCreateReq.SecurityScore,
		Comment:       ratingCreateReq.Comment,
	}
	fmt.Println("1 SERVICE")

	userRating, err := s.userRatingRepository.CreateRating(userRatingEntity)
	if err != nil {
		return nil, err
	}
	fmt.Println("2 SERVICE")
	fmt.Printf("userRating model: %+v\n", userRating)

	return userRating.ToUserRatingModel(), nil
}
