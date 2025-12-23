# Variables de Entorno

Este documento describe todas las variables de entorno utilizadas en el proyecto EduGo Dev Environment.

## Resumen de Prefijos por Servicio

Cada servicio utiliza un prefijo diferente para sus variables de entorno debido a cómo Viper (la librería de configuración de Go) lee las variables con `AutomaticEnv`:

| Servicio | Prefijo | Ejemplo |
|----------|---------|---------|
| API Mobile | `DATABASE_*`, `MESSAGING_*`, `STORAGE_*` | `DATABASE_POSTGRES_HOST` |
| API Admin | `EDUGO_ADMIN_*` | `EDUGO_ADMIN_DATABASE_POSTGRES_HOST` |
| Worker | `EDUGO_WORKER_*` | `EDUGO_WORKER_DATABASE_POSTGRES_HOST` |

---

## Variables Base (Infraestructura)

Estas variables se definen en `.env` y son usadas por Docker Compose para la infraestructura:

### PostgreSQL

| Variable | Descripción | Default |
|----------|-------------|---------|
| `POSTGRES_HOST` | Host del servidor PostgreSQL | `postgres` |
| `POSTGRES_PORT` | Puerto de PostgreSQL | `5432` |
| `POSTGRES_USER` | Usuario de PostgreSQL | `edugo` |
| `POSTGRES_PASSWORD` | Contraseña de PostgreSQL | `edugo123` |
| `POSTGRES_DB` | Nombre de la base de datos | `edugo` |

### MongoDB

| Variable | Descripción | Default |
|----------|-------------|---------|
| `MONGO_HOST` | Host del servidor MongoDB | `mongodb` |
| `MONGO_PORT` | Puerto de MongoDB | `27017` |
| `MONGO_USER` | Usuario de MongoDB | `edugo` |
| `MONGO_PASSWORD` | Contraseña de MongoDB | `edugo123` |
| `MONGO_DB` | Nombre de la base de datos | `edugo` |

### RabbitMQ

| Variable | Descripción | Default |
|----------|-------------|---------|
| `RABBITMQ_HOST` | Host del servidor RabbitMQ | `rabbitmq` |
| `RABBITMQ_PORT` | Puerto AMQP de RabbitMQ | `5672` |
| `RABBITMQ_MGMT_PORT` | Puerto de management UI | `15672` |
| `RABBITMQ_USER` | Usuario de RabbitMQ | `edugo` |
| `RABBITMQ_PASSWORD` | Contraseña de RabbitMQ | `edugo123` |

---

## Mapeo de Variables por Servicio

### PostgreSQL

| Concepto | API Mobile | API Admin | Worker |
|----------|------------|-----------|--------|
| Host | `DATABASE_POSTGRES_HOST` | `EDUGO_ADMIN_DATABASE_POSTGRES_HOST` | `EDUGO_WORKER_DATABASE_POSTGRES_HOST` |
| Port | `DATABASE_POSTGRES_PORT` | `EDUGO_ADMIN_DATABASE_POSTGRES_PORT` | `EDUGO_WORKER_DATABASE_POSTGRES_PORT` |
| User | `DATABASE_POSTGRES_USER` | `EDUGO_ADMIN_DATABASE_POSTGRES_USER` | `EDUGO_WORKER_DATABASE_POSTGRES_USER` |
| Password | `DATABASE_POSTGRES_PASSWORD` | `POSTGRES_PASSWORD` | `POSTGRES_PASSWORD` |
| Database | `DATABASE_POSTGRES_DATABASE` | `EDUGO_ADMIN_DATABASE_POSTGRES_DATABASE` | `EDUGO_WORKER_DATABASE_POSTGRES_DATABASE` |
| SSL Mode | `DATABASE_POSTGRES_SSLMODE` | `EDUGO_ADMIN_DATABASE_POSTGRES_SSLMODE` | `EDUGO_WORKER_DATABASE_POSTGRES_SSLMODE` |
| Max Conn | `DATABASE_POSTGRES_MAX_CONNECTIONS` | `EDUGO_ADMIN_DATABASE_POSTGRES_MAX_CONNECTIONS` | `EDUGO_WORKER_DATABASE_POSTGRES_MAX_CONNECTIONS` |

### MongoDB

| Concepto | API Mobile | API Admin | Worker |
|----------|------------|-----------|--------|
| URI | `DATABASE_MONGODB_URI` | `MONGODB_URI` | `MONGODB_URI` |
| Database | `DATABASE_MONGODB_DATABASE` | `EDUGO_ADMIN_DATABASE_MONGODB_DATABASE` | `EDUGO_WORKER_DATABASE_MONGODB_DATABASE` |
| Timeout | `DATABASE_MONGODB_TIMEOUT` | `EDUGO_ADMIN_DATABASE_MONGODB_TIMEOUT` | `EDUGO_WORKER_DATABASE_MONGODB_TIMEOUT` |

