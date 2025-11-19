# Sprint 3 - Documentación y Validación Opcional

**Proyecto:** edugo-dev-environment  
**Sprint:** 3 (Mejoras Opcionales)  
**Duración:** 2-3 horas  
**Prioridad:** BAJA (solo si quieres mejorar)  
**Fecha:** 19 de Noviembre, 2025

---

## 🎯 Objetivo del Sprint

Mejorar la experiencia del desarrollador mediante:
1. Mejor documentación
2. Scripts de validación local
3. Pre-commit hooks opcionales

**⚠️ IMPORTANTE:** Este sprint es OPCIONAL. Si el proyecto ya funciona bien, NO es necesario ejecutarlo.

---

## 📋 Checklist General

- [ ] Tarea 3.1: Mejorar README.md (30-45 min)
- [ ] Tarea 3.2: Script de validación YAML (30 min)
- [ ] Tarea 3.3: Pre-commit hook opcional (30 min)
- [ ] Tarea 3.4: Documentar decisión de NO CI/CD (15 min)
- [ ] Tarea 3.5: Crear ejemplo end-to-end (30-45 min)

**Total:** 2-3 horas

---

## Tarea 3.1: Mejorar README.md

**⏱️ Tiempo estimado:** 30-45 minutos  
**🎯 Objetivo:** README.md más completo y útil para nuevos desarrolladores  
**📍 Ubicación:** `/repos-separados/edugo-dev-environment/README.md`

### Subtareas

#### 1.1 Agregar Sección de Requisitos Previos

**Archivo:** `README.md`

**Agregar al inicio:**

```markdown
## 📋 Requisitos Previos

Antes de usar este entorno, asegúrate de tener instalado:

### Obligatorio
- **Docker:** v20.10 o superior
- **Docker Compose:** v2.0 o superior

### Opcional
- **Git:** v2.30 o superior
- **Make:** Para usar shortcuts (solo Linux/Mac)

### Verificar Instalación

```bash
# Verificar Docker
docker --version
# Esperado: Docker version 20.10.x o superior

# Verificar Docker Compose
docker-compose --version
# Esperado: Docker Compose version 2.x.x o superior

# Verificar que Docker está corriendo
docker ps
# Esperado: Lista de contenedores (puede estar vacía)
```

### Instalación de Requisitos

**macOS:**
```bash
# Instalar Docker Desktop
brew install --cask docker
```

**Linux (Ubuntu/Debian):**
```bash
# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Instalar Docker Compose
sudo apt-get install docker-compose-plugin
```

**Windows:**
- Descargar Docker Desktop desde: https://www.docker.com/products/docker-desktop
```

**Checklist:**
- [ ] Sección de requisitos agregada
- [ ] Comandos de verificación incluidos
- [ ] Instrucciones por SO incluidas

---

#### 1.2 Agregar Sección de Troubleshooting

**Archivo:** `README.md`

**Agregar al final:**

```markdown
## 🔧 Solución de Problemas Comunes

### Problema 1: Puerto Ya en Uso

**Error:**
```
ERROR: for postgres  Cannot start service postgres: 
Bind for 0.0.0.0:5432 failed: port is already allocated
```

**Solución:**
```bash
# Opción A: Detener el servicio local
sudo systemctl stop postgresql  # Linux
brew services stop postgresql    # macOS

# Opción B: Cambiar puerto en docker-compose.yml
# Editar: "5432:5432" → "5433:5432"

# Opción C: Identificar y matar proceso
lsof -ti:5432 | xargs kill -9
```

---

### Problema 2: Permisos de Volúmenes

**Error:**
```
permission denied while trying to connect to the Docker daemon socket
```

**Solución:**
```bash
# Linux: Agregar usuario a grupo docker
sudo usermod -aG docker $USER
newgrp docker

# macOS: Reiniciar Docker Desktop
```

---

### Problema 3: Contenedores No Inician

**Error:**
```
ERROR: Service 'postgres' failed to build: 
The command '/bin/sh -c apt-get update' returned a non-zero code
```

**Solución:**
```bash
# Limpiar caché de Docker
docker system prune -a

