# Guía de Setup Detallada - EduGo Dev Environment

**Tiempo estimado:** 15-20 minutos

---

## Pre-requisitos

### 1. Docker Desktop

**macOS:**
```bash
# Descargar desde:
https://docs.docker.com/desktop/install/mac-install/

# Verificar instalación
docker --version
# Output: Docker version 24.x.x

docker info
# Debe conectar sin errores
```

### 2. GitHub Personal Access Token

**Generar token:**

1. Ir a: https://github.com/settings/tokens
2. Click **"Generate new token (classic)"**
3. Scopes necesarios:
   - ✅ `read:packages`
4. Copiar token (formato: `ghp_...`)
5. Guardar en lugar seguro

---

## Setup Paso a Paso

### Paso 1: Clonar Repositorio

```bash
git clone https://github.com/medinatello/edugo-dev-environment.git
cd edugo-dev-environment
```

### Paso 2: Ejecutar Script de Setup

```bash
./scripts/setup.sh
```

**El script hará:**
- Verificar que Docker está instalado
- Verificar que Docker está corriendo
- Crear archivo `.env` desde `.env.example`
- Pedir tu GitHub token
- Hacer login en ghcr.io
- Descargar imágenes Docker

**Cuando pida token, pegar tu GitHub Personal Access Token**

### Paso 3: Personalizar Configuración (Opcional)

```bash
# Editar variables de entorno
nano docker/.env

# Cambiar:
# - Puertos (si tienes conflictos)
# - Credenciales
# - Versiones de imágenes
# - OPENAI_API_KEY (obligatorio para worker)
```

### Paso 4: Levantar Servicios

```bash
cd docker
docker-compose up -d
```

**Tiempo de inicio:** 30-60 segundos

### Paso 5: Verificar Servicios

```bash
# Ver estado
docker-compose ps

# Todos deben mostrar "Up" y "healthy"
```

### Paso 6: Verificar Endpoints

```bash
# API Mobile
curl http://localhost:8081/health
# Esperado: {"status":"ok"}

# API Administración
curl http://localhost:8082/health
# Esperado: {"status":"ok"}

# RabbitMQ Management UI
open http://localhost:15672
# Login: edugo / edugo123
```

---

## Verificación Completa

### Checklist

- [ ] Docker Desktop instalado
- [ ] Docker corriendo
- [ ] Repositorio clonado
- [ ] Script `setup.sh` ejecutado
- [ ] Login a ghcr.io exitoso
- [ ] Archivo `.env` creado
- [ ] OPENAI_API_KEY configurada (si usas worker)
- [ ] `docker-compose up -d` ejecutado
- [ ] Todos los servicios "Up"
- [ ] API Mobile responde en :8081
- [ ] API Admin responde en :8082
- [ ] RabbitMQ UI accesible en :15672

---

## Comandos de Diagnóstico

```bash
# Ver logs en tiempo real
docker-compose logs -f

# Ver recursos usados
docker stats

# Ver redes
docker network ls | grep edugo

# Ver volúmenes
docker volume ls | grep edugo

# Inspeccionar servicio
docker-compose exec api-mobile sh
```

---

## Actualización de Imágenes

```bash
# Cuando hay nuevas versiones en ghcr.io
cd ..
./scripts/update-images.sh

cd docker
docker-compose down
docker-compose up -d
```

---

## Limpiar Ambiente

```bash
# Usar script interactivo
cd ..
./scripts/cleanup.sh

# Manual: Detener y eliminar todo
cd docker
docker-compose down -v  # -v elimina volúmenes (datos)
```

---

Ver más en: [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
