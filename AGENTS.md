# AGENTS.md — edugo-dev-environment

> Detalle local. Reglas globales del ecosistema en `../../AGENTS.md` (no las repitas).
> Guía operativa de comandos en `README.md`.

## Propósito

Ambiente de desarrollo de EduGo. Hoy la pieza activa es **`migrator/`**: migraciones de **PostgreSQL
(Neon)** + seeds E2E para levantar datos de prueba del ecosistema. (MongoDB fue retirado del
ecosistema — plan 037, D-037.11.)

La segunda pieza es **`test/`**: los tests de integración cross-API (identity, academic y platform
reales contra un Postgres de testcontainers).

## Estructura

```
edugo-dev-environment/
  migrator/    # migraciones Postgres + seeds E2E (cmd/, internal/{config,orchestrator,postgres,ontop,provision})
  test/        # tests de integración cross-API (módulo Go aparte; gate LOCAL, no de CI)
  archived/    # material del ecosistema previo (no en uso) — ver archived/README.md
```

## Cómo usar

Toda la actividad sucede dentro de `migrator/`:

```bash
cd migrator
make help     # lista de comandos
```

Detalle de comandos y flujo en `README.md`.

**Escenarios on-top (plan 053 T5.2):** los seis escenarios con usuarios no van por `--playground-v2`.
Flujo local en cinco pasos: `make docker-recreate`, arrancar `edugo-api-identity` con
`AUTH_IDENTITY_SYNC_ENABLED=true`, `make docker-provision-users P=<escenario>` (alta por API;
`cmd/provision_users`), `make docker-on-top P=<escenario>` (`migrator --on-top`, sin force) y dejar que
identity drene el outbox. Detalle en `migrator/README.md`. No hay targets cloud para el alta.

## `test/` es un gate local, no de CI

`test/go.mod` resuelve las tres APIs con `replace` a `../../edugo-api-{academic,identity,platform}`: son
ejecutables, no librerías versionadas, y sus repos son privados. Por eso el módulo solo compila con los
repos hermanos clonados al lado, como en `EduBack/`. La CI de este repo no los clona, y el workflow
`cross-api-integration.yml` que lo intentaba se borró (plan 054, ola S.1; deuda 078): la integración
cross-API se corre en la Mac.

Desde la raíz de este repo:

```bash
make cross-api-check   # gofmt + go vet + go build + govulncheck -test (sin Docker)
make cross-api-test    # los 11 paquetes de test/integration (Docker en marcha; varios minutos)
```

Lo que hacen por dentro, por si hay que correr un solo paquete:

```bash
cd test
GOWORK=off go vet -tags=integration,notest_swagger ./...
GOWORK=off ENABLE_INTEGRATION_TESTS=true \
  go test -tags=integration,notest_swagger -count=1 -timeout=1200s -p 1 ./integration/...
```

- **El tag `integration` es obligatorio.** Todos los paquetes lo llevan: sin él, `./...` no casa ningún
  paquete y `go vet`, `go build` y `govulncheck` pasan en vacío.
- **`GOWORK=off` siempre**: el `go.work` local enmascara las versiones del `go.mod`.
- Sin `ENABLE_INTEGRATION_TESTS=true` los tests salen con 0 sin ejecutar nada.
- Las APIs se compilan desde la rama que tengan en disco. Si cambia una firma de sus `cmd/builder` o el
  contrato de permisos del seed, este módulo se rompe sin que ninguna CI avise: córrelo al tocar el
  login, el contexto o los grants.

## Reglas locales clave

- Código en inglés; documentación y logs en español.
- **No tocar migraciones ni seeds sin confirmar** (regla global del ecosistema): cambian datos compartidos.
- UTC en BD y transporte; zona local solo al renderizar.