# Reconstruir sin caché
docker-compose build --no-cache
docker-compose up -d
```

---

### Problema 4: Variables de Entorno No Cargadas

**Error:**
```
POSTGRES_PASSWORD is required
```

**Solución:**
```bash
# Crear .env desde ejemplo
cp .env.example .env

# Editar .env con tus valores
nano .env
```

---

### Problema 5: Logs de Errores

**Ver logs de un servicio específico:**
```bash
docker-compose logs postgres
docker-compose logs mongodb
docker-compose logs rabbitmq
```

**Ver logs en tiempo real:**
```bash
docker-compose logs -f
```

**Ver últimas 50 líneas:**
```bash
docker-compose logs --tail=50
```
```

**Checklist:**
- [ ] Sección de troubleshooting agregada
- [ ] Al menos 5 problemas comunes cubiertos
- [ ] Soluciones con comandos copy-paste

---

#### 1.3 Agregar Sección de Arquitectura

**Archivo:** `README.md`

**Agregar antes de "Uso":**

```markdown
## 🏗️ Arquitectura del Entorno

Este entorno proporciona la infraestructura completa para desarrollo:

```
┌─────────────────────────────────────────────────┐
│         edugo-dev-environment                    │
│                                                  │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  │
│  │PostgreSQL │  │ MongoDB   │  │ RabbitMQ  │  │
│  │   :5432   │  │  :27017   │  │  :5672    │  │
│  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘  │
│        │              │              │         │
└────────┼──────────────┼──────────────┼─────────┘
         │              │              │
    ┌────▼──────────────▼──────────────▼────┐
    │     APIs y Worker de EduGo             │
    │  (api-mobile, api-admin, worker)       │
    └────────────────────────────────────────┘
```

### Componentes

| Servicio | Puerto | Propósito | URL de Acceso |
|----------|--------|-----------|---------------|
| PostgreSQL | 5432 | Base de datos relacional | `postgresql://localhost:5432/edugo` |
| MongoDB | 27017 | Base de datos documentos | `mongodb://localhost:27017/edugo` |
| RabbitMQ | 5672 | Sistema de mensajería | `amqp://localhost:5672` |
| RabbitMQ UI | 15672 | Interfaz web RabbitMQ | http://localhost:15672 |

### Persistencia

Los datos se persisten en volúmenes Docker:
- `edugo_postgres_data` - Datos de PostgreSQL
- `edugo_mongo_data` - Datos de MongoDB
- `edugo_rabbitmq_data` - Datos de RabbitMQ

**⚠️ IMPORTANTE:** Los datos persisten entre reinicios de contenedores, pero se pierden si eliminas los volúmenes.
```

**Checklist:**
- [ ] Diagrama de arquitectura agregado
- [ ] Tabla de componentes agregada
- [ ] Nota sobre persistencia incluida

---

### Validación de Tarea 3.1

**Checklist final:**
- [ ] README.md tiene sección de requisitos previos
- [ ] README.md tiene troubleshooting con 5+ problemas
- [ ] README.md tiene diagrama de arquitectura
- [ ] Formato markdown correcto
- [ ] Enlaces funcionan correctamente

**Comando de validación:**
```bash
# Ver preview del README
cat README.md | grep "##"  # Ver todas las secciones

# Validar markdown
npx markdownlint README.md  # Si tienes Node.js
```

---

## Tarea 3.2: Script de Validación YAML

**⏱️ Tiempo estimado:** 30 minutos  
**🎯 Objetivo:** Script para validar sintaxis de docker-compose.yml localmente  
**📍 Ubicación:** `/repos-separados/edugo-dev-environment/scripts/`

### Subtareas

#### 2.1 Crear Script de Validación

**Archivo:** `scripts/validate.sh`

**Contenido:**

```bash
#!/bin/bash

# edugo-dev-environment - Script de Validación
# Valida sintaxis de docker-compose.yml

