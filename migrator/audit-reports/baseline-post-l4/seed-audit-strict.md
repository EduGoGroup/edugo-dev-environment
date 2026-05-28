# Reporte de Auditoría del Seed

- **Fuente**: `production`
- **Generado**: 2026-05-12T00:16:19Z
- **Schema**: 1.0.0

## Estadísticas

| Colección | Conteo |
|---|---:|
| Resources | 31 |
| Permissions | 89 |
| Roles | 5 |
| RolePermissions | 178 |
| ResourceScreens | 64 |
| ScreenInstances | 68 |
| ConceptTypes | 5 |
| ConceptDefinitions | 50 |

## Resumen

- **Errores**: 39
- **Advertencias**: 9
- **Informativos**: 0

### Conteo por código

| Código | Total |
|---|---:|
| `PERMISSION_ZOMBIE` | 3 |
| `PERM_RESOURCE_MISSING` | 3 |
| `RESOURCE_ORPHAN` | 3 |
| `ROLE_PERM_PERMISSION_MISSING` | 19 |
| `RS_NO_DEFAULT` | 3 |
| `SLOT_REF_MISSING` | 17 |

## Violaciones

### Errores (39)

| Código | Entidad | EntityID | Mensaje | Referencias | Path |
|---|---|---|---|---|---|
| `PERM_RESOURCE_MISSING` | Permission | 6358de3d-11ef-49c0-be42-da51ffcdfbc1 | El permiso "materials:delete" referencia un recurso inexistente. | missing_resource_id=b3000000-0000-0000-0000-000000000001, permission_name=materials:delete |  |
| `PERM_RESOURCE_MISSING` | Permission | 9aba2d20-ca23-403e-b127-6bc967eec751 | El permiso "materials:publish" referencia un recurso inexistente. | missing_resource_id=b3000000-0000-0000-0000-000000000001, permission_name=materials:publish |  |
| `PERM_RESOURCE_MISSING` | Permission | 9b0c10e0-0a7b-4e73-af9a-2d5adf99790f | El permiso "materials:download" referencia un recurso inexistente. | missing_resource_id=b3000000-0000-0000-0000-000000000001, permission_name=materials:download |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 017037ed-e261-5b3e-b828-aa19259ec110 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000004 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 38d18335-8799-56db-a226-b4bd844323f0 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000004, role_id=b4000000-0001-0000-0000-000000000004 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 653cbc7d-7775-537f-8c27-765f6e05f29d | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000004, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 724c0cde-dcad-59b9-bf36-052dec996cf7 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000001 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 731371b5-29c8-5da3-8ed4-1f1c751f4536 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000001, role_id=b4000000-0001-0000-0000-000000000003 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 832f4d7d-f6f3-5611-b79d-631148557a71 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000003, role_id=b4000000-0001-0000-0000-000000000004 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 873e2e69-85dc-567b-8dd4-cedaa40cb769 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000003, role_id=b4000000-0001-0000-0000-000000000002 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 91966af3-53a9-5a3e-8194-2c299771e828 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000003, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | 921d4bd2-ee20-559d-b456-548e36d27caa | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000001, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | a7acbeb8-e6dc-50f8-a7bc-9df7dacaacf6 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000001, role_id=b4000000-0001-0000-0000-000000000001 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | b3de7866-1e1b-53b7-bb9c-ebc7ba5465e2 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | b5ac29be-1fd0-5b9b-9631-903a8ec21f98 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000003, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | b86823a6-2665-5bfd-8527-31270f095d48 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000004, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | b938b297-f586-5cf5-9515-2b30f89d9cb1 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000002 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | c3116e14-800f-53f7-97e2-55379342fd9e | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000005 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | e089f097-5ee9-5768-b2f6-30386d299e98 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000002 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | e5be2d61-8678-536e-82bf-512b3cda8524 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000002, role_id=b4000000-0001-0000-0000-000000000003 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | f0fe4966-b67a-5ef0-9c2b-e6f60a2d9ff4 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=b3000000-0000-0000-0000-000000000004, role_id=b4000000-0001-0000-0000-000000000002 |  |
| `ROLE_PERM_PERMISSION_MISSING` | RolePermission | fc567838-dea0-500b-9184-9cca6cd94a55 | Una asignación role_permission referencia un permiso inexistente. | missing_permission_id=20000000-0000-0000-0000-000000000001, role_id=b4000000-0001-0000-0000-000000000002 |  |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000024 | La pantalla "roles-list" referencia un permiso inexistente "roles:create". | missing_permission=roles:create, ref_kind=permission, screen_key=roles-list | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000024 | La pantalla "roles-list" referencia un permiso inexistente "roles:update". | missing_permission=roles:update, ref_kind=permission, screen_key=roles-list | $.actions[1].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000024 | La pantalla "roles-list" referencia un permiso inexistente "roles:delete". | missing_permission=roles:delete, ref_kind=permission, screen_key=roles-list | $.actions[2].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000025 | La pantalla "roles-form" referencia un permiso inexistente "roles:create". | missing_permission=roles:create, ref_kind=permission, screen_key=roles-form | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000025 | La pantalla "roles-form" referencia un permiso inexistente "roles:update". | missing_permission=roles:update, ref_kind=permission, screen_key=roles-form | $.actions[1].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000025 | La pantalla "roles-form" referencia un permiso inexistente "roles:delete". | missing_permission=roles:delete, ref_kind=permission, screen_key=roles-form | $.actions[2].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000026 | La pantalla "permissions-list" referencia un permiso inexistente "permissions_mgmt:create". | missing_permission=permissions_mgmt:create, ref_kind=permission, screen_key=permissions-list | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000027 | La pantalla "permissions-form" referencia un permiso inexistente "permissions_mgmt:create". | missing_permission=permissions_mgmt:create, ref_kind=permission, screen_key=permissions-form | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000035 | La pantalla "concept-types-list" referencia un permiso inexistente "concept_types:create". | missing_permission=concept_types:create, ref_kind=permission, screen_key=concept-types-list | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000035 | La pantalla "concept-types-list" referencia un permiso inexistente "concept_types:update". | missing_permission=concept_types:update, ref_kind=permission, screen_key=concept-types-list | $.actions[1].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000036 | La pantalla "concept-types-form" referencia un permiso inexistente "concept_types:create". | missing_permission=concept_types:create, ref_kind=permission, screen_key=concept-types-form | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000036 | La pantalla "concept-types-form" referencia un permiso inexistente "concept_types:update". | missing_permission=concept_types:update, ref_kind=permission, screen_key=concept-types-form | $.actions[1].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-000000000074 | La pantalla "attendance-form" referencia un permiso inexistente "attendance:update". | missing_permission=attendance:update, ref_kind=permission, screen_key=attendance-form | $.actions[1].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-0000000000d0 | La pantalla "school-concepts-list" referencia un permiso inexistente "concept_types:create". | missing_permission=concept_types:create, ref_kind=permission, screen_key=school-concepts-list | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-0000000000d0 | La pantalla "school-concepts-list" referencia un permiso inexistente "concept_types:update". | missing_permission=concept_types:update, ref_kind=permission, screen_key=school-concepts-list | $.actions[1].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-0000000000d1 | La pantalla "school-concepts-form" referencia un permiso inexistente "concept_types:create". | missing_permission=concept_types:create, ref_kind=permission, screen_key=school-concepts-form | $.actions[0].permission |
| `SLOT_REF_MISSING` | ScreenInstance | b4400000-0000-0000-0000-0000000000d1 | La pantalla "school-concepts-form" referencia un permiso inexistente "concept_types:update". | missing_permission=concept_types:update, ref_kind=permission, screen_key=school-concepts-form | $.actions[1].permission |

