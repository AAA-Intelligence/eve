package manager

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/AAA-Intelligence/eve/manager/bots"

	"github.com/AAA-Intelligence/eve/db"
)

//indexHandler serves HTML index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// make sure request is really for index page
	if len(r.URL.Path) > 1 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	user := GetUserFromRequest(r)
	bot := GetBotFromRequest(r)
	messages := &[]db.Message{}
	if bot != nil {
		if msgs, err := db.GetMessagesForBot(bot.ID); err == nil {
			messages = msgs
		} else {
			log.Println(err)
		}

	}
	tpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading template:", err.Error())
		return
	}
	bots, err := db.GetBotsForUser(user.ID)
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error getting bots for user:", err.Error())
		return
	}
	err = saveExecute(w, tpl, struct {
		User      *db.User
		Bots      *[]db.Bot
		ActiveBot *db.Bot
		Messages  *[]db.Message
	}{
		User:      user,
		Bots:      bots,
		ActiveBot: bot,
		Messages:  messages,
	})
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error executing template:", err.Error())
		return
	}
}

//registerHandler serves HTML page for user registration
func registerHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/register.gohtml")
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading template:", err)
		return
	}
	err = saveExecute(w, tpl, nil)
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error executing template:", err)
		return
	}
}

func getRandomName(res http.ResponseWriter, req *http.Request) {
	log.Println("Trigger")
	sex := req.URL.Query().Get("sex")
	sexId, err := strconv.Atoi(sex)
	if err != nil {
		http.Error(res, "invalid sex", http.StatusBadRequest)
		return
	}
	names, err := db.GetNames(sexId)
	if err != nil || len(*names) < 1 {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading names")
		return
	}

	tpl, err := template.ParseFiles("templates/register.gohtml")

	err = saveExecute(res, tpl, struct {
		Name db.Name
	}{
		Name: (*names)[rand.Intn(len(*names))],
	})
	if err != nil {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error executing template:", err)
		return
	}

}

func createBot(res http.ResponseWriter, req *http.Request) {

	// tpl,err2 := template.ParseFiles("templates/index.gohtml")
	// if err2 != nil {
	// 	http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
	// 	log.Println("error loading template:", err2.Error())
	// 	return
	// }
	// err2 = saveExecute(w, tpl, nil)
	// if err2 != nil {
	// 	http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
	// 	log.Println("error executing template:", err2)
	// 	return
	// }

	err := db.CreateBot(&db.Bot{
		Name:   names[rand.Intn(9)],
		Image:  "hässlich.png",
		Gender: db.Female,
		User:   GetUserFromRequest(r).ID,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func onShutdown() {
	log.Println("shutting down...")
	log.Println("killing bots...")
	botPool.Close()
	botPool.Wait()
	log.Println("closing database connection...")
	err := db.Close()
	if err != nil {
		log.Panic("error closing connection to database: ", err)
		return
	}
	log.Println("shutdown complete")
}

// StartWebServer creates a handler for incomming http requests on the given host and port
// The method only returns if the server is shut down or runs into an error
func StartWebServer(host string, httpPort int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/createUser", createUser)

	mux.HandleFunc("/", basicAuth(indexHandler))
	mux.HandleFunc("/createBot", basicAuth(createBot))
	mux.HandleFunc("/getRandomName", basicAuth(getRandomName))
	mux.HandleFunc("/ws", basicAuth(webSocket))

	// handle static files like css
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := http.Server{
		Addr:    host + ":" + strconv.Itoa(httpPort),
		Handler: mux,
	}
	//go startBot()

	log.Println("Starting web server")
	server.RegisterOnShutdown(onShutdown)
	botPool = bots.NewBotPool(4)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// GetBotFromRequest checks if there is a bot id in the request (HTTP GET e.g. ?bot=2)
// If not the first bot for the user is returned
// If an error occures or there is no bot in the request or database nil is returned
func GetBotFromRequest(r *http.Request) *db.Bot {
	user := GetUserFromRequest(r)
	if user == nil {
		return nil
	}
	idString := r.URL.Query().Get("bot")
	bots, err := db.GetBotsForUser(user.ID)
	if err != nil || len(*bots) < 1 {
		return nil
	}
	if len(idString) < 1 {
		return &(*bots)[0]
	}
	if id, err := strconv.Atoi(idString); err == nil {
		for _, b := range *bots {
			if b.ID == id {
				return &b
			}
		}
	}
	return nil
}
