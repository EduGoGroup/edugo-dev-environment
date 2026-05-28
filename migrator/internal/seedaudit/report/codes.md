# Catálogo de códigos de violación

Cada constante `CodeXxx` declarada en `codes.go` corresponde a un
patrón de inconsistencia detectable por el auditor estático. La
severidad efectiva se obtiene con `report.SeverityFor(code)`.

| Código                         | Severidad | Requirement | Qué detecta |
|--------------------------------|-----------|-------------|-------------|
| `PERM_RESOURCE_MISSING`        | error     | A-REQ-1.2   | `Permission.ResourceID` apunta a un `Resource` que no existe en el seed. |
| `PERM_DUPLICATE_ACTION`        | error     | A-REQ-1.3   | Dos `Permission` distintas declaran el mismo `<resource>:<action>`. |
| `ROLE_PERM_ROLE_MISSING`       | error     | A-REQ-2.2   | `RolePermission.RoleID` apunta a un `Role` inexistente. |
| `ROLE_PERM_PERMISSION_MISSING` | error     | A-REQ-2.2   | `RolePermission.PermissionID` apunta a un `Permission` inexistente. |
| `ROLE_PERM_DUPLICATE`          | error     | A-REQ-2.3   | Dos `RolePermission` para el mismo par `(role, permission)`. |
| `RS_DUPLICATE_DEFAULT`         | error     | A-REQ-3.2   | Dos `ResourceScreen` con `IsDefault=true` para el mismo `(resource, screen_type)`. |
| `RS_SCREEN_MISSING`            | error     | A-REQ-3.3   | `ResourceScreen.ScreenKey` no resuelve a un `ScreenInstance` seedado. |
| `RS_NO_DEFAULT`                | warning   | A-REQ-3.4   | Un `Resource` visible en menú no tiene ningún `ResourceScreen` `IsDefault=true`. |
| `SLOT_INVALID_JSON`            | error     | A-REQ-4.1   | `ScreenInstance.SlotData` no parsea como JSON. El walker se omite para esa pantalla. |
| `SLOT_REF_MISSING`             | error     | A-REQ-4.5   | Una clave canónica de `slot_data` (`permission`, `permissions`, `requires`, `resource`, `resource_key`) referencia un valor inexistente en el seed. El campo `path` reporta el JSONPath simplificado. |
| `CONCEPT_TYPE_MISSING`         | error     | A-REQ-5.1   | `ConceptDefinition.ConceptTypeID` apunta a un `ConceptType` inexistente. |
| `CONCEPT_DUPLICATE_KEY`        | error     | A-REQ-5.2   | Dos `ConceptDefinition` con el mismo `(concept_type_id, term_key)`. |
| `RESOURCE_ORPHAN`              | warning   | A-REQ-6.1   | `Resource` sin permisos asociados ni `ResourceScreen`. Indica recurso definido pero no integrado. |
| `PERMISSION_ZOMBIE`            | warning   | A-REQ-6.2   | `Permission` que ningún `RolePermission` asigna y que ningún `slot_data` referencia. |
| `ROLE_NO_DEFAULT_SCREEN`       | warning   | A-REQ-6.3   | Un `Role` no tiene acceso a ninguna pantalla `IsDefault=true` en su scope. |
| `MENU_PARENT_MISSING`          | error     | A-REQ-7.1   | `Resource.ParentID` apunta a un padre que no existe. |
| `MENU_CYCLE`                   | error     | A-REQ-7.2   | Ciclo en la jerarquía de `Resource.ParentID`. Detectado vía DFS de tres colores; los IDs del ciclo se reportan en `references`. |
| `INTERNAL_ERROR`               | error     | —           | Panic capturado por `RunAll`. Indica un bug del auditor, no del seed. |

## Cómo extender el catálogo

1. Definir la constante en `codes.go`.
2. Mapear su severidad en `SeverityFor` (default es `error`).
3. Implementar la detección en algún `validators/*.go` y registrarlo
   en `validators/registry.go`.
4. Añadir un test positivo y uno negativo en `validators/*_test.go`.
5. Documentar el código en esta tabla con su requirement asociado.
