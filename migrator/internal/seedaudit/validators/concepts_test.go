package validators

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestValidateConcepts_Clean(t *testing.T) {
	snap := minimalSnapshot()
	assertEmpty(t, ValidateConcepts(snap))
}

func TestValidateConcepts_TypeMissing(t *testing.T) {
	snap := minimalSnapshot()
	snap.ConceptDefinitions = append(snap.ConceptDefinitions, entities.ConceptDefinition{
		ID:            uuidFromInt(500),
		ConceptTypeID: uuidFromInt(999),
		TermKey:       "ghost.term",
		TermValue:     "Ghost",
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateConcepts(snap)
	if v := findCode(t, got, report.CodeConceptTypeMissing); v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeConceptTypeMissing, got)
	} else if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
}

func TestValidateConcepts_DuplicateKey(t *testing.T) {
	snap := minimalSnapshot()
	snap.ConceptDefinitions = append(snap.ConceptDefinitions, entities.ConceptDefinition{
		ID:            uuidFromInt(501),
		ConceptTypeID: snap.ConceptDefinitions[0].ConceptTypeID,
		TermKey:       snap.ConceptDefinitions[0].TermKey,
		TermValue:     "Otra",
	})
	snap = loader.NewSnapshot(snap.Resources, snap.Permissions, snap.Roles, snap.ResourceScreens, snap.ScreenInstances, snap.ScreenTemplates, snap.ConceptTypes, snap.ConceptDefinitions)

	got := ValidateConcepts(snap)
	if v := findCode(t, got, report.CodeConceptDuplicateKey); v == nil {
		t.Fatalf("expected %s, got %+v", report.CodeConceptDuplicateKey, got)
	} else if v.Severity != report.SeverityError {
		t.Errorf("severity=%s", v.Severity)
	}
}
