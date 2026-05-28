# `baseline/` — baseline archivado (pre-Fase-6, 2026-05-08)

> **ESTADO: ARCHIVADO**. Este baseline corresponde a la fotografía del
> sistema **antes** del rebuild de seeds (Fase 6 — capas L0..L4 +
> borrado del catálogo legacy). Se mantiene en disco como referencia
> histórica.

## Baseline vigente

El baseline vigente post-Fase-6 vive en:

```
audit-reports/baseline-post-l4/
```

Contiene:

- `contract-check-baseline.{json,md}` — snapshot del cross-checker
  FE↔BE tras aplicar L0..L4.
- `contract-check-NOTES.md` — análisis y diff vs este baseline previo.
- `seed-audit-strict.{json,md}` — snapshot del auditor estático del
  seed corrido con `--strict` (referencia para el ticket TC-5).

## Por qué no se borró este directorio

1. Permite diff explícito (cross-checker) entre estado pre-Fase-6
   (legacy aplicado) y post-Fase-6 (catálogo L0..L4 únicamente).
2. `seed-audit-baseline.{json,md}` aún es la referencia "feliz" del
   auditor antes de que TC-5 (extender accessors a L0..L3) se
   resuelva. Una vez cerrado TC-5, el baseline strict del auditor
   debería volver a 0 errors/0 warnings y este directorio podrá
   archivarse en branch.

## Cómo CI consume estos archivos

CI lee `baseline-post-l4/` como fuente de verdad. Si un PR cambia
algo del seed, el cross-checker debe producir el MISMO output que
está en `baseline-post-l4/contract-check-baseline.json` (o
documentar el delta en `contract-check-NOTES.md`).

## Referencias

- `audit-reports/README.md` — convención general del directorio.
- `phase-6-layer-l4/decisions-log.md` — registro de decisiones de
  Fase 6 (incluye B7 que archivó el baseline post-L4).
