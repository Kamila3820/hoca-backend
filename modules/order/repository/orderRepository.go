package repository

import "github.com/Kamila3820/hoca-backend/entities"

type OrderRepository interface {
	CreatingOrder(orderEntity *entities.Order) (*entities.Order, error)

	FindPostByID(postID uint64) (*entities.Post, error)
}
