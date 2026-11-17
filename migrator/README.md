# Migrator - EduGo

Servicio que aplica migraciones de esquema y datos iniciales para PostgreSQL al iniciar el entorno de desarrollo.

> **Nota (plan 037, D-037.11):** MongoDB fue **retirado** del ecosistema. El migrator ya no
> aplica migraciones ni seeds de Mongo; las variables `MONGO_*`, el servicio docker `mongo` y
> los tests de integración de Mongo se eliminaron.

## Estado actual

| Base de datos | Estado |
|---|---|
| PostgreSQL | Funcional al 100% — migraciones, seeds y versionamiento |

## Dependencias

```
github.com/EduGoGroup/edugo-infrastructure/postgres v0.65.0
```

## Variables de entorno

### PostgreSQL

| Variable | Default | Descripcion |
|---|---|---|
| `POSTGRES_URI` | — | URI completa (alternativa a las variables individuales) |
| `POSTGRES_HOST` | `localhost` | Host |
| `POSTGRES_PORT` | `5432` | Puerto |
| `POSTGRES_DB` | `edugo` | Base de datos |
| `POSTGRES_USER` | `edugo` | Usuario |
| `POSTGRES_PASSWORD` | `edugo123` | Contrasena |
| `POSTGRES_SSLMODE` | `disable` | Modo SSL |

> Las variables `MONGO_*` fueron retiradas (plan 037, D-037.11).

### Flags de control

| Variable | Default | Descripcion |
|---|---|---|
| `FORCE_MIGRATION` | `false` | Elimina y recrea todas las bases de datos |
| `APPLY_MOCK_DATA` | `true` | Aplica datos de desarrollo (popula `--seed-demo` si el flag no es explícito) |
| `STATUS_ONLY` | `false` | Muestra estado actual sin aplicar cambios |

### Flags CLI del binario migrator

| Flag | Tipo | Default | Descripcion |
|---|---|---|---|
| `--seed-up-to-layer` | string | `""` (todas) | Aplica capas del sistema hasta la capa indicada. Ejemplo: `--seed-up-to-layer=legacy`. |
| `--seed-demo` | bool | `true` | Aplica seeds de desarrollo (`demo.ApplyDemo`). Si `APPLY_MOCK_DATA=false` en ENV y el flag no se pasa explícito, toma el valor de la variable de entorno. |

| `--playground-v2` | string | `""` (siembra `base`) | Recrea la BD (force) y aplica un fixture autónomo: `base`, `mp10_material`, `f6_representante` o `all` (los tres). Un escenario on-top por aquí falla listando sus usuarios. |
| `--on-top` | string | `""` | Aplica un escenario on-top **sin force** sobre una BD ya migrada a la `SchemaVersion` actual. Excluyente con los otros dos flags y con `FORCE_MIGRATION=true`. Ver «Escenarios on-top». |

> **Nota de migración**: `seeds.ApplyProduction` y `seeds.ApplyDevelopment` fueron eliminados.
> Usar `system.ApplySystem(db, upTo)` y `demo.ApplyDemo(gdb)` respectivamente.
> `MigrateOptions.ApplyMock` fue renombrado a `MigrateOptions.SeedDemo`.

## Uso con docker-compose

El migrator esta integrado en el `docker-compose.yml` principal del entorno de desarrollo:

```bash
# Levantar solo infraestructura (postgres + migrator)
docker compose up

# Levantar entorno completo (infraestructura + todas las apps)
docker compose --profile full up

# Ver logs del migrator
docker compose logs migrator
```

El migrator corre una sola vez y termina. Si las bases de datos ya tienen datos, las migraciones se omiten (comportamiento idempotente). Para forzar una recreacion completa:

```bash
FORCE_MIGRATION=true docker compose up
```

## Escenarios on-top (plan 053 T5.2)

Los seis escenarios con usuarios (`onboarding`, `n1_inscripcion`, `n17_secciones`, `n4_evaluacion`,
`multi_unidad`, `n0n1_escuelas`) ya no crean usuarios: sus usuarios los da de alta la API de identidad,
con el UUID que asigna identity-core, y el escenario se siembra encima de `base`. Solo hay targets
locales; contra Neon lo decide el dueño.

1. **Migrar con base:** `make docker-recreate`.
2. **Arrancar `edugo-api-identity` en local** contra esa BD con `AUTH_IDENTITY_SYNC_ENABLED=true` (y el
   exchange) apuntando al identity-core de desarrollo.
3. **Alta de los usuarios:** `make docker-provision-users P=<escenario>` (`YES=1` salta la
   confirmación). Por cada usuario hace el signup de `edugo-api-identity`, el de identity-core (que
   adopta la cuenta con la contraseña de prueba) y comprueba que el id de `auth.users` es el de
   identity-core.
