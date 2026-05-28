package kmp

import (
	"path/filepath"
	"testing"
	"time"
)

// BenchmarkExtract satisfies B-REQ-10.3 / task 7.5: a full extract
// over the in-tree multi-module fixture must complete in well under
// 5 s. We use multi-module because it's the largest fixture and the
// closest analogue to a real KMP repo with several `commonMain`
// trees.
func BenchmarkExtract(b *testing.B) {
	root, err := filepath.Abs("../testdata/kmp/multi-module")
	if err != nil {
		b.Fatalf("abs: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := Extract([]string{root}); err != nil {
			b.Fatalf("Extract failed: %v", err)
		}
	}
}

// TestExtract_PerformanceBudget guards the 5 s ceiling outside of
// `go test -bench`: a regular `go test` run measures one Extract
// call and fails if it crosses the budget. Keeps the contract
// enforced even when nobody runs benchmarks.
func TestExtract_PerformanceBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("performance check skipped under -short")
	}
	root, err := filepath.Abs("../testdata/kmp/multi-module")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	start := time.Now()
	if _, _, err := Extract([]string{root}); err != nil {
		t.Fatalf("Extract failed: %v", err)
	}
	if d := time.Since(start); d > 5*time.Second {
		t.Fatalf("Extract took %s, expected <5s", d)
	}
}
