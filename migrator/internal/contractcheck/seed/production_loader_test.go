package seed

import (
	"context"
	"testing"
)

// TestProductionLoader_Load_Smoke ensures the adapter projects the
// production seed without losing rows. The expected counts mirror the
// values reported by Phase A's baseline (2026-05-08); refresh them if
// the seed is intentionally changed.
func TestProductionLoader_Load_Smoke(t *testing.T) {
	loader := NewProductionLoader("")
	snap, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	wantMin := map[string]int{
		"resources":        33,
		"permissions":      104,
		"roles":            12,
		"resource_screens": 62,
		"screen_instances": 47,
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
