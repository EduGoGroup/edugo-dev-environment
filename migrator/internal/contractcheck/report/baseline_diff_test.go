package report

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

// TestComputeDiff_RegressionWhenCurrAddsDrift: caso A del Bloque 5.6.
func TestComputeDiff_RegressionWhenCurrAddsDrift(t *testing.T) {
	prev := newResultWithDrifts()
	curr := newResultWithDrifts(validate.Drift{
		Category:   validate.CategoryScreenKeyPhantom,
		Severity:   validate.SeverityError,
		Identifier: "new-screen",
		Detail:     "drift nuevo",
	})

	diff := ComputeDiff(prev, curr)
	if len(diff.Regressions) != 1 || diff.Regressions[0].Identifier != "new-screen" {
		t.Fatalf("expected one regression for 'new-screen', got %+v", diff.Regressions)
	}
	if len(diff.Fixes) != 0 {
		t.Fatalf("expected no fixes, got %+v", diff.Fixes)
	}
}

// TestComputeDiff_FixWhenPrevHadDrift: caso B.
func TestComputeDiff_FixWhenPrevHadDrift(t *testing.T) {
	prev := newResultWithDrifts(validate.Drift{
		Category:   validate.CategoryRolePhantom,
		Severity:   validate.SeverityError,
		Identifier: "principal",
		Detail:     "drift en prev",
	})
	curr := newResultWithDrifts()

	diff := ComputeDiff(prev, curr)
	if len(diff.Fixes) != 1 || diff.Fixes[0].Identifier != "principal" {
		t.Fatalf("expected one fix for 'principal', got %+v", diff.Fixes)
	}
	if len(diff.Regressions) != 0 {
		t.Fatalf("expected no regressions, got %+v", diff.Regressions)
	}
}

// TestComputeDiff_StableWhenBothHaveDrift: caso C.
func TestComputeDiff_StableWhenBothHaveDrift(t *testing.T) {
	d := validate.Drift{
		Category:   validate.CategoryServicePrefixMismatch,
		Severity:   validate.SeverityError,
		Identifier: "announcements",
		Detail:     "estable",
	}
	prev := newResultWithDrifts(d)
	curr := newResultWithDrifts(d)

	diff := ComputeDiff(prev, curr)
	if len(diff.Regressions) != 0 || len(diff.Fixes) != 0 {
		t.Fatalf("expected empty diff when drift is stable, got %+v", diff)
	}
}

// TestComputeDiff_NilPrev: caso D — sin baseline.
func TestComputeDiff_NilPrev(t *testing.T) {
	curr := newResultWithDrifts(validate.Drift{
		Category:   validate.CategoryScreenKeyPhantom,
		Severity:   validate.SeverityError,
		Identifier: "any",
	})
	diff := ComputeDiff(nil, curr)
	if len(diff.Regressions) != 0 || len(diff.Fixes) != 0 {
		t.Fatalf("expected empty diff when prev is nil, got %+v", diff)
	}
}

// TestComputeDiff_SeverityChangeIsNotRegression: el match es por
// (Category, Identifier), así que un cambio de warning→error sobre el
// MISMO (cat,id) NO se reporta como regresión por el algoritmo. Esto
// es intencional (queda capturado en design §7.3 + tasks 5.3).
func TestComputeDiff_SeverityChangeIsNotRegression(t *testing.T) {
	prev := newResultWithDrifts(validate.Drift{
		Category:   validate.CategoryScreenKeyDead,
		Severity:   validate.SeverityWarning,
		Identifier: "legacy",
	})
	curr := newResultWithDrifts(validate.Drift{
		Category:   validate.CategoryScreenKeyDead,
		Severity:   validate.SeverityError, // escalado
		Identifier: "legacy",
	})
	diff := ComputeDiff(prev, curr)
	if len(diff.Regressions) != 0 {
		t.Fatalf("severity-only change should not produce regressions, got %+v", diff.Regressions)
	}
	if len(diff.Fixes) != 0 {
		t.Fatalf("severity-only change should not produce fixes, got %+v", diff.Fixes)
	}
}

func TestLoadBaseline_MissingReturnsNilNil(t *testing.T) {
	dir := t.TempDir()
	r, err := LoadBaseline(filepath.Join(dir, "does-not-exist.json"))
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if r != nil {
		t.Fatalf("expected nil result for missing file, got %+v", r)
	}
}

func TestLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	want := newResultWithDrifts(validate.Drift{
		Category:   validate.CategoryRolePhantom,
		Severity:   validate.SeverityError,
		Identifier: "x",
	})
	if err := UpdateBaseline(want, path); err != nil {
		t.Fatalf("UpdateBaseline: %v", err)
	}
	got, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if got == nil {
		t.Fatal("LoadBaseline returned nil")
	}
	if len(got.Drifts) != len(want.Drifts) {
		t.Fatalf("drift count mismatch: got %d want %d", len(got.Drifts), len(want.Drifts))
	}
}

// newResultWithDrifts construye un Result mínimo con los drifts dados.
func newResultWithDrifts(extra ...validate.Drift) *Result {
	ts := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	k := kmp.Snapshot{ScreenKeys: map[string][]kmp.Location{}, Permissions: map[string][]kmp.Location{}, Roles: map[string][]kmp.Location{}}
	s := seed.Snapshot{}
	return NewResult(ts, k, s, extra)
}
