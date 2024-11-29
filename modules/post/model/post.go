package model

type (
	Post struct {
		ID             uint64       `json:"id"`
		OwnerID        string       `json:"owner_id"`
		Name           string       `json:"name"`
		Description    string       `json:"description"`
		Avatar         string       `json:"avatar"`
		Location       string       `json:"location"`
		LocationLat    string       `json:"latitude"`
		LocationLong   string       `json:"longitude"`
		Price          float64      `json:"price"`
		PromptPay      string       `json:"prompt_pay"`
		Distance       string       `json:"distance"`
		DistanceFee    string       `json:"distance_fee"`
		PhoneNumber    string       `json:"phone_number"`
		Gender         string       `json:"gender"`
		AmountFamily   string       `json:"amount_family"`
		AvailableStart string       `json:"available_start"`
		AvailableEnd   string       `json:"available_end"`
		Duration       string       `json:"duration"`
		TotalScore     float64      `json:"total_score"`
		ActiveStatus   bool         `json:"active_status"`
		Categories     []Category   `json:"categories"`
		PlaceTypes     []PlaceType  `json:"place_types"`
		UserRatings    []UserRating `json:"user_ratings"`
		CreatedAt      string       `json:"created_at"`
		UpdatedAt      string       `json:"updated_at"`
	}

	PlaceType struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	Category struct {
		ID          uint64 `json:"id"`
		GroupID     uint64 `json:"group_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UserRating struct {
		ID            uint64 `json:"id"`
		UserID        string `json:"user_id"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		WorkScore     int    `json:"work_score"`
		SecurityScore int    `json:"security_score"`
		Comment       string `json:"comment"`
		CreatedAt     string `json:"created_at"`
	}

	PostCreatingReq struct {
		OwnerID        string
		Name           string   `json:"name" validate:"required,max=64"`
		Description    string   `json:"description" validate:"omitempty,max=128"`
		Avatar         string   `json:"avatar" validate:"omitempty"`
		CategoryIDs    []uint64 `json:"category_ids" validate:"required"`
		PlaceTypeIDs   []uint64 `json:"placetype_ids" validate:"required"`
		Location       string   `json:"location" validate:"required"`
		LocationLat    string   `json:"latitude" validate:"required"`
		LocationLong   string   `json:"longtitude" validate:"required"`
		Price          float64  `json:"price" validate:"required"`
		PromptPay      string   `json:"prompt_pay" validate:"omitempty,max=128"`
		PhoneNumber    string   `json:"phone_number" validate:"required"`
		Gender         string   `json:"gender" validate:"required"`
		AmountFamily   string   `json:"amount_family" validate:"required"`
		AvailableStart string   `json:"available_start" validate:"required"`
		AvailableEnd   string   `json:"available_end" validate:"required"`
		Duration       string   `json:"duration" validate:"required"`
	}

	PostEditingReq struct {
		OwnerID        string
		Name           string   `json:"name" validate:"omitempty,max=64"`
		Description    string   `json:"description" validate:"omitempty,max=128"`
		Avatar         string   `json:"avatar" validate:"omitempty"`
		CategoryIDs    []uint64 `json:"category_id" validate:"omitempty"`
		PlaceTypeIDs   []uint64 `json:"placetype_ids" validate:"omitempty"`
		Location       string   `json:"location" validate:"omitempty"`
		LocationLat    string   `json:"latitude" validate:"omitempty"`
		LocationLong   string   `json:"longtitude" validate:"omitempty"`
		Price          float64  `json:"price" validate:"omitempty"`
		PromptPay      string   `json:"prompt_pay" validate:"omitempty,max=128"`
		PhoneNumber    string   `json:"phone_number" validate:"omitempty"`
		Gender         string   `json:"gender" validate:"omitempty"`
		AmountFamily   string   `json:"amount_family" validate:"omitempty"`
		AvailableStart string   `json:"available_start" validate:"omitempty"`
		AvailableEnd   string   `json:"available_end" validate:"omitempty"`
		Duration       string   `json:"duration" validate:"omitempty"`
	}
)
