package bots

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/AAA-Intelligence/eve/db"
)

// A instance of the bot python script
// This struct can be used to communicate with the pyhton script via stdin and stdout
type botInstance struct {
	// The command that started the script
	cmd *exec.Cmd

	// everything that is written to this writer is sent to the script
	writer io.WriteCloser

	// read the scripts output
	reader *bufio.Reader

	// last time the bot instance was used / handled a requests
	lastUsed time.Time
}

// MessageData is used so send all needed information to the bot instance
type MessageData struct {
	Text            string    `json:"text"`
	PreviousPattern *int      `json:"previous_pattern,omitempty"`
	Mood            float64   `json:"mood"`
	Affection       float64   `json:"affection"`
	Gender          db.Gender `json:"bot_gender"`
	Name            string    `json:"bot_name"`
	Birthdate       int64     `json:"bot_birthdate"` // Unix timestamp
	FavoriteColor   string    `json:"bot_favorite_color"`
	FatherName      string    `json:"father_name"`
	FatherAge       int       `json:"father_age"`
	MotherName      string    `json:"mother_name"`
	MotherAge       int       `json:"mother_age"`
}

// BotAnswer is the answer returned by the bot instance
type BotAnswer struct {
	Text      string  `json:"text"`
	Pattern   *int    `json:"pattern,omitempty"`
	Mood      float64 `json:"mood"`
	Affection float64 `json:"affection"`
}

// sends request to bot instance, waits for the response and returns it.
// the function always returns an answer
func (b *botInstance) sendRequest(data MessageData) *BotAnswer {
	b.lastUsed = time.Now()
	writer := b.writer
	serialized, err := json.Marshal(data)
	_, err = fmt.Fprintln(writer, string(serialized))
	if err != nil {
		log.Println("error writing to bot pipt:", err)
		return errorBotAnswer(data.Mood, data.Affection)
	}
	response, _, err := b.reader.ReadLine()
	if err != nil {
		log.Println("error reading from pipe:", err)
		return errorBotAnswer(data.Mood, data.Affection)
	}
	msg := &BotAnswer{}
	err = json.Unmarshal(response, msg)
	if err != nil {
		if string(response) == "error" {
			// bot instance returns "error" if the request could not be processed
			log.Println("an error in the bot instance occurred")
		} else {
			log.Println("error reading response:", err)
		}
		return errorBotAnswer(data.Mood, data.Affection)
	}
	return msg
}

// creates new instance of python script that handles message requests
// if no error occures the bot instance struct is returned
func newBotInstance(python string) (*botInstance, error) {
	cmd := exec.Command(python, "-m", "bot")

	writer, err := cmd.StdinPipe()
	if err != nil {
		//log.Fatal("error creating stdin pipe:", err)
		return nil, err
	}
	reader, err := cmd.StdoutPipe()
	if err != nil {
		//log.Fatal("error creating stdout pipe:", err)
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}
	return &botInstance{
		cmd:    cmd,
		writer: writer,
		reader: bufio.NewReader(reader),
	}, nil
}

// The answer that is returned when the bot instance returns no valid answer
// The Text is "Ok" and the given mood and affection is used
func errorBotAnswer(mood, affection float64) *BotAnswer {
	return &BotAnswer{
		Text:      "Ok",
		Mood:      mood,
		Affection: affection,
	}
}

// Close closes the pipe to the python script and thereby the process is stoped
func (b *botInstance) Close() {
	b.writer.Close()
}
