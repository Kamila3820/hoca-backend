package entities

type PostPlaceType struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(128);not null;unique"`
	Description string `gorm:"type:text"`
}