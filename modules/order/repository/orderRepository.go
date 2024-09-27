package repository

import "github.com/Kamila3820/hoca-backend/entities"

type OrderRepository interface {
	CreatingOrder(orderEntity *entities.Order) (*entities.Order, error)
	FindOrderByID(orderID uint64) (*entities.Order, error)
	UpdateOrder(orderEntity *entities.Order) error
	CreatingHistory(historyEntity *entities.History) (*entities.History, error)

	FindPostByID(postID uint64) (*entities.Post, error)
	FindUserByID(userID string) (*entities.User, error)
	UpdatePost(postEntity *entities.Post) error
}
