package report

import "time"

// SchemaVersion identifies the JSON contract of AuditReport.
// Bump on breaking changes; consumers pin against this string.
const SchemaVersion = "1.0.0"

// Severity classifies a Violation. SeverityError blocks --strict;
// SeverityWarning is reported but non-blocking; SeverityInfo is
// reserved for diagnostics (unused in v1, see Decision D-3).
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Violation is a single inconsistency detected by a validator.
//
// References carries free-form context (role_id, permission_name, ...)
// to keep the report self-explanatory without forcing readers back
// into the seed source. Path is populated only when the violation
// originates inside a slot_data JSON document, using a simplified
// JSONPath such as "$.actions[2].permission".
type Violation struct {
	Severity   Severity          `json:"severity"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Entity     string            `json:"entity"`
	EntityID   string            `json:"entity_id,omitempty"`
	References map[string]string `json:"references,omitempty"`
	Path       string            `json:"path,omitempty"`
}

// Stats holds absolute counts per collection. Useful for week-over-week
// deltas and as a sanity check on the loader.
type Stats struct {
	Resources          int `json:"resources"`
	Permissions        int `json:"permissions"`
	Roles              int `json:"roles"`
	ResourceScreens    int `json:"resource_screens"`
	ScreenInstances    int `json:"screen_instances"`
	ConceptTypes       int `json:"concept_types"`
	ConceptDefinitions int `json:"concept_definitions"`
}

// Summary aggregates Violations by severity and by code.
type Summary struct {
	Errors   int            `json:"errors"`
	Warnings int            `json:"warnings"`
	Infos    int            `json:"infos"`
	ByCode   map[string]int `json:"by_code"`
}

// AuditReport is the final document persisted as JSON and Markdown.
// The structure is reproducible: identical inputs (excluding
// GeneratedAt) produce identical bytes.
type AuditReport struct {
	SchemaVersion string      `json:"schema_version"`
	GeneratedAt   time.Time   `json:"generated_at"`
	SeedSource    string      `json:"seed_source"`
	Stats         Stats       `json:"stats"`
	Summary       Summary     `json:"summary"`
	Violations    []Violation `json:"violations"`
}
