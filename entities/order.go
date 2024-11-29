package entities

import (
	"time"

	_orderModel "github.com/Kamila3820/hoca-backend/modules/order/model"
)

type Order struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement"`
	UserID             string    `gorm:"not null"`
	WorkerPostID       uint64    `gorm:"not null"`
	ContactName        string    `gorm:"type:varchar(64)"`
	ContactPhone       string    `gorm:"type:varchar(20)"`
	PaymentType        string    `gorm:"type:varchar(64);not null"`
	SpecificPlace      string    `gorm:"type:varchar(64)"`
	Note               string    `gorm:"type:varchar(64)"`
	Duration           string    `gorm:"type:varchar(128);not null"`
	OrderStatus        string    `gorm:"type:varchar(64);not null"`
	Price              float64   `gorm:"not null"`
	Paid               bool      `gorm:"not null;default:false;"`
	IsCancel           bool      `gorm:"not null;default:false;"`
	CancellationReason string    `gorm:"type:varchar(64)"` // Reason for cancellation, if any
	CancelledBy        string    `gorm:"type:varchar(64)"` // User or Worker
	CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"not null;autoUpdateTime"`
	User               *User     `gorm:"foreignKey:UserID"`
	Post               *Post     `gorm:"foreignKey:WorkerPostID"`
}

func (o *Order) ToOrderModel() *_orderModel.Order {
	return &_orderModel.Order{
		ID:                 o.ID,
		UserID:             o.UserID,
		WorkerPostID:       o.WorkerPostID,
		ContactName:        o.ContactName,
		ContactPhone:       o.ContactPhone,
		PaymentType:        o.PaymentType,
		SpecificPlace:      o.SpecificPlace,
		Note:               o.Note,
		Duration:           o.Duration,
		OrderStatus:        o.OrderStatus,
		Price:              o.Price,
		Paid:               o.Paid,
		IsCancel:           o.IsCancel,
		CancellationReason: o.CancellationReason,
		CancelledBy:        o.CancelledBy,
		CreatedAt:          o.CreatedAt,
		UpdatedAt:          o.UpdatedAt,
	}
}

func (o *Order) ToUserOrder() *_orderModel.UserOrder {
	return &_orderModel.UserOrder{
		ID:           o.ID,
		UserID:       o.UserID,
		WorkerPostID: o.WorkerPostID,
		WorkerName:   o.Post.Name,
		WorkerPhone:  o.Post.PhoneNumber,
		WorkerAvatar: o.Post.Avatar,
		Price:        o.Price,
		Paid:         o.Paid,
		PaymentType:  o.PaymentType,
		OrderStatus:  o.OrderStatus,
	}
}

func (o *Order) ToWorkerOrder() *_orderModel.WorkerOrder {
	return &_orderModel.WorkerOrder{
		ID:            o.ID,
		UserID:        o.UserID,
		WorkerPostID:  o.WorkerPostID,
		ContactName:   o.ContactName,
		ContactPhone:  o.ContactPhone,
		UserAvatar:    o.User.Avatar,
		PaymentType:   o.PaymentType,
		SpecificPlace: o.SpecificPlace,
		Location:      o.User.Location,
		Note:          o.Note,
		Duration:      o.Duration,
		Price:         o.Price,
		Paid:          o.Paid,
		OrderStatus:   o.OrderStatus,
		CreatedAt:     o.CreatedAt.Format("2006-01-02 15:04"),
		EndedAt:       o.CreatedAt.Add(4 * time.Minute).Format("2006-01-02 15:04"),
	}
}
