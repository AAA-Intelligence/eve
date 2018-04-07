package main

import (
	"log"
	"os"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/AAA-Intelligence/eve/manager"
)

// database filename (absolute or relative to working directory)
const dbFile = "eve.sqlite"

// usage: eve -host [hostname] -http [http port] -https [https port]
func main() {
	// get web server config from program arguments
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
