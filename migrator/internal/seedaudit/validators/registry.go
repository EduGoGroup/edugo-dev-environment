package validators

import (
	"fmt"
	"runtime/debug"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
)

// Validator is the canonical signature shared by every static check.
type Validator func(*loader.SeedSnapshot) []report.Violation

// NamedValidator pairs a Validator with a stable identifier; the order
// of `registered` defines the run order for RunAll.
type NamedValidator struct {
	Name string
	Fn   Validator
}

// registered enumerates the validators in the order the registry must
// invoke them. Reorder with care: golden files and reports depend on
// stable emission order before the post-sort step in the reporter.
var registered = []NamedValidator{
	{Name: "permissions", Fn: ValidatePermissions},
	// P4-1 (plan B): el validador `role_permissions` fue retirado
	// porque la tabla iam.role_permissions ya no existe. Se reemplaza
	// en P4-2 por un validador sobre iam.role_grants.
	{Name: "resource_screens", Fn: ValidateResourceScreens},
	{Name: "slot_data", Fn: ValidateSlotData},
	{Name: "concepts", Fn: ValidateConcepts},
	{Name: "inverse_coverage", Fn: ValidateInverseCoverage},
	{Name: "menu_hierarchy", Fn: ValidateMenuHierarchy},
}

// RunAll invokes every registered validator over snap and aggregates
// their violations in declaration order. A panic in any single
// validator is recovered, recorded as an INTERNAL_ERROR violation, and
// the remaining validators continue running.
func RunAll(snap *loader.SeedSnapshot) []report.Violation {
	if snap == nil {
		return nil
	}

	out := make([]report.Violation, 0)
	for _, v := range registered {
		out = append(out, runOne(snap, v)...)
	}
	return out
}

func runOne(snap *loader.SeedSnapshot, v NamedValidator) (vs []report.Violation) {
	defer func() {
		if r := recover(); r != nil {
			vs = []report.Violation{
				{
					Severity: report.SeverityFor(report.CodeInternalError),
					Code:     report.CodeInternalError,
					Message:  fmt.Sprintf("El validador %q entró en panic: %v", v.Name, r),
					Entity:   "Validator",
					EntityID: v.Name,
					References: map[string]string{
						"validator": v.Name,
						"stack":     string(debug.Stack()),
					},
				},
			}
		}
	}()
	return v.Fn(snap)
}
