# Contract Check Report

- **Generado:** 2026-05-12T02:33:38Z
- **Schema:** 1.0.0

## Resumen ejecutivo

| Severidad | Conteo |
|-----------|--------|
| Errores   | 17 |
| Warnings  | 60 |
| Infos     | 5 |

### Conteo por categoría

| Categoría | Drifts |
|-----------|--------|
| permission_phantom | 22 |
| permission_zombie | 30 |
| role_phantom | 5 |
| role_unused | 10 |
| screen_key_dead | 3 |
| screen_key_phantom | 4 |
| service_prefix_mismatch | 8 |

### Estadísticas del input

| Origen | Métrica | Valor |
|--------|---------|-------|
| KMP    | screenKeys | 70 |
| KMP    | permisos | 21 |
| KMP    | roles | 8 |
| KMP    | contratos | 70 |
| Seed   | resources | 33 |
| Seed   | permissions | 105 |
| Seed   | roles | 13 |
| Seed   | role_permissions | 556 |
| Seed   | resource_screens | 68 |
| Seed   | screen_instances | 73 |

## Drifts por categoría

### permission_phantom

_FE consume un permiso (literal o inferido) que el seed no contiene. Degrada a warning si el resource existe pero la acción inferida no._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| warning | FE_ONLY | attendance:delete | FE infiere permiso "attendance:delete" (resource="attendance" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AttendanceBatchContract.kt:11, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AttendanceListContract.kt:17, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AttendanceSummaryContract.kt:12 |
| warning | FE_ONLY | audit:create | FE infiere permiso "audit:create" (resource="audit" action="create"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AuditListContract.kt:13 |
| warning | FE_ONLY | audit:delete | FE infiere permiso "audit:delete" (resource="audit" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AuditListContract.kt:13 |
| warning | FE_ONLY | audit:update | FE infiere permiso "audit:update" (resource="audit" action="update"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AuditListContract.kt:13 |
| warning | FE_ONLY | concept_types:delete | FE infiere permiso "concept_types:delete" (resource="concept_types" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ConceptTypesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ConceptTypesListContract.kt:13 |
| warning | FE_ONLY | grades:delete | FE infiere permiso "grades:delete" (resource="grades" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GradesFormContract.kt:11, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GradesListContract.kt:17 |
| warning | FE_ONLY | guardian_relations:create | FE infiere permiso "guardian_relations:create" (resource="guardian_relations" action="create"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsListContract.kt:11, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRequestsListContract.kt:12 |
| warning | FE_ONLY | guardian_relations:delete | FE infiere permiso "guardian_relations:delete" (resource="guardian_relations" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsListContract.kt:11, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRequestsListContract.kt:12 |
| warning | FE_ONLY | guardian_relations:update | FE infiere permiso "guardian_relations:update" (resource="guardian_relations" action="update"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsListContract.kt:11, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRequestsListContract.kt:12 |
| warning | FE_ONLY | permissions_mgmt:delete | FE infiere permiso "permissions_mgmt:delete" (resource="permissions_mgmt" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/PermissionsFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/PermissionsListContract.kt:9 |
| warning | FE_ONLY | progress:create | FE infiere permiso "progress:create" (resource="progress" action="create"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardGuardianContract.kt:4 |
| warning | FE_ONLY | progress:delete | FE infiere permiso "progress:delete" (resource="progress" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardGuardianContract.kt:4 |
| warning | FE_ONLY | reports:create | FE infiere permiso "reports:create" (resource="reports" action="create"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ReportCardContract.kt:8 |
| warning | FE_ONLY | reports:delete | FE infiere permiso "reports:delete" (resource="reports" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ReportCardContract.kt:8 |
| warning | FE_ONLY | reports:update | FE infiere permiso "reports:update" (resource="reports" action="update"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ReportCardContract.kt:8 |
| warning | FE_ONLY | screens:create | FE infiere permiso "screens:create" (resource="screens" action="create"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesListContract.kt:13, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenTemplatesListContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreensFormContract.kt:9 |
| warning | FE_ONLY | screens:delete | FE infiere permiso "screens:delete" (resource="screens" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesListContract.kt:13, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenTemplatesListContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreensFormContract.kt:9 |
| warning | FE_ONLY | screens:update | FE infiere permiso "screens:update" (resource="screens" action="update"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesListContract.kt:13, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenTemplatesListContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreensFormContract.kt:9 |
| warning | FE_ONLY | stats:create | FE infiere permiso "stats:create" (resource="stats" action="create"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSchooladminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardStudentContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSuperadminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardTeacherContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ProgressDashboardContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/StatsDashboardContract.kt:4 |
| warning | FE_ONLY | stats:delete | FE infiere permiso "stats:delete" (resource="stats" action="delete"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSchooladminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardStudentContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSuperadminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardTeacherContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ProgressDashboardContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/StatsDashboardContract.kt:4 |
| warning | FE_ONLY | stats:read | FE infiere permiso "stats:read" (resource="stats" action="read"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSchooladminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardStudentContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSuperadminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardTeacherContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ProgressDashboardContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/StatsDashboardContract.kt:4 |
| warning | FE_ONLY | stats:update | FE infiere permiso "stats:update" (resource="stats" action="update"): el resource existe en el seed pero la acción no está declarada en iam.permissions. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSchooladminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardStudentContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSuperadminContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardTeacherContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ProgressDashboardContract.kt:4, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/StatsDashboardContract.kt:4 |

### permission_zombie

_Permiso seedado que ningún role_permission asigna y que ningún slot_data ni FE referencia._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| warning | BE_ONLY | assessments:attempt | Permiso "assessments:attempt" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | assessments:view_results | Permiso "assessments:view_results" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | assessments_student:read | Permiso "assessments_student:read" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | audit:export | Permiso "audit:export" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | context:browse_schools | Permiso "context:browse_schools" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | context:browse_units | Permiso "context:browse_units" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | dashboard:view | Permiso "dashboard:view" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | grades:finalize | Permiso "grades:finalize" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | materials:delete | Permiso "materials:delete" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | materials:download | Permiso "materials:download" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | materials:publish | Permiso "materials:publish" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | materials:read | Permiso "materials:read" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | menu:full_read | Permiso "menu:full_read" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | menu:read | Permiso "menu:read" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | notifications:read | Permiso "notifications:read" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | periods:activate | Permiso "periods:activate" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | progress:read:own | Permiso "progress:read:own" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | schools:manage | Permiso "schools:manage" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | screen_templates:create | Permiso "screen_templates:create" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | screen_templates:delete | Permiso "screen_templates:delete" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | screen_templates:update | Permiso "screen_templates:update" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | stats:global | Permiso "stats:global" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | stats:school | Permiso "stats:school" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | stats:unit | Permiso "stats:unit" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | system_settings:read | Permiso "system_settings:read" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | system_settings:settings | Permiso "system_settings:settings" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | users:read:own | Permiso "users:read:own" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| warning | BE_ONLY | users:update:own | Permiso "users:update:own" está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend. |  |
| info | BE_ONLY | screen_instances:read | Permiso "screen_instances:read" seedado sin role_permissions, sin referencias en KMP ni en slot_data. Candidato a poda. |  |
| info | BE_ONLY | screen_templates:read | Permiso "screen_templates:read" seedado sin role_permissions, sin referencias en KMP ni en slot_data. Candidato a poda. |  |

### role_phantom

_FE menciona un role.code que el seed no declara._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | FE_ONLY | home | El FE referencia role.code "home" (literal o sufijo dashboard-) pero el seed no lo declara en iam.roles. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/LoginContract.kt:46, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/MainScreen.kt:118, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/navigation/RouteRegistry.kt:14, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/navigation/Routes.kt:77 |
| error | FE_ONLY | name | El FE referencia role.code "name" (literal o sufijo dashboard-) pero el seed no lo declara en iam.roles. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/MembershipsListContract.kt:53, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/UnitDirectoryContract.kt:47 |
| error | FE_ONLY | schooladmin | El FE referencia role.code "schooladmin" (literal o sufijo dashboard-) pero el seed no lo declara en iam.roles. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSchooladminContract.kt:7, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/dashboard/HybridDashboardContainer.kt:124, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/DynamicDashboardScreen.kt:35, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/DynamicDashboardScreen.kt:36, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/DynamicDashboardScreen.kt:37, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/DynamicDashboardScreen.kt:38 |
| error | FE_ONLY | subjects | El FE referencia role.code "subjects" (literal o sufijo dashboard-) pero el seed no lo declara en iam.roles. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/MembershipsListContract.kt:49, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/UnitDirectoryContract.kt:43 |
| error | FE_ONLY | superadmin | El FE referencia role.code "superadmin" (literal o sufijo dashboard-) pero el seed no lo declara en iam.roles. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/DashboardSuperadminContract.kt:7, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/dashboard/HybridDashboardContainer.kt:117, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/DynamicDashboardScreen.kt:33, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/screens/DynamicDashboardScreen.kt:34 |

### role_unused

_Rol seedado que el FE nunca atiende. Escala a error si scope=system._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | BE_ONLY | platform_admin | El seed declara role.code "platform_admin" pero ningún composable KMP lo atiende. scope=system: un rol de sistema sin UI explícita es bloqueante. |  |
| error | BE_ONLY | super_admin | El seed declara role.code "super_admin" pero ningún composable KMP lo atiende. scope=system: un rol de sistema sin UI explícita es bloqueante. |  |
| warning | BE_ONLY | announcement_viewer | El seed declara role.code "announcement_viewer" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | assistant_teacher | El seed declara role.code "assistant_teacher" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | observer | El seed declara role.code "observer" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | readonly_auditor | El seed declara role.code "readonly_auditor" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | school_admin | El seed declara role.code "school_admin" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | school_assistant | El seed declara role.code "school_assistant" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | school_coordinator | El seed declara role.code "school_coordinator" pero ningún composable KMP lo atiende. |  |
| warning | BE_ONLY | school_director | El seed declara role.code "school_director" pero ningún composable KMP lo atiende. |  |

### screen_key_dead

_El seed declara un screen_key sin implementación KMP. Escala a error si screen_type=dashboard o is_default=true._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | BE_ONLY | materials-list | El seed declara resource_screens.screen_key "materials-list" pero ningún composable KMP lo implementa. screen_type="list" is_default=true — pantalla crítica inalcanzable. |  |
| warning | BE_ONLY | announcement-form | El seed declara resource_screens.screen_key "announcement-form" pero ningún composable KMP lo implementa. |  |
| warning | BE_ONLY | material-form | El seed declara resource_screens.screen_key "material-form" pero ningún composable KMP lo implementa. |  |

### screen_key_phantom

_FE referencia un screenKey que el seed no declara; el FE quedaría sin pantalla en runtime._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | FE_ONLY | announcements-form | El frontend declara screenKey "announcements-form" pero el seed de producción no tiene ningún resource_screens.screen_key con ese valor. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AnnouncementsFormContract.kt:13 |
| error | FE_ONLY | attendance-form | El frontend declara screenKey "attendance-form" pero el seed de producción no tiene ningún resource_screens.screen_key con ese valor. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AttendanceFormContract.kt:4 |
| error | FE_ONLY | guardian-relations-list | El frontend declara screenKey "guardian-relations-list" pero el seed de producción no tiene ningún resource_screens.screen_key con ese valor. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsListContract.kt:15 |
| error | FE_ONLY | guardian_relations-form | El frontend declara screenKey "guardian_relations-form" pero el seed de producción no tiene ningún resource_screens.screen_key con ese valor. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRelationsFormAliasContract.kt:4 |

### service_prefix_mismatch

_El apiPrefix declarado por el FE no coincide con la tabla canónica resource→servicio._

| Severidad | Dirección | Identificador | Detalle | Evidencia |
|-----------|-----------|---------------|---------|-----------|
| error | MISMATCH | audit | resource "audit" declarado con apiPrefix=[identity:] en KMP; la tabla canónica espera "iam:". | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/AuditListContract.kt:13 |
| error | MISMATCH | guardian_relations | resource "guardian_relations" declarado con apiPrefix=[learning:] en KMP; la tabla canónica espera "academic:". | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/GuardianRequestsListContract.kt:12 |
| error | MISMATCH | roles | resource "roles" declarado con apiPrefix=[identity:] en KMP; la tabla canónica espera "iam:". | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/RolesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/RolesListContract.kt:9 |
| error | MISMATCH | screens | resource "screens" declarado con apiPrefix=[platform:] en KMP; la tabla canónica espera "iam:". | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenInstancesListContract.kt:13, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreenTemplatesListContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ScreensFormContract.kt:9 |
| error | MISMATCH | users | resource "users" declarado con apiPrefix=[identity:] en KMP; la tabla canónica espera "iam:". | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/UserRolesContract.kt:12, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/UsersFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/UsersListContract.kt:13 |
| info | MISMATCH | concept_types | resource "concept_types" no está clasificado en serviceRoutingTable (apiPrefix observado=[academic:]). Revisión humana: ampliar la tabla canónica. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ConceptTypesFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ConceptTypesListContract.kt:13 |
| info | MISMATCH | permissions_mgmt | resource "permissions_mgmt" no está clasificado en serviceRoutingTable (apiPrefix observado=[identity:]). Revisión humana: ampliar la tabla canónica. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/PermissionsFormContract.kt:9, EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/PermissionsListContract.kt:9 |
| info | MISMATCH | reports | resource "reports" no está clasificado en serviceRoutingTable (apiPrefix observado=[academic:]). Revisión humana: ampliar la tabla canónica. | EduUI/edugo-ui-kmp/kmp-screens/src/commonMain/kotlin/com/edugo/kmp/screens/dynamic/contracts/ReportCardContract.kt:8 |

