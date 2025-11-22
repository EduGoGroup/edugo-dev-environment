# EduGo - Ambiente de Desarrollo Local

**Versión:** 1.0.0
**Última actualización:** 18 de Noviembre, 2025

Este repositorio contiene todo lo necesario para ejecutar **EduGo** localmente usando Docker Compose.

---

## 📖 Guías Disponibles

| Guía | Descripción | Cuándo Usar |
|------|-------------|-------------|
| **[🚀 Quick Start](docker/QUICK_START.md)** | Inicio rápido (5 min) | Primera vez, quiero empezar YA |
| **[📝 Ejemplo End-to-End](docs/EXAMPLE.md)** | Guía paso a paso completa | Quiero ver todo el flujo en detalle |
| **[📚 Guía Completa](docker/README.md)** | Documentación detallada | Necesito entender todo el sistema |
| **[✅ Reporte de Validación](docker/RESULTADO_VALIDACION.md)** | Estado y troubleshooting | Tengo problemas técnicos |

**¿Primera vez usando este proyecto?** → Comienza con [Quick Start](docker/QUICK_START.md) o [Ejemplo End-to-End](docs/EXAMPLE.md)

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
git clone git@github.com:EduGoGroup/edugo-dev-environment.git
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

Los siguientes servicios se levantarán automáticamente:
- **API Mobile** (8081)
- **API Administración** (8082)
- **Worker** (background)
- **PostgreSQL** (5432)
- **MongoDB** (27017)
- **RabbitMQ** (5672, 15672)
- **Migrator** (ejecuta migraciones automáticas)
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

Por defecto, se usan las imágenes `latest` de cada servicio desde GitHub Container Registry.

**Imágenes disponibles:**
- `ghcr.io/edugogroup/edugo-api-mobile`
- `ghcr.io/edugogroup/edugo-api-administracion`
- `ghcr.io/edugogroup/edugo-worker`

**Usar versiones específicas:**

```bash
# En docker/.env
API_MOBILE_VERSION=1.0.0          # Versión específica
API_MOBILE_VERSION=1.0            # Último patch de 1.0
API_MOBILE_VERSION=1              # Último minor de 1.x
API_MOBILE_VERSION=latest         # Última versión publicada

# También puedes usar:
API_ADMIN_VERSION=1.0.0
WORKER_VERSION=1.0.0
```

**Ver versiones disponibles:**
- https://github.com/orgs/EduGoGroup/packages

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

### Problema: "pull access denied for ghcr.io/edugogroup/api-mobile"

**Solución:**
```bash
# Login nuevamente con tu GitHub token
echo "TU_GITHUB_TOKEN" | docker login ghcr.io -u TU_USUARIO_GITHUB --password-stdin

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

### Problema: "Error de conexión a base de datos"

**Error:**
```
dial tcp [::1]:5432: connect: connection refused
```

**Solución:**
```bash
# Verificar que PostgreSQL está corriendo
docker-compose ps postgres

# Si no está corriendo, iniciarlo
docker-compose up -d postgres

# Verificar logs
docker-compose logs postgres

# Probar conexión manual
docker exec -it edugo-dev-environment-postgres-1 psql -U edugo -d edugo -c "SELECT 1;"
```

### Problema: "Imágenes Docker no se descargan"

**Error:**
```
Error response from daemon: pull access denied for ghcr.io/edugogroup/...
```

**Solución:**
```bash
# 1. Verificar autenticación
docker login ghcr.io

# 2. Verificar token tiene permisos read:packages
echo $GITHUB_TOKEN | docker login ghcr.io -u TU_USUARIO --password-stdin

# 3. Si el problema persiste, re-ejecutar setup
./scripts/setup.sh

# 4. Verificar que puedes ver el paquete en GitHub
open https://github.com/orgs/EduGoGroup/packages
```

### Problema: "Migraciones no se ejecutan"

**Síntomas:**
- Las tablas no existen en PostgreSQL
- Error "relation does not exist"

**Solución:**
```bash
# Verificar logs del migrator
docker-compose logs migrator

# Ejecutar migraciones manualmente
docker-compose run --rm migrator

# Verificar tablas creadas
docker exec -it edugo-dev-environment-postgres-1 psql -U edugo -d edugo -c "\dt"

# Si sigue fallando, limpiar y reiniciar
docker-compose down -v
docker-compose up -d
```

### Problema: "Espacio en disco lleno"

**Error:**
```
no space left on device
```

**Solución:**
```bash
# Ver uso de espacio de Docker
docker system df

# Limpiar contenedores detenidos
docker container prune

# Limpiar imágenes sin usar
docker image prune -a

# Limpiar volúmenes sin usar (⚠️ borra datos)
docker volume prune

# Limpieza completa (⚠️ borra todo)
docker system prune -a --volumes
```

### Problema: "API responde 500 Internal Server Error"

**Solución:**
```bash
# 1. Ver logs de la API
docker-compose logs -f api-mobile

# 2. Verificar variables de entorno
docker-compose exec api-mobile env | grep -E "DATABASE|MONGO|RABBITMQ"

