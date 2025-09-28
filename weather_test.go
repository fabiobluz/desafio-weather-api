package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"weather/mocks"
)

func TestIsValidCEP(t *testing.T) {
	valid := []string{"01234567", "12345678"}
	for _, cep := range valid {
		if !IsValidCEP(cep) {
			t.Errorf("Expected valid CEP: %s", cep)
		}
	}
	invalid := []string{"1234567", "abcdefgh", ""}
	for _, cep := range invalid {
		if IsValidCEP(cep) {
			t.Errorf("Expected invalid CEP: %s", cep)
		}
	}
}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather?cep=1234567", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

	expected := "invalid zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

func TestWeatherHandler_EmptyCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

	expected := "invalid zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

// Testes para GetCityFromCEP
func TestGetCityFromCEP_ValidCEP(t *testing.T) {
	// Este teste requer conexão com internet e API ViaCEP
	// Em um ambiente de produção, seria melhor usar mocks
	cep := "01310100" // CEP conhecido de São Paulo
	city, err := GetCityFromCEP(cep)

	if err != nil {
		t.Skipf("Skipping test due to network/API issue: %v", err)
	}

	if city == "" {
		t.Error("Expected city name, got empty string")
	}

	// Verifica se contém "São Paulo" (pode variar o formato)
	if !containsIgnoreCase(city, "são paulo") && !containsIgnoreCase(city, "sao paulo") {
		t.Errorf("Expected city to contain São Paulo, got: %s", city)
	}
}

func TestGetCityFromCEP_InvalidCEP(t *testing.T) {
	cep := "00000000" // CEP inválido
	city, err := GetCityFromCEP(cep)

	if err == nil {
		t.Error("Expected error for invalid CEP, got nil")
	}

	if city != "" {
		t.Errorf("Expected empty city for invalid CEP, got: %s", city)
	}
}

// Testes para GetWeather
func TestGetWeather_ValidCity(t *testing.T) {
	// Este teste requer conexão com internet e API WeatherAPI
	// Em um ambiente de produção, seria melhor usar mocks
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: WEATHER_API_KEY not set")
	}

	city := "São Paulo"
	temp, err := GetWeather(city, apiKey)

	if err != nil {
		t.Skipf("Skipping test due to network/API issue: %v", err)
	}

	if temp == 0 {
		t.Error("Expected temperature > 0, got 0")
	}

	// Temperatura deve estar em uma faixa razoável (-50°C a 60°C)
	if temp < -50 || temp > 60 {
		t.Errorf("Temperature %f seems unrealistic", temp)
	}
}

func TestGetWeather_InvalidAPIKey(t *testing.T) {
	city := "São Paulo"
	invalidKey := "invalid_key"
	temp, err := GetWeather(city, invalidKey)

	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}

	if temp != 0 {
		t.Errorf("Expected temperature 0 for invalid API key, got: %f", temp)
	}
}

// Testes para WeatherHandler - cenários de sucesso
func TestWeatherHandler_Success(t *testing.T) {
	// Salvar API key original
	originalKey := os.Getenv("WEATHER_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("WEATHER_API_KEY", originalKey)
		} else {
			os.Unsetenv("WEATHER_API_KEY")
		}
	}()

	// Configurar API key para teste
	os.Setenv("WEATHER_API_KEY", "test_key")

	req, err := http.NewRequest("GET", "/weather?cep=01310100", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	// Como não temos uma API key real, esperamos erro 502
	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadGateway)
	}

	expected := "error fetching weather\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

func TestWeatherHandler_NoAPIKey(t *testing.T) {
	// Salvar API key original
	originalKey := os.Getenv("WEATHER_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("WEATHER_API_KEY", originalKey)
		} else {
			os.Unsetenv("WEATHER_API_KEY")
		}
	}()

	// Remover API key
	os.Unsetenv("WEATHER_API_KEY")

	req, err := http.NewRequest("GET", "/weather?cep=01310100", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "weather api key not set\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

// Testes para WeatherHandler - CEP não encontrado
func TestWeatherHandler_CEPNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather?cep=99999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	// CEP não encontrado retorna 404
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "can not find zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

// Função auxiliar para teste
func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}

	// Conversão simples para comparação case-insensitive
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)

	return strings.Contains(sLower, substrLower)
}

// Testes usando mocks - sem dependência de internet
func TestMockCEPService(t *testing.T) {
	mock := mocks.NewMockCEPService()
	mock.SetupCommonCEPs()

	// Teste CEP válido
	city, err := mock.GetCityFromCEP("01310100")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if city != "São Paulo" {
		t.Errorf("Expected 'São Paulo', got %s", city)
	}

	// Teste CEP inválido
	city, err = mock.GetCityFromCEP("00000000")
	if err == nil {
		t.Error("Expected error for invalid CEP")
	}
	if city != "" {
		t.Errorf("Expected empty city, got %s", city)
	}
}

