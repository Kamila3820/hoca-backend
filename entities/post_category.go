package entities

type PostCategory struct {
	PostID     uint64 `gorm:"primaryKey;autoIncrement:false"` // Foreign key to Post
	CategoryID uint64 `gorm:"primaryKey;autoIncrement:false"` // Foreign key to Category
}
