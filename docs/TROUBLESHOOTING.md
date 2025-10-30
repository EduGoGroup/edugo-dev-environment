# Troubleshooting - EduGo Dev Environment

---

## Docker

### ❌ "Cannot connect to Docker daemon"

**Causa:** Docker Desktop no está corriendo

**Solución:**
```bash
# Iniciar Docker Desktop
open -a Docker

# Esperar a que aparezca el ícono en la barra de menú
# Verificar
docker ps
```

---

## GitHub Container Registry

### ❌ "pull access denied for ghcr.io/medinatello/..."

**Causa:** No has hecho login o token expiró

**Solución:**
```bash
# Login con tu GitHub token
echo "TU_GITHUB_TOKEN" | docker login ghcr.io -u medinatello --password-stdin

# Verificar
docker pull ghcr.io/medinatello/api-mobile:latest
```

### ❌ "manifest unknown" al hacer pull

**Causa:** La imagen todavía no existe en ghcr.io (no se ha hecho deploy)

**Solución:**
- Esperar a que el pipeline de GitLab construya y suba las imágenes
- O usar versión específica que exista: `API_MOBILE_VERSION=develop`

---

## Puertos

### ❌ "Port 5432 already in use"

**Causa:** PostgreSQL local ya está corriendo

**Solución opción 1:**
```bash
# Detener PostgreSQL local
brew services stop postgresql
```

**Solución opción 2:**
```bash
# Cambiar puerto en docker/.env
echo "POSTGRES_PORT=5433" >> docker/.env

# Reiniciar
docker-compose down
docker-compose up -d
```

### ❌ "Port 27017 already in use"

**Causa:** MongoDB local corriendo

**Solución:**
```bash
brew services stop mongodb-community
```

---

## Servicios

### ❌ Servicio no inicia (status "unhealthy")

**Solución:**
```bash
# Ver logs del servicio
docker-compose logs nombre-servicio

# Reiniciar servicio específico
docker-compose restart nombre-servicio

# Recrear desde cero
docker-compose down -v
docker-compose up -d
```

### ❌ "worker keeps restarting"

**Causa:** Probablemente falta OPENAI_API_KEY

**Solución:**
```bash
# Verificar configuración
grep OPENAI_API_KEY docker/.env

# Editar y agregar API key válida
nano docker/.env

# Reiniciar worker
docker-compose restart worker

# Ver logs
docker-compose logs -f worker
```

---

## RabbitMQ

### ❌ Worker no consume mensajes

**Diagnóstico:**
```bash
# Ver logs del worker
docker-compose logs -f worker

# Verificar RabbitMQ UI
open http://localhost:15672

# Login: edugo / edugo123
# Verificar que existen las queues
```

**Solución:**
- Verificar que RabbitMQ está "healthy"
- Verificar credenciales en docker/.env
- Reiniciar worker

---

## Volúmenes

### ❌ Datos persisten después de `docker-compose down`

**Causa:** Los volúmenes no se eliminan por defecto

**Solución:**
```bash
# Eliminar volúmenes explícitamente
docker-compose down -v

# Ver volúmenes
docker volume ls | grep edugo

# Eliminar manualmente
docker volume rm edugo-postgres-data
```

---

## Rendimiento

### ❌ Servicios muy lentos

**Solución:**
```bash
# Ver uso de recursos
docker stats

# Aumentar recursos de Docker Desktop:
# Docker Desktop > Settings > Resources
# - CPUs: 4+
# - Memory: 8GB+
```

---

## Imágenes

### ❌ "No space left on device"

**Solución:**
```bash
# Limpiar imágenes no usadas
docker image prune -a -f

# Ver espacio usado
docker system df
```

---

Para más ayuda, consulta: [SETUP.md](SETUP.md)
