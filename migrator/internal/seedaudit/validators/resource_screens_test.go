package validators

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestValidateResourceScreens_Clean(t *testing.T) {
	snap := minimalSnapshot()
	assertEmpty(t, ValidateResourceScreens(snap))
}

func TestValidateResourceScreens_DuplicateDefault(t *testing.T) {
	snap := minimalSnapshot()
	snap.ResourceScreens = append(snap.ResourceScreens, entities.ResourceScreen{
		ID:          uuidFromInt(300),
		ResourceID:  snap.ResourceScreens[0].ResourceID,
		ResourceKey: snap.ResourceScreens[0].ResourceKey,
		ScreenKey:   snap.ResourceScreens[0].ScreenKey,
		ScreenType:  snap.ResourceScreens[0].ScreenType,
		IsDefault:   true,
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateResourceScreens(snap)
	if v := findCode(t, got, report.CodeRSDuplicateDefault); v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeRSDuplicateDefault, got)
	} else if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
}

func TestValidateResourceScreens_ScreenMissing(t *testing.T) {
	snap := minimalSnapshot()
	snap.ResourceScreens = append(snap.ResourceScreens, entities.ResourceScreen{
		ID:          uuidFromInt(301),
		ResourceID:  snap.Resources[0].ID,
		ResourceKey: snap.Resources[0].Key,
		ScreenKey:   "ghost-screen",
		ScreenType:  "detail",
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateResourceScreens(snap)
	if v := findCode(t, got, report.CodeRSScreenMissing); v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeRSScreenMissing, got)
	} else if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
}

func TestValidateResourceScreens_NoDefault(t *testing.T) {
	snap := minimalSnapshot()
	// Unset the default flag so the menu-visible resource has no default screen.
	snap.ResourceScreens[0].IsDefault = false
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateResourceScreens(snap)
	v := findCode(t, got, report.CodeRSNoDefault)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeRSNoDefault, got)
	}
	if v.Severity != report.SeverityWarning {
		t.Errorf("severity=%s want warning", v.Severity)
	}
}
