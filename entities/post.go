package entities

import (
	"fmt"
	"strconv"
	"time"

	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
)

type Post struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement"`
	Owner          *User         `gorm:"foreignKey:OwnerID"`
	OwnerID        string        `gorm:"type:varchar(64);not null"`
	Name           string        `gorm:"type:varchar(128);not null"`
	Description    string        `gorm:"type:TEXT; not null"`
	Avatar         string        `gorm:"type:TEXT; not null"`
	Location       string        `gorm:"type:varchar(64);not null"`
	LocationLat    string        `gorm:"type:varchar(64);not null"`
	LocationLong   string        `gorm:"type:varchar(64);not null"`
	Price          float64       `gorm:"not null"`
	PromptPay      string        `gorm:"type:varchar(64)"`
	Distance       string        `gorm:"type:varchar(64);not null"`
	DistanceFee    string        `gorm:"type:varchar(64)"`
	PhoneNumber    string        `gorm:"type:varchar(64);not null"`
	Gender         string        `gorm:"type:varchar(64);not null"`
	AmountFamily   string        `gorm:"type:varchar(64);not null"`
	AvailableStart string        `gorm:"type:varchar(64);not null"`
	AvailableEnd   string        `gorm:"type:varchar(64);not null"`
	Duration       string        `gorm:"type:varchar(128);not null"`
	TotalScore     float64       `gorm:"not null"`
	ActiveStatus   bool          `gorm:"not null;default:true"`
	IsReserved     bool          `gorm:"not null;default:false"`
	Banned         bool          `gorm:"not null;default:false"`
	CreatedAt      time.Time     `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time     `gorm:"not null;autoUpdateTime"`
	Categories     []*Category   `gorm:"many2many:post_categories;"`
	PlaceTypes     []*PlaceType  `gorm:"many2many:post_place_types;"`
	UserRatings    []*UserRating `gorm:"foreignKey:WorkerPostID"`
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

	var categories []_postModel.Category
	for _, category := range p.Categories {
		categories = append(categories, _postModel.Category{
			ID:          category.ID,
			GroupID:     category.GroupID,
			Name:        category.Name,
			Description: category.Description,
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
		ID:             p.ID,
		OwnerID:        p.OwnerID,
		Name:           p.Name,
		Description:    p.Description,
		Avatar:         p.Avatar,
		Price:          p.Price,
		PromptPay:      p.PromptPay,
		Distance:       p.Distance,
		DistanceFee:    p.DistanceFee,
		Location:       p.Location,
		LocationLat:    p.LocationLat,
		LocationLong:   p.LocationLong,
		TotalScore:     summaryTotalScore,
		PhoneNumber:    p.PhoneNumber,
		Gender:         p.Gender,
		AmountFamily:   p.AmountFamily,
		Duration:       p.Duration,
		AvailableStart: p.AvailableStart,
		AvailableEnd:   p.AvailableEnd,
		ActiveStatus:   p.ActiveStatus,
		Categories:     categories,
		PlaceTypes:     placeTypes,
		UserRatings:    userRatings,
		CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02 15:04"),
	}
}
