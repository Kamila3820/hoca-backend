package entities

import "time"

type OrderQrpayment struct {
	User          *User     `gorm:"foreignKey:UserId"`
	UserId        string    `gorm:"not null"`
	OrderID       uint64    `gorm:"not null"`
	Amount        uint64    `gorm:"not null"`
	TransactionID string    `gorm:"not null"`
	Paid          bool      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	Order         *Order    `gorm:"foreignKey:OrderID"`
}
