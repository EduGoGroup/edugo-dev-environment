# Reporte de Auditoría del Seed

- **Fuente**: `production`
- **Generado**: 2026-05-08T22:05:58Z
- **Schema**: 1.0.0

## Estadísticas

| Colección | Conteo |
|---|---:|
| Resources | 33 |
| Permissions | 104 |
| Roles | 12 |
| RolePermissions | 428 |
| ResourceScreens | 63 |
| ScreenInstances | 75 |
| ConceptTypes | 5 |
| ConceptDefinitions | 50 |

## Resumen

- **Errores**: 0
- **Advertencias**: 6
- **Informativos**: 0

### Conteo por código

| Código | Total |
|---|---:|
| `RESOURCE_ORPHAN` | 3 |
| `RS_NO_DEFAULT` | 3 |

## Violaciones

### Advertencias (6)

| Código | Entidad | EntityID | Mensaje | Referencias | Path |
|---|---|---|---|---|---|
| `RESOURCE_ORPHAN` | Resource | 20000000-0000-0000-0000-000000000002 | El recurso "admin" no tiene permisos ni resource_screens asociados. | resource_key=admin |  |
| `RESOURCE_ORPHAN` | Resource | 20000000-0000-0000-0000-000000000003 | El recurso "academic" no tiene permisos ni resource_screens asociados. | resource_key=academic |  |
| `RESOURCE_ORPHAN` | Resource | 20000000-0000-0000-0000-000000000004 | El recurso "content" no tiene permisos ni resource_screens asociados. | resource_key=content |  |
| `RS_NO_DEFAULT` | Resource | 20000000-0000-0000-0000-000000000002 | El recurso "admin" es visible en menú pero no tiene ningún ResourceScreen marcado como default. | resource_key=admin |  |
| `RS_NO_DEFAULT` | Resource | 20000000-0000-0000-0000-000000000003 | El recurso "academic" es visible en menú pero no tiene ningún ResourceScreen marcado como default. | resource_key=academic |  |
| `RS_NO_DEFAULT` | Resource | 20000000-0000-0000-0000-000000000004 | El recurso "content" es visible en menú pero no tiene ningún ResourceScreen marcado como default. | resource_key=content |  |

