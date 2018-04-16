package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	// driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

// The connection, which is used for all database requests
// It is not accessable from outside the package to make sure all interaction with the database is made over the defined functions
var dbConnection struct {
	// Path is the path of the sqlite file
	Path string

	// connection to database driver
	db *sql.DB
}

// User represents a user entry in the database
// Every person that uses eve needs to have a user entry in the database, because it is used for authentication.
type User struct {

	// UserID in database
	ID int

	// Username
	Name string
}

// Connect creates a connection to a sqlite3 database
// The given path is the location of the datbase file
// If the function runs without errors, the database is ready for requests
func Connect(path string) error {
	dbConnection.Path = path
	_, err := os.Stat(path)
	dbExists := !os.IsNotExist(err)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	if !dbExists {
		if data, err := ioutil.ReadFile("db-script.sql"); err == nil {
			if _, err := db.Exec(string(data)); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	dbConnection.db = db
	return nil
}

// Close closes the connection to the database
func Close() error {
	// check if database is conntected
	if dbConnection.db == nil {
		return ErrConnectionClosed
	}
	// remove connection object to avoid requests to closed connection
	defer func() { dbConnection.db = nil }()
	return dbConnection.db.Close()
}

// checks if the given query returns at least one row
func rowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := dbConnection.db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

// Name represents database entry
type Name struct {
	ID     int
	Text   string
	Gender int
}

// GetNames returns all bots which belong to the given user
func GetNames(gender int) (*[]Name, error) {
	rows, err := dbConnection.db.Query(`
		SELECT 	NameID,
				Text,
				Gender
		FROM Name 
		WHERE Gender = $1`, gender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []Name
	var cursor Name
	for rows.Next() {
		if err := rows.Scan(&cursor.ID, &cursor.Text, &cursor.Gender); err == nil {
			names = append(names, cursor)
		} else {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &names, nil
}

// GetName returns all bots which belong to the given user
func GetName(id int) (*Name, error) {
	name := Name{
		ID: id,
	}
	err := dbConnection.db.QueryRow(`
		SELECT 	Text,
				Gender
		FROM Name 
		WHERE NameID = $1`, name.ID).Scan(&name.Text, &name.Gender)

	if err != nil {
		return nil, err
	}
	return &name, nil
}

// Image represents database entry
type Image struct {
	ImageID int
	Gender  Gender
	Path    string
}

// GetImages returns image object with given id
func GetImages(gender int) (*[]Image, error) {
	rows, err := dbConnection.db.Query(`
		SELECT 	ImageID,
				Gender,
				Path
		FROM Image 
		WHERE Gender =$1`, gender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []Image
	var cursor Image
	for rows.Next() {
		if err := rows.Scan(&cursor.ImageID, &cursor.Gender, &cursor.Path); err == nil {
			images = append(images, cursor)
		} else {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &images, nil
}

// GetImage returns image object with given id
func GetImage(id int) (*Image, error) {
	image := Image{
		ImageID: id,
	}
	err := dbConnection.db.QueryRow(`
		SELECT 	Gender,
				Path
		FROM Image 
		WHERE ImageID = $1`, image.ImageID).Scan(&image.Gender, &image.Path)

	if err != nil {
		return nil, err
	}
	return &image, nil
}
