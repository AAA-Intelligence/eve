package main

import (
	"math/rand"
	"time"

	"github.com/AAA-Intelligence/eve/manager"
)

// usage: eve -host [hostname] -http [http port] -https [https port]
func main() {
	// get web server config from program arguments
	config := loadConfig()
	// set seed for random generations
	rand.Seed(time.Now().Unix())

	manager.StartWebServer(config.Host, config.HTTP)
}
