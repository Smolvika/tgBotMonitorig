package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type infoUSD struct {
	CostUSD          float64
	changeCostPrUSD  float64
	changeCostRubUSD float64
	isIncreaseUSD    bool
}

type infoEUR struct {
	CostEUR          float64
	changeCostPrEUR  float64
	changeCostRubEUR float64
	isIncreaseEUR    bool
}
type infoCurrency struct {
	infoUSD
	infoEUR
}

func parsAllInfoCurrency() (infoCurrency, error) {

	usd, err := parsUSD()
	if err != nil {
		return infoCurrency{}, err
	}
	eur, err := parsEUR()
	if err != nil {
		return infoCurrency{}, err
	}
	fmt.Println(usd, eur)
	return infoCurrency{
		infoUSD: usd,
		infoEUR: eur,
	}, nil
}

func parsUSD() (infoUSD, error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	url := "https://bankiros.ru/ajax/moex-rate-current?currency_code=USDRUB_TOM"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return infoUSD{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return infoUSD{}, err
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	fmt.Println(string(body))
	if readErr != nil {
		return infoUSD{}, err
	}
	jsCur1 := jsUsd{}
	jsonErr := json.Unmarshal(body, &jsCur1)
	if jsonErr != nil {
		return infoUSD{}, err
	}

	info := infoUSD{}
	for _, j := range jsCur1.Data {
		info.CostUSD, err = strconv.ParseFloat(j.Last, 64)
		if err != nil {
			return info, err
		}
		info.changeCostRubUSD, err = strconv.ParseFloat(strings.ReplaceAll(j.Change, "-", ""), 64)
		if err != nil {
			return info, err
		}
		info.changeCostPrUSD, err = strconv.ParseFloat(j.ChangePercent, 64)
		if err != nil {
			return info, err
		}
		if j.Change[0] == '-' {
			info.isIncreaseUSD = true
		}
	}
	return info, nil
}

func parsEUR() (infoEUR, error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	url := "https://bankiros.ru/ajax/moex-rate-current?currency_code=EURRUB_TOM"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return infoEUR{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return infoEUR{}, err
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return infoEUR{}, err
	}
	jsCur1 := jsEur{}
	jsonErr := json.Unmarshal(body, &jsCur1)
	if jsonErr != nil {
		return infoEUR{}, err
	}
	info := infoEUR{}
	for _, j := range jsCur1.Data {
		info.CostEUR, err = strconv.ParseFloat(j.Last, 64)
		if err != nil {
			return info, err
		}
		info.changeCostRubEUR, err = strconv.ParseFloat(strings.ReplaceAll(j.Change, "-", ""), 64)
		if err != nil {
			return info, err
		}
		info.changeCostPrEUR, err = strconv.ParseFloat(j.ChangePercent, 64)
		if err != nil {
			return info, err
		}
		if j.Change[0] == '-' {
			info.isIncreaseEUR = true
		}
	}
	return info, nil
}
