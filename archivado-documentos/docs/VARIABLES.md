# Variables de Entorno - EduGo Dev Environment

Ver documentación completa en: [`../docker/.env.example`](../docker/.env.example)

---

## PostgreSQL

| Variable | Default | Descripción |
|----------|---------|-------------|
| `POSTGRES_DB` | `edugo` | Nombre de base de datos |
| `POSTGRES_USER` | `edugo` | Usuario de PostgreSQL |
| `POSTGRES_PASSWORD` | `edugo123` | Contraseña |
| `POSTGRES_PORT` | `5432` | Puerto expuesto |

---

## MongoDB

| Variable | Default | Descripción |
|----------|---------|-------------|
| `MONGO_USER` | `edugo` | Usuario root |
| `MONGO_PASSWORD` | `edugo123` | Contraseña |
| `MONGO_DB` | `edugo` | Base de datos |
| `MONGO_PORT` | `27017` | Puerto expuesto |

---

## RabbitMQ

| Variable | Default | Descripción |
|----------|---------|-------------|
| `RABBITMQ_USER` | `edugo` | Usuario |
| `RABBITMQ_PASSWORD` | `edugo123` | Contraseña |
| `RABBITMQ_PORT` | `5672` | Puerto AMQP |
| `RABBITMQ_MGMT_PORT` | `15672` | Puerto Management UI |

---

## JWT

| Variable | Default | Descripción |
|----------|---------|-------------|
| `JWT_SECRET` | `dev-secret-key...` | Secret para firma de tokens |

---

## APIs

| Variable | Default | Descripción |
|----------|---------|-------------|
| `API_MOBILE_PORT` | `8081` | Puerto API Mobile |
| `API_ADMIN_PORT` | `8082` | Puerto API Admin |

---

## Versiones Docker

| Variable | Default | Descripción |
|----------|---------|-------------|
| `API_MOBILE_VERSION` | `latest` | Tag de imagen (latest, develop, SHA) |
| `API_ADMIN_VERSION` | `latest` | Tag de imagen |
| `WORKER_VERSION` | `latest` | Tag de imagen |

---

## OpenAI (Worker)

| Variable | Default | Descripción |
|----------|---------|-------------|
| `OPENAI_API_KEY` | `sk-...` | API Key de OpenAI |

**⚠️ Obligatorio** para que el worker funcione.

Obtener en: https://platform.openai.com/api-keys
