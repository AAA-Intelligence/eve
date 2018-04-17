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
	bot, err := request.User.GetBot(request.Bot)
	if err != nil {
		log.Println("error loading bot data from db:", err)
		return "Ok"
	}

	err = bot.StoreMessages(request.User, []db.Message{
		db.Message{
			Sender:    db.UserIsSender,
			Content:   request.Message,
			Timestamp: time.Now(),
		},
	})
	botAnswer := botPool.HandleRequest(bots.MessageData{
		Text:            request.Message,
		Mood:            bot.Mood,
		Affection:       bot.Affection,
		Gender:          bot.Gender,
		Name:            bot.Name,
		PreviousPattern: bot.Pattern,
		Birthdate:       bot.Birthdate.Unix(),
		FavoriteColor:   bot.GetFavoriteColor(),
		FatherName:      bot.GetFatherName(),
		FatherAge:       bot.FatherAge,
		MotherName:      bot.GetMotherName(),
		MotherAge:       bot.MotherAge,
	})
	if err = bot.UpdateContext(botAnswer.Affection, botAnswer.Mood, botAnswer.Pattern); err != nil {
		log.Println("error updating bot:", err)
	}

	// store sent messages
	err = bot.StoreMessages(request.User, []db.Message{
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
