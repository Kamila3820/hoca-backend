package entities

import (
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
)

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

func (u *User) ToUserModel() *_userModel.User {
	return &_userModel.User{
		ID:           u.ID,
		UserName:     u.UserName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Avatar:       u.Avatar,
		Password:     u.Password,
		PhoneNumber:  u.PhoneNumber,
		IDCard:       u.IDCard,
		VerifyStatus: u.VerifyStatus,
		Location:     u.Location,
		Latitude:     u.Latitude,
		Longtitude:   u.Longtitude,
	}
}
