package model

type RegisterRequest struct {
	Username        *string `json:"username" validate:"required"`
	Password        *string `json:"password" validate:"required"`
	PhoneNumber     *string `json:"phone_number" validate:"required"`
	Email           *string `json:"email" validate:"required"`
	ConfirmPassword *string `json:"confirmPassword" validate:"required"`
	IDcard          *string `json:"id_card"`
}

type RegisterResponse struct {
	UserID *string `json:"userID"`
}

type LoginRequest struct {
	Email    *string `json:"email" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token *string `json:"token"`
}
