package manager

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/AAA-Intelligence/eve/manager/bots"
)

// Global bot pool that is used for all requests
var botPool *bots.BotPool

// MessageRequest represents a message request from the user
// The client sends a message as json string.
// The json contains a message and the id of the bot the message is sent to.
type MessageRequest struct {
	Message string `json:"message"`
}

// Takes incoming message requests, sends them to the bot instance and returns the bot's answer
// All messages are stored in the database in the Message table.
// If any error occurs the string "Ok" is returned
func handleMessage(request MessageRequest) string {
	start := time.Now()
	botAnswer := botPool.HandleRequest(bots.MessageData{
		Text:            request.Message,
		Mood:            1,
		Affection:       1,
		Gender:          1,
		Name:            "Emma",
		PreviousPattern: nil,
		Birthdate:       22,
		FavoriteColor:   "Red",
		FatherName:      "Peter",
		FatherAge:       50,
		MotherName:      "Helena",
		MotherAge:       50,
	})

	elapsed := time.Since(start)
	// bot answers should take between 2 and 5 seconds. So we just wait a random time if the bot instance was to quick
	if diff := int64(time.Second*5 - elapsed); diff > 0 {
		min := int64(2*time.Second - elapsed)
		sleepDuration := rand.Int63n(diff-min) + min
		time.Sleep(time.Duration(sleepDuration))
	}
	return botAnswer.Text
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
func httpMessageInterface(w http.ResponseWriter, r *http.Request) {

	if strings.ToLower(r.Method) != "post" {
		http.Error(w, "HTTP POST only", http.StatusMethodNotAllowed)
		return
	}
	var request MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	answer := handleMessage(request)
	fmt.Fprint(w, answer)
}
