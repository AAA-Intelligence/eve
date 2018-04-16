package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		request.User = GetUserFromRequest(r)
		if request.User == nil {
			log.Println("error reading user from context")
			break
		}
		answer := handleMessage(request)
		err = c.WriteMessage(websocket.TextMessage, []byte(answer))
		if err != nil {
			log.Println("error writing:", err)
			break
		}
	}
}

// Answers messages via http requests.
// A request can be made with a HTTP POST request.
// The body must be a json which contains a message and the bot id.
// Example for request body:
// 	{
//		"message":"message string",
//		"bot": 1
//	}
// The HTTP response body contains the answer as plain text.
//
// IMPORTANT!
// Authentification is needed!
// HTTP header musst contain basic auth credentials or a session key as cookie.
// see basicAuth(...) for more information
func httpMessageInterface(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "POST" {
		http.Error(w, "HTTP POST only", http.StatusMethodNotAllowed)
		return
	}
	var request MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	request.User = GetUserFromRequest(r)
	if request.User == nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	answer := handleMessage(request)
	fmt.Fprint(w, answer)
}
