package telegram

import (
	"log"
	"time"
)

func updateInfoAboutCurrency(updateRate time.Duration, bitcoinNow *infoCurrency, errInfoBitcoinPars *error) {
	for {
		time.Sleep(updateRate)
		*(bitcoinNow), *errInfoBitcoinPars = parsAllInfoCurrency()
		if *errInfoBitcoinPars != nil {
			log.Printf("Problem with parsing currency information: %v\n", *errInfoBitcoinPars)
		}
	}
}
