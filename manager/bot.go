package manager

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
		User:   GetUserFromRequest(r).ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
