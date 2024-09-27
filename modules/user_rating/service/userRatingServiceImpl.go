package service

import (
	"errors"
	"strconv"

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

func (s *userRatingServiceImpl) CreateRating(raterID string, historyID uint64, ratingCreateReq *_userRatingModel.UserRatingCreateReq) (*_userRatingModel.UserRating, error) {
	history, err := s.userRatingRepository.GetHistoryByID(historyID)
	if err != nil {
		return nil, errors.New("service: cannot find history by id")
	}

	if history.Status == "cancelled" {
		return nil, errors.New("service: cannot rate the order that has been cancelled")
	}

	if history.IsRated {
		return nil, errors.New("service: cannot rate the order that has been rated")
	}

	orderID, err := strconv.ParseUint(history.OrderID, 10, 64)
	if err != nil {
		return nil, errors.New("service: cannot find history by id")
	}

	order, err := s.userRatingRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("service: cannot find order by the history")
	}

	if raterID != order.UserID {
		return nil, errors.New("service: you have no permission to rate the worker")
	}

	postID := strconv.Itoa(int(order.WorkerPostID))

	userRatingEntity := &entities.UserRating{
		UserID:        order.UserID,
		WorkerPostID:  postID,
		WorkScore:     ratingCreateReq.WorkScore,
		SecurityScore: ratingCreateReq.SecurityScore,
		Comment:       ratingCreateReq.Comment,
	}

	userRating, err := s.userRatingRepository.CreateRating(userRatingEntity)
	if err != nil {
		return nil, errors.New("service: cannot rate the worker")
	}

	history.IsRated = true
	s.userRatingRepository.UpdateHistoryByID(history)

	return userRating.ToUserRatingModel(), nil
}
