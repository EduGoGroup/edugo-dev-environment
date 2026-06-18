package validators

import (
	"fmt"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/google/uuid"
)

// ValidatePermissions covers A-REQ-1: every Permission must reference an
// existing Resource (PERM_RESOURCE_MISSING) and the (ResourceID, Action)
// tuple must be unique (PERM_DUPLICATE_ACTION).
func ValidatePermissions(s *loader.SeedSnapshot) []report.Violation {
	if s == nil {
		return nil
	}

	violations := make([]report.Violation, 0)

	type actionKey struct {
		resourceID uuid.UUID
		action     string
	}
	seen := make(map[actionKey]string, len(s.Permissions))

	for i := range s.Permissions {
		p := s.Permissions[i]

		if _, ok := s.ResourceByID[p.ResourceID]; !ok {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodePermResourceMissing),
				Code:     report.CodePermResourceMissing,
				Message:  fmt.Sprintf("El permiso %q referencia un recurso inexistente.", p.Name),
				Entity:   "Permission",
				EntityID: p.ID.String(),
				References: map[string]string{
					"permission_name":     p.Name,
					"missing_resource_id": p.ResourceID.String(),
				},
			})
		}

		key := actionKey{resourceID: p.ResourceID, action: p.Action}
		if prev, ok := seen[key]; ok {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodePermDuplicateAction),
				Code:     report.CodePermDuplicateAction,
				Message:  fmt.Sprintf("Permiso duplicado para la tupla (resource_id, action) entre %q y %q.", prev, p.Name),
				Entity:   "Permission",
				EntityID: p.ID.String(),
				References: map[string]string{
					"permission_name":     p.Name,
					"resource_id":         p.ResourceID.String(),
					"action":              p.Action,
					"existing_permission": prev,
				},
			})
			continue
		}
		seen[key] = p.Name
	}

	return violations
}
