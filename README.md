# Weather API

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## 🚀 API em Execução

A Weather API está rodando no Google Cloud Run e pode ser testada diretamente:

**🌐 URL da API:** https://weather-api-231779291153.us-central1.run.app

**🧪 Teste rápido:**
```bash
curl "https://weather-api-231779291153.us-central1.run.app/weather?cep=01310100"
```

**📊 Status:** ✅ Online e funcionando

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

**URL da API:** https://weather-api-231779291153.us-central1.run.app/weather

**Parâmetros:**
- `cep`: CEP de 8 dígitos (ex: 01234567)

**Exemplo de uso:**
```bash
curl "https://weather-api-231779291153.us-central1.run.app/weather?cep=01310100"
```

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
- **Total de Testes:** 17 testes
- **Mocks:** Implementados para testes offline

### 🔧 Testes com Mocks

Este projeto inclui mocks para serviços externos, permitindo executar testes unitários sem dependência de internet. Para mais detalhes sobre como usar os mocks, consulte:

📖 **[Documentação dos Mocks](./mocks/README.md)**

### 🚀 Benefícios dos Mocks

- ⚡ **Velocidade:** Testes executam 2x mais rápido
- 🔒 **Confiabilidade:** Sem falhas de rede
- 🏠 **Offline:** Funcionam sem internet
- 🎯 **Previsibilidade:** Respostas sempre iguais

## 🌐 API em Produção

A Weather API está atualmente rodando no Google Cloud Run e pode ser acessada publicamente:

### 📍 Informações do Serviço
- **URL:** https://weather-api-231779291153.us-central1.run.app
- **Região:** us-central1
- **Plataforma:** Google Cloud Run
- **Status:** ✅ Online e funcionando
- **Última atualização:** 28/09/2025

### 🧪 Testes Rápidos
```bash
# CEP de São Paulo
curl "https://weather-api-231779291153.us-central1.run.app/weather?cep=01310100"

# CEP do Rio de Janeiro
curl "https://weather-api-231779291153.us-central1.run.app/weather?cep=20040020"

# CEP de Brasília
curl "https://weather-api-231779291153.us-central1.run.app/weather?cep=70040900"
```

### 📊 Monitoramento
- **Logs:** Disponível no Google Cloud Console
- **Métricas:** Monitoramento automático do Cloud Run
- **Disponibilidade:** 99.9% SLA do Google Cloud Run

## Deploy no Google Cloud Run

### 🚀 Deploy Rápido

```bash
# Deploy básico
gcloud run deploy weather-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=sua_chave_weatherapi \
  --port 8080
```

### 📚 Documentação Completa

Para instruções detalhadas, configurações avançadas, scripts automatizados e troubleshooting, consulte:

📖 **[Guia Completo de Deploy](./DEPLOY.md)**

### 🖥️ Deploy no Windows

Para usuários do Windows, temos guias específicos:

📖 **[Guia de Deploy para Windows](./DEPLOY_WINDOWS.md)**

**Opções disponíveis:**
- Comandos manuais (recomendado)
- Google Cloud Console (interface gráfica)
- WSL (Windows Subsystem for Linux)
- GitHub Actions (CI/CD)

**Inclui:**
- Instalação do Google Cloud CLI
- Configurações avançadas
- Monitoramento e troubleshooting
- Otimizações de performance
- Configurações de segurança