package entities

type Category struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	GroupID     uint64 `gorm:"not null"`
	Name        string `gorm:"type:varchar(128);not null;unique"`
	Description string `gorm:"type:varchar(128)"`
}
