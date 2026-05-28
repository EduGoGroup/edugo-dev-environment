# TODO — pendientes derivados del baseline del auditor

> Lista de mejoras detectadas durante la primera corrida del auditor
> estático (`seed-audit`) sobre el production seed. Estos ítems **no**
> son urgentes: el seed pasa `--strict` hoy. Son seguimiento de calidad
> para reducir falsos positivos y enriquecer el catálogo.

## Prioridad media

### T-1. Refinar validadores para distinguir contenedores de menú

**Síntoma**: las 6 advertencias del baseline se concentran en 3
recursos (`admin`, `academic`, `content`) que son contenedores del
menú, no recursos accionables. El validador los marca como
`RESOURCE_ORPHAN` y `RS_NO_DEFAULT` aunque no son inconsistencias
reales (ver `NOTES.md`).

**Opciones**:

1. **Heurística en el validador** (preferida en esta sesión).
   Considerar contenedor a un `Resource` que cumple **todas**:
   - tiene al menos un hijo (otro `Resource` con `parent_id` apuntando a él),
   - no tiene ningún `Permission` asociado,
   - no tiene ningún `ResourceScreen`.

   Tocaría `validators/inverse.go` y `validators/resource_screens.go`
   para excluirlos de las dos reglas.

2. **Marca explícita en el modelo**. Añadir `Resource.IsContainer bool`
   en `edugo-infrastructure/postgres/entities` + actualizar el seed.
   Más limpio pero impacta `edugo-infrastructure` y migraciones.

**Recomendación**: opción 1. Es el cambio mínimo y reversible.

**Bloquea a**: nada. Conveniente antes de Fase E (CI) para que el
baseline de CI sea "0 violaciones".

**Ubicación esperada del fix**: `internal/seedaudit/validators/inverse.go`
y `internal/seedaudit/validators/resource_screens.go`. Tests en sus
respectivos `*_test.go`.

## Prioridad baja

### T-2. Validar que screen_template_id (cuando aplique) resuelve

`design.md §5.1` lista "resource_screens con screen_template_id inválido"
como riesgo. El catálogo de códigos v1 no cubre este caso explícitamente.
Si en el seed real existen `ResourceScreen` con `screen_template_id`
poblado, añadir un nuevo código `RS_TEMPLATE_MISSING` y su validador.

**Verificar primero**: ¿el seed actual usa `screen_template_id`? Si
todos son `nil`, este TODO no aplica.

### T-3. Documentar el delta con findings.md

`findings.md §1` dice 47 ScreenInstances. El baseline del auditor
reporta 75. La diferencia sugiere que findings.md está desactualizado
o cuenta solo un subconjunto. Reconciliar y actualizar findings.md o
añadir una nota explicando la diferencia.

## Cómo registrar progreso

Cuando un ítem se complete:

1. Eliminarlo de este archivo.
2. Refrescar el baseline (ver `NOTES.md` "Cómo refrescar este baseline").
3. En el commit, referenciar el ID del TODO completado.
