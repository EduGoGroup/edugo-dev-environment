package validate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// frontendPermissionSet construye el set de permisos que el FE declara
// con literales `requiredPermission = "..."` (B-REQ-3.2). El segundo
// retorno mapea cada permiso a las locations donde aparece.
//
// Ya no infiere `<resource>:{read,create,update,delete}` por cada
// contrato (antigua B-REQ-3.1): desde el ADR 0003 el permiso de una
// pantalla sale solo de `slot.permission` del seed, y en KMP `resource`
// solo da el `formScreenKey()` por defecto. La inferencia producía
// fantasmas como `records:read` (deuda 194).
func frontendPermissionSet(k kmp.Snapshot) (map[string]struct{}, map[string][]kmp.Location) {
	perms := make(map[string]struct{})
	evidence := make(map[string][]kmp.Location)
	for code, locs := range k.Permissions {
		if code == "" {
			continue
		}
		perms[code] = struct{}{}
		evidence[code] = append(evidence[code], locs...)
	}
	return perms, evidence
}

// detectPhantomPermissions implementa B-REQ-3:
//
//   - severity=error si el resource del permiso ni siquiera existe en
//     iam.permissions / iam.resources del seed.
//   - severity=warning si el resource sí existe pero la acción concreta
//     no fue seedada (sobre-aproximación de la inferencia canónica).
func detectPhantomPermissions(k kmp.Snapshot, s seed.Snapshot) []Drift {
	seedPerms := make(map[string]struct{}, len(s.Permissions))
	seedResources := make(map[string]struct{})
	for _, p := range s.Permissions {
		if p.Code == "" {
			continue
		}
		seedPerms[p.Code] = struct{}{}
		if resource, _, ok := splitPermission(p.Code); ok {
			seedResources[resource] = struct{}{}
		}
	}
	for _, r := range s.Resources {
		if r.Key != "" {
			seedResources[r.Key] = struct{}{}
		}
	}

	fePerms, evidence := frontendPermissionSet(k)
	codes := make([]string, 0, len(fePerms))
	for code := range fePerms {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	out := make([]Drift, 0)
	for _, code := range codes {
		if _, ok := seedPerms[code]; ok {
			continue
		}
		resource, action, valid := splitPermission(code)
		severity := SeverityError
		var detail string
		if valid {
			if _, ok := seedResources[resource]; ok {
				severity = SeverityWarning
				detail = fmt.Sprintf(
					"FE infiere permiso %q (resource=%q action=%q): el resource existe en el seed pero la acción no está declarada en iam.permissions.",
					code, resource, action,
				)
			} else {
				detail = fmt.Sprintf(
					"FE consume permiso %q pero el resource %q no existe en el seed (ni en iam.resources ni en iam.permissions).",
					code, resource,
				)
			}
		} else {
			detail = fmt.Sprintf(
				"FE consume permiso %q (formato no canónico) y el seed no lo declara.",
				code,
			)
		}
		out = append(out, Drift{
			Direction:  DirectionFEOnly,
			Category:   CategoryPermissionPhantom,
			Severity:   severity,
			Identifier: code,
			Detail:     detail,
			Evidence:   dedupLocations(evidence[code]),
		})
	}
	return out
}

// splitPermission extrae (resource, action) del código canónico
// `resource:action`. Retorna ok=false si el separador no existe o si
// alguna parte queda vacía.
func splitPermission(code string) (string, string, bool) {
	idx := strings.IndexByte(code, ':')
	if idx <= 0 || idx == len(code)-1 {
		return "", "", false
	}
	return code[:idx], code[idx+1:], true
}