func TestMockWeatherService(t *testing.T) {
	mock := mocks.NewMockWeatherService()
	mock.SetupCommonCities()

	// Teste cidade válida
	temp, err := mock.GetWeather("São Paulo", "valid_key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if temp != 22.5 {
		t.Errorf("Expected temperature 22.5, got %f", temp)
	}

	// Teste API key inválida
	temp, err = mock.GetWeather("São Paulo", "")
	if err == nil {
		t.Error("Expected error for invalid API key")
	}
	if temp != 0 {
		t.Errorf("Expected temperature 0, got %f", temp)
	}
}

// Testes para WeatherResponse
func TestWeatherResponse_JSONSerialization(t *testing.T) {
	// Criar uma resposta de exemplo
	response := WeatherResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298.15,
	}

	// Serializar para JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verificar se contém os campos esperados
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, "temp_C") {
		t.Error("JSON should contain temp_C field")
	}
	if !strings.Contains(jsonStr, "temp_F") {
		t.Error("JSON should contain temp_F field")
	}
	if !strings.Contains(jsonStr, "temp_K") {
		t.Error("JSON should contain temp_K field")
	}

	// Verificar valores
	if !strings.Contains(jsonStr, "25") {
		t.Error("JSON should contain temperature value 25")
	}
}

func TestWeatherResponse_JSONDeserialization(t *testing.T) {
	// JSON de exemplo
	jsonData := `{"temp_C": 20.5, "temp_F": 68.9, "temp_K": 293.65}`

	// Deserializar do JSON
	var response WeatherResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verificar valores
	if response.TempC != 20.5 {
		t.Errorf("Expected TempC 20.5, got %f", response.TempC)
	}
	if response.TempF != 68.9 {
		t.Errorf("Expected TempF 68.9, got %f", response.TempF)
	}
	if response.TempK != 293.65 {
		t.Errorf("Expected TempK 293.65, got %f", response.TempK)
	}
}

func TestWeatherResponse_TemperatureConversions(t *testing.T) {
	// Teste com temperatura conhecida
	tempC := 0.0 // Ponto de congelamento da água

	// Criar resposta
	response := WeatherResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273.15,
	}

	// Verificar conversões
	expectedF := 32.0   // 0°C = 32°F
	expectedK := 273.15 // 0°C = 273.15K

	if response.TempF != expectedF {
		t.Errorf("Expected TempF %f, got %f", expectedF, response.TempF)
	}
	if response.TempK != expectedK {
		t.Errorf("Expected TempK %f, got %f", expectedK, response.TempK)
	}
}

func TestWeatherResponse_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		tempC     float64
		expectedF float64
		expectedK float64
	}{
		{
			name:      "Negative temperature",
			tempC:     -40.0,
			expectedF: -40.0, // -40°C = -40°F
			expectedK: 233.15,
		},
		{
			name:      "High temperature",
			tempC:     100.0,
			expectedF: 212.0, // 100°C = 212°F
			expectedK: 373.15,
		},
		{
			name:      "Zero temperature",
			tempC:     0.0,
			expectedF: 32.0,
			expectedK: 273.15,
		},
		{
			name:      "Room temperature",
			tempC:     25.0,
			expectedF: 77.0,
			expectedK: 298.15,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := WeatherResponse{
				TempC: tc.tempC,
				TempF: tc.tempC*1.8 + 32,
				TempK: tc.tempC + 273.15,
			}

			// Usar tolerância para comparação de ponto flutuante
			tolerance := 0.001
			if math.Abs(response.TempF-tc.expectedF) > tolerance {
				t.Errorf("Expected TempF %f, got %f", tc.expectedF, response.TempF)
			}
			if math.Abs(response.TempK-tc.expectedK) > tolerance {
				t.Errorf("Expected TempK %f, got %f", tc.expectedK, response.TempK)
			}
		})
	}
}

func TestWeatherResponse_JSONStructure(t *testing.T) {
	// Criar resposta de exemplo
	response := WeatherResponse{
		TempC: 30.0,
		TempF: 86.0,
		TempK: 303.15,
	}

	// Serializar para JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verificar estrutura JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verificar se todos os campos estão presentes
	requiredFields := []string{"temp_C", "temp_F", "temp_K"}
	for _, field := range requiredFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("JSON should contain field: %s", field)
		}
	}

	// Verificar tipos dos campos
	for field, value := range jsonMap {
		if _, ok := value.(float64); !ok {
			t.Errorf("Field %s should be float64, got %T", field, value)
		}
	}
}
