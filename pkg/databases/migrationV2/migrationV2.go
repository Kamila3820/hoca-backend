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

	// userAdding(tx)
	postCategoryAdding(tx)
	// placeTypeAdding(tx)
	// postsAdding(tx)
	// userRatingAdding(tx)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(tx.Error)
	}
}

func postsAdding(tx *gorm.DB) {
	posts := []entities.Post{
		// {
		// 	ID:           44,
		// 	OwnerID:      "33bc9842-2801-41de-814d-bbd7fb172a1b",
		// 	Name:         "2 Test Service",
		// 	Description:  "This is a test service.",
		// 	Avatar:       "", // Empty string
		// 	Location:     "Test Location",
		// 	LocationLat:  "13.7563",
		// 	LocationLong: "100.5018",
		// 	Price:        100.00,
		// 	Distance:     "5.0",
		// 	PhoneNumber:  "0123456789",
		// 	Gender:       "Male",
		// 	AmountFamily: "4",
		// 	TotalScore:   9.5,
		// 	ActiveStatus: true,
		// },
		{
			ID:           55,
			OwnerID:      "205118840060769110477",
			Name:         "3 Test Service",
			Description:  "This is a test service.",
			Avatar:       "", // Empty string
			Location:     "Test Location",
			LocationLat:  "13.7563",
			LocationLong: "100.4900",
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
	postCategory := []entities.Category{
		// {
		// 	ID:          1,
		// 	GroupID:     1,
		// 	Name:        "Deep cleaning",
		// 	Description: "Cleaning",
		// },
		// {
		// 	ID:          4,
		// 	GroupID:     2,
		// 	Name:        "Laundry",
		// 	Description: "Clothes",
		// },
		// {
		// 	ID:          9,
		// 	GroupID:     4,
		// 	Name:        "Pet sitting",
		// 	Description: "Pets",
		// },
		{
			ID:          6,
			GroupID:     3,
			Name:        "Lawn Mowing",
			Description: "Garden",
		},
		{
			ID:          8,
			GroupID:     3,
			Name:        "Yard cleanup",
			Description: "Garden",
		},
	}

	tx.CreateInBatches(postCategory, len(postCategory))
}

func placeTypeAdding(tx *gorm.DB) {
	postPlace := []entities.PlaceType{
		{
			Name:        "House",
			Description: "House Home",
		},
		{
			Name:        "Room & Condo",
			Description: "Room & Condo",
		},
		{
			Name:        "Dormitory",
			Description: "Dormitory",
		},
	}

	tx.CreateInBatches(postPlace, len(postPlace))
}

func userAdding(tx *gorm.DB) {
	users := []entities.User{
		{
			ID:           "205118840060769110477",
			UserName:     "Fang",
			FirstName:    "Attapin",
			LastName:     "Pinya",
			Email:        "attpinya@gmail.com",
			Avatar:       "https://lh3.googleusercontent.com/a/ACg8ocLyvY_troho1V-6qhTv6gyWrBKoOUcZwI9VCd6EUYc7MpURVgMQ=s96-c",
			Password:     "your_secure_password", // Add a placeholder or hash the password
			PhoneNumber:  "0258505514",
			IDCard:       "1234567890123", // Add a unique ID card value
			VerifyStatus: false,
			Latitude:     "",
			Longtitude:   "",
		},
	}

	tx.CreateInBatches(users, len(users))
}

func userRatingAdding(tx *gorm.DB) {
	userRating := []entities.UserRating{
		{
			ID:            30,
			UserID:        "205118840060769110477",
			WorkerPostID:  "44",
			WorkScore:     3,
			SecurityScore: 0,
			Comment:       "Terrible work",
		},
	}

	tx.CreateInBatches(userRating, len(userRating))
}