set -e  # Exit on error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Validando configuración de Docker Compose..."
echo ""

# Verificar que docker-compose está instalado
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ ERROR: docker-compose no está instalado${NC}"
    echo "Instalar desde: https://docs.docker.com/compose/install/"
    exit 1
fi

echo -e "${GREEN}✅ docker-compose instalado${NC}"

# Verificar que docker-compose.yml existe
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}❌ ERROR: docker-compose.yml no encontrado${NC}"
    echo "Ejecutar desde el directorio raíz del proyecto"
    exit 1
fi

echo -e "${GREEN}✅ docker-compose.yml encontrado${NC}"

# Validar sintaxis YAML
echo ""
echo "📝 Validando sintaxis YAML..."
if docker-compose -f docker-compose.yml config > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Sintaxis YAML válida${NC}"
else
    echo -e "${RED}❌ ERROR: Sintaxis YAML inválida${NC}"
    echo ""
    docker-compose -f docker-compose.yml config
    exit 1
fi

# Verificar servicios definidos
echo ""
echo "🔍 Servicios encontrados:"
docker-compose config --services | while read service; do
    echo -e "  ${GREEN}✓${NC} $service"
done

# Verificar volúmenes definidos
echo ""
echo "💾 Volúmenes encontrados:"
docker-compose config --volumes | while read volume; do
    echo -e "  ${GREEN}✓${NC} $volume"
done

# Verificar puertos expuestos
echo ""
echo "🌐 Puertos expuestos:"
docker-compose config | grep -A 1 "ports:" | grep -o "[0-9]*:[0-9]*" | sort -u | while read port; do
    echo -e "  ${GREEN}✓${NC} $port"
done

# Verificar que .env existe (si es requerido)
if [ -f ".env.example" ] && [ ! -f ".env" ]; then
    echo ""
    echo -e "${YELLOW}⚠️  ADVERTENCIA: .env no existe${NC}"
    echo "Crear desde .env.example:"
    echo "  cp .env.example .env"
fi

echo ""
echo -e "${GREEN}✅ Validación completada exitosamente${NC}"
echo ""
echo "Próximo paso:"
echo "  docker-compose up -d"
```

**Checklist:**
- [ ] Script creado en `scripts/validate.sh`
- [ ] Permisos de ejecución agregados
- [ ] Script probado localmente

**Comandos:**
```bash
# Crear y hacer ejecutable
touch scripts/validate.sh
chmod +x scripts/validate.sh

# Probar
./scripts/validate.sh
```

---

#### 2.2 Agregar Documentación del Script

**Archivo:** `scripts/README.md`

**Crear/Actualizar con:**

```markdown
# Scripts de Utilidad

## validate.sh

Valida la sintaxis de `docker-compose.yml` sin levantar contenedores.

### Uso

```bash
./scripts/validate.sh
```

### ¿Qué Valida?

- ✅ Sintaxis YAML correcta
- ✅ Servicios definidos
- ✅ Volúmenes definidos
- ✅ Puertos expuestos
- ⚠️ Existencia de .env

### Salida Esperada

```
🔍 Validando configuración de Docker Compose...

✅ docker-compose instalado
✅ docker-compose.yml encontrado

📝 Validando sintaxis YAML...
✅ Sintaxis YAML válida

🔍 Servicios encontrados:
  ✓ postgres
  ✓ mongodb
  ✓ rabbitmq

💾 Volúmenes encontrados:
  ✓ edugo_postgres_data
  ✓ edugo_mongo_data
  ✓ edugo_rabbitmq_data

🌐 Puertos expuestos:
  ✓ 5432:5432
  ✓ 15672:15672
  ✓ 27017:27017
  ✓ 5672:5672

✅ Validación completada exitosamente
```

### Errores Comunes

**Error: docker-compose no instalado**
```bash
❌ ERROR: docker-compose no está instalado
```
Solución: Instalar Docker Compose

