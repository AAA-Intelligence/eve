package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"text/template"

	"github.com/AAA-Intelligence/eve/db"
)

// Config configures web server
type Config struct {
	Host  string
	HTTP  int
	HTTPS int
}

func loadConfig() *Config {
	var config Config
	flag.StringVar(&config.Host, "host", "", "hostname")
	flag.IntVar(&config.HTTP, "http", 80, "HTTP port")
	flag.IntVar(&config.HTTPS, "https", 443, "HTTPS port")
	flag.Parse()
	return &config
}

//IndexHandler serves index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	tpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading template:", err.Error())
		return
	}
	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, struct {
		User *db.User
	}{
		User: GetUserFromRequest(r),
	})
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error executing template:", err.Error())
		return
	}
	w.Write([]byte(buffer.String()))
}

//RegisterHandler serves index page
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/register.gohtml")
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error loading template:", err)
		return
	}
	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, nil)
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		log.Println("error executing template:", err)
		return
	}
	w.Write([]byte(buffer.String()))
}

func main() {
	config := loadConfig()
	err := db.Connect("eve.sqlite")
	if err != nil {
		log.Panic("error connecting to database: ", err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicAuth(IndexHandler))
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/createUser", createUser)
	mux.HandleFunc("/ws", basicAuth(webSocket))

	// handle static files like css
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := http.Server{
		Addr:    config.Host + ":" + strconv.Itoa(config.HTTP),
		Handler: mux,
	}
	//go startBot()

	log.Println("Starting web server")
	server.RegisterOnShutdown(onShutdown)
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

func onShutdown() {
	log.Println("shutting down...")
	err := db.Close()
	if err != nil {
		log.Panic("error closing connection to database: ", err)
		return
	}
}

func startBot() {
	//generate new message
	log.Println("(almost) bot was succesfully started")
	cmd := exec.Command("python", "bot/request_handler.py")
	cmd.Run()

}
