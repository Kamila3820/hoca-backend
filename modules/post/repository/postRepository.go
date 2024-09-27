package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
)

type PostRepository interface {
	FindPost() ([]*entities.Post, error)
	FindPostByID(postID uint64) (*entities.Post, error)
	FindPostByUserID(userID string) (*entities.Post, error)
	CreatingPost(postEntity *entities.Post) (*entities.Post, error)
	GetPlaceTypesByIds(placeTypeIDs []uint64) ([]*entities.PlaceType, error)
	EditingPost(postID uint64, postEditingReq *_postModel.PostEditingReq) (uint64, error)
	DeletePost(postID uint64) error

	Activate(postID uint64) error
	UnActivate(postID uint64) error
}
