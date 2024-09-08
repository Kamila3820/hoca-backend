package entities

import "time"

type Notification struct {
	Id               *uint64 `gorm:"primaryKey"`
	Trigger          *User   `gorm:"foreignKey:TriggerId"`
	TriggerId        *string
	Triggee          *User `gorm:"foreignKey:TriggeeId"`
	TriggeeId        *string
	Post             *Post `gorm:"foreignKey:PostId"`
	PostId           *string
	NotificationType string     `gorm:"type:varchar(128); not null"`
	CreatedAt        *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt        *time.Time `gorm:"not null"` // Embedded field
}

// 	NotificationType *enum.Notification `gorm:"type:ENUM('comment','like','user_donate','post_donate'); not null"`
