package validators

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestValidatePermissions_Clean(t *testing.T) {
	snap := minimalSnapshot()
	assertEmpty(t, ValidatePermissions(snap))
}

func TestValidatePermissions_PermResourceMissing(t *testing.T) {
	snap := minimalSnapshot()
	snap.Permissions = append(snap.Permissions, entities.Permission{
		ID:         uuidFromInt(100),
		Name:       "ghost:read",
		ResourceID: uuidFromInt(999),
		Action:     "read",
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidatePermissions(snap)
	v := findCode(t, got, report.CodePermResourceMissing)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodePermResourceMissing, got)
	}
	if v.Severity != report.SeverityError {
		t.Errorf("severity=%s want %s", v.Severity, report.SeverityError)
	}
	if v.EntityID != uuidFromInt(100).String() {
		t.Errorf("entity_id=%s want %s", v.EntityID, uuidFromInt(100))
	}
}

func TestValidatePermissions_PermDuplicateAction(t *testing.T) {
	snap := minimalSnapshot()
	snap.Permissions = append(snap.Permissions, entities.Permission{
		ID:         uuidFromInt(101),
		Name:       "alpha:read_again",
		ResourceID: snap.Permissions[0].ResourceID,
		Action:     snap.Permissions[0].Action,
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidatePermissions(snap)
	v := findCode(t, got, report.CodePermDuplicateAction)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodePermDuplicateAction, got)
	}
	if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
}
