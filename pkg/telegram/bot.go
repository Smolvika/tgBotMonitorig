package telegram

import (
	"github.com/Smolvika/tgBotMonitorig/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	db  *repository.Repository
}

func NewBot(bot *tgbotapi.BotAPI, db *repository.Repository) *Bot {
	return &Bot{
		bot: bot,
		db:  db,
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s\n", b.bot.Self.UserName)
	updates := b.initUpdatesChannel()
	currencyNow, err := parsAllInfoCurrency()
	if err != nil {
		return err
	}
	b.handleUpdates(updates, currencyNow)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel, currencyNow infoCurrency) {
	var errInfoCurrencyPars error
	status := make(map[int64]string, 0)
	go func() {
		for update := range updates {
			if update.Message == nil && update.CallbackQuery == nil {
				continue
			} else if update.CallbackQuery != nil {
				b.isCallbackQuery(&update, status)
			} else if update.Message.IsCommand() {
				b.isCommandCase(&update, currencyNow, &errInfoCurrencyPars)
			} else if !update.Message.IsCommand() {
				b.isUsualMessage(&update, status)
			}
		}
	}()
	go b.sendMessageUserAboutError(&errInfoCurrencyPars)
	go updateInfoAboutCurrency(30*time.Second, &currencyNow, &errInfoCurrencyPars)
	go b.sendMessageAboutCostCurrency(currencyNow, 60*time.Second, &errInfoCurrencyPars)
	b.sendMessageCurrency(currencyNow, 60*time.Minute, &errInfoCurrencyPars)
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}
