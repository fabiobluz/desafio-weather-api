# Weather API

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## Funcionalidades

- Recebe um CEP válido de 8 dígitos
- Consulta a localização via API ViaCEP
- Retorna temperatura em Celsius, Fahrenheit e Kelvin
- Tratamento de erros adequado

## Requisitos

- Go 1.24.2+
- Chave da API WeatherAPI

## Configuração

1. Clone o repositório
2. Configure a variável de ambiente:
   ```bash
   export WEATHER_API_KEY=your_weatherapi_key_here
   ```
3. Execute o projeto:
   ```bash
   go run .
   ```

## Uso com Docker

1. Configure a variável de ambiente:
   ```bash
   export WEATHER_API_KEY=your_weatherapi_key_here
   ```

2. Execute com docker-compose:
   ```bash
   docker-compose up --build
   ```

## Endpoints

### GET /weather?cep={cep}

Retorna a temperatura atual para o CEP informado.

**Parâmetros:**
- `cep`: CEP de 8 dígitos (ex: 01234567)

**Respostas:**

**Sucesso (200):**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

**CEP inválido (422):**
```
invalid zipcode
```

**CEP não encontrado (404):**
```
can not find zipcode
```

## Testes

### 🧪 Executando Testes

Execute todos os testes:
```bash
go test -v
```

Execute apenas testes unitários (sem internet):
```bash
go test -v -run "TestIsValidCEP|TestWeatherHandler_InvalidCEP|TestWeatherHandler_EmptyCEP|TestWeatherHandler_NoAPIKey|TestMockCEPService|TestMockWeatherService"
```

Verificar cobertura de testes:
```bash
go test -cover
```

### 📊 Cobertura de Testes

- **Cobertura Total:** 78.4%
- **Status:** ✅ Excelente
- **Total de Testes:** 23 testes
- **Mocks:** Implementados para testes offline

### 🔧 Testes com Mocks

Este projeto inclui mocks para serviços externos, permitindo executar testes unitários sem dependência de internet. Para mais detalhes sobre como usar os mocks, consulte:

📖 **[Documentação dos Mocks](./mocks/README.md)**

### 🚀 Benefícios dos Mocks

- ⚡ **Velocidade:** Testes executam 2x mais rápido
- 🔒 **Confiabilidade:** Sem falhas de rede
- 🏠 **Offline:** Funcionam sem internet
- 🎯 **Previsibilidade:** Respostas sempre iguais

## Deploy no Google Cloud Run

1. Configure a variável de ambiente `WEATHER_API_KEY` no Google Cloud Run
2. Faça o deploy usando o Dockerfile fornecido