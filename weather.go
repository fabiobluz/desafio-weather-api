package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func IsValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func GetCityFromCEP(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status returned from viacep")
	}
	type viaCEP struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro"`
	}
	var v viaCEP
	json.NewDecoder(resp.Body).Decode(&v)
	if v.Erro || v.Localidade == "" {
		return "", fmt.Errorf("not found")
	}
	return strings.TrimSpace(v.Localidade), nil
}

func GetWeather(city, apiKey string) (float64, error) {
	// URL encode the city name to handle special characters
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
		apiKey, encodedCity,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("weather not found")
	}
	type WeatherAPIResp struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	var wa WeatherAPIResp
	json.NewDecoder(resp.Body).Decode(&wa)
	return wa.Current.TempC, nil
}
