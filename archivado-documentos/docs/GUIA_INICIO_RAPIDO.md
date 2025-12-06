# 📱 Guía para Programadores Frontend/Mobile - EduGo

**Fecha:** 30 de Octubre, 2025
**Versión:** 1.0
**Organización:** EduGoGroup

---

## 🎯 ¿Qué necesitas para empezar?

Simplemente necesitas:
1. **Docker Desktop** instalado
2. **GitHub token** con acceso a la organización EduGoGroup
3. **10 minutos** para setup

---

## ⚡ Setup Rápido (Primera Vez)

### PASO 1: Obtener GitHub Token

1. Ve a: https://github.com/settings/tokens/new
2. Nombre del token: `EduGo Dev`
3. Selecciona scope: **`read:packages`** ← IMPORTANTE
4. Click "Generate token"
5. **Copia el token** (se muestra solo una vez, guárdalo)

### PASO 2: Login en GitHub Container Registry

```bash
# Guarda tu token (reemplaza TU_TOKEN con el token real)
export GITHUB_TOKEN="ghp_tu_token_aqui"

# Login en ghcr.io
echo $GITHUB_TOKEN | docker login ghcr.io -u TU_USUARIO_GITHUB --password-stdin
```

Deberías ver: `Login Succeeded` ✅

### PASO 3: Descargar las Imágenes Docker

```bash
# Descargar las 3 APIs (tarda ~2-3 minutos)
docker pull ghcr.io/edugogroup/edugo-api-mobile:latest
docker pull ghcr.io/edugogroup/edugo-api-administracion:latest
docker pull ghcr.io/edugogroup/edugo-worker:latest
```

### PASO 4: Clonar el Ambiente de Desarrollo

```bash
git clone git@github.com:EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment
```

### PASO 5: Configurar Variables

```bash
# Copiar archivo de ejemplo
cp .env.example .env

# Editar .env si necesitas cambiar algo (opcional)
nano .env
```

### PASO 6: Levantar Servicios

```bash
docker-compose up -d
```

### PASO 7: Verificar que Todo Funciona

```bash
# Ver estado de servicios
docker-compose ps

# Todos deben estar "Up" ✅

# Probar APIs
curl http://localhost:8081/health
curl http://localhost:8082/health
```

---

## ✅ ¡Listo para Desarrollar!

### Endpoints Disponibles

| Servicio | URL | Documentación |
|----------|-----|---------------|
| **API Mobile** | http://localhost:8081 | http://localhost:8081/swagger/index.html |
| **API Admin** | http://localhost:8082 | http://localhost:8082/swagger/index.html |
| **RabbitMQ UI** | http://localhost:15672 | Usuario: `edugo` / Pass: `edugo123` |

---

## 🔄 Actualizar a Nueva Versión

Cuando el equipo backend publique nuevas versiones:

```bash
# 1. Descargar nuevas imágenes
docker-compose pull

# 2. Recrear contenedores
docker-compose down
docker-compose up -d

# 3. Verificar
docker-compose ps
```

---

## 📱 Conectar tu App Móvil

### iOS (Simulator)

```swift
let baseURL = "http://localhost:8081"
```

### Android (Emulator)

```kotlin
val baseURL = "http://10.0.2.2:8081"  // IP especial del emulador
```

### Dispositivo Físico (mismo WiFi)

```
# Encuentra tu IP local (macOS):
ifconfig | grep "inet " | grep -v 127.0.0.1

# Usa esa IP:
http://192.168.X.X:8081
```

---

## 🔍 Comandos Útiles

### Ver Logs en Tiempo Real

```bash
# Todos los servicios
docker-compose logs -f

# Solo API Mobile
docker-compose logs -f api-mobile

# Solo Worker
docker-compose logs -f worker
```

### Reiniciar un Servicio

```bash
# Si una API se comporta raro
docker-compose restart api-mobile
```

### Detener Todo

```bash
# Detener (mantiene datos)
docker-compose stop

# Detener y eliminar contenedores (mantiene datos)
docker-compose down

# Limpiar TODO (incluyendo datos)
docker-compose down -v
```

---

## 🐛 Problemas Comunes

### "pull access denied for ghcr.io/edugogroup/..."

**Solución:**
```bash
# Hacer login nuevamente
echo $GITHUB_TOKEN | docker login ghcr.io -u TU_USUARIO --password-stdin
```

### Puerto 5432 ya está en uso

**Solución:**
```bash
# Opción 1: Detener PostgreSQL local (macOS)
brew services stop postgresql

# Opción 2: Cambiar puerto en .env
echo "POSTGRES_PORT=5433" >> .env
```

### "Container unhealthy"

**Solución:**
```bash
# Ver qué servicio falla
docker-compose ps

# Ver logs del servicio
docker-compose logs postgres
docker-compose logs mongodb
docker-compose logs rabbitmq

# Reiniciar desde cero
docker-compose down -v
docker-compose up -d
```

### Las APIs no responden

**Solución:**
```bash
# Espera ~30 segundos a que healthchecks pasen
docker-compose ps

# Cuando todos digan "healthy", intenta de nuevo
curl http://localhost:8081/health
```

---

## 📊 Imágenes Docker Disponibles

### Últimas Versiones

