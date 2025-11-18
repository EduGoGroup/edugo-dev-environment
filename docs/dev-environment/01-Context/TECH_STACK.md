# Tech Stack - Infraestructura Docker de EduGo

## 📋 Overview del Stack Tecnológico

edugo-dev-environment orquesta una arquitectura **multi-contenedor** basada en las siguientes tecnologías:

---

## 🗄️ Almacenamiento de Datos

### PostgreSQL 15 (Base de Datos Relacional Principal)

**Versión:** 15.x (Alpine)  
**Imagen:** `postgres:15-alpine`  
**Puerto:** 5432

#### Características
- Soporte de CTEs (Common Table Expressions) para consultas recursivas
- JSONB para datos semi-estructurados
- Full-text search
- UUID generación nativa
- Triggers y stored procedures
- Connection pooling

#### Uso en EduGo
| Tabla | Descripción | Proyecto |
|-------|-------------|----------|
| `users` | Usuarios del sistema | Todos |
| `schools` | Instituciones educativas | api-admin |
| `academic_units` | Estructura jerárquica académica | api-admin |
| `materials` | Contenido educativo (PDFs) | api-mobile |
| `assessment_attempts` | Intentos de evaluación | api-mobile |
| `memberships` | Asignaciones usuario-rol-unidad | api-admin |
| `processing_status` | Estado de procesamiento de materiales | worker |

#### Volumen Persistente
```
postgres_data:/var/lib/postgresql/data
Tamaño estimado: 5GB para producción
```

#### Configuración Inicial
```sql
-- Usuario de aplicación
CREATE ROLE edugo_user WITH PASSWORD 'password' LOGIN;

-- Base de datos
CREATE DATABASE edugo_dev OWNER edugo_user;

-- Extensiones
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

---

### MongoDB 7.0 (Base de Datos de Documentos)

**Versión:** 7.0.x  
**Imagen:** `mongo:7.0`  
**Puerto:** 27017

#### Características
- Replica set para consistencia
- Transacciones ACID multi-documento
- Índices compostos
- Aggregation framework
- Change streams (observar cambios)
- Compresión nativa

#### Uso en EduGo
| Colección | Descripción | Proyecto |
|-----------|-------------|----------|
| `material_summary` | Resúmenes generados por IA | worker |
| `material_assessment` | Quizzes generados por IA | worker |
| `material_events` | Eventos de procesamiento de materiales | worker |

#### Volumen Persistente
```
mongodb_data:/data/db
Tamaño estimado: 2GB para desarrollo
```

#### Configuración Inicial
```javascript
// Replica set simple
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongo:27017" }
  ]
});

// Base de datos y usuario
db.createUser({
  user: "edugo_user",
  pwd: "password",
  roles: ["readWrite"]
});
```

#### Índices Recomendados
```javascript
// material_summary
db.material_summary.createIndex({ material_id: 1, created_at: -1 });
db.material_summary.createIndex({ material_id: 1 }, { unique: true });

// material_assessment
db.material_assessment.createIndex({ material_id: 1, created_at: -1 });
db.material_assessment.createIndex({ question_count: 1 });
```

---

## 📨 Mensajería

### RabbitMQ 3.12 (Message Broker)

**Versión:** 3.12.x con Management  
**Imagen:** `rabbitmq:3.12-management`  
**Puertos:** 5672 (AMQP), 15672 (Management UI)

#### Características
- Topic exchanges para pub/sub
- Message acknowledgments
- Persistent messages
- Dead letter queues (DLQ)
- Consumer prefetch control
- Management UI integrado

#### Exchanges Predefinidos
```
Name: material-events
Type: topic
Durable: true
Auto-delete: false
```

#### Queues Predefinidas
```
Queue: material.processing
Exchange: material-events
Routing key: material.#
Durable: true
Consumer: worker

Queue: material.processed
Exchange: material-events
Routing key: processing.#
Durable: true
Consumer: api-mobile (opcional)
```

#### Flujos de Eventos
```
1. API publica evento:
   "material.created" → material-events exchange
   
2. Worker consume:
   material-events → material.processing queue
   
3. Worker procesa y publica:
   "processing.completed" → material-events exchange
   
4. API consume resultado (opcional):
   material-events → material.processed queue
