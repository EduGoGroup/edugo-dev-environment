# PROJECT OVERVIEW - Dev Environment

## Información General

**Proyecto:** EduGo Dev Environment  
**Tipo:** Infraestructura Docker Compose  
**Tecnologías:** Docker, Docker Compose, Bash Scripts  
**Especificación de Origen:** spec-05-dev-environment  
**Estado:** En Desarrollo (Sprint 1/3)

---

## Propósito del Proyecto

Dev Environment es el orquestador de infraestructura para desarrollo local de EduGo. Proporciona un `docker-compose.yml` que levanta toda la pila de servicios necesarios (PostgreSQL, MongoDB, RabbitMQ, APIs, Worker) con un simple comando.

### Responsabilidades Principales
- Definir topología de servicios en Docker Compose
- Configurar volúmenes para persistencia
- Gestionar redes y comunicación entre contenedores
- Proporcionar scripts de setup/teardown
- Seed databases con datos de prueba
- Configurar health checks
- Documentar ambiente
- Perfil de desarrollo con hot-reload

---

## Estructura del Proyecto

```
dev-environment/
├── docker-compose.yml           # Definición principal
├── docker-compose.override.yml  # Overrides locales
├── docker/
│   ├── postgres/
│   │   ├── Dockerfile
│   │   ├── init.sql
│   │   └── seed.sql
│   ├── mongo/
│   │   ├── Dockerfile
│   │   ├── init.js
│   │   └── seed.js
│   ├── rabbitmq/
│   │   └── rabbitmq.conf
│   └── redis/
│       └── Dockerfile
├── scripts/
│   ├── setup.sh              # Iniciar ambiente
│   ├── teardown.sh           # Detener ambiente
│   ├── seed-data.sh          # Cargar datos de prueba
│   ├── clean.sh              # Limpiar volúmenes
│   ├── logs.sh               # Ver logs
│   └── health-check.sh       # Verificar salud
├── docs/
│   ├── SETUP.md              # Guía de instalación
│   ├── TROUBLESHOOTING.md    # Solución de problemas
│   └── SERVICES.md           # Detalle de servicios
└── README.md
```

---

## Docker Compose Stack

### Servicios Principales

```yaml
version: '3.8'

services:
  # PostgreSQL - Base de datos relacional
  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=edugo_user
      - POSTGRES_PASSWORD=edugo_pass
      - POSTGRES_DB=edugo_mobile
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./docker/postgres/init.sql:/docker-entrypoint-initdb.d/01-init.sql
      - ./docker/postgres/seed.sql:/docker-entrypoint-initdb.d/02-seed.sql

  # MongoDB - Base de datos de documentos
  mongo:
    image: mongo:7.0
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=admin
      - MONGO_INITDB_ROOT_PASSWORD=admin
    volumes:
      - mongo_data:/data/db
      - ./docker/mongo/init.js:/docker-entrypoint-initdb.d/01-init.js

  # RabbitMQ - Message broker
  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      - RABBITMQ_DEFAULT_USER=guest
      - RABBITMQ_DEFAULT_PASS=guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
      - ./docker/rabbitmq/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf

  # Redis - Caché (opcional)
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # API Mobile (8080)
  api-mobile:
    build:
      context: ../repos-separados/edugo-api-mobile
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=edugo_user
      - DB_PASSWORD=edugo_pass
      - DB_NAME=edugo_mobile
      - MONGO_URI=mongodb://admin:admin@mongo:27017
      - MONGO_DB_NAME=edugo_assessments
      - RABBITMQ_HOST=rabbitmq
      - RABBITMQ_PORT=5672
      - LOG_LEVEL=debug
    depends_on:
      postgres:
        condition: service_healthy
      mongo:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  # API Admin (8081)
  api-admin:
    build:
      context: ../repos-separados/edugo-api-administracion
      dockerfile: docker/Dockerfile
    ports:
      - "8081:8081"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=edugo_user
      - DB_PASSWORD=edugo_pass
      - DB_NAME=edugo_admin
      - LOG_LEVEL=debug
    depends_on:
      postgres:
        condition: service_healthy

  # Worker (Async processing)
  worker:
    build:
      context: ../repos-separados/edugo-worker
      dockerfile: docker/Dockerfile
    environment:
      - RABBITMQ_HOST=rabbitmq
      - RABBITMQ_PORT=5672
      - MONGO_URI=mongodb://admin:admin@mongo:27017
      - MONGO_DB_NAME=edugo_assessments
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
      - LOG_LEVEL=debug
    depends_on:
      - rabbitmq
      - mongo

volumes:
  postgres_data:
  mongo_data:
  rabbitmq_data:
  redis_data:

networks:
  default:
    name: edugo-network
```

---

## Scripts Principales

