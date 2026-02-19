# EduGo - Configuraci贸n Cloud

> Infraestructura de EduGo en servicios cloud (Neon, Atlas, Upstash)

## Servicios Soportados

| Servicio | Proveedor | Prop贸sito |
|----------|-----------|-----------|
| PostgreSQL | [Neon](https://neon.tech) | Base de datos relacional |
| MongoDB | [MongoDB Atlas](https://www.mongodb.com/atlas) | Base de datos de documentos |
| Redis | [Upstash](https://upstash.com) | Cache y sesiones |
| RabbitMQ | Docker local | Mensajer铆a (opcional) |

## Configuraci贸n

### 1. Crear archivo de variables

```bash
cp docker/.env.cloud.example docker/.env.cloud
```

Edita `docker/.env.cloud` con las credenciales de tus servicios cloud.

### 2. Ejecutar migraciones

```bash
# Migraciones idempotentes (si ya existen tablas/colecciones, no hace nada)
make db-migrate-cloud

# Recrear desde cero (DESTRUYE DATOS)
make db-recreate-cloud
```

### 3. Levantar APIs con servicios cloud

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --env-file .env.cloud up -d --profile full
```

## Variables de Entorno

### PostgreSQL (Neon)

Se puede configurar de dos formas:

**Opci贸n A: URI directa** (recomendada para Neon)
```bash
POSTGRES_URI=postgresql://user:password@host:5432/edugo?sslmode=require
```

**Opci贸n B: Variables individuales**
```bash
POSTGRES_HOST=<tu-host>.neon.tech
POSTGRES_PORT=5432
POSTGRES_USER=<tu-usuario>
POSTGRES_PASSWORD=<tu-password>
POSTGRES_DB=edugo
POSTGRES_SSLMODE=require
```

### MongoDB (Atlas)

```bash
MONGO_URI=mongodb+srv://<usuario>:<password>@<cluster>.mongodb.net/?appName=Edugo
MONGO_DB_NAME=edugo
```

### Redis (Upstash)

```bash
REDIS_URL=redis://default:<password>@<host>.upstash.io:6379
```

## Comandos Disponibles

| Comando | Descripci贸n |
|---------|-------------|
| `make db-migrate-cloud` | Ejecutar migraciones (idempotente) |
| `make db-recreate-cloud` | Recrear BD desde cero (DESTRUYE DATOS) |
| `make db-migrate` | Migraciones en Docker local |
| `make db-recreate` | Recrear BD en Docker local |

## L铆mites de Planes Gratuitos

| Servicio | Almacenamiento | Notas |
|----------|---------------|-------|
| Neon | 0.5 GB | 100 CU-horas/mes |
| Atlas | 512 MB | Plan M0, 500 conexiones |
| Upstash | 256 MB | 500,000 comandos/mes |

## Troubleshooting

### Error de conexi贸n SSL (PostgreSQL)
Aseg煤rate de incluir `sslmode=require` en tu URI o variable `POSTGRES_SSLMODE`.

### L铆mite de almacenamiento
Limpia datos de prueba con `APPLY_MOCK_DATA=false` al recrear:
```bash
APPLY_MOCK_DATA=false make db-recreate-cloud
```

## Referencias

- [Documentaci贸n de Neon](https://neon.tech/docs)
- [Documentaci贸n de Atlas](https://www.mongodb.com/docs/atlas/)
- [Documentaci贸n de Upstash](https://upstash.com/docs)
