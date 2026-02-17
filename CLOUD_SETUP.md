# 🚀 EduGo - Configuración Completa en la Nube

> Toda la infraestructura de EduGo ahora en servicios cloud gratuitos

## ✅ Servicios Configurados

### PostgreSQL - Neon (Nube)
La base de datos PostgreSQL de EduGo ha sido migrada exitosamente a Neon (plan gratuito).

### Redis - Upstash (Nube)
Redis está configurado en Upstash para caché y sesiones (plan gratuito).

### MongoDB - Atlas (Nube)
MongoDB está configurado en MongoDB Atlas para almacenamiento de documentos (plan gratuito M0).

### 📊 Detalles de PostgreSQL (Neon)

- **Proyecto Neon**: MCPEco (gentle-shadow-07969581)
- **Base de datos**: `edugo`
- **Host**: `ep-green-frost-ado4abbi-pooler.c-2.us-east-1.aws.neon.tech`
- **Puerto**: 5432
- **Usuario**: `neondb_owner`
- **SSL**: Requerido

### 📊 Detalles de Redis (Upstash)

- **Nombre**: edugo-redis
- **Host**: `living-wildcat-41131.upstash.io`
- **Puerto**: 6379
- **Usuario**: `default`
- **Password**: `AaCrAAIncDJmMTFjYjJiOGU1M2U0YmM3YWIxMDQyZTA2ZjdlZDgxZXAyNDExMzE`
- **TLS**: Habilitado
- **URL**: `redis://default:AaCrAAIncDJmMTFjYjJiOGU1M2U0YmM3YWIxMDQyZTA2ZjdlZDgxZXAyNDExMzE@living-wildcat-41131.upstash.io:6379`

### 📊 Detalles de MongoDB (Atlas)

- **Cluster**: `edugo.alxme5j.mongodb.net`
- **Base de datos**: `edugo`
- **Usuario**: `medinatello_db_user`
- **Password**: `6NQjJDaOkN4nvldT`
- **Plan**: M0 (Free Tier)
- **URI**: `mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/?appName=Edugo`
- **URI con DB**: `mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/edugo?appName=Edugo`

### 📦 Contenido Migrado

**PostgreSQL (Neon):**
✅ **Estructura de base de datos** (todas las tablas, índices, constraints)
✅ **Datos iniciales (Seeds)** (roles, permisos, configuraciones del sistema)
✅ **Datos de prueba (Mock Data)** (usuarios, cursos, etc. para desarrollo)

**MongoDB (Atlas):**
✅ **Estructura de colecciones** (9 colecciones con schemas y validaciones)
⚠️ **Datos de prueba** (tienen problemas de validación - no aplicados)

## 🔧 Cómo Usar Neon en tu Desarrollo Local

### Opción 1: Usando Variables de Entorno

1. **Copia el archivo de ejemplo**:
   ```bash
   cp .env.neon .env
   ```

2. **Actualiza tu aplicación** para que lea las variables de entorno:
   - ✅ PostgreSQL ya NO necesita Docker (está en Neon)
   - ✅ Redis ya NO necesita Docker (está en Upstash)
   - ✅ MongoDB ya NO necesita Docker (está en Atlas)
   - ⚠️ Solo RabbitMQ sigue en Docker local

3. **Levanta solo RabbitMQ (opcional si usas mensajería)**:
   ```bash
   # Solo RabbitMQ (todo lo demás está en la nube)
   cd docker
   docker-compose up -d rabbitmq
   ```

   O si no usas RabbitMQ:
   ```bash
   # ¡No necesitas levantar nada! Todo está en la nube 🎉
   ```

### Opción 2: Configuración Manual

Usa la siguiente cadena de conexión en tu aplicación:

```
postgresql://neondb_owner:npg_sC2u9pTVwQJI@ep-green-frost-ado4abbi-pooler.c-2.us-east-1.aws.neon.tech:5432/edugo?sslmode=require
```

O las variables individuales:

```bash
POSTGRES_HOST=ep-green-frost-ado4abbi-pooler.c-2.us-east-1.aws.neon.tech
POSTGRES_PORT=5432
POSTGRES_USER=neondb_owner
POSTGRES_PASSWORD=npg_sC2u9pTVwQJI
POSTGRES_DB=edugo
POSTGRES_SSLMODE=require
```

## 🎯 Beneficios

**PostgreSQL (Neon):**
- ✅ **No necesitas levantar PostgreSQL con Docker**
- ✅ **Base de datos persistente** entre reinicios
- ✅ **Acceso desde cualquier lugar** (desarrollo remoto, colaboración)
- ✅ **Backups automáticos** (6 horas de point-in-time recovery)
- ✅ **0.5 GB de almacenamiento** gratuito

**Redis (Upstash):**
- ✅ **No necesitas levantar Redis con Docker**
- ✅ **256 MB de memoria** para caché y sesiones
- ✅ **500,000 comandos/mes** en plan gratuito
- ✅ **TLS habilitado** por seguridad
- ✅ **Acceso desde cualquier lugar**

