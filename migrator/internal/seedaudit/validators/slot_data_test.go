package validators

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestValidateSlotData_Clean(t *testing.T) {
	snap := minimalSnapshot()
	assertEmpty(t, ValidateSlotData(snap))
}

func TestValidateSlotData_InvalidJSON(t *testing.T) {
	snap := minimalSnapshot()
	snap.ScreenInstances[0].SlotData = json.RawMessage(`{ this is not json `)
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateSlotData(snap)
	v := findCode(t, got, report.CodeSlotInvalidJSON)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeSlotInvalidJSON, got)
	}
	if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
	// On invalid JSON the walker should not fire any SLOT_REF_MISSING.
	if countCode(got, report.CodeSlotRefMissing) != 0 {
		t.Errorf("walker should not run on invalid JSON, got: %+v", got)
	}
}

func TestValidateSlotData_RefMissingPermission(t *testing.T) {
	snap := minimalSnapshot()
	snap.ScreenInstances[0].SlotData = json.RawMessage(`{"actions":[{"permission":"ghost:read"}]}`)
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateSlotData(snap)
	v := findCode(t, got, report.CodeSlotRefMissing)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeSlotRefMissing, got)
	}
	if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
	if !strings.Contains(v.Path, "permission") {
		t.Errorf("expected JSONPath to contain 'permission', got %q", v.Path)
	}
}

func TestValidateSlotData_RefMissingResource(t *testing.T) {
	snap := minimalSnapshot()
	snap.ScreenInstances[0].SlotData = json.RawMessage(`{"resource":"ghost"}`)
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateSlotData(snap)
	v := findCode(t, got, report.CodeSlotRefMissing)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeSlotRefMissing, got)
	}
}

func TestValidateSlotData_PermissionsArrayResolves(t *testing.T) {
	snap := minimalSnapshot()
	snap.ScreenInstances[0].SlotData = json.RawMessage(`{"requires":["alpha:read","ghost:write"]}`)
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateSlotData(snap)
	if countCode(got, report.CodeSlotRefMissing) != 1 {
		t.Fatalf("expected exactly one SLOT_REF_MISSING (ghost:write), got %+v", got)
	}
}

func TestValidateSlotData_NilAndEmptyAreSafe(t *testing.T) {
	snap := minimalSnapshot()
	snap.ScreenInstances[0].SlotData = nil
	snap.ScreenInstances = append(snap.ScreenInstances, entities.ScreenInstance{
		ID:        uuidFromInt(400),
		ScreenKey: "alpha-empty",
		SlotData:  json.RawMessage(``),
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateSlotData(snap)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %+v", got)
	}
}
