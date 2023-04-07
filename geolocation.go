package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Geolocation struct {
	Name       string     `json:"name"`
	Code       string     `json:"code"`
	Lat        float64    `json:"lat"`
	Lon        float64    `json:"lon"`
	Currencies []Currency `json:"currencies"`
}

type Currency struct {
	Iso            string  `json:"iso"`
	Symbol         string  `json:"symbol"`
	ConversionRate float64 `json:"conversion_rate"`
}

func getGeolocation(ip string) (*Geolocation, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?key=%s", ip, ipApiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geolocation Geolocation
	if err := json.NewDecoder(resp.Body).Decode(&geolocation); err != nil {
		return nil, err
	}

	return &geolocation, nil
}
