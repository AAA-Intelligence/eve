package main

import (
	"log"

	"github.com/AAA-Intelligence/eve/manager"

	"github.com/AAA-Intelligence/eve/db"
)

func main() {
	config := loadConfig()
	err := db.Connect("eve.sqlite")
	if err != nil {
		log.Panic("error connecting to database: ", err)
		return
	}
	manager.StartWebServer(config.Host, config.HTTP)
}
