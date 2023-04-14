package main

import (
	"github.com/Smolvika/tgBotMonitorig/pkg/repository"
	"github.com/Smolvika/tgBotMonitorig/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	botToken := "6151390825:AAHQtZr6ublDAhtJ7Q4yYuHQ24G48rvwkSM"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = true
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "ghjghj",
	})
	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}
	telegramBot := telegram.NewBot(bot, db)
	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Problem with parsing currency information:%v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("Failed to close database:%v\n", err)
	}

}
