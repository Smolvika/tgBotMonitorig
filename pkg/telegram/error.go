package telegram

import (
	"github.com/Smolvika/tgBotMonitorig/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"log"
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
	placeCallbackQuery            = "CallbackQuery message"
	placeMessageCommand           = "MessageCommand message"
	placeMessageNotCommand        = "MessageNotCommand message"
	placeSendMessageAboutCurrency = "SendMessageAboutCurrency message"
	placeSendMessageAboutError    = "SendMessageAboutError message"
	InfoCurrencyDB                = "problem in work InfoCurrencyDB"
	ChatIdCostDB                  = "problem in work ChatIdCostDB"
	ChatIdChangeCostDB            = "problem in work ChatIdChangeCostDB"
	giveInfo                      = "giveInfo function not work"
	addInfo                       = "addInfo function not work"
	deleteInfo                    = "deleteInfo function not work"
)

func errorsMessage(place string, err error, msgConf tgbotapi.MessageConfig, db *sqlx.DB) {
	log.Printf("error %s: %v\n", place, err)
	log.Printf("message Config: %v\n", msgConf)
	if err.Error() == "Forbidden: bot was blocked by the user" {
		err = repository.DeleteUserChatIdChangeCostDB(int(msgConf.ChatID), db)
		errorsWorkDB(InfoCurrencyDB, deleteInfo, err)
		err = repository.DeleteUserChatIdCostDB(int(msgConf.ChatID), db)
		errorsWorkDB(ChatIdCostDB, deleteInfo, err)
	}
}

func errorsWorkDB(place string, operation string, err error) {
	if err != nil {
		log.Printf("%s %s %v\n", place, operation, err)
	}
}

// Сообщение о проблемах с работой сайта
func (b *Bot) sendMessageUserAboutError(errInfoBitcoinPars *error) {
	if *errInfoBitcoinPars != nil {
		usersChatId, err := repository.AllUserChatIdCostDB(b.db)
		if err != nil {
			log.Printf("error getting all possible users from CostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Проблемы с доступом к сайту биржи
ежечасные оповещения о цене валюты некоторое время будут недоступны, 
вы также можете настраивать уведомления, они будут работать корректно после 
решения проблемы с доступом`)
			_, err = b.bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, b.db)
			}
		}
		usersChatId, err = repository.AllUserChatIdChangeCostDB(b.db)
		if err != nil {
			log.Printf("error getting all possible users from ChangeCostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Проблемы с доступом к сайту биржи
оповещения о дохождении до фиксированной  цены некоторое время будут недоступны, 
вы также можете настраивать уведомления, они будут работать корректно после 
решения проблемы с доступом`)
			_, err = b.bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, b.db)
			}
		}
		//ожидание исправления ошибки
		for *errInfoBitcoinPars != nil {
			time.Sleep(1 * time.Minute)
		}
		//ошибка исправлена отправляем оповещение
		usersChatId, err = repository.AllUserChatIdCostDB(b.db)
		if err != nil {
			log.Printf("error getting all possible users from CostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Доступ к бирже возобновлен ежечасные оповещения о цене снова доступны`)
			_, err = b.bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, b.db)
			}
		}
		usersChatId, err = repository.AllUserChatIdChangeCostDB(b.db)
		if err != nil {
			log.Printf("error getting all possible users from ChangeCostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Доступ к бирже возобновлен оповещения о достижении фиксированной цены снова доступны`)
			_, err = b.bot.Send(msg)
			if err != nil {
				errorsMessage(placeSendMessageAboutError, err, msg, b.db)
			}
		}
	}
}
