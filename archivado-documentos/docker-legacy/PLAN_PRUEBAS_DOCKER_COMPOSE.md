# Plan de Trabajo: Correcciones Docker Compose y Pruebas

## Contexto

Con la eliminacion del entrypoint en api-mobile (PR #86, release v0.14.0), es necesario:
1. Corregir todos los docker-compose para eliminar referencias al entrypoint obsoleto
2. Probar cada docker-compose para validar que funciona correctamente

---

## Fase 1: Correcciones a Docker Compose

### 1.1 docker-compose.yml (Principal)

**Archivo:** `docker/docker-compose.yml`

**Cambios:**
- Eliminar variables del entrypoint script en api-mobile:
  ```yaml
  # ELIMINAR estas lineas del servicio api-mobile:
  - POSTGRES_HOST=postgres
  - POSTGRES_PORT=5432
  - MONGO_HOST=mongodb
  - MONGO_PORT=27017
  - RABBIT_HOST=rabbitmq
  - RABBIT_PORT=5672
  ```

**Razon:** Estas variables eran usadas por `docker-entrypoint.sh` para wait-for, que ya no existe.

---

### 1.2 docker-compose-mock.yml

**Archivo:** `docker/docker-compose-mock.yml`

**Cambios:**
- Eliminar override de entrypoint en api-mobile:
  ```yaml
  # ELIMINAR estas lineas del servicio api-mobile:
  entrypoint: ["/root/main"]
  working_dir: /root
  ```

**Razon:** La nueva imagen (v0.14.0) ya no tiene entrypoint, el CMD es `./main` directamente.

**NOTA:** Este cambio solo debe aplicarse DESPUES de que la imagen v0.14.0 este publicada.

---

### 1.3 docker-compose-apps.yml

**Archivo:** `docker/docker-compose-apps.yml`

**Cambios:**
- Eliminar variables del entrypoint script en api-mobile:
  ```yaml
  # ELIMINAR estas lineas del servicio api-mobile:
  - POSTGRES_HOST=${POSTGRES_HOST:-postgres}
  - POSTGRES_PORT=${POSTGRES_PORT:-5432}
  - MONGO_HOST=${MONGO_HOST:-mongodb}
  - MONGO_PORT=${MONGO_PORT:-27017}
  - RABBIT_HOST=${RABBITMQ_HOST:-rabbitmq}
  - RABBIT_PORT=${RABBITMQ_PORT:-5672}
  ```

**Razon:** Igual que docker-compose.yml, estas variables eran para el wait-for eliminado.

---

### 1.4 Archivos que NO requieren cambios

- `docker-compose-infrastructure.yml` - Solo infraestructura, sin APIs
- `docker-compose.migrate.yml` - Solo migrator + bases de datos
- `migrator/docker-compose.migrator.yml` - Template/propuesta

---

## Fase 2: Plan de Pruebas Detallado

### Pre-requisitos

1. Imagen api-mobile v0.14.0 publicada en GHCR
2. Docker Desktop corriendo
3. Puertos disponibles: 5432, 27017, 5672, 15672, 8081, 8082

### Limpieza Pre-Pruebas

```bash
# Ejecutar antes de cada prueba
docker-compose down -v --remove-orphans 2>/dev/null
docker-compose -f docker-compose-mock.yml down -v --remove-orphans 2>/dev/null
docker-compose -f docker-compose-infrastructure.yml down -v --remove-orphans 2>/dev/null
docker-compose -f docker-compose-apps.yml down -v --remove-orphans 2>/dev/null

# Limpiar imagenes locales (opcional, para prueba limpia)
docker system prune -f
```

---

### Prueba 1: docker-compose.yml (Principal)

**Objetivo:** Validar stack completo con infraestructura + APIs + migrator

**Pasos:**

```bash
cd docker/

# 1. Levantar stack completo
docker-compose up -d

# 2. Esperar que todos los servicios esten healthy (max 2 min)
docker-compose ps

# 3. Verificar logs del migrator (debe completar sin errores)
docker-compose logs migrator

# 4. Probar health de APIs
curl -s http://localhost:8081/health | jq .  # api-mobile
curl -s http://localhost:8082/health | jq .  # api-admin

# 5. Probar autenticacion
curl -s -X POST http://localhost:8082/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@edugo.test", "password": "edugo2024"}' | jq .

# 6. Probar validacion cruzada de JWT
TOKEN=$(curl -s -X POST http://localhost:8082/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@edugo.test", "password": "edugo2024"}' | jq -r '.data.token')

curl -s http://localhost:8081/api/v1/courses \
  -H "Authorization: Bearer $TOKEN" | jq .

# 7. Verificar Swagger
curl -s http://localhost:8081/swagger/index.html -o /dev/null -w "%{http_code}"  # Debe ser 200
curl -s http://localhost:8082/swagger/index.html -o /dev/null -w "%{http_code}"  # Debe ser 200
```

**Resultado Esperado:**
- [ ] Todos los contenedores en estado "healthy"
- [ ] Migrator termina exitosamente
- [ ] Health endpoints responden OK
- [ ] Login retorna token JWT
- [ ] Token de api-admin funciona en api-mobile
- [ ] Swagger accesible en ambas APIs

**Limpiar:**
```bash
docker-compose down -v
```

---

### Prueba 2: docker-compose-mock.yml

**Objetivo:** Validar APIs en modo mock (sin infraestructura)

**Pasos:**

```bash
cd docker/

# 1. Levantar solo APIs en modo mock
docker-compose -f docker-compose-mock.yml up -d

# 2. Verificar que NO hay contenedores de infraestructura
docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "(postgres|mongodb|rabbitmq)"
# (No debe mostrar nada)

# 3. Esperar startup (~5 segundos)
sleep 5

# 4. Verificar health
curl -s http://localhost:8081/health | jq .
curl -s http://localhost:8082/health | jq .

# 5. Probar login con usuario mock
curl -s -X POST http://localhost:8082/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@edugo.test", "password": "edugo2024"}' | jq .

# 6. Probar endpoints con datos mock
TOKEN=$(curl -s -X POST http://localhost:8082/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@edugo.test", "password": "edugo2024"}' | jq -r '.data.token')

# Listar instituciones (datos mock)
curl -s http://localhost:8082/v1/institutions \
  -H "Authorization: Bearer $TOKEN" | jq .

# Listar cursos desde api-mobile
curl -s http://localhost:8081/api/v1/courses \
  -H "Authorization: Bearer $TOKEN" | jq .

# 7. Verificar memoria (debe ser menor que stack completo)
docker stats --no-stream
```

**Resultado Esperado:**
- [ ] Solo contenedores api-mobile-mock y api-admin-mock corriendo
- [ ] Sin contenedores de infraestructura (postgres, mongodb, rabbitmq)
- [ ] Health OK en ambas APIs
- [ ] Login funciona con usuario mock
- [ ] Datos mock retornados en endpoints
- [ ] Memoria total < 500MB (vs ~1.5GB del stack completo)

**Limpiar:**
```bash
docker-compose -f docker-compose-mock.yml down
```

---

### Prueba 3: docker-compose-infrastructure.yml + docker-compose-apps.yml

**Objetivo:** Validar modo separado (infra primero, luego apps)

**Pasos:**

```bash
cd docker/

# 1. Levantar solo infraestructura
docker-compose -f docker-compose-infrastructure.yml up -d

# 2. Esperar que esten healthy
docker-compose -f docker-compose-infrastructure.yml ps

# 3. Verificar conectividad a servicios
docker exec edugo-postgres pg_isready -U edugo
docker exec edugo-mongodb mongosh --eval "db.adminCommand('ping')"
docker exec edugo-rabbitmq rabbitmq-diagnostics ping

# 4. Levantar aplicaciones (conectan a infra existente)
docker-compose -f docker-compose-apps.yml up -d

# 5. Verificar que api-mobile se levanto
docker-compose -f docker-compose-apps.yml ps

# 6. Probar health
curl -s http://localhost:8081/health | jq .
```

**Resultado Esperado:**
- [ ] Infraestructura levanta primero sin problemas
- [ ] Apps conectan a infraestructura existente
- [ ] api-mobile healthy
- [ ] Network edugo-network compartida

**Limpiar:**
```bash
docker-compose -f docker-compose-apps.yml down
docker-compose -f docker-compose-infrastructure.yml down -v
```

---

### Prueba 4: docker-compose.migrate.yml

**Objetivo:** Validar migracion forzada de bases de datos

**Pasos:**

```bash
cd docker/

# 1. Levantar con migracion forzada (ELIMINA datos existentes)
docker-compose -f docker-compose.migrate.yml up

# 2. Verificar logs del migrator
docker-compose -f docker-compose.migrate.yml logs migrator

# 3. El migrator debe terminar con exit 0
docker-compose -f docker-compose.migrate.yml ps migrator
```

**Resultado Esperado:**
- [ ] PostgreSQL y MongoDB levantan
- [ ] Migrator ejecuta con FORCE_MIGRATION=true
- [ ] Schemas recreados desde cero
- [ ] Migrator termina con exit 0

**Limpiar:**
```bash
docker-compose -f docker-compose.migrate.yml down -v
```

---

## Fase 3: Matriz de Compatibilidad

| Docker Compose | Usa Imagen | Requiere Infra | Modo Mock | Migrator |
|----------------|------------|----------------|-----------|----------|
| docker-compose.yml | GHCR latest | Si | No | Si |
| docker-compose-mock.yml | GHCR latest | No | Si | No |
| docker-compose-apps.yml | Build local | Externa | No | No |
| docker-compose-infrastructure.yml | N/A | Es infra | N/A | No |
| docker-compose.migrate.yml | Build local | Si (interno) | No | Si (forzado) |

---

## Fase 4: Checklist Final

### Correcciones Aplicadas
- [ ] docker-compose.yml - Variables entrypoint eliminadas
- [ ] docker-compose-mock.yml - Override entrypoint eliminado
- [ ] docker-compose-apps.yml - Variables entrypoint eliminadas

### Pruebas Exitosas
- [ ] Prueba 1: docker-compose.yml (stack completo)
- [ ] Prueba 2: docker-compose-mock.yml (modo mock)
- [ ] Prueba 3: Infraestructura + Apps separados
- [ ] Prueba 4: Migracion forzada

### Documentacion Actualizada
- [ ] README.md actualizado
- [ ] QUICK_START.md actualizado
- [ ] Comentarios en docker-compose actualizados

---

## Notas Importantes

1. **Orden de Pruebas:** Ejecutar en orden (1-4) para evitar conflictos de puertos
2. **Limpieza:** Siempre ejecutar `down -v` entre pruebas
3. **Imagen v0.14.0:** Las correcciones al docker-compose-mock.yml solo funcionaran cuando la imagen este publicada
4. **Tiempos:** El stack completo toma ~2 min en levantar, mock ~5 segundos
