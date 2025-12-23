# Servicios - EduGo Dev Environment

## Resumen de Servicios

| Servicio | Imagen | Puerto Local | Puerto Container |
|----------|--------|--------------|------------------|
| PostgreSQL | postgres:16-alpine | 5432 | 5432 |
| MongoDB | mongo:7.0 | 27017 | 27017 |
| RabbitMQ | rabbitmq:3.12-management-alpine | 5672, 15672 | 5672, 15672 |
| API Mobile | ghcr.io/edugogroup/edugo-api-mobile | 8081 | 8080 |
| API Admin | ghcr.io/edugogroup/edugo-api-administracion | 8082 | 8081 |
| Worker | ghcr.io/edugogroup/edugo-worker | - | - |
| Migrator | Build local desde ./migrator | - | - |

---

## PostgreSQL

Base de datos relacional principal.

### Configuración

```yaml
image: postgres:16-alpine
container_name: edugo-postgres
ports: 5432:5432
```

### Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `POSTGRES_DB` | edugo | Nombre de la base de datos |
| `POSTGRES_USER` | edugo | Usuario |
| `POSTGRES_PASSWORD` | edugo123 | Contraseña |

### Conexión Directa

```bash
docker exec -it edugo-postgres psql -U edugo -d edugo
```

---

## MongoDB

Base de datos de documentos para materiales y contenido procesado.

### Configuración

```yaml
image: mongo:7.0
container_name: edugo-mongodb
ports: 27017:27017
```

### Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `MONGO_USER` | edugo | Usuario root |
| `MONGO_PASSWORD` | edugo123 | Contraseña |
| `MONGO_DB` | edugo | Base de datos |

### Conexión Directa

```bash
docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 edugo --authSource admin
```

---

## RabbitMQ

Cola de mensajes para procesamiento asíncrono.

### Configuración

```yaml
image: rabbitmq:3.12-management-alpine
container_name: edugo-rabbitmq
ports:
  - 5672:5672   # AMQP
  - 15672:15672 # Management UI
```

### Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `RABBITMQ_USER` | edugo | Usuario |
| `RABBITMQ_PASSWORD` | edugo123 | Contraseña |

### UI de Administración

- **URL**: http://localhost:15672
- **Usuario**: edugo
- **Contraseña**: edugo123

### Colas Principales

| Cola | Propósito |
|------|-----------|
| `material.uploaded` | Materiales subidos pendientes de procesar |
| `assessment.attempt` | Intentos de evaluación |

---

## API Mobile

API principal para aplicaciones móviles y frontend web.

### Configuración

```yaml
image: ghcr.io/edugogroup/edugo-api-mobile:latest
container_name: edugo-api-mobile
ports: 8081:8080
```

### Endpoints Principales

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/health` | GET | Estado del servicio |
| `/v1/auth/login` | POST | Autenticación |
| `/v1/auth/register` | POST | Registro |
| `/v1/courses` | GET | Listar cursos |
| `/v1/materials` | POST | Subir material |
| `/swagger/index.html` | GET | Documentación Swagger |

### Variables Importantes

| Variable | Descripción |
|----------|-------------|
| `DATABASE_POSTGRES_*` | Conexión PostgreSQL |
| `DATABASE_MONGODB_URI` | Conexión MongoDB |
| `MESSAGING_RABBITMQ_URL` | Conexión RabbitMQ |
| `AUTH_JWT_SECRET` | Secreto para tokens JWT |

---

## API Administración

API para el panel de administración.

### Configuración

```yaml
image: ghcr.io/edugogroup/edugo-api-administracion:latest
container_name: edugo-api-administracion
ports: 8082:8081
```

### Endpoints Principales

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/health` | GET | Estado del servicio |
| `/api/v1/admin/users` | GET | Gestión de usuarios |
| `/api/v1/admin/institutions` | GET | Gestión de instituciones |
| `/swagger/index.html` | GET | Documentación Swagger |

### Variables Importantes

Las variables usan prefijo `EDUGO_ADMIN_*` para ser leídas por Viper.

---

## Worker

Servicio de procesamiento en background para PDFs y generación de contenido con IA.

### Configuración

```yaml
image: ghcr.io/edugogroup/edugo-worker:latest
container_name: edugo-worker
```

### Procesamiento

1. Escucha cola `material.uploaded` en RabbitMQ
2. Descarga y procesa el PDF
3. Usa OpenAI para generar resúmenes/preguntas
4. Guarda resultado en MongoDB

### Variables Importantes

| Variable | Descripción |
|----------|-------------|
| `OPENAI_API_KEY` | **Requerido** para procesamiento de IA |
| `NLP_PROVIDER` | openai |
| `NLP_MODEL` | gpt-4 |
| `NLP_MAX_TOKENS` | 2000 |
| `NLP_TEMPERATURE` | 0.7 |

---

## Migrator

Ejecuta migraciones de base de datos automáticamente al inicio.

### Configuración

```yaml
build:
  context: ../migrator
  dockerfile: Dockerfile
container_name: edugo-migrator
restart: "no"
```

### Comportamiento

- Se ejecuta una vez cuando `postgres` y `mongodb` están saludables
- Aplica migraciones de `edugo-infrastructure`
- Sale con código 0 al terminar (no se reinicia)

### Dependencias

Utiliza `edugo-infrastructure`:
- `github.com/EduGoGroup/edugo-infrastructure/postgres v0.12.0`
- `github.com/EduGoGroup/edugo-infrastructure/mongodb v0.11.0`

> **Última actualización:** Diciembre 2025

---

## Archivos Docker Compose Disponibles

| Archivo | Uso |
|---------|-----|
| `docker-compose.yml` | Configuración principal (todo incluido) |
| `docker-compose-apps.yml` | Solo aplicaciones (requiere infra externa) |
| `docker-compose-infrastructure.yml` | Solo bases de datos y RabbitMQ |
| `docker-compose-mock.yml` | Con servicios mock para testing |
| `docker-compose.migrate.yml` | Solo migrator |

---

## Perfiles de Docker Compose

El `setup.sh` soporta perfiles:

| Perfil | Servicios |
|--------|-----------|
| `full` | Todo (default) |
| `db-only` | PostgreSQL + MongoDB + RabbitMQ |
| `api-only` | Infraestructura + APIs (sin Worker) |
| `mobile-only` | Infraestructura + API Mobile |
| `admin-only` | Infraestructura + API Admin |
| `worker-only` | Infraestructura + Worker |

Uso:
```bash
./scripts/setup.sh --profile db-only
```

---

**Ver también:** [GUIA-RAPIDA.md](./GUIA-RAPIDA.md) para instrucciones de uso.
