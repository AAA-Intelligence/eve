package main

import (
	"strings"

	"github.com/AAA-Intelligence/eve/db"
)

func handleMessage(message string, user *db.User) (answer string) {
	if strings.HasPrefix(message, "nenn mich") {
		user.Name = strings.TrimPrefix(message, "nenn mich")
	}
	answer = "ok, " + user.Name
	return
}
