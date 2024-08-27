package entities

import "time"

type Order struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement"`
	UserID             uint64    `gorm:"not null"`
	WorkerPostID       uint64    `gorm:"not null"`
	PaymentType        string    `gorm:"type:varchar(64);not null"`
	SpecificPlace      string    `gorm:"type:varchar(64)"`
	Note               string    `gorm:"type:varchar(64)"`
	OrderStatus        string    `gorm:"type:varchar(64);not null"`
	Price              float64   `gorm:"not null"`
	IsCancel           bool      `gorm:"not null;default:false;"`
	CancellationReason string    `gorm:"type:text"`        // Reason for cancellation, if any
	CancelledBy        string    `gorm:"type:varchar(64)"` // User or Worker
	CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"not null;autoUpdateTime"`
	User               *User     `gorm:"foreignKey:UserID"`
	Post               *Post     `gorm:"foreignKey:WorkerPostID"`
}
