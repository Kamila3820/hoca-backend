package entities

import (
	"time"

	_notiModel "github.com/Kamila3820/hoca-backend/modules/notification/model"
)

type Notification struct {
	Id               *uint64 `gorm:"primaryKey"`
	Trigger          *User   `gorm:"foreignKey:TriggerID"`
	TriggerID        *string
	Triggee          *User `gorm:"foreignKey:TriggeeID"`
	TriggeeID        *string
	Order            *Order `gorm:"foreignKey:OrderID"`
	OrderID          *uint64
	UserRating       *UserRating `gorm:"foreignKey:UserRatingID"`
	UserRatingID     *uint64
	NotificationType *_notiModel.NotificationEnum `gorm:"type:notification_enum; not null"`
	CreatedAt        *time.Time                   `gorm:"not null"`
	UpdatedAt        *time.Time                   `gorm:"not null"`
}

func (n *Notification) ToNotificationResponseModel() *_notiModel.NotificationResponse {
	return &_notiModel.NotificationResponse{
		Id:               *n.Id,
		UserID:           n.Trigger.ID,
		Username:         n.Trigger.UserName,
		Avatar:           n.Trigger.Avatar,
		TriggerID:        *n.TriggerID,
		NotificationType: *n.NotificationType,
		CreatedAt:        n.CreatedAt,
	}
}
