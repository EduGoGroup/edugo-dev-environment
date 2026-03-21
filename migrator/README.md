# Migrator - EduGo

Servicio que aplica migraciones de esquema y datos iniciales para PostgreSQL y MongoDB al iniciar el entorno de desarrollo.

## Estado actual

| Base de datos | Estado |
|---|---|
| PostgreSQL | Funcional al 100% — migraciones, seeds y versionamiento |
| MongoDB | Conecta correctamente, aplica migraciones y seeds |

## Dependencias

```
github.com/EduGoGroup/edugo-infrastructure/postgres v0.65.0
github.com/EduGoGroup/edugo-infrastructure/mongodb  v0.55.0
```

## Variables de entorno

### PostgreSQL

| Variable | Default | Descripcion |
|---|---|---|
| `POSTGRES_URI` | — | URI completa (alternativa a las variables individuales) |
| `POSTGRES_HOST` | `localhost` | Host |
| `POSTGRES_PORT` | `5432` | Puerto |
| `POSTGRES_DB` | `edugo` | Base de datos |
| `POSTGRES_USER` | `edugo` | Usuario |
| `POSTGRES_PASSWORD` | `edugo123` | Contrasena |
| `POSTGRES_SSLMODE` | `disable` | Modo SSL |

### MongoDB

| Variable | Default | Descripcion |
|---|---|---|
| `MONGO_URI` | — | URI completa (alternativa a las variables individuales) |
| `MONGO_HOST` | `localhost` | Host |
| `MONGO_PORT` | `27017` | Puerto |
| `MONGO_USER` | `edugo` | Usuario |
| `MONGO_PASSWORD` | `edugo123` | Contrasena |
| `MONGO_DB_NAME` | `edugo` | Base de datos |

### Flags de control

| Variable | Default | Descripcion |
|---|---|---|
| `FORCE_MIGRATION` | `false` | Elimina y recrea todas las bases de datos |
| `APPLY_MOCK_DATA` | `true` | Aplica datos de desarrollo |
| `POSTGRES_ONLY` | `false` | Ejecuta solo migraciones de PostgreSQL |
| `MONGO_ONLY` | `false` | Ejecuta solo migraciones de MongoDB |
| `STATUS_ONLY` | `false` | Muestra estado actual sin aplicar cambios |

## Uso con docker-compose

El migrator esta integrado en el `docker-compose.yml` principal del entorno de desarrollo:

```bash
# Levantar solo infraestructura (postgres, mongodb, rabbitmq + migrator)
docker compose up

# Levantar entorno completo (infraestructura + todas las apps)
docker compose --profile full up

# Ver logs del migrator
docker compose logs migrator
```

El migrator corre una sola vez y termina. Si las bases de datos ya tienen datos, las migraciones se omiten (comportamiento idempotente). Para forzar una recreacion completa:

```bash
FORCE_MIGRATION=true docker compose up
```

## Tests de integracion

Los tests usan testcontainers para crear contenedores temporales de PostgreSQL y MongoDB. Docker debe estar corriendo.

```bash
# Ejecutar todos los tests
go test -v ./tests

# Solo PostgreSQL
go test -v ./tests -run TestPostgresIntegration

# Solo MongoDB
go test -v ./tests -run TestMongoDBIntegration
```

## Compilar y ejecutar local

```bash
# Descargar dependencias
go mod download

# Compilar
go build -o bin/migrator ./cmd

# Ejecutar (requiere PostgreSQL y MongoDB corriendo)
./bin/migrator
```

## Referencias

- [docs/DOCKER_INTEGRATION_GUIDE.md](docs/DOCKER_INTEGRATION_GUIDE.md) — guia detallada de integracion con docker-compose
- [docs/MONGODB_LIMITATIONS.md](docs/MONGODB_LIMITATIONS.md) — limitaciones conocidas de MongoDB y soluciones propuestas