4. **Sembrar el escenario:** `make docker-on-top P=<escenario>` (`migrator --on-top=<escenario>`). Exige
   la BD migrada a la `SchemaVersion` actual y falla listando los usuarios que falten.
5. **identity drena el outbox** (`PUT /users/{id}/systems` en identity-core) por su cuenta.

Variables del paso 3 (además de las `POSTGRES_*` de arriba, para leer `auth.users`):

| Variable | Default | Descripción |
|---|---|---|
| `EDUGO_IDENTITY_URL` | `http://localhost:8070` | Base de `edugo-api-identity` |
| `IDENTITY_CORE_BASE_URL` | — (obligatoria) | Base del identity-core de desarrollo |
| `SCENARIO_TEST_PASSWORD` | — (obligatoria) | Contraseña de prueba: solo del entorno, nunca en un fichero versionado. Entre 12 caracteres y 72 bytes tras NFKC, la política de identity-core |

Flags de `provision_users`: `-scenario` (obligatorio), `-yes`, `-interval` (2s, pausa entre signups de
identity-core), `-retry-wait` (30s, espera tras un 429; identity-core no manda `Retry-After`, y si lo
manda se respeta), `-max-retries` (5) y `-timeout` (15s por petición).

Cómo lee cada respuesta:

- `edugo-api-identity` 201 = creado; 409 `EMAIL_ALREADY_EXISTS` = ya existía (re-ejecución, no es fallo);
  409 `IDENTITY_USER_MISMATCH` o 502 `IDENTITY_SYNC_FAILED` = fallo, y no se llama a identity-core para
  ese usuario.
- identity-core 201 `{"id"}` = cuenta adoptada (caso B) o ya reconocida (caso C): el cuerpo es el mismo y
  el resumen los cuenta juntos. 409 = otra contraseña (caso D) o cuenta bloqueada o inactiva: se reporta,
  no se reintenta (sumaría al lockout) y se sigue con el resto. 429 = se reintenta.
- Si el id de `auth.users` no coincide con el de identity-core, el sync estaba apagado: la herramienta
  **aborta** sin tocar a los usuarios que quedan. Recrea la BD y repite con el sync encendido.

El resumen final cuenta creados, ya existían, identity-core 201 y fallidos; sale con código distinto de
0 si algo falló. La contraseña no aparece en ningún log.

## Tests de integracion

Los tests usan testcontainers para crear contenedores temporales de PostgreSQL. Docker debe estar corriendo.

```bash
# Ejecutar todos los tests
go test -v ./tests

# Solo PostgreSQL
go test -v ./tests -run TestPostgresIntegration
```

## Compilar y ejecutar local

```bash
# Descargar dependencias
go mod download

# Compilar
go build -o bin/migrator ./cmd

# Ejecutar (requiere PostgreSQL corriendo)
./bin/migrator
```

## Auditor estático del seed (`seed-audit`)

Binario complementario que carga el production seed en memoria, ejecuta
siete validadores estáticos (referenciales, cobertura inversa, jerarquía
de menú, slot_data) y persiste un reporte JSON + Markdown bajo
`audit-reports/`. No abre conexión a PostgreSQL ni modifica el seed.

Spec completa: [`EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-a-static-auditor/`](../../../EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-a-static-auditor/).

### Cuándo correrlo

- Al modificar `seeds/system/l4/*.go` o `seeds/system/layers/l*_*.go` —
  para detectar FK rotas, duplicados, huérfanos o `slot_data` que
  referencia permisos inexistentes.
- Antes de mergear cambios al seed (target `seed-audit-strict` falla CI
  si aparecen violaciones de severidad `error`).
- Periódicamente, para comparar contra el `audit-reports/baseline/` y
  detectar regresiones (Fase E cubrirá esto en CI).

### Comandos

| Target                  | Qué hace                                                                  |
|-------------------------|---------------------------------------------------------------------------|
| `make seed-audit`       | Modo `--report-only`. Genera reportes en `audit-reports/`, exit 0 siempre |
| `make seed-audit-strict`| Modo `--strict`. Exit 1 si hay violaciones `error`. Pensado para CI       |
| `make seed-audit-test`  | Tests del auditor; falla si la cobertura agregada cae bajo 80 %           |

### Flags del binario

```text
seed-audit [--seed-source=production] [--output-dir=./audit-reports]
           [--format=json|md|both] [--strict] [--report-only] [--version]
```

- `--strict` y `--report-only` son mutuamente excluyentes; gana `--strict`
  con un aviso a stderr (Decisión D-7).
