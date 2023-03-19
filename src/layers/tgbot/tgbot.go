package tgbot

import (
	"MyTelegramAssistentAI/src/config"
	messagelayer "MyTelegramAssistentAI/src/layers/messageLayer"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start() {
	bot, err := tgbotapi.NewBotAPI(config.GetValue("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	messagelayer.Start()
	listener(bot)
}

func listener(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			response, attachment := messagelayer.SendReuest(update.Message.Chat.ID, &update.Message.Text)

			var msg tgbotapi.MessageConfig

			if *attachment != "" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, *response+"\n<a href=\""+*attachment+"\">Img link</a>")
				msg.ParseMode = "HTML"
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, *response)
			}

			bot.Send(msg)
		}
	}
}
