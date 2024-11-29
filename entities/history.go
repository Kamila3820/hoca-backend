package entities

import (
	"time"

	_historyModel "github.com/Kamila3820/hoca-backend/modules/history/model"
)

type History struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement"`
	UserID             string    `gorm:"type:varchar(64);not null"`
	OrderID            string    `gorm:"type:varchar(64);not null"`
	Status             string    `gorm:"type:varchar(64);not null"`
	IsRated            bool      `gorm:"not null;default:false"`
	CancellationReason string    `gorm:"type:varchar(64)"` // Reason for cancellation, if any
	CancelledBy        string    `gorm:"type:varchar(64)"` // User or Worker
	CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	Order              *Order    `gorm:"foreignKey:OrderID"`
}

func (h *History) ToHistoryModel() *_historyModel.History {
	return &_historyModel.History{
		ID:                 h.ID,
		UserID:             h.UserID,
		Name:               h.Order.Post.Name,
		Price:              h.Order.Price,
		OrderID:            h.OrderID,
		Status:             h.Status,
		CancellationReason: h.CancellationReason,
		CancelledBy:        h.CancelledBy,
		IsRated:            h.IsRated,
		CreatedAt:          h.CreatedAt.Format("2006-01-02 15:04"),
	}
}

func (h *History) ToWorkingHistoryModel() *_historyModel.WorkingHistory {
	return &_historyModel.WorkingHistory{
		ID:                 h.ID,
		UserID:             h.UserID,
		Name:               h.Order.User.UserName,
		Price:              h.Order.Price,
		OrderID:            h.OrderID,
		Status:             h.Status,
		CancellationReason: h.CancellationReason,
		CancelledBy:        h.CancelledBy,
		IsRated:            h.IsRated,
		CreatedAt:          h.CreatedAt.Format("2006-01-02 15:04"),
	}
}
