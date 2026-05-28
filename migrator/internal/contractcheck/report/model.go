package report

import (
	"sort"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

// SchemaVersion identifica el contrato del archivo JSON emitido por el
// reporter. Se versiona aparte del binario para que cualquier consumidor
// (CI, dashboards) pueda decidir si el archivo es compatible.
const SchemaVersion = "1.0.0"

// Result agrupa todo lo que el reporter persiste o entrega al diff. La
// serialización omite por defecto los snapshots completos KMP/Seed para
// reducir el tamaño del archivo (ver Compact). Stats + Summary + Drifts
// son siempre serializados.
type Result struct {
	SchemaVersion string           `json:"schema_version"`
	Timestamp     time.Time        `json:"generated_at"`
	KMPSnapshot   *kmp.Snapshot    `json:"kmp_snapshot,omitempty"`
	SeedSnapshot  *seed.Snapshot   `json:"seed_snapshot,omitempty"`
	Stats         Stats            `json:"stats"`
	Summary       Summary          `json:"summary"`
	Drifts        []validate.Drift `json:"drifts"`
	BaselineDiff  *BaselineDiff    `json:"baseline_diff,omitempty"`
}

// Stats expone los conteos del input (KMP + Seed). Son metadatos para
// diagnosticar runs vacíos engañosos (e.g. "0 KMP files found") sin
// requerir adjuntar los snapshots completos.
type Stats struct {
	KMPScreenKeys       int `json:"kmp_screen_keys"`
	KMPPermissions      int `json:"kmp_permissions"`
	KMPRoles            int `json:"kmp_roles"`
	KMPContracts        int `json:"kmp_contracts"`
	SeedResources       int `json:"seed_resources"`
	SeedPermissions     int `json:"seed_permissions"`
	SeedRoles           int `json:"seed_roles"`
	SeedRolePermissions int `json:"seed_role_permissions"`
	SeedResourceScreens int `json:"seed_resource_screens"`
	SeedScreenInstances int `json:"seed_screen_instances"`
}

// Summary resume los drifts del run por severidad y por categoría. Los
// maps se serializan con llaves ordenadas alfabéticamente (ver json.go).
type Summary struct {
	Errors     int            `json:"errors"`
	Warnings   int            `json:"warnings"`
	Infos      int            `json:"infos"`
	ByCategory map[string]int `json:"by_category"`
	BySeverity map[string]int `json:"by_severity"`
}

// BaselineDiff lista los drifts que aparecieron (regresiones) o
// desaparecieron (fixes) respecto al baseline previo.
type BaselineDiff struct {
	Regressions []validate.Drift `json:"regressions"`
	Fixes       []validate.Drift `json:"fixes"`
}

// NewResult construye un Result a partir de los snapshots y los drifts
// computados, calculando Stats y Summary de forma determinista. NO
// adjunta los snapshots completos al Result devuelto: si el caller los
// necesita en el JSON, debe asignarlos a r.KMPSnapshot / r.SeedSnapshot
// explícitamente. Esto sigue la decisión del Bloque 5 (snapshots
// omitidos por defecto para mantener el JSON manejable).
func NewResult(ts time.Time, k kmp.Snapshot, s seed.Snapshot, drifts []validate.Drift) *Result {
	r := &Result{
		SchemaVersion: SchemaVersion,
		Timestamp:     ts.UTC(),
		Drifts:        append([]validate.Drift(nil), drifts...),
	}
	r.Stats = computeStats(k, s)
	r.Summary = computeSummary(r.Drifts)
	return r
}

func computeStats(k kmp.Snapshot, s seed.Snapshot) Stats {
	return Stats{
		KMPScreenKeys:       len(k.ScreenKeys),
		KMPPermissions:      len(k.Permissions),
		KMPRoles:            len(k.Roles),
		KMPContracts:        len(k.Contracts),
		SeedResources:       len(s.Resources),
		SeedPermissions:     len(s.Permissions),
		SeedRoles:           len(s.Roles),
		SeedRolePermissions: len(s.RolePermissions),
		SeedResourceScreens: len(s.ResourceScreens),
		SeedScreenInstances: len(s.ScreenInstances),
	}
}

func computeSummary(drifts []validate.Drift) Summary {
	sum := Summary{
		ByCategory: map[string]int{},
		BySeverity: map[string]int{},
	}
	for _, d := range drifts {
		sum.ByCategory[d.Category]++
		sum.BySeverity[string(d.Severity)]++
		switch d.Severity {
		case validate.SeverityError:
			sum.Errors++
		case validate.SeverityWarning:
			sum.Warnings++
		case validate.SeverityInfo:
			sum.Infos++
		}
	}
	return sum
}

// sortedKeys devuelve las llaves de un map[string]int en orden
// alfabético. Helper compartido por la serialización JSON y la plantilla
// Markdown para garantizar determinismo (B-REQ-10).
func sortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
