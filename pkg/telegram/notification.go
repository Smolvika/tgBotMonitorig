package telegram

import (
	"fmt"
	"github.com/Smolvika/tgBotMonitorig/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

func (b *Bot) sendMessageCurrency(currencyNow infoCurrency, updateRate time.Duration, errInfoBitcoinPars *error) {
	for {
		time.Sleep(updateRate)
		if *errInfoBitcoinPars == nil {
			var msgText string
			if currencyNow.isIncreaseUSD {
				msgText = fmt.Sprintf("Цена на данный момент: 1$ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostUSD, currencyNow.changeCostRubUSD, currencyNow.changeCostPrUSD)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1$ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostUSD, currencyNow.changeCostRubUSD, currencyNow.changeCostPrUSD)
			}
			allChatId, err := repository.AllChatIdCostDB("USD", b.db)
			errorsWorkDB(ChatIdCostDB, giveInfo, err)
			for _, chatId := range allChatId {
				msg := tgbotapi.NewMessage(chatId.ChatId, msgText)
				_, err = b.bot.Send(msg)
				if err != nil {
					errorsMessage(placeSendMessageAboutCurrency, err, msg, b.db)
				}
			}
			if currencyNow.isIncreaseEUR {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			}
			allChatId, err = repository.AllChatIdCostDB("EUR", b.db)
			errorsWorkDB(InfoCurrencyDB, giveInfo, err)
			for _, chatId := range allChatId {
				msg := tgbotapi.NewMessage(chatId.ChatId, msgText)
				_, err = b.bot.Send(msg)
				if err != nil {
					errorsMessage(placeSendMessageAboutCurrency, err, msg, b.db)
				}
			}
		}
	}
}

func (b *Bot) sendMessageAboutCostCurrency(currencyNow infoCurrency, updateRate time.Duration, errInfoBitcoinPars *error) {
	for {
		time.Sleep(updateRate)
		if *errInfoBitcoinPars == nil {
			users, err := repository.AllChatIdChangeCostDB("USD", currencyNow.CostUSD, b.db)
			errorsWorkDB(ChatIdChangeCostDB, giveInfo, err)
			for _, user := range users {
				msg := tgbotapi.NewMessage(int64(user.ChatId), "USD достиг стоимости в "+
					strconv.FormatFloat(user.Cost, 'f', 2, 64)+" ₽\n"+
					"Сейчас цена составляет "+strconv.FormatFloat(currencyNow.CostUSD, 'f', 2, 64)+" ₽.")
				_, err = b.bot.Send(msg)
				if err != nil {
					errorsMessage(placeMessageNotCommand, err, msg, b.db)
				}
			}
			users, err = repository.AllChatIdChangeCostDB("EUR", currencyNow.CostUSD, b.db)
			errorsWorkDB(ChatIdCostDB, giveInfo, err)
			for _, user := range users {
				msg := tgbotapi.NewMessage(int64(user.ChatId), "EUR достиг стоимости в "+
					strconv.FormatFloat(user.Cost, 'f', 2, 64)+" ₽\n"+
					"Сейчас цена составляет "+strconv.FormatFloat(currencyNow.CostUSD, 'f', 2, 64)+" ₽.")
				_, err = b.bot.Send(msg)
				if err != nil {
					errorsMessage(placeMessageNotCommand, err, msg, b.db)
				}
			}
		}

	}
}