- Exit codes: `0` éxito (o `--report-only`), `1` violaciones bajo
  `--strict`, `2` error interno (seed no compila, output-dir no escribible,
  flag inválido).

### Lectura de un reporte

`audit-reports/seed-audit-<timestamp>.md` agrupa:

- **Estadísticas** del seed cargado (counts por colección).
- **Resumen** de violaciones por severidad y por código.
- **Violaciones** agrupadas por severidad, con código, entidad afectada,
  mensaje en español y referencias adicionales.

El catálogo completo de códigos vive en
[`internal/seedaudit/report/codes.md`](internal/seedaudit/report/codes.md).

### Baseline

`audit-reports/baseline/seed-audit-baseline.{json,md}` se commitea como
referencia inmutable del estado conocido del seed. El resto de
`audit-reports/*` queda excluido por `.gitignore`.

## Cross-checker FE↔BE (`contract-check`)

Binario complementario que detecta drift entre el contrato del frontend
KMP (`screenKey`, `apiPrefix`, `requiredPermission`, roles) y el seed de
producción (`resource_screens`, permisos, roles, `slot_data`). No abre
conexión a PostgreSQL ni modifica código.

Spec completa: [`EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-b-contract-checker/`](../../../EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-b-contract-checker/).

### Cuándo correrlo

- Al introducir un nuevo `screenKey` en KMP — para validar que el seed lo declara.
- Al modificar `apiPrefix` o `resource` en un `ContractDecl` del KMP.
- Al añadir/eliminar roles, permisos o `resource_screens` en el seed.
- Antes de mergear cambios al seed o al kmp (target `contract-check-strict` falla CI si hay errores).

### Comandos

| Target                              | Qué hace                                                                |
|-------------------------------------|-------------------------------------------------------------------------|
| `make contract-check`               | Reporte informativo con TODAS las severidades. Exit 0 siempre.          |
| `make contract-check-strict`        | Filtra `--severity=error`. Exit 1 si hay drifts error después del filtro. CI/precommit. |
| `make contract-check-update-baseline` | Refresca `audit-reports/baseline/contract-check-baseline.json`.       |
| `make contract-check-test`          | Tests; falla si la cobertura agregada cae bajo 80 %.                   |

### Detecciones (7 categorías)

| Categoría | Severidad | Qué detecta |
|-----------|-----------|-------------|
| `screen_key_phantom`     | error    | KMP referencia un `screenKey` que el seed no declara. |
| `screen_key_dead`        | warning¹ | El seed declara un `screen_key` sin implementación KMP. |
| `permission_phantom`     | error²   | KMP infiere/literaliza un permiso que el seed no contiene. |
| `permission_zombie`      | info³    | Permiso seedado sin role_permission ni referencia FE/slot_data. |
| `role_phantom`           | error    | KMP menciona un role.code que el seed no declara. |
| `role_unused`            | warning⁴ | Rol seedado que el FE nunca atiende. |
| `service_prefix_mismatch`| error⁵   | El `apiPrefix` declarado por el FE no coincide con la tabla canónica. |

¹ Escala a error si `screen_type=dashboard` o `is_default=true`.
² Degrada a warning si el resource existe pero la acción inferida no.
³ Sube a warning si el permiso aparece en role_permissions pero no en FE.
⁴ Escala a error si `scope=system`.
⁵ Cae a info si el resource no está clasificado en la tabla.

### Flags del binario

```text
contract-check [--kmp-roots=path1,path2,...] [--severity=error|warning|info]
               [--output-dir=./audit-reports] [--update-baseline]
               [--seed-source=production]
```

- Las paths default de `--kmp-roots` son **relativas al monorepo `EduGo/`**;
  por eso los targets `make` cambian a la raíz antes de ejecutar.
- Exit codes: `0` éxito (o sin filtro de severidad), `1` drifts error
  bajo `--severity=error`, `2` error de entorno (seed loader, kmp roots).

### Lectura del reporte

`audit-reports/contract-check-<timestamp>.md` agrupa los drifts por categoría,
cada uno con su evidencia (archivos KMP + número de línea). El JSON
acompañante está pensado para diff y consumo programático.

### Baseline

`audit-reports/baseline/contract-check-baseline.json` se commitea como
referencia inmutable. Cada corrida calcula `Regresiones` (drifts nuevos
vs. baseline) y `Fixes` (drifts del baseline ya resueltos), que aparecen
en el Markdown como secciones dedicadas.

## Referencias

- [docs/DOCKER_INTEGRATION_GUIDE.md](docs/DOCKER_INTEGRATION_GUIDE.md) — guia detallada de integracion con docker-compose
