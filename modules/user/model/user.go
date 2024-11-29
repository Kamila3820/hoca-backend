package model

type User struct {
	ID           string `gorm:"primaryKey;type:varchar(64);"`
	UserName     string `gorm:"type:varchar(128);not null;"`
	FirstName    string `gorm:"type:varchar(128);not null;"`
	LastName     string `gorm:"type:varchar(128);not null;"`
	Email        string `gorm:"type:varchar(128);unique;not null;"`
	Avatar       string `gorm:"type:varchar(256);not null;default:'';"`
	Password     string `gorm:"type:varchar(64);not null;"`
	PhoneNumber  string `gorm:"type:varchar(64);unique;not null;"`
	IDCard       string `gorm:"type:varchar(128);unique;not null;"`
	VerifyStatus bool   `gorm:"not null;"`
	Location     string `gorm:"type:varchar(64);not null"`
	Latitude     string `gorm:"type:varchar(64);default:''"`
	Longtitude   string `gorm:"type:varchar(64);default:''"`
}

type UserCreatingReq struct {
	ID       string `gorm:"primaryKey;type:varchar(64);"`
	UserName string `gorm:"type:varchar(128);not null;"`
	Avatar   string `gorm:"type:varchar(256);not null;default:'';"`
	Email    string `gorm:"type:varchar(128);unique;not null;"`
}

type ProfileUser struct {
	ID          string `json:"id"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
}

type UserEditingReq struct {
	ID          string
	UserName    string `json:"user_name" validate:"omitempty,max=64"`
	Avatar      string `json:"avatar" validate:"omitempty,max=64"`
	Email       string `json:"email" validate:"omitempty,max=128"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,max=64"`
}

type UserLocation struct {
	Location   string `json:"location"`
	Latitude   string `json:"latitude"`
	Longtitude string `json:"longtitude"`
}

type UserLocationReq struct {
	Location   string `json:"location"`
	Latitude   string `json:"latitude"`
	Longtitude string `json:"longtitude"`
}
