# Notas del baseline del auditor estático del seed

**Generado**: 2026-05-08
**Reporte fuente**: `seed-audit-baseline.{json,md}` (timestamp original `20260508T220558Z`)
**Versión del schema**: 1.0.0
**Versión del binario**: `seed-audit 0.1.0-dev`

## Resumen

- **Errores**: 0
- **Advertencias**: 6 (todas clasificadas como **aceptadas conocidas** — ver abajo)
- **Informativos**: 0

El seed pasaría `make seed-audit-strict` hoy.

## Clasificación de cada violación

Las 6 advertencias se reducen a 3 recursos × 2 reglas. Los tres son
**contenedores de menú** que por diseño no tienen permisos ni pantalla
default propia; existen para agrupar a sus hijos en la barra lateral.

| Resource    | EntityID                              | Códigos disparados                        | Razón estructural                                                    | Clasificación        |
|-------------|---------------------------------------|-------------------------------------------|----------------------------------------------------------------------|----------------------|
| `admin`     | `20000000-0000-0000-0000-000000000002`| `RESOURCE_ORPHAN`, `RS_NO_DEFAULT`        | Contenedor del menú "Administración" (8 hijos: usuarios, escuelas, roles, permisos, templates, instancias, auditoría, tipos de concepto). | Aceptada conocida    |
| `academic`  | `20000000-0000-0000-0000-000000000003`| `RESOURCE_ORPHAN`, `RS_NO_DEFAULT`        | Contenedor del menú "Académico" (10 hijos: unidades, miembros, materias, períodos, calificaciones, asistencia, horarios, anuncios, calendario, vínculos guardian). | Aceptada conocida    |
| `content`   | `20000000-0000-0000-0000-000000000004`| `RESOURCE_ORPHAN`, `RS_NO_DEFAULT`        | Contenedor del menú "Contenido" (3 hijos: materiales, evaluaciones, tomar evaluación). | Aceptada conocida    |

## Por qué no se fixea ahora

1. **No son inconsistencias del seed.** El validador actual no distingue
   un recurso accionable de un nodo contenedor de menú; ambos tienen
   `is_menu_visible=true` y forma idéntica en la tabla `resources`.
2. **El fix correcto vive en el validador, no en el seed.** Cambiar el
   seed (ej. desactivar `is_menu_visible`) rompería la navegación del
   FE. Cambiar el modelo (ej. añadir `IsContainer`) toca el repo
   `edugo-infrastructure` y excede el alcance de la Fase A.
3. **Documentado vs. silenciado.** Aceptarlas como conocidas en este
   baseline preserva la señal: si en el futuro aparece una **séptima**
   violación o cambia la naturaleza de las existentes, el diff contra
   este baseline lo evidencia.

## Acciones pendientes (no para esta sesión)

Ver `audit-reports/baseline/TODO.md` para el seguimiento.

## Cómo refrescar este baseline

Cuando una violación se resuelva (o se acepte una nueva):

```bash
make seed-audit
# inspeccionar audit-reports/seed-audit-<timestamp>.md
mv audit-reports/seed-audit-<timestamp>.json audit-reports/baseline/seed-audit-baseline.json
mv audit-reports/seed-audit-<timestamp>.md   audit-reports/baseline/seed-audit-baseline.md
# editar este NOTES.md para reflejar el cambio
git add audit-reports/baseline/
git commit -m "chore(audit): refresh seed baseline (<motivo>)"
```
