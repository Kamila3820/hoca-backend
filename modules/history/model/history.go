package model

import "time"

type (
	History struct {
		ID                 uint64    `gorm:"primaryKey;autoIncrement"`
		UserID             string    `gorm:"type:varchar(64);not null"`
		OrderID            string    `gorm:"type:varchar(64);not null"`
		Status             string    `gorm:"type:varchar(64);not null"`
		IsRated            bool      `gorm:"not null;default:false"`
		CancellationReason string    `gorm:"type:varchar(64)"` // Reason for cancellation, if any
		CancelledBy        string    `gorm:"type:varchar(64)"` // User or Worker
		Name               string    `json:"name"`
		Price              float64   `json:"price"`
		CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	}

	WorkingHistory struct {
		ID                 uint64    `gorm:"primaryKey;autoIncrement"`
		UserID             string    `gorm:"type:varchar(64);not null"`
		OrderID            string    `gorm:"type:varchar(64);not null"`
		Status             string    `gorm:"type:varchar(64);not null"`
		IsRated            bool      `gorm:"not null;default:false"`
		CancellationReason string    `gorm:"type:varchar(64)"` // Reason for cancellation, if any
		CancelledBy        string    `gorm:"type:varchar(64)"` // User or Worker
		Name               string    `json:"name"`
		Price              float64   `json:"price"`
		CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	}
)