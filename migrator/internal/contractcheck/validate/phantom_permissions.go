package validate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// canonicalActions enumera las acciones que el FE infiere automáticamente
// para cada `resource` declarado en un BaseCrudContract (B-REQ-3.1 +
// findings.md §3.2).
var canonicalActions = []string{"read", "create", "update", "delete"}

// frontendPermissionSet construye el set unión de permisos que el FE
// declara o infiere:
//
//   - literales `requiredPermission = "..."` capturados por el extractor
//     KMP (B-REQ-3.2).
//   - inferencia canónica `<resource>:{read,create,update,delete}` por
//     cada ContractDecl con resource no vacío (B-REQ-3.1).
//
// El segundo retorno mapea cada permiso a las locations donde se origina
// (literal explícito o ContractDecl.File para los inferidos), de modo
// que el reporte pueda mostrar evidencia concreta.
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
	for _, c := range k.Contracts {
		if c.Resource == "" {
			continue
		}
		for _, action := range canonicalActions {
			code := c.Resource + ":" + action
			perms[code] = struct{}{}
			evidence[code] = append(evidence[code], c.File)
		}
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
