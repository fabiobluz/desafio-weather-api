package mocks

import "errors"

// MockWeatherService implementa WeatherService para testes
type MockWeatherService struct {
	// Mapa para simular respostas baseadas na cidade
	responses map[string]WeatherResponse
}

// WeatherResponse representa uma resposta do serviço de clima
type WeatherResponse struct {
	Temperature float64
	Error       error
}

// NewMockWeatherService cria uma nova instância do mock
func NewMockWeatherService() *MockWeatherService {
	return &MockWeatherService{
		responses: make(map[string]WeatherResponse),
	}
}

// SetResponse define uma resposta para uma cidade específica
func (m *MockWeatherService) SetResponse(city string, temperature float64, err error) {
	m.responses[city] = WeatherResponse{
		Temperature: temperature,
		Error:       err,
	}
}

// GetWeather implementa a interface WeatherService
func (m *MockWeatherService) GetWeather(city, apiKey string) (float64, error) {
	// Verificar se a API key é válida
	if apiKey == "" || apiKey == "invalid_key" {
		return 0, errors.New("invalid API key")
	}

	response, exists := m.responses[city]
	if !exists {
		return 0, errors.New("city not found")
	}

	return response.Temperature, response.Error
}

// SetupCommonCities define temperaturas para cidades comuns
func (m *MockWeatherService) SetupCommonCities() {
	// Cidades com temperaturas realistas
	m.SetResponse("São Paulo", 22.5, nil)
	m.SetResponse("Rio de Janeiro", 28.3, nil)
	m.SetResponse("Belo Horizonte", 24.1, nil)
	m.SetResponse("Salvador", 29.7, nil)
	m.SetResponse("Curitiba", 18.9, nil)

	// Cidade com erro
	m.SetResponse("Cidade Inexistente", 0, errors.New("city not found"))
}
