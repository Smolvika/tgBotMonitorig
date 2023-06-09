package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

var setRegime = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1", "trackingRegime1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "trackingRegime2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "trackingRegime3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "trackingRegime4"),
	),
)
var (
	helloMessage = `Я слежу за ценой валюты USD и EUR .
/rate_usd - текущая цена USD
/rate_eur - текущая цена EUR

Также вы можете настроить систему мониторинга стоимоси валюты , а именно уведомления о цене и её изменении 2-х типов:
1) уведомления раз в час
2) уведомление при повышении цены 

/tracking - настроить оповещения
/stop_tracking - отменить все оповещения`
	trackingMessage = `Здесь вы можете настроить оповещения.
    На выбор предоставляются два режима:
  1. Ежечасное оповещение о цене USD
  2. Ежечасное оповещение о цене EUR
  3. Оповещение при достижении определенной цены USD
  4. Оповещение при достижении определенной цены EUR
  Для включения/настройки одного из режимов нажмите на кнопку с соответствующим номером`
)
var (
	eur = "EUR"
	usd = "USD"
)

func (b *Bot) isCommandCase(update *tgbotapi.Update, currencyNow infoCurrency, errInfoCurrencyNowPars *error) {
	start := time.Now()
	defer func() {
		observeRequest(time.Since(start), update.Message.Command())
	}()
	var err error
	cmdText := update.Message.Command()
	switch cmdText {
	case "help", "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageCommand, err, msg)
		}
	case "tracking":
		msgText := trackingMessage
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.ReplyMarkup = setRegime
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageCommand, err, msg)
		}
	case "rate_usd":
		var msgText string
		if *errInfoCurrencyNowPars == nil {
			if currencyNow.isIncreaseUSD {
				msgText = fmt.Sprintf("Цена на данный момент: 1$ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostUSD, currencyNow.changeCostRubUSD, currencyNow.changeCostPrUSD)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1$ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostUSD, currencyNow.changeCostRubUSD, currencyNow.changeCostPrUSD)
			}
		} else {
			msgText = "В данный момент имеются проблемы с доступом к сайту биржи,\nвы сможете ознакомиться с курсом валюты, как только проблема будет решена.\nПриносим извинения за доставленные неудобства."
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageCommand, err, msg)
		}
	case "rate_eur":
		var msgText string
		if *errInfoCurrencyNowPars == nil {
			if currencyNow.isIncreaseEUR {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			}
		} else {
			msgText = "В данный момент имеются проблемы с доступом к сайту биржи,\nвы сможете ознакомиться с курсом валюты, как только проблема будет решена.\nПриносим извинения за доставленные неудобства."
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageCommand, err, msg)
		}
	case "stop_tracking":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Оповещения отключены ")
		err = b.db.Cost.DeleteUserChatIdCostDB(int(update.Message.Chat.ID))
		errorsWorkDB(ChatIdCostDB, deleteInfo, err)
		err = b.db.ChangeCost.DeleteUserChatIdChangeCostDB(int(update.Message.Chat.ID))
		errorsWorkDB(ChatIdChangeCostDB, deleteInfo, err)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageCommand, err, msg)
		}
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageNotCommand, err, msg)
		}
	}
}

func (b *Bot) isCallbackQuery(update *tgbotapi.Update, status map[int64]string) {
	var err error
	switch update.CallbackQuery.Data {
	case "trackingRegime1":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✔ Режим отслеживания включен")
		err = b.db.Cost.AddNewUser(int(update.CallbackQuery.Message.Chat.ID), usd)
		errorsWorkDB(ChatIdCostDB, addInfo, err)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeCallbackQuery, err, msg)
		}
	case "trackingRegime2":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✔ Режим отслеживания включен")
		err = b.db.Cost.AddNewUser(int(update.CallbackQuery.Message.Chat.ID), eur)
		errorsWorkDB(ChatIdCostDB, addInfo, err)
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeCallbackQuery, err, msg)
		}
	case "trackingRegime3":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Укажите стоимость USD, о которой нужно сообщить для этого отправьте число в формате: '123.456'")
		status[update.CallbackQuery.Message.Chat.ID] = "CostUSD"
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeCallbackQuery, err, msg)
		}
	case "trackingRegime4":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Укажите стоимость EUR, о которой нужно сообщить для этого отправьте число в формате: '123.456'")
		status[update.CallbackQuery.Message.Chat.ID] = "CostEUR"
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeCallbackQuery, err, msg)
		}

	}
}

func (b *Bot) isUsualMessage(update *tgbotapi.Update, status map[int64]string) {
	var err error
	switch status[update.Message.Chat.ID] {
	case "CostUSD":
		changeCost, ok := validAndPrepare(update.Message.Text)
		var msg tgbotapi.MessageConfig
		if ok {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Когда скачок цены USD  окажется больше чем %s₽ вы получите уведомление.", strconv.FormatFloat(changeCost, 'f', 2, 64)))
			err = b.db.ChangeCost.AddUserChangeCostDB(int(update.Message.Chat.ID), usd, changeCost)
			errorsWorkDB(ChatIdChangeCostDB, addInfo, err)
			delete(status, update.Message.Chat.ID)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
		}
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageNotCommand, err, msg)
		}
	case "CostEUR":
		changeCost, ok := validAndPrepare(update.Message.Text)
		var msg tgbotapi.MessageConfig
		if ok {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Когда скачок цены EUR  окажется больше чем %s₽ вы получите уведомление.", strconv.FormatFloat(changeCost, 'f', 2, 64)))
			err = b.db.ChangeCost.AddUserChangeCostDB(int(update.Message.Chat.ID), eur, changeCost)
			errorsWorkDB(ChatIdChangeCostDB, addInfo, err)
			delete(status, update.Message.Chat.ID)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
		}
		_, err = b.bot.Send(msg)
		if err != nil {
			b.ErrorsMessage(placeMessageNotCommand, err, msg)
		}
	}
}
