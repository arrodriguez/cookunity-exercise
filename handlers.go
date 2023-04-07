package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleTraces(c *gin.Context) {
	// Parse IP address from request body
	var requestBody struct {
		IP string `json:"ip"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get geolocation information
	geolocation, err := getGeolocation(requestBody.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get geolocation information"})
		return
	}

	// Get currency conversion rates
	currencyConversion, err := getCurrencyConversion(geolocation.Currencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get currency conversion rates"})
		return
	}

	// Calculate distance to USA
	distance := calculateDistance(geolocation.Lat, geolocation.Lon)

	// Construct response
	response := gin.H{
		"ip":              requestBody.IP,
		"name":            geolocation.Name,
		"code":            geolocation.Code,
		"lat":             geolocation.Lat,
		"lon":             geolocation.Lon,
		"currencies":      currencyConversion,
		"distance_to_usa": distance,
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

func handleStatistics(c *gin.Context) {
	// Get all traced countries
	var tracedCountries []string
	db.Find(&Trace{}).Pluck("country", &tracedCountries)

	// Count the number of traces for each country
	countryCounts := make(map[string]int)
	for _, country := range tracedCountries {
		countryCounts[country]++
	}

	// Find the most traced country
	mostTraced := struct {
		Country string `json:"country"`
		Value   int    `json:"value"`
	}{}
	for country, count := range countryCounts {
		if count > mostTraced.Value {
			mostTraced.Country = country
			mostTraced.Value = count
		}
	}

	// Find the longest distance from requested traces
	var traces []Trace
	db.Order("distance DESC").Find(&traces)
	longestDistance := struct {
		Country string  `json:"country"`
		Value   float64 `json:"value"`
	}{}
	if len(traces) > 0 {
		longestDistance.Country = traces[0].Country
		longestDistance.Value = traces[0].Distance
	}

	// Construct response
	response := gin.H{
		"most_traced":      mostTraced,
		"longest_distance": longestDistance,
	}

	// Return response
	c.JSON(http.StatusOK, response)
}
