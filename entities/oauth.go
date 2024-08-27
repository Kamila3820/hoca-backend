package entities

import "time"

type Oauth struct {
	ID           string    `gorm:"primaryKey;type:varchar(64);"`
	UserID       string    `gorm:"type:varchar(64);not null;index"`
	AccessToken  string    `gorm:"type:text;not null;"`
	RefreshToken string    `gorm:"type:text;not null;"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}
