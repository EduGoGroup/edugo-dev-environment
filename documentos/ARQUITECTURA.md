# Arquitectura - EduGo Dev Environment

## Visión General

EduGo es una plataforma educativa compuesta por múltiples microservicios. Este repositorio proporciona la configuración Docker para ejecutar todo el ecosistema localmente.

## Diagrama de Arquitectura

```
                    ┌─────────────────────────────────────┐
                    │      Aplicaciones Frontend          │
                    │   (React, Vue, Angular, Mobile)     │
                    └─────────────────┬───────────────────┘
                                      │ HTTP REST
                                      ▼
    ┌─────────────────────────────────────────────────────────────┐
    │                    Docker Compose                            │
    │                                                              │
    │  ┌─────────────────┐     ┌─────────────────┐                │
    │  │   API Mobile    │     │  API Administr. │                │
    │  │   :8081         │     │   :8082         │                │
    │  └────────┬────────┘     └────────┬────────┘                │
    │           │                       │                          │
    │           └───────────┬───────────┘                          │
    │                       │                                      │
    │           ┌───────────┼───────────┐                          │
    │           ▼           ▼           ▼                          │
    │  ┌─────────────┐ ┌─────────┐ ┌─────────┐                    │
    │  │ PostgreSQL  │ │ MongoDB │ │RabbitMQ │                    │
    │  │   :5432     │ │  :27017 │ │  :5672  │                    │
    │  └─────────────┘ └─────────┘ └────┬────┘                    │
    │                                    │                         │
    │                           ┌────────▼────────┐                │
    │                           │     Worker      │                │
    │                           │ (Procesa PDFs)  │                │
    │                           └─────────────────┘                │
    │                                                              │
    │  ┌─────────────────┐                                        │
    │  │    Migrator     │  (Ejecuta una vez al inicio)           │
    │  └─────────────────┘                                        │
    └─────────────────────────────────────────────────────────────┘
```

## Componentes

### APIs

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| **API Mobile** | 8081 | API principal para aplicaciones móviles y web |
| **API Administración** | 8082 | API para panel de administración |

### Bases de Datos

| Servicio | Puerto | Uso |
|----------|--------|-----|
| **PostgreSQL** | 5432 | Datos relacionales (usuarios, cursos, instituciones) |
| **MongoDB** | 27017 | Documentos no estructurados (materiales, PDFs procesados) |

### Mensajería

| Servicio | Puerto | Uso |
|----------|--------|-----|
| **RabbitMQ** | 5672 | Cola de mensajes para procesamiento asíncrono |
| **RabbitMQ UI** | 15672 | Interfaz de administración de colas |

### Servicios de Soporte

| Servicio | Descripción |
|----------|-------------|
| **Worker** | Procesa PDFs y genera contenido con IA (OpenAI) |
| **Migrator** | Ejecuta migraciones de base de datos al inicio |

## Flujo de Datos

### 1. Autenticación
```
Frontend → API Mobile (/auth/login) → PostgreSQL → JWT Token
```

### 2. Consulta de Datos
```
Frontend → API Mobile (/courses) → PostgreSQL → Response JSON
```

### 3. Subida de PDF
```
Frontend → API Mobile → RabbitMQ (cola) → Worker → OpenAI → MongoDB
```

### 4. Panel de Admin
```
Admin Panel → API Administración → PostgreSQL → Response JSON
```

## Redes Docker

Todos los servicios están conectados a la red `edugo-network` que permite comunicación interna usando nombres de servicio como hostnames:

- `postgres` - Host para PostgreSQL
- `mongodb` - Host para MongoDB
- `rabbitmq` - Host para RabbitMQ

## Volúmenes

Los datos persisten en volúmenes Docker:

| Volumen | Propósito |
|---------|-----------|
| `postgres-data` | Datos de PostgreSQL |
| `mongodb-data` | Datos de MongoDB |
| `rabbitmq-data` | Configuración y colas de RabbitMQ |

## Health Checks

Cada servicio de infraestructura tiene health checks configurados:

- **PostgreSQL**: `pg_isready` cada 10 segundos
- **MongoDB**: `mongosh ping` cada 10 segundos
- **RabbitMQ**: `rabbitmq-diagnostics ping` cada 10 segundos

Las APIs esperan que la infraestructura esté saludable antes de iniciar.

## Dependencias del Ecosistema EduGo

```
edugo-dev-environment (este repo)
    │
    ├── edugo-api-mobile          → Imagen Docker desde ghcr.io
    ├── edugo-api-administracion  → Imagen Docker desde ghcr.io
    ├── edugo-worker              → Imagen Docker desde ghcr.io
    │
    └── edugo-infrastructure      → Migraciones (vía migrator)
        ├── postgres/migrations
        └── mongodb/migrations
```

---

**Ver también:** [SERVICIOS.md](./SERVICIOS.md) para configuración detallada de cada servicio.
