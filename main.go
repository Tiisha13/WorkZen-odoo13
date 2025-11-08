package main

import (
	"log"

	"api.workzen.odoo/config"
	"api.workzen.odoo/constants"
	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/seed"
	"api.workzen.odoo/routers"
	"github.com/Delta456/box-cli-maker/v2"
)

func init() {
	config.AppConfig = config.GetConfig()
	boxConfig := box.Config{
		Px:           10,
		Py:           1,
		Type:         "Round",
		Color:        "Yellow",
		TitleColor:   "Cyan",
		ContentAlign: "Center",
	}

	Box := box.New(boxConfig)

	if _, err := databases.InitDB(); err != nil {
		Box.Print("WorkZen - Backend Server", "Database Connection Failed!")
		panic(err)
	} else {
		Box.Print("WorkZen - Backend Server Started!", "Database Connected! \nServer Mode: "+constants.ServerMode)

		// Run database seeding
		if err := seed.SeedDatabase(databases.GetMongoDBDatabase()); err != nil {
			log.Printf("⚠️  Warning: Database seeding failed: %v\n", err)
		}
	}
}

func main() {
	defer func() {
		if _, err := databases.CloseDB(); err != nil {
			panic(err)
		}
	}()

	router := routers.Init()
	router.Listen(constants.ServerPort)
}
