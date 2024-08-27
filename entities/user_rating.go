package entities

import "time"

type UserRating struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	UserID        uint64    `gorm:"not null"`
	WorkerPostID  uint64    `gorm:"not null"`
	WorkScore     int       `gorm:"not null"` // Rating value, 1-10
	SecurityScore int       `gorm:"not null"` // Rating value, 1-10
	Comment       string    `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	User          *User     `gorm:"foreignKey:UserID"`
	Post          *Post     `gorm:"foreignKey:WorkerPostID"`
}