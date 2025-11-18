# Etapa 1: Builder - Compila el código Go
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/main.go

# Etapa 2: Runtime - Alpine minimalista
FROM alpine:latest

# Instalar dependencias básicas
RUN apk add --no-cache \
    git \
    postgresql-client \
    ca-certificates \
    bash

# Instalar Go runtime (necesario para ejecutar go run en las migraciones)
COPY --from=builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR /app

# Copiar binario compilado
COPY --from=builder /app/migrator .

# Ejecutar migrator
CMD ["./migrator"]
