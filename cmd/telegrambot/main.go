package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Environment variable %s is not set", key)
	}
	return val
}

func handleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to Hilal Scope Bot!")
	bot.Send(msg)
}

func handleHilal(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hilal data is not yet available.")
	bot.Send(msg)
}

func main() {
	token := getEnv("BOT_TOKEN")
	if token == "" {
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			handleStart(bot, update)
		case "hilal":
			handleHilal(bot, update)
		}
	}
}
