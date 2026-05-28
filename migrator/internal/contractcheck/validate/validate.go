package validate

import (
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// severityRank ordena severidades de mayor a menor (error > warning >
// info). Se usa al ordenar la salida del validador (B-REQ-10.1).
func severityRank(s Severity) int {
	switch s {
	case SeverityError:
		return 0
	case SeverityWarning:
		return 1
	case SeverityInfo:
		return 2
	default:
		return 3
	}
}

// Validate ejecuta los siete detectores de drift sobre los snapshots
// del frontend KMP y del seed de producción y devuelve la lista
// completa de drifts ordenada de forma determinista por
// (Category asc, Severity desc, Identifier asc) — B-REQ-10.1.
//
// El método es puro: no toca disco, no muta los snapshots de entrada y
// es seguro de invocar varias veces sobre el mismo input.
func Validate(k kmp.Snapshot, s seed.Snapshot) []Drift {
	out := make([]Drift, 0, 64)
	out = append(out, detectPhantomScreenKeys(k, s)...)
	out = append(out, detectDeadScreenKeys(k, s)...)
	out = append(out, detectPhantomPermissions(k, s)...)
	out = append(out, detectZombiePermissions(k, s)...)
	out = append(out, detectPhantomRoles(k, s)...)
	out = append(out, detectUnusedRoles(k, s)...)
	out = append(out, detectServicePrefixMismatch(k)...)

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		ri, rj := severityRank(out[i].Severity), severityRank(out[j].Severity)
		if ri != rj {
			return ri < rj
		}
		return out[i].Identifier < out[j].Identifier
	})
	return out
}
