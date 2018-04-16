package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/AAA-Intelligence/eve/manager"
)

// database filename (absolute or relative to working directory)
const dbFile = "eve.sqlite"

// usage: eve -host [hostname] -http [http port] -https [https port]
func main() {
	// get web server config from program arguments
	config := loadConfig()
	err := db.Connect(dbFile)
	if err != nil {
		log.Fatalln("error connecting to database: ", err)
		return
	}
	// set seed for random generations
	rand.Seed(time.Now().Unix())

	manager.StartWebServer(config.Host, config.HTTP)
}
