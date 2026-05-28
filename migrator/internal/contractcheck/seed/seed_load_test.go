package seed

import (
	"context"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestFixtureLoader_HappySnapshot(t *testing.T) {
	loader := NewFixtureLoader(filepath.Join("..", "testdata", "seed", "happy_snapshot.json"))
	snap, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if len(snap.Resources) == 0 {
		t.Fatal("expected resources to be populated")
	}
	if len(snap.Permissions) == 0 {
		t.Fatal("expected permissions to be populated")
	}
	if len(snap.Roles) == 0 {
		t.Fatal("expected roles to be populated")
	}
	if len(snap.ResourceScreens) == 0 {
		t.Fatal("expected resource_screens to be populated")
	}

	// Encontrar dashboard-teacher debe traer is_default=true.
	var found bool
	for _, rs := range snap.ResourceScreens {
		if rs.ScreenKey == "dashboard-teacher" {
			if !rs.IsDefault || rs.ScreenType != "dashboard" {
				t.Errorf("dashboard-teacher metadata mismatch: %+v", rs)
			}
			found = true
		}
	}
	if !found {
		t.Fatal("dashboard-teacher missing from happy_snapshot")
	}

	// SlotData debe ser json.RawMessage (no normalizado).
	if len(snap.ScreenInstances) == 0 || len(snap.ScreenInstances[0].SlotData) == 0 {
		t.Fatal("expected ScreenInstances[0].SlotData to carry raw JSON")
	}
}

func TestFixtureLoader_MissingFile(t *testing.T) {
	loader := NewFixtureLoader(filepath.Join("..", "testdata", "seed", "this-does-not-exist.json"))
	if _, err := loader.Load(context.Background()); err == nil {
		t.Fatal("expected error for missing fixture, got nil")
	}
}

func TestFixtureLoader_EmptyPath(t *testing.T) {
	loader := NewFixtureLoader("")
	if _, err := loader.Load(context.Background()); err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestPermissionsReferencedInSlots(t *testing.T) {
	loader := NewFixtureLoader(filepath.Join("..", "testdata", "seed", "happy_snapshot.json"))
	snap, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	got := snap.PermissionsReferencedInSlots()
	want := []string{
		"menu:read",
		"schools:delete",
		"schools:read",
		"schools:update",
	}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PermissionsReferencedInSlots\n got: %v\nwant: %v", got, want)
	}
}

func TestPermissionsReferencedInSlots_HandlesEmpty(t *testing.T) {
	snap := Snapshot{ScreenInstances: nil}
	if got := snap.PermissionsReferencedInSlots(); len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

// La interfaz Loader está satisfecha por FixtureLoader; este test
// asegura que cualquier refactor mantenga el contrato.
func TestFixtureLoader_ImplementsLoader(t *testing.T) {
	var _ Loader = (*FixtureLoader)(nil)
}
