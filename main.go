package main

import (
	"MyTelegramAssistentAI/src/layers/tgbot"
	"os"
)

// "MyTelegramAssistentAI/src/models"
// "MyTelegramAssistentAI/src/services/dataRepository"
// "log"
// messagehandler "MyTelegramAssistentAI/src/messageHandler"
// "fmt"

func main() {

	e := os.Remove("dbfile.db")
	if e != nil {
		panic(e)
	}
	tgbot.Start()
}
