package model

import "time"

type (
	Order struct {
		ID                 uint64    `json:"id"`
		UserID             string    `json:"user_id"`
		WorkerPostID       uint64    `json:"worker_post_id"`
		PaymentType        string    `json:"payment_type"`
		SpecificPlace      string    `json:"specific_place"`
		Note               string    `json:"note"`
		OrderStatus        string    `json:"order_status"`
		Price              float64   `json:"price"`
		IsCancel           bool      `json:"is_cancel"`
		CancellationReason string    `json:"cancellation_reason"` // Reason for cancellation, if any
		CancelledBy        string    `json:"cancelled_by"`        // User or Worker
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_by"`
	}

	OrderReq struct {
		UserID        string
		PaymentType   string `json:"payment_type" validate:"required"`
		SpecificPlace string `json:"specific_place,omitempty"`
		Note          string `json:"note,omitempty"`
	}
)
