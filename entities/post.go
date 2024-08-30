package entities

import "time"

type Post struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	UserID       uint64         `gorm:"not null"`
	Name         string         `gorm:"type:varchar(128);not null"`
	Description  string         `gorm:"type:text;not null"`
	Avatar       string         `gorm:"type:varchar(256);not null;default:'';"`
	CategoryID   uint64         `gorm:"not null"` // Foreign key to PostCategory
	PlaceTypeID  uint64         `gorm:"not null"` // Foreign key to PostPlaceType
	LocationLat  string         `gorm:"type:varchar(64);not null"`
	LocationLong string         `gorm:"type:varchar(64);not null"`
	Price        float64        `gorm:"not null"`
	PhoneNumber  string         `gorm:"type:varchar(64);not null"`
	AmountFamily string         `gorm:"type:varchar(64);not null"`
	TotalScore   float64        `gorm:"not null"`
	ActiveStatus bool           `gorm:"not null;default:true"`
	CreatedAt    time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"not null;autoUpdateTime"`
	User         *User          `gorm:"foreignKey:UserID"`
	Category     *PostCategory  `gorm:"foreignKey:CategoryID"`
	PlaceType    *PostPlaceType `gorm:"foreignKey:PlaceTypeID"` // Corrected this line
	UserRatings  []UserRating   `gorm:"foreignKey:WorkerPostID"`
}
