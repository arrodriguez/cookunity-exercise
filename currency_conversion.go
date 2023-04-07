package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CurrencyConversion struct {
	Iso            string  `json:"iso"`
	Symbol         string  `json:"symbol"`
	ConversionRate float64 `json:"conversion_rate"`
}

func getCurrencyConversion(currencies []Currency) ([]CurrencyConversion, error) {
	var currencyConversions []CurrencyConversion

	for _, currency := range currencies {
		url := fmt.Sprintf("http://data.fixer.io/api/convert?access_key=%s&from=%s&to=USD&amount=1", fixerApiKey, currency.Iso)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var response struct {
			Success bool `json:"success"`
			Rates   struct {
				USD float64 `json:"USD"`
			} `json:"rates"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		conversionRate := response.Rates.USD * currency.ConversionRate

		currencyConversion := CurrencyConversion{
			Iso:            currency.Iso,
			Symbol:         currency.Symbol,
			ConversionRate: conversionRate,
		}

		currencyConversions = append(currencyConversions, currencyConversion)
	}

	return currencyConversions, nil
}