**Error: sintaxis YAML inválida**
```bash
❌ ERROR: Sintaxis YAML inválida
```
Solución: Revisar indentación y sintaxis en docker-compose.yml
```

**Checklist:**
- [ ] README.md de scripts actualizado
- [ ] Documentación del script completa
- [ ] Ejemplos de uso incluidos

---

### Validación de Tarea 3.2

**Checklist final:**
- [ ] `scripts/validate.sh` creado y ejecutable
- [ ] Script valida sintaxis YAML
- [ ] Script imprime servicios/volúmenes/puertos
- [ ] Documentación en `scripts/README.md`
- [ ] Probado localmente con éxito

**Comando de validación:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment
./scripts/validate.sh
```

---

## Tarea 3.3: Pre-commit Hook Opcional

**⏱️ Tiempo estimado:** 30 minutos  
**🎯 Objetivo:** Hook opcional para validar antes de commit  
**📍 Ubicación:** `/repos-separados/edugo-dev-environment/.githooks/`

### Subtareas

#### 3.1 Crear Pre-commit Hook

**Archivo:** `.githooks/pre-commit`

**Contenido:**

```bash
#!/bin/bash

# edugo-dev-environment - Pre-commit Hook
# Valida docker-compose.yml antes de permitir commit

echo "🔍 Ejecutando validaciones pre-commit..."
echo ""

# Validar docker-compose.yml si fue modificado
if git diff --cached --name-only | grep -q "docker-compose.yml"; then
    echo "📝 docker-compose.yml modificado, validando..."
    
    if ! ./scripts/validate.sh; then
        echo ""
        echo "❌ Validación falló. Commit bloqueado."
        echo ""
        echo "Soluciones:"
        echo "  1. Corregir errores en docker-compose.yml"
        echo "  2. Omitir validación: git commit --no-verify"
        exit 1
    fi
    
    echo ""
    echo "✅ docker-compose.yml válido"
fi

# Validar que .env no se commitee
if git diff --cached --name-only | grep -q "^\.env$"; then
    echo ""
    echo "⚠️  ADVERTENCIA: Intentando commitear .env"
    echo ""
    echo "El archivo .env contiene secrets y NO debe commitearse."
    echo ""
    echo "Soluciones:"
    echo "  1. git reset HEAD .env"
    echo "  2. Agregar .env a .gitignore"
    echo "  3. Omitir validación: git commit --no-verify (NO recomendado)"
    exit 1
fi

echo ""
echo "✅ Todas las validaciones pasaron"
exit 0
```

**Checklist:**
- [ ] Hook creado en `.githooks/pre-commit`
- [ ] Permisos de ejecución agregados
- [ ] Valida docker-compose.yml
- [ ] Valida que .env no se commitee

**Comandos:**
```bash
mkdir -p .githooks
touch .githooks/pre-commit
chmod +x .githooks/pre-commit
```

---

#### 3.2 Documentar Cómo Activar el Hook

**Archivo:** `.githooks/README.md`

**Crear con:**

```markdown
# Git Hooks

Este directorio contiene hooks opcionales de Git.

## Pre-commit Hook

Valida configuración antes de permitir commits.

### ¿Qué Valida?

- ✅ Sintaxis de docker-compose.yml (si fue modificado)
- ✅ Que .env NO se commitee accidentalmente

### Activar el Hook

**Opción A: Manual (recomendado)**
```bash
git config core.hooksPath .githooks
```

**Opción B: Symlink**
```bash
ln -s ../../.githooks/pre-commit .git/hooks/pre-commit
```

### Desactivar el Hook

```bash
git config --unset core.hooksPath
```

### Omitir el Hook (Una Vez)

Si necesitas hacer un commit urgente sin validación:

```bash
git commit --no-verify
```

**⚠️ ADVERTENCIA:** Solo usa `--no-verify` si estás seguro.

### Probar el Hook

```bash
# Modificar docker-compose.yml con error
echo "invalid yaml:" >> docker-compose.yml

# Intentar commit
git add docker-compose.yml
git commit -m "test"

# Esperado: Commit bloqueado con mensaje de error

