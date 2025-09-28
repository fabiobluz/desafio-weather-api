package mocks

import "errors"

// MockCEPService implementa CEPService para testes
type MockCEPService struct {
	// Mapa para simular respostas baseadas no CEP
	responses map[string]CEPResponse
}

// CEPResponse representa uma resposta do serviço de CEP
type CEPResponse struct {
	City  string
	Error error
}

// NewMockCEPService cria uma nova instância do mock
func NewMockCEPService() *MockCEPService {
	return &MockCEPService{
		responses: make(map[string]CEPResponse),
	}
}

// SetResponse define uma resposta para um CEP específico
func (m *MockCEPService) SetResponse(cep string, city string, err error) {
	m.responses[cep] = CEPResponse{
		City:  city,
		Error: err,
	}
}

// GetCityFromCEP implementa a interface CEPService
func (m *MockCEPService) GetCityFromCEP(cep string) (string, error) {
	response, exists := m.responses[cep]
	if !exists {
		return "", errors.New("CEP not found in mock")
	}
	return response.City, response.Error
}

// Predefined responses para CEPs comuns
func (m *MockCEPService) SetupCommonCEPs() {
	// CEPs válidos
	m.SetResponse("01310100", "São Paulo", nil)
	m.SetResponse("20040020", "Rio de Janeiro", nil)
	m.SetResponse("30112000", "Belo Horizonte", nil)
	m.SetResponse("40070110", "Salvador", nil)
	m.SetResponse("80010000", "Curitiba", nil)

	// CEPs inválidos
	m.SetResponse("00000000", "", errors.New("CEP not found"))
	m.SetResponse("99999999", "", errors.New("CEP not found"))
}
