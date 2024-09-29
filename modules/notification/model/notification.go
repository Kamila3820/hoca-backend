package model

import (
	"time"
)

type (
	NotificationResponse struct {
		Id               uint64              `json:"id"`
		UserID           string              `json:"user_id"`
		Username         string              `json:"username"`
		Avatar           string              `json:"avatar"`
		TriggerID        string              `json:"trigger_id"`
		Order            *OrderResponse      `json:"order"`
		OrderID          uint64              `json:"order_id"`
		UserRating       *UserRatingResponse `json:"user_rating"`
		UserRatingID     uint64              `json:"user_rating_id"`
		NotificationType NotificationEnum    `json:"type"`
		CreatedAt        *time.Time          `json:"created_at"`
	}

	OrderResponse struct {
		OrderID uint64 `json:"order_id"`
	}

	UserRatingResponse struct {
		UserRatingID uint64 `json:"user_rating_id"`
	}
)
