# TODO — pendientes derivados del baseline del cross-checker

> Lista de mejoras y drifts a resolver detectados en la primera corrida
> de `contract-check` contra el KMP + seed reales (2026-05-08). Cada
> ítem identifica el responsable conceptual (FE / BE / herramienta).

## Prioridad alta

### ~~TC-A. Resources `assessment-questions` y `user_roles` faltan en seed~~ — RESUELTO (parcial) 2026-05-08

**Decisión aplicada**: usar permisos del recurso padre. Ningún cambio
en el seed.

- `AssessmentQuestionsListContract.kt` y `AssessmentQuestionFormContract.kt`:
  `resource = "assessment-questions"` → `resource = "assessments"`.
  El override `permissionFor` ya emitía `assessments:read/update`; ahora
  el campo `resource` coincide con la fuente de verdad.
- `UserRolesContract.kt`: `resource = "user_roles"` → `resource = "users"`.
  Se añadió `permissionFor` que emite `users:read/update` (sin override
  granular). Los handlers `assign-role` y `revoke-role` ahora exigen
  `requiredPermission = "users:update"`, no permisos `user_roles:*`
  inexistentes.

**Errores eliminados**: 9 (4 assessment-questions + 4 user_roles + 1 service_prefix asociado).

**Pendiente derivado**: la pantalla `user-roles` sigue como
`screen_key_phantom` porque el seed no la declara — eso es ahora parte
de TC-B (no de TC-A). Decisión de producto: ¿se seedará la pantalla?

### TC-A-bis. Aplicar las mismas decisiones cuando `apple_new` entre en alcance

**Contexto**: `apple_new` es un producto separado, hoy out-of-scope
(ver `e2e-integration-plan/00-out-of-scope.md §B6`). Cuando se retome,
debe replicar las mismas decisiones tomadas en este TC-A para evitar
recrear los mismos drifts en la nueva codebase.

**Cambios requeridos en `apple_new` (cuando aplique)**:

1. Pantallas de preguntas de evaluación: declarar `resource = "assessments"`
   (no `"assessment-questions"`) y heredar permisos `assessments:read/update`.
2. Pantalla de asignación/revocación de roles a usuarios: declarar
   `resource = "users"` (no `"user_roles"`) y exigir `users:update`
   en los handlers de assign/revoke (no `user_roles:create`/`user_roles:delete`).

**Razón**: el seed no tiene resources `assessment-questions` ni
`user_roles`. La política acordada es que esas pantallas heredan
permisos del recurso padre (`assessments`, `users`).

**Bloqueador para destildar**: que `apple_new` entre en scope y que el
contract-check (o equivalente) corra contra esa codebase.

**Owner**: equipo de `apple_new` cuando lo retomen.

### TC-B. 14 screen_key_phantom — pantallas KMP sin seed

**Síntoma**: 14 screenKeys hardcodeados en KMP que el seed no expone como `screen_key` en `screen_instances`. Lista exacta en `contract-check-baseline.md` sección `screen_key_phantom`.

**Decisión a tomar**: para cada uno, añadir al seed (caso normal) o eliminar la pantalla del KMP (caso de pantalla muerta).

**Owner**: Backend (cargar al seed).

### TC-C. service_prefix_mismatch — 9 errors de routing

**Síntoma**: contratos KMP que declaran un `apiPrefix` distinto al de la `serviceRoutingTable` canónica. Replica el mismo bug histórico que F2·H3.a (anuncios — `academic:` vs `platform:`).

**Decisión a tomar**: para cada uno, decidir si el FE corrige el `apiPrefix` o si la `serviceRoutingTable` (en `validate/service_prefix.go`) debe ampliarse porque la tabla está incompleta.

**Owner**: FE (corregir contratos) o herramienta (extender tabla).

### ~~TC-D. Mismatch de naming de roles~~ — DESCARTADO 2026-05-08

**Conclusión tras inspección**: NO ES UN DRIFT. El FE consulta el
rol con su nombre canónico (`hasRole("super_admin")`, con underscore)
y mapea al screenKey `"dashboard-superadmin"` (sin underscore) —
ambos consistentes con el seed. Los 5 errors `role_phantom` son
falsos positivos por heurística incorrecta del extractor (asume que
el sufijo de `dashboard-X` es un `role.code`, lo cual no se cumple
para esta codebase).

