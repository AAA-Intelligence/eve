package main

import (
	"log"
	"net/http"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// MessageRequest represents a message request from the user
type MessageRequest struct {
	Message string `json:"message"`
	Bot     int    `json:"bot"`
	User    *db.User
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var request MessageRequest
		err := c.ReadJSON(&request)
		if err != nil {
			//log.Println("error reading:", err)
			break
		}
		user := getUser(r.Context())
		if user == nil {
			log.Println("error reading user from context")
			break
		}
		request.User = user
		answer := handleMessage(request)
		err = c.WriteMessage(websocket.TextMessage, []byte(answer))
		if err != nil {
			log.Println("error writing:", err)
			break
		}
	}
}
