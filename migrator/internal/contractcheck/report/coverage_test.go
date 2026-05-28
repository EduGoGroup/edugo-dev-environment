package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

// TestWriteMarkdown_RoundTrip exercises WriteMarkdown end-to-end:
// it lands a Markdown file in the requested directory and the file
// contains the required headings.
func TestWriteMarkdown_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	r := newSyntheticResult()

	path, err := WriteMarkdown(r, dir)
	if err != nil {
		t.Fatalf("WriteMarkdown error: %v", err)
	}
	if !strings.HasSuffix(path, ".md") {
		t.Fatalf("expected .md path, got %s", path)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written file: %v", err)
	}
	for _, want := range []string{
		"# Contract Check Report",
		"## Resumen ejecutivo",
		"### Conteo por categoría",
		"## Drifts por categoría",
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("Markdown missing heading %q", want)
		}
	}
}

// TestWriteMarkdown_FileNamingMatchesJSON ensures both reports share
// the same RFC3339-compact timestamp so a single run produces a
// matched pair.
func TestWriteMarkdown_FileNamingMatchesJSON(t *testing.T) {
	dir := t.TempDir()
	r := newSyntheticResult()

	jsonPath, err := WriteJSON(r, dir)
	if err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	mdPath, err := WriteMarkdown(r, dir)
	if err != nil {
		t.Fatalf("WriteMarkdown: %v", err)
	}
	jsonBase := strings.TrimSuffix(filepath.Base(jsonPath), ".json")
	mdBase := strings.TrimSuffix(filepath.Base(mdPath), ".md")
	if jsonBase != mdBase {
		t.Fatalf("expected matching basenames, got json=%s md=%s", jsonBase, mdBase)
	}
}

// TestSortDrifts_StableOrder covers the stable-sort branch of
// sortDrifts (Category asc, Severity desc, Identifier asc).
func TestSortDrifts_StableOrder(t *testing.T) {
	in := []validate.Drift{
		{Category: validate.CategoryRoleUnused, Severity: validate.SeverityWarning, Identifier: "z"},
		{Category: validate.CategoryRoleUnused, Severity: validate.SeverityError, Identifier: "a"},
		{Category: validate.CategoryRoleUnused, Severity: validate.SeverityError, Identifier: "b"},
		{Category: validate.CategoryRolePhantom, Severity: validate.SeverityError, Identifier: "x"},
	}
	sortDrifts(in)
	want := []string{
		"role_phantom/x",
		"role_unused/a",
		"role_unused/b",
		"role_unused/z",
	}
	for i, d := range in {
		got := string(d.Category) + "/" + d.Identifier
		if got != want[i] {
			t.Errorf("sorted[%d] = %s, want %s", i, got, want[i])
		}
	}
}

// TestUpdateBaseline_OmitsBaselineDiff confirms that calling
// UpdateBaseline drops the BaselineDiff section so the next run
// computes a fresh diff against the new snapshot.
func TestUpdateBaseline_OmitsBaselineDiff(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "contract-check-baseline.json")

	r := newSyntheticResult()
	r.BaselineDiff = &BaselineDiff{
		Regressions: []validate.Drift{{Category: validate.CategoryRolePhantom, Identifier: "stale"}},
	}
	if err := UpdateBaseline(r, path); err != nil {
		t.Fatalf("UpdateBaseline: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read baseline: %v", err)
	}
	var roundtrip map[string]any
	if err := json.Unmarshal(raw, &roundtrip); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, present := roundtrip["baseline_diff"]; present {
		t.Fatalf("baseline_diff should be omitted; got %s", string(raw))
	}
}

func newSyntheticResult() *Result {
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	drifts := []validate.Drift{
		{
			Direction:  validate.DirectionFEOnly,
			Category:   validate.CategoryScreenKeyPhantom,
			Severity:   validate.SeverityError,
			Identifier: "ghost-screen",
			Detail:     "Test fixture",
		},
	}
	return NewResult(ts, kmp.Snapshot{}, seed.Snapshot{}, drifts)
}
