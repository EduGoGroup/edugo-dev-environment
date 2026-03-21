# Guía de Integración con Docker Compose

## 📋 Pasos para Integrar el Migrator en Docker Compose

### Paso 1: Agregar el Servicio al docker-compose.yml

Abre el archivo `docker/docker-compose.yml` y agrega el siguiente servicio después de los servicios de base de datos:

```yaml
  # ========================================
  # MIGRATOR - Ejecuta migraciones automáticas
  # ========================================
  migrator:
    build:
      context: ../migrator
      dockerfile: Dockerfile
    image: edugogroup-migrator:latest
    container_name: edugo-migrator
    profiles:
      - full
      - db-only
    depends_on:
      postgres:
        condition: service_healthy
      mongodb:
        condition: service_healthy
    environment:
      # PostgreSQL
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=${POSTGRES_DB:-edugo}
      - DB_USER=${POSTGRES_USER:-edugo}
      - DB_PASSWORD=${POSTGRES_PASSWORD:-edugo123}
      
      # MongoDB
      - MONGO_HOST=mongodb
      - MONGO_PORT=27017
      - MONGO_USER=${MONGO_USER:-edugo}
      - MONGO_PASSWORD=${MONGO_PASSWORD:-edugo123}
      - MONGO_DB_NAME=${MONGO_DB:-edugo}
    networks:
      - edugo-network
    restart: "no"
```

### Paso 2: Actualizar .gitignore

Asegúrate de que `.infrastructure/` esté en el `.gitignore` del repositorio:

```bash
echo "migrator/.infrastructure/" >> .gitignore
```

### Paso 3: Construir la Imagen

```bash
cd docker
docker compose build migrator
```

### Paso 4: Ejecutar las Migraciones

**Opción A - Con todo el stack:**
```bash
docker compose --profile full up -d
```

**Opción B - Solo infraestructura:**
```bash
docker compose --profile db-only up -d
```

**Opción C - Solo migrator:**
```bash
docker compose up migrator
```

### Paso 5: Verificar Ejecución

```bash
# Ver logs del migrator
docker compose logs migrator

# Ver estado
docker compose ps migrator
```

## 🔄 Flujo de Trabajo

### Iniciar el Sistema Completo
```bash
cd docker
docker compose --profile full up -d
```

**Orden de ejecución:**
1. PostgreSQL inicia y pasa healthcheck ✅
2. MongoDB inicia y pasa healthcheck ✅
3. RabbitMQ inicia y pasa healthcheck ✅
4. **Migrator ejecuta migraciones** ✅
5. API Mobile inicia ✅
6. API Admin inicia ✅
7. Worker inicia ✅

### Re-ejecutar Migraciones

Si necesitas volver a ejecutar las migraciones:

```bash
# Opción 1: Restart del contenedor
docker compose restart migrator

# Opción 2: Recrear el contenedor
docker compose up --force-recreate migrator

# Opción 3: Ejecutar manualmente
docker compose run --rm migrator
```

### Actualizar Migraciones del Repositorio

El migrator automáticamente hace `git pull` cada vez que se ejecuta, por lo que siempre usa la última versión de las migraciones.

## 🐛 Troubleshooting

### El migrator no se ejecuta automáticamente

**Causa**: No está en el profile activo

**Solución**:
```bash
# Usar profile que incluya migrator
docker compose --profile full up -d
# o
docker compose --profile db-only up -d
```

### Migraciones fallan con "database does not exist"

**Causa**: PostgreSQL/MongoDB no están listos

**Solución**: El `depends_on` con `condition: service_healthy` debería manejarlo. Si persiste:
```bash
# Verificar healthcheck
docker compose ps

# Esperar manualmente y re-ejecutar
sleep 10
docker compose up migrator
```

### Error "permission denied" en git

**Causa**: El repositorio de infraestructura requiere autenticación

**Solución**: El repositorio es público, no debería ocurrir. Si ocurre, verificar conectividad de red del contenedor.

### Migraciones no se actualizan

**Causa**: El directorio `.infrastructure/` está cacheado

**Solución**:
```bash
# Eliminar el volumen/caché y recrear
docker compose down migrator
docker volume prune  # cuidado: elimina volúmenes no usados
docker compose up migrator
```

## ⚙️ Configuración Avanzada

### Ejecutar solo migraciones de PostgreSQL

Modifica temporalmente `cmd/main.go` comentando la sección de MongoDB.

### Usar una versión específica de infraestructura

Modifica `cmd/main.go` para hacer checkout de un tag específico:

```go
cmd := exec.Command("git", "checkout", "v1.2.3")
cmd.Dir = infraDir
cmd.Run()
```

### Agregar logs más detallados

Modifica las funciones de migración para incluir más output.

## 📊 Monitoreo

### Ver progreso en tiempo real

```bash
docker compose logs -f migrator
```

### Verificar tablas creadas

**PostgreSQL:**
```bash
docker compose exec postgres psql -U edugo -d edugo -c "\dt"
```

**MongoDB:**
```bash
docker compose exec mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo --eval "show collections"
```

## 🎯 Recomendaciones

1. **Siempre revisar logs**: `docker compose logs migrator` después de levantar el stack
2. **No modificar `.infrastructure/`**: Se sobreescribe en cada ejecución
3. **Probar migraciones localmente**: Usar `go run cmd/main.go` antes de integrar con Docker
4. **Mantener credenciales sincronizadas**: Las variables de entorno deben coincidir con las bases de datos

## 🚀 Next Steps

Una vez integrado y funcionando:

1. ✅ Las migraciones se ejecutan automáticamente en cada `docker compose up`
2. ✅ Siempre usan la última versión de los scripts
3. ✅ No requiere intervención manual
4. ✅ Logs disponibles para debugging
