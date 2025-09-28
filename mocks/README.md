# Mocks para Testes Unitários

Este pacote contém mocks para os serviços externos utilizados na aplicação Weather API, permitindo executar testes unitários sem dependência de internet.

## 📁 Estrutura

```
mocks/
├── interfaces.go      # Interfaces para injeção de dependência
├── cep_mock.go        # Mock para serviço de CEP (ViaCEP)
├── weather_mock.go    # Mock para serviço de clima (WeatherAPI)
└── README.md          # Esta documentação
```

## 🔧 Interfaces

### CEPService
```go
type CEPService interface {
    GetCityFromCEP(cep string) (string, error)
}
```

### WeatherService
```go
type WeatherService interface {
    GetWeather(city, apiKey string) (float64, error)
}
```

## 🚀 Como Usar

### 1. MockCEPService

```go
// Criar instância do mock
mock := mocks.NewMockCEPService()

// Configurar respostas predefinidas
mock.SetupCommonCEPs()

// Configurar resposta customizada
mock.SetResponse("01310100", "São Paulo", nil)
mock.SetResponse("00000000", "", errors.New("CEP not found"))

// Usar o mock
city, err := mock.GetCityFromCEP("01310100")
```

### 2. MockWeatherService

```go
// Criar instância do mock
mock := mocks.NewMockWeatherService()

// Configurar respostas predefinidas
mock.SetupCommonCities()

// Configurar resposta customizada
mock.SetResponse("São Paulo", 22.5, nil)
mock.SetResponse("Rio de Janeiro", 28.3, nil)

// Usar o mock
temp, err := mock.GetWeather("São Paulo", "valid_key")
```

## 📊 CEPs Predefinidos

| CEP | Cidade | Status |
|-----|--------|--------|
| 01310100 | São Paulo | ✅ Válido |
| 20040020 | Rio de Janeiro | ✅ Válido |
| 30112000 | Belo Horizonte | ✅ Válido |
| 40070110 | Salvador | ✅ Válido |
| 80010000 | Curitiba | ✅ Válido |
| 00000000 | - | ❌ Inválido |
| 99999999 | - | ❌ Inválido |

## 🌡️ Cidades Predefinidas

| Cidade | Temperatura | Status |
|--------|-------------|--------|
| São Paulo | 22.5°C | ✅ Válido |
| Rio de Janeiro | 28.3°C | ✅ Válido |
| Belo Horizonte | 24.1°C | ✅ Válido |
| Salvador | 29.7°C | ✅ Válido |
| Curitiba | 18.9°C | ✅ Válido |
| Cidade Inexistente | - | ❌ Erro |

## 🧪 Exemplo de Teste

```go
func TestWeatherAPIWithMocks(t *testing.T) {
    // Configurar mocks
    cepService := mocks.NewMockCEPService()
    weatherService := mocks.NewMockWeatherService()
    
    // Configurar respostas
    cepService.SetResponse("01310100", "São Paulo", nil)
    weatherService.SetResponse("São Paulo", 22.5, nil)
    
    // Executar teste
    city, err := cepService.GetCityFromCEP("01310100")
    if err != nil {
        t.Fatal(err)
    }
    
    temp, err := weatherService.GetWeather(city, "valid_key")
    if err != nil {
        t.Fatal(err)
    }
    
    // Verificar resultado
    if temp != 22.5 {
        t.Errorf("Expected 22.5, got %f", temp)
    }
}
```

## ✅ Benefícios

- **Sem Internet**: Testes executam offline
- **Rápidos**: Sem latência de rede
- **Confiáveis**: Respostas previsíveis
- **Isolados**: Não dependem de serviços externos
- **Reproduzíveis**: Mesmos resultados sempre

## 🚀 Executando Testes com Mocks

```bash
# Executar apenas testes com mocks (sem internet)
go test -v -run "TestIsValidCEP|TestWeatherHandler_InvalidCEP|TestWeatherHandler_EmptyCEP|TestWeatherHandler_NoAPIKey|TestMockCEPService|TestMockWeatherService"

# Executar todos os testes
go test -v

# Verificar cobertura
go test -cover
```
