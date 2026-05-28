package validators

import (
	"fmt"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/google/uuid"
)

// ValidateConcepts covers A-REQ-5: every ConceptDefinition must point
// to an existing ConceptType (CONCEPT_TYPE_MISSING) and the
// (ConceptTypeID, TermKey) tuple must be unique
// (CONCEPT_DUPLICATE_KEY).
func ValidateConcepts(s *loader.SeedSnapshot) []report.Violation {
	if s == nil {
		return nil
	}

	violations := make([]report.Violation, 0)

	type conceptKey struct {
		typeID  uuid.UUID
		termKey string
	}
	seen := make(map[conceptKey]string, len(s.ConceptDefinitions))

	for i := range s.ConceptDefinitions {
		cd := s.ConceptDefinitions[i]

		if _, ok := s.ConceptTypeByID[cd.ConceptTypeID]; !ok {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeConceptTypeMissing),
				Code:     report.CodeConceptTypeMissing,
				Message:  fmt.Sprintf("La definición %q referencia un concept_type inexistente.", cd.TermKey),
				Entity:   "ConceptDefinition",
				EntityID: cd.ID.String(),
				References: map[string]string{
					"term_key":               cd.TermKey,
					"missing_concept_type_id": cd.ConceptTypeID.String(),
				},
			})
		}

		key := conceptKey{typeID: cd.ConceptTypeID, termKey: cd.TermKey}
		if prev, dup := seen[key]; dup {
			violations = append(violations, report.Violation{
				Severity: report.SeverityFor(report.CodeConceptDuplicateKey),
				Code:     report.CodeConceptDuplicateKey,
				Message:  fmt.Sprintf("Definición de concepto duplicada para term_key=%q dentro del mismo concept_type.", cd.TermKey),
				Entity:   "ConceptDefinition",
				EntityID: cd.ID.String(),
				References: map[string]string{
					"term_key":          cd.TermKey,
					"concept_type_id":   cd.ConceptTypeID.String(),
					"existing_definition": prev,
				},
			})
			continue
		}
		seen[key] = cd.ID.String()
	}

	return violations
}
