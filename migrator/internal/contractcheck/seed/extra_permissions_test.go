package seed

import (
	"context"
	"testing"
)

// TestFixtureLoader_ExtraPermissions covers task 7.3: load the
// `extra_permissions.json` fixture and assert it surfaces the
// orphan permissions used by zombie-detection tests downstream.
func TestFixtureLoader_ExtraPermissions(t *testing.T) {
	l := NewFixtureLoader("../testdata/seed/extra_permissions.json")
	snap, err := l.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	wantOrphans := map[string]bool{
		"users:delete_all":   false,
		"audit:export_v1":    false,
		"deprecated:cleanup": false,
	}
	for _, p := range snap.Permissions {
		if _, ok := wantOrphans[p.Code]; ok {
			wantOrphans[p.Code] = true
		}
	}
	for code, present := range wantOrphans {
		if !present {
			t.Errorf("orphan permission %q missing from fixture", code)
		}
	}
	// Each orphan must have NO entry in role_permissions, otherwise
	// it would not be a zombie candidate.
	rpIndex := make(map[string]struct{}, len(snap.RolePermissions))
	for _, rp := range snap.RolePermissions {
		rpIndex[rp.PermissionCode] = struct{}{}
	}
	for code := range wantOrphans {
		if _, used := rpIndex[code]; used {
			t.Errorf("permission %q should be unused (zombie), found in role_permissions", code)
		}
	}
}
