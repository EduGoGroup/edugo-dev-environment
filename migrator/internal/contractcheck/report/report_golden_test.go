package report

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

// updateGolden, cuando es true, regenera los archivos golden en lugar
// de compararlos. Se controla con `go test -update-golden`.
var updateGolden = flag.Bool("update-golden", false, "regenerate golden files in testdata/golden")

// fixedTimestamp es el instante "petrificado" usado en goldens para que
// dos ejecuciones del test consecutivas comparen byte-a-byte.
var fixedTimestamp = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestGolden_Happy(t *testing.T) {
	r := newHappyResult()
	assertGolden(t, "happy", r)
}

func TestGolden_WithDrifts(t *testing.T) {
	r := newWithDriftsResult()
	assertGolden(t, "with_drifts", r)
}

// TestWriteJSON_Deterministic: dos invocaciones consecutivas con el
// mismo Result producen archivos byte-idénticos. Cubre B-REQ-10.
func TestWriteJSON_Deterministic(t *testing.T) {
	r := newWithDriftsResult()
	dir := t.TempDir()
	pathA, err := WriteJSON(r, dir)
	if err != nil {
		t.Fatalf("WriteJSON A: %v", err)
	}
	a, err := os.ReadFile(pathA)
	if err != nil {
		t.Fatalf("read A: %v", err)
	}

	// Mover el archivo para no chocar.
	if err := os.Remove(pathA); err != nil {
		t.Fatalf("remove A: %v", err)
	}
	pathB, err := WriteJSON(r, dir)
	if err != nil {
		t.Fatalf("WriteJSON B: %v", err)
	}
	b, err := os.ReadFile(pathB)
	if err != nil {
		t.Fatalf("read B: %v", err)
	}
	if string(a) != string(b) {
		t.Fatal("WriteJSON output is not byte-identical across invocations")
	}
}

func TestUpdateBaseline_OmitsDiff(t *testing.T) {
	r := newWithDriftsResult()
	r.BaselineDiff = &BaselineDiff{
		Regressions: []validate.Drift{{Category: "x", Identifier: "y", Severity: validate.SeverityError, Detail: "regression"}},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "contract-check-baseline.json")
	if err := UpdateBaseline(r, path); err != nil {
		t.Fatalf("UpdateBaseline: %v", err)
	}
	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if loaded.BaselineDiff != nil {
		t.Fatalf("expected baseline_diff omitted, got %+v", loaded.BaselineDiff)
	}
	if len(loaded.Drifts) != len(r.Drifts) {
		t.Fatalf("expected %d drifts in baseline, got %d", len(r.Drifts), len(loaded.Drifts))
	}
}

// assertGolden compara la salida JSON + Markdown del Result contra los
// archivos en testdata/golden/<name>.{json,md}. Con -update-golden,
// los regenera.
func assertGolden(t *testing.T, name string, r *Result) {
	t.Helper()
	jsonBytes, err := marshalResult(r)
	if err != nil {
		t.Fatalf("marshalResult: %v", err)
	}
	mdBytes, err := renderMarkdown(r)
	if err != nil {
		t.Fatalf("renderMarkdown: %v", err)
	}

	jsonPath := filepath.Join("testdata", "golden", name+".json")
	mdPath := filepath.Join("testdata", "golden", name+".md")

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
			t.Fatalf("write golden json: %v", err)
		}
		if err := os.WriteFile(mdPath, mdBytes, 0o644); err != nil {
			t.Fatalf("write golden md: %v", err)
		}
		return
	}

	wantJSON, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("read golden json (%s): %v — run with -update-golden to create it", jsonPath, err)
	}
	if string(wantJSON) != string(jsonBytes) {
		t.Errorf("JSON golden mismatch for %s\n--- want ---\n%s\n--- got ---\n%s", name, wantJSON, jsonBytes)
	}

	wantMD, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("read golden md (%s): %v — run with -update-golden to create it", mdPath, err)
	}
	if string(wantMD) != string(mdBytes) {
		t.Errorf("Markdown golden mismatch for %s\n--- want ---\n%s\n--- got ---\n%s", name, wantMD, mdBytes)
	}
}