```

#### Management UI
- URL: http://localhost:15672
- Usuario: guest
- Contraseña: guest
- Funcionalidad: Ver colas, mensajes, conexiones

#### Volumen Persistente
```
rabbitmq_data:/var/lib/rabbitmq
Tamaño estimado: 500MB
```

---

## 💾 Cache y Sesiones

### Redis 7.0 (In-Memory Cache)

**Versión:** 7.0.x (Alpine)  
**Imagen:** `redis:7.0-alpine`  
**Puerto:** 6379

#### Características
- Strings, Hashes, Lists, Sets, Sorted Sets
- TTL automático para expiración
- Persistencia RDB y AOF (opcional)
- Pub/Sub simple
- Lua scripting

#### Uso en EduGo (Planeado)
```
key: "session:{user_id}"        → Token de sesión (TTL: 24h)
key: "assessment:{attempt_id}"  → Cache de intento (TTL: 1h)
key: "cache:hierarchy:{unit_id}" → Caché de jerarquía (TTL: 1h)
```

#### Comandos Útiles
```bash
# Conectar
redis-cli -h localhost -p 6379

# Verificar keys
KEYS *

# Ver tipo de dato
TYPE key_name

# Ver TTL
TTL key_name
```

#### Volumen Persistente
```
redis_data:/data
Tamaño estimado: 100MB
```

---

## 🖥️ Interfaces Web

### PgAdmin 4 (Cliente PostgreSQL)

**Versión:** Última (weekly release)  
**Imagen:** `dpage/pgadmin4:latest`  
**Puerto:** 5050

#### Características
- Gestión visual de bases de datos
- Editor SQL con syntax highlighting
- Backup y restore
- Estadísticas y monitoreo
- Diseñador de esquemas

#### Acceso
- URL: http://localhost:5050
- Email: admin@edugo.local
- Contraseña: admin
- Servidor pre-configurado: postgres:5432

#### Tareas Típicas
```
1. Navegar a Servers → postgres → Databases → edugo_dev
2. Expandir "Schemas" para ver tablas
3. Tools → Query Tool para escribir SQL
4. Tools → Backup para hacer respaldos
```

---

### Mongo Express (Cliente MongoDB)

**Versión:** Última  
**Imagen:** `mongo-express:latest`  
**Puerto:** 8081

#### Características
- Visualización de bases de datos
- Editor de documentos JSON
- Creación de colecciones
- Índices y validadores
- Importar/Exportar datos

#### Acceso
- URL: http://localhost:8081
- Usuario: admin
- Contraseña: pass

#### Tareas Típicas
```
1. Seleccionar base de datos "edugo_dev"
2. Navegar a colecciones
3. Hacer clic en documentos para editar
4. Crear índices desde interfaz
5. Importar datos JSON
```

---

### RabbitMQ Management (Panel de Administración)

**Versión:** 3.12.x (integrado)  
**Puerto:** 15672

#### Características
- Visualización de exchanges, queues, bindings
- Monitoreo de mensajes
- Usuarios y permisos
- Estadísticas de rendimiento
- Purgar colas

#### Acceso
- URL: http://localhost:15672
- Usuario: guest
- Contraseña: guest

#### Tareas Típicas
```
1. Ir a "Queues" para ver estado
2. Monitorear "Overview" para tráfico
3. Crear exchanges/queues manualmente (si es necesario)
4. Ver mensajes pendientes
5. Purgar colas en desarrollo
```

---

## 🔌 Networking

### Docker Compose Network

**Tipo:** Bridge network (default)  
**Nombre:** `edugo_default`

#### Conectividad
```
Contenedor → Contenedor: hostname interno
postgres:5432
mongo:27017
rabbitmq:5672
redis:6379

Host → Contenedor: localhost:port
localhost:5432  → postgres
localhost:27017 → mongo
localhost:5672  → rabbitmq
localhost:6379  → redis
localhost:5050  → pgadmin
localhost:8081  → mongo-express
localhost:15672 → rabbitmq-management
```

#### DNS Interno
```
# Desde dentro de un contenedor (ej: worker)
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
DATABASE_URL=postgres://user:pass@postgres:5432/edugo_dev
MONGO_URI=mongodb://mongo:27017
REDIS_URL=redis://:password@redis:6379
```

---

## 🔐 Seguridad

### Credenciales por Defecto

⚠️ **SOLO PARA DESARROLLO LOCAL**

```bash
PostgreSQL
  Usuario: postgres
  Contraseña: postgres
  Admin app: edugo_user / secure_password_change_in_prod

MongoDB
  Usuario: root
  Contraseña: mongo
  App user: edugo_user / secure_password_change_in_prod

RabbitMQ
  Usuario: guest
  Contraseña: guest

Redis
  Contraseña: (sin contraseña por defecto)

PgAdmin
  Usuario: admin@edugo.local
  Contraseña: admin

