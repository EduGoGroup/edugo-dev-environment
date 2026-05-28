package report

import (
	"sort"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
)

// reportClock is the clock used by Build to stamp GeneratedAt. Tests
// override it through SetClock to keep golden files reproducible. The
// default returns the current UTC instant.
var reportClock func() time.Time = func() time.Time { return time.Now().UTC() }

// SetClock replaces the package-level clock used by Build. It returns
// the previous clock so callers can restore it (typical pattern in
// tests). The clock is process-global; do not call from production
// code paths.
func SetClock(c func() time.Time) func() time.Time {
	prev := reportClock
	reportClock = c
	return prev
}

// severityRank gives a deterministic ordering: errors first, then
// warnings, then infos. Unknown severities sort last so they remain
// visible without breaking the ordering invariant.
func severityRank(s Severity) int {
	switch s {
	case SeverityError:
		return 0
	case SeverityWarning:
		return 1
	case SeverityInfo:
		return 2
	}
	return 3
}

// Build assembles an AuditReport from a loader snapshot, the raw
// violations emitted by validators and the seed source identifier.
//
// Determinism guarantees (A-REQ-10.2/10.3):
//   - Violations are sorted by (severity rank desc, code asc,
//     entity_id asc, path asc); two calls with the same set of
//     violations in any order produce a byte-identical Violations
//     slice.
//   - Stats and Summary are computed from the inputs only — no map
//     iteration leaks into the JSON output (Summary.ByCode is a map,
//     but encoding/json sorts keys alphabetically).
//   - SchemaVersion and SeedSource are copied verbatim; GeneratedAt
//     is the only non-deterministic field and is set from
//     reportClock so tests can pin it.
//
// snap may be nil; in that case Stats is the zero value. This keeps
// the reporter usable from unit tests that build violations by hand
// without going through the loader.
func Build(snap *loader.SeedSnapshot, violations []Violation, seedSource string) *AuditReport {
	stats := statsFromSnapshot(snap)

	summary := Summary{ByCode: make(map[string]int, len(violations))}
	for _, v := range violations {
		switch v.Severity {
		case SeverityError:
			summary.Errors++
		case SeverityWarning:
			summary.Warnings++
		case SeverityInfo:
			summary.Infos++
		}
		summary.ByCode[v.Code]++
	}

	// Copy violations before sorting so we never mutate the caller's
	// slice (validators may keep a reference for their own tests).
	sorted := make([]Violation, len(violations))
	copy(sorted, violations)
	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		if ra, rb := severityRank(a.Severity), severityRank(b.Severity); ra != rb {
			return ra < rb
		}
		if a.Code != b.Code {
			return a.Code < b.Code
		}
		if a.EntityID != b.EntityID {
			return a.EntityID < b.EntityID
		}
		return a.Path < b.Path
	})

	return &AuditReport{
		SchemaVersion: SchemaVersion,
		GeneratedAt:   reportClock(),
		SeedSource:    seedSource,
		Stats:         stats,
		Summary:       summary,
		Violations:    sorted,
	}
}

// statsFromSnapshot is a nil-safe helper that derives Stats from a
// SeedSnapshot. Defined here (not in loader) so the report package
// remains the single owner of the wire model and the loader stays
// agnostic of its consumers.
func statsFromSnapshot(snap *loader.SeedSnapshot) Stats {
	if snap == nil {
		return Stats{}
	}
	return Stats{
		Resources:          len(snap.Resources),
		Permissions:        len(snap.Permissions),
		Roles:              len(snap.Roles),
		ResourceScreens:    len(snap.ResourceScreens),
		ScreenInstances:    len(snap.ScreenInstances),
		ConceptTypes:       len(snap.ConceptTypes),
		ConceptDefinitions: len(snap.ConceptDefinitions),
	}
}
