package db

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	// driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

// default databaseConnection to use
var dbConection struct {
	Path string
	db   *sql.DB
}

// User holds user data
type User struct {
	ID   int
	Name string
}

// Connect creates a connection to database
func Connect(path string) error {
	dbConection.Path = path
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	dbConection.db = db
	return nil
}

// Close closes a connection to database
func Close() error {
	if dbConection.db == nil {
		return ErrConnectionClosed
	}
	defer func() { dbConection.db = nil }()
	return dbConection.db.Close()
}

// CreateUser adds user to database
func CreateUser(userName, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("cannot hash password: " + err.Error())
	}
	_, err = dbConection.db.Exec("INSERT INTO User(Name,PasswordHash) VALUES($1,$2)", userName, hash)
	if err != nil {
		return err
	}
	return nil
}

//CheckCredentials checks if combination of userName and password is valid
func CheckCredentials(userName, password string) (*User, error) {
	var hash []byte
	user := User{
		Name: userName,
	}
	var sessionKey sql.NullString
	err := dbConection.db.QueryRow("SELECT UserID,PasswordHash,SessionKey FROM User WHERE Name=$1", user.Name).Scan(&user.ID, &hash, &sessionKey)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotExists
	} else if err != nil {
		return nil, ErrInternalServerError
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrWrongPassword
	} else if err != nil {
		return nil, ErrInternalServerError
	}
	return &user, nil
}

//StoreSessionKey saves the key in the database for the given user and saves the key to the user struct
func StoreSessionKey(user *User, key string) bool {
	result, err := dbConection.db.Exec("UPDATE User SET SessionKey=$1 WHERE UserId = $2", key, user.ID)
	if err != nil {
		log.Println("cannot update session key for user ", user.ID, err)
		return false
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("session key was not stored:", err)
		return false
	} else if rows != 1 {
		log.Println("cannot update session key for user: ", user.ID)
		return false
	}
	return true
}

//GetUserForSession gets the user associated with the given session key
func GetUserForSession(sessionKey string) *User {
	var user User
	err := dbConection.db.QueryRow("SELECT UserID,Name FROM User WHERE SessionKey=$1", sessionKey).Scan(&user.ID, &user.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("cannot get user for sessionID:", err)
		}
		return nil
	}
	return &user
}
