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

	// userMigration(tx)
	// categoryMigration(tx)
	// postCategoryMigration(tx)
	// placeTypeMigration(tx)
	// postPlaceTypeMigration(tx)
	// postMigration(tx)
	// orderMigration(tx)
	// orderQRMigration(tx)
	// historyMigration(tx)
	// userRatingMigration(tx)
	// notificationMigration(tx)
	workerFeeMigration(tx)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(tx.Error)
	}
}

func userMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.User{})
}

func postMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Post{})
}

func categoryMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Category{})
}

func postCategoryMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PostCategory{})
}

func placeTypeMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PlaceType{})
}

func postPlaceTypeMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PostPlaceType{})
}

func orderMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Order{})
}

func orderQRMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.OrderQrpayment{})
}

func historyMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.History{})
}

func userRatingMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.UserRating{})
}

func notificationMigration(tx *gorm.DB) {
	// enumCreateQuery := `
	// CREATE TYPE notification_enum AS ENUM (
	//     'confirmation',
	//     'preparing',
	//     'working',
	//     'complete',
	//     'user_cancel',
	//     'worker_cancel',
	//     'user_rating',
	// 	'system_cancel'
	// );
	// `

	// if err := tx.Exec(enumCreateQuery).Error; err != nil {
	// 	log.Printf("Error creating ENUM type: %v", err)
	// }

	tx.Migrator().CreateTable(&entities.Notification{})
}

func workerFeeMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.WorkerFeepayment{})
}