# Revertir cambio
git checkout docker-compose.yml
```

## Otros Hooks

Actualmente solo hay pre-commit. Otros hooks pueden agregarse según necesidad.
```

**Checklist:**
- [ ] README.md de hooks creado
- [ ] Instrucciones de activación incluidas
- [ ] Instrucciones de desactivación incluidas
- [ ] Ejemplo de uso incluido

---

#### 3.3 Actualizar README.md Principal

**Archivo:** `README.md`

**Agregar sección:**

```markdown
## 🪝 Pre-commit Hooks (Opcional)

Este proyecto incluye hooks opcionales que validan configuración antes de commits.

### Activar Hooks

```bash
git config core.hooksPath .githooks
```

### ¿Qué Validan?

- ✅ Sintaxis de `docker-compose.yml`
- ✅ Que `.env` no se commitee accidentalmente

### Desactivar Hooks

```bash
git config --unset core.hooksPath
```

### Más Información

Ver [.githooks/README.md](.githooks/README.md)
```

**Checklist:**
- [ ] Sección de hooks agregada al README principal
- [ ] Links a documentación incluidos

---

### Validación de Tarea 3.3

**Checklist final:**
- [ ] `.githooks/pre-commit` creado y ejecutable
- [ ] `.githooks/README.md` con documentación completa
- [ ] README.md principal actualizado
- [ ] Hook probado (activar, intentar commit inválido, desactivar)

**Comando de validación:**
```bash
# Activar hook
git config core.hooksPath .githooks

# Probar con cambio inválido
echo "test" >> docker-compose.yml
git add docker-compose.yml
git commit -m "test"  # Debería fallar

# Limpiar
git checkout docker-compose.yml
git config --unset core.hooksPath
```

---

## Tarea 3.4: Documentar Decisión de NO CI/CD

**⏱️ Tiempo estimado:** 15 minutos  
**🎯 Objetivo:** Explicar por qué este proyecto NO tiene workflows  
**📍 Ubicación:** `/repos-separados/edugo-dev-environment/README.md`

### Subtareas

#### 4.1 Agregar Sección al README

**Archivo:** `README.md`

**Agregar sección (después de Arquitectura):**

```markdown
## 🤔 ¿Por Qué NO Hay CI/CD?

Este proyecto **intencionalmente NO tiene workflows de CI/CD**.

### Razón

Este es un repositorio de **configuración**, no de **código**:

- ✅ Contiene `docker-compose.yml` (configuración)
- ✅ Contiene scripts de setup (utilidades)
- ❌ NO contiene código de aplicación
- ❌ NO tiene tests unitarios
- ❌ NO genera builds
- ❌ NO se despliega a producción

### Filosofía

> "No uses CI/CD para todo. Úsalo solo donde agregue valor."

**CI/CD es útil para:**
- ✅ Proyectos con tests (api-mobile, api-administracion, worker)
- ✅ Proyectos con builds (imágenes Docker, binarios)
- ✅ Proyectos con releases (shared, infrastructure)

**CI/CD NO es útil para:**
- ❌ Repos de configuración (este proyecto)
- ❌ Repos de documentación pura
- ❌ Repos de scripts de utilidad

### Validación Local en Lugar de CI/CD

En lugar de workflows, usamos **validación local**:

```bash
# Validar sintaxis (instantáneo)
./scripts/validate.sh

# Validar que funciona (levanta servicios)
docker-compose up -d
docker-compose ps
```

**Ventajas:**
- ⚡ Feedback instantáneo (sin esperar CI)
- 💰 Cero consumo de minutos de GitHub Actions
- 🎯 Más simple y directo

### Referencias

Para más detalles sobre esta decisión:
- [Análisis de Estado Actual CI/CD](../../00-Projects-Isolated/cicd-analysis/01-ANALISIS-ESTADO-ACTUAL.md#edugo-dev-environment)
- [Plan de Implementación](../../00-Projects-Isolated/cicd-analysis/implementation-plans/06-dev-environment/README.md)
```

