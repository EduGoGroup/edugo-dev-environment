# AGENTS.md — edugo-dev-environment

> Detalle local. Reglas globales del ecosistema en `../../AGENTS.md` (no las repitas).
> Guía operativa de comandos en `README.md`.

## Propósito

Ambiente de desarrollo de EduGo. Hoy la pieza activa es **`migrator/`**: migraciones de **PostgreSQL
(Neon)** y **MongoDB (Atlas)** + seeds E2E para levantar datos de prueba del ecosistema.

## Estructura

```
edugo-dev-environment/
  migrator/    # migraciones Postgres/Mongo + seeds E2E (cmd/, internal/{config,orchestrator,postgres})
  archived/    # material del ecosistema previo (no en uso) — ver archived/README.md
```

## Cómo usar

Toda la actividad sucede dentro de `migrator/`:

```bash
cd migrator
make help     # lista de comandos
```

Detalle de comandos y flujo en `README.md`.

## Reglas locales clave

- Código en inglés; documentación y logs en español.
- **No tocar migraciones ni seeds sin confirmar** (regla global del ecosistema): cambian datos compartidos.
- UTC en BD y transporte; zona local solo al renderizar.
