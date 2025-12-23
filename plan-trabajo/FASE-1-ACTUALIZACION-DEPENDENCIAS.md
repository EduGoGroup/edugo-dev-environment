# FASE 1: Actualización de Dependencias

**Prioridad:** Alta  
**Estimación:** 30 minutos  
**Riesgo:** Bajo  

---

## Objetivo

Actualizar los paquetes de `edugo-infrastructure` a las últimas versiones disponibles.

## Estado Actual vs Objetivo

| Paquete | Versión Actual | Versión Objetivo | Delta |
|---------|----------------|------------------|-------|
| edugo-infrastructure/postgres | v0.11.1 | v0.12.0 | +1 minor |
| edugo-infrastructure/mongodb | v0.10.1 | v0.11.0 | +1 minor |

## Cambios en las Nuevas Versiones

### postgres v0.12.0

Según los commits recientes del repositorio:
- fix(postgres): agregar función `update_updated_at_column` faltante
- feat(database): FASE 1 UI Database + Feature Flags - 5 nuevas tablas:
  - `user_active_context` (011)
  - `user_favorites` (012)
  - `user_activity_log` (013)
  - `feature_flags` (014)
  - `feature_flag_overrides` (015)

### mongodb v0.11.0

- Sincronización con cambios de mock-generator
- Mejoras en entidades tipadas

---

## Pasos de Implementación

### Paso 1.1: Actualizar go.mod

**Archivo:** `migrator/go.mod`

```bash
cd migrator
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.12.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.11.0
```

### Paso 1.2: Ejecutar go mod tidy

```bash
go mod tidy
```

### Paso 1.3: Verificar compilación

```bash
go build ./...
```

### Paso 1.4: Ejecutar tests (si existen)

```bash
go test ./...
```

---

## Validación

- [ ] `go.mod` muestra versiones actualizadas
- [ ] `go mod tidy` ejecuta sin errores
- [ ] `go build ./...` compila correctamente
- [ ] Tests pasan (si existen)

---

## Rollback

Si hay problemas, revertir a versiones anteriores:

```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.11.1
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.10.1
go mod tidy
```

---

## Commit Sugerido

```
chore(migrator): actualizar edugo-infrastructure a últimas versiones

- postgres: v0.11.1 → v0.12.0
- mongodb: v0.10.1 → v0.11.0

Incluye nuevas tablas para UI Database y Feature Flags
```

---

## Flujo de Trabajo Git

### 1. Crear rama desde dev

```bash
git checkout dev
git pull origin dev
git checkout -b fase-1-actualizacion-dependencias
```

### 2. Realizar los cambios

Ejecutar los pasos de implementación descritos arriba.

### 3. Crear PR hacia dev

```bash
git add .
git commit -m "chore(migrator): actualizar edugo-infrastructure a últimas versiones"
git push origin fase-1-actualizacion-dependencias
# Crear PR en GitHub hacia dev
```

---

## Documentación a Actualizar

Al completar esta fase, actualizar los siguientes documentos:

| Documento | Cambio Requerido |
|-----------|------------------|
| `documentos/SERVICIOS.md` | Actualizar versiones de dependencias del Migrator |
| `documentos/DEPRECADO-MEJORAS.md` | Marcar esta mejora como completada |
| `documentos/ARQUITECTURA.md` | Actualizar si hay nuevas tablas relevantes |

### Checklist de Cierre

- [ ] Código completado y compilando
- [ ] `documentos/SERVICIOS.md` actualizado
- [ ] `documentos/DEPRECADO-MEJORAS.md` actualizado  
- [ ] PR creado hacia `dev`
- [ ] PR revisado y aprobado
- [ ] PR mergeado a `dev`
