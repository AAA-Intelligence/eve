package db

import (
	"log"
	"math/rand"
	"time"
)

// Year is the duraction of one year in nano seconds
const Year time.Duration = time.Hour * 24 * 365


func randomBirthDate(minAge, maxAge int) time.Time {
	nowYear := time.Now().Year()
	min := time.Date(nowYear-maxAge, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(nowYear-minAge, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// returns primary key of random color
func randomColor() int {
	var color int
	err := dbConnection.db.QueryRow("SELECT ColorID FROM Color ORDER BY RANDOM() LIMIT 1").Scan(&color)
	if err != nil {
		log.Println("error generating random color:", err)
		return 0
	}
	return color
}

// returns primary key of random name
func randomName(gender Gender) int {
	var name int
	err := dbConnection.db.QueryRow("SELECT NameId FROM Name WHERE Gender = $1 ORDER BY RANDOM() LIMIT 1", gender).Scan(&name)
	if err != nil {
		log.Println("error generating random name:", err)
		return 0
	}
	return name
}
