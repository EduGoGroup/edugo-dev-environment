# `audit-reports/` — reportes de auditoría

Esta carpeta es la salida de los binarios de la spec
`system-data-quality-spec/` (`seed-audit` de Fase A y `contract-check`
de Fase B). Cumple dos funciones:

1. **Working area**: cada corrida deja archivos
   `seed-audit-<timestamp>.{json,md}` o `contract-check-<timestamp>.{json,md}`
   que el desarrollador inspecciona y descarta.
2. **Baseline versionado**: la subcarpeta `baseline/` contiene el
   "estado conocido" del sistema y se commitea al repo. Los runs
   futuros computan diff contra ese baseline para distinguir drifts
   nuevos (regresiones) de drifts ya aceptados.

## Convención del baseline

```
baseline/
├── seed-audit-baseline.{json,md}        # estado del seed
├── seed-audit-NOTES.md                  # análisis y clasificación
├── seed-audit-TODO.md                   # mejoras pendientes (T-1, T-2, …)
├── contract-check-baseline.{json,md}    # estado del drift FE↔BE
├── contract-check-NOTES.md              # clasificación por categoría
└── contract-check-TODO.md               # mejoras pendientes (TC-A, TC-B, …)
```

- **JSON** es la fuente de verdad estructural (lo consumen los binarios
  para computar diff).
- **Markdown** es el rendering humano (resumen ejecutivo, tablas por
  categoría, evidencia con archivos:línea).
- **NOTES.md** clasifica cada hallazgo: drift real vs. falso positivo,
  decisión tomada, justificación.
- **TODO.md** lista las acciones derivadas con IDs estables (`T-1`,
  `TC-A`) para referencia desde commits y PRs.

El resto de archivos generados (`*-<timestamp>.{json,md}`) está
ignorado por `.gitignore`. Solo `baseline/` se commitea.

## Cómo aceptar un baseline en una revisión de PR

Cuando un PR introduce intencionalmente un cambio que altera el
baseline (ej. añadir un permiso nuevo al seed, o renombrar un
screenKey en KMP):

1. Después de aplicar el cambio funcional, correr
   `make seed-audit` o `make contract-check` localmente y revisar
   el reporte fresco.
2. Si los nuevos drifts son esperados, regenerar el baseline:
   - Fase A: `make seed-audit` y mover los `*-<timestamp>` a
     `baseline/seed-audit-baseline.*` con commit explícito.
   - Fase B: `make contract-check-update-baseline`.
3. Actualizar el `NOTES.md` correspondiente con la justificación del
   delta y, si aplica, cerrar entradas en `TODO.md`.
4. Incluir la regeneración del baseline en el mismo PR como un
   commit separado con título tipo
   `chore(audit): refresh seed baseline (<motivo>)`.

Esta convención asegura que cualquier revisor pueda ver en el diff:
- El cambio funcional.
- El delta exacto del baseline (qué drifts aparecen, qué desaparecen).
- La justificación humana del delta.

## Cómo CI consume estos reportes

A partir de Fase E:

- **`make seed-audit-strict`** falla CI si el reporte contiene
  cualquier `Violation` con `severity=error`. Esto se evalúa contra
  el reporte fresco, no contra el baseline.
- **`make contract-check-strict`** filtra a severidad `error` y falla
  si quedan drifts. Se compara contra `baseline.json` para reportar
  regresiones explícitas.

Hasta que Fase E aterrice, los targets `*-strict` están disponibles
manualmente para desarrolladores y precommit hooks.

## Limpieza

Los archivos `*-<timestamp>.{json,md}` no se rotan automáticamente.
Si la carpeta crece demasiado, el desarrollador puede borrarlos con
seguridad — el baseline está protegido por la regla del `.gitignore`.

```bash
# Borrar reportes ad-hoc dejando intacto el baseline
rm -f audit-reports/seed-audit-*.{json,md}
rm -f audit-reports/contract-check-*.{json,md}
```

## Referencias

- [Spec `system-data-quality-spec`](../../../EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/README.md)
- Fase A — auditor estático del seed: [`README.md` § "Auditor estático del seed"](../README.md)
- Fase B — cross-checker FE↔BE: [`README.md` § "Cross-checker FE↔BE"](../README.md)
