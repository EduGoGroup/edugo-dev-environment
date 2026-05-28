package validators

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// P4-1 (plan B): los chequeos ROLE_NO_DEFAULT_SCREEN y la rama del
// PERMISSION_ZOMBIE basada en role_permissions fueron retirados. Los
// tests correspondientes (ROLE_NO_DEFAULT_SCREEN, role-link branch del
// ZOMBIE) se eliminan; se mantienen los tests de RESOURCE_ORPHAN y
// PERMISSION_ZOMBIE-via-slot_data/required_permission.

func TestValidateInverseCoverage_Clean(t *testing.T) {
	snap := minimalSnapshot()
	assertEmpty(t, ValidateInverseCoverage(snap))
}

func TestValidateInverseCoverage_ResourceOrphan(t *testing.T) {
	snap := minimalSnapshot()
	snap.Resources = append(snap.Resources, entities.Resource{
		ID:            uuidFromInt(600),
		Key:           "lonely",
		IsMenuVisible: false,
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateInverseCoverage(snap)
	v := findCode(t, got, report.CodeResourceOrphan)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeResourceOrphan, got)
	}
	if v.Severity != report.SeverityWarning {
		t.Errorf("severity=%s want warning", v.Severity)
	}
	if v.EntityID != uuidFromInt(600).String() {
		t.Errorf("entity_id=%s", v.EntityID)
	}
}

func TestValidateInverseCoverage_PermissionZombie(t *testing.T) {
	snap := minimalSnapshot()
	// Permiso nuevo que no es referenciado por ningún slot_data ni
	// required_permission: debe reportarse como zombie.
	snap.Permissions = append(snap.Permissions, entities.Permission{
		ID:         uuidFromInt(601),
		Name:       "alpha:zombie",
		ResourceID: snap.Resources[0].ID,
		Action:     "zombie",
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateInverseCoverage(snap)
	v := findCode(t, got, report.CodePermissionZombie)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodePermissionZombie, got)
	}
	if v.Severity != report.SeverityWarning {
		t.Errorf("severity=%s want warning", v.Severity)
	}
}

func TestValidateInverseCoverage_RequiredPermissionKeepsAlive(t *testing.T) {
	// Un permiso sólo referenciado vía ScreenInstance.RequiredPermission
	// NO debe reportarse como zombie.
	snap := minimalSnapshot()
	required := "alpha:read"
	snap.ScreenInstances[0].RequiredPermission = &required
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateInverseCoverage(snap)
	if v := findCode(t, got, report.CodePermissionZombie); v != nil {
		t.Fatalf("required_permission should keep the permission alive, got zombie: %+v", v)
	}
}