**Checklist:**
- [ ] Sección agregada al README
- [ ] Razones explicadas claramente
- [ ] Alternativas documentadas
- [ ] Referencias incluidas

---

### Validación de Tarea 3.4

**Checklist final:**
- [ ] README.md tiene sección "¿Por Qué NO Hay CI/CD?"
- [ ] Razones técnicas explicadas
- [ ] Filosofía documentada
- [ ] Alternativa de validación local explicada
- [ ] Referencias a documentos de análisis incluidas

---

## Tarea 3.5: Crear Ejemplo End-to-End

**⏱️ Tiempo estimado:** 30-45 minutos  
**🎯 Objetivo:** Guía completa de uso del entorno  
**📍 Ubicación:** `/repos-separados/edugo-dev-environment/docs/`

### Subtareas

#### 5.1 Crear EXAMPLE.md

**Archivo:** `docs/EXAMPLE.md`

**Contenido:**

```markdown
# Ejemplo End-to-End - edugo-dev-environment

Esta guía te lleva paso a paso para levantar el entorno completo de desarrollo.

---

## 🎯 Objetivo

Al final de esta guía tendrás:
- ✅ PostgreSQL corriendo en puerto 5432
- ✅ MongoDB corriendo en puerto 27017
- ✅ RabbitMQ corriendo en puertos 5672 y 15672

---

## 📋 Paso 1: Clonar el Repositorio

```bash
# Clonar
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment

# Verificar contenido
ls -la
```

**Esperado:**
```
docker-compose.yml
.env.example
scripts/
README.md
```

---

## 📋 Paso 2: Configurar Variables de Entorno

```bash
# Copiar ejemplo
cp .env.example .env

# Editar (opcional)
nano .env
```

**Contenido de .env:**
```env
# PostgreSQL
POSTGRES_USER=edugo
POSTGRES_PASSWORD=edugo_dev_2024
POSTGRES_DB=edugo

# MongoDB
MONGO_INITDB_ROOT_USERNAME=edugo
MONGO_INITDB_ROOT_PASSWORD=edugo_dev_2024
MONGO_INITDB_DATABASE=edugo

# RabbitMQ
RABBITMQ_DEFAULT_USER=edugo
RABBITMQ_DEFAULT_PASS=edugo_dev_2024
```

**⚠️ NOTA:** Estas son credenciales de desarrollo. NO usar en producción.

---

## 📋 Paso 3: Validar Configuración (Opcional)

```bash
# Ejecutar validación
./scripts/validate.sh
```

**Esperado:**
```
🔍 Validando configuración de Docker Compose...

✅ docker-compose instalado
✅ docker-compose.yml encontrado
✅ Sintaxis YAML válida

🔍 Servicios encontrados:
  ✓ postgres
  ✓ mongodb
  ✓ rabbitmq

✅ Validación completada exitosamente
```

---

## 📋 Paso 4: Levantar Servicios

```bash
# Levantar en background
docker-compose up -d

# Ver logs
docker-compose logs -f
```

**Esperado:**
```
[+] Running 3/3
 ✔ Container edugo-postgres   Started
 ✔ Container edugo-mongodb    Started
 ✔ Container edugo-rabbitmq   Started
```

---

## 📋 Paso 5: Verificar Servicios

```bash
# Ver estado
docker-compose ps
```

**Esperado:**
```
NAME                STATUS        PORTS
edugo-postgres      Up 10s        0.0.0.0:5432->5432/tcp
edugo-mongodb       Up 10s        0.0.0.0:27017->27017/tcp
edugo-rabbitmq      Up 10s        0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
```

---

## 📋 Paso 6: Probar Conexiones

### PostgreSQL

```bash
# Usando psql (si está instalado)
psql -h localhost -U edugo -d edugo

# O usando Docker
docker exec -it edugo-postgres psql -U edugo -d edugo
```

**Esperado:**
```sql
edugo=# SELECT version();
-- PostgreSQL 15.x ...

edugo=# \l
-- Lista de bases de datos

