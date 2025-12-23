# Plan de Trabajo - EduGo Dev Environment

Este directorio contiene el plan detallado de mejoras y correcciones para el proyecto.

## Resumen Ejecutivo

El plan está organizado en **4 fases**, ordenadas por prioridad e impacto:

| Fase | Nombre | Prioridad | Estimación |
|------|--------|-----------|------------|
| 1 | Actualización de Dependencias | Alta | 30 min |
| 2 | Mejoras de Alta Prioridad | Alta | 4-6 horas |
| 3 | Mejoras de Media Prioridad | Media | 1-2 días |
| 4 | Mejoras de Baja Prioridad y Deuda Técnica | Baja | 2-3 días |

---

## Flujo de Trabajo por Fase

### Reglas Obligatorias

1. **Rama de trabajo:** Cada fase DEBE trabajarse en una rama creada desde `dev`
   ```bash
   git checkout dev
   git pull origin dev
   git checkout -b fase-X-descripcion
   ```

2. **Actualización de documentación:** Al completar cada fase, se DEBE actualizar la documentación afectada en `documentos/` para evitar documentos desactualizados

3. **Pull Request:** Al finalizar cada fase, se DEBE crear un PR hacia `dev`
   - Incluir resumen de cambios
   - Incluir documentación actualizada
   - Solicitar revisión si aplica

### Flujo Visual

```
dev ─────────────────────────────────────────────────────────►
      │                    ▲
      │ crear rama         │ PR
      ▼                    │
      fase-1-dependencias ─┘
      
dev ─────────────────────────────────────────────────────────►
      │                    ▲
      │ crear rama         │ PR
      ▼                    │
      fase-2-health-checks ┘

... (repetir para cada fase)
```

### Nomenclatura de Ramas

| Fase | Nombre de Rama Sugerido |
|------|-------------------------|
| 1 | `fase-1-actualizacion-dependencias` |
| 2 | `fase-2-mejoras-alta-prioridad` |
| 3 | `fase-3-mejoras-media-prioridad` |
| 4 | `fase-4-deuda-tecnica` |

---

## Checklist de Cierre de Fase

Antes de crear el PR, verificar:

- [ ] Todos los cambios de código están completos
- [ ] Tests pasan (si aplica)
- [ ] Documentación en `documentos/` actualizada
- [ ] README principal actualizado (si aplica)
- [ ] Archivos DEPRECADO-MEJORAS.md actualizado (marcar como completado)
- [ ] Commit messages siguen convención (feat, fix, docs, chore)
- [ ] PR creado hacia `dev` con descripción clara

---

## Índice de Documentos

| Documento | Descripción |
|-----------|-------------|
| [FASE-1-ACTUALIZACION-DEPENDENCIAS.md](./FASE-1-ACTUALIZACION-DEPENDENCIAS.md) | Actualizar paquetes de infraestructura |
| [FASE-2-MEJORAS-ALTA-PRIORIDAD.md](./FASE-2-MEJORAS-ALTA-PRIORIDAD.md) | Health checks, seed data automático |
| [FASE-3-MEJORAS-MEDIA-PRIORIDAD.md](./FASE-3-MEJORAS-MEDIA-PRIORIDAD.md) | Variables de entorno, Apple Silicon, seeds |
| [FASE-4-DEUDA-TECNICA.md](./FASE-4-DEUDA-TECNICA.md) | CI/CD, backup/restore, consolidación |

---

## Estado Actual

| Fase | Estado |
|------|--------|
| FASE-1 | ✅ Completada |
| FASE-2 | ⏳ Pendiente |
| FASE-3 | ⏳ Pendiente |
| FASE-4 | ⏳ Pendiente |

### Versiones de Infraestructura

- **postgres:** v0.12.0 ✅
- **mongodb:** v0.11.0 ✅

---

**Creado:** Diciembre 2025
