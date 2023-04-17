package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (b *Bot) ErrorsMessage(place string, err error, msgConf tgbotapi.MessageConfig) {
	log.Printf("error %s: %v\n", place, err)
	log.Printf("message Config: %v\n", msgConf)
	if err.Error() == "Forbidden: bot was blocked by the user" {
		err = b.db.ChangeCost.DeleteUserChatIdChangeCostDB(int(msgConf.ChatID))
		errorsWorkDB(InfoCurrencyDB, deleteInfo, err)
		err = b.db.Cost.DeleteUserChatIdCostDB(int(msgConf.ChatID))
		errorsWorkDB(ChatIdCostDB, deleteInfo, err)
	}
}

func errorsWorkDB(place string, operation string, err error) {
	if err != nil {
		log.Printf("%s %s %v\n", place, operation, err)
	}
}

func (b *Bot) sendMessageUserAboutError(errInfoCurrencyNowPars *error) {
	if *errInfoCurrencyNowPars != nil {
		usersChatId, err := b.db.Cost.AllUserChatIdCostDB()
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
				b.ErrorsMessage(placeSendMessageAboutError, err, msg)
			}
		}
		usersChatId, err = b.db.ChangeCost.AllUserChatIdChangeCostDB()
		if err != nil {
			log.Printf("error getting all possible users from ChangeCostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Проблемы с доступом к сайту биржи
оповещения о дохождении валюты до фиксированной  цены некоторое время будут недоступны, 
вы также можете настраивать уведомления, они будут работать корректно после 
решения проблемы с доступом`)
			_, err = b.bot.Send(msg)
			if err != nil {
				b.ErrorsMessage(placeSendMessageAboutError, err, msg)
			}
		}
		for *errInfoCurrencyNowPars != nil {
			time.Sleep(1 * time.Minute)
		}
		usersChatId, err = b.db.Cost.AllUserChatIdCostDB()
		if err != nil {
			log.Printf("error getting all possible users from CostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Доступ к бирже возобновлен ежечасные оповещения о цене валюты снова доступны`)
			_, err = b.bot.Send(msg)
			if err != nil {
				b.ErrorsMessage(placeSendMessageAboutError, err, msg)
			}
		}
		usersChatId, err = b.db.ChangeCost.AllUserChatIdChangeCostDB()
		if err != nil {
			log.Printf("error getting all possible users from ChangeCostDB:%v\n", err)
		}
		for userChatId := range usersChatId {
			msg := tgbotapi.NewMessage(int64(userChatId), `Доступ к бирже возобновлен оповещения о достижении валюты фиксированной цены снова доступны`)
			_, err = b.bot.Send(msg)
			if err != nil {
				b.ErrorsMessage(placeSendMessageAboutError, err, msg)
			}
		}
	}
}
