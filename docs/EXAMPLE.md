# Ejemplo End-to-End - EduGo Dev Environment

**Objetivo:** Guía paso a paso para levantar y usar el ambiente de desarrollo completo de EduGo.

**Tiempo estimado:** 15-20 minutos (primera vez)

---

## 🎯 Lo Que Lograrás

Al final de esta guía tendrás:
- ✅ Todos los servicios corriendo (PostgreSQL, MongoDB, RabbitMQ, APIs, Worker)
- ✅ Datos de prueba cargados
- ✅ APIs respondiendo en http://localhost:8081 y :8082
- ✅ Worker procesando mensajes
- ✅ Ambiente listo para desarrollo

---

## 📋 Paso 1: Verificar Requisitos Previos

### 1.1 Docker Desktop

```bash
# Verificar que Docker está instalado
docker --version
# Esperado: Docker version 20.10.x o superior

# Verificar que Docker está corriendo
docker ps
# Esperado: Listado de contenedores (puede estar vacío)
```

**Si falla:**
```bash
# macOS
open -a Docker

# Esperar a que inicie (ver ícono en barra de menú)
```

### 1.2 Autenticación GitHub

Necesitas acceso a GitHub Container Registry para descargar las imágenes.

```bash
# Verificar autenticación
docker login ghcr.io

# Si no estás autenticado:
# Username: tu-usuario-github
# Password: tu-personal-access-token
```

**Crear token si no tienes:**
1. https://github.com/settings/tokens
2. Generate new token (classic)
3. Scope: `read:packages`
4. Copiar token

---

## 📋 Paso 2: Clonar y Configurar

### 2.1 Clonar Repositorio

```bash
# Clonar
git clone https://github.com/EduGoGroup/edugo-dev-environment.git

# Entrar al directorio
cd edugo-dev-environment

# Verificar contenido
ls -la
```

**Esperado:**
```
docker/
docs/
scripts/
migrator/
README.md
...
```

### 2.2 Ejecutar Setup

```bash
# Ejecutar script de setup completo
./scripts/setup.sh
```

**El script hará:**
1. Verificar Docker Desktop corriendo
2. Solicitar credenciales de GitHub (si no estás autenticado)
3. Descargar imágenes Docker
4. Crear archivo `.env` desde `.env.example`
5. Levantar todos los servicios
6. Ejecutar migraciones automáticas

**Output esperado:**
```
🚀 EduGo - Setup de Ambiente de Desarrollo

✅ Docker Desktop está corriendo
✅ Autenticado en ghcr.io

📥 Descargando imágenes Docker...
✅ Imágenes descargadas

📄 Creando archivo .env...
✅ Archivo .env creado

🚀 Levantando servicios...
[+] Running 7/7
 ✔ Container edugo-postgres    Started
 ✔ Container edugo-mongodb     Started
 ✔ Container edugo-rabbitmq    Started
 ✔ Container edugo-migrator    Started
 ✔ Container edugo-api-mobile  Started
 ✔ Container edugo-api-admin   Started
 ✔ Container edugo-worker      Started

✅ Setup completado exitosamente

🎉 Ambiente listo!

Verificar servicios:
  cd docker && docker-compose ps

Ver logs:
  docker-compose logs -f

Detener:
  ./scripts/stop.sh
```

---

## 📋 Paso 3: Verificar Servicios

### 3.1 Ver Estado de Contenedores

```bash
cd docker
docker-compose ps
```

**Esperado:**
```
NAME                    STATUS        PORTS
edugo-postgres          Up 2min       0.0.0.0:5432->5432/tcp
edugo-mongodb           Up 2min       0.0.0.0:27017->27017/tcp
edugo-rabbitmq          Up 2min       0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
edugo-api-mobile        Up 1min       0.0.0.0:8081->8081/tcp
edugo-api-admin         Up 1min       0.0.0.0:8082->8082/tcp
edugo-worker            Up 1min       
edugo-migrator          Exited (0)    
```

**Notas:**
- `edugo-migrator` debe mostrar `Exited (0)` - Esto es correcto, ejecutó migraciones y terminó
- Si algún servicio muestra `unhealthy` o `Restarting`, ver logs: `docker-compose logs [servicio]`

### 3.2 Verificar Logs

```bash
# Ver logs de todos los servicios
docker-compose logs --tail=50

# Ver logs en tiempo real
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f api-mobile
```

**Buscar:**
- ✅ PostgreSQL: `database system is ready to accept connections`
- ✅ MongoDB: `Waiting for connections`
- ✅ RabbitMQ: `Server startup complete`
- ✅ API Mobile: `Starting server on :8081`
- ✅ API Admin: `Starting server on :8082`
- ✅ Worker: `Worker started, waiting for messages`

