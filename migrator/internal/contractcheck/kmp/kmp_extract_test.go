package kmp

import (
	"path/filepath"
	"testing"
)

func TestExtract_Happy(t *testing.T) {
	root := filepath.Join("..", "testdata", "kmp", "happy")
	snap, errs, err := Extract([]string{root})
	if err != nil {
		t.Fatalf("Extract returned error: %v (errs=%v)", err, errs)
	}
	if len(errs) != 0 {
		t.Fatalf("expected no extract errors, got: %v", errs)
	}

	wantScreenKeys := []string{"schools-form"}
	for _, k := range wantScreenKeys {
		if _, ok := snap.ScreenKeys[k]; !ok {
			t.Errorf("missing screenKey %q in snapshot; got=%v", k, keys(snap.ScreenKeys))
		}
	}

	// La constante en el comentario NO debe estar presente.
	if _, ok := snap.ScreenKeys["phantom-from-comment"]; ok {
		t.Errorf("comment literal phantom-from-comment leaked into ScreenKeys")
	}
	if _, ok := snap.ScreenKeys["phantom-from-block"]; ok {
		t.Errorf("block-comment literal phantom-from-block leaked into ScreenKeys")
	}

	wantPerms := []string{"schools:create", "schools:update"}
	for _, p := range wantPerms {
		if _, ok := snap.Permissions[p]; !ok {
			t.Errorf("missing permission %q; got=%v", p, keys(snap.Permissions))
		}
	}

	wantRoles := []string{"teacher", "student", "guardian", "home"}
	for _, r := range wantRoles {
		if _, ok := snap.Roles[r]; !ok {
			t.Errorf("missing role %q; got=%v", r, keys(snap.Roles))
		}
	}

	if len(snap.Contracts) != 1 {
		t.Fatalf("expected 1 ContractDecl, got %d (%v)", len(snap.Contracts), snap.Contracts)
	}
	c := snap.Contracts[0]
	if c.ScreenKey != "schools-form" || c.APIPrefix != "academic:" || c.BasePath != "/api/v1/schools" || c.Resource != "schools" {
		t.Errorf("ContractDecl mismatch: %+v", c)
	}
}

func TestExtract_PhantomScreen(t *testing.T) {
	root := filepath.Join("..", "testdata", "kmp", "phantom-screen")
	snap, errs, err := Extract([]string{root})
	if err != nil {
		t.Fatalf("Extract returned error: %v (errs=%v)", err, errs)
	}
	if len(errs) != 0 {
		t.Fatalf("expected no extract errors, got: %v", errs)
	}
	if _, ok := snap.ScreenKeys["ghost-form"]; !ok {
		t.Errorf("expected ghost-form in ScreenKeys, got=%v", keys(snap.ScreenKeys))
	}
	if len(snap.Contracts) != 1 {
		t.Fatalf("expected 1 ContractDecl, got %d", len(snap.Contracts))
	}
}

func TestExtract_MultiModule(t *testing.T) {
	rootA := filepath.Join("..", "testdata", "kmp", "multi-module", "moduleA")
	rootB := filepath.Join("..", "testdata", "kmp", "multi-module", "moduleB")
	snap, errs, err := Extract([]string{rootA, rootB})
	if err != nil {
		t.Fatalf("Extract returned error: %v (errs=%v)", err, errs)
	}
	if len(errs) != 0 {
		t.Fatalf("expected no extract errors, got: %v", errs)
	}
	for _, want := range []string{"alpha-list", "beta-list"} {
		if _, ok := snap.ScreenKeys[want]; !ok {
			t.Errorf("missing %q in ScreenKeys; got=%v", want, keys(snap.ScreenKeys))
		}
	}
	// El archivo IgnoreMeTest.kt no debe contribuir nada.
	if _, ok := snap.ScreenKeys["should-not-be-extracted"]; ok {
		t.Errorf("Test.kt suffix file leaked screenKey")
	}
	if len(snap.Contracts) != 2 {
		t.Fatalf("expected 2 ContractDecls (one per module), got %d (%v)", len(snap.Contracts), snap.Contracts)
	}
}

func TestExtract_MissingRoot(t *testing.T) {
	// Path inexistente: devuelve ExtractError pero no aborta si hay otro
	// root válido (B-REQ-9.3).
	good := filepath.Join("..", "testdata", "kmp", "happy")
	bad := filepath.Join("..", "testdata", "kmp", "this-does-not-exist")
	snap, errs, err := Extract([]string{good, bad})
	if err != nil {
		t.Fatalf("Extract should not abort with one good root: %v", err)
	}
	if len(errs) == 0 {
		t.Fatalf("expected at least one ExtractError for missing root")
	}
	if len(snap.ScreenKeys) == 0 {
		t.Fatalf("expected good root to still produce screen keys")
	}

	// Todos los roots inexistentes -> error sistémico.
	_, _, err = Extract([]string{bad})
	if err == nil {
		t.Fatalf("expected systemic error when no root exists")
	}
}

func TestStripComments(t *testing.T) {
	src := `// line comment
val a = "keep this"
/* block
   comment */
val b = "https://example.com" // trailing
`
	cleaned, lineMap := stripComments(src)
	if got := cleaned; len(got) == 0 {
		t.Fatal("cleaned string is empty")
	}
	// El URL adentro del string no debe interpretarse como comentario.
	if !contains(cleaned, `"https://example.com"`) {
		t.Errorf("URL inside string was eaten by comment stripper: %q", cleaned)
	}
	// El comentario "// line comment" no debe sobrevivir.
	if contains(cleaned, "line comment") {
		t.Errorf("// comment leaked: %q", cleaned)
	}
	// El bloque /* ... */ no debe sobrevivir.
	if contains(cleaned, "block") {
		t.Errorf("/* */ comment leaked: %q", cleaned)
	}
	// lineMap debe coincidir con número de líneas del cleaned.
	if len(lineMap) == 0 {
		t.Fatal("lineMap is empty")
	}
}

// helpers

func keys(m map[string][]Location) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
