# Notas del baseline del cross-checker FE↔BE

**Generado**: 2026-05-08
**Reporte fuente**: `contract-check-baseline.{json,md}` (timestamp original `20260508T223721Z`)
**Versión del schema**: 1.0.0
**Versión del binario**: `contract-check` (Fase B)

## Resumen

- **Errores**: 36
- **Advertencias**: 68
- **Infos**: 3

El binario falla `make contract-check-strict` con el baseline actual.

## Conteo por categoría

| Categoría | Drifts | Estado |
|-----------|--------|--------|
| `permission_phantom`     | 30 (8 error + 22 warning) | Mezcla de drifts reales y heurística |
| `permission_zombie`      | 36 (warning/info) | Real — permisos seedados sin uso |
| `role_phantom`           | 5 (5 error) | Mezcla — 2 drifts reales de naming, 3 falsos positivos remanentes |
| `role_unused`            | 9 (warning) | Real — roles del seed no atendidos por el FE |
| `screen_key_dead`        | 4 (warning) | Real — screenKeys del seed sin pantalla KMP |
| `screen_key_phantom`     | 14 (error) | Real — pantallas KMP sin entrada en el seed |
| `service_prefix_mismatch` | 9 (error) | Real — apiPrefix divergente con la tabla canónica |

## Análisis crítico

### Categorías con drifts 100 % reales

**`screen_key_phantom` (14 errors)** — pantallas implementadas en KMP que el seed no expone. Sin fix el FE las pinta pero el backend devolverá 404 al pedir su `ScreenInstance`. Fix: añadir entradas correspondientes en `seeds/system/legacy/data.go`. Lista en el reporte Markdown.

**`service_prefix_mismatch` (9 errors)** — el apiPrefix declarado por `BaseCrudContract` no coincide con la `serviceRoutingTable` del cross-validator. Ejemplos típicos: contrato declara `academic:` pero el resource pertenece al servicio `iam:` o `platform:`.

**`screen_key_dead` (4 warnings)** — el seed declara `screen_key` que ningún composable atiende. Probablemente pantallas planeadas pero no implementadas.

**`permission_zombie` (36 warnings)** — permisos seedados que ni `role_permissions`, ni `slot_data`, ni el FE referencian. Candidatos a limpieza del seed (o al revés, indicio de que algún rol no fue cableado).

**`role_unused` (9 warnings)** — roles del seed que el FE no atiende. Coherente con que el FE solo implementa 6 roles (student, teacher, guardian, admin, school_admin, superadmin) frente a los 12 seedados.

### Categorías con drifts mixtos

**`permission_phantom` (8 errors + 22 warnings)**

- **Errors reales**: `assessment-questions:×4` y `user_roles:×4`. El FE implementa pantallas para ambos resources que NO existen en el seed (`AssessmentQuestionsListContract.kt`, `UserRolesContract.kt`). Hay que decidir: añadir esos resources al seed, o renombrar el screenKey/contract en el FE.
- **Warnings**: 22 acciones inferidas (`<resource>:create/update/delete`) cuyas resources existen pero las acciones no. La mayoría son de pantallas read-only (`audit`, `progress`, `reports`, `stats`, `dashboard-*`) donde el FE muestra un botón de crear/editar que el backend no soporta. Decidir: implementar la acción en el seed o eliminar el botón en el FE.

**`role_phantom` (5 errors — TODOS falsos positivos, ver TC-D)**

Tras inspección del código real (2026-05-08), los 5 errors son
artefactos de la heurística "el sufijo de `dashboard-X` es un código
de rol". El seed usa **convenciones de naming distintas** entre
roles (`super_admin`, `school_admin` con underscore) y screenKeys
(`dashboard-superadmin`, `dashboard-schooladmin` sin underscore).

El FE en `DynamicDashboardScreen.kt:33-38` hace lo correcto:

```kotlin
context?.hasRole("super_admin") == true -> "dashboard-superadmin"
context?.hasRole("school_admin") == true -> "dashboard-schooladmin"
```

`hasRole("super_admin")` consulta el rol con su nombre canónico
(matchea el seed) y mapea al screenKey `"dashboard-superadmin"`
(también matchea el seed). No hay drift real.

Detalle por identifier:

| Identifier | Origen del falso positivo | Acción |
|------------|---------------------------|--------|
| `schooladmin`, `superadmin` | Sufijos de screenKeys `dashboard-schooladmin` y `dashboard-superadmin`, malinterpretados como `role.code`. | Aceptado conocido |
| `home`, `name`, `subjects` | Strings con prefijo `"dashboard-"` que no son ni roles ni screenKeys (rutas, tags de UI). | Aceptado conocido |

**No hay nada que cambiar en seed ni en FE para estos 5 casos**. La
heurística del extractor es la que necesita refinamiento (TC-1 o
descarte total) — fuera de alcance de esta sesión.

## Acciones aplicadas en esta sesión (2026-05-08)

| Acción | Archivos tocados | Errors eliminados |
|--------|------------------|-------------------|
| `assessment-questions` → `assessments` (resource del contrato) | `AssessmentQuestionsListContract.kt`, `AssessmentQuestionFormContract.kt` | 4 |
| `user_roles:create/delete` → `users:update` + resource del contrato a `users` | `UserRolesContract.kt` | 5 (4 permission_phantom + 1 service_prefix_mismatch) |

**Resultado**: 36 errors → 27 errors. Los permisos ahora resuelven contra resources que sí existen en el seed.

> ⚠️ **Pendiente para `apple_new`**: cuando el repo `apple_new` (producto
> separado, hoy out-of-scope — ver `e2e-integration-plan/00-out-of-scope.md
> §B6`) entre en alcance, **debe aplicarse el mismo cambio**:
>
> - `UserRolesContract.kt` (o su equivalente nativo iOS): `requiredPermission`
>   de los handlers `assign-role` y `revoke-role` debe ser `users:update`,
>   no `user_roles:create`/`user_roles:delete` (esos permisos no existen
>   en el seed y no se planea seedarlos).
> - `AssessmentQuestionsListContract.kt` y `AssessmentQuestionFormContract.kt`
>   (o sus equivalentes): `resource = "assessments"`, no
>   `"assessment-questions"`.
>
> La razón es la misma que aplicó al KMP: el seed no tiene resources
> `user_roles` ni `assessment-questions` y el FE debe heredar permisos
> del recurso padre (`users`, `assessments`).

## Acciones recomendadas próximas

1. **TC-B**: 14 screen_key_phantom — caso por caso decidir si se seedan o se eliminan del FE. Una de ellas (`user-roles`) es la pantalla que toca el cambio de arriba; al seedarla se cubre tanto TC-B como cierra el ciclo del cambio anterior.
2. **TC-C**: 8 service_prefix_mismatch — auditar la `serviceRoutingTable` contra el routing real del backend.

## Cómo refrescar este baseline

Cuando una categoría se resuelva:

```bash
make contract-check
# inspeccionar audit-reports/contract-check-<timestamp>.md
make contract-check-update-baseline
git add audit-reports/baseline/
git commit -m "chore(audit): refresh contract-check baseline (<motivo>)"
```
