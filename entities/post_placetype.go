package entities

type PostPlaceType struct {
	PostID      uint64 `gorm:"primaryKey;autoIncrement:false"` // Foreign key to Post
	PlaceTypeID uint64 `gorm:"primaryKey;autoIncrement:false"` // Foreign key to PlaceType
}
