package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/AAA-Intelligence/eve/db"
	"golang.org/x/crypto/bcrypt"
)

// HTTP handler for creating a user
// request METHOD needs to be post
// request BODY needs to contain username and password
// e.g. username=peter&password=super_secret_password
//
// if the creation was successful the response directs to the index page
func createUser(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		http.Error(w, "HTTP POST only", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		log.Println("error:", err.Error())
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 1 || len(password) < 1 {
		http.Error(w, "missing username or password", http.StatusBadRequest)
		return
	}
	err = db.CreateUser(username, password)
	if err != nil {
		http.Error(w, "cannot register user", http.StatusInternalServerError)
		log.Println("error:", err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type key string

// UserContextKey is the key for the user data that is stored in the request context
const UserContextKey key = "user"

// middleware for handler, that authenticates the user
// basic access authentication is used
// see: https://en.wikipedia.org/wiki/Basic_access_authentication
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		username, password, authOK := r.BasicAuth()
		if !authOK {
			http.Error(w, "invalid user credentials", http.StatusUnauthorized)
			return
		}

		user, err := db.CheckCredentials(username, password)
		if err != nil {
			if err != bcrypt.ErrMismatchedHashAndPassword {
				log.Println(err)
			}
			http.Error(w, "invalid user credentials", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))

		w.Header().Del("WWW-Authenticate")
		h.ServeHTTP(w, r)
	}
}

// gets user data from request context
func getUser(ctx context.Context) *db.User {
	user, ok := ctx.Value(UserContextKey).(*db.User)
	if !ok {
		return nil
	}
	return user
}
