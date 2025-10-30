# EduGo - Ambiente de Desarrollo Local

**Versión:** 1.0.0
**Última actualización:** 30 de Octubre, 2025

Este repositorio contiene todo lo necesario para ejecutar **EduGo** localmente usando Docker Compose.

---

## 🚀 Inicio Rápido

### Pre-requisitos

- ✅ [Docker Desktop](https://docs.docker.com/desktop/install/mac-install/) instalado y corriendo
- ✅ Git instalado
- ✅ Acceso a GitHub Container Registry (ghcr.io)
- ✅ GitHub Personal Access Token con scope `read:packages`

### Setup Inicial (Primera vez)

```bash
# 1. Clonar este repositorio
git clone https://github.com/medinatello/edugo-dev-environment.git
cd edugo-dev-environment

# 2. Ejecutar script de setup
./scripts/setup.sh
# Te pedirá tu GitHub Personal Access Token

# 3. Levantar servicios
cd docker
docker-compose up -d

# 4. Verificar que todo está corriendo
docker-compose ps
# Todos los servicios deben mostrar "Up"
```

---

## 📦 Servicios Incluidos

| Servicio | Puerto Local | URL | Estado |
|----------|-------------|-----|--------|
| **API Mobile** | 8081 | http://localhost:8081 | Backend REST API |
| **API Administración** | 8082 | http://localhost:8082 | Backend Admin Panel |
| **Worker** | - | (background) | Procesador de PDFs |
| **PostgreSQL** | 5432 | localhost:5432 | Base de datos relacional |
| **MongoDB** | 27017 | localhost:27017 | Base de datos NoSQL |
| **RabbitMQ** | 5672, 15672 | http://localhost:15672 | Message Queue + UI |

### Endpoints de Health Check

```bash
# API Mobile
curl http://localhost:8081/health

# API Administración
curl http://localhost:8082/health

# RabbitMQ Management UI
open http://localhost:15672
# Usuario: edugo
# Password: edugo123
```

---

## 🔄 Comandos Útiles

### Ver logs de todos los servicios

```bash
cd docker
docker-compose logs -f
```

### Ver logs de un servicio específico

```bash
docker-compose logs -f api-mobile
docker-compose logs -f worker
docker-compose logs -f postgres
```

### Reiniciar un servicio

```bash
docker-compose restart api-mobile
```

### Detener servicios (mantiene datos)

```bash
docker-compose stop
```

### Detener y eliminar contenedores (mantiene datos)

```bash
docker-compose down
```

### Actualizar a última versión de las imágenes

```bash
# Desde raíz de edugo-dev-environment
./scripts/update-images.sh

# Luego reiniciar
cd docker
docker-compose down
docker-compose up -d
```

### Limpiar ambiente completo

```bash
# Desde raíz de edugo-dev-environment
./scripts/cleanup.sh

# El script preguntará si deseas:
# - Eliminar volúmenes (datos de BD)
# - Limpiar imágenes no usadas
# - Eliminar imágenes de EduGo
```

---

## 🔐 Credenciales por Defecto (Desarrollo)

### PostgreSQL
- **Usuario:** `edugo`
- **Password:** `edugo123`
- **Database:** `edugo`
- **Puerto:** 5432

### MongoDB
- **Usuario:** `edugo`
- **Password:** `edugo123`
- **Database:** `edugo`
- **Puerto:** 27017

### RabbitMQ
- **Usuario:** `edugo`
- **Password:** `edugo123`
- **Puerto AMQP:** 5672
- **Puerto Management UI:** 15672
- **Management UI:** http://localhost:15672

### JWT Secret (Desarrollo)
- **Secret:** `dev-secret-key-change-in-production`

---

## ⚙️ Configuración Personalizada

### Editar variables de entorno

```bash
# Copiar ejemplo si no existe
cp docker/.env.example docker/.env

# Editar configuración
nano docker/.env
```

### Variables Importantes

| Variable | Descripción | Default |
|----------|-------------|---------|
| `POSTGRES_PASSWORD` | Password de PostgreSQL | `edugo123` |
| `MONGO_PASSWORD` | Password de MongoDB | `edugo123` |
| `RABBITMQ_PASSWORD` | Password de RabbitMQ | `edugo123` |
| `JWT_SECRET` | Secret para tokens JWT | `dev-secret-key...` |
| `OPENAI_API_KEY` | API Key de OpenAI (para worker) | `sk-...` |
| `API_MOBILE_VERSION` | Versión de imagen Docker | `latest` |
| `API_ADMIN_VERSION` | Versión de imagen Docker | `latest` |
| `WORKER_VERSION` | Versión de imagen Docker | `latest` |

**Ver archivo completo:** [`docker/.env.example`](docker/.env.example)

---

## 🐳 Versiones de Imágenes

Por defecto, se usan las imágenes `latest` de cada servicio. Puedes usar versiones específicas:

```bash
# En docker/.env
API_MOBILE_VERSION=develop        # Usar versión de develop
API_MOBILE_VERSION=a1b2c3d         # Usar SHA específico
API_MOBILE_VERSION=v1.2.3          # Usar tag de versión
```

---

## 🔍 Troubleshooting

### Problema: "Cannot connect to Docker daemon"

**Solución:**
```bash
# Verificar que Docker Desktop está corriendo
open -a Docker

# Esperar a que inicie (ícono en la barra de menú)
# Reintentar: docker ps
```

### Problema: "pull access denied for ghcr.io/medinatello/api-mobile"

**Solución:**
```bash
# Login nuevamente con tu GitHub token
echo "TU_GITHUB_TOKEN" | docker login ghcr.io -u medinatello --password-stdin

# Verificar login
docker info | grep ghcr.io
```

### Problema: "Port 5432 already in use"

**Solución:**
```bash
# Opción 1: Detener PostgreSQL local
brew services stop postgresql

# Opción 2: Cambiar puerto en docker/.env
echo "POSTGRES_PORT=5433" >> docker/.env
```

### Problema: "Servicios no arrancan (unhealthy)"

**Solución:**
```bash
# Ver logs del servicio problemático
cd docker
docker-compose logs postgres
docker-compose logs mongodb
docker-compose logs rabbitmq

# Reiniciar desde cero
docker-compose down -v  # Elimina volúmenes
docker-compose up -d    # Recrea todo
```

### Problema: "Worker no procesa mensajes"

**Solución:**
1. Verificar RabbitMQ:
   ```bash
   docker-compose logs -f rabbitmq
   open http://localhost:15672  # Ver UI
   ```

2. Verificar configuración de OPENAI_API_KEY:
   ```bash
   grep OPENAI_API_KEY docker/.env
   ```

3. Ver logs del worker:
   ```bash
   docker-compose logs -f worker
   ```

---

## 📚 Documentación Adicional

- 📖 [Configuración Detallada](docs/SETUP.md)
- 📖 [Variables de Entorno](docs/VARIABLES.md)
- 📖 [Troubleshooting Completo](docs/TROUBLESHOOTING.md)

---

## ⚠️ Notas Importantes

- ⚠️ **Este ambiente es SOLO para desarrollo local**
- ⚠️ **NO usar estas credenciales en producción**
- ⚠️ Las imágenes se descargan de GitHub Container Registry (ghcr.io)
- ⚠️ Necesitas estar autenticado en ghcr.io para descargar imágenes
- ⚠️ El worker requiere OPENAI_API_KEY válida para funcionar

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────┐
│                  GITHUB CONTAINER REGISTRY               │
│                     (ghcr.io/medinatello)                │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ api-mobile   │  │ api-admin    │  │   worker     │  │
│  │   :latest    │  │   :latest    │  │   :latest    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────┬────────────┬────────────┬──────────────────┘
             │            │            │
             │  docker pull (en setup.sh)
             ↓            ↓            ↓
┌────────────────────────────────────────────────────────┐
│           DOCKER COMPOSE (tu Mac local)                │
│                                                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│  │PostgreSQL│  │ MongoDB  │  │ RabbitMQ │            │
│  │  :5432   │  │  :27017  │  │:5672/15672           │
│  └──────────┘  └──────────┘  └──────────┘            │
│                                                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│  │API Mobile│  │API Admin │  │  Worker  │            │
│  │  :8081   │  │  :8082   │  │(background)          │
│  └──────────┘  └──────────┘  └──────────┘            │
└────────────────────────────────────────────────────────┘
```

---

## 📞 Soporte

Si encuentras problemas:

1. Revisa la documentación en [`docs/`](docs/)
2. Verifica logs: `docker-compose logs -f`
3. Consulta troubleshooting: [`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md)

---

## 📝 Licencia

Privado - EduGo © 2025

---

**Última actualización:** 30 de Octubre, 2025
**Mantenedor:** Equipo EduGo
