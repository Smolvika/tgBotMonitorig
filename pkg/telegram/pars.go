package telegram

import (
	"encoding/json"
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
	client := http.Client{Timeout: 3 * time.Second}
	usd, err := parsUSD(&client)
	if err != nil {
		return infoCurrency{
			infoUSD: infoUSD{},
			infoEUR: infoEUR{},
		}, err
	}
	eur, err := parsEUR(&client)
	if err != nil {
		return infoCurrency{
			infoUSD: infoUSD{},
			infoEUR: infoEUR{},
		}, err
	}
	return infoCurrency{
		infoUSD: usd,
		infoEUR: eur,
	}, nil
}

func parsUSD(client *http.Client) (infoUSD, error) {
	url := "https://bankiros.ru/ajax/moex-rate-current?currency_code=USDRUB_TOM"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if err != nil {
		return infoUSD{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return infoUSD{}, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	if readErr != nil {
		return infoUSD{}, err
	}
	jsCur1 := jsCurrency{}
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

func parsEUR(client *http.Client) (infoEUR, error) {
	url := "https://bankiros.ru/ajax/moex-rate-current?currency_code=EURRUB_TOM"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return infoEUR{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	res, err := client.Do(req)
	if err != nil {
		return infoEUR{}, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	if readErr != nil {
		return infoEUR{}, err
	}
	jsCur1 := jsCurrency{}
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
