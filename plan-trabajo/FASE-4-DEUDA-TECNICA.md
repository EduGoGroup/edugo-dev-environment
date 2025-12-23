# FASE 4: Deuda Técnica y Mejoras de Baja Prioridad

**Prioridad:** Baja  
**Estimación:** 2-3 días  
**Riesgo:** Bajo  

---

## Objetivo

Abordar la deuda técnica acumulada y mejoras que, aunque útiles, no son críticas.

---

## 4.1 Tests de Integración para Docker Compose (CI/CD)

### Problema Actual

No hay CI/CD que valide que los servicios levantan correctamente.

### Solución Propuesta

Crear GitHub Actions que valide el ambiente Docker.

### Pasos de Implementación

#### Paso 4.1.1: Crear workflow de GitHub Actions

**Archivo:** `.github/workflows/docker-compose-test.yml`

```yaml
name: Docker Compose Test

on:
  push:
    branches: [main, dev]
    paths:
      - 'docker/**'
      - 'migrator/**'
      - 'scripts/**'
  pull_request:
    branches: [main]
    paths:
      - 'docker/**'
      - 'migrator/**'
      - 'scripts/**'

jobs:
  test-docker-compose:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
          
      - name: Create .env file
        run: |
          cd docker
          cp .env.example .env
          echo "OPENAI_API_KEY=sk-test-not-real" >> .env
          
      - name: Start infrastructure only
        run: |
          cd docker
          docker-compose -f docker-compose-infrastructure.yml up -d
          
      - name: Wait for health checks
        run: |
          echo "Esperando PostgreSQL..."
          timeout 60 bash -c 'until docker exec edugo-postgres pg_isready; do sleep 2; done'
          echo "Esperando MongoDB..."
          timeout 60 bash -c 'until docker exec edugo-mongodb mongosh --eval "db.adminCommand(\"ping\")" -u edugo -p edugo123 --authSource admin; do sleep 2; done'
          echo "Esperando RabbitMQ..."
          timeout 60 bash -c 'until docker exec edugo-rabbitmq rabbitmq-diagnostics ping; do sleep 2; done'
          
      - name: Run migrator
        run: |
          cd docker
          docker-compose -f docker-compose.migrate.yml up --build
          
      - name: Verify migrations
        run: |
          docker exec edugo-postgres psql -U edugo -d edugo -c "\dt" | grep -q "users"
          echo "Tabla users existe ✓"
          
      - name: Cleanup
        if: always()
        run: |
          cd docker
          docker-compose down -v
```

#### Paso 4.1.2: Agregar badge al README

**Archivo:** `README.md`

```markdown
![Docker Compose Test](https://github.com/EduGoGroup/edugo-dev-environment/actions/workflows/docker-compose-test.yml/badge.svg)
```

### Validación

- [ ] Workflow ejecuta en PR a main
- [ ] Infrastructure levanta correctamente
- [ ] Migrator ejecuta sin errores
- [ ] Badge muestra estado

### Commit Sugerido

```
ci: agregar workflow de pruebas de Docker Compose

- Crear docker-compose-test.yml para GitHub Actions
- Verificar health checks de infraestructura
- Ejecutar migraciones
- Agregar badge de estado al README
```

---

## 4.2 Soporte para Windows/WSL

### Problema Actual

Documentación solo menciona macOS.

### Solución Propuesta

Agregar sección en guía rápida para WSL2.

### Pasos de Implementación

#### Paso 4.2.1: Agregar sección en GUIA-RAPIDA.md

**Archivo:** `documentos/GUIA-RAPIDA.md`

