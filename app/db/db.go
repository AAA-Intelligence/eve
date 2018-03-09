package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	// driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

// Errors used by the Database handler
var (
	ErrConnectionClosed = errors.New("connection is closed or not established jet")
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

// DB returns database object
func DB() *sql.DB {
	return dbConection.db
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
	var userID int
	err := dbConection.db.QueryRow("SELECT UserID,PasswordHash FROM User WHERE Name=$1", userName).Scan(&userID, &hash)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return nil, err
	}

	return &User{
		ID:   userID,
		Name: userName,
	}, nil
}
