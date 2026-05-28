package validators

import (
	"encoding/json"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// uuidFromInt builds a deterministic UUID from a small integer for
// readable test fixtures.
func uuidFromInt(n int) uuid.UUID {
	var b [16]byte
	b[15] = byte(n)
	b[14] = byte(n >> 8)
	return uuid.UUID(b)
}

// minimalSnapshot returns a self-consistent fixture: one menu-visible
// resource, one permission, one role granted that permission, one
// resource_screen marked as default and a screen instance whose
// slot_data references the permission and resource. No validator
// should emit a violation for this snapshot.
func minimalSnapshot() *loader.SeedSnapshot {
	resID := uuidFromInt(1)
	permID := uuidFromInt(2)
	roleID := uuidFromInt(3)
	screenID := uuidFromInt(5)
	rsID := uuidFromInt(6)
	conceptTypeID := uuidFromInt(7)
	conceptDefID := uuidFromInt(8)

	resources := []entities.Resource{{
		ID:            resID,
		Key:           "alpha",
		DisplayName:   "Alpha",
		IsMenuVisible: true,
		IsActive:      true,
	}}
	permissions := []entities.Permission{{
		ID:          permID,
		Name:        "alpha:read",
		DisplayName: "Read alpha",
		ResourceID:  resID,
		Action:      "read",
		IsActive:    true,
	}}
	roles := []entities.Role{{
		ID:          roleID,
		Name:        "viewer",
		DisplayName: "Viewer",
		IsActive:    true,
	}}
	slot := json.RawMessage(`{"permission":"alpha:read","resource":"alpha","actions":[{"permission":"alpha:read"}]}`)
	screens := []entities.ScreenInstance{{
		ID:        screenID,
		ScreenKey: "alpha-list",
		Name:      "Alpha list",
		SlotData:  slot,
		IsActive:  true,
	}}
	rs := []entities.ResourceScreen{{
		ID:          rsID,
		ResourceID:  resID,
		ResourceKey: "alpha",
		ScreenKey:   "alpha-list",
		ScreenType:  "list",
		IsDefault:   true,
		IsActive:    true,
	}}
	cts := []entities.ConceptType{{
		ID:   conceptTypeID,
		Name: "Primary",
		Code: "primary_school",
	}}
	cdefs := []entities.ConceptDefinition{{
		ID:            conceptDefID,
		ConceptTypeID: conceptTypeID,
		TermKey:       "org.name_singular",
		TermValue:     "Escuela",
	}}

	return loader.NewSnapshot(resources, permissions, roles, rs, screens, nil, cts, cdefs)
}

// findCode returns the first violation whose Code matches; helps tests
// assert presence without depending on slice order.
func findCode(t *testing.T, vs []report.Violation, code string) *report.Violation {
	t.Helper()
	for i := range vs {
		if vs[i].Code == code {
			return &vs[i]
		}
	}
	return nil
}

// countCode reports how many violations carry the given Code.
func countCode(vs []report.Violation, code string) int {
	n := 0
	for _, v := range vs {
		if v.Code == code {
			n++
		}
	}
	return n
}

// assertEmpty fails the test if the slice contains any violation.
func assertEmpty(t *testing.T, vs []report.Violation) {
	t.Helper()
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(vs), vs)
	}
}