### RabbitMQ

| Concepto | API Mobile | Worker |
|----------|------------|--------|
| URL | `MESSAGING_RABBITMQ_URL` | `RABBITMQ_URL` |
| Prefetch | `MESSAGING_RABBITMQ_PREFETCH_COUNT` | `EDUGO_WORKER_MESSAGING_RABBITMQ_PREFETCH_COUNT` |
| Queue Material | `MESSAGING_RABBITMQ_QUEUES_MATERIAL_UPLOADED` | `EDUGO_WORKER_MESSAGING_RABBITMQ_QUEUES_MATERIAL_UPLOADED` |
| Queue Assessment | `MESSAGING_RABBITMQ_QUEUES_ASSESSMENT_ATTEMPT` | `EDUGO_WORKER_MESSAGING_RABBITMQ_QUEUES_ASSESSMENT_ATTEMPT` |
| Exchange | `MESSAGING_RABBITMQ_EXCHANGES_MATERIALS` | `EDUGO_WORKER_MESSAGING_RABBITMQ_EXCHANGES_MATERIALS` |

---

## Variables Específicas por Servicio

### API Mobile

| Variable | Descripción | Default |
|----------|-------------|---------|
| `API_MOBILE_PORT` | Puerto expuesto del servicio | `8081` |
| `STORAGE_S3_ACCESS_KEY_ID` | Access Key de S3 | - |
| `STORAGE_S3_SECRET_ACCESS_KEY` | Secret Key de S3 | - |
| `STORAGE_S3_BUCKET_NAME` | Nombre del bucket S3 | `edugo-materials-dev-local` |
| `STORAGE_S3_REGION` | Región de S3 | `us-east-1` |
| `STORAGE_S3_ENDPOINT` | Endpoint custom (MinIO) | - |
| `BOOTSTRAP_OPTIONAL_RESOURCES_S3` | S3 como recurso opcional | `true` |
| `BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ` | RabbitMQ como recurso opcional | `false` |

### API Admin

| Variable | Descripción | Default |
|----------|-------------|---------|
| `API_ADMIN_PORT` | Puerto expuesto del servicio | `8082` |
| `EDUGO_ADMIN_SERVER_PORT` | Puerto interno | `8081` |
| `EDUGO_ADMIN_SERVER_HOST` | Host de binding | `0.0.0.0` |
| `EDUGO_ADMIN_SERVER_READ_TIMEOUT` | Timeout de lectura | `30s` |
| `EDUGO_ADMIN_SERVER_WRITE_TIMEOUT` | Timeout de escritura | `30s` |

### Worker

| Variable | Descripción | Default |
|----------|-------------|---------|
| `OPENAI_API_KEY` | API Key de OpenAI | (requerido) |
| `EDUGO_WORKER_NLP_PROVIDER` | Proveedor NLP | `openai` |
| `EDUGO_WORKER_NLP_MODEL` | Modelo a usar | `gpt-4` |
| `EDUGO_WORKER_NLP_MAX_TOKENS` | Max tokens | `2000` |
| `EDUGO_WORKER_NLP_TEMPERATURE` | Temperature | `0.7` |

---

## Variables Globales

| Variable | Descripción | Default |
|----------|-------------|---------|
| `APP_ENV` | Ambiente de ejecución | `development` |
| `ENV` | Alias de APP_ENV | `development` |
| `JWT_SECRET` | Secreto para tokens JWT | `dev-secret-key-change-in-production-edugo-2024` |
| `AUTH_JWT_SECRET` | Alias para API Mobile | (usa JWT_SECRET) |
| `LOG_LEVEL` | Nivel de logging | `debug` |
| `LOG_FORMAT` | Formato de logs | `json` |

---

## Construcción de URIs

### MongoDB URI

```
mongodb://${MONGO_USER}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/${MONGO_DB}?authSource=admin
```

Ejemplo:
```
mongodb://edugo:edugo123@mongodb:27017/edugo?authSource=admin
```

### RabbitMQ URL

```
amqp://${RABBITMQ_USER}:${RABBITMQ_PASSWORD}@${RABBITMQ_HOST}:${RABBITMQ_PORT}/
```

Ejemplo:
```
amqp://edugo:edugo123@rabbitmq:5672/
```

---

## Notas Importantes

1. **Prefijos diferentes por diseño**: Cada servicio usa prefijos distintos porque así está configurado Viper en cada aplicación.

2. **Variables duplicadas**: Algunas variables como `POSTGRES_PASSWORD` se pasan sin prefijo además de con prefijo para compatibilidad.

3. **Valores por defecto**: Docker Compose aplica valores por defecto usando la sintaxis `${VAR:-default}`.

4. **S3 Opcional**: La API Mobile puede funcionar sin S3 si `BOOTSTRAP_OPTIONAL_RESOURCES_S3=true`.

5. **OpenAI Requerido**: El Worker requiere `OPENAI_API_KEY` válido para funcionar.