# 3. Verificar conectividad a servicios
docker-compose exec api-mobile ping -c 2 postgres
docker-compose exec api-mobile ping -c 2 mongodb
docker-compose exec api-mobile ping -c 2 rabbitmq

# 4. Reiniciar API
docker-compose restart api-mobile
```

---

## 📚 Documentación Adicional

### Docker Compose
- 🚀 **[Quick Start](docker/QUICK_START.md)** ← EMPIEZA AQUÍ
- 📚 [Guía Completa Docker](docker/README.md) - 3 archivos docker-compose disponibles
- ✅ [Reporte de Validación](docker/RESULTADO_VALIDACION.md) - Estado actual y soluciones

### Documentación del Proyecto
- 📖 [Documentación Dev Environment](docs/dev-environment/) - Especificaciones técnicas
- 📖 [Templates de Workflow](docs/workflow-templates/) - Metodología de trabajo

---

## ⚠️ Notas Importantes

- ⚠️ **Este ambiente es SOLO para desarrollo local**
- ⚠️ **NO usar estas credenciales en producción**
- ⚠️ Las imágenes se descargan de GitHub Container Registry (ghcr.io)
- ⚠️ Necesitas estar autenticado en ghcr.io para descargar imágenes
- ⚠️ El worker requiere OPENAI_API_KEY válida para funcionar

---

## 🏗️ Arquitectura

### Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────┐
│                  GITHUB CONTAINER REGISTRY               │
│                     (ghcr.io/edugogroup)                 │
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

### Flujo de Datos

```
┌──────────────┐
│ App Móvil    │
│ (Flutter)    │
└──────┬───────┘
       │ HTTP REST
       ↓
┌──────────────┐      ┌──────────────┐
│ API Mobile   │─────→│ PostgreSQL   │
│ (Go)         │←─────│ (Datos)      │
└──────┬───────┘      └──────────────┘
       │
       │ Publica mensaje
       ↓
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│ RabbitMQ     │─────→│ Worker       │─────→│ MongoDB      │
│ (Queue)      │      │ (Go + AI)    │      │ (PDFs)       │
└──────────────┘      └──────────────┘      └──────────────┘
       ↑
       │ Consume mensajes
       │
┌──────────────┐      ┌──────────────┐
│ API Admin    │─────→│ PostgreSQL   │
│ (Go)         │←─────│ (Config)     │
└──────────────┘      └──────────────┘
       ↑
       │ HTTP REST
       ┌──────────────┐
       │ Panel Admin  │
       │ (Web)        │
       └──────────────┘