```markdown
## Notas para Windows (WSL2)

### Prerequisitos

1. Windows 10/11 con WSL2 habilitado
2. Distribución Linux instalada (Ubuntu recomendado)
3. Docker Desktop for Windows con integración WSL2

### Configuración

1. Habilitar WSL2:
   ```powershell
   wsl --install
   ```

2. Instalar Docker Desktop con WSL2 backend

3. En Docker Desktop:
   - Settings → Resources → WSL Integration
   - Habilitar para tu distribución

4. Clonar y ejecutar desde WSL:
   ```bash
   # En terminal WSL (Ubuntu)
   cd /mnt/c/proyectos  # o tu directorio preferido
   git clone https://github.com/EduGoGroup/edugo-dev-environment.git
   cd edugo-dev-environment
   ./scripts/setup.sh
   ```

### Notas Importantes

- Ejecutar siempre desde terminal WSL, no PowerShell
- Los archivos deben estar en el filesystem de WSL para mejor rendimiento
- Acceder a servicios desde Windows: `localhost:8081`
```

### Validación

- [ ] Documentación WSL agregada
- [ ] Instrucciones probadas en WSL2

### Commit Sugerido

```
docs: agregar soporte y documentación para Windows/WSL2
```

---

## 4.3 Script de Backup/Restore

### Problema Actual

No hay forma fácil de respaldar datos de desarrollo.

### Solución Propuesta

Crear scripts para backup y restore de datos.

### Pasos de Implementación

#### Paso 4.3.1: Crear script de backup

**Archivo:** `scripts/backup.sh`

```bash
#!/bin/bash
set -e

BACKUP_DIR="backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="edugo_backup_$TIMESTAMP"

mkdir -p "$BACKUP_DIR"

echo "=== Respaldando datos de EduGo ==="

# Backup PostgreSQL
echo "Respaldando PostgreSQL..."
docker exec edugo-postgres pg_dump -U edugo edugo > "$BACKUP_DIR/${BACKUP_NAME}_postgres.sql"

# Backup MongoDB
echo "Respaldando MongoDB..."
docker exec edugo-mongodb mongodump --uri="mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin" --archive > "$BACKUP_DIR/${BACKUP_NAME}_mongodb.archive"

echo ""
echo "Backup completado en: $BACKUP_DIR/"
echo "  - ${BACKUP_NAME}_postgres.sql"
echo "  - ${BACKUP_NAME}_mongodb.archive"
```

#### Paso 4.3.2: Crear script de restore

**Archivo:** `scripts/restore.sh`

```bash
#!/bin/bash
set -e

if [ -z "$1" ]; then
  echo "Uso: ./scripts/restore.sh <nombre_backup>"
  echo ""
  echo "Backups disponibles:"
  ls -1 backups/*.sql 2>/dev/null | sed 's/_postgres.sql//' | sed 's/backups\///' | sort -u
  exit 1
fi

BACKUP_NAME=$1

echo "=== Restaurando datos de EduGo ==="
echo "ADVERTENCIA: Esto sobrescribirá los datos actuales."
read -p "¿Continuar? (y/N) " confirm

if [ "$confirm" != "y" ]; then
  echo "Cancelado."
  exit 0
fi

# Restore PostgreSQL
echo "Restaurando PostgreSQL..."
docker exec -i edugo-postgres psql -U edugo -d edugo < "backups/${BACKUP_NAME}_postgres.sql"

# Restore MongoDB
echo "Restaurando MongoDB..."
docker exec -i edugo-mongodb mongorestore --uri="mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin" --archive < "backups/${BACKUP_NAME}_mongodb.archive"

echo ""
echo "Restore completado."
```

#### Paso 4.3.3: Agregar al Makefile

**Archivo:** `Makefile`

```makefile
backup: ## Respaldar bases de datos
	./scripts/backup.sh

restore: ## Restaurar bases de datos (usar: make restore BACKUP=nombre)
	./scripts/restore.sh $(BACKUP)
```

#### Paso 4.3.4: Agregar .gitignore para backups

**Archivo:** `.gitignore`

```
backups/
```

### Validación

- [ ] `./scripts/backup.sh` crea archivos de respaldo
- [ ] `./scripts/restore.sh` restaura datos
- [ ] `make backup` funciona
- [ ] Backups no se suben a git

