package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !IsValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := GetCityFromCEP(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, "weather api key not set", http.StatusInternalServerError)
		return
	}

	tempC, err := GetWeather(city, apiKey)
	if err != nil {
		http.Error(w, "error fetching weather", http.StatusBadGateway)
		return
	}
	// Convert temperatures using correct formulas
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	resp := WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
