package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	var errInfoCurrencyNowPars error
	botToken := "6151390825:AAHQtZr6ublDAhtJ7Q4yYuHQ24G48rvwkSM"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	//cfg := Config{
	//	Host:     "localhost",
	//	Port:     "5432",
	//	Username: "postgres",
	//	DBName:   "postgres",
	//	SSLMode:  "disable",
	//	Password: os.Getenv("db_password"),
	//})
	db, err := sql.Open("mysql", "root:password@/userschat")
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Printf("Не удалось закрыть БД:%v ", err)
		}
	}()
	//массив с данными о статусе пользователя
	status := make(map[int64]string, 0)
	//Получение информации с сайта, до запуска бота
	currencyNow, err := parsAllInfoCurrency()
	if err != nil {
		log.Panic("Проблема с парсингом информации о валюте: ", err)
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { //пустое обновление
			continue
		} else if update.CallbackQuery != nil { //обработка callback ответа
			isCallbackQuery(&update, bot, db, status)
		} else if update.Message.IsCommand() { //обработка команд
			isCommandCase(&update, bot, db, currencyNow, &errInfoCurrencyNowPars)
		} else if !update.Message.IsCommand() { //обработка сообщений
			isUsualMessage(&update, bot, db, status)
		}

	}
}
