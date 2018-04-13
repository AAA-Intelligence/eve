package db

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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
func (u *User) StoreSessionKey(key string) bool {
	result, err := dbConnection.db.Exec("UPDATE User SET SessionKey=$1 WHERE UserId = $2", key, u.ID)
	if err != nil {
		log.Println("cannot update session key for user ", u.ID, err)
		return false
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("session key was not stored:", err)
		return false
	} else if rows != 1 {
		log.Println("cannot update session key for user: ", u.ID)
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

// GetBots returns all bots which belong to the user
func (u *User) GetBots() (*[]Bot, error) {
	rows, err := dbConnection.db.Query(`
		SELECT	b.BotID,
				b.Name,
				b.Image,
				b.Gender,
				b.Affection,
				b.Mood,
				b.Pattern,
				b.Birthdate,
				b.FavoriteColor,
				b.FatherName,
				b.FatherAge,
				b.MotherName,
				b.MotherAge
		FROM	Bot b
		WHERE b.User=$2`, u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bots []Bot
	var cursor Bot
	for rows.Next() {
		if err := rows.Scan(&cursor.ID, &cursor.Name, &cursor.Image, &cursor.Gender, &cursor.Affection, &cursor.Mood, &cursor.Pattern, &cursor.Birthdate,
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

// GetBot returns bot entry from database if bot belongs to user
// This funtion can be used to check if a bot belongs to the given user.
func (u *User) GetBot(botID int) (*Bot, error) {
	bot := Bot{
		ID:   botID,
		User: u.ID,
	}
	err := dbConnection.db.QueryRow(`
		SELECT	b.Name,
				b.Image,
				b.Gender,
				b.Affection,
				b.Mood,
				b.Pattern,
				b.Birthdate,
				b.FavoriteColor,
				b.FatherName,
				b.FatherAge,
				b.MotherName,
				b.MotherAge
		FROM	Bot b
		WHERE	b.BotID = $1 AND b.User = $2`, bot.ID, bot.User).Scan(
		&bot.Name, &bot.Image, &bot.Gender, &bot.Affection, &bot.Mood, &bot.Pattern, &bot.Birthdate,
		&bot.FavoriteColor, &bot.FatherName, &bot.FatherAge, &bot.MotherName, &bot.MotherAge)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}