edugo=# \q
-- Salir
```

### MongoDB

```bash
# Usando mongosh (si está instalado)
mongosh mongodb://edugo:edugo_dev_2024@localhost:27017/edugo

# O usando Docker
docker exec -it edugo-mongodb mongosh -u edugo -p edugo_dev_2024
```

**Esperado:**
```javascript
edugo> show dbs
admin   0.000GB
config  0.000GB
edugo   0.000GB

edugo> db.test.insertOne({message: "Hello EduGo"})
{ acknowledged: true, insertedId: ... }

edugo> exit
```

### RabbitMQ

**Interfaz Web:**
1. Abrir navegador: http://localhost:15672
2. Login:
   - Usuario: `edugo`
   - Password: `edugo_dev_2024`
3. Explorar dashboard

**Esperado:** Dashboard de RabbitMQ con pestañas Overview, Connections, Channels, Exchanges, Queues.

---

## 📋 Paso 7: Usar con las APIs de EduGo

### 7.1 Clonar Repos de APIs

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/

# Si no están clonados
git clone https://github.com/EduGoGroup/edugo-api-mobile.git
git clone https://github.com/EduGoGroup/edugo-api-administracion.git
git clone https://github.com/EduGoGroup/edugo-worker.git
```

### 7.2 Configurar APIs

Cada API usa las mismas variables de entorno:

**Archivo:** `.env` (en cada repo)
```env
# Database
DATABASE_URL=postgresql://edugo:edugo_dev_2024@localhost:5432/edugo

# MongoDB
MONGO_URI=mongodb://edugo:edugo_dev_2024@localhost:27017/edugo

# RabbitMQ
RABBITMQ_URL=amqp://edugo:edugo_dev_2024@localhost:5672/
```

### 7.3 Ejecutar Migraciones

```bash
cd edugo-api-mobile/

# Ejecutar migraciones
make migrate-up
```

**Esperado:**
```
Running migrations...
✅ Migrations completed successfully
```

### 7.4 Levantar API

```bash
# API Mobile
cd edugo-api-mobile/
go run main.go

# En otra terminal: API Administración
cd edugo-api-administracion/
go run main.go

# En otra terminal: Worker
cd edugo-worker/
go run main.go
```

**Esperado:**
```
Starting edugo-api-mobile...
Listening on :8080
```

---

## 📋 Paso 8: Verificar Todo Funciona

### Verificar API Mobile

```bash
curl http://localhost:8080/health
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

### Verificar API Administración

```bash
curl http://localhost:8081/health
```

### Verificar Worker

```bash
# Ver logs
docker-compose logs edugo-worker
```

---

## 📋 Paso 9: Detener Servicios

```bash
# Detener y remover contenedores
docker-compose down

# Detener, remover contenedores Y volúmenes (⚠️ borra datos)
docker-compose down -v
```

---

## 🎉 ¡Listo!

Has completado el setup completo del entorno de desarrollo EduGo.

### Próximos Pasos

1. Explorar las APIs: http://localhost:8080/swagger
2. Crear datos de prueba
3. Desarrollar nuevas funcionalidades

### Comandos Útiles

```bash
# Ver logs de todos los servicios
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f postgres

# Reiniciar un servicio
docker-compose restart mongodb

# Ver uso de recursos
docker stats
```

---

## 🔧 Troubleshooting

Si algo falla, ver [README.md - Solución de Problemas](../README.md#-solución-de-problemas-comunes).
```

**Checklist:**
- [ ] EXAMPLE.md creado en `docs/`
- [ ] Ejemplo completo paso a paso
- [ ] Comandos copy-paste listos
- [ ] Output esperado documentado
- [ ] Troubleshooting referenciado

---

#### 5.2 Agregar Link en README Principal

**Archivo:** `README.md`

**Agregar en sección de "Uso":**

```markdown
## 🚀 Uso

### Quick Start

```bash
# 1. Configurar
cp .env.example .env

# 2. Validar (opcional)
./scripts/validate.sh

# 3. Levantar
docker-compose up -d

