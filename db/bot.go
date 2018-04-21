package db

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

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

	// Pattern is the recognized pattern id in the last message
	Pattern *int

	Birthdate     time.Time

	// FavoriteColor is the primary key of the bots favorite color
	FavoriteColor int
	
	// FatherName is the primary key of the bots fathers name
	FatherName    int

	// FatherAge is the bots fathers age in years
	FatherAge     int

	// MotherName is the primary key of the bots mothers name
	MotherName    int

	// MotherAge is the bots mothers age in years
	MotherAge     int
	
	// CreationDate is the point in time when the bot was created by a user
	CreationDate  time.Time
}

// Create creates a bot entry in the database
// The following fields in the bot struct need to be filled: Name, Image, Gender, User, Affection and Mood
// If the insertion was successful the generated bot id is saved in the given bot struct.
func (b *Bot) Create() error {
	// random values
	b.FavoriteColor = randomColor()
	b.Birthdate = randomBirthDate(20, 30)
	b.FatherName = randomName(Male)
	b.FatherAge = rand.Intn(20) + 40
	b.MotherName = randomName(Female)
	b.MotherAge = b.FatherAge + rand.Intn(10) - 5
	b.CreationDate = time.Now()

	v, err := dbConnection.db.Exec(`
		INSERT INTO Bot(Name,Image,Gender,User,Affection,Mood,Pattern,Birthdate,FavoriteColor,FatherName,FatherAge,MotherName,MotherAge,CreationDate) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		b.Name, b.Image, b.Gender, b.User, b.Affection, b.Mood, b.Pattern, b.Birthdate, b.FavoriteColor, b.FatherName, b.FatherAge, b.MotherName, b.MotherAge, b.CreationDate)
	if err != nil {
		log.Println("error inserting new bot:", err)
		return ErrInternalServerError
	}
	botID, err := v.LastInsertId()
	if err != nil {
		log.Println("error receiving inserted rows:", err)
		return ErrInternalServerError
	}
	b.ID = int(botID)
	return nil
}

// Delete removes bot from database
func (b *Bot) Delete() error {
	result, err := dbConnection.db.Exec(`
		DELETE FROM Bot
		WHERE BotID = $1`, b.ID)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err == nil {
		if rows < 1 {
			return errors.New("no bot wiht id" + strconv.Itoa(b.ID) + " found")
		}
	} else {
		return err
	}
	return nil
}

// GetMotherName returns the name of the mother as string
func (b *Bot) GetMotherName() string {
	var name string
	err := dbConnection.db.QueryRow(`
		SELECT	n.Text 
		FROM Name n 
		INNER JOIN Bot b ON (n.NameID = b.MotherName)
		WHERE b.BotID = $1`, b.ID).Scan(&name)
	if err != nil {
		log.Println("error getting mother name:", err)
		return "Eva"
	}
	return name
}

// GetFatherName returns the name of the father as string
func (b *Bot) GetFatherName() string {
	var name string
	err := dbConnection.db.QueryRow(`
		SELECT	n.Text 
		FROM Name n 
		INNER JOIN Bot b ON (n.NameID = b.FatherName)
		WHERE b.BotID = $1`, b.ID).Scan(&name)
	if err != nil {
		log.Println("error getting father name:", err)
		return "Adam"
	}
	return name
}

// GetFavoriteColor returns the favorite color as string
func (b *Bot) GetFavoriteColor() string {
	var name string
	err := dbConnection.db.QueryRow(`
		SELECT	c.Name 
		FROM Color c 
		INNER JOIN Bot b ON (c.ColorID = b.FavoriteColor)
		WHERE b.BotID = $1`, b.ID).Scan(&name)
	if err != nil {
		log.Println("error getting favorite color:", err)
		return "weiß"
	}
	return name
}

// UpdateContext updates the bots fields affection and mood in the database and saves the value in the struct
func (b *Bot) UpdateContext(affection, mood float64, pattern *int) error {
	b.Mood, b.Affection = mood, affection
	b.Pattern = pattern
	result, err := dbConnection.db.Exec(`
		UPDATE Bot
		SET	Affection = $1,
			Mood = $2,
			Pattern = $3
		WHERE BotID = $4`, b.Affection, b.Mood, b.Pattern, b.ID)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err == nil {
		if rows < 1 {
			return errors.New("no bot with id " + strconv.Itoa(b.ID) + " found")
		}
	} else {
		return err
	}
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
}

// MessageMaxLength defines the maximum message length that a user can send to a bot
const MessageMaxLength = 200

// StoreMessages stores messages in database
// The messages need to be sent between the given bot and user
func (b *Bot) StoreMessages(user *User, msgs []Message) error {
	// check if bot belongs to user
	exists, err := rowExists("SELECT * FROM Bot WHERE BotID=$1 AND User=$2", b.ID, user.ID)
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
		INSERT INTO Message(Bot,Sender,Timestamp,Content)
					 VALUES($1,$2,$3,$4)`)
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
		_, err := stmt.Exec(b.ID, m.Sender, m.Timestamp, m.Content)
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

// GetMessages returns a list of all messages, that the user and bot sent each other
func (b *Bot) GetMessages() (*[]Message, error) {
	rows, err := dbConnection.db.Query(`
		SELECT 	Timestamp,
				Content,
				Sender
		FROM Message 
		WHERE Bot=$1`, b.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := []Message{}
	var cursor Message
	for rows.Next() {
		if err := rows.Scan(&cursor.Timestamp, &cursor.Content, &cursor.Sender); err == nil {
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
