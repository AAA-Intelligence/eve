package manager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/drhodes/golorem"

	"github.com/AAA-Intelligence/eve/db"
)

func createBot(w http.ResponseWriter, r *http.Request) {
	err := db.CreateBot(&db.Bot{
		Name:   lorem.Word(3, 10),
		Image:  "hässlich.png",
		Gender: "apache",
		User:   GetUserFromRequest(r).ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type BotInstance struct {
	cmd    *exec.Cmd
	writer io.Writer
	reader *bufio.Reader
	mutex  *sync.Mutex
}

type MessageData struct {
	text    string
	user_id int
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

func createBotToPython() (botInstance BotInstance) {
	cmd := exec.Command("python", "bot/__main__.py")
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

func getMessages(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	messages, err := db.GetMessagesForUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(*messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
