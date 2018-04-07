package manager

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/AAA-Intelligence/eve/db"
)

// regex rule for user names. all names must match this pattern at registration
// see https://regex101.com/r/WXXRUl/1 for testing and closer explanation
const userNameRule = "^[a-zA-Z0-9]+([^-\\s]?[a-zA-Z0-9])*$"

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
	if ok, err := regexp.Match(userNameRule, []byte(username)); !ok || err != nil {
		http.Error(w, "username must have form of: "+userNameRule, http.StatusBadRequest)
		return
	}
	err = db.CreateUser(username, password)
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
			http.Error(w, "username allready in use", http.StatusInternalServerError)
		} else {
			http.Error(w, "cannot register user", http.StatusInternalServerError)
			log.Println("error:", err.Error())
		}
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// a key is simply a string
type key string

// UserContextKey is the key for the user data that is stored in the request context
const UserContextKey key = "user"

// SessionKey is the key for the cookie in the HTTP header
const SessionKey = "eve-session"

// middleware for handler, that authenticates the user
// basic access authentication is used
// also a session cookie is used, because basic auth does not work for sockets in safari.
// the session key is stored in the User table in the database.
// see: https://en.wikipedia.org/wiki/Basic_access_authentication
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if user is allready authenticated
		cookie, err := r.Cookie(SessionKey)
		if err == nil {
			user := db.GetUserForSession(cookie.Value)
			if user != nil {
				r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))
				next.ServeHTTP(w, r)
				return
			}
			// delete invalid session key
			c := &http.Cookie{
				Name:    SessionKey,
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),
			}
			http.SetCookie(w, c)
		}

		// get credentials via basic auth
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		username, password, authOK := r.BasicAuth()
		if !authOK {
			http.Error(w, db.ErrNoUserCredentials.Error(), http.StatusUnauthorized)
			return
		}
		user, err := db.CheckCredentials(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// set session key for safari
		sessionKey, err := GenerateRandomString(11)
		if err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:    SessionKey,
				Value:   sessionKey,
				Expires: time.Now().Add(365 * 24 * time.Hour),
				Path:    "/",
			})
			db.StoreSessionKey(user, sessionKey)
		} else {
			log.Println("cannot generate session key:", err)
		}
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))

		w.Header().Del("WWW-Authenticate")
		next.ServeHTTP(w, r)
	}
}

// GetUserFromRequest extracts the user stored in the request context
// If no user is stored, nil is returned
func GetUserFromRequest(r *http.Request) *db.User {
	user, ok := r.Context().Value(UserContextKey).(*db.User)
	if !ok {
		return nil
	}
	return user
}
