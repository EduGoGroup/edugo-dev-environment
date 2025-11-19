# EduGo - Guía de Uso Docker Compose

Esta guía explica cómo usar los 3 archivos docker-compose disponibles en este proyecto.

## 📋 Archivos Disponibles

1. **docker-compose.yml** - Instalación completa (infraestructura + aplicaciones)
2. **docker-compose-infrastructure.yml** - Solo servicios externos (PostgreSQL, MongoDB, RabbitMQ, Redis)
3. **docker-compose-apps.yml** - Solo aplicaciones EduGo (APIs + Worker)
4. **docker-compose.migrate.yml** - ⚠️ Actualización forzada de base de datos (elimina y recrea)

---

## 🚀 Opción 1: Instalación Completa

**Cuándo usar**: Desarrollo local completo, testing end-to-end, demos

```bash
# Levantar todo el ecosistema
cd docker
docker-compose --profile full up -d

# Ver logs
docker-compose logs -f

# Detener todo
docker-compose down
```

**Servicios incluidos**:
- PostgreSQL (puerto 5432)
- MongoDB (puerto 27017)
- RabbitMQ (puertos 5672, 15672)
- Redis (puerto 6379) - opcional con profile
- API Mobile (puerto 8081)
- API Administración (puerto 8082) - requiere config.yaml
- Worker - requiere config.yaml

**URLs importantes**:
- API Mobile Swagger: http://localhost:8081/swagger/index.html
- API Mobile Health: http://localhost:8081/health
- RabbitMQ Management: http://localhost:15672 (usuario: edugo, password: edugo123)

---

## 🗄️ Opción 2: Solo Infraestructura

**Cuándo usar**: Desarrollo de APIs en local pero corriendo contra bases de datos en Docker

```bash
# Levantar solo servicios externos
cd docker
docker-compose -f docker-compose-infrastructure.yml up -d

# Opcional: incluir Redis
docker-compose -f docker-compose-infrastructure.yml --profile redis up -d

# Ver estado
docker-compose -f docker-compose-infrastructure.yml ps

# Detener
docker-compose -f docker-compose-infrastructure.yml down
```

**Servicios incluidos**:
- PostgreSQL (puerto 5432)
- MongoDB (puerto 27017)
- RabbitMQ (puertos 5672, 15672)
- Redis (puerto 6379) - solo con `--profile redis`

**Conexiones desde aplicaciones locales**:
```bash
# PostgreSQL
postgresql://edugo:edugo123@localhost:5432/edugo

# MongoDB
mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin

# RabbitMQ
amqp://edugo:edugo123@localhost:5672/

# Redis
redis://localhost:6379
```

---

## 🔧 Opción 3: Solo Aplicaciones

**Cuándo usar**: Cuando ya tienes la infraestructura corriendo (desde Opción 2 u otro ambiente)

**Prerequisitos**: 
1. La infraestructura debe estar corriendo (Opción 2)
2. La red `edugo-network` debe existir

```bash
# Paso 1: Asegurar que la red existe
docker network create edugo-network 2>/dev/null || true

# Paso 2: Levantar aplicaciones
cd docker
docker-compose -f docker-compose-apps.yml up -d api-mobile

# Para API Admin y Worker (requieren config.yaml)
# docker-compose -f docker-compose-apps.yml --profile admin up -d
# docker-compose -f docker-compose-apps.yml --profile worker up -d

# Ver logs
docker-compose -f docker-compose-apps.yml logs -f

# Detener
docker-compose -f docker-compose-apps.yml down
```

**Servicios incluidos**:
- API Mobile (puerto 8081)
- API Administración (puerto 8082) - requiere config.yaml, usar `--profile admin`
- Worker - requiere config.yaml, usar `--profile worker`

---

## 🔄 Opción 4: Actualización de Base de Datos

**Cuándo usar**: Cuando el equipo de backend actualiza el esquema de base de datos

⚠️ **ADVERTENCIA**: Este proceso **ELIMINA** completamente las bases de datos y las recrea desde cero.

```bash
# Actualizar esquema de base de datos (elimina y recrea)
cd docker
docker-compose -f docker-compose.migrate.yml up

# Reiniciar servicios
docker-compose up -d
```

**¿Qué hace?**:
1. ✅ Construye imagen del migrator (obtiene últimas dependencias)
2. ✅ Clona/actualiza repositorio edugo-infrastructure
3. 🔥 Elimina PostgreSQL schema y MongoDB database
4. ✅ Recrea con estructura más reciente
5. ✅ Carga datos de prueba actualizados

