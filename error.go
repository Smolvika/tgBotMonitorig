package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

func validAndPrepare(costStr string) (float64, bool) {
	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		return 0, false
	} else {
		return cost, true
	}
}

// места возможного возникновения ошибок в коде с парсингом сайта
var (
	placeCallbackQuery           = "CallbackQuery message"
	placeMessageCommand          = "MessageCommand message"
	placeMessageNotCommand       = "MessageNotCommand message"
	placeSendMessageAboutBitcoin = "SendMessageAboutBitcoin message"
	placeSendMessageAboutError   = "SendMessageAboutError message"
)

func errorsMessage(place string, err error, msgConf tgbotapi.MessageConfig, db *sql.DB) {
	fmt.Println("error", place, ":", err)
	fmt.Println("message Config:", msgConf)
	if err.Error() == "Forbidden: bot was blocked by the user" {
		err = deleteUserChatIdChangeCostDB(int(msgConf.ChatID), db)
		errorsWorkDB(InfoBitcoinDB, deleteInfo, err)
		err = deleteUserChatIdCostDB(int(msgConf.ChatID), db)
		errorsWorkDB(ChatIdCostDB, deleteInfo, err)
	}
}

// места возможного возникновения ошибок в коде с базой данных
var (
	InfoBitcoinDB      = "problem in work InfoBitcoinDB "
	ChatIdCostDB       = "problem in work ChatIdCostDB "
	ChatIdChangeCostDB = "problem in work ChatIdChangeCostDB "
)

// операции с базой данных
var (
	changeInfo = "changeInfo function not work"
	giveInfo   = "giveInfo function not work"
	addInfo    = "addInfo function not work"
	deleteInfo = "deleteInfo function not work"
)

func errorsWorkDB(place string, operation string, err error) {
	if err != nil {
		fmt.Println(place, operation, ": ", err)
	}
}

// Сообщение о проблемах с работой сайта
func sendMessageUserAboutError(db *sql.DB, bot *tgbotapi.BotAPI, errInfoBitcoinPars *error) {
	if *errInfoBitcoinPars != nil {
		usersChatId, err := allUserChatIdCostDB(db)
		if err != nil {
			fmt.Println("ошибка получения всех уникальных юзеров из CostDB:", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Проблемы с доступом к сайту биржи
оповещения о достижении фиксированной цены некоторое время будут недоступны, 
вы также можете настраивать уведомления, они будут работать корректно после 
решения проблемы с доступом`)
			_, err = bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, db)
			}
		}
		usersChatId, err = allUserChatIdChangeCostDB(db)
		if err != nil {
			fmt.Println("ошибка получения всех уникальных юзеров из ChangeCostDB:", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Проблемы с доступом к сайту биржи
оповещения о фиксированном скачке цены некоторое время будут недоступны, 
вы также можете настраивать уведомления, они будут работать корректно после 
решения проблемы с доступом`)
			_, err = bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, db)
			}
		}
		usersChatId, err = allUserChatIdCostDB(db)
		if err != nil {
			fmt.Println("ошибка получения всех уникальных юзеров из ChangeCostDB:", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Проблемы с доступом к сайту биржи
ежечасные оповещения о цене некоторое время будут недоступны`)
			_, err = bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, db)
			}
		}
		//ожидание исправления ошибки
		for *errInfoBitcoinPars != nil {
			time.Sleep(1 * time.Minute)
		}
		//ошибка исправлена отправляем оповещение
		usersChatId, err = allUserChatIdCostDB(db)
		if err != nil {
			fmt.Println("ошибка получения всех уникальных юзеров из CostDB:", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Доступ к бирже возобновлен
оповещения о достижении фиксированной цены снова доступны`)
			_, err = bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, db)
			}
		}
		usersChatId, err = allUserChatIdChangeCostDB(db)
		if err != nil {
			fmt.Println("ошибка получения всех уникальных юзеров из ChangeCostDB:", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Доступ к бирже возобновлен
оповещения о скачке цены снова доступны`)
			_, err = bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, db)
			}
		}
	}
}
