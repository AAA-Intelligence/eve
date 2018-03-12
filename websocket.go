package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func webSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			//log.Println("error reading:", err)
			break
		}

		// cut messages longer than 200 characters
		if len(message) > 200 {
			message = message[:200]
		}

		answer := handleMessage(string(message), GetUserFromRequest(r))
		err = c.WriteMessage(mt, []byte(answer))
		if err != nil {
			log.Println("error writing:", err)
			break
		}
	}
}
