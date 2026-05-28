# Notas del baseline post-L4 del cross-checker FE↔BE

**Generado**: 2026-05-11
**Reporte fuente**: `contract-check-20260511T234657Z.{json,md}` (timestamp original)
**Versión del schema**: 1.0.0
**Versión del binario**: `contract-check` (Fase 6 — rebuild seeds L0..L4)
**Baseline previo**: `audit-reports/baseline/contract-check-baseline.{json,md}` (2026-05-08)

## Resumen

| Severidad | Conteo |
|-----------|--------|
| Errores   | 27 |
| Warnings  | 68 |
| Infos     | 3 |

**Diff vs baseline previo**: 0 drifts nuevos, 0 drifts eliminados, 0 drifts modificados. La salida es idéntica al baseline 2026-05-08 (verificado con `jq -S '.drifts' | diff`).

## Conteo por categoría (idéntico al baseline previo)

| Categoría | Drifts | Estado |
|-----------|--------|--------|
| `permission_phantom`      | 22 (warning) | Aceptado — heurística infiere CRUD sobre dashboards read-only. Ver B2-D3. |
| `permission_zombie`       | 36 (warning) | Aceptado — permisos asignados a roles activos, usados por backend JWT. Ver B2-D2. |
| `role_phantom`            | 5 (error)   | Falso positivo conocido — heurística captura sufijos `dashboard-*` que no son `role.code`. Ver TC-D (NOTES baseline previo). |
| `role_unused`             | 9 (warning) | Aceptado — 7 roles legacy + 2 roles `scope=system` sin UI dedicada. Ver B2-D1. |
| `screen_key_dead`         | 4 (warning) | Aceptado heredado del baseline previo. |
| `screen_key_phantom`      | 14 (error)  | Aceptado heredado del baseline previo. |
| `service_prefix_mismatch` | 8 (error+info) | Aceptado — tabla canónica del cross-checker obsoleta excepto `guardian_relations` (corregido en seed; ticket FE separado para `GuardianRequestsListContract.kt`). Ver B4. |

## Justificación de las decisiones en Fase 6

Las decisiones aplicadas durante la Fase 6 (rebuild de seeds L0..L4) se documentan en `decisions-log.md`:

- **B2-D1**: `role_unused` aceptados (9 warnings) — 7 roles legacy no implementados + 2 system roles `platform_admin` y `super_admin`.
- **B2-D2**: `permission_zombie` (36 warnings) — todos conservados; 32 son consumidos por backend (JWT claims/autz), 4 zombies revisados y justificados.
- **B2-D3**: `permission_phantom` (22 warnings) — NO se siembran. Acción correcta = refinar la herramienta + limpiar FE.
- **B4**: `service_prefix_mismatch` (9 warnings/errors) — 8 documentados como "seed correcto, tabla canónica del cross-checker obsoleta"; 1 corregido (`guardian_relations: learning → academic`) con ticket FE separado.

## Cómo refrescar este baseline

Igual que el baseline previo:

```bash
make contract-check
# inspeccionar audit-reports/contract-check-<timestamp>.md
make contract-check-update-baseline
git add audit-reports/baseline/
```

Si se quiere refrescar específicamente este snapshot post-L4, copiar manualmente el último reporte a `audit-reports/baseline-post-l4/`.
