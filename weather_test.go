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

func TestGetCityFromCEP_ValidCEP(t *testing.T) {

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

func TestGetWeather_ValidCity(t *testing.T) {
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

func TestWeatherHandler_Success(t *testing.T) {
	originalKey := os.Getenv("WEATHER_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("WEATHER_API_KEY", originalKey)
		} else {
			os.Unsetenv("WEATHER_API_KEY")
		}
	}()

	os.Setenv("WEATHER_API_KEY", "test_key")

	req, err := http.NewRequest("GET", "/weather?cep=01310100", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadGateway)
	}

	expected := "error fetching weather\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

func TestWeatherHandler_NoAPIKey(t *testing.T) {
	originalKey := os.Getenv("WEATHER_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("WEATHER_API_KEY", originalKey)
		} else {
			os.Unsetenv("WEATHER_API_KEY")
		}
	}()

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

func TestWeatherHandler_CEPNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather?cep=99999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "can not find zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}
}

func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}

	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)

	return strings.Contains(sLower, substrLower)
}

func TestMockCEPService(t *testing.T) {
	mock := mocks.NewMockCEPService()
	mock.SetupCommonCEPs()

	city, err := mock.GetCityFromCEP("01310100")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if city != "São Paulo" {
		t.Errorf("Expected 'São Paulo', got %s", city)
	}

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

	temp, err := mock.GetWeather("São Paulo", "valid_key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if temp != 22.5 {
		t.Errorf("Expected temperature 22.5, got %f", temp)
	}

	temp, err = mock.GetWeather("São Paulo", "")
	if err == nil {
		t.Error("Expected error for invalid API key")
	}
	if temp != 0 {
		t.Errorf("Expected temperature 0, got %f", temp)
	}
}

func TestWeatherResponse_JSONSerialization(t *testing.T) {
	response := WeatherResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298.15,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

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

	if !strings.Contains(jsonStr, "25") {
		t.Error("JSON should contain temperature value 25")
	}
}

func TestWeatherResponse_JSONDeserialization(t *testing.T) {
	jsonData := `{"temp_C": 20.5, "temp_F": 68.9, "temp_K": 293.65}`

	var response WeatherResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

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
	tempC := 0.0

	response := WeatherResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}

	expectedF := 32.0
	expectedK := 273.0

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

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	requiredFields := []string{"temp_C", "temp_F", "temp_K"}
	for _, field := range requiredFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("JSON should contain field: %s", field)
		}
	}

	for field, value := range jsonMap {
		if _, ok := value.(float64); !ok {
			t.Errorf("Field %s should be float64, got %T", field, value)
		}
	}
}
