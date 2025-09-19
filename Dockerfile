FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar dependencias primero para aprovechar cache
COPY go.mod go.sum ./
RUN go mod download

# Copiar código y compilar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o admira ./cmd/api

# Imagen final minimal
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario y archivos de configuración
COPY --from=builder /app/admira .

EXPOSE 8080

CMD ["./admira"]