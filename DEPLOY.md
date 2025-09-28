# Deploy no Google Cloud Run

Este guia mostra como fazer o deploy da Weather API no Google Cloud Run usando o Google Cloud CLI.

## 📋 Pré-requisitos

1. **Google Cloud CLI instalado**
2. **Conta Google Cloud ativa**
3. **Projeto Google Cloud configurado**
4. **API WeatherAPI key**

## 🚀 Passo a Passo

### 1. Configurar Google Cloud CLI

```bash
# Fazer login no Google Cloud
gcloud auth login

# Configurar projeto (substitua PROJECT_ID pelo seu ID)
gcloud config set project PROJECT_ID

# Habilitar APIs necessárias
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
```

### 2. Configurar Variáveis de Ambiente

```bash
# Definir variáveis de ambiente
export PROJECT_ID="seu-projeto-id"
export SERVICE_NAME="weather-api"
export REGION="us-central1"
export WEATHER_API_KEY="sua-chave-weatherapi"
```

### 3. Deploy usando Dockerfile

```bash
# Build e deploy da aplicação
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080
```

### 4. Deploy usando Container Registry

```bash
# Build da imagem
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# Deploy da imagem
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080
```

### 5. Deploy usando Artifact Registry (Recomendado)

```bash
# Criar repositório no Artifact Registry
gcloud artifacts repositories create weather-api-repo \
  --repository-format=docker \
  --location=$REGION

# Configurar Docker para Artifact Registry
gcloud auth configure-docker $REGION-docker.pkg.dev

# Build e push da imagem
docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/weather-api-repo/$SERVICE_NAME .
docker push $REGION-docker.pkg.dev/$PROJECT_ID/weather-api-repo/$SERVICE_NAME

# Deploy da imagem
gcloud run deploy $SERVICE_NAME \
  --image $REGION-docker.pkg.dev/$PROJECT_ID/weather-api-repo/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080
```

## 🔧 Configurações Avançadas

### Deploy com Configurações Personalizadas

```bash
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300 \
  --max-instances 10 \
  --min-instances 0
```

### Deploy com Autenticação

```bash
# Deploy com autenticação obrigatória
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --no-allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080
```

## 📊 Monitoramento

### Verificar Status do Deploy

```bash
# Listar serviços
gcloud run services list

# Ver detalhes do serviço
gcloud run services describe $SERVICE_NAME --region $REGION

# Ver logs
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=$SERVICE_NAME" --limit 50
```

### Testar a API

```bash
# Obter URL do serviço
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region $REGION --format 'value(status.url)')

# Testar endpoint
curl "$SERVICE_URL/weather?cep=01310100"
```

## 🛠️ Scripts de Automação

### Script de Deploy Completo

```bash
#!/bin/bash
# deploy.sh

set -e

# Configurações
PROJECT_ID="seu-projeto-id"
SERVICE_NAME="weather-api"
REGION="us-central1"
WEATHER_API_KEY="sua-chave-weatherapi"

echo "🚀 Iniciando deploy da Weather API..."

# Configurar projeto
gcloud config set project $PROJECT_ID

# Habilitar APIs
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# Deploy
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080

# Obter URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region $REGION --format 'value(status.url)')

echo "✅ Deploy concluído!"
echo "🌐 URL da API: $SERVICE_URL"
echo "🧪 Teste: curl \"$SERVICE_URL/weather?cep=01310100\""
```

### Script PowerShell (Windows)

```powershell
# deploy.ps1

$PROJECT_ID = "seu-projeto-id"
$SERVICE_NAME = "weather-api"
$REGION = "us-central1"
$WEATHER_API_KEY = "sua-chave-weatherapi"

Write-Host "🚀 Iniciando deploy da Weather API..." -ForegroundColor Green

# Configurar projeto
gcloud config set project $PROJECT_ID

# Habilitar APIs
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# Deploy
gcloud run deploy $SERVICE_NAME `
  --source . `
  --platform managed `
  --region $REGION `
  --allow-unauthenticated `
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY `
  --port 8080

# Obter URL
$SERVICE_URL = gcloud run services describe $SERVICE_NAME --region $REGION --format 'value(status.url)'

Write-Host "✅ Deploy concluído!" -ForegroundColor Green
Write-Host "🌐 URL da API: $SERVICE_URL" -ForegroundColor Cyan
Write-Host "🧪 Teste: curl `"$SERVICE_URL/weather?cep=01310100`"" -ForegroundColor Yellow
```

## 🔍 Troubleshooting

### Problemas Comuns

1. **Erro de autenticação:**
   ```bash
   gcloud auth login
   gcloud auth application-default login
   ```

2. **Erro de permissões:**
   ```bash
   gcloud projects add-iam-policy-binding $PROJECT_ID \
     --member="user:seu-email@gmail.com" \
     --role="roles/run.admin"
   ```

3. **Erro de build:**
   ```bash
   # Verificar logs do build
   gcloud logging read "resource.type=build" --limit 10
   ```

4. **Erro de variáveis de ambiente:**
   ```bash
   # Verificar variáveis configuradas
   gcloud run services describe $SERVICE_NAME --region $REGION
   ```

## 📈 Otimizações

### Performance

```bash
# Deploy com configurações otimizadas
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080 \
  --memory 1Gi \
  --cpu 2 \
  --timeout 300 \
  --max-instances 100 \
  --min-instances 1 \
  --concurrency 1000
```

### Segurança

```bash
# Deploy com HTTPS obrigatório
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --port 8080 \
  --ingress all \
  --no-allow-unauthenticated
```

## 🎯 Próximos Passos

1. **Configurar CI/CD** com GitHub Actions
2. **Implementar monitoramento** com Cloud Monitoring
3. **Configurar alertas** para falhas
4. **Implementar cache** para melhor performance
5. **Adicionar rate limiting** para proteção

## 📚 Recursos Adicionais

- [Google Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Google Cloud CLI Reference](https://cloud.google.com/sdk/gcloud/reference)
- [Cloud Run Pricing](https://cloud.google.com/run/pricing)
- [Best Practices](https://cloud.google.com/run/docs/best-practices)