### setup.sh - Iniciar Ambiente
```bash
#!/bin/bash

echo "Iniciando EduGo Dev Environment..."

# 1. Construir imágenes (si es necesario)
docker-compose build

# 2. Iniciar servicios
docker-compose up -d

# 3. Esperar que servicios estén listos
echo "Esperando que servicios estén saludables..."
sleep 10

# 4. Ejecutar health check
./scripts/health-check.sh

echo "✅ Ambiente listo!"
echo "PostgreSQL: localhost:5432"
echo "MongoDB: localhost:27017"
echo "RabbitMQ: http://localhost:15672"
echo "API Mobile: http://localhost:8080"
echo "API Admin: http://localhost:8081"
```

### seed-data.sh - Cargar Datos de Prueba
```bash
#!/bin/bash

echo "Cargando datos de prueba..."

# 1. PostgreSQL seeds
psql -h localhost -U edugo_user -d edugo_mobile \
  -f ./docker/postgres/seed.sql

# 2. MongoDB seeds
mongosh --host localhost:27017 \
  -u admin -p admin \
  < ./docker/mongo/seed.js

# 3. RabbitMQ declarations
rabbitmqctl declare_exchange assessment.requests type:direct durable:true
rabbitmqctl declare_queue worker.assessment.requests durable:true
# ... más declarations

echo "✅ Datos de prueba cargados"
```

### logs.sh - Ver Logs
```bash
#!/bin/bash

if [ -z "$1" ]; then
  docker-compose logs -f
else
  docker-compose logs -f $1
fi
```

---

## Perfil de Desarrollo (Hot Reload)

### docker-compose.override.yml
```yaml
services:
  api-mobile:
    build:
      context: ../repos-separados/edugo-api-mobile
      dockerfile: docker/Dockerfile.dev
    volumes:
      - ../repos-separados/edugo-api-mobile:/app
    command: go run ./cmd/api-mobile

  api-admin:
    build:
      context: ../repos-separados/edugo-api-administracion
      dockerfile: docker/Dockerfile.dev
    volumes:
      - ../repos-separados/edugo-api-administracion:/app
    command: go run ./cmd/api-admin

  worker:
    build:
      context: ../repos-separados/edugo-worker
      dockerfile: docker/Dockerfile.dev
    volumes:
      - ../repos-separados/edugo-worker:/app
    command: go run ./cmd/worker
```

### Dockerfile.dev (para cada API)
```dockerfile
FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

# Instalar air para hot reload
RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]
```

---

## Datos de Prueba

### Escuela de Ejemplo (seed.sql)
```sql
-- Escuela
INSERT INTO schools (name, city) VALUES ('Instituto EduGo', 'Madrid');

-- Facultad
INSERT INTO academic_units (school_id, parent_id, type, name)
VALUES (1, NULL, 'faculty', 'Facultad de Ingeniería');

-- Departamento
INSERT INTO academic_units (school_id, parent_id, type, name)
VALUES (1, 1, 'department', 'Departamento de Sistemas');

-- Usuario de prueba
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES ('test@edugo.com', 'hashed_password', 'Test', 'User');
```

### Datos MongoDB (seed.js)
```javascript
use edugo_assessments;

db.evaluation_results.insertOne({
  evaluation_id: 1,
  student_id: 42,
  answers: [],
  total_score: 0,
  max_score: 100,
  submitted_at: new Date()
});
```

---

## Health Checks

Cada servicio implementa health checks:

```bash
# PostgreSQL
curl -f http://localhost:5432/health || exit 1

# MongoDB
mongosh --host localhost:27017 --eval "db.adminCommand('ping')"

# RabbitMQ
rabbitmq-diagnostics ping

# APIs
curl -f http://localhost:8080/api/v1/health || exit 1
curl -f http://localhost:8081/api/v1/health || exit 1
```

---

## Troubleshooting

### Puertos en uso
```bash
# Ver qué proceso usa puerto 5432
lsof -i :5432

# Liberar puerto
kill -9 <PID>
```

### Volúmenes dañados
```bash
# Limpiar todo
./scripts/clean.sh

# O manualmente
docker-compose down -v
docker volume prune
```

### Logs
```bash
# Ver todos los logs
./scripts/logs.sh

# Log de servicio específico
./scripts/logs.sh postgres
```

---

## Sprint Planning (3 Sprints)

| Sprint | Funcionalidad | Duración |
|--------|---------------|----------|
| 1 | docker-compose.yml + scripts básicos | 1 semana |
| 2 | Health checks + seed data | 1 semana |
| 3 | Perfil dev + documentación | 1 semana |

---

## Requisitos Previos

- Docker Desktop (o Docker + Docker Compose)
- Docker versión 20.10+
- Docker Compose versión 2.0+
- 8GB RAM disponible
- 20GB almacenamiento libre

---

## Contacto y Referencias

- **Repositorio GitHub:** https://github.com/EduGoGroup/edugo-dev-environment
- **Especificación Completa:** docs/ESTADO_PROYECTO.md (repo análisis)
- **Documentación Técnica:** Este directorio (01-Context/)
