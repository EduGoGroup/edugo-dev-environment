# EduGo Migrator

Microproyecto en Go para ejecutar migraciones de base de datos automáticamente, utilizando el repositorio [edugo-infrastructure](https://github.com/EduGoGroup/edugo-infrastructure).

## 🎯 Propósito

Este migrator automatiza la ejecución de migraciones de PostgreSQL y MongoDB, sincronizándose automáticamente con los últimos scripts del repositorio de infraestructura mediante `git clone/pull`.

## 🚀 Uso

### Ejecución Manual

```bash
cd migrator
go run cmd/main.go
```

### Variables de Entorno

El migrator utiliza las siguientes variables de entorno (con valores por defecto para desarrollo local):

**PostgreSQL:**
- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_NAME` (default: `edugo`)
- `DB_USER` (default: `edugo`)
- `DB_PASSWORD` (default: `edugo123`)

**MongoDB:**
- `MONGO_HOST` (default: `localhost`)
- `MONGO_PORT` (default: `27017`)
- `MONGO_USER` (default: `edugo`)
- `MONGO_PASSWORD` (default: `edugo123`)
- `MONGO_DB_NAME` (default: `edugo`)

## 📋 Funcionamiento

1. **Sincronización**: Clona o actualiza el repositorio `edugo-infrastructure` en `.infrastructure/`
2. **PostgreSQL**: Ejecuta migraciones pendientes usando `postgres/migrate.go`
3. **MongoDB**: Ejecuta migraciones pendientes usando `mongodb/migrate.go`

## 🔧 Integración con Docker Compose

El migrator puede ejecutarse como un servicio en docker-compose para aplicar migraciones automáticamente al iniciar el stack:

```yaml
services:
  migrator:
    build:
      context: ../migrator
      dockerfile: Dockerfile
    image: edugogroup-migrator:latest
    container_name: edugo-migrator
    depends_on:
      postgres:
        condition: service_healthy
      mongodb:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_NAME=edugo
      - DB_USER=edugo
      - DB_PASSWORD=${POSTGRES_PASSWORD:-edugo123}
      - MONGO_HOST=mongodb
      - MONGO_USER=edugo
      - MONGO_PASSWORD=${MONGO_PASSWORD:-edugo123}
    profiles:
      - infrastructure
```

## 📝 Notas

- El migrator siempre obtiene la última versión de los scripts de migración
- Las migraciones ya aplicadas son detectadas y omitidas automáticamente
- Si una migración falla en PostgreSQL, el proceso continúa con MongoDB
- El directorio `.infrastructure/` se crea automáticamente y debe añadirse a `.gitignore`

## 🐛 Troubleshooting

### Error: "database does not exist"
Asegúrate de que las bases de datos PostgreSQL y MongoDB estén creadas y accesibles.

### Error: "password authentication failed"
Verifica que las credenciales en las variables de entorno coincidan con las configuradas en PostgreSQL/MongoDB.

### Error: "mongosh: executable file not found"
Las migraciones de MongoDB requieren `mongosh` instalado. En Docker, esto se maneja automáticamente.