### Commit Sugerido

```
feat(scripts): agregar scripts de backup y restore

- Crear backup.sh para PostgreSQL y MongoDB
- Crear restore.sh con confirmación
- Agregar comandos make backup/restore
- Ignorar carpeta backups en git
```

---

## 4.4 Consolidar Docker Compose Files

### Problema Actual

Existen múltiples archivos docker-compose con propósitos solapados.

### Solución Propuesta

Consolidar usando perfiles de Docker Compose v2.

### Pasos de Implementación

#### Paso 4.4.1: Evaluar archivos actuales

| Archivo | Propósito | Acción |
|---------|-----------|--------|
| docker-compose.yml | Principal | Mantener como base |
| docker-compose-apps.yml | Solo apps | Integrar con profiles |
| docker-compose-infrastructure.yml | Solo infra | Integrar con profiles |
| docker-compose-mock.yml | Con mocks | Evaluar si se necesita |
| docker-compose.migrate.yml | Solo migrator | Mantener separado |

#### Paso 4.4.2: Unificar con profiles

**Archivo:** `docker/docker-compose.yml`

```yaml
version: '3.8'

services:
  postgres:
    # Sin profile - infraestructura siempre disponible
    image: postgres:16-alpine
    ...
    
  mongodb:
    # Sin profile - infraestructura siempre disponible
    image: mongo:7.0
    ...
    
  rabbitmq:
    # Sin profile - infraestructura siempre disponible
    image: rabbitmq:3.12-management-alpine
    ...
    
  api-mobile:
    profiles: ["apps", "full"]
    image: ghcr.io/edugogroup/edugo-api-mobile:latest
    ...
    
  api-administracion:
    profiles: ["apps", "admin", "full"]
    image: ghcr.io/edugogroup/edugo-api-administracion:latest
    ...
    
  worker:
    profiles: ["apps", "worker", "full"]
    image: ghcr.io/edugogroup/edugo-worker:latest
    ...
```

#### Paso 4.4.3: Actualizar scripts

Actualizar `setup.sh` para usar profiles:

```bash
# Infraestructura solamente
docker-compose up -d

# Con todas las apps
docker-compose --profile full up -d

# Solo con API Mobile
docker-compose --profile apps up -d api-mobile
```

#### Paso 4.4.4: Archivar archivos redundantes

Mover a `archivado-documentos/docker-legacy/`:
- `docker-compose-apps.yml`
- `docker-compose-infrastructure.yml`
- `docker-compose-mock.yml` (si no se usa)

### Validación

- [ ] docker-compose.yml unificado funciona
- [ ] Profiles se activan correctamente
- [ ] Documentación actualizada
- [ ] Archivos legacy archivados

### Commit Sugerido

```
refactor(docker): consolidar docker-compose files usando profiles

- Unificar servicios en docker-compose.yml principal
- Implementar profiles: apps, admin, worker, full
- Archivar archivos redundantes
- Actualizar documentación
```

---

## 4.5 Mover Documentación de docker/ a documentos/

### Problema Actual

Hay documentación dentro de `docker/` que debería estar centralizada.

### Solución Propuesta

Mover documentación a `documentos/` o archivar si está desactualizada.

### Pasos de Implementación

#### Paso 4.5.1: Evaluar documentos en docker/

| Archivo | Estado | Acción |
|---------|--------|--------|
| ACTUALIZAR_BASE_DATOS.md | Revisar | Integrar en docs o archivar |
| PLAN_PRUEBAS_DOCKER_COMPOSE.md | Temporal | Archivar |
| QUICK_START.md | Redundante | Archivar (existe GUIA-RAPIDA) |
| README.md | Útil | Mantener (específico de docker) |
| RESULTADO_VALIDACION.md | Temporal | Archivar |

#### Paso 4.5.2: Ejecutar archivado