**MongoDB (Atlas):**
- ✅ **No necesitas levantar MongoDB con Docker**
- ✅ **512 MB de almacenamiento** gratuito
- ✅ **Backups automáticos**
- ✅ **Acceso desde cualquier lugar**
- ✅ **Alta disponibilidad** (3 nodos replica set)

## 📝 Límites de los Planes Gratuitos

**Neon (PostgreSQL):**
- **Almacenamiento**: 0.5 GB
- **Cómputo**: 100 CU-horas/mes (~400 horas con 0.25 CU)
- **Transferencia de datos**: 5 GB/mes
- **Proyectos**: Hasta 20

**Upstash (Redis):**
- **Almacenamiento**: 256 MB
- **Comandos**: 500,000/mes
- **Ancho de banda**: 200 GB/mes gratis
- **Bases de datos**: Hasta 10

**MongoDB Atlas:**
- **Almacenamiento**: 512 MB
- **RAM**: Compartida
- **Conexiones**: 500 simultáneas
- **Clusters**: 1 por proyecto (plan M0)

## 🔄 Recrear las Bases de Datos

### PostgreSQL (Neon)

#### Opción 1: Con datos de prueba (por defecto)
```bash
cd migrator
./recreate_neon_db.sh
```

### Opción 2: Solo estructura y seeds (SIN datos de prueba)
```bash
cd migrator
APPLY_MOCK_DATA=false ./recreate_neon_db.sh
```

### Opción 3: Manualmente
```bash
cd migrator
FORCE_MIGRATION=true go run migrate_to_neon.go              # Con datos de prueba
FORCE_MIGRATION=true APPLY_MOCK_DATA=false go run migrate_to_neon.go  # Sin datos de prueba
```

⚠️ **ADVERTENCIA**: Esto eliminará TODOS los datos y recreará la base de datos.

### MongoDB (Atlas)

#### Opción 1: Solo estructura (recomendado - sin datos de prueba)
```bash
cd migrator
APPLY_MOCK_DATA=false ./recreate_atlas_db.sh
```

#### Opción 2: Manualmente
```bash
cd migrator
FORCE_MIGRATION=true APPLY_MOCK_DATA=false go run migrate_to_atlas.go
```

⚠️ **NOTA**: Los datos de prueba de MongoDB tienen problemas de validación. Se recomienda usar `APPLY_MOCK_DATA=false`.

### 🔄 Equivalencia con Docker

| Docker (antes) | Nube (ahora) |
|----------------|--------------|
| **PostgreSQL**: `docker-compose down -v && up` | `./recreate_neon_db.sh` |
| **MongoDB**: `docker-compose down -v && up` | `./recreate_atlas_db.sh` |
| Elimina contenedor + volúmenes | Elimina base de datos + recrea |
| Migrator aplica todo automáticamente | Scripts aplican todo |

## ✅ Probar las Conexiones

### PostgreSQL (Neon)
```bash
cd migrator
go run migrate_to_neon.go  # Si ya existe, verá "migraciones omitidas"
```

### Redis (Upstash)
```bash
cd migrator
go run test_redis_connection.go
```

Deberías ver:
- ✅ PING exitoso: PONG
- ✅ SET exitoso
- ✅ GET exitoso
- ✅ TTL
- ✅ Clave eliminada

### MongoDB (Atlas)
```bash
cd migrator
go run test_mongodb_connection.go
```

Deberías ver:
- ✅ PING exitoso
- ✅ Versión de MongoDB: 8.0.19
- ✅ Colecciones encontradas
- ✅ Documento insertado/recuperado/eliminado

## 📂 Archivos de Configuración

- `neon-config.yaml` - Configuración completa (PostgreSQL, Redis, MongoDB)
- `.env.neon` - Variables de entorno listas para usar

**Scripts de PostgreSQL (Neon):**
- `migrator/migrate_to_neon.go` - Script de migración PostgreSQL a Neon
- `migrator/recreate_neon_db.sh` - Script de recreación PostgreSQL

**Scripts de Redis (Upstash):**
- `migrator/test_redis_connection.go` - Script de prueba Redis Upstash

**Scripts de MongoDB (Atlas):**
- `migrator/migrate_to_atlas.go` - Script de migración MongoDB a Atlas
- `migrator/test_mongodb_connection.go` - Script de prueba MongoDB Atlas
- `migrator/recreate_atlas_db.sh` - Script de recreación MongoDB

## 🆘 Troubleshooting

### Error de conexión SSL

Asegúrate de incluir `sslmode=require` en tu cadena de conexión.

### Límite de almacenamiento

Si alcanzas el límite de 0.5 GB:
1. Limpia datos de prueba innecesarios
2. Considera crear un proyecto Neon dedicado
3. Evalúa actualizar al plan pagado

### MongoDB sigue necesitando Docker

Correcto. Por ahora solo PostgreSQL está en Neon. MongoDB sigue corriendo localmente en Docker.

## 📚 Referencias

- [Documentación de Neon](https://neon.tech/docs)
- [Planes de Neon](https://neon.tech/pricing)
- [edugo-infrastructure](https://github.com/EduGoGroup/edugo-infrastructure)
