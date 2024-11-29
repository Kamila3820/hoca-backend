package model

type (
	History struct {
		ID                 uint64  `json:"id"`
		UserID             string  `json:"user_id"`
		OrderID            string  `json:"order_id"`
		Status             string  `json:"status"`
		IsRated            bool    `json:"is_rated"`
		CancellationReason string  `json:"cancellation_reason"` // Reason for cancellation, if any
		CancelledBy        string  `json:"cancelled_by"`        // User or Worker
		Name               string  `json:"name"`
		Price              float64 `json:"price"`
		CreatedAt          string  `json:"created_at"`
	}

	WorkingHistory struct {
		ID                 uint64  `json:"id"`
		UserID             string  `json:"user_id"`
		OrderID            string  `json:"order_id"`
		Status             string  `json:"status"`
		IsRated            bool    `json:"is_rated"`
		CancellationReason string  `json:"cancellation_reason"` // Reason for cancellation, if any
		CancelledBy        string  `json:"cancelled_by"`        // User or Worker
		Name               string  `json:"name"`
		Price              float64 `json:"price"`
		CreatedAt          string  `json:"created_at"`
	}
)
