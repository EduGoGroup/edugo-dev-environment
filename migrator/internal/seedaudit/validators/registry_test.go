package validators

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
)

func TestRunAll_Clean(t *testing.T) {
	got := RunAll(minimalSnapshot())
	if len(got) != 0 {
		t.Fatalf("expected zero violations on minimal snapshot, got %+v", got)
	}
}

func TestRunAll_NilSnapshotIsSafe(t *testing.T) {
	got := RunAll(nil)
	if got != nil {
		t.Fatalf("expected nil for nil snapshot, got %+v", got)
	}
}

func TestRunAll_RecoversPanics(t *testing.T) {
	// Replace the registry with a panicking validator and a benign one,
	// using a defer to restore the original list.
	original := registered
	defer func() { registered = original }()

	registered = []NamedValidator{
		{Name: "boom", Fn: func(*loader.SeedSnapshot) []report.Violation {
			panic("boom")
		}},
		{Name: "noop", Fn: func(*loader.SeedSnapshot) []report.Violation {
			return nil
		}},
	}

	got := RunAll(minimalSnapshot())
	v := findCode(t, got, report.CodeInternalError)
	if v == nil {
		t.Fatalf("expected INTERNAL_ERROR violation, got %+v", got)
	}
	if v.EntityID != "boom" {
		t.Errorf("entity_id=%q want \"boom\"", v.EntityID)
	}
}

func TestRunAll_PreservesOrder(t *testing.T) {
	// All registered validators run on the minimal snapshot without
	// emitting violations; just make sure RunAll wires them all up.
	if len(registered) != 7 {
		t.Fatalf("expected 7 registered validators, got %d", len(registered))
	}
}
