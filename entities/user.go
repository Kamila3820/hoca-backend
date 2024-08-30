package entities

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
	Latitude     string `gorm:"type:varchar(64);default:''"`
	Longtitude   string `gorm:"type:varchar(64);default:''"`
}