Mongo Express
  Usuario: admin
  Contraseña: pass
```

### Política de Seguridad para Producción

```yaml
❌ NO usar credenciales por defecto
✅ Usar secrets management (SOPS, Vault)
✅ Habilitar autenticación en todos los servicios
✅ Usar TLS/SSL para comunicaciones
✅ Limitar acceso a puertos de management (5050, 8081, 15672)
✅ Usar redes privadas (no exponer en internet)
✅ Implementar backups encriptados
```

---

## 📊 Recursos Recomendados

### Por Ambiente

**Desarrollo Local (Docker Desktop)**
```
CPU: 4 cores
RAM: 8GB (mínimo), 16GB (recomendado)
Disk: 10GB libres
```

**Testing/Staging**
```
CPU: 8 cores
RAM: 16GB
Disk: 50GB
```

**Producción**
```
CPU: 16+ cores
RAM: 32GB+
Disk: 500GB+ (SSD)
```

### Asignación por Contenedor
```
PostgreSQL:  2GB RAM, 2 CPUs
MongoDB:     2GB RAM, 2 CPUs
RabbitMQ:    1GB RAM, 1 CPU
Redis:       512MB RAM, 1 CPU
PgAdmin:     512MB RAM
Mongo Exp:   256MB RAM
```

---

## 🔄 Ciclo de Vida de Datos

### Volúmenes Persistentes

```
docker volume ls

# Estructura
edugo_postgres_data   → /var/lib/postgresql/data
edugo_mongodb_data    → /data/db
edugo_rabbitmq_data   → /var/lib/rabbitmq
edugo_redis_data      → /data
```

### Backup Strategy

```bash
# PostgreSQL dump
docker-compose exec postgres pg_dump -U edugo_user edugo_dev > backup.sql

# MongoDB dump
docker-compose exec mongo mongodump --db edugo_dev --out /backup

# Restaurar PostgreSQL
cat backup.sql | docker-compose exec -T postgres psql -U edugo_user edugo_dev

# Restaurar MongoDB
docker-compose exec mongo mongorestore /backup
```

### Reset Completo (Desarrollo)

```bash
# ⚠️ PELIGRO: Borra todos los datos
docker-compose down -v
docker-compose up -d
# Datos iniciales se recrean automáticamente
```

---

## 🚀 Performance Tuning

### PostgreSQL
```sql
-- Aumentar conexiones simultáneas
max_connections = 200

-- Mejorar cache
shared_buffers = 256MB
effective_cache_size = 1GB

-- WAL optimization
wal_level = replica
max_wal_senders = 10
```

### MongoDB
```javascript
// Usar índices apropiados
db.material_summary.createIndex({ material_id: 1 });
db.material_summary.createIndex({ created_at: -1 });

// Aggregation pipeline optimization
// - $match early
// - $project para limitar campos
// - $limit antes de $lookup
```

### RabbitMQ
```bash
# Consumer prefetch
channel.basicQos(1);  // Procesar 1 mensaje a la vez

# Persistent queues
durable = true
```

### Redis
```bash
# Configuración
maxmemory 2gb
maxmemory-policy allkeys-lru
```

---

## 📚 Versiones Específicas

```yaml
PostgreSQL:
  Estable: 15.x
  Soporte: Hasta octubre 2026
  Features críticas: CTEs, JSONB, Full-text search

MongoDB:
  Estable: 7.0.x
  Soporte: Hasta septiembre 2027
  Features críticas: Replica sets, Transactions

RabbitMQ:
  Estable: 3.12.x
  Soporte: Hasta abril 2027
  Features críticas: Topic exchanges, DLQ

Redis:
  Estable: 7.0.x
  Soporte: Hasta junio 2025
  Features críticas: Strings, Hashes, TTL
```

---

## 🔗 Referencias y Documentación

- [PostgreSQL 15 Docs](https://www.postgresql.org/docs/15/)
- [MongoDB 7.0 Docs](https://docs.mongodb.com/manual/release-notes/7.0/)
- [RabbitMQ 3.12 Docs](https://www.rabbitmq.com/documentation.html)
- [Redis 7.0 Docs](https://redis.io/docs/)
- [Docker Compose Docs](https://docs.docker.com/compose/)
- [PgAdmin Docs](https://www.pgadmin.org/docs/)
- [Mongo Express Docs](https://github.com/mongo-express/mongo-express)

---

**Última actualización:** 15 de Noviembre, 2025  
**Generado con:** Claude Code  
**Proyecto:** edugo-dev-environment  
**Tipo:** Especificación técnica de stack
