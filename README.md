# edugo-dev-environment

Repositorio del ambiente de desarrollo de EduGo. Hoy contiene una sola pieza activa: `migrator/`.

## Estructura

```
edugo-dev-environment/
├── migrator/      # Migraciones de PostgreSQL (Neon) + seeds E2E
└── archived/      # Material del ecosistema previo (no en uso) — ver archived/README.md
```

## Empezar

Toda la actividad sucede dentro de `migrator/`:

```bash
cd migrator
make help
```

Comandos clave (ejecutar desde `migrator/`):

| Comando | Descripción |
|---|---|
| `make build` | Compilar el binario |
| `make docker-migrate` | Levantar Postgres local y aplicar migraciones |
| `make cloud-migrate` | Aplicar migraciones en Neon (requiere `migrator/.env.cloud`) |
| `make cloud-status` | Ver estado de la BD cloud |
| `make cloud-recreate` | Recrear BDs cloud desde cero (BORRA DATOS) |
| `make cloud-seed-scenario SCENARIO=legacy_e2e` | Aplicar fixtures E2E (scenario canónico) |
| `make cloud-seed-layer LAYER=legacy` | Aplicar capas del sistema hasta la indicada |
| `make check` | Pipeline completo: fmt + vet + lint + test + build |

Documentación detallada en [`migrator/README.md`](./migrator/README.md).

## Estado del ecosistema

- **APIs por dominio** (identity, academic, learning, platform): **desplegadas en Cloud Run** desde
  2026-07-03. `edugo-api-messaging` es la 5ª API, con BD propia aislada.
- **`edugo-worker`**: vivo y vigente (imagen lista en Artifact Registry); su despliegue en Cloud Run está
  pospuesto por costo, no archivado.
- **APIs viejas** (`api-mobile`, `api-administracion`): deprecadas. El stack `docker-compose` que las
  levantaba está en [`archived/docker/`](./archived/docker/) y no representa el estado actual.

Dentro de este repo, `migrator/` es la superficie activa (migraciones + seeds). La orquestación de las
APIs vive en `../Makefile` + `../process-compose.yaml` (backend local) y en Cloud Run (nube).