```bash
mkdir -p archivado-documentos/docker-docs
mv docker/PLAN_PRUEBAS_DOCKER_COMPOSE.md archivado-documentos/docker-docs/
mv docker/QUICK_START.md archivado-documentos/docker-docs/
mv docker/RESULTADO_VALIDACION.md archivado-documentos/docker-docs/
```

### Validación

- [ ] Documentación reorganizada
- [ ] docker/ contiene solo archivos necesarios
- [ ] Archivos archivados mantienen historial

### Commit Sugerido

```
chore(docs): reorganizar documentación de docker/

- Archivar documentos temporales/redundantes
- Mantener README.md específico de docker
- Centralizar documentación en documentos/
```

---

## 4.6 Publicar Imagen del Migrator en ghcr.io

### Problema Actual

El migrator se construye localmente aumentando tiempo de setup.

### Solución Propuesta

Publicar imagen pre-construida en GitHub Container Registry.

### Pasos de Implementación

Esta tarea requiere:
1. Crear workflow de CI para construir y publicar imagen
2. Modificar docker-compose para usar imagen publicada
3. Mantener opción de build local para desarrollo

**Nota:** Esta es una tarea más compleja que puede requerir coordinación con el equipo de DevOps.

### Commit Sugerido

```
feat(migrator): publicar imagen en ghcr.io

- Crear workflow de CI para build y push
- Actualizar docker-compose para usar imagen
- Documentar proceso de actualización
```

---

## Resumen de Commits de Fase 4

1. `ci: agregar workflow de pruebas de Docker Compose`
2. `docs: agregar soporte y documentación para Windows/WSL2`
3. `feat(scripts): agregar scripts de backup y restore`
4. `refactor(docker): consolidar docker-compose files usando profiles`
5. `chore(docs): reorganizar documentación de docker/`
6. `feat(migrator): publicar imagen en ghcr.io` (opcional)

---

## Dependencias

- Fases 1-3 completadas (recomendado)
- Acceso a GitHub Actions (para CI/CD)
- Permisos en ghcr.io (para publicar migrator)

---

## Flujo de Trabajo Git

### 1. Crear rama desde dev

```bash
git checkout dev
git pull origin dev
git checkout -b fase-4-deuda-tecnica
```

### 2. Realizar los cambios

Ejecutar los pasos de implementación descritos arriba. Dado que esta fase tiene múltiples tareas independientes, se recomienda hacer commits atómicos por cada mejora (4.1, 4.2, 4.3, etc.).

### 3. Crear PR hacia dev

```bash
git push origin fase-4-deuda-tecnica
# Crear PR en GitHub hacia dev
```

---

## Documentación a Actualizar

Al completar esta fase, actualizar los siguientes documentos:

| Documento | Cambio Requerido |
|-----------|------------------|
| `documentos/GUIA-RAPIDA.md` | Agregar sección Windows/WSL2, comandos backup |
| `documentos/FAQ.md` | Agregar preguntas sobre backup/restore |
| `documentos/SERVICIOS.md` | Actualizar referencias a docker-compose consolidado |
| `documentos/DEPRECADO-MEJORAS.md` | Marcar todas las mejoras de deuda técnica como completadas |
| `README.md` | Agregar badge de CI/CD |

### Checklist de Cierre

- [ ] Workflow de CI/CD creado y funcionando
- [ ] Documentación WSL2 agregada
- [ ] Scripts backup/restore creados
- [ ] Docker Compose files consolidados
- [ ] Documentación de docker/ reorganizada
- [ ] Imagen de migrator publicada (opcional)
- [ ] `documentos/GUIA-RAPIDA.md` actualizado
- [ ] `documentos/FAQ.md` actualizado
- [ ] `documentos/SERVICIOS.md` actualizado
- [ ] `documentos/DEPRECADO-MEJORAS.md` actualizado
- [ ] PR creado hacia `dev`
- [ ] PR revisado y aprobado
- [ ] PR mergeado a `dev`
