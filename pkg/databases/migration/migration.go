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

	userMigration(tx)
	postCategoryMigration(tx)
	postPlaceTypeMigration(tx)
	postMigration(tx)
	orderMigration(tx)
	historyMigration(tx)
	userRatingMigration(tx)
	notificationMigration(tx)

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

func postCategoryMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PostCategory{})
}

func postPlaceTypeMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PostPlaceType{})
}

func orderMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Order{})
}

func historyMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.History{})
}

func userRatingMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.UserRating{})
}

func notificationMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Notification{})
}
