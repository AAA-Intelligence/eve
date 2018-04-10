package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// CreateUser adds a new user to the database
// It creates a new entry in the table "User"
// The password is hashed with bcrypt
func CreateUser(userName, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("cannot hash password: " + err.Error())
	}
	_, err = dbConnection.db.Exec("INSERT INTO User(Name,PasswordHash) VALUES($1,$2)", userName, hash)
	return err
}

// CheckCredentials verifies if the combination of userName and password is valid.
// The function checks if a user with the given name exists and compares the password with the one in the database.
// The complete user data is returned if the credentials are valid.
func CheckCredentials(userName, password string) (*User, error) {
	var hash []byte
	user := User{
		Name: userName,
	}
	err := dbConnection.db.QueryRow("SELECT UserID,PasswordHash FROM User WHERE Name=$1", user.Name).Scan(&user.ID, &hash)
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

// StoreSessionKey saves the session key in the database in the "User" table
// This session key authenticates the user in further requests
// The function returns true if the storing was successfull
func StoreSessionKey(user *User, key string) bool {
	result, err := dbConnection.db.Exec("UPDATE User SET SessionKey=$1 WHERE UserId = $2", key, user.ID)
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

// GetUserForSession checks if the given sesssion key is associated with any user in the database.
// If the key exists in the database the User, which the session belongs to, is returned.
// An invalid key resolves in the return of a nil pointer
func GetUserForSession(sessionKey string) *User {
	var user User
	err := dbConnection.db.QueryRow("SELECT UserID,Name FROM User WHERE SessionKey=$1", sessionKey).Scan(&user.ID, &user.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("cannot get user for sessionID:", err)
		}
		return nil
	}
	return &user
}

// Gender defines the gender of a bot to be male or female
type Gender int

const (
	// Male is the constant that defines that a bot is male (value = 0)
	Male = Gender(0)
	// Female is the constant that defines that a bot is female (value = 1)
	Female = Gender(1)
)

// Bot represents a bots entry in the database
// A bot is only accessible for one user
// The entry holds personal information about the bot but also information about the current mood and affection to the user.
type Bot struct {
	// The bots unique id
	ID int

	// The bots name, which the user can see
	Name string

	// The path to the bots profile picture
	Image string

	Gender Gender

	// The user the bot belongs to. Only this user can communicate with the bot
	User int

	// The bots current affection to the user
	Affection float64

	// The bots current mood
	Mood float64

	Birthdate     time.Time
	FavoriteColor int
	FatherName    int
	FatherAge     int
	MotherName    int
	MotherAge     int
}

// CreateBot creates a bot entry in the database
// The following fields in the bot struct need to be filled: Name, Image, Gender, User, Affection and Mood
// If the insertion was successful the generated bot id is saved in the given bot struct.
func CreateBot(bot *Bot) error {
	// random values
	bot.FavoriteColor = randomColor()
	bot.Birthdate = randomBirthDate(20, 30)
	bot.FatherName = randomName(Male)
	bot.FatherAge = rand.Intn(20) + 40
	bot.MotherName = randomName(Female)
	bot.MotherAge = bot.FatherAge + rand.Intn(10) - 5

	v, err := dbConnection.db.Exec(`
		INSERT INTO Bot(Name,Image,Gender,User,Affection,Mood,Birthdate,FavoriteColor,FatherName,FatherAge,MotherName,MotherAge) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		bot.Name, bot.Image, bot.Gender, bot.User, bot.Affection, bot.Mood, bot.Birthdate, bot.FavoriteColor, bot.FatherName, bot.FatherAge, bot.MotherName, bot.MotherAge)
	if err != nil {
		log.Println("error inserting new bot:", err)
		return ErrInternalServerError
	}
	botID, err := v.LastInsertId()
	if err != nil {
		log.Println("error receiving inserted rows:", err)
		return ErrInternalServerError
	}
	bot.ID = int(botID)
	return nil
}

// MessageSender defines who send a message
type MessageSender int

const (
	// BotIsSender says that the message was sent by the bot
	BotIsSender = 0

	// UserIsSender says that the message was sent by the user
	UserIsSender = 1
)

// Message represents database entry of a message
// A messsage is a text that was sent between the user and a bot
// Since a bot can only communicate with one user a message is always associated with the bot
type Message struct {
	// ID is a unique id
	ID int

	// The bot who sent the message or received it
	Bot int

	// The sender of the message
	Sender MessageSender

	// The point in time that the message was sent
	Timestamp time.Time

	// The text that was sent
	Content string

	// Affection value of the message
	Affection float64

	// Mood value of the message
	Mood float64
}

// MessageMaxLength defines the maximum message length that a user can send to a bot
const MessageMaxLength = 200

// StoreMessages stores messages in database
// The messages need to be sent between the given bot and user
func StoreMessages(userID, botID int, msgs []Message) error {
	// check if bot belongs to user
	exists, err := rowExists("SELECT * FROM Bot WHERE BotID=$1 AND User=$2", botID, userID)
	if err != nil {
		log.Println("cannot check for bot:", err)
		return ErrInternalServerError
	} else if !exists {
		return ErrBotDoesNotBelongToUser
	}

	// start a new database transaction
	tx, err := dbConnection.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO Message(Bot,Sender,Timestamp,Content,Affection,Mood)
					 VALUES($1,$2,$3,$4,$5,$6)`)
	if err != nil {
		tx.Rollback()
		log.Println("cannot rollback:", err)
		return ErrInternalServerError
	}
	defer stmt.Close()
	for _, m := range msgs {
		if len(m.Content) > MessageMaxLength {
			tx.Rollback()
			return ErrMessageToLong
		}
		_, err := stmt.Exec(botID, m.Sender, m.Timestamp, m.Content, m.Affection, m.Mood)
		if err != nil {
			tx.Rollback()
			log.Println("cannot store message:", err)
			return ErrInternalServerError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("cannot commit changes:", err)
		return ErrInternalServerError
	}
	return nil
}

// GetMessagesForBot returns a list of all messages, that the user and bot sent each other
func GetMessagesForBot(bot int) (*[]Message, error) {
	rows, err := dbConnection.db.Query(`
		SELECT 	Timestamp,
				Content,
				Sender,
				Affection,
				Mood 
		FROM Message 
		WHERE Bot=$1`, bot)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := []Message{}
	var cursor Message
	for rows.Next() {
		if err := rows.Scan(&cursor.Timestamp, &cursor.Content, &cursor.Sender, &cursor.Affection, &cursor.Mood); err == nil {
			messages = append(messages, cursor)
		} else {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &messages, nil
}

// GetBotsForUser returns all bots which belong to the given user
func GetBotsForUser(userID int) (*[]Bot, error) {
	rows, err := dbConnection.db.Query(`
		SELECT	b.BotID,
				b.Name,
				b.Image,
				b.Gender,
				b.Affection,
				b.Mood,
				b.Birthdate,
				b.FavoriteColor,
				b.FatherName,
				b.FatherAge,
				b.MotherName,
				b.MotherAge
		FROM	Bot b
		WHERE b.User=$2`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bots []Bot
	var cursor Bot
	for rows.Next() {
		if err := rows.Scan(&cursor.ID, &cursor.Name, &cursor.Image, &cursor.Gender, &cursor.Affection, &cursor.Mood, &cursor.Birthdate,
			&cursor.FavoriteColor, &cursor.FatherName, &cursor.FatherAge, &cursor.MotherName, &cursor.MotherAge); err == nil {
			bots = append(bots, cursor)
		} else {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &bots, nil
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

// GetBot returns bot entry from database if bot belongs to user
// This funtion can be used to check if a bot belongs to the given user.
func GetBot(botID, userID int) (*Bot, error) {
	bot := Bot{
		ID:   botID,
		User: userID,
	}
	err := dbConnection.db.QueryRow(`
		SELECT	b.Name,
				b.Image,
				b.Gender,
				b.Affection,
				b.Mood,
				b.Birthdate,
				b.FavoriteColor,
				b.FatherName,
				b.FatherAge,
				b.MotherName,
				b.MotherAge
		FROM	Bot b
		WHERE	b.BotID = $1 AND b.User = $2`, bot.ID, bot.User).Scan(
		&bot.Name, &bot.Image, &bot.Gender, &bot.Affection, &bot.Mood, &bot.Birthdate,
		&bot.FavoriteColor, &bot.FatherName, &bot.FatherAge, &bot.MotherName, &bot.MotherAge)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// Name represents database entry
type Name struct {
	ID     int
	Text   string
	Gender int
}

// GEtNames returns all bots which belong to the given user
func GetNames(gender int) (*[]Name, error) {
	rows, err := dbConnection.db.Query(`
		SELECT 	*
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

	name := Name{}
	err := dbConnection.db.QueryRow(`
		SELECT 	*
		FROM Name 
		WHERE NameID = $1`, id).Scan(&name.ID, &name.Text, &name.Gender)

	if err != nil {
		return nil, err
	}
	return &name, nil
}

// Image represents database entry
type Image struct {
	ImageID int
	Gender  int
	Path    string
}

// GetImage returns image object with given id
func GetImages(gender int) (*[]Image, error) {

	rows, err := dbConnection.db.Query(`
		SELECT 	*
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

	image := Image{}
	err := dbConnection.db.QueryRow(`
		SELECT 	*
		FROM Image 
		WHERE ImageID = $1`, id).Scan(&image.ImageID, &image.Gender, &image.Path)

	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetMotherName returns the name of the mother as string
func (bot *Bot) GetMotherName() string {
	var name string
	err := dbConnection.db.QueryRow(`
		SELECT	n.Text 
		FROM Name n 
		INNER JOIN Bot b ON (n.NameID = b.MotherName)
		WHERE b.BotID = $1`, bot.ID).Scan(&name)
	if err != nil {
		log.Println("error getting mother name:", err)
		return "Eva"
	}
	return name
}

// GetFatherName returns the name of the father as string
func (bot *Bot) GetFatherName() string {
	var name string
	err := dbConnection.db.QueryRow(`
		SELECT	n.Text 
		FROM Name n 
		INNER JOIN Bot b ON (n.NameID = b.FatherName)
		WHERE b.BotID = $1`, bot.ID).Scan(&name)
	if err != nil {
		log.Println("error getting father name:", err)
		return "Adam"
	}
	return name
}

// GetFavoriteColor returns the favorite color as string
func (bot *Bot) GetFavoriteColor() string {
	var name string
	err := dbConnection.db.QueryRow(`
		SELECT	c.Name 
		FROM Color c 
		INNER JOIN Bot b ON (c.ColorID = b.FavoriteColor)
		WHERE b.BotID = $1`, bot.ID).Scan(&name)
	if err != nil {
		log.Println("error getting favorite color:", err)
		return "weiß"
	}
	return name
}

// UpdateContext updates the bots fields affection and mood in the database and saves the value in the struct
func (bot *Bot) UpdateContext(affection, mood float64) error {
	bot.Mood, bot.Affection = mood, affection
	result, err := dbConnection.db.Exec(`
		UPDATE Bot
		SET	Affection = $1,
			Mood = $2
		WHERE BotID = $3`, bot.Affection, bot.Mood, bot.ID)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err == nil {
		if rows < 1 {
			return errors.New("not bot with id " + strconv.Itoa(bot.ID) + " found")
		}
	} else {
		return err
	}
	return nil
}
