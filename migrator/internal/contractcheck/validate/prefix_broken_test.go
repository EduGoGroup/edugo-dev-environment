package validate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// TestValidate_PrefixBrokenFixture exercises task 7.2: the
// `prefix-broken` fixture reproduces bug F2·H3.a (apiPrefix=academic:
// declared for resource=announcements, which actually lives behind
// the platform: service). Validate must emit a
// service_prefix_mismatch drift with severity error.
func TestValidate_PrefixBrokenFixture(t *testing.T) {
	root, err := filepath.Abs("../testdata/kmp/prefix-broken")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	kmpSnap, _, err := kmp.Extract([]string{root})
	if err != nil {
		t.Fatalf("kmp.Extract: %v", err)
	}
	if len(kmpSnap.Contracts) == 0 {
		t.Fatal("expected ContractDecl from prefix-broken fixture")
	}

	// Minimal seed where `announcements` resource exists but with a
	// different routing prefix than the contract declares.
	seedSnap := seed.Snapshot{
		Resources: []seed.Resource{
			{Key: "announcements", Name: "Anuncios"},
		},
	}

	drifts := Validate(kmpSnap, seedSnap)

	var hit *Drift
	for i := range drifts {
		if drifts[i].Category == CategoryServicePrefixMismatch && strings.Contains(drifts[i].Identifier, "announcements") {
			hit = &drifts[i]
			break
		}
	}
	if hit == nil {
		t.Fatalf("expected service_prefix_mismatch drift for announcements; got %+v", drifts)
	}
	if hit.Severity != SeverityError {
		t.Errorf("expected severity error, got %s", hit.Severity)
	}
}
