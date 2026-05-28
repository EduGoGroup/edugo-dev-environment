package loader

import (
	"encoding/json"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

func TestLoad_ProductionSeedHydratesAndIndexes(t *testing.T) {
	snap, err := Load(RunOptions{SeedSource: SeedSourceProduction})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(snap.Resources) == 0 {
		t.Fatal("expected at least one resource")
	}
	if len(snap.Permissions) == 0 {
		t.Fatal("expected at least one permission")
	}
	if len(snap.Roles) == 0 {
		t.Fatal("expected at least one role")
	}
	if len(snap.ConceptTypes) == 0 {
		t.Fatal("expected at least one concept_type")
	}

	if got, want := len(snap.ResourceByID), len(snap.Resources); got != want {
		t.Fatalf("ResourceByID size=%d want %d", got, want)
	}
	if got, want := len(snap.ResourceByKey), len(snap.Resources); got != want {
		t.Fatalf("ResourceByKey size=%d want %d", got, want)
	}
	if got, want := len(snap.PermissionByID), len(snap.Permissions); got != want {
		t.Fatalf("PermissionByID size=%d want %d", got, want)
	}
	if got, want := len(snap.PermissionByName), len(snap.Permissions); got != want {
		t.Fatalf("PermissionByName size=%d want %d", got, want)
	}
	if got, want := len(snap.RoleByID), len(snap.Roles); got != want {
		t.Fatalf("RoleByID size=%d want %d", got, want)
	}
	if got, want := len(snap.ConceptTypeByID), len(snap.ConceptTypes); got != want {
		t.Fatalf("ConceptTypeByID size=%d want %d", got, want)
	}

	first := snap.Resources[0]
	got, ok := snap.ResourceByID[first.ID]
	if !ok {
		t.Fatalf("ResourceByID missing first.ID=%s", first.ID)
	}
	if got.Key != first.Key {
		t.Fatalf("ResourceByID alias mismatch: got %q want %q", got.Key, first.Key)
	}

	for _, p := range snap.Permissions {
		if _, ok := snap.ResourceByID[p.ResourceID]; !ok {
			// Not asserting integrity here (that's the validator's job),
			// just sanity-checking the maps have the production IDs.
			t.Logf("permission %s references resource %s missing from index (will be flagged by validators)", p.Name, p.ResourceID)
		}
	}
}

func TestLoad_RejectsUnknownSource(t *testing.T) {
	if _, err := Load(RunOptions{SeedSource: "bogus"}); err == nil {
		t.Fatal("expected error for unknown seed source")
	}
}

func TestLoad_DefaultsToProductionWhenSourceEmpty(t *testing.T) {
	if _, err := Load(RunOptions{}); err != nil {
		t.Fatalf("Load with empty SeedSource: %v", err)
	}
}

func TestNewSnapshot_BuildsIndexesFromSyntheticFixture(t *testing.T) {
	resA := entities.Resource{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Key: "alpha"}
	resB := entities.Resource{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Key: "beta"}
	permA := entities.Permission{ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), Name: "alpha:read", ResourceID: resA.ID, Action: "read"}
	role := entities.Role{ID: uuid.MustParse("44444444-4444-4444-4444-444444444444"), Name: "viewer"}
	screen := entities.ScreenInstance{ID: uuid.MustParse("66666666-6666-6666-6666-666666666666"), ScreenKey: "alpha-list", SlotData: json.RawMessage(`{}`)}
	rs := entities.ResourceScreen{ID: uuid.MustParse("77777777-7777-7777-7777-777777777777"), ResourceID: resA.ID, ResourceKey: "alpha", ScreenKey: "alpha-list", ScreenType: "list", IsDefault: true}
	ctype := entities.ConceptType{ID: uuid.MustParse("88888888-8888-8888-8888-888888888888"), Code: "primary_school"}
	cdef := entities.ConceptDefinition{ID: uuid.MustParse("99999999-9999-9999-9999-999999999999"), ConceptTypeID: ctype.ID, TermKey: "org.name_singular", TermValue: "Escuela"}

	snap := NewSnapshot(
		[]entities.Resource{resA, resB},
		[]entities.Permission{permA},
		[]entities.Role{role},
		[]entities.ResourceScreen{rs},
		[]entities.ScreenInstance{screen},
		nil,
		[]entities.ConceptType{ctype},
		[]entities.ConceptDefinition{cdef},
	)

	if got := snap.ResourceByID[resA.ID]; got == nil || got.Key != "alpha" {
		t.Fatalf("ResourceByID[alpha] = %+v", got)
	}
	if got := snap.ResourceByKey["beta"]; got == nil || got.ID != resB.ID {
		t.Fatalf("ResourceByKey[beta] = %+v", got)
	}
	if got := snap.PermissionByID[permA.ID]; got == nil || got.Name != "alpha:read" {
		t.Fatalf("PermissionByID = %+v", got)
	}
	if got := snap.PermissionByName["alpha:read"]; got == nil || got.ID != permA.ID {
		t.Fatalf("PermissionByName = %+v", got)
	}
	if got := snap.RoleByID[role.ID]; got == nil || got.Name != "viewer" {
		t.Fatalf("RoleByID = %+v", got)
	}
	if got := snap.ScreenByKey["alpha-list"]; got == nil || got.ID != screen.ID {
		t.Fatalf("ScreenByKey = %+v", got)
	}
	if got := snap.ConceptTypeByID[ctype.ID]; got == nil || got.Code != "primary_school" {
		t.Fatalf("ConceptTypeByID = %+v", got)
	}
}