---

## 📋 Paso 4: Probar Conexiones

### 4.1 Health Checks de APIs

```bash
# API Mobile
curl http://localhost:8081/health
```

**Esperado:**
```json
{
  "status": "ok",
  "database": "connected",
  "mongodb": "connected",
  "rabbitmq": "connected"
}
```

```bash
# API Administración
curl http://localhost:8082/health
```

**Esperado:**
```json
{
  "status": "ok",
  "database": "connected",
  "mongodb": "connected",
  "rabbitmq": "connected"
}
```

### 4.2 RabbitMQ Management UI

```bash
# Abrir en navegador
open http://localhost:15672
```

**Credenciales:**
- Usuario: `edugo`
- Password: `edugo123`

**Verificar:**
- ✅ Dashboard carga correctamente
- ✅ Connections muestra las APIs y Worker conectados
- ✅ Queues muestra las colas configuradas

### 4.3 PostgreSQL

```bash
# Conectar usando psql (si está instalado)
psql -h localhost -U edugo -d edugo

# O usando Docker
docker exec -it postgres psql -U edugo -d edugo
```

**Dentro de psql:**
```sql
-- Ver tablas
\dt

-- Debe mostrar tablas de migraciones:
users
institutions
courses
...

-- Verificar datos
SELECT COUNT(*) FROM users;

-- Salir
\q
```

### 4.4 MongoDB

```bash
# Conectar usando mongosh (si está instalado)
mongosh mongodb://edugo:edugo123@localhost:27017/edugo

# O usando Docker
docker exec -it mongodb mongosh -u edugo -p edugo123
```

**Dentro de mongosh:**
```javascript
// Ver colecciones
show collections

// Debe estar vacía o tener colecciones iniciales
// (se llenarán cuando el worker procese PDFs)

// Salir
exit
```

---

## 📋 Paso 5: Cargar Datos de Prueba (Opcional)

```bash
# Volver a la raíz del proyecto
cd ..

# Ejecutar script de seed
./scripts/seed-data.sh
```

**El script cargará:**
- Usuarios de prueba (admin, profesores, estudiantes)
- Instituciones de ejemplo
- Cursos de prueba
- Configuración inicial

**Verificar datos cargados:**
```bash
# PostgreSQL
docker exec -it postgres psql -U edugo -d edugo -c "SELECT COUNT(*) FROM users;"
```

**Esperado:**
```
 count 
-------
    10
(1 row)
```

---

## 📋 Paso 6: Probar Funcionalidad End-to-End

### 6.1 Crear Usuario Vía API

```bash
# Crear un nuevo usuario
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@edugo.com",
    "password": "test123",
    "firstName": "Test",
    "lastName": "User",
    "role": "student"
  }'
```

**Esperado:**
```json
{
  "id": "uuid-aqui",
  "email": "test@edugo.com",
  "firstName": "Test",
  "lastName": "User",
  "role": "student",
  "createdAt": "2025-11-22T..."
}
```

### 6.2 Login

```bash
# Hacer login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@edugo.com",
    "password": "test123"
  }'
```

**Esperado:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "test@edugo.com",
    ...
  }
}
```

### 6.3 Subir PDF (Trigger Worker)

```bash
# Guardar el token de arriba
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Subir un PDF de prueba
curl -X POST http://localhost:8081/api/v1/documents \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/test.pdf" \
  -F "title=Test Document"
```

**Esperado:**
```json
{
  "id": "doc-uuid",
  "title": "Test Document",
  "status": "processing",
  "createdAt": "..."
}
```

**Verificar que Worker procesa:**
```bash
# Ver logs del worker
docker-compose logs -f worker

# Buscar:
# "Processing document: doc-uuid"
# "Document processed successfully"
```

**Verificar en MongoDB:**
```bash
docker exec -it mongodb mongosh -u edugo -p edugo123 edugo

