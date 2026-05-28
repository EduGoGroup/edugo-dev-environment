package validators

import (
	"fmt"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/google/uuid"
)

// ValidateResourceScreens covers A-REQ-3:
//   - RS_SCREEN_MISSING when ScreenKey doesn't resolve to a ScreenInstance.
//   - RS_DUPLICATE_DEFAULT when a (ResourceID, ScreenType) tuple has more
//     than one IsDefault=true row.
//   - RS_NO_DEFAULT (warning) when an IsMenuVisible Resource lacks any
//     IsDefault=true ResourceScreen.
func ValidateResourceScreens(s *loader.SeedSnapshot) []report.Violation {
	if s == nil {
		return nil
	}

	violations := make([]report.Violation, 0)

	type defKey struct {
		resourceID uuid.UUID
		screenType string
	}
	defaultsByResourceType := make(map[defKey]string, len(s.ResourceScreens))
	hasAnyDefaultByResource := make(map[uuid.UUID]bool, len(s.Resources))

	for i := range s.ResourceScreens {
		rs := s.ResourceScreens[i]

		if _, ok := s.ScreenByKey[rs.ScreenKey]; !ok {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeRSScreenMissing),
				Code:     report.CodeRSScreenMissing,
				Message:  fmt.Sprintf("ResourceScreen referencia screen_key %q sin ScreenInstance correspondiente.", rs.ScreenKey),
				Entity:   "ResourceScreen",
				EntityID: rs.ID.String(),
				References: map[string]string{
					"resource_id":         rs.ResourceID.String(),
					"resource_key":        rs.ResourceKey,
					"missing_screen_key":  rs.ScreenKey,
					"screen_type":         rs.ScreenType,
				},
			})
		}

		if rs.IsDefault {
			hasAnyDefaultByResource[rs.ResourceID] = true

			key := defKey{resourceID: rs.ResourceID, screenType: rs.ScreenType}
			if prev, ok := defaultsByResourceType[key]; ok {
				violations = append(violations, report.Violation{
					Severity: report.SeverityFor(report.CodeRSDuplicateDefault),
					Code:     report.CodeRSDuplicateDefault,
					Message:  fmt.Sprintf("Más de un ResourceScreen marcado como default para (resource_id=%s, screen_type=%q).", rs.ResourceID, rs.ScreenType),
					Entity:   "ResourceScreen",
					EntityID: rs.ID.String(),
					References: map[string]string{
						"resource_id":               rs.ResourceID.String(),
						"screen_type":               rs.ScreenType,
						"existing_resource_screen":  prev,
						"current_resource_screen":   rs.ID.String(),
					},
				})
				continue
			}
			defaultsByResourceType[key] = rs.ID.String()
		}
	}

	for i := range s.Resources {
		r := s.Resources[i]
		if !r.IsMenuVisible {
			continue
		}
		if !hasAnyDefaultByResource[r.ID] {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeRSNoDefault),
				Code:     report.CodeRSNoDefault,
				Message:  fmt.Sprintf("El recurso %q es visible en menú pero no tiene ningún ResourceScreen marcado como default.", r.Key),
				Entity:   "Resource",
				EntityID: r.ID.String(),
				References: map[string]string{
					"resource_key": r.Key,
				},
			})
		}
	}

	return violations
}
