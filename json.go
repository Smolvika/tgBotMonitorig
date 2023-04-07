package main

type jsEur struct {
	Error bool `json:"error"`
	Data  []struct {
		Id            string `json:"id"`
		DateParse     string `json:"date_parse"`
		ShortName     string `json:"short_name"`
		FullName      string `json:"full_name"`
		CurrencyId    string `json:"currency_id"`
		Info          string `json:"info"`
		LotSize       string `json:"lot_size"`
		Qty           string `json:"qty"`
		Code          string `json:"code"`
		Last          string `json:"last"`
		Open          string `json:"open"`
		Change        string `json:"change"`
		ChangePercent string `json:"change_percent"`
		Low           string `json:"low"`
		High          string `json:"high"`
	} `json:"data"`
	Params    string `json:"params"`
	DataGraph struct {
		Field1 string `json:"4"`
	} `json:"dataGraph"`
}

type jsUsd struct {
	Error bool `json:"error"`
	Data  []struct {
		Id            string `json:"id"`
		DateParse     string `json:"date_parse"`
		ShortName     string `json:"short_name"`
		FullName      string `json:"full_name"`
		CurrencyId    string `json:"currency_id"`
		Info          string `json:"info"`
		LotSize       string `json:"lot_size"`
		Qty           string `json:"qty"`
		Code          string `json:"code"`
		Last          string `json:"last"`
		Open          string `json:"open"`
		Change        string `json:"change"`
		ChangePercent string `json:"change_percent"`
		Low           string `json:"low"`
		High          string `json:"high"`
	} `json:"data"`
	Params    string `json:"params"`
	DataGraph struct {
		Field1 string `json:"1"`
	} `json:"dataGraph"`
}
