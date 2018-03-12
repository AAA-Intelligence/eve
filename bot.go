package main

import (
	"net/http"

	"github.com/drhodes/golorem"

	"github.com/AAA-Intelligence/eve/db"
)

func createBot(w http.ResponseWriter, r *http.Request) {
	err := db.CreateBot(&db.Bot{
		Name:   lorem.Word(3, 10),
		Image:  "hässlich.png",
		Gender: "apache",
		User:   getUser(r.Context()).ID,
	})
	if err != nil {
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
