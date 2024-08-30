package model

type UserCreatingReq struct {
	ID       string `gorm:"primaryKey;type:varchar(64);"`
	UserName string `gorm:"type:varchar(128);not null;"`
	Avatar   string `gorm:"type:varchar(256);not null;default:'';"`
	Email    string `gorm:"type:varchar(128);unique;not null;"`
}
