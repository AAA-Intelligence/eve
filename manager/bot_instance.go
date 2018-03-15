package manager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

type BotInstance struct {
	cmd    *exec.Cmd
	writer io.Writer
	reader *bufio.Reader
	mutex  *sync.Mutex
}

type MessageData struct {
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

func (b BotInstance) sendRequest(data MessageData) {
	writer := b.writer
	serialized, err := json.Marshal(data)
	_, err = writer.Write(serialized)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(writer)
	response, _, err := b.reader.ReadLine()
	json.Unmarshal(response, &MessageData{})
}

func newBotInstance() BotInstance {
	cmd := exec.Command("python", "-m", "bot")
	cmd.Start()

	writer, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	return BotInstance{writer: writer, cmd: cmd, mutex: &sync.Mutex{}, reader: bufio.NewReader(reader)}
}
