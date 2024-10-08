package main

import (
	"currency-exchange-service/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/currency_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected to the database.")

	fetchCurrencyRates()
}

func fetchCurrencyRates() {
	resp, err := http.Get("https://api.nbrb.by/exrates/rates?periodicity=0")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var rates []models.CurrencyRate
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		log.Fatal(err)
	}

	for _, rate := range rates {
		date, err := formatDate(rate.Date)
		if err == nil {
			rate.Date = date
			saveCurrencyRate(rate)
		} else {
			log.Printf("Error formatting date for rate: %v\n", err)
		}
	}
}

func formatDate(dateStr string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil
}

func saveCurrencyRate(rate models.CurrencyRate) {
	query := "INSERT INTO exchange_rates (date, currency_name, currency_abbreviation, scale, official_rate) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, rate.Date, rate.CurrencyName, rate.CurrencyAbbreviation, rate.Scale, rate.OfficialRate)
	if err != nil {
		log.Printf("Error saving rate for %s: %v\n", rate.Date, err)
	} else {
		log.Printf("Saved rate for %s\n", rate.Date)
	}
}
