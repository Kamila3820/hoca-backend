package repository

import "github.com/Kamila3820/hoca-backend/entities"

type UserRatingRepository interface {
	ListRatingByPost(postID uint64) ([]*entities.UserRating, error)
	CreateRating(ratingEntity *entities.UserRating) (*entities.UserRating, error)

	GetHistoryByID(historyID uint64) (*entities.History, error)
	UpdateHistoryByID(historyEntity *entities.History) error
	FindOrderByID(orderID uint64) (*entities.Order, error)

	CreateNotification(notiEntityy *entities.Notification) error
}
