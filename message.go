package main

import (
	"log"
	"time"

	"github.com/AAA-Intelligence/eve/db"
)

func handleMessage(request MessageRequest) (answer string) {
	// store user message
	err := db.StoreMessage(request.User.ID, db.Message{
		Sender:    db.UserIsSender,
		Bot:       request.Bot,
		Content:   request.Message,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return "oh no i peed my pants"
	}
	answer = "ok, " + request.User.Name

	// store answer
	err = db.StoreMessage(request.User.ID, db.Message{
		Sender:    db.BotIsSender,
		Bot:       request.Bot,
		Content:   request.Message,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return "oh no i peed my pants"
	}
	return
}
