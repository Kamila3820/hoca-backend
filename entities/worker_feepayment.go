package entities

import "time"

type WorkerFeepayment struct {
	Post          *Post     `gorm:"foreignKey:PostId"`
	PostId        uint64    `gorm:"not null"`
	OrderCount    uint64    `gorm:"not null"`
	Amount        uint64    `gorm:"not null"`
	TransactionID string    `gorm:"not null"`
	Paid          bool      `gorm:"not null"`
	StartFrom     string    `gorm:"not null"`
	EndFrom       string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
}
