package main

import (
	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"gorm.io/gorm"
)

func main() {
	conf := config.ConfigGetting()
	db := databases.NewPostgresDatabase(conf.Database)

	tx := db.Connect().Begin()

	// postCategoryAdding(tx)
	// postPlaceTypeAdding(tx)
	postsAdding(tx)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(tx.Error)
	}
}

func postsAdding(tx *gorm.DB) {
	posts := []entities.Post{
		{
			ID:           4444,
			OwnerID:      "105118840060769110477",
			Name:         "2 Test Service",
			Description:  "This is a test service.",
			Avatar:       "", // Empty string
			CategoryID:   1,
			PlaceTypeID:  1,
			Location:     "Test Location",
			LocationLat:  "14.7563",
			LocationLong: "100.5018",
			Price:        100.00,
			Distance:     "5.0",
			PhoneNumber:  "0123456789",
			Gender:       "Male",
			AmountFamily: "4",
			TotalScore:   9.5,
			ActiveStatus: true,
		},
		{
			ID:           5555,
			OwnerID:      "105118840060769110477",
			Name:         "3 Test Service",
			Description:  "This is a test service.",
			Avatar:       "", // Empty string
			CategoryID:   1,
			PlaceTypeID:  1,
			Location:     "Test Location",
			LocationLat:  "14.7563",
			LocationLong: "102.5018",
			Price:        100.00,
			Distance:     "5.0",
			PhoneNumber:  "0123456789",
			Gender:       "Male",
			AmountFamily: "4",
			TotalScore:   9.5,
			ActiveStatus: true,
		},
	}

	tx.CreateInBatches(posts, len(posts))
}

func postCategoryAdding(tx *gorm.DB) {
	postCategory := []entities.PostCategory{
		{
			ID:          1,
			Name:        "Cleaning",
			Description: "House cleaning",
		},
		{
			ID:          2,
			Name:        "Clothes",
			Description: "Clothes care",
		},
		{
			ID:          3,
			Name:        "Pets",
			Description: "Pet sitting",
		},
		{
			ID:          4,
			Name:        "Gaedening",
			Description: "House gardening",
		},
	}

	tx.CreateInBatches(postCategory, len(postCategory))
}

func postPlaceTypeAdding(tx *gorm.DB) {
	postPlace := []entities.PostPlaceType{
		{
			ID:          1,
			Name:        "House",
			Description: "House Home",
		},
		{
			ID:          2,
			Name:        "Room & Condo",
			Description: "Room & Condo",
		},
		{
			ID:          3,
			Name:        "Dormitory",
			Description: "Dormitory",
		},
	}

	tx.CreateInBatches(postPlace, len(postPlace))
}
