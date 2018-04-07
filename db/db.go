package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	// driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

// default databaseConnection to use
var dbConection struct {
	// Path is the path of the sqlite file
	Path string

	// connection to database driver
	db *sql.DB
}

// User holds user data
type User struct {

	// UserID in database
	ID int

	// Username
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
	// check if database is conntected
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
	return err
}

// CheckCredentials checks if combination of userName and password is valid
// User object data is returned if credentials are valid
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

// Gender defines the gender of a bot to be male or female
type Gender int

const (
	// Male is the constant that defines that a bot is male (value = 0)
	Male = Gender(0)
	// Female is the constant that defines that a bot is female (value = 1)
	Female = Gender(1)
)

// Bot represents database entry of a bot
type Bot struct {
	ID        int
	Name      string
	Image     string
	Gender    Gender
	User      int
	Affection float64
	Mood      float64
}

// CreateBot creates a bot entry in the database and fills the empty values in the given bot struct
func CreateBot(bot *Bot) error {
	v, err := dbConection.db.Exec("INSERT INTO Bot(Name,Image,Gender,User,Affection,Mood) VALUES($1,$2,$3,$4,$5,$6)", bot.Name, bot.Image, bot.Gender, bot.User, bot.Affection, bot.Mood)
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
type Message struct {
	ID        int
	Bot       int
	Sender    MessageSender
	Timestamp time.Time
	Content   string
	Rating    float64
}

// MessageMaxLength defines the maximum message length
const MessageMaxLength = 200

// StoreMessages saves message in database
func StoreMessages(userID, botID int, msgs []Message) error {
	// check if bot belongs to user
	exists, err := rowExists("SELECT * FROM Bot WHERE BotID=$1 AND User=$2", botID, userID)
	if err != nil {
		log.Println("cannot check for bot:", err)
		return ErrInternalServerError
	} else if !exists {
		return ErrBotDoesNotBelongToUser
	}

	tx, err := dbConection.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO Message(Bot,Sender,Timestamp,Content,Rating)
					 VALUES($1,$2,$3,$4,$5)`)
	if err != nil {
		tx.Rollback()
		return ErrInternalServerError
	}
	defer stmt.Close()
	for _, m := range msgs {
		if len(m.Content) > MessageMaxLength {
			tx.Rollback()
			return ErrMessageToLong
		}
		_, err := stmt.Exec(botID, m.Sender, m.Timestamp, m.Content, m.Rating)
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
	rows, err := dbConection.db.Query(`
		SELECT 	Timestamp,
				Content,
				Sender,
				Rating 
		FROM Message 
		WHERE Bot=$1`, bot)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := []Message{}
	var cursor Message
	for rows.Next() {
		if err := rows.Scan(&cursor.Timestamp, &cursor.Content, &cursor.Sender, &cursor.Rating); err == nil {
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
	rows, err := dbConection.db.Query(`
		SELECT 	BotID,
				Name,
				Image,
				Gender,
				Affection, 
				Mood
		FROM Bot 
		WHERE User=$2`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bots []Bot
	var cursor Bot
	for rows.Next() {
		if err := rows.Scan(&cursor.ID, &cursor.Name, &cursor.Image, &cursor.Gender, &cursor.Affection, &cursor.Mood); err == nil {
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

//GetMessagesForUser returns alle messages the user has send with any bot
func GetMessagesForUser(user *User) (*map[int][]Message, error) {
	rows, err := dbConection.db.Query(`
		SELECT	m.MessageID,
				m.Sender,
				m.Timestamp,
				m.Content,
				m.Rating,
				b.BotID
		FROM  Message m
		INNER JOIN Bot b ON (b.BotID = m.Bot)
		WHERE b.User = $1
		ORDER BY b.BotID,m.Timestamp`, user.ID)
	if err != nil {
		log.Println("cannot get messages for user", user.ID, ":", err)
		return nil, ErrInternalServerError
	}
	defer rows.Close()
	var cursor Message
	messages := make(map[int][]Message)
	for rows.Next() {
		if err = rows.Scan(&cursor.ID, &cursor.Sender, &cursor.Timestamp, &cursor.Content, &cursor.Rating, &cursor.Bot); err == nil {
			messages[cursor.Bot] = append(messages[cursor.Bot], cursor)
		} else {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, ErrInternalServerError
	}
	return &messages, nil
}

// checks if the given query returns at least one row
func rowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := dbConection.db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

// GetBot returns bot entry from database if bot belongs to user
func GetBot(botID, userID int) (*Bot, error) {
	bot := Bot{
		ID:   botID,
		User: userID,
	}
	err := dbConection.db.QueryRow(`
		SELECT	b.Name,
				b.Image,
				b.Gender,
				b.Affection,
				b.Mood
		FROM	Bot b
		WHERE	b.BotID = $1 AND b.User = $2`, bot.ID, bot.User).Scan(
		&bot.Name, &bot.Image, &bot.Gender, &bot.Affection, &bot.Mood)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// Name represents database entry
type Name struct {
	ID        	int
	Name      	string
	Sex     	int
}

// GEtNames returns all bots which belong to the given user
func GetNames(sex int) (*[]Name, error) {
	rows, err := dbConection.db.Query(`
		SELECT 	*
		FROM Name 
		WHERE Sex = $1`,sex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []Name
	var cursor Name
	for rows.Next() {
		if err := rows.Scan(&cursor.ID, &cursor.Name, &cursor.Sex); err == nil {
			bots = append(bots, cursor)
		} else {
			log.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &names, nil
}

