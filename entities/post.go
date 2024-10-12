package entities

import (
	"fmt"
	"strconv"
	"time"

	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
)

type Post struct {
	ID           uint64        `gorm:"primaryKey;autoIncrement"`
	Owner        *User         `gorm:"foreignKey:OwnerID"`
	OwnerID      string        `gorm:"type:varchar(64);not null"`
	Name         string        `gorm:"type:varchar(128);not null"`
	Description  string        `gorm:"type:varchar(256);not null"`
	Avatar       string        `gorm:"type:varchar(256);not null;default:'';"`
	CategoryID   uint64        `gorm:"not null"` // Foreign key to PostCategory
	Location     string        `gorm:"type:varchar(64);not null"`
	LocationLat  string        `gorm:"type:varchar(64);not null"`
	LocationLong string        `gorm:"type:varchar(64);not null"`
	Price        float64       `gorm:"not null"`
	Distance     string        `gorm:"type:varchar(64);not null"`
	PhoneNumber  string        `gorm:"type:varchar(64);not null"`
	Gender       string        `gorm:"type:varchar(64);not null"`
	AmountFamily string        `gorm:"type:varchar(64);not null"`
	TotalScore   float64       `gorm:"not null"`
	ActiveStatus bool          `gorm:"not null;default:true"`
	IsReserved   bool          `gorm:"not null;default:false"`
	CreatedAt    time.Time     `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"not null;autoUpdateTime"`
	Category     *PostCategory `gorm:"foreignKey:CategoryID"`
	PlaceTypes   []*PlaceType  `gorm:"many2many:post_place_types;"`
	UserRatings  []*UserRating `gorm:"foreignKey:WorkerPostID"`
}

func (p *Post) ToPostModel() *_postModel.Post {
	var placeTypes []_postModel.PlaceType
	for _, pt := range p.PlaceTypes {
		placeTypes = append(placeTypes, _postModel.PlaceType{
			ID:          pt.ID,
			Name:        pt.Name,
			Description: pt.Description,
		})
	}

	totalWorkScore := 0.0
	totalSecurityScore := 0.0
	count := len(p.UserRatings)

	var userRatings []_postModel.UserRating
	for _, ur := range p.UserRatings {
		userRatings = append(userRatings, _postModel.UserRating{
			ID:            ur.ID,
			UserID:        ur.UserID,
			Username:      ur.User.UserName,
			Avatar:        ur.User.Avatar,
			WorkScore:     ur.WorkScore,
			SecurityScore: ur.SecurityScore,
			Comment:       ur.Comment,
			CreatedAt:     ur.CreatedAt.Format("2006-01-02 15:04"),
		})

		totalWorkScore += float64(ur.WorkScore)
		totalSecurityScore += float64(ur.SecurityScore)
	}

	// Compute average total score
	var totalScore float64
	if count > 0 {
		totalScore = (totalWorkScore + totalSecurityScore) / float64(2*count)
	}

	summaryTotalScore, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", totalScore), 64)

	return &_postModel.Post{
		ID:           p.ID,
		OwnerID:      p.OwnerID,
		Name:         p.Name,
		Description:  p.Description,
		Avatar:       p.Avatar,
		CategoryID:   p.CategoryID,
		Price:        p.Price,
		Distance:     p.Distance,
		Location:     p.Location,
		LocationLat:  p.LocationLat,
		LocationLong: p.LocationLong,
		TotalScore:   summaryTotalScore,
		PhoneNumber:  p.PhoneNumber,
		Gender:       p.Gender,
		AmountFamily: p.AmountFamily,
		ActiveStatus: p.ActiveStatus,
		PlaceTypes:   placeTypes,
		UserRatings:  userRatings,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
