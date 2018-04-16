package manager

import (
	"log"
	"net/http"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/gorilla/websocket"
)

// a simple upgrader to establish a WebSocket from a HTTP request
var upgrader = websocket.Upgrader{}

// MessageRequest represents a message request from the user
// The client sends a message as json string.
// The json contains a message and the id of the bot the message is sent to.
type MessageRequest struct {
	Message string `json:"message"`
	Bot     int    `json:"bot"`
	User    *db.User
}

// webSocket upgrades the HTTP request to the WebSocket protocol.
// all messages sent between the server and clint are communicated over the websocket.
func webSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading websocket:", err)
		return
	}
	defer c.Close()
	// wait for messsages
	for {
		var request MessageRequest
		err := c.ReadJSON(&request)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println("error reading:", err)
			}
			break
		}
		user := GetUserFromRequest(r)
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
