# TODO — Regression test del plan E2E previo (Bloque 6.2 de Fase C)

> Spec: `EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-c-fixtures-refactor/tasks.md` § Bloque 6.2.
> Requirement: `C-REQ-5.4` — los `.md` del plan E2E previo deben seguir pasando contra una BD seedeada con `scenario_legacy_e2e`.
> Estado actual: **scaffolding pendiente**. La infraestructura para correrlo en CI vive en Fase E.

## Objetivo

Tras aplicar `scenarios.LegacyE2E` sobre una BD limpia (testcontainers
postgres:15-alpine + production seed + `scenario_legacy_e2e`), correr
los pasos 1-3 (login + smoke API) de cada uno de estos `.md` y validar
que cada `curl` devuelve el status HTTP documentado:

- `EduUI/edugo-ui-kmp/e2e-integration-plan/05-fase-2-platform-screen-config.md`
- `EduUI/edugo-ui-kmp/e2e-integration-plan/06-fase-3-academic.md`
- `EduUI/edugo-ui-kmp/e2e-integration-plan/07-fase-4-learning.md`

Si algún paso devuelve un status distinto al documentado, la build CI
falla con un diff legible.

## Especificación detallada

### Forma del test

- **Lenguaje**: Go puro. Vivirá en
  `EduBack/edugo-infrastructure/postgres/seeds/e2e/integration/legacy_e2e_regression_test.go`
  (paquete `legacye2e_test`) bajo build tag `integration`.
- **Helper de BD**: reusar
  `seeds/e2e/internal/testdb.StartPostgres(t)` (ya existente desde Fase
  C). El helper levanta postgres:15-alpine y aplica
  migrations + production seed.
- **Aplicación del scenario**: `framework.NewComposer(reg, ...).Apply(gdb, "legacy_e2e")`.
- **Servicios HTTP**: el test asume que las 4 APIs Go
  (`identity`, `academic`, `platform`, `learning`) están corriendo
  contra la misma BD. Si el job de CI las arranca con
  `docker compose`, se pasan las URLs vía env (`IDENTITY_URL`,
  `ACADEMIC_URL`, `PLATFORM_URL`, `LEARNING_URL`) y el test las usa.
  Si las URLs no están definidas, el test se skipea con mensaje
  `set <X>_URL to enable regression`.

### Parser de pasos

Cada `.md` tiene secciones del estilo:

```markdown
### Paso 1: Login

```bash
curl -s -X POST "$IDENTITY_URL/api/v1/auth/login" -d '{"email":"...","password":"..."}'
```

Esperado: 200 OK con `access_token` no vacío.
```

El parser debe:

1. Tokenizar bloques de código bash que empiecen con `curl`.
2. Reconocer la línea inmediatamente posterior con el patrón
   `Esperado: <status> <texto>` o equivalente.
3. Sustituir variables de entorno (`$IDENTITY_URL`, etc.) y los
   tokens dinámicos (`$TOKEN`) que vengan de pasos anteriores.

Para minimizar el alcance, **sólo se procesan los primeros 3 pasos de
cada `.md`** (login + 2 smoke API), tal como dicta el bloque 6.2.

### Caso esperado por archivo

| Archivo                                | Pasos a ejecutar                                       |
|----------------------------------------|--------------------------------------------------------|
| `05-fase-2-platform-screen-config.md`  | 1. login admin → 2. GET /screen-config → 3. GET /menu |
| `06-fase-3-academic.md`                | 1. login admin → 2. GET /schools → 3. GET /units      |
| `07-fase-4-learning.md`                | 1. login user  → 2. GET /materials → 3. GET /assessments |

(Los detalles exactos de los endpoints viven en cada `.md`; el parser
los extrae literalmente.)

### Validación

El test compara `httpClient.Do(req).StatusCode` contra el `Esperado` y
falla con un diff de la forma:

```
regression failed: 06-fase-3-academic.md paso 2:
  curl: GET /api/v1/schools
  expected: 200
  got:      404
  body:     {"error":"route not found"}
```

## Por qué se deja como TODO

- Implementarlo requiere arrancar las 4 APIs Go en CI (compose o k8s);
  esa infra es el job nightly de Fase E.
- Para una sesión de cierre de Fase C, el costo (~día de trabajo) no
  cabe sin contaminar el alcance.
- El parser `.md → curl` es no trivial: hay que decidir cómo tratar
  `$TOKEN`, jq sobre respuestas, retries, etc.

Cuando Fase E arme el pipeline con las APIs corriendo, este test se
implementa desde aquí. La lección clave es que `scenario_legacy_e2e`
ya garantiza paridad bit-a-bit con los UUIDs/codes/emails históricos
(probado por `TestLegacyCompat_Parity` en `seeds/e2e/legacy_compat_integration_test.go`),
así que cualquier fallo en los `.md` posteriores será de la API, no
del seed.

## Ownership

- Spec original: Fase C (este TODO).
- Implementación: Fase E (pipeline CI con las 4 APIs).
- Mantenimiento del parser: equipo de calidad (cuando los `.md` cambien).

## Última revisión

- 2026-05-08 — TODO creado al cerrar Fase C. Resto de bloques 0..9 verdes.