```bash
# Última versión (recomendada para desarrollo)
ghcr.io/edugogroup/edugo-api-mobile:latest
ghcr.io/edugogroup/edugo-api-administracion:latest
ghcr.io/edugogroup/edugo-worker:latest

# Versión estable v1.0.0
ghcr.io/edugogroup/edugo-api-mobile:v1.0.0
ghcr.io/edugogroup/edugo-api-administracion:v1.0.0
ghcr.io/edugogroup/edugo-worker:v1.0.0
```

### Versionamiento

**Sistema de tags:**
- `latest` - Última versión de main (puede cambiar frecuentemente)
- `v1.0.0` - Versión estable específica (recomendada)
- `develop` - Versión de desarrollo (cuando esté disponible)
- `abc1234` - SHA específico del commit

**Para desarrollo frontend:**
→ Usa `latest` (siempre la más reciente)

---

## 🔐 Credenciales de Desarrollo

### PostgreSQL
```
Host: localhost
Port: 5432
User: edugo
Password: edugo123
Database: edugo
```

### MongoDB
```
Host: localhost
Port: 27017
User: edugo
Password: edugo123
Database: edugo
```

### RabbitMQ
```
Host: localhost
AMQP Port: 5672
Management UI: http://localhost:15672
User: edugo
Password: edugo123
```

### JWT Secret
```
dev-secret-key-change-in-production
```

⚠️ **Estas credenciales son SOLO para desarrollo local**

---

## 📚 Documentación de APIs

### API Mobile (Puerto 8081)

**Swagger UI:** http://localhost:8081/swagger/index.html

**Endpoints principales:**
- `POST /api/v1/auth/login` - Autenticación
- `GET /api/v1/materials` - Listar materiales
- `POST /api/v1/assessments` - Crear evaluación
- `GET /api/v1/progress/:userId` - Ver progreso

### API Administración (Puerto 8082)

**Swagger UI:** http://localhost:8082/swagger/index.html

**Endpoints principales:**
- `POST /api/v1/admin/users` - Crear usuario
- `GET /api/v1/admin/schools` - Listar escuelas
- `POST /api/v1/admin/subjects` - Crear materia
- `GET /api/v1/admin/stats` - Ver estadísticas

---

## 🚀 Flujo de Trabajo Típico

### Día Normal de Desarrollo

```bash
# 1. Ir al directorio
cd edugo-dev-environment

# 2. Actualizar imágenes (si hay nuevas versiones)
docker-compose pull

# 3. Levantar servicios
docker-compose up -d

# 4. Ver logs (opcional)
docker-compose logs -f

# 5. Desarrollar tu app móvil
# Conectarte a http://localhost:8081

# 6. Al terminar el día
docker-compose stop
```

### Cuando Backend Publique Nueva Versión

Te avisarán en Slack/Teams. Entonces:

```bash
cd edugo-dev-environment

# Actualizar
docker-compose pull
docker-compose down
docker-compose up -d

# Verificar
curl http://localhost:8081/health
```

---

## 💡 Tips y Trucos

### Verificar Versión de la Imagen

```bash
# Ver qué versión estás usando
docker inspect ghcr.io/edugogroup/edugo-api-mobile:latest | grep Created

# Ver todas las versiones disponibles
gh api /orgs/EduGoGroup/packages/container/edugo-api-mobile/versions
```

### Probar una Versión Específica

```bash
# En docker-compose.yml, cambiar:
image: ghcr.io/edugogroup/edugo-api-mobile:v1.0.0
# En lugar de:
image: ghcr.io/edugogroup/edugo-api-mobile:latest
```

### Ver Espacio Usado por Docker

```bash
docker system df

# Limpiar si ocupa mucho espacio
docker system prune -a
```

---

## 🆘 Soporte

### Problemas con el Ambiente

1. **Ver logs completos:**
   ```bash
   docker-compose logs --tail=100 api-mobile
   ```

2. **Recrear desde cero:**
   ```bash
   docker-compose down -v
   docker-compose up -d
   ```

3. **Verificar Docker Desktop:**
   ```bash
   docker info
   ```

### Reportar Bugs en las APIs

- **API Mobile:** https://github.com/EduGoGroup/edugo-api-mobile/issues
- **API Admin:** https://github.com/EduGoGroup/edugo-api-administracion/issues
- **Worker:** https://github.com/EduGoGroup/edugo-worker/issues

---

## ✅ Checklist de Setup

Verifica que completaste todo:

- [ ] GitHub token generado con scope `read:packages`
- [ ] Login exitoso en ghcr.io
- [ ] Imágenes descargadas (3 en total)
- [ ] Repositorio edugo-dev-environment clonado
- [ ] Archivo .env creado
- [ ] `docker-compose up -d` ejecutado exitosamente
- [ ] Todos los servicios están "Up" (verificado con `docker-compose ps`)
- [ ] APIs responden (curl http://localhost:8081/health)
- [ ] Puedes abrir Swagger UI en el navegador

---

## 🎉 ¡Ya Puedes Desarrollar!

Con esto configurado, puedes:

✅ Desarrollar tu app móvil conectándote a las APIs
✅ Probar autenticación y endpoints
✅ Crear usuarios, materiales, evaluaciones
✅ Ver cómo el worker procesa trabajos en background
✅ Actualizar a nuevas versiones cuando backend las publique

---

**¿Dudas?** Contacta al equipo backend o crea un issue en GitHub.

**Última actualización:** 30 de Octubre, 2025
**Mantenedor:** Equipo Backend EduGo
