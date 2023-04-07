package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

// Кнопки для выбора режима
var setRegime = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1", "trackingRegime1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "trackingRegime2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "trackingRegime3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "trackingRegime4"),
	),
)

// Вступительное сообщение
var (
	helloMessage = `Я слежу за ценой валюты USD и EUR .
/rateUSD - текущяя цена USD
/rateEUR - текущяя цена EUR

Также вы можете настроить систему мониторинга стоимоси валюты , а именно уведомления о цене и её изменении 2-х типов:
1) уведомления раз в час
2) уведомление при повышении цены 

/tracking - настроить оповещения
/stop_tracking - отменить все оповещения`
	trackingMessage = `Здесь вы можете настроить оповещения.
		На выбор предоставляются два режима:
		1.Ежечасное оповещение о цене USD
        2.Ежечасное оповещение о цене EUR
		3.Оповещение при достижении определенной цены USD
        4.Оповещение при достижении определенной цены EUR
		Для включения/настройки одного из режимов нажмите на кнопку с соответствующим номером`
)
var (
	eur = "EUR"
	usd = "USD"
)

// ___________________________________
// Обработка Команд
// ___________________________________
func isCommandCase(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, currencyNow infoCurrency, errInfoCurrencyNowPars *error) {
	var err error
	cmdText := update.Message.Command()
	switch cmdText {
	case "help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageCommand, err, msg, db)
		}
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
		err = addNewUser(int(update.Message.Chat.ID), db)
		errorsWorkDB(InfoBitcoinDB, addInfo, err)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageCommand, err, msg, db)
		}
	case "tracking":
		msgText := trackingMessage
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		//добавление кнопок
		msg.ReplyMarkup = setRegime
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageCommand, err, msg, db)
		}
	case "rateUSD":
		var msgText string
		if *errInfoCurrencyNowPars == nil {
			if currencyNow.isIncreaseUSD {
				msgText = fmt.Sprintf("Цена на данный момент: 1$ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostUSD, currencyNow.changeCostRubUSD, currencyNow.changeCostPrUSD)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1$ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostUSD, currencyNow.changeCostRubUSD, currencyNow.changeCostPrUSD)
			}
		} else {
			msgText = "В данный момент имеются проблемы с доступом к сайту биржи,\nвы сможете ознакомиться с курсом криптовалюты, как только проблема будет решена.\nПриносим извинения за доставленные неудобства."
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageCommand, err, msg, db)
		}
	case "rateEUR":
		var msgText string
		if *errInfoCurrencyNowPars == nil {
			if currencyNow.isIncreaseEUR {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена снизилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			} else {
				msgText = fmt.Sprintf("Цена на данный момент: 1€ = %v₽ \nЗа последние 24 часа цена повысилась на %v₽ (%v%%)", currencyNow.CostEUR, currencyNow.changeCostRubEUR, currencyNow.changeCostPrEUR)
			}
		} else {
			msgText = "В данный момент имеются проблемы с доступом к сайту биржи,\nвы сможете ознакомиться с курсом криптовалюты, как только проблема будет решена.\nПриносим извинения за доставленные неудобства."
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageCommand, err, msg, db)
		}
	case "stop_tracking":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Оповещения отключены ")
		err = deleteUserChatIdCostDB(int(update.Message.Chat.ID), db)
		errorsWorkDB(ChatIdCostDB, deleteInfo, err)
		err = deleteUserChatIdChangeCostDB(int(update.Message.Chat.ID), db)
		errorsWorkDB(ChatIdChangeCostDB, deleteInfo, err)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageCommand, err, msg, db)
		}
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageNotCommand, err, msg, db)
		}
	}
}

// ___________________________________
// Обработка Callback ответов
// ___________________________________
func isCallbackQuery(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, status map[int64]string) {
	var err error
	switch update.CallbackQuery.Data {
	case "trackingRegime1":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✔ Режим отслеживания включен")
		err = changeInformation(int(update.CallbackQuery.Message.Chat.ID), usd, db)
		errorsWorkDB(InfoBitcoinDB, changeInfo, err)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeCallbackQuery, err, msg, db)
		}
		//EUR
	case "trackingRegime2":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✔ Режим отслеживания включен")
		err = changeInformation(int(update.CallbackQuery.Message.Chat.ID), eur, db)
		errorsWorkDB(InfoBitcoinDB, changeInfo, err)
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeCallbackQuery, err, msg, db)
		}
	//USD
	case "trackingRegime3":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Укажите стоимость USD, о которой нужно сообщить для этого отправьте число в формате: '123.456'")
		//status
		status[update.CallbackQuery.Message.Chat.ID] = "CostUSD"
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeCallbackQuery, err, msg, db)
		}
		//EUR
	case "trackingRegime4":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Укажите стоимость EUR, о которой нужно сообщить для этого отправьте число в формате: '123.456'")
		//status
		status[update.CallbackQuery.Message.Chat.ID] = "CostEUR"
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeCallbackQuery, err, msg, db)
		}

	}
}

// ____________________________________
// Обработка обычных сообщений
// ___________________________________
func isUsualMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, status map[int64]string) {
	var err error
	switch status[update.Message.Chat.ID] {
	case "CostUSD": // установление цены для оповещений
		changeCost, ok := validAndPrepare(update.Message.Text)
		var msg tgbotapi.MessageConfig
		if ok {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Когда скачок цены USD  окажется больше чем %s₽ вы получите уведомление.", strconv.FormatFloat(changeCost, 'f', 2, 64)))
			//добавление в базу данных
			err = addUserChangeCostDB(int(update.Message.Chat.ID), usd, changeCost, db)
			errorsWorkDB(ChatIdCostDB, addInfo, err)
			delete(status, update.Message.Chat.ID)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
		}
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageNotCommand, err, msg, db)
		}
	case "CostEUR": // установление изменения цены для оповещений
		changeCost, ok := validAndPrepare(update.Message.Text)
		var msg tgbotapi.MessageConfig
		if ok {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Когда скачок цены EUR  окажется больше чем %s₽ вы получите уведомление.", strconv.FormatFloat(changeCost, 'f', 2, 64)))

			err = addUserChangeCostDB(int(update.Message.Chat.ID), eur, changeCost, db)
			errorsWorkDB(ChatIdChangeCostDB, addInfo, err)
			delete(status, update.Message.Chat.ID)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
		}
		_, err = bot.Send(msg)
		if err != nil {
			errorsMessage(placeMessageNotCommand, err, msg, db)
		}
	}
}
