# Contract Check Report

- **Generado:** 2026-01-01T00:00:00Z
- **Schema:** 1.0.0

## Resumen ejecutivo

| Severidad | Conteo |
|-----------|--------|
| Errores   | 4 |
| Warnings  | 2 |
| Infos     | 1 |

### Conteo por categoría

| Categoría | Drifts |
|-----------|--------|
| permission_phantom | 1 |
| permission_zombie | 1 |
| role_phantom | 1 |
| role_unused | 1 |
| screen_key_dead | 1 |
| screen_key_phantom | 1 |
| service_prefix_mismatch | 1 |

### Estadísticas del input

| Origen | Métrica | Valor |
|--------|---------|-------|
| KMP    | screenKeys | 2 |
| KMP    | permisos | 1 |
| KMP    | roles | 1 |
| KMP    | contratos | 1 |
| Seed   | resources | 1 |
| Seed   | permissions | 2 |
| Seed   | roles | 2 |
| Seed   | role_permissions | 0 |
| Seed   | resource_screens | 1 |
| Seed   | screen_instances | 0 |

## Drifts por categoría

### permission_phantom

_FE consume un permiso (literal o inferido) que el seed no contiene. Degrada a warning si el resource existe pero la acción inferida no._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | FE_ONLY | ghost:read | Permiso inferido por FE no existe en iam.permissions. | kmp/Foo.kt:42 |

### permission_zombie

_Permiso seedado que ningún role_permission asigna y que ningún slot_data ni FE referencia._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| info | BE_ONLY | audit:purge | Permiso seedado sin role_permission ni referencia en slot_data. |  |

### role_phantom

_FE menciona un role.code que el seed no declara._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | FE_ONLY | principal | FE menciona el rol 'principal' que no existe en iam.roles. |  |

### role_unused

_Rol seedado que el FE nunca atiende. Escala a error si scope=system._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| warning | BE_ONLY | auditor | Rol 'auditor' seedado sin uso en KMP. |  |

### screen_key_dead

_El seed declara un screen_key sin implementación KMP. Escala a error si screen_type=dashboard o is_default=true._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| warning | BE_ONLY | legacy-list | Seed declara 'legacy-list' pero KMP no atiende esa pantalla. |  |

### screen_key_phantom

_FE referencia un screenKey que el seed no declara; el FE quedaría sin pantalla en runtime._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | FE_ONLY | ghost-form | FE declara screenKey 'ghost-form' que el seed no contiene. | kmp/Foo.kt:42 |

### service_prefix_mismatch

_El apiPrefix declarado por el FE no coincide con la tabla canónica resource→servicio._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | MISMATCH | announcements | apiPrefix=academic: declarado pero el ruteo canónico es platform:. | kmp/Foo.kt:42 |

