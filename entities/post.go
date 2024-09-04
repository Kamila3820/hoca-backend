package entities

import (
	"time"

	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
)

type Post struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	Owner        *User          `gorm:"foreignKey:OwnerID"`
	OwnerID      string         `gorm:"type:varchar(64);not null"`
	Name         string         `gorm:"type:varchar(128);not null"`
	Description  string         `gorm:"type:varchar(256);not null"`
	Avatar       string         `gorm:"type:varchar(256);not null;default:'';"`
	CategoryID   uint64         `gorm:"not null"` // Foreign key to PostCategory
	PlaceTypeID  uint64         `gorm:"not null"` // Foreign key to PostPlaceType
	Location     string         `gorm:"type:varchar(64);not null"`
	LocationLat  string         `gorm:"type:varchar(64);not null"`
	LocationLong string         `gorm:"type:varchar(64);not null"`
	Price        float64        `gorm:"not null"`
	Distance     string         `gorm:"type:varchar(64);not null"`
	PhoneNumber  string         `gorm:"type:varchar(64);not null"`
	Gender       string         `gorm:"type:varchar(64);not null"`
	AmountFamily string         `gorm:"type:varchar(64);not null"`
	TotalScore   float64        `gorm:"not null"`
	ActiveStatus bool           `gorm:"not null;default:true"`
	CreatedAt    time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"not null;autoUpdateTime"`
	Category     *PostCategory  `gorm:"foreignKey:CategoryID"`
	PlaceType    *PostPlaceType `gorm:"foreignKey:PlaceTypeID"`
	UserRatings  []UserRating   `gorm:"foreignKey:WorkerPostID"`
}

func (p *Post) ToPostModel() *_postModel.Post {
	return &_postModel.Post{
		ID:           p.ID,
		OwnerID:      p.OwnerID,
		Name:         p.Name,
		Description:  p.Description,
		Avatar:       p.Avatar,
		CategoryID:   p.CategoryID,
		PlaceTypeID:  p.PlaceTypeID,
		Price:        p.Price,
		Distance:     p.Distance,
		Location:     p.Location,
		TotalScore:   p.TotalScore,
		PhoneNumber:  p.PhoneNumber,
		Gender:       p.Gender,
		AmountFamily: p.AmountFamily,
		CreatedAt:    p.CreatedAt,
	}
}
