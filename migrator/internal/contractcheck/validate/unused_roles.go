package validate

import (
	"fmt"
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// detectUnusedRoles implementa B-REQ-6: roles seedados que el FE no
// menciona. Severidad por defecto `warning`; eleva a `error` si el rol
// tiene `scope=system` (un rol de sistema sin UI explícita es bug, no
// deuda) — B-REQ-6.2.
func detectUnusedRoles(k kmp.Snapshot, s seed.Snapshot) []Drift {
	feRoles := make(map[string]struct{}, len(k.Roles))
	for code := range k.Roles {
		feRoles[code] = struct{}{}
	}

	type meta struct {
		scope string
	}
	grouped := make(map[string]meta)
	order := make([]string, 0)
	seen := make(map[string]struct{})
	for _, r := range s.Roles {
		if r.Code == "" {
			continue
		}
		if _, dup := seen[r.Code]; dup {
			continue
		}
		seen[r.Code] = struct{}{}
		if _, used := feRoles[r.Code]; used {
			continue
		}
		grouped[r.Code] = meta{scope: r.Scope}
		order = append(order, r.Code)
	}
	sort.Strings(order)

	out := make([]Drift, 0, len(order))
	for _, code := range order {
		m := grouped[code]
		severity := SeverityFor(CategoryRoleUnused)
		detail := fmt.Sprintf(
			"El seed declara role.code %q pero ningún composable KMP lo atiende.",
			code,
		)
		if m.scope == "system" {
			severity = SeverityError
			detail += " scope=system: un rol de sistema sin UI explícita es bloqueante."
		}
		out = append(out, Drift{
			Direction:  DirectionBEOnly,
			Category:   CategoryRoleUnused,
			Severity:   severity,
			Identifier: code,
			Detail:     detail,
		})
	}
	return out
}
