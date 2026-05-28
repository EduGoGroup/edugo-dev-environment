package validate

import (
	"fmt"
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// detectZombiePermissions implementa B-REQ-4: permisos seedados que
// ningún role_permission asigna y que el FE no consume (ni como literal
// ni vía inferencia canónica) ni aparecen referenciados en slot_data.
//
// Severidad:
//   - `info` cuando el permiso está completamente huérfano (no asignado
//     a ningún rol) — B-REQ-4.1.
//   - `warning` cuando el permiso sí está en role_permissions pero no
//     hay UI/slot que lo consuma — B-REQ-4.2.
func detectZombiePermissions(k kmp.Snapshot, s seed.Snapshot) []Drift {
	fePerms, _ := frontendPermissionSet(k)

	rolePerms := make(map[string]struct{}, len(s.RolePermissions))
	for _, rp := range s.RolePermissions {
		if rp.PermissionCode != "" {
			rolePerms[rp.PermissionCode] = struct{}{}
		}
	}

	slotPerms := make(map[string]struct{})
	for _, p := range s.PermissionsReferencedInSlots() {
		slotPerms[p] = struct{}{}
	}

	codes := make([]string, 0, len(s.Permissions))
	seenCodes := make(map[string]struct{}, len(s.Permissions))
	for _, p := range s.Permissions {
		if p.Code == "" {
			continue
		}
		if _, dup := seenCodes[p.Code]; dup {
			continue
		}
		seenCodes[p.Code] = struct{}{}
		codes = append(codes, p.Code)
	}
	sort.Strings(codes)

	out := make([]Drift, 0)
	for _, code := range codes {
		if _, used := fePerms[code]; used {
			continue
		}
		if _, used := slotPerms[code]; used {
			continue
		}
		_, assigned := rolePerms[code]
		var severity Severity
		var detail string
		if assigned {
			severity = SeverityWarning
			detail = fmt.Sprintf(
				"Permiso %q está asignado a algún role_permission pero ningún composable KMP ni slot_data lo consume. Revisión humana: posible uso interno backend.",
				code,
			)
		} else {
			severity = SeverityInfo
			detail = fmt.Sprintf(
				"Permiso %q seedado sin role_permissions, sin referencias en KMP ni en slot_data. Candidato a poda.",
				code,
			)
		}
		out = append(out, Drift{
			Direction:  DirectionBEOnly,
			Category:   CategoryPermissionZombie,
			Severity:   severity,
			Identifier: code,
			Detail:     detail,
		})
	}
	return out
}
