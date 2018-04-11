package manager

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/AAA-Intelligence/eve/db"
	"github.com/AAA-Intelligence/eve/manager/bots"
	"github.com/gorilla/schema"
)

// indexHandler serves HTML index page
// template file: templates/index.gohtml
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// make sure request is really for index page
	if len(r.URL.Path) > 1 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	user := GetUserFromRequest(r)
	bot := GetBotFromRequest(r)

	// load all messages sent between user and bot
	messages := &[]db.Message{}
	if bot != nil {
		if msgs, err := db.GetMessagesForBot(bot.ID); err == nil {
			messages = msgs
		} else {
			log.Println(err)
		}
	}
	tpl := template.New("index").Funcs(template.FuncMap{
		"time":  formatTime,
		"years": yearsSince,
	})
	tpl, err := tpl.ParseFiles("templates/index.gohtml")
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
	err = saveExecute(w, tpl.Lookup("index.gohtml"), struct {
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

// registerHandler serves HTML page for user registration
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

	sex := req.URL.Query().Get("sex")
	sexID, err := strconv.Atoi(sex)
	if err != nil {
		http.Error(res, "invalid sex", http.StatusBadRequest)
		return
	}
	names, err := db.GetNames(sexID)
	if err != nil || len(*names) < 1 {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading names")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode((*names)[rand.Intn(len(*names))])

}

func getRandomImage(res http.ResponseWriter, req *http.Request) {
	sex := req.URL.Query().Get("sex")
	sexID, err := strconv.Atoi(sex)
	if err != nil {
		http.Error(res, "invalid sex", http.StatusBadRequest)
		return
	}
	images, err := db.GetImages(sexID)
	if err != nil || len(*images) < 1 {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading names")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode((*images)[rand.Intn(len(*images))])

}

func getImages(res http.ResponseWriter, req *http.Request) {
	sex := req.URL.Query().Get("sex")
	sexID, err := strconv.Atoi(sex)
	if err != nil {
		http.Error(res, "invalid sex", http.StatusBadRequest)
		return
	}
	images, err := db.GetImages(sexID)
	if err != nil || len(*images) < 1 {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading names")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode((*images))

}

type Params struct {
	Name  int       `schema:"nameID"`
	Image int       `schema:"imageID"`
	Sex   db.Gender `schema:"sex"`
}

var decoder = schema.NewDecoder()

func createBot(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("error parsing form:", err)
		http.Error(res, "invalid request", http.StatusBadRequest)
		return
	}

	var params Params

	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&params, req.PostForm)
	if err != nil {
		log.Println("error decoding form:", err)
		http.Error(res, "invalid request", http.StatusBadRequest)
		return
	}
	name, err := db.GetName(params.Name)
	if err != nil {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading name")
		return
	}

	image, err := db.GetImage(params.Image)
	if err != nil {
		http.Error(res, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading image")
		return
	}

	err = db.CreateBot(&db.Bot{
		Name:   name.Text,
		Image:  image.Path,
		Gender: params.Sex,
		User:   GetUserFromRequest(req).ID,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// if the webserver is shut down, all bot instances are killed
// and the connection to the database is closed
func onShutdown() {
	log.Println("shutting down...")
	log.Println("killing bots...")
	botPool.Close()
	// wait until all bots finished their running tasks
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
	mux.HandleFunc("/getRandomImage", basicAuth(getRandomImage))
	mux.HandleFunc("/getImages", basicAuth(getImages))
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
// It also checks if the bot belongs to the authenticated user
// If no bot id is provided in the request the first bot for the user is returned
// If an error occures or there is no bot in the request or database nil is returned
func GetBotFromRequest(r *http.Request) *db.Bot {
	user := GetUserFromRequest(r)
	if user == nil {
		// user not authenticated
		return nil
	}
	idString := r.URL.Query().Get("bot")
	bots, err := db.GetBotsForUser(user.ID)
	if err != nil || len(*bots) < 1 {
		// user has no bots
		return nil
	}
	if len(idString) < 1 {
		// return first bot that belongs to user
		return &(*bots)[0]
	}
	// check if the given bot belongs to the user
	if id, err := strconv.Atoi(idString); err == nil {
		for _, b := range *bots {
			if b.ID == id {
				return &b
			}
		}
	}
	return nil
}
