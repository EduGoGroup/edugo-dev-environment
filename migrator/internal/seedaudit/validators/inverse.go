package validators

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/google/uuid"
)

// ValidateInverseCoverage covers A-REQ-6:
//   - RESOURCE_ORPHAN: Resource without Permissions and without
//     ResourceScreens.
//   - PERMISSION_ZOMBIE: Permission no referenciado por SlotData ni
//     RequiredPermission.
//
// P4-1 (plan B): el chequeo ROLE_NO_DEFAULT_SCREEN y el match por
// RolePermissions del PERMISSION_ZOMBIE se eliminan porque la tabla
// iam.role_permissions ya no existe. El modelo nuevo (iam.role_grants
// con patterns wildcard) hará efectivo cualquier permiso del catálogo
// que matchee algún pattern del rol; rehacer este cruce requiere un
// permission_matcher equivalente al de runtime y está fuera del alcance
// de P4-1. Se reintroducirá en P4-2 sobre el modelo nuevo.
func ValidateInverseCoverage(s *loader.SeedSnapshot) []report.Violation {
	if s == nil {
		return nil
	}

	violations := make([]report.Violation, 0)

	resourcesWithPerm := make(map[uuid.UUID]bool, len(s.Resources))
	for i := range s.Permissions {
		resourcesWithPerm[s.Permissions[i].ResourceID] = true
	}

	resourcesWithScreen := make(map[uuid.UUID]bool, len(s.Resources))
	for i := range s.ResourceScreens {
		resourcesWithScreen[s.ResourceScreens[i].ResourceID] = true
	}

	for i := range s.Resources {
		r := s.Resources[i]
		if !resourcesWithPerm[r.ID] && !resourcesWithScreen[r.ID] {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeResourceOrphan),
				Code:     report.CodeResourceOrphan,
				Message:  fmt.Sprintf("El recurso %q no tiene permisos ni resource_screens asociados.", r.Key),
				Entity:   "Resource",
				EntityID: r.ID.String(),
				References: map[string]string{
					"resource_key": r.Key,
				},
			})
		}
	}

	// PERMISSION_ZOMBIE: ahora sólo detecta permisos NO referenciados
	// por SlotData ni RequiredPermission. No considera role_grants
	// (pendiente para P4-2).
	usedPermissionNames := collectSlotPermissionRefs(s)
	for i := range s.ScreenInstances {
		si := s.ScreenInstances[i]
		if si.RequiredPermission != nil && *si.RequiredPermission != "" {
			usedPermissionNames[*si.RequiredPermission] = true
		}
	}

	for i := range s.Permissions {
		p := s.Permissions[i]
		if usedPermissionNames[p.Name] {
			continue
		}
		violations = append(violations, report.Violation{
			Severity: report.SeverityFor(report.CodePermissionZombie),
			Code:     report.CodePermissionZombie,
			Message:  fmt.Sprintf("El permiso %q no está referenciado por ninguna pantalla (slot_data ni required_permission).", p.Name),
			Entity:   "Permission",
			EntityID: p.ID.String(),
			References: map[string]string{
				"permission_name": p.Name,
			},
		})
	}

	return violations
}

// collectSlotPermissionRefs walks every parseable SlotData JSON looking
// for canonical permission references and returns a set of names. This
// is intentionally lenient: malformed JSON is skipped silently here
// (ValidateSlotData already reports it under SLOT_INVALID_JSON).
func collectSlotPermissionRefs(s *loader.SeedSnapshot) map[string]bool {
	out := make(map[string]bool)
	for i := range s.ScreenInstances {
		raw := s.ScreenInstances[i].SlotData
		if len(raw) == 0 {
			continue
		}
		var root interface{}
		if err := json.Unmarshal(raw, &root); err != nil {
			continue
		}
		collectPermsFromNode(root, out)
	}
	return out
}

func collectPermsFromNode(node interface{}, out map[string]bool) {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if kind := knownRefKeys[key]; kind == refKindPermission {
				switch raw := value.(type) {
				case string:
					if raw != "" {
						out[raw] = true
					}
				case []interface{}:
					for _, item := range raw {
						if str, ok := item.(string); ok && str != "" {
							out[str] = true
						}
					}
				}
			}
			collectPermsFromNode(value, out)
		}
	case []interface{}:
		for _, item := range v {
			collectPermsFromNode(item, out)
		}
	}
}
