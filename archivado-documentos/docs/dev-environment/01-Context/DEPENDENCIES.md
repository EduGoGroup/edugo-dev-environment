# DEPENDENCIES - Dev Environment

## Matriz de Dependencias

```
┌──────────────────────────────────────────────┐
│        Dev Environment (Orchestrator)        │
├──────────────────────────────────────────────┤
│ Dependencias de Sistema (Usuario)            │
│ ├─ Docker Desktop o Docker + Docker Compose │
│ ├─ 8GB RAM mínimo                            │
│ ├─ 20GB almacenamiento libre                 │
│ └─ Puertos 5432, 27017, 5672, 8080, 8081    │
│                                              │
│ Dependencias de Contenedores (Docker)        │
│ ├─ PostgreSQL:15-alpine                      │
│ ├─ MongoDB:7.0                               │
│ ├─ RabbitMQ:3.12-management-alpine           │
│ ├─ Redis:7-alpine (opcional)                 │
│ └─ Imágenes de APIs (build locales)          │
│                                              │
│ Dependencias de Scripts (Bash)               │
│ ├─ docker                                    │
│ ├─ docker-compose                            │
│ ├─ curl                                      │
│ ├─ psql (PostgreSQL client)                  │
│ ├─ mongosh (MongoDB client)                  │
│ └─ rabbitmq-diagnostics                      │
└──────────────────────────────────────────────┘
```

---

## Dependencias de Sistema

### Docker Desktop (recomendado)
```bash
# Versión mínima
Docker Desktop 4.0+ (incluye Docker Engine 20.10+ y Docker Compose 2.0+)

# Verificar
docker --version
docker-compose --version

# Instalación
# https://www.docker.com/products/docker-desktop
```

### Docker + Docker Compose (alternativa Linux)
```bash
# Docker Engine 20.10+
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Docker Compose 2.0+
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Requisitos de Hardware
```
Mínimo:
- RAM: 8GB
- CPU: 2 cores
- Almacenamiento: 20GB libre

Recomendado:
- RAM: 16GB (8GB para Docker, 8GB para sistema)
- CPU: 4+ cores
- Almacenamiento: 50GB libre
- SSD para mejor performance
```

### Puertos Requeridos
```
5432   - PostgreSQL
27017  - MongoDB
5672   - RabbitMQ AMQP
15672  - RabbitMQ Management UI
6379   - Redis (opcional)
8080   - API Mobile
8081   - API Admin
```

---

## Dependencias de Contenedores

### PostgreSQL 15 Alpine
```bash
# Imagen
postgres:15-alpine

# Configuración
- POSTGRES_USER=edugo_user
- POSTGRES_PASSWORD=edugo_pass
- POSTGRES_DB=edugo_mobile
- Volumen: postgres_data:/var/lib/postgresql/data

# Health check
HEALTHCHECK: pg_isready -U edugo_user
```

### MongoDB 7.0
```bash
# Imagen
mongo:7.0

# Configuración
- MONGO_INITDB_ROOT_USERNAME=admin
- MONGO_INITDB_ROOT_PASSWORD=admin
- Volumen: mongo_data:/data/db

# Health check
HEALTHCHECK: mongosh --eval "db.adminCommand('ping')"
```

### RabbitMQ 3.12 Management
```bash
# Imagen
rabbitmq:3.12-management-alpine

# Configuración
- RABBITMQ_DEFAULT_USER=guest
- RABBITMQ_DEFAULT_PASS=guest
- Puerto 5672 (AMQP)
- Puerto 15672 (Management UI)
- Volumen: rabbitmq_data:/var/lib/rabbitmq

# Health check
HEALTHCHECK: rabbitmq-diagnostics ping
```

### Redis 7 (Opcional)
```bash
# Imagen
redis:7-alpine

# Configuración
- Puerto 6379
- Volumen: redis_data:/data

# Health check
HEALTHCHECK: redis-cli ping
```

### APIs (Build Local)
```bash
# api-mobile
- Build: ../repos-separados/edugo-api-mobile/docker/Dockerfile
- Depende: PostgreSQL, MongoDB, RabbitMQ
- Puerto: 8080

# api-admin
- Build: ../repos-separados/edugo-api-administracion/docker/Dockerfile
- Depende: PostgreSQL
- Puerto: 8081

# worker
- Build: ../repos-separados/edugo-worker/docker/Dockerfile
- Depende: RabbitMQ, MongoDB, PostgreSQL
- No expone puerto (backend)
```

---

## Dependencias de Scripts

### curl (para health checks)
```bash
# Verificar instalación
which curl

