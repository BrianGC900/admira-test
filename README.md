# Admira Backend - Prueba Técnica

Servicio en Go para procesamiento ETL de datos de Ads y CRM, generando métricas de marketing y revenue.

## Características

- Consumo de datos desde APIs externas (Ads y CRM)
- Transformación y cruce de datos por UTM parameters
- Cálculo de métricas: CPC, CPA, CVR, ROAS
- API REST para consulta de métricas
- Health checks y logging estructurado
- Dockerizado para fácil despliegue

## Datos de entrada
- Se utilizaron datos de entrada locales, ya que la URLs de API de mockio.io parecian estar caidas.

# Ejecución

## Local(Go)
- Clonar el repositorio
```bash
    git clone https://github.com/BrianGC900/admira-test
    cd admira-test
```

## Instalar dependencias
```bash
    go mod tidy
    go run main.go
```

## Con Docker
- Construir imagen y levantar servicio
```bash
    docker compose up --build
```
# Probar EndPoints 
### Probar mock ads
- curl http://localhost:3001

### Probar mock crm  
- curl http://localhost:3002

### Probar aplicación principal
- curl.exe -X POST http://localhost:8080/ingest/run
- curl.exe http://localhost:8080/healthz
- curl.exe http://localhost:8080/readyz
- curl.exe http://localhost:8080/metrics/channel

### Requisitos

- Go 1.21+
- Docker (opcional)

### Variables de entorno

Crear un archivo `.env` 

```env
ADS_API_URL=http://mock-ads:3001
CRM_API_URL=http://mock-crm:3002
SINK_URL=
SINK_SECRET=
PORT=8080
LOG_LEVEL=info
MAX_RETRIES=3
RETRY_BACKOFF_MS=1000
```