# En mongosh:
db.documents.find().pretty()
```

---

## 📋 Paso 7: Explorar APIs

### 7.1 API Mobile - Swagger

```bash
# Abrir documentación Swagger
open http://localhost:8081/swagger
```

**Explorar endpoints:**
- `/api/v1/auth/*` - Autenticación
- `/api/v1/users/*` - Usuarios
- `/api/v1/courses/*` - Cursos
- `/api/v1/documents/*` - Documentos

### 7.2 API Administración - Swagger

```bash
open http://localhost:8082/swagger
```

**Explorar endpoints admin:**
- `/api/v1/admin/institutions/*` - Gestión instituciones
- `/api/v1/admin/users/*` - Gestión usuarios
- `/api/v1/admin/reports/*` - Reportes

---

## 📋 Paso 8: Desarrollo Local

### 8.1 Modificar Código de API

Si quieres modificar el código de las APIs:

```bash
# Clonar repo de API
git clone https://github.com/EduGoGroup/edugo-api-mobile.git
cd edugo-api-mobile

# Detener contenedor de API en dev-environment
cd ../edugo-dev-environment/docker
docker-compose stop api-mobile

# Correr API localmente (conectará a DBs del dev-environment)
cd ../../edugo-api-mobile
go run main.go
```

**Configurar .env local:**
```env
DATABASE_URL=postgresql://edugo:edugo123@localhost:5432/edugo
MONGO_URI=mongodb://edugo:edugo123@localhost:27017/edugo
RABBITMQ_URL=amqp://edugo:edugo123@localhost:5672/
```

### 8.2 Ver Logs en Tiempo Real

```bash
# Terminal 1: Logs de API Mobile
docker-compose logs -f api-mobile

# Terminal 2: Logs de Worker
docker-compose logs -f worker

# Terminal 3: Logs de PostgreSQL
docker-compose logs -f postgres
```

---

## 📋 Paso 9: Detener Ambiente

### Opción A: Detener (Mantiene Datos)

```bash
cd docker
docker-compose stop
```

**Resultado:**
- ✅ Contenedores detenidos
- ✅ Datos en volúmenes preservados
- ✅ Próximo `docker-compose up -d` inicia rápido

### Opción B: Detener y Eliminar Contenedores

```bash
docker-compose down
```

**Resultado:**
- ✅ Contenedores eliminados
- ✅ Datos en volúmenes preservados
- ⚠️ Próximo inicio un poco más lento

### Opción C: Reset Completo (Elimina Datos)

```bash
# Usando script
cd ..
./scripts/cleanup.sh

# O manualmente
cd docker
docker-compose down -v
```

**Resultado:**
- ✅ Contenedores eliminados
- ❌ Volúmenes eliminados (datos borrados)
- ⚠️ Requiere re-ejecutar migraciones y seed

---

## 📋 Paso 10: Comandos Útiles

### Ver Uso de Recursos

```bash
# Ver CPU y memoria de contenedores
docker stats

# Ver espacio usado por Docker
docker system df
```

### Reiniciar Servicio Específico

```bash
cd docker

# Reiniciar solo API Mobile
docker-compose restart api-mobile

# Rebuild y reiniciar
docker-compose up -d --no-deps --build api-mobile
```

### Limpiar Logs

```bash
# Los logs se acumulan, limpiar periódicamente
docker system prune

# Ver tamaño de logs
du -sh $(docker inspect --format='{{.LogPath}}' $(docker ps -qa))
```

### Actualizar Imágenes

```bash
cd ..
./scripts/update-images.sh
```

---

## 🐛 Troubleshooting

### Problema: "Cannot connect to Docker daemon"

**Solución:**
```bash
open -a Docker
# Esperar a que inicie
```

### Problema: "pull access denied"

**Solución:**
```bash
docker login ghcr.io
# Usuario: tu-github-username
# Password: tu-personal-access-token
```

### Problema: "Port already in use"

**Solución:**
```bash
# Ver qué está usando el puerto
lsof -ti:5432 | xargs kill -9  # PostgreSQL
lsof -ti:8081 | xargs kill -9  # API Mobile
lsof -ti:8082 | xargs kill -9  # API Admin
```

### Problema: Servicios "unhealthy"

**Solución:**
```bash
# Ver logs del servicio
docker-compose logs postgres

# Reiniciar servicio
docker-compose restart postgres

# Si persiste, reset completo
docker-compose down -v
docker-compose up -d
```

---

## 🎉 ¡Listo!

Has completado el setup completo del ambiente de desarrollo EduGo.

### Próximos Pasos

1. **Explorar APIs** - http://localhost:8081/swagger
2. **Modificar código** - Clona los repos de APIs
3. **Crear features** - Desarrolla nuevas funcionalidades
4. **Probar Worker** - Sube PDFs y ve el procesamiento

### Recursos Adicionales

- **Quick Start:** [docker/QUICK_START.md](../docker/QUICK_START.md)
- **Documentación Completa:** [docker/README.md](../docker/README.md)
- **Scripts:** [scripts/README.md](../scripts/README.md)
- **Hooks:** [.githooks/README.md](../.githooks/README.md)

---

**Última actualización:** 22 de Noviembre, 2025  
**Versión:** 1.0  
**Mantenedor:** Equipo EduGo