### Advertencias (9)

| Código | Entidad | EntityID | Mensaje | Referencias | Path |
|---|---|---|---|---|---|
| `PERMISSION_ZOMBIE` | Permission | 52011396-5981-4c59-a772-1f353d10a3e9 | El permiso "screen_templates:create" no está asignado a ningún rol ni referenciado por ninguna pantalla. | permission_name=screen_templates:create |  |
| `PERMISSION_ZOMBIE` | Permission | b6db1991-4a2c-429a-9c45-0ed177b6e3ed | El permiso "screen_templates:delete" no está asignado a ningún rol ni referenciado por ninguna pantalla. | permission_name=screen_templates:delete |  |
| `PERMISSION_ZOMBIE` | Permission | e5bf88e6-73ff-40d4-93a4-8c787d3930af | El permiso "screen_templates:update" no está asignado a ningún rol ni referenciado por ninguna pantalla. | permission_name=screen_templates:update |  |
| `RESOURCE_ORPHAN` | Resource | b4000000-0000-0000-0000-000000000002 | El recurso "admin" no tiene permisos ni resource_screens asociados. | resource_key=admin |  |
| `RESOURCE_ORPHAN` | Resource | b4000000-0000-0000-0000-000000000003 | El recurso "academic" no tiene permisos ni resource_screens asociados. | resource_key=academic |  |
| `RESOURCE_ORPHAN` | Resource | b4000000-0000-0000-0000-000000000004 | El recurso "content" no tiene permisos ni resource_screens asociados. | resource_key=content |  |
| `RS_NO_DEFAULT` | Resource | b4000000-0000-0000-0000-000000000002 | El recurso "admin" es visible en menú pero no tiene ningún ResourceScreen marcado como default. | resource_key=admin |  |
| `RS_NO_DEFAULT` | Resource | b4000000-0000-0000-0000-000000000003 | El recurso "academic" es visible en menú pero no tiene ningún ResourceScreen marcado como default. | resource_key=academic |  |
| `RS_NO_DEFAULT` | Resource | b4000000-0000-0000-0000-000000000004 | El recurso "content" es visible en menú pero no tiene ningún ResourceScreen marcado como default. | resource_key=content |  |

