# Migrator - EduGo

Servicio que aplica migraciones de esquema y datos iniciales para PostgreSQL y MongoDB al iniciar el entorno de desarrollo.

## Estado actual

| Base de datos | Estado |
|---|---|
| PostgreSQL | Funcional al 100% — migraciones, seeds y versionamiento |
| MongoDB | Conecta correctamente, aplica migraciones y seeds |

## Dependencias

```
github.com/EduGoGroup/edugo-infrastructure/postgres v0.65.0
github.com/EduGoGroup/edugo-infrastructure/mongodb  v0.55.0
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

### MongoDB

| Variable | Default | Descripcion |
|---|---|---|
| `MONGO_URI` | — | URI completa (alternativa a las variables individuales) |
| `MONGO_HOST` | `localhost` | Host |
| `MONGO_PORT` | `27017` | Puerto |
| `MONGO_USER` | `edugo` | Usuario |
| `MONGO_PASSWORD` | `edugo123` | Contrasena |
| `MONGO_DB_NAME` | `edugo` | Base de datos |

### Flags de control

| Variable | Default | Descripcion |
|---|---|---|
| `FORCE_MIGRATION` | `false` | Elimina y recrea todas las bases de datos |
| `APPLY_MOCK_DATA` | `true` | Aplica datos de desarrollo (popula `--seed-demo` si el flag no es explícito) |
| `POSTGRES_ONLY` | `false` | Ejecuta solo migraciones de PostgreSQL |
| `MONGO_ONLY` | `false` | Ejecuta solo migraciones de MongoDB |
| `STATUS_ONLY` | `false` | Muestra estado actual sin aplicar cambios |

### Flags CLI del binario migrator

| Flag | Tipo | Default | Descripcion |
|---|---|---|---|
| `--seed-up-to-layer` | string | `""` (todas) | Aplica capas del sistema hasta la capa indicada. Ejemplo: `--seed-up-to-layer=legacy`. |
| `--seed-demo` | bool | `true` | Aplica seeds de desarrollo (`demo.ApplyDemo`). Si `APPLY_MOCK_DATA=false` en ENV y el flag no se pasa explícito, toma el valor de la variable de entorno. |

> **Nota de migración**: `seeds.ApplyProduction` y `seeds.ApplyDevelopment` fueron eliminados.
> Usar `system.ApplySystem(db, upTo)` y `demo.ApplyDemo(gdb)` respectivamente.
> `MigrateOptions.ApplyMock` fue renombrado a `MigrateOptions.SeedDemo`.

## Uso con docker-compose

El migrator esta integrado en el `docker-compose.yml` principal del entorno de desarrollo:

```bash
# Levantar solo infraestructura (postgres, mongodb, rabbitmq + migrator)
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

## Tests de integracion

Los tests usan testcontainers para crear contenedores temporales de PostgreSQL y MongoDB. Docker debe estar corriendo.

```bash
# Ejecutar todos los tests
go test -v ./tests

# Solo PostgreSQL
go test -v ./tests -run TestPostgresIntegration

# Solo MongoDB
go test -v ./tests -run TestMongoDBIntegration
```

## Compilar y ejecutar local

```bash
# Descargar dependencias
go mod download

# Compilar
go build -o bin/migrator ./cmd

# Ejecutar (requiere PostgreSQL y MongoDB corriendo)
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
- [docs/MONGODB_LIMITATIONS.md](docs/MONGODB_LIMITATIONS.md) — limitaciones conocidas de MongoDB y soluciones propuestas
