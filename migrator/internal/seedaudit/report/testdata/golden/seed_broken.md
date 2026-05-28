# Reporte de Auditoría del Seed

- **Fuente**: `production`
- **Generado**: 2026-01-01T00:00:00Z
- **Schema**: 1.0.0

## Estadísticas

| Colección | Conteo |
|---|---:|
| Resources | 4 |
| Permissions | 6 |
| Roles | 2 |
| RolePermissions | 5 |
| ResourceScreens | 3 |
| ScreenInstances | 3 |
| ConceptTypes | 1 |
| ConceptDefinitions | 2 |

## Resumen

- **Errores**: 3
- **Advertencias**: 1
- **Informativos**: 1

### Conteo por código

| Código | Total |
|---|---:|
| `DIAG_NOTE` | 1 |
| `PERM_RESOURCE_MISSING` | 2 |
| `RESOURCE_ORPHAN` | 1 |
| `SLOT_REF_MISSING` | 1 |

## Violaciones

### Errores (3)

| Código | Entidad | EntityID | Mensaje | Referencias | Path |
|---|---|---|---|---|---|
| `PERM_RESOURCE_MISSING` | Permission | perm:alpha.read | Permission referencia un Resource inexistente | permission_name=alpha.read, resource_id=00000000-0000-0000-0000-0000000000aa |  |
| `PERM_RESOURCE_MISSING` | Permission | perm:beta.write | Permission referencia un Resource inexistente | permission_name=beta.write |  |
| `SLOT_REF_MISSING` | ScreenInstance | screen:dashboard | slot_data referencia un permission inexistente | permission_name=missing.permission | $.actions[2].permission |

### Advertencias (1)

| Código | Entidad | EntityID | Mensaje | Referencias | Path |
|---|---|---|---|---|---|
| `RESOURCE_ORPHAN` | Resource | resource:zeta | Resource sin permissions ni resource_screens | resource_key=zeta |  |

### Informativos (1)

| Código | Entidad | EntityID | Mensaje | Referencias | Path |
|---|---|---|---|---|---|
| `DIAG_NOTE` | Diagnostic | diag:1 | Nota informativa de ejemplo |  |  |

