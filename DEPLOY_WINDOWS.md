# Deploy no Windows - Weather API

Este guia mostra como fazer o deploy da Weather API no Google Cloud Run usando Windows.

## 📋 Pré-requisitos

1. **Google Cloud CLI instalado**
2. **Conta Google Cloud ativa**
3. **Projeto Google Cloud configurado**
4. **API WeatherAPI key**

## 🚀 Opções de Deploy no Windows

### Opção 1: Comandos Manuais (Recomendado)

Execute os comandos diretamente no terminal:

```bash
# 1. Configurar projeto
gcloud config set project SEU_PROJECT_ID

# 2. Habilitar APIs
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# 3. Deploy
gcloud run deploy weather-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=SUA_CHAVE_WEATHERAPI \
  --port 8080
```

### Opção 2: Google Cloud Console (Interface Gráfica)

Para usuários que preferem interface gráfica:

1. Acesse o [Google Cloud Console](https://console.cloud.google.com/)
2. Navegue para **Cloud Run**
3. Clique em **Criar Serviço**
4. Configure:
   - **Fonte:** Fonte de código
   - **Região:** us-central1
   - **Autenticação:** Permitir tráfego não autenticado
   - **Variáveis de ambiente:** WEATHER_API_KEY=sua_chave

### Opção 3: WSL (Windows Subsystem for Linux)

Se você tem WSL instalado, pode usar comandos Linux:

```bash
# No WSL
gcloud config set project SEU_PROJECT_ID
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud run deploy weather-api --source . --platform managed --region us-central1 --allow-unauthenticated --set-env-vars WEATHER_API_KEY=SUA_CHAVE_WEATHERAPI --port 8080
```

### Opção 4: GitHub Actions (CI/CD)

Para deploy automático, configure GitHub Actions:

```yaml
# .github/workflows/deploy.yml
name: Deploy to Cloud Run

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: google-github-actions/setup-gcloud@v0
      with:
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        project_id: ${{ secrets.GCP_PROJECT }}
    
    - run: gcloud run deploy weather-api --source . --region us-central1 --allow-unauthenticated --set-env-vars WEATHER_API_KEY=${{ secrets.WEATHER_API_KEY }}
```

## 🔧 Instalação do Google Cloud CLI no Windows

### Método 1: Instalador Oficial (Recomendado)

1. Baixe o instalador em: https://cloud.google.com/sdk/docs/install
2. Execute o instalador como administrador
3. Siga as instruções de instalação
4. Reinicie o terminal

### Método 2: Chocolatey

```powershell
# Instalar Chocolatey (se não tiver)
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Instalar Google Cloud CLI
choco install gcloudsdk
```

### Método 3: Scoop

```powershell
# Instalar Scoop (se não tiver)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Instalar Google Cloud CLI
scoop install gcloud
```

## 🎯 Exemplo Completo de Deploy

### Passo a Passo Detalhado

```bash
# 1. Fazer login no Google Cloud
gcloud auth login

# 2. Configurar projeto
gcloud config set project ${{PROJECT_ID}}

# 3. Habilitar APIs necessárias
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# 4. Verificar se está no diretório correto
ls -la  # Deve mostrar Dockerfile, main.go, etc.

# 5. Deploy com configurações básicas
gcloud run deploy weather-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=SUA_CHAVE_WEATHERAPI \
  --port 8080

# 6. Deploy com configurações avançadas (opcional)
gcloud run deploy weather-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=SUA_CHAVE_WEATHERAPI \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300 \
  --max-instances 10 \
  --min-instances 0

# 7. Obter URL do serviço
gcloud run services describe weather-api --region us-central1 --format 'value(status.url)'
```

## 📝 Configuração Inicial

Na primeira execução, você precisará:

1. **Fazer login no Google Cloud:**
   ```bash
   gcloud auth login
   ```

2. **Configurar o projeto padrão:**
   ```bash
   gcloud config set project SEU_PROJECT_ID
   ```

3. **Habilitar APIs necessárias:**
   ```bash
   gcloud services enable run.googleapis.com
   gcloud services enable cloudbuild.googleapis.com
   ```

## 🔍 Troubleshooting

### Problema: "gcloud não é reconhecido"

**Solução:**
- Verifique se o Google Cloud CLI está instalado
- Reinicie o terminal
- Verifique se o PATH está configurado corretamente

### Problema: "Dockerfile não encontrado"

**Solução:**
- Certifique-se de estar no diretório correto do projeto
- Verifique se o arquivo Dockerfile existe
- Execute `ls -la` para listar os arquivos

### Problema: "Erro de permissões"

**Solução:**
- Verifique se está logado: `gcloud auth list`
- Verifique se o projeto está correto: `gcloud config get-value project`
- Verifique se as APIs estão habilitadas

### Problema: "Erro de build"

**Solução:**
- Verifique os logs do build no Google Cloud Console
- Certifique-se de que o Dockerfile está correto
- Verifique se todas as dependências estão no go.mod

## 📊 Monitoramento

Após o deploy, você pode monitorar usando:

```bash
# Ver status do serviço
gcloud run services describe weather-api --region us-central1

# Ver logs
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=weather-api" --limit 50

# Listar serviços
gcloud run services list

# Ver métricas
gcloud monitoring metrics list --filter="resource.type=cloud_run_revision"
```

## 🧪 Testando a API

Após o deploy, teste a API:

```bash
# Obter URL do serviço
SERVICE_URL=$(gcloud run services describe weather-api --region us-central1 --format 'value(status.url)')

# Testar com CEP de São Paulo
curl "$SERVICE_URL/weather?cep=01310100"

# Testar com CEP do Rio de Janeiro
curl "$SERVICE_URL/weather?cep=20040020"
```

## 🎉 Próximos Passos

1. **Configurar domínio personalizado** (opcional)
2. **Implementar monitoramento** com Cloud Monitoring
3. **Configurar alertas** para falhas
4. **Implementar cache** para melhor performance
5. **Adicionar rate limiting** para proteção

## 💡 Dicas

- Use comandos manuais para ter controle total sobre o deploy
- Configure um projeto de teste primeiro
- Mantenha suas credenciais seguras
- Use variáveis de ambiente para configurações sensíveis
- O WSL oferece a melhor experiência para desenvolvimento
- Sempre teste a API após o deploy

## 🔄 Deploy Automático com GitHub Actions

Para automatizar o deploy, você pode usar GitHub Actions:

```yaml
# .github/workflows/deploy.yml
name: Deploy to Cloud Run

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: google-github-actions/setup-gcloud@v0
      with:
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        project_id: ${{ secrets.GCP_PROJECT }}
    
    - run: |
        gcloud services enable run.googleapis.com
        gcloud services enable cloudbuild.googleapis.com
        gcloud run deploy weather-api --source . --region us-central1 --allow-unauthenticated --set-env-vars WEATHER_API_KEY=${{ secrets.WEATHER_API_KEY }}
```

## 📚 Recursos Adicionais

- [Google Cloud CLI para Windows](https://cloud.google.com/sdk/docs/install)
- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [WSL Documentation](https://docs.microsoft.com/windows/wsl/)
- [GitHub Actions Documentation](https://docs.github.com/actions)