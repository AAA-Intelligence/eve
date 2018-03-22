package main

import (
	"log"
	"os"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/AAA-Intelligence/eve/manager"
)

const dbFile = "eve.sqlite"

func main() {
	config := loadConfig()
	// check if db file exists
	if _, err := os.Stat(dbFile); err != nil {
		log.Fatalln("cannot find database file:", dbFile)
		return
	}
	err := db.Connect(dbFile)
	if err != nil {
		log.Fatalln("error connecting to database: ", err)
		return
	}
	manager.StartWebServer(config.Host, config.HTTP)
}
