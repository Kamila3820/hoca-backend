package main

import (
	"fmt"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/Kamila3820/hoca-backend/server"
)

func main() {
	conf := config.ConfigGetting()
	db := databases.NewPostgresDatabase(conf.Database)
	server := server.NewEchoServer(conf, db)

	google := conf.Google.ApiKey

	fmt.Println("Hello, Test!")
	fmt.Println(google)
	server.Start()
}
