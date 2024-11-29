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

type User struct {
	ID           string `json:"id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	PhoneNumber  string `json:"phone_number"`
	VerifyStatus bool   `json:"verify_status"`
	Location     string `json:"location"`
	Latitude     string `json:"latitude"`
	Longtitude   string `json:"longtitude"`
}

type UserEditingReq struct {
	ID          string
	UserName    string `json:"user_name" validate:"omitempty,max=64"`
	Avatar      string `json:"avatar" validate:"omitempty,max=64"`
	Email       string `json:"email" validate:"omitempty,max=128"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,max=64"`
}
