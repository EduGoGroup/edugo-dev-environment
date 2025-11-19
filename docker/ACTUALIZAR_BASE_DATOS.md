# 🔄 Actualización de Base de Datos

Este documento explica cómo actualizar el esquema de las bases de datos cuando hay cambios en el repositorio `edugo-infrastructure`.

---

## ⚠️ ADVERTENCIA IMPORTANTE

**Este proceso ELIMINA COMPLETAMENTE las bases de datos y las recrea desde cero.**

- ❌ **Se perderán TODOS los datos existentes**
- ✅ Solo usar en entorno de **DESARROLLO**
- ✅ Las bases de datos se recrearán con datos de prueba actualizados

---

## 📋 ¿Cuándo usar esto?

Usa este proceso cuando:

1. El equipo de backend actualiza el esquema de base de datos
2. Se agregan nuevas tablas o colecciones
3. Se modifican estructuras de datos existentes
4. Necesitas resetear a un estado limpio con datos de prueba frescos

---

## 🚀 Proceso de Actualización

### Opción 1: Usando docker-compose.migrate.yml (Recomendado)

Este método construye la imagen del migrator (obteniendo las últimas dependencias) y ejecuta la migración forzada.

```bash
# 1. Ir al directorio docker
cd docker

# 2. Ejecutar migración forzada
docker-compose -f docker-compose.migrate.yml up

# 3. Reiniciar los servicios principales
docker-compose up -d
```

**¿Qué hace este comando?**
1. ✅ Construye la imagen del migrator (`build` en vez de `image`)
2. ✅ Ejecuta `go mod download` para obtener últimas dependencias
3. ✅ Clona/actualiza el repositorio `edugo-infrastructure`
4. ✅ Elimina completamente PostgreSQL schema y MongoDB database
5. ✅ Recrea desde cero con la estructura más reciente
6. ✅ Carga datos de prueba actualizados

---

### Opción 2: Variable de entorno manual

Si prefieres más control, puedes ejecutar el migrator directamente:

```bash
# 1. Reconstruir imagen del migrator
cd docker
docker-compose build migrator

# 2. Ejecutar con FORCE_MIGRATION=true
docker-compose run --rm -e FORCE_MIGRATION=true migrator

# 3. Reiniciar servicios
docker-compose up -d
```

---

## 📊 Salida Esperada

Cuando ejecutes la migración forzada, deberías ver:

```
=== EduGo Migrator ===
Iniciando proceso de migraciones...
⚠️  MODO FORZADO ACTIVADO - Se eliminarán y recrearán todas las bases de datos

📦 Obteniendo repositorio de infraestructura...
✅ Repositorio de infraestructura listo

--- PostgreSQL Migrations ---
🔥 Eliminando schema público de PostgreSQL...
✅ Schema eliminado exitosamente
Ejecutando runner de PostgreSQL (estructura, constraints, seeds y testing)...
✓ Conectado a PostgreSQL: edugo@postgres:5432/edugo

═══════════════════════════════════════════════════════════════
  CAPA: STRUCTURE
═══════════════════════════════════════════════════════════════
▸ Ejecutando: 001_create_users.sql
  ✓ Éxito
...
✅ Migraciones de PostgreSQL completadas

--- MongoDB Migrations ---
🔥 Eliminando base de datos MongoDB...
✅ Base de datos MongoDB eliminada exitosamente
Ejecutando runner de MongoDB (estructura y constraints)...
🏗️  Ejecutando Structure...
✅ 001_material_assessment
...
✅ Migraciones de MongoDB completadas

✅ Todas las migraciones se ejecutaron correctamente
```

---

## ✅ Verificación

Después de la migración, verifica que todo esté correcto:

### PostgreSQL
```bash
# Ver tablas
docker exec edugo-postgres psql -U edugo -d edugo -c "\dt"

# Ver cantidad de datos
docker exec edugo-postgres psql -U edugo -d edugo -c "
SELECT 'users' as tabla, COUNT(*) FROM users
UNION ALL SELECT 'schools', COUNT(*) FROM schools
UNION ALL SELECT 'materials', COUNT(*) FROM materials;"
```

**Esperado**: 8 tablas, con datos de prueba en users, schools, academic_units, materials, memberships

### MongoDB
```bash
# Ver colecciones
docker exec edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo --quiet --eval "db.getCollectionNames()"
```

**Esperado**: 9 colecciones creadas con validación de esquema e índices

---

## 🔧 Troubleshooting

### Error: "no se pudo conectar a PostgreSQL"
**Solución**: Asegúrate de que los servicios de base de datos estén corriendo:
```bash
docker-compose up -d postgres mongodb
# Espera 5-10 segundos
docker-compose -f docker-compose.migrate.yml up
```

### Error: "error eliminando schema"
**Solución**: Puede haber conexiones activas. Detén todos los servicios primero:
```bash
docker-compose down
docker-compose up -d postgres mongodb
# Espera 5-10 segundos
docker-compose -f docker-compose.migrate.yml up
```

### Quiero mantener los datos actuales
**Solución**: Este proceso NO es para ti. El migrator normal (sin FORCE_MIGRATION) es idempotente y no elimina datos:
```bash
docker-compose up  # Inicio normal
```

---

## 📝 Notas Técnicas

### ¿Por qué construir la imagen?

El `docker-compose.migrate.yml` usa `build` en vez de `image`:

```yaml
migrator:
  build:
    context: ../migrator
    dockerfile: Dockerfile
```

**Beneficios**:
1. Ejecuta `go mod download` → obtiene últimas versiones de `edugo-infrastructure`
2. Clona/actualiza el repo dentro del contenedor → estructura más reciente
3. Compila en el momento → sin necesidad de imagen pre-publicada

### ¿Qué hace FORCE_MIGRATION=true?

Cuando `FORCE_MIGRATION=true`:
1. **PostgreSQL**: Ejecuta `DROP SCHEMA public CASCADE` + `CREATE SCHEMA public`
2. **MongoDB**: Ejecuta `db.dropDatabase()`
3. Luego ejecuta las migraciones normalmente

Cuando `FORCE_MIGRATION` no está definida (modo normal):
1. Verifica si existen tablas/colecciones
2. Si existen → salta migraciones (idempotente)
3. Si no existen → ejecuta migraciones

---

## 🎯 Flujo de Trabajo Recomendado

Para el equipo de desarrollo:

```bash
# Desarrollo diario (NO elimina datos)
docker-compose up

# Actualización de esquema (elimina y recrea)
docker-compose -f docker-compose.migrate.yml up
docker-compose up -d
```

Para programadores que reciben instrucción "actualiza tu base de datos":

```bash
cd edugo-dev-environment/docker
docker-compose -f docker-compose.migrate.yml up
docker-compose up -d
```

**Simple, rápido, y siempre con la estructura más reciente** ✅
