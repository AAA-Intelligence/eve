package manager

import (
	"log"
	"net/http"
	"strconv"
	"github.com/rs/cors"
	"github.com/AAA-Intelligence/eve/manager/bots"
)

// if the webserver is shut down, all bot instances are killed
// and the connection to the database is closed
func onShutdown() {
	log.Println("shutting down...")
	log.Println("killing bots...")
	botPool.Close()
	// wait until all bots finished their running tasks
	botPool.Wait()
	log.Println("shutdown complete")
}

// StartWebServer creates a handler for incomming http requests on the given host and port
// The method only returns if the server is shut down or runs into an error
func StartWebServer(host string, httpPort int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/message-api", httpMessageInterface)
	handler := cors.Default().Handler(mux)
	server := http.Server{
		Addr:    host + ":" + strconv.Itoa(httpPort),
		Handler: handler,
	}
	//go startBot()

	log.Println("Starting web server")
	server.RegisterOnShutdown(onShutdown)
	// start as many bot instances as cpu has cores
	botPool = bots.NewBotPool(1)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
