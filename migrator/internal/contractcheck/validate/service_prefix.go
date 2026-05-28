package validate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
)

// serviceRoutingTable mapea cada `resource` declarado en BaseCrudContract
// al prefijo de servicio canónico `<servicio>:` esperado por el ruteo
// real entre microservicios (B-REQ-7.2 + findings.md §3.4).
//
// Convenciones (orden de la tabla canónica):
//   - academic: schools, units, subjects, periods, calendar, schedules,
//     attendance, grades, memberships, guardian_relations, concepts.
//   - learning: assessments, materials, courses.
//   - iam:      users, roles, permissions, audit, screens,
//     system_settings, screen_instances, screen_templates, user_roles.
//   - platform: announcements, notifications, menu, screen_config.
//
// Si un `resource` no aparece en este map, detectServicePrefixMismatch
// emite un drift `info` ("ruteo no clasificado") en lugar de `error`,
// invitando a expandir la tabla cuando aparece un nuevo dominio.
var serviceRoutingTable = map[string]string{
	// academic
	"schools":            "academic:",
	"units":              "academic:",
	"subjects":           "academic:",
	"periods":            "academic:",
	"academic_periods":   "academic:",
	"calendar":           "academic:",
	"schedules":          "academic:",
	"attendance":         "academic:",
	"grades":             "academic:",
	"memberships":        "academic:",
	"guardian_relations": "academic:",
	"concepts":           "academic:",

	// learning
	"assessments": "learning:",
	"materials":   "learning:",
	"courses":     "learning:",

	// iam / identity
	"users":             "iam:",
	"roles":             "iam:",
	"permissions":       "iam:",
	"user_roles":        "iam:",
	"audit":             "iam:",
	"screens":           "iam:",
	"screen_instances":  "iam:",
	"screen_templates":  "iam:",
	"system_settings":   "iam:",

	// platform
	"announcements":  "platform:",
	"notifications":  "platform:",
	"menu":           "platform:",
	"screen_config":  "platform:",
}

// detectServicePrefixMismatch implementa B-REQ-7. Itera sobre cada
// ContractDecl con apiPrefix definido y compara contra la tabla
// canónica:
//   - error: el resource está en la tabla y el prefijo difiere.
//   - info:  el resource no está en la tabla — no se puede saber.
//   - sin drift: prefijos coinciden.
//
// El identifier del drift es `<resource>` (no el screenKey) para que
// múltiples contratos con el mismo error agreguen evidencia bajo una
// sola entrada.
func detectServicePrefixMismatch(k kmp.Snapshot) []Drift {
	type bucket struct {
		expected string
		found    map[string]struct{}
		ev       []kmp.Location
	}
	mismatches := make(map[string]*bucket)
	mismatchOrder := make([]string, 0)

	type infoBucket struct {
		prefixes map[string]struct{}
		ev       []kmp.Location
	}
	unclassified := make(map[string]*infoBucket)
	infoOrder := make([]string, 0)

	for _, c := range k.Contracts {
		resource := c.Resource
		prefix := c.APIPrefix
		if resource == "" || prefix == "" {
			continue
		}
		// Normalizar: aceptar "academic" y "academic:" como equivalentes.
		if !strings.HasSuffix(prefix, ":") {
			prefix += ":"
		}

		expected, known := serviceRoutingTable[resource]
		if !known {
			b, ok := unclassified[resource]
			if !ok {
				b = &infoBucket{prefixes: map[string]struct{}{}}
				unclassified[resource] = b
				infoOrder = append(infoOrder, resource)
			}
			b.prefixes[prefix] = struct{}{}
			b.ev = append(b.ev, c.File)
			continue
		}
		if prefix == expected {
			continue
		}
		b, ok := mismatches[resource]
		if !ok {
			b = &bucket{expected: expected, found: map[string]struct{}{}}
			mismatches[resource] = b
			mismatchOrder = append(mismatchOrder, resource)
		}
		b.found[prefix] = struct{}{}
		b.ev = append(b.ev, c.File)
	}

	out := make([]Drift, 0, len(mismatches)+len(unclassified))

	sort.Strings(mismatchOrder)
	for _, resource := range mismatchOrder {
		b := mismatches[resource]
		found := setToSorted(b.found)
		out = append(out, Drift{
			Direction:  DirectionMismatch,
			Category:   CategoryServicePrefixMismatch,
			Severity:   SeverityError,
			Identifier: resource,
			Detail: fmt.Sprintf(
				"resource %q declarado con apiPrefix=%v en KMP; la tabla canónica espera %q.",
				resource, found, b.expected,
			),
			Evidence: dedupLocations(b.ev),
		})
	}

	sort.Strings(infoOrder)
	for _, resource := range infoOrder {
		b := unclassified[resource]
		found := setToSorted(b.prefixes)
		out = append(out, Drift{
			Direction:  DirectionMismatch,
			Category:   CategoryServicePrefixMismatch,
			Severity:   SeverityInfo,
			Identifier: resource,
			Detail: fmt.Sprintf(
				"resource %q no está clasificado en serviceRoutingTable (apiPrefix observado=%v). Revisión humana: ampliar la tabla canónica.",
				resource, found,
			),
			Evidence: dedupLocations(b.ev),
		})
	}
	return out
}

func setToSorted(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