# 4. Verificar
docker-compose ps
```

### Guía Completa

Para un tutorial paso a paso detallado, ver:
👉 [docs/EXAMPLE.md](docs/EXAMPLE.md)
```

**Checklist:**
- [ ] Link a EXAMPLE.md agregado
- [ ] Quick start incluido

---

### Validación de Tarea 3.5

**Checklist final:**
- [ ] `docs/EXAMPLE.md` creado con guía completa
- [ ] Todos los pasos incluyen comandos copy-paste
- [ ] Output esperado documentado
- [ ] README.md principal actualizado con link
- [ ] Probado siguiendo la guía

---

## ✅ Validación del Sprint Completo

### Checklist General

- [ ] Tarea 3.1: README.md mejorado con requisitos, troubleshooting, arquitectura
- [ ] Tarea 3.2: Script `validate.sh` funcional
- [ ] Tarea 3.3: Pre-commit hook opcional creado
- [ ] Tarea 3.4: Decisión de NO CI/CD documentada
- [ ] Tarea 3.5: Ejemplo end-to-end creado

### Archivos Creados/Modificados

**Nuevos:**
- [ ] `scripts/validate.sh`
- [ ] `scripts/README.md`
- [ ] `.githooks/pre-commit`
- [ ] `.githooks/README.md`
- [ ] `docs/EXAMPLE.md`

**Modificados:**
- [ ] `README.md` (varias secciones agregadas)

### Comandos de Validación Final

```bash
# Ir al repo
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment

# Verificar archivos
ls scripts/validate.sh
ls .githooks/pre-commit
ls docs/EXAMPLE.md

# Probar validación
./scripts/validate.sh

# Probar pre-commit (sin activarlo)
./.githooks/pre-commit

# Levantar entorno
docker-compose up -d
docker-compose ps
docker-compose down
```

---

## 📊 Resumen del Sprint

### Métricas

| Métrica | Valor |
|---------|-------|
| Tareas completadas | 5 |
| Archivos creados | 5 |
| Archivos modificados | 1 |
| Scripts bash | 2 |
| Tiempo total | 2-3 horas |
| Líneas de código | ~500 |
| Líneas de documentación | ~800 |

### Antes y Después

**Antes del Sprint:**
- ✅ docker-compose.yml funcional
- ❌ Poca documentación
- ❌ Sin validación local
- ❌ Sin hooks

**Después del Sprint:**
- ✅ docker-compose.yml funcional
- ✅ Documentación completa
- ✅ Validación local con script
- ✅ Pre-commit hooks opcionales
- ✅ Ejemplo end-to-end
- ✅ Decisión de NO CI/CD documentada

---

## 🎯 Próximos Pasos (Post-Sprint)

### Opcional: Si Quieres Mejorar Más

1. **Agregar más scripts de utilidad**
   - `scripts/backup.sh` - Backup de volúmenes
   - `scripts/restore.sh` - Restore de backups
   - `scripts/clean.sh` - Limpiar volúmenes

2. **Crear docker-compose.prod.yml**
   - Configuración para producción
   - Sin puertos expuestos
   - Con secrets reales

3. **Agregar monitoreo**
   - Prometheus para métricas
   - Grafana para dashboards

**⚠️ NOTA:** Estas mejoras son opcionales y NO críticas.

---

## 📝 Notas Finales

### Este Sprint es OPCIONAL

Si el proyecto ya funciona bien y está documentado, **NO es necesario ejecutar este sprint**.

### Cuándo SÍ Ejecutarlo

- ❓ README.md es confuso o incompleto
- ❓ Nuevos devs tienen problemas al setup
- ❓ No hay validación de docker-compose.yml
- ❓ docker-compose.yml tiene errores frecuentes

### Cuándo NO Ejecutarlo

- ✅ README.md es claro y completo
- ✅ docker-compose.yml funciona bien
- ✅ Devs no tienen problemas al setup
- ✅ No hay tiempo para mejoras opcionales

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Prioridad:** BAJA (opcional)
