package repository

import "github.com/Kamila3820/hoca-backend/entities"

type UserRatingRepository interface {
	ListRatingByPost(postID uint64) ([]*entities.UserRating, error)
}