📖 **Documentación completa**: Ver [ACTUALIZAR_BASE_DATOS.md](./ACTUALIZAR_BASE_DATOS.md)

---

## 🔍 Validación y Testing

### Validar API Mobile

```bash
# Health check
curl http://localhost:8081/health | jq

# Swagger UI
open http://localhost:8081/swagger/index.html
```

### Validar Infraestructura

```bash
# PostgreSQL
docker exec -it edugo-postgres psql -U edugo -d edugo -c "SELECT version();"

# MongoDB
docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin --eval "db.adminCommand('ping')"

# RabbitMQ
curl -u edugo:edugo123 http://localhost:15672/api/overview | jq
```

---

## 📝 Variables de Entorno

Todas las variables están en el archivo `.env`. Puedes copiar `.env.example` para crear tu propio `.env`:

```bash
cp .env.example .env
```

**Variables importantes**:
- `POSTGRES_PASSWORD` - Contraseña de PostgreSQL
- `MONGO_PASSWORD` - Contraseña de MongoDB
- `S3_ACCESS_KEY` - Access key de S3 (temporal, marcado como opcional)
- `OPENAI_API_KEY` - API key de OpenAI para el worker
- `BOOTSTRAP_OPTIONAL_RESOURCES_S3=true` - Permite que API Mobile arranque sin S3

---

## 🐛 Troubleshooting

### Error: Puerto en uso

```bash
# Liberar puerto 8081
lsof -ti:8081 | xargs kill -9

# Liberar puerto 8082
lsof -ti:8082 | xargs kill -9
```

### Error: API Admin/Worker no inician

**Causa**: Requieren archivo `config.yaml` en sus Dockerfiles

**Solución temporal**: Usar solo API Mobile

**Solución permanente**: Ver [RESULTADO_VALIDACION.md](./RESULTADO_VALIDACION.md) sección "Problemas Encontrados"

### Ver logs de un servicio específico

```bash
# Con docker-compose.yml
docker-compose logs -f api-mobile

# Con archivo específico
docker-compose -f docker-compose-infrastructure.yml logs -f postgres
```

### Reiniciar un servicio

```bash
docker-compose restart api-mobile
```

### Eliminar volúmenes y empezar de cero

```bash
# CUIDADO: Esto elimina todos los datos
docker-compose down -v
docker volume rm edugo-postgres-data edugo-mongodb-data edugo-rabbitmq-data edugo-redis-data 2>/dev/null || true
```

---

## 📊 Estado Actual del Proyecto

| Servicio | Estado | Puerto | Swagger | Notas |
|----------|--------|--------|---------|-------|
| PostgreSQL | ✅ Funcionando | 5432 | - | Health check OK |
| MongoDB | ✅ Funcionando | 27017 | - | Health check OK |
| RabbitMQ | ✅ Funcionando | 5672, 15672 | - | Management UI disponible |
| Redis | ✅ Funcionando | 6379 | - | Opcional (profile redis) |
| API Mobile | ✅ Funcionando | 8081 | ✅ | Completamente operativa |
| API Admin | ⚠️ Requiere config | 8082 | - | Necesita config.yaml en Dockerfile |
| Worker | ⚠️ Requiere config | - | - | Necesita config.yaml en Dockerfile |

Para más detalles sobre problemas y soluciones, consulta [RESULTADO_VALIDACION.md](./RESULTADO_VALIDACION.md)

---

## 🎯 Workflows Recomendados

### Desarrollo Full-Stack Local
```bash
# Levantar todo
docker-compose --profile full up -d
# Trabajar en tu código
# Ver logs cuando necesites
docker-compose logs -f api-mobile
```

### Desarrollo Backend (Go APIs)
```bash
# Levantar solo infraestructura
docker-compose -f docker-compose-infrastructure.yml up -d
# Correr tus APIs localmente en tu IDE
# Apuntar a localhost:5432, localhost:27017, etc.
```

### Testing de Integración
```bash
# Infraestructura + API Mobile
docker-compose -f docker-compose-infrastructure.yml up -d
docker-compose -f docker-compose-apps.yml up -d api-mobile
# Ejecutar tests
# Limpiar
docker-compose -f docker-compose-apps.yml down
docker-compose -f docker-compose-infrastructure.yml down
```

---

## 📚 Documentación Adicional

- [RESULTADO_VALIDACION.md](./RESULTADO_VALIDACION.md) - Reporte completo de validación
- [../docs/dev-environment/](../docs/dev-environment/) - Documentación del proyecto
- Variables de entorno: `.env.example`
