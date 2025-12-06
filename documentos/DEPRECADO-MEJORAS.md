# Deprecado y Mejoras - EduGo Dev Environment

Este documento registra código deprecado, mejoras pendientes y deuda técnica del proyecto.

---

## Código Deprecado

### Archivos Movidos a `archivado-documentos/`

Los siguientes archivos fueron archivados por estar desactualizados o ser redundantes:

| Archivo | Razón | Acción |
|---------|-------|--------|
| `CLAUDE.md` | Reglas de sprint específicas de Claude, no generales | Archivado |
| `SETUP-PARA-FLOJOS.md` | Redundante con nueva guía rápida | Archivado |
| `PLAN-POST-MERGE.md` | Documento de planificación temporal | Archivado |
| `CHANGELOG-PASSWORDS.md` | Historial específico, no documentación general | Archivado |
| `RESUMEN-IMPLEMENTACION-PASSWORDS.md` | Documento de implementación completada | Archivado |
| `VALIDACION-FINAL-COMPLETA.md` | Reporte de validación puntual | Archivado |
| `docs/` (carpeta completa) | Documentación fragmentada y desactualizada | Archivado |

### Configuraciones Deprecadas

| Configuración | Estado | Alternativa |
|---------------|--------|-------------|
| Perfiles de docker-compose | Parcialmente implementado | Usar `--profile` en setup.sh |
| Mock services | En `docker-compose-mock.yml` | Evaluar si aún se necesitan |

---

## Mejoras Implementadas

### Diciembre 2025

#### ✅ Makefile para comandos simplificados

**Implementado:** Se creó `Makefile` en la raíz con comandos:
- `make up` / `make down` / `make restart`
- `make logs` / `make logs-api` / `make logs-admin`
- `make status` / `make health`
- `make psql` / `make mongo`
- `make setup` / `make diagnose` / `make clean`

#### ✅ Script de Diagnóstico

**Implementado:** Se creó `scripts/diagnose.sh` que verifica:
- Estado de Docker
- Contenedores corriendo y su health
- Puertos en uso
- Conectividad de APIs
- Configuración (.env, OPENAI_API_KEY)
- Autenticación en ghcr.io
- Errores recientes en logs

---

## Mejoras Pendientes

### Alta Prioridad

#### 1. Validación de Health Checks en Setup

**Problema:** `setup.sh` espera 10 segundos fijos en lugar de verificar health checks.

**Mejora propuesta:**
```bash
# En lugar de:
sleep 10

# Usar:
until docker-compose exec postgres pg_isready; do
  echo "Esperando PostgreSQL..."
  sleep 2
done
```

**Complejidad:** Baja (2-3 horas)

#### 2. Seed Data Automático

**Problema:** `seed-data.sh` existe pero no se ejecuta por defecto en setup.

**Mejora propuesta:** Integrar seed en el migrator o como paso opcional automático.

**Complejidad:** Media (1 día)

### Media Prioridad

#### 3. Centralizar Variables de Entorno

**Problema:** Algunas variables se repiten entre servicios con nombres diferentes:
- API Mobile usa `DATABASE_POSTGRES_*`
- API Admin usa `EDUGO_ADMIN_DATABASE_POSTGRES_*`
- Worker usa `EDUGO_WORKER_DATABASE_POSTGRES_*`

**Mejora propuesta:** Documentar claramente o unificar en futuras versiones de las APIs.

**Complejidad:** Alta (requiere cambios en repos de APIs)

#### 4. Soporte para Apple Silicon (M1/M2/M3)

**Problema:** No hay documentación específica para arquitectura ARM.

**Mejora propuesta:** Verificar y documentar compatibilidad.

**Complejidad:** Baja (investigación y documentación)

#### 5. Seeds más Completos

**Problema:** Los seeds actuales son mínimos (solo 2 archivos SQL y 1 JS).

**Mejora propuesta:** Agregar más datos de prueba:
- Más usuarios con diferentes roles
- Cursos con materiales
- Instituciones completas
- Documentos de ejemplo

**Complejidad:** Media (2-3 horas)

### Baja Prioridad

#### 6. Tests de Integración para Docker Compose

**Problema:** No hay CI/CD que valide que los servicios levantan correctamente.

**Mejora propuesta:** GitHub Actions que:
1. Ejecute `docker-compose up -d`
2. Espere health checks
3. Ejecute curl a endpoints de health
4. Reporte resultado

**Complejidad:** Media (1 día)

#### 7. Soporte para Windows/WSL

**Problema:** Documentación solo menciona macOS.

**Mejora propuesta:** Agregar sección en guía rápida para WSL2.

**Complejidad:** Baja (documentación)

#### 8. docker-compose-apps.yml tiene Profiles Rotos

**Problema:** Usa `profiles: ["with-admin"]` pero no hay documentación ni funciona bien.

**Mejora propuesta:** Corregir o eliminar profiles no funcionales.

**Complejidad:** Baja (30 min)

#### 9. Script de Backup/Restore

**Problema:** No hay forma fácil de respaldar datos de desarrollo.

**Mejora propuesta:** Crear `scripts/backup.sh` y `scripts/restore.sh`.

**Complejidad:** Media (2 horas)

---

## Deuda Técnica

### 1. Docker Compose Files Redundantes

**Situación:** Existen múltiples archivos docker-compose con propósitos solapados:
- `docker-compose.yml` - Principal
- `docker-compose-apps.yml` - Solo apps
- `docker-compose-infrastructure.yml` - Solo infra
- `docker-compose-mock.yml` - Con mocks
- `docker-compose.migrate.yml` - Solo migrator

**Impacto:** Confusión sobre cuál usar, mantenimiento de múltiples archivos.

**Solución propuesta:** Consolidar usando perfiles de Docker Compose v2.

### 2. Documentación en Carpeta Docker

**Situación:** Hay documentación dentro de `docker/`:
- `ACTUALIZAR_BASE_DATOS.md`
- `PLAN_PRUEBAS_DOCKER_COMPOSE.md`
- `QUICK_START.md`
- `README.md`
- `RESULTADO_VALIDACION.md`

**Impacto:** Documentación dispersa en múltiples ubicaciones.

**Solución propuesta:** Mover a `documentos/` o archivar si están desactualizados.

### 3. Migrator como Build Local

**Situación:** El migrator se construye localmente en lugar de usar imagen pre-construida.

**Impacto:** Tiempo de setup más largo, requiere contexto de build.

**Solución propuesta:** Publicar imagen del migrator en ghcr.io.

---

## Registro de Cambios Recientes

### Diciembre 2025

- Reorganización de documentación
- Creación de carpeta `documentos/` con estructura limpia
- Archivado de documentación antigua en `archivado-documentos/`
- Nuevo README.md principal simplificado
- **Nuevo:** Makefile con comandos simplificados
- **Nuevo:** Script de diagnóstico `scripts/diagnose.sh`

---

## Cómo Contribuir

Si identificas código deprecado o tienes propuestas de mejora:

1. Documenta el problema claramente
2. Propón una solución con estimación de complejidad
3. Agrega a este documento
4. Crea issue en GitHub si es significativo

---

**Última actualización:** Diciembre 2025
