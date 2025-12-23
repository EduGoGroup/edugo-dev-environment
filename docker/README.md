# EduGo - Docker Compose

Este directorio contiene la configuración de Docker Compose para el ambiente de desarrollo.

## Archivos

| Archivo | Propósito |
|---------|-----------|
| `docker-compose.yml` | Archivo principal consolidado con profiles |
| `docker-compose.migrate.yml` | Migrator para actualización de base de datos |
| `.env.example` | Plantilla de variables de entorno |
| `.env` | Variables de entorno (no versionado) |

## Profiles Disponibles

El archivo `docker-compose.yml` usa profiles para flexibilidad:

| Profile | Servicios |
|---------|-----------|
| (ninguno) | Solo infraestructura: postgres, mongodb, rabbitmq |
| `apps` | Infraestructura + API Mobile |
| `admin` | Infraestructura + API Admin |
| `worker` | Infraestructura + Worker |
| `full` | Todo: infraestructura + todas las apps + redis |
| `with-redis` | Agrega Redis |

## Uso Rápido

```bash
# Solo infraestructura (desarrollo de APIs localmente)
docker-compose up -d

# Infraestructura + API Mobile
docker-compose --profile apps up -d

# Todo el ecosistema
docker-compose --profile full up -d

# Detener
docker-compose down

# Detener y eliminar datos
docker-compose down -v
```

## Servicios y Puertos

### Infraestructura (siempre disponible)

| Servicio | Puerto | Credenciales |
|----------|--------|--------------|
| PostgreSQL | 5432 | edugo / edugo123 |
| MongoDB | 27017 | edugo / edugo123 |
| RabbitMQ AMQP | 5672 | edugo / edugo123 |
| RabbitMQ UI | 15672 | edugo / edugo123 |

### Aplicaciones (requieren profile)

| Servicio | Puerto | Profile |
|----------|--------|---------|
| API Mobile | 8081 | apps, full |
| API Admin | 8082 | admin, full |
| Worker | - | worker, full |
| Redis | 6379 | with-redis, full |

## URLs Importantes

- **API Mobile Swagger:** http://localhost:8081/swagger/index.html
- **API Admin Swagger:** http://localhost:8082/swagger/index.html
- **RabbitMQ Management:** http://localhost:15672

## Ejecutar Migraciones

```bash
# Usando el migrator
docker-compose -f docker-compose.migrate.yml up --build

# O con make desde la raíz
make migrate
```

## Variables de Entorno

Copiar `.env.example` a `.env` y ajustar:

```bash
cp .env.example .env
```

Ver documentación completa en `documentos/VARIABLES-ENTORNO.md`.

## Archivos Archivados

Los siguientes archivos fueron consolidados y movidos a `archivado-documentos/docker-legacy/`:
- `docker-compose-apps.yml`
- `docker-compose-infrastructure.yml`
- `docker-compose-mock.yml`
- Documentos temporales de validación

---

**Ver también:** `documentos/GUIA-RAPIDA.md` para instrucciones completas de setup.
