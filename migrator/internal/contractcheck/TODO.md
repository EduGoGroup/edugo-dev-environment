# contractcheck — TODO de cableado pendiente

Ítems intencionalmente omitidos de la tanda Fase B inicial. Cablearlos
cuando la implementación de validate/report esté lista.

## 1.3 Target Makefile

Añadir al `Makefile` del migrator (`EduBack/edugo-dev-environment/migrator/Makefile`):

```makefile
.PHONY: contract-check contract-check-strict contract-check-update-baseline

contract-check:
	go build -o bin/contract-check ./cmd/contract-check
	./bin/contract-check

contract-check-strict:
	go build -o bin/contract-check ./cmd/contract-check
	./bin/contract-check --severity=error

contract-check-update-baseline:
	go build -o bin/contract-check ./cmd/contract-check
	./bin/contract-check --update-baseline
```

Requirements: B-REQ-9.2, B-REQ-12.

## 1.4 Sección README.md

Añadir al `migrator/README.md` un bloque "Audit tools" que documente:

- Para qué sirve `contract-check` (detectar drift FE↔BE).
- Flags soportadas (`--kmp-roots`, `--severity`, `--update-baseline`,
  `--output-dir`, `--seed-source`).
- Ubicación de los outputs (`audit-reports/contract-check-<ts>.{json,md}`).
- Convención del baseline (`audit-reports/contract-check-baseline.json`,
  versionado en git).
- Convención de exit codes (0/1/2).

Requirements: B-REQ-12.1, B-REQ-12.2.

## 3.1 Migración a `internal/seedaudit/loader`

El paquete `seed` declara hoy una interfaz local `seed.Loader` y un mock
alimentado por `testdata/seed/happy_snapshot.json` para no acoplarse al
trabajo en curso de Fase A. Cuando Fase A exponga su API estable:

1. Importar `internal/seedaudit/loader`.
2. Crear un adapter mínimo que implemente `seed.Loader` envolviendo el
   loader real.
3. Inyectarlo desde `cmd/contract-check/main.go`.
4. Mantener el mock disponible bajo `seed.NewFixtureLoader(path string)`
   para los tests.

Requirements: B-REQ-1.2, B-REQ-2.1, B-REQ-3.3, B-REQ-4.1, B-REQ-5.3, B-REQ-6.1.
