package report

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// updateGolden, when true, rewrites the golden files instead of
// asserting against them. Toggle with `-update-golden=true`.
var updateGolden = flag.Bool("update-golden", false, "regenerate golden files in testdata/golden")

// fixedClock returns a deterministic timestamp used by Build during
// the golden tests. Goldens encode this exact value so reruns are
// byte-stable across machines.
var fixedClock = func() time.Time {
	return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
}

// goldenCase describes one snapshot/golden pair. Violations are
// hand-built so the test is independent of the validators package
// (per the Bloque 4 contract).
type goldenCase struct {
	name       string
	snapshot   *loader.SeedSnapshot
	violations []Violation
}

func cases() []goldenCase {
	return []goldenCase{
		{
			name:       "seed_minimal",
			snapshot:   &loader.SeedSnapshot{}, // zero stats, no violations
			violations: nil,
		},
		{
			name: "seed_broken",
			// Stats picked to be obviously synthetic but stable.
			// We do not need real entities: Build only reads len()
			// of each slice. Empty slices of the right length
			// would suffice but nil leaves them as 0; we want
			// non-zero stats to exercise the table.
			snapshot: &loader.SeedSnapshot{
				Resources:          make([]entities.Resource, 4),
				Permissions:        make([]entities.Permission, 6),
				Roles:              make([]entities.Role, 2),
				ResourceScreens:    make([]entities.ResourceScreen, 3),
				ScreenInstances:    make([]entities.ScreenInstance, 3),
				ConceptTypes:       make([]entities.ConceptType, 1),
				ConceptDefinitions: make([]entities.ConceptDefinition, 2),
			},
			violations: []Violation{
				// Intentionally unsorted to exercise Build's
				// deterministic ordering.
				{
					Severity: SeverityWarning,
					Code:     CodeResourceOrphan,
					Message:  "Resource sin permissions ni resource_screens",
					Entity:   "Resource",
					EntityID: "resource:zeta",
					References: map[string]string{
						"resource_key": "zeta",
					},
				},
				{
					Severity: SeverityError,
					Code:     CodePermResourceMissing,
					Message:  "Permission referencia un Resource inexistente",
					Entity:   "Permission",
					EntityID: "perm:alpha.read",
					References: map[string]string{
						"permission_name": "alpha.read",
						"resource_id":     "00000000-0000-0000-0000-0000000000aa",
					},
				},
				{
					Severity: SeverityError,
					Code:     CodeSlotRefMissing,
					Message:  "slot_data referencia un permission inexistente",
					Entity:   "ScreenInstance",
					EntityID: "screen:dashboard",
					References: map[string]string{
						"permission_name": "missing.permission",
					},
					Path: "$.actions[2].permission",
				},
				{
					Severity: SeverityError,
					Code:     CodePermResourceMissing,
					Message:  "Permission referencia un Resource inexistente",
					Entity:   "Permission",
					EntityID: "perm:beta.write",
					References: map[string]string{
						"permission_name": "beta.write",
					},
				},
				{
					Severity: SeverityInfo,
					Code:     "DIAG_NOTE",
					Message:  "Nota informativa de ejemplo",
					Entity:   "Diagnostic",
					EntityID: "diag:1",
				},
			},
		},
	}
}

func TestGolden(t *testing.T) {
	restore := SetClock(fixedClock)
	t.Cleanup(func() { SetClock(restore) })

	for _, tc := range cases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			report := Build(tc.snapshot, tc.violations, "production")

			tmp := t.TempDir()
			jsonPath, err := WriteJSON(report, tmp)
			if err != nil {
				t.Fatalf("WriteJSON: %v", err)
			}
			mdPath, err := WriteMarkdown(report, tmp)
			if err != nil {
				t.Fatalf("WriteMarkdown: %v", err)
			}

			gotJSON, err := os.ReadFile(jsonPath)
			if err != nil {
				t.Fatalf("read json: %v", err)
			}
			gotMD, err := os.ReadFile(mdPath)
			if err != nil {
				t.Fatalf("read md: %v", err)
			}

			goldenDir := filepath.Join("testdata", "golden")
			goldenJSON := filepath.Join(goldenDir, tc.name+".json")
			goldenMD := filepath.Join(goldenDir, tc.name+".md")

			if *updateGolden {
				if err := os.MkdirAll(goldenDir, 0o755); err != nil {
					t.Fatalf("mkdir golden: %v", err)
				}
				if err := os.WriteFile(goldenJSON, gotJSON, 0o644); err != nil {
					t.Fatalf("write golden json: %v", err)
				}
				if err := os.WriteFile(goldenMD, gotMD, 0o644); err != nil {
					t.Fatalf("write golden md: %v", err)
				}
				t.Logf("updated goldens for %q", tc.name)
				return
			}

			wantJSON, err := os.ReadFile(goldenJSON)
			if err != nil {
				t.Fatalf("read golden json (%s): %v — run with -update-golden=true", goldenJSON, err)
			}
			wantMD, err := os.ReadFile(goldenMD)
			if err != nil {
				t.Fatalf("read golden md (%s): %v — run with -update-golden=true", goldenMD, err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("JSON mismatch for %q.\n--- got ---\n%s\n--- want ---\n%s", tc.name, gotJSON, wantJSON)
			}
			if string(gotMD) != string(wantMD) {
				t.Errorf("Markdown mismatch for %q.\n--- got ---\n%s\n--- want ---\n%s", tc.name, gotMD, wantMD)
			}
		})
	}
}

// TestBuild_Reproducibility asserts A-REQ-10.2: feeding the same set
// of violations in different orders must yield identical Violations
// after Build.
func TestBuild_Reproducibility(t *testing.T) {
	restore := SetClock(fixedClock)
	t.Cleanup(func() { SetClock(restore) })

	tc := cases()[1] // seed_broken has multiple violations
	a := Build(tc.snapshot, tc.violations, "production")

	// Reverse the input slice and rebuild.
	rev := make([]Violation, len(tc.violations))
	for i, v := range tc.violations {
		rev[len(tc.violations)-1-i] = v
	}
	b := Build(tc.snapshot, rev, "production")

	if len(a.Violations) != len(b.Violations) {
		t.Fatalf("len mismatch: %d vs %d", len(a.Violations), len(b.Violations))
	}
	for i := range a.Violations {
		if a.Violations[i].Code != b.Violations[i].Code ||
			a.Violations[i].EntityID != b.Violations[i].EntityID ||
			a.Violations[i].Path != b.Violations[i].Path ||
			a.Violations[i].Severity != b.Violations[i].Severity {
			t.Errorf("ordering mismatch at %d: %+v vs %+v", i, a.Violations[i], b.Violations[i])
		}
	}
}
