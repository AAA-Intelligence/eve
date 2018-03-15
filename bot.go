package main

import (
	"encoding/json"
	"net/http"

	"github.com/drhodes/golorem"

	"github.com/AAA-Intelligence/eve/db"
)

func createBot(w http.ResponseWriter, r *http.Request) {
	err := db.CreateBot(&db.Bot{
		Name:   lorem.Word(3, 10),
		Image:  "hässlich.png",
		Gender: "apache",
		User:   GetUserFromRequest(r).ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	messages, err := db.GetMessagesForUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(*messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
