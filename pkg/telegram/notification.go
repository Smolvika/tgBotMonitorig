package telegram

import (
	"fmt"
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
			allChatId, err := b.db.AllChatIdCostDB("USD")
			errorsWorkDB(ChatIdCostDB, giveInfo, err)
			for _, chatId := range allChatId {
				msg := tgbotapi.NewMessage(chatId.ChatId, msgText)
				_, err = b.bot.Send(msg)
				if err != nil {
					b.ErrorsMessage(placeSendMessageAboutCurrency, err, msg)
				}
			}
			if currencyNow.isIncreaseEUR {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			}
			allChatId, err = b.db.Cost.AllChatIdCostDB("EUR")
			errorsWorkDB(InfoCurrencyDB, giveInfo, err)
			for _, chatId := range allChatId {
				msg := tgbotapi.NewMessage(chatId.ChatId, msgText)
				_, err = b.bot.Send(msg)
				if err != nil {
					b.ErrorsMessage(placeSendMessageAboutCurrency, err, msg)
				}
			}
		}
	}
}

func (b *Bot) sendMessageAboutCostCurrency(currencyNow infoCurrency, updateRate time.Duration, errInfoBitcoinPars *error) {
	for {
		time.Sleep(updateRate)
		if *errInfoBitcoinPars == nil {
			users, err := b.db.ChangeCost.AllChatIdChangeCostDB("USD", currencyNow.CostUSD)
			errorsWorkDB(ChatIdChangeCostDB, giveInfo, err)
			for _, user := range users {
				msg := tgbotapi.NewMessage(int64(user.ChatId), "USD достиг стоимости в "+
					strconv.FormatFloat(user.Cost, 'f', 2, 64)+" ₽.\n"+
					"Сейчас цена составляет "+strconv.FormatFloat(currencyNow.CostUSD, 'f', 2, 64)+" ₽.")
				_, err = b.bot.Send(msg)
				if err != nil {
					b.ErrorsMessage(placeMessageNotCommand, err, msg)
				}
			}
			users, err = b.db.ChangeCost.AllChatIdChangeCostDB("EUR", currencyNow.CostUSD)
			errorsWorkDB(ChatIdCostDB, giveInfo, err)
			for _, user := range users {
				msg := tgbotapi.NewMessage(int64(user.ChatId), "EUR достиг стоимости в "+
					strconv.FormatFloat(user.Cost, 'f', 2, 64)+" ₽.\n"+
					"Сейчас цена составляет "+strconv.FormatFloat(currencyNow.CostUSD, 'f', 2, 64)+" ₽.")
				_, err = b.bot.Send(msg)
				if err != nil {
					b.ErrorsMessage(placeMessageNotCommand, err, msg)
				}
			}
		}

	}
}
