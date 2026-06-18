package validate

import "github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"

// Direction classifies which side of the FE↔BE contract owns the
// drift. FEOnly = FE declares it, seed does not. BEOnly = seed declares
// it, FE does not. Mismatch = both declare it but disagree (today only
// service_prefix_mismatch lives here).
type Direction string

const (
	DirectionFEOnly   Direction = "FE_ONLY"
	DirectionBEOnly   Direction = "BE_ONLY"
	DirectionMismatch Direction = "MISMATCH"
)

// Severity follows Phase A's three-level model. Error blocks
// `--severity=error` runs (exit 1). Warning surfaces but doesn't block.
// Info is reserved for diagnostic context (unclassified routes, etc.).
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Category names (constants) — design.md §4 "Categorías canónicas".
// Each maps to one B-REQ and one detector function.
const (
	CategoryScreenKeyPhantom      = "screen_key_phantom"      // B-REQ-1
	CategoryScreenKeyDead         = "screen_key_dead"         // B-REQ-2
	CategoryPermissionPhantom     = "permission_phantom"      // B-REQ-3
	CategoryPermissionZombie      = "permission_zombie"       // B-REQ-4
	CategoryRolePhantom           = "role_phantom"            // B-REQ-5
	CategoryRoleUnused            = "role_unused"             // B-REQ-6
	CategoryServicePrefixMismatch = "service_prefix_mismatch" // B-REQ-7
)

// Drift is one detected divergence between the FE contract and the
// seed. Identifier carries the human-readable key (screenKey, permission
// code, role code, route…). Detail spells out the reason in Spanish.
// Evidence points back to KMP source lines when applicable.
type Drift struct {
	Direction  Direction      `json:"direction"`
	Category   string         `json:"category"`
	Severity   Severity       `json:"severity"`
	Identifier string         `json:"identifier"`
	Detail     string         `json:"detail"`
	Evidence   []kmp.Location `json:"evidence,omitempty"`
}
