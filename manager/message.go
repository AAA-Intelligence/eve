package manager

import (
	"log"
	"time"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/AAA-Intelligence/eve/manager/bots"
)

// Global bot pool that is used for all requests
var botPool *bots.BotPool

// Takes incoming message requests, sends them to the bot instance and returns the bot's answer
// All messages are stored in the database in the Message table.
// If any error occurs the string "Ok" is returned
func handleMessage(request MessageRequest) string {
	bot, err := db.GetBot(request.Bot, request.User.ID)
	if err != nil {
		log.Println("error loading bot data from db:", err)
		return "Ok"
	}
	botAnswer := botPool.HandleRequest(bots.MessageData{
		Text:            request.Message,
		Mood:            bot.Mood,
		Affection:       bot.Affection,
		Gender:          bot.Gender,
		Name:            bot.Name,
		PreviousPattern: nil,
		Birthdate:       bot.Birthdate.Unix(),
		FavoriteColor:   bot.GetFavoriteColor(),
		FatherName:      bot.GetFatherName(),
		FatherAge:       bot.FatherAge,
		MotherName:      bot.GetMotherName(),
		MotherAge:       bot.MotherAge,
	})
	if err = bot.UpdateContext(botAnswer.Affection, botAnswer.Mood); err != nil {
		log.Println("error updating bot:", err)
	}

	// store sent messages
	err = db.StoreMessages(request.User.ID, bot.ID, []db.Message{
		db.Message{
			Sender:    db.UserIsSender,
			Content:   request.Message,
			Timestamp: time.Now(),
		},
		db.Message{
			Sender:    db.BotIsSender,
			Content:   botAnswer.Text,
			Timestamp: time.Now(),
		},
	})
	if err != nil {
		log.Println("error storing message:", err)
		return "Ok"
	}
	return botAnswer.Text
}
