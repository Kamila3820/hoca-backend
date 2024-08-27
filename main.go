package main

import (
	"fmt"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/databases"
	"github.com/Kamila3820/hoca-backend/server"
)

func main() {
	conf := config.ConfigGetting()
	db := databases.NewPostgresDatabase(conf.Database)
	server := server.NewEchoServer(conf, db.ConnectionGetting())

	fmt.Println("Hello, World!")
	server.Start()
}