```

### Componentes Detallados

| Componente | Tecnología | Propósito | Datos Persistentes |
|------------|------------|-----------|-------------------|
| **API Mobile** | Go 1.21+ | Backend para app móvil | PostgreSQL |
| **API Admin** | Go 1.21+ | Backend para panel admin | PostgreSQL |
| **Worker** | Go 1.21+ | Procesamiento asíncrono PDFs | MongoDB |
| **PostgreSQL** | PostgreSQL 15 | BD relacional principal | Volumen Docker |
| **MongoDB** | MongoDB 7.0 | BD documentos (PDFs) | Volumen Docker |
| **RabbitMQ** | RabbitMQ 3.12 | Cola de mensajes | Volumen Docker |
| **Migrator** | Go (custom) | Migraciones automáticas | N/A (init) |

---

## 🤔 ¿Por Qué Este Proyecto NO Tiene CI/CD?

**Pregunta común:** ¿Por qué no hay workflows de GitHub Actions en este repositorio?

**Respuesta:** Este proyecto **intencionalmente NO tiene CI/CD** porque es un repositorio de **configuración**, no de **código**.

### Análisis Técnico

| Aspecto | Este Proyecto | Proyectos con CI/CD |
|---------|---------------|---------------------|
| **Tipo** | Configuración Docker | Código fuente (Go/Python/etc) |
| **Contenido** | docker-compose.yml, scripts | Aplicaciones con lógica |
| **Tests** | ❌ No aplica | ✅ Tests unitarios/integración |
| **Builds** | ❌ No genera artefactos | ✅ Binarios, imágenes Docker |
| **Despliegue** | ❌ Solo para desarrollo local | ✅ Staging/Production |
| **Validación** | ✅ Local (instantánea) | ✅ CI/CD (distribuido) |

### Razones Específicas

1. **No hay código que testear**
   - Los archivos YAML no tienen tests unitarios
   - Los scripts bash son utilidades simples
   - No hay lógica de negocio

2. **La validación es mejor localmente**
   - `docker-compose config` valida sintaxis al instante
   - `./scripts/validate.sh` ejecuta en segundos
   - Feedback inmediato vs esperar queue de CI

3. **No hay despliegues automáticos**
   - Este ambiente es solo para desarrollo local
   - No se despliega a staging ni producción
   - No se publican imágenes Docker

4. **Costo vs Beneficio**
   ```
   Costo de CI/CD:
   - ~50-100 minutos/mes de GitHub Actions
   - Mantenimiento de workflows
   - Complejidad adicional
   
   Beneficio:
   - Validar sintaxis YAML (se hace mejor local)
   - ¿?
   
   Conclusión: Costo > Beneficio
   ```

### Enfoque Alternativo: Validación Local

En lugar de CI/CD completo, usamos **herramientas locales**:

#### 1. Script de Validación

```bash
./scripts/validate.sh
```

**Qué hace:**
- ✅ Valida sintaxis de todos los docker-compose
- ✅ Lista servicios, volúmenes y puertos
- ✅ Verifica variables de entorno
- ⚡ Resultado en segundos

#### 2. Pre-commit Hooks (Opcional)

```bash
git config core.hooksPath .githooks
```

**Qué hace:**
- ✅ Valida automáticamente antes de commit
- ✅ Previene commit de archivos `.env`
- ✅ Asegura permisos correctos en scripts
- ⚡ Bloquea commits con errores

#### 3. Validación Nativa de Docker

```bash
cd docker
docker-compose config
```

**Qué hace:**
- ✅ Valida sintaxis YAML
- ✅ Expande variables de entorno
- ✅ Muestra configuración final
- ⚡ Herramienta oficial de Docker

### Filosofía

> **"No uses CI/CD para todo. Úsalo solo donde agregue valor."**

**CI/CD es excelente para:**
- ✅ `edugo-api-mobile` - Tests, builds, deploys
- ✅ `edugo-api-administracion` - Tests, builds, deploys
- ✅ `edugo-worker` - Tests, builds, deploys
- ✅ `edugo-shared` - Tests, releases de paquetes

**CI/CD NO es necesario para:**
- ❌ Repos de configuración (este proyecto)
- ❌ Repos de documentación pura
- ❌ Repos de scripts de utilidad

### Comparación: Con CI/CD vs Sin CI/CD

#### Opción A: CON CI/CD (No Recomendado)

**Workflows que podríamos crear:**
```yaml
# .github/workflows/validate.yml
name: Validate
on: [push, pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate YAML
        run: docker-compose config
```

**Problemas:**
- ❌ Esperar 2-5 minutos por resultado
- ❌ Consumir minutos de GitHub Actions
- ❌ Validación que se hace mejor local
- ❌ Mantenimiento de workflow

#### Opción B: SIN CI/CD (Recomendado) ✅

**Validación local:**
```bash
./scripts/validate.sh  # 2 segundos
```

**Beneficios:**
- ✅ Feedback instantáneo
- ✅ Cero minutos de GitHub Actions
- ✅ Menos complejidad
- ✅ Mejor experiencia de desarrollo

### Casos Especiales

**¿Cuándo SÍ agregar CI/CD a este proyecto?**

Solo si cambia su propósito:

1. **Si genera imágenes Docker propias**
   - Actualmente: Pull de `ghcr.io/edugogroup/*`
   - Si cambia a build local → Sí CI/CD

2. **Si se despliega a cloud**
   - Actualmente: Solo desarrollo local
   - Si se despliega a AWS/GCP → Sí CI/CD

3. **Si tiene tests de integración complejos**
   - Actualmente: No hay tests
   - Si se agregan tests E2E → Considerar CI/CD

### Decisión Documentada

**Fecha:** 22 de Noviembre, 2025  
**Decisión:** NO implementar CI/CD en este repositorio  
**Razón:** Es configuración, no código  
**Alternativa:** Validación local con scripts  
**Revisar decisión:** Solo si el propósito del repo cambia  

### Referencias

Para más contexto sobre esta decisión:
- Ver análisis completo: [docs/cicd/README.md](docs/cicd/README.md)
- Ver plan de implementación: [docs/cicd/sprints/SPRINT-3-TASKS.md](docs/cicd/sprints/SPRINT-3-TASKS.md)

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

## 🚀 Perfiles Disponibles (Opcional)

Si deseas usar perfiles específicos para levantamientos parciales, puedes ejecutar:

```bash
# Solo bases de datos (sin APIs ni worker)
cd docker
docker-compose --profile db-only up -d

# APIs sin worker
docker-compose --profile api-only up -d

# Solo worker
docker-compose --profile worker-only up -d

# Solo Mobile API
docker-compose --profile mobile-only up -d

# Solo Admin API
docker-compose --profile admin-only up -d
```

### Perfiles Disponibles

| Profile | Servicios | Uso Recomendado |
|---------|-----------|-----------------|
| (sin profile) | Todos los servicios | Desarrollo completo (DEFAULT) |
| `db-only` | PostgreSQL + MongoDB + RabbitMQ + Migrator | Testing de migraciones |
| `api-only` | DBs + APIs + Migrator | Desarrollo de APIs |
| `mobile-only` | DBs + API Mobile + Migrator | App móvil |
| `admin-only` | DBs + API Admin + Migrator | Panel admin |
| `worker-only` | DBs + Worker + Migrator | Testing de workers |

Ver [docs/PROFILES.md](docs/PROFILES.md) para más detalles.

## 🛑 Detener Servicios

```bash
# Detener todo
./scripts/stop.sh

# Detener perfil específico
./scripts/stop.sh --profile db-only

# Eliminar volúmenes
./scripts/stop.sh --volumes
```
