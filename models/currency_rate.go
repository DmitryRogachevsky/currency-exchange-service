package models

type CurrencyRate struct {
	Date                 string  `json:"Date"`
	CurrencyName         string  `json:"Cur_Name"`
	CurrencyAbbreviation string  `json:"Cur_Abbreviation"`
	Scale                int     `json:"Cur_Scale"`
	OfficialRate         float64 `json:"Cur_OfficialRate"`
}
