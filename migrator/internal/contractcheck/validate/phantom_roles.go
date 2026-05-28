package validate

import (
	"fmt"
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// detectPhantomRoles implementa B-REQ-5: roles que el FE menciona en
// `dashboard-{role}` o en `when (role.code)` y que el seed no declara.
// Severidad fija: `error`.
func detectPhantomRoles(k kmp.Snapshot, s seed.Snapshot) []Drift {
	seedRoles := make(map[string]struct{}, len(s.Roles))
	for _, r := range s.Roles {
		if r.Code != "" {
			seedRoles[r.Code] = struct{}{}
		}
	}

	codes := make([]string, 0, len(k.Roles))
	for code := range k.Roles {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	out := make([]Drift, 0)
	for _, code := range codes {
		if _, ok := seedRoles[code]; ok {
			continue
		}
		out = append(out, Drift{
			Direction:  DirectionFEOnly,
			Category:   CategoryRolePhantom,
			Severity:   SeverityFor(CategoryRolePhantom),
			Identifier: code,
			Detail: fmt.Sprintf(
				"El FE referencia role.code %q (literal o sufijo dashboard-) pero el seed no lo declara en iam.roles.",
				code,
			),
			Evidence: dedupLocations(k.Roles[code]),
		})
	}
	return out
}
