# Dockerfile para desarrollo - Migrator de EduGo
# Este contenedor incluye el código fuente y herramientas de desarrollo
FROM golang:1.25-alpine

# Instalar herramientas necesarias
RUN apk add --no-cache \
    git \
    make \
    bash \
    postgresql-client \
    ca-certificates

# Directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum primero para aprovechar cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente del migrator
COPY . .

# El CMD ejecutará directamente el código Go
# El migrator clonará edugo-infrastructure y ejecutará las migraciones
CMD ["go", "run", "./cmd/main.go"]
