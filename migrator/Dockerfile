# Etapa 1: Builder
FROM golang:1.24-alpine AS builder

# Instalar git (necesario para clonar repositorio de infraestructura)
RUN apk add --no-cache git

WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY cmd/ ./cmd/

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/main.go

# Etapa 2: Runtime
FROM alpine:latest

# Instalar dependencias mínimas
RUN apk add --no-cache \
    git \
    postgresql-client \
    ca-certificates

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
