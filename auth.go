package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/AAA-Intelligence/eve/db"
	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	if getUser(r.Context()) != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
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
func basicAuth(h http.HandlerFunc, needsAuth bool) http.HandlerFunc {
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
			if needsAuth {
				http.Error(w, "invalid user credentials", http.StatusUnauthorized)
				return
			}
		} else {
			// save user data in context
			r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))
		}

		w.Header().Del("WWW-Authenticate")
		h.ServeHTTP(w, r)
	}
}

func getUser(ctx context.Context) *db.User {
	user, ok := ctx.Value(UserContextKey).(*db.User)
	if !ok {
		return nil
	}
	return user
}
