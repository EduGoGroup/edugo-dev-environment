package validators

import (
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

func TestValidateMenuHierarchy_Clean(t *testing.T) {
	snap := minimalSnapshot()

	// Add a child resource pointing at the existing alpha resource.
	parentID := snap.Resources[0].ID
	snap.Resources = append(snap.Resources, entities.Resource{
		ID:            uuidFromInt(700),
		Key:           "alpha-child",
		ParentID:      &parentID,
		IsMenuVisible: false,
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	assertEmpty(t, ValidateMenuHierarchy(snap))
}

func TestValidateMenuHierarchy_ParentMissing(t *testing.T) {
	snap := minimalSnapshot()

	missing := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	snap.Resources = append(snap.Resources, entities.Resource{
		ID:       uuidFromInt(701),
		Key:      "orphan-child",
		ParentID: &missing,
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateMenuHierarchy(snap)
	v := findCode(t, got, report.CodeMenuParentMissing)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeMenuParentMissing, got)
	}
	if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
}

func TestValidateMenuHierarchy_Cycle(t *testing.T) {
	// Build a tiny snapshot with a 2-node cycle: A -> B -> A.
	idA := uuidFromInt(800)
	idB := uuidFromInt(801)
	resources := []entities.Resource{
		{ID: idA, Key: "node-a", ParentID: &idB},
		{ID: idB, Key: "node-b", ParentID: &idA},
	}
	snap := loader.NewSnapshot(resources, nil, nil, nil, nil, nil, nil, nil)

	got := ValidateMenuHierarchy(snap)
	v := findCode(t, got, report.CodeMenuCycle)
	if v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeMenuCycle, got)
	}
	if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
	if !strings.Contains(v.References["cycle"], idA.String()) || !strings.Contains(v.References["cycle"], idB.String()) {
		t.Errorf("cycle should list both nodes, got %q", v.References["cycle"])
	}
}

func TestValidateMenuHierarchy_SelfLoopIsCycle(t *testing.T) {
	id := uuidFromInt(802)
	resources := []entities.Resource{
		{ID: id, Key: "self", ParentID: &id},
	}
	snap := loader.NewSnapshot(resources, nil, nil, nil, nil, nil, nil, nil)

	got := ValidateMenuHierarchy(snap)
	if v := findCode(t, got, report.CodeMenuCycle); v == nil {
		t.Fatalf("expected MENU_CYCLE for self-loop, got %+v", got)
	}
}
