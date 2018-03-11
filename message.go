package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func handleMessage(message string) (answer string) {
	reply := "messages/bot_msg.json"
	//fmt.Println(reply)
	f, err := ioutil.ReadFile(reply)
	if err != nil {
		log.Fatal(err)
	}
	byt := []byte(string(f))

	//string to json
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		log.Fatal(err)
	}
	answer = dat["content"].(string)

	//fmt.Println(text)
	//	answer = "ok"
	return
}
