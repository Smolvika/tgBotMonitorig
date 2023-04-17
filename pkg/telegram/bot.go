package telegram

import (
	"context"
	"github.com/Smolvika/tgBotMonitorig/metrics"
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

func (b *Bot) Start(ctx context.Context) error {
	log.Printf("Authorized on account %s\n", b.bot.Self.UserName)
	updates := b.initUpdatesChannel()
	currencyNow, err := parsAllInfoCurrency()
	if err != nil {
		return err
	}
	go func() {
		if err := metrics.Listen("127.0.0.1:8082"); err != nil {
			log.Println(err)
		}
	}()
	b.handleUpdates(ctx, updates, currencyNow)
	return nil
}

func (b *Bot) handleUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel, currencyNow infoCurrency) {
	var errInfoCurrencyPars error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				if update.Message == nil && update.CallbackQuery == nil {
					continue
				}
				b.nonEmptyUpdate(update, currencyNow, errInfoCurrencyPars)
			}
		}
	}()
	go b.sendMessageUserAboutError(&errInfoCurrencyPars)
	go updateInfoAboutCurrency(ctx, 30*time.Second, &currencyNow, &errInfoCurrencyPars)
	go b.sendMessageAboutCostCurrency(ctx, currencyNow, 60*time.Second, &errInfoCurrencyPars)
	b.sendMessageCurrency(ctx, currencyNow, 60*time.Minute, &errInfoCurrencyPars)
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}
func (b *Bot) nonEmptyUpdate(update tgbotapi.Update, currencyNow infoCurrency, errInfoCurrencyPars error) {
	status := make(map[int64]string, 0)
	if update.CallbackQuery != nil {
		b.isCallbackQuery(&update, status)
	} else if update.Message.IsCommand() {
		b.isCommandCase(&update, currencyNow, &errInfoCurrencyPars)
	} else if !update.Message.IsCommand() {
		b.isUsualMessage(&update, status)
	}
}
