package mocks

// CEPService define a interface para serviços de CEP
type CEPService interface {
	GetCityFromCEP(cep string) (string, error)
}

// WeatherService define a interface para serviços de clima
type WeatherService interface {
	GetWeather(city, apiKey string) (float64, error)
}