Documentado en `contract-check-NOTES.md` sección `role_phantom`.

## Prioridad media

### TC-E. 22 warnings permission_phantom — acciones inferidas inválidas

**Síntoma**: el cross-validator infiere `<resource>:{create,update,delete}` por convención y muchas acciones inferidas no están seedadas para resources read-only (`audit`, `progress`, `reports`, `stats`, dashboards).

**Opciones**:
1. **Limpieza FE**: el FE no debería pretender escribir en resources read-only. Quitar botones / inferencia de write en dashboards.
2. **Limpieza herramienta**: si un resource solo aparece en pantallas tipo `dashboard` o `detail`, no inferir crud completo. Refinar `phantom_permissions.go`.
3. **Aceptar warnings** y documentarlas en el baseline (status quo).

**Recomendación**: mezcla — refinar la herramienta para no inferir crud en dashboards, y dejar el resto como warnings reales que el FE debe limpiar.

### TC-F. 36 permission_zombie — limpieza del seed

**Síntoma**: permisos seedados que nadie consume. Posibles causas: (a) FE no implementa la pantalla todavía, (b) permiso histórico que ya no se usa.

**Owner**: Backend, en colaboración con product (decidir qué pantallas planeadas siguen vivas).

### TC-G. 9 role_unused — roles seedados sin pantalla

**Síntoma**: roles como `assistant_teacher`, `school_assistant`, `school_director`, `school_coordinator`, etc., están seedados pero el FE no implementa dashboards ni pantallas específicas para ellos. Coherente con findings.md §3.3 (FE solo implementa 6 de 12 roles).

**Decisión**: aceptar como conocido, o priorizar implementación de UI por rol.

## Prioridad baja (herramienta)

### TC-1. Falsos positivos remanentes en role_phantom (3 casos)

**Síntoma**: la regex `"dashboard-X"` captura `home`, `name`, `subjects` que NO son códigos de rol sino strings con prefijo "dashboard-" que aparecen en otros contextos:

- `LoginContract.kt:41`, `MainScreen.kt:118`, `Routes.kt:77` → string `"dashboard-home"` como ruta.
- `MembershipsListContract.kt:53`, `UnitDirectoryContract.kt:47` → string `"dashboard-name"` como tag de UI.
- `MembershipsListContract.kt:49`, `UnitDirectoryContract.kt:43` → similar para `subjects`.

**Fix propuesto**: enriquecer el extractor para distinguir entre strings con prefijo `"dashboard-"` que efectivamente son códigos de rol (típicamente concatenados en `setOf("dashboard-student", "dashboard-teacher")` o usados en `when (role.code)`) vs strings que son rutas/tags.

Una heurística simple: solo aceptar `"dashboard-X"` como rol cuando aparezca en un archivo `*Contract.kt` y el string esté dentro de un literal como `screenKey` o un set de `requiredRoles`/`screensFor`. Las rutas en archivos `Routes.kt` o `RouteRegistry.kt` deberían ignorarse.

**Owner**: herramienta (`internal/contractcheck/kmp/extract.go`).

### TC-2. Reducir capacidad de inferencia para resources que solo tienen `screen_type=dashboard`

Ver TC-E. Si un resource solo tiene pantallas tipo dashboard / detail, no inferir create/update/delete; en su lugar emitir un único drift "FE implementa lectura sobre resource X — verificar que no quiere escritura".

### TC-3. Confirmar la `serviceRoutingTable` con backend

La tabla actual (`internal/contractcheck/validate/service_prefix.go`) se construyó por inferencia desde `findings.md`. Antes de aplicar TC-C, confirmar caso por caso con el dueño del backend que la tabla es correcta. Si está incompleta, ampliarla para que casos legítimos no aparezcan como drift.

## Cómo registrar progreso

Cuando un ítem se resuelva:

1. Eliminarlo de este archivo (o marcar `~~tachado~~`).
2. Refrescar el baseline (ver `contract-check-NOTES.md`).
3. En el commit, referenciar el ID del TODO completado (ej. `chore(audit): TC-A — añadir resources assessment-questions al seed`).
