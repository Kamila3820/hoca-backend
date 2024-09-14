package model

import "time"

type (
	UserRating struct {
		ID            uint64    `json:"id"`
		UserID        string    `json:"user_id"`
		WorkerPostID  string    `json:"post_id"`
		Username      string    `json:"username"`
		UserAvatar    string    `json:"user_avatar"`
		WorkScore     int       `json:"work_score"`     // Rating value, 1-10
		SecurityScore int       `json:"security_score"` // Rating value, 1-10
		Comment       string    `json:"comment"`
		CreatedAt     time.Time `json:"created_at"`
	}

	UserRatingCreateReq struct {
		UserID        string
		WorkerPostID  string `json:"post_id"`
		WorkScore     int    `json:"work_score" validate:"required"`     // Rating value, 1-10
		SecurityScore int    `json:"security_score" validate:"required"` // Rating value, 1-10
		Comment       string `json:"comment" validate:"omitempty"`
	}
)
