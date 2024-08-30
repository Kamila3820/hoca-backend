package main

import (
	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
)

func main() {
	conf := config.ConfigGetting()
	db := databases.NewPostgresDatabase(conf.Database)

	userMigration(db)
	oauthMigration(db)
	postMigration(db)
	postCategoryMigration(db)
	postPlaceTypeMigration(db)
	orderMigration(db)
	historyMigration(db)
	userRatingMigration(db)
	notificationMigration(db)
}

func userMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.User{})
}

func oauthMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.Oauth{})
}

func postMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.Post{})
}

func postCategoryMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.PostCategory{})
}

func postPlaceTypeMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.PostPlaceType{})
}

func orderMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.Order{})
}

func historyMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.History{})
}

func userRatingMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.UserRating{})
}

func notificationMigration(db databases.Database) {
	db.Connect().Migrator().CreateTable(&entities.Notification{})
}
