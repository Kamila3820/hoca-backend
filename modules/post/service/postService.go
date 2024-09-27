package service

import (
	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
)

type PostService interface {
	FindPostByDistance(userLat, userLong float64) ([]*_postModel.Post, error)
	GetPostByUserID(userID string) (*_postModel.Post, error)
	CreatingPost(postCreatingReq *_postModel.PostCreatingReq) (*_postModel.Post, error)
	EditingPost(postID uint64, postEditingReq *_postModel.PostEditingReq) (*_postModel.Post, error)
	DeletePost(postID uint64, userID string) error

	Activate(postID uint64) error
	UnActivate(postID uint64) error
}
