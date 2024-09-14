package entities

import (
	"time"

	_userRatingModel "github.com/Kamila3820/hoca-backend/modules/user_rating/model"
)

type UserRating struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	UserID        string    `gorm:"type:varchar(64);not null"`
	WorkerPostID  string    `gorm:"type:varchar(64);not null"`
	WorkScore     int       `gorm:"not null"` // Rating value, 1-10
	SecurityScore int       `gorm:"not null"` // Rating value, 1-10
	Comment       string    `gorm:"type:varchar(128)"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	User          *User     `gorm:"foreignKey:UserID"`
	Post          *Post     `gorm:"foreignKey:WorkerPostID"`
}

func (r *UserRating) ToUserRatingModel() *_userRatingModel.UserRating {
	var userName, userAvatar string
	if r.User != nil {
		userName = r.User.UserName
		userAvatar = r.User.Avatar
	}

	return &_userRatingModel.UserRating{
		ID:            r.ID,
		UserID:        r.UserID,
		Username:      userName,
		UserAvatar:    userAvatar,
		WorkerPostID:  r.WorkerPostID,
		WorkScore:     r.WorkScore,
		SecurityScore: r.SecurityScore,
		Comment:       r.Comment,
		CreatedAt:     r.CreatedAt,
	}
}
