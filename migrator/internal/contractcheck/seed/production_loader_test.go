package seed

import (
	"context"
	"testing"
)

// TestProductionLoader_Load_Smoke ensures the adapter projects the
// production seed without losing rows. The expected counts mirror the
// baseline refreshed 2026-06-15 (3.58.0): after the SDUI admin-tool
// screens were pruned and the L3 accessor mirror was fixed to include
// materials-list, resource_screens settled at 53 and screen_instances at 48.
// Refresh these counts whenever the seed is intentionally changed.
func TestProductionLoader_Load_Smoke(t *testing.T) {
	loader := NewProductionLoader("")
	snap, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	wantMin := map[string]int{
		"resources":        33,
		"permissions":      98,
		"roles":            12,
		"resource_screens": 53,
		"screen_instances": 48,
	}
	got := map[string]int{
		"resources":        len(snap.Resources),
		"permissions":      len(snap.Permissions),
		"roles":            len(snap.Roles),
		"resource_screens": len(snap.ResourceScreens),
		"screen_instances": len(snap.ScreenInstances),
	}
	for k, min := range wantMin {
		if got[k] < min {
			t.Errorf("snap.%s = %d, expected at least %d", k, got[k], min)
		}
	}
}
