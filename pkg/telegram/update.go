package telegram

import (
	"context"
	"log"
	"time"
)

func updateInfoAboutCurrency(ctx context.Context, updateRate time.Duration, currencycoinNow *infoCurrency, errInfoCurrencyNowPars *error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(updateRate):
			*(currencycoinNow), *errInfoCurrencyNowPars = parsAllInfoCurrency()
			if *errInfoCurrencyNowPars != nil {
				log.Printf("Problem with parsing currency information: %v\n", *errInfoCurrencyNowPars)
			}
		}
	}
}