// newHappyResult produce un Result sin drifts: solo Stats + Summary
// vacío. Sirve para verificar el layout "todo verde".
func newHappyResult() *Result {
	k := kmp.Snapshot{
		ScreenKeys:  map[string][]kmp.Location{"dashboard-teacher": {{FilePath: "kmp/Foo.kt", Line: 10}}},
		Permissions: map[string][]kmp.Location{},
		Roles:       map[string][]kmp.Location{},
		Contracts:   nil,
	}
	s := seed.Snapshot{
		Resources:       []seed.Resource{{Key: "schools"}},
		Permissions:     []seed.Permission{{Code: "schools:read"}},
		Roles:           []seed.Role{{Code: "teacher"}},
		ResourceScreens: []seed.ResourceScreen{{ScreenKey: "dashboard-teacher", ScreenType: "dashboard", IsDefault: true}},
	}
	return NewResult(fixedTimestamp, k, s, nil)
}

// newWithDriftsResult produce un Result con un drift por categoría
// (uno error, uno warning, uno info) cubriendo las 7 categorías
// canónicas para que el golden ejercite todas las ramas del template.
func newWithDriftsResult() *Result {
	loc := []kmp.Location{{FilePath: "kmp/Foo.kt", Line: 42, Snippet: "override val screenKey = \"ghost-form\""}}
	drifts := []validate.Drift{
		{
			Direction:  validate.DirectionFEOnly,
			Category:   validate.CategoryScreenKeyPhantom,
			Severity:   validate.SeverityError,
			Identifier: "ghost-form",
			Detail:     "FE declara screenKey 'ghost-form' que el seed no contiene.",
			Evidence:   loc,
		},
		{
			Direction:  validate.DirectionBEOnly,
			Category:   validate.CategoryScreenKeyDead,
			Severity:   validate.SeverityWarning,
			Identifier: "legacy-list",
			Detail:     "Seed declara 'legacy-list' pero KMP no atiende esa pantalla.",
		},
		{
			Direction:  validate.DirectionFEOnly,
			Category:   validate.CategoryPermissionPhantom,
			Severity:   validate.SeverityError,
			Identifier: "ghost:read",
			Detail:     "Permiso inferido por FE no existe en iam.permissions.",
			Evidence:   loc,
		},
		{
			Direction:  validate.DirectionBEOnly,
			Category:   validate.CategoryPermissionZombie,
			Severity:   validate.SeverityInfo,
			Identifier: "audit:purge",
			Detail:     "Permiso seedado sin role_permission ni referencia en slot_data.",
		},
		{
			Direction:  validate.DirectionFEOnly,
			Category:   validate.CategoryRolePhantom,
			Severity:   validate.SeverityError,
			Identifier: "principal",
			Detail:     "FE menciona el rol 'principal' que no existe en iam.roles.",
		},
		{
			Direction:  validate.DirectionBEOnly,
			Category:   validate.CategoryRoleUnused,
			Severity:   validate.SeverityWarning,
			Identifier: "auditor",
			Detail:     "Rol 'auditor' seedado sin uso en KMP.",
		},
		{
			Direction:  validate.DirectionMismatch,
			Category:   validate.CategoryServicePrefixMismatch,
			Severity:   validate.SeverityError,
			Identifier: "announcements",
			Detail:     "apiPrefix=academic: declarado pero el ruteo canónico es platform:.",
			Evidence:   loc,
		},
	}

	k := kmp.Snapshot{
		ScreenKeys: map[string][]kmp.Location{
			"ghost-form":        {{FilePath: "kmp/Foo.kt", Line: 42}},
			"dashboard-teacher": {{FilePath: "kmp/Bar.kt", Line: 1}},
		},
		Permissions: map[string][]kmp.Location{"ghost:read": loc},
		Roles:       map[string][]kmp.Location{"principal": {{FilePath: "kmp/Baz.kt", Line: 5}}},
		Contracts:   []kmp.ContractDecl{{ScreenKey: "announcements", APIPrefix: "academic:", Resource: "announcements"}},
	}
	s := seed.Snapshot{
		Resources:       []seed.Resource{{Key: "schools"}},
		Permissions:     []seed.Permission{{Code: "schools:read"}, {Code: "audit:purge"}},
		Roles:           []seed.Role{{Code: "teacher"}, {Code: "auditor", Scope: "system"}},
		ResourceScreens: []seed.ResourceScreen{{ScreenKey: "legacy-list", ScreenType: "list"}},
	}
	return NewResult(fixedTimestamp, k, s, drifts)
}
