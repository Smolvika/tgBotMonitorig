package main

import (
	"context"
	"github.com/Smolvika/tgBotMonitorig/config"
	"github.com/Smolvika/tgBotMonitorig/pkg/repository"
	"github.com/Smolvika/tgBotMonitorig/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	cfg := config.New()
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = true
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.ConfigDB.Host,
		Port:     cfg.ConfigDB.Port,
		Username: cfg.ConfigDB.Username,
		DBName:   cfg.ConfigDB.DBName,
		SSLMode:  cfg.ConfigDB.SSLMode,
		Password: cfg.ConfigDB.Password,
	})
	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	repos := repository.New(db)
	telegramBot := telegram.NewBot(bot, repos)
	go func() {
		if err := telegramBot.Start(ctx); err != nil {
			log.Fatalf("Problem with parsing currency information:%v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
	log.Println(" tgBot Shutting Down")

	if err := db.Close(); err != nil {
		log.Printf("Failed to close database:%v\n", err)
	}
}
