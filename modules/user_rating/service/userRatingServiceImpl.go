package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Kamila3820/hoca-backend/entities"
	_notiModel "github.com/Kamila3820/hoca-backend/modules/notification/model"
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

	go func() {
		history.IsRated = true
		s.userRatingRepository.UpdateHistoryByID(history)

		ratingType := _notiModel.NotificationRating
		notification := &entities.Notification{
			Trigger:          nil,
			TriggerID:        &raterID,
			Triggee:          nil,
			TriggeeID:        &order.Post.OwnerID,
			Order:            nil,
			OrderID:          nil,
			UserRating:       nil,
			UserRatingID:     &userRating.ID,
			NotificationType: &ratingType,
			CreatedAt:        nil,
		}

		if err := s.userRatingRepository.CreateNotification(notification); err != nil {
			fmt.Printf("service: unable to create notification %v", err.Error)
		}
	}()

	return userRating.ToUserRatingModel(), nil
}

func (s *userRatingServiceImpl) GetRatingMetrics(postID uint64) (*_userRatingModel.RatingMetricsResponse, error) {
	ratings, err := s.userRatingRepository.ListRatingByPost(postID)
	if err != nil {
		return nil, err
	}

	count := len(ratings)

	if count == 0 {
		return nil, errors.New("no ratings found for the specified post")
	}

	totalWorkScore, totalSecurityScore, maxScore, minScore := 0, 0, 0, 20
	for _, rating := range ratings {
		workScore := rating.WorkScore
		securityScore := rating.SecurityScore
		totalWorkScore += workScore
		totalSecurityScore += securityScore

		combinedScore := workScore + securityScore
		if combinedScore > maxScore {
			maxScore = combinedScore
		}
		if combinedScore < minScore {
			minScore = combinedScore
		}
	}

	totalScore := float64(totalWorkScore+totalSecurityScore) / float64(2*count)
	avgWorkScore := float64(totalWorkScore) / float64(len(ratings))
	avgSecurityScore := float64(totalSecurityScore) / float64(len(ratings))
	max := float64(maxScore) / float64(2)
	min := float64(minScore) / float64(2)

	summaryTotalScore, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", totalScore), 64)
	summaryAvgWorkScore, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgWorkScore), 64)
	summaryAvgSecurityScore, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgSecurityScore), 64)
	summaryMaxScore, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", max), 64)
	summaryMinScore, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", min), 64)

	metrics := &_userRatingModel.RatingMetricsResponse{
		TotalScore:       summaryTotalScore,
		AvgWorkScore:     summaryAvgWorkScore,
		AvgSecurityScore: summaryAvgSecurityScore,
		MaxCombinedScore: summaryMaxScore,
		MinCombinedScore: summaryMinScore,
		TotalRatings:     count,
	}

	return metrics, nil
}