# Instalar si falta
# macOS: brew install curl
# Ubuntu: sudo apt-get install curl
```

### PostgreSQL Client (psql)
```bash
# Verificar instalación
which psql

# Instalar si falta
# macOS: brew install postgresql
# Ubuntu: sudo apt-get install postgresql-client
```

### MongoDB Client (mongosh)
```bash
# Verificar instalación
which mongosh

# Instalar si falta
# https://www.mongodb.com/try/download/shell
```

### RabbitMQ Diagnostics
```bash
# Dentro del contenedor (ya instalado)
docker exec rabbitmq rabbitmq-diagnostics ping

# O instalar localmente:
# https://www.rabbitmq.com/cli.html
```

---

## Volúmenes

```yaml
volumes:
  postgres_data:
    driver: local
  mongo_data:
    driver: local
  rabbitmq_data:
    driver: local
  redis_data:
    driver: local
```

**Propósito:**
- Persistencia de datos entre restarts
- No se pierden datos si contenedor se detiene
- `docker-compose down -v` elimina volúmenes (reset)

---

## Network

```yaml
networks:
  default:
    name: edugo-network
    driver: bridge
```

**Características:**
- Los contenedores se comunican por nombre (postgres, mongo, rabbitmq)
- Aislamiento de otros servicios Docker
- Bridge network automático

---

## Configuración de Ejemplo

### .env (opcional, para overrides)
```bash
# Database
DB_USER=edugo_user
DB_PASSWORD=edugo_pass

# MongoDB
MONGO_USER=admin
MONGO_PASSWORD=admin

# OpenAI (para Worker)
OPENAI_API_KEY=sk-...

# AWS (para Worker)
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
```

---

## Checklist de Verificación Previa

```bash
# Antes de docker-compose up:

# 1. Docker disponible
docker --version

# 2. Docker Compose disponible
docker-compose --version

# 3. Puertos libres
lsof -i :5432    # PostgreSQL
lsof -i :27017   # MongoDB
lsof -i :5672    # RabbitMQ
lsof -i :8080    # API Mobile
lsof -i :8081    # API Admin

# 4. Espacio en disco
df -h | grep -E "/$|/home"  # Al menos 20GB libres

# 5. RAM disponible
free -h  # o 'vm_stat' en macOS
```

---

## Resolución de Problemas de Dependencias

### "docker: command not found"
```bash
# Docker no instalado
# Ir a https://www.docker.com/products/docker-desktop
# O ejecutar instalador Linux

# O verificar PATH
echo $PATH
which docker
```

### "docker-compose: command not found"
```bash
# Docker Compose no instalado
# Docker Desktop ya lo incluye

# Si instalaste Docker manualmente:
sudo curl -L "..." -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### "Insufficient memory"
```bash
# Docker no tiene suficiente RAM asignada

# Docker Desktop:
# Settings > Resources > Memory: aumentar a 8GB+

# Ver memoria actual usada
docker stats

# Liberar espacio
docker system prune -a
```

### Puerto ya en uso
```bash
# Ver qué proceso usa puerto
lsof -i :5432

# Liberar puerto
kill -9 <PID>

# O usar puerto diferente en docker-compose.yml
# ports:
#   - "5433:5432"  # Puerto local:container
```

### Volúmenes corruptos
```bash
# Ver volúmenes
docker volume ls

# Eliminar volumen específico
docker volume rm dev-environment_postgres_data

# O desde docker-compose
docker-compose down -v  # Elimina todos volúmenes del proyecto
```

---

## Checklist de Instalación Completa

```markdown
## Sistema
- [ ] Docker Desktop 4.0+ o Docker 20.10+ + Docker Compose 2.0+
- [ ] 8GB RAM mínimo
- [ ] 20GB almacenamiento libre
- [ ] Puertos libres (5432, 27017, 5672, 8080, 8081)

## Cliente Tools
- [ ] curl instalado
- [ ] PostgreSQL client (psql) instalado
- [ ] MongoDB client (mongosh) instalado

## Repositorio Dev Environment
- [ ] Clonado desde GitHub
- [ ] docker-compose.yml presente
- [ ] Scripts en directorio scripts/

## Pre-startup
- [ ] Verificar puertos libres
- [ ] Verificar RAM disponible
- [ ] Verificar espacio en disco

## Post-startup
- [ ] docker-compose ps muestra todos servicios UP
- [ ] ./scripts/health-check.sh pasa
- [ ] APIs responden en puertos correctos
```

---

## Actualizaciones de Dependencias

```bash
# Actualizar imágenes Docker
docker-compose pull

# Reconstruir imágenes locales
docker-compose build --no-cache

# Limpiar imágenes no usadas
docker image prune -a
```
