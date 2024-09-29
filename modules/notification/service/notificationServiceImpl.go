package service

import (
	"errors"

	_notiModel "github.com/Kamila3820/hoca-backend/modules/notification/model"
	_notificationRepository "github.com/Kamila3820/hoca-backend/modules/notification/repository"
)

type notificationServiceImpl struct {
	notificationRepository _notificationRepository.NotificationRepository
}

func NewNotificationServiceImpl(notificationRepository _notificationRepository.NotificationRepository) NotificationService {
	return &notificationServiceImpl{
		notificationRepository: notificationRepository,
	}
}

func (s *notificationServiceImpl) GetNotificationsByUser(userID string) ([]*_notiModel.NotificationResponse, error) {
	notifications, err := s.notificationRepository.GetNotificationByUser(userID)
	if err != nil {
		return nil, errors.New("service: cannot find notifications by user")
	}

	mappedNotifications := make([]*_notiModel.NotificationResponse, 0)
	for _, notification := range notifications {
		if notification.TriggerID == &userID {
			continue
		}

		resp := notification.ToNotificationResponseModel()

		if notification.OrderID != nil {
			resp.Order = &_notiModel.OrderResponse{
				OrderID: *notification.OrderID,
			}
			resp.OrderID = *notification.OrderID
		} else if notification.UserRatingID != nil {
			resp.UserRating = &_notiModel.UserRatingResponse{
				UserRatingID: *notification.UserRatingID,
			}
			resp.UserRatingID = *notification.UserRatingID
		}

		mappedNotifications = append(mappedNotifications, resp)

	}

	return mappedNotifications, nil
}
