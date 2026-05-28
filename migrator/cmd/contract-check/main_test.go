package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

func TestParseFlags_Defaults(t *testing.T) {
	var stderr bytes.Buffer
	cfg, err := parseFlags(nil, &stderr)
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}
	if len(cfg.kmpRoots) == 0 {
		t.Fatal("expected default KMP roots")
	}
	if cfg.severity != "" {
		t.Fatalf("default severity should be empty, got %q", cfg.severity)
	}
}

func TestParseFlags_InvalidSeverity(t *testing.T) {
	var stderr bytes.Buffer
	_, err := parseFlags([]string{"--severity", "loud"}, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid severity")
	}
	if !strings.Contains(err.Error(), "--severity inválido") {
		t.Fatalf("error %q lacks expected prefix", err.Error())
	}
}

func TestParseFlags_KMPRootsCSV(t *testing.T) {
	var stderr bytes.Buffer
	cfg, err := parseFlags([]string{"--kmp-roots", "a, b ,c"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags: %v", err)
	}
	if len(cfg.kmpRoots) != 3 {
		t.Fatalf("expected 3 roots, got %v", cfg.kmpRoots)
	}
	for i, w := range []string{"a", "b", "c"} {
		if cfg.kmpRoots[i] != w {
			t.Errorf("roots[%d] = %q, want %q", i, cfg.kmpRoots[i], w)
		}
	}
}

func TestRealMain_UnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := realMain([]string{"--bogus"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

// TestRealMain_AgainstFixtures runs the full pipeline against a tiny
// KMP fixture so we exercise extractor + validator + reporter without
// pulling in the real seed.
//
// We use an in-tree fixture under contractcheck/kmp/testdata; the seed
// loader still hits production, which is fine: the test only checks
// that the CLI completes, writes both files, and respects --severity=info
// so the abundant real drifts don't fail the test.
func TestRealMain_AgainstFixturesNoStrict(t *testing.T) {
	dir := t.TempDir()
	args := []string{
		"--kmp-roots", "../../internal/contractcheck/testdata/kmp/happy",
		"--output-dir", dir,
		"--severity", "info",
	}
	var stdout, stderr bytes.Buffer
	if err := realMain(args, &stdout, &stderr); err != nil {
		t.Fatalf("run failed: %v\nstdout=%s\nstderr=%s", err, stdout.String(), stderr.String())
	}
	json, _ := filepath.Glob(filepath.Join(dir, "contract-check-*.json"))
	md, _ := filepath.Glob(filepath.Join(dir, "contract-check-*.md"))
	if len(json) != 1 || len(md) != 1 {
		t.Fatalf("expected 1 JSON + 1 MD, got %d/%d (stdout=%s)", len(json), len(md), stdout.String())
	}
}

// TestRealMain_StrictAgainstRealSeed expects a non-zero count of error
// drifts today and asserts that --severity=error returns errDriftErrors.
// If the seed/FE eventually align, the test should be re-checked.
func TestRealMain_StrictAgainstRealSeed_BaselineFailsToday(t *testing.T) {
	dir := t.TempDir()
	args := []string{
		"--kmp-roots", "../../internal/contractcheck/testdata/kmp/phantom-screen",
		"--output-dir", dir,
		"--severity", "error",
	}
	var stdout, stderr bytes.Buffer
	err := realMain(args, &stdout, &stderr)
	// Phantom-screen fixture references a screenKey not in seed → at
	// least one error drift expected.
	if err == nil {
		t.Fatalf("expected errDriftErrors against phantom-screen fixture; stdout=%s", stdout.String())
	}
	if !errors.Is(err, errDriftErrors) {
		t.Fatalf("expected errDriftErrors, got %v", err)
	}
}

func TestRealMain_OutputDirCreated(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "nested", "audit-reports")
	args := []string{
		"--kmp-roots", "../../internal/contractcheck/testdata/kmp/happy",
		"--output-dir", target,
		"--severity", "info",
	}
	var stdout, stderr bytes.Buffer
	if err := realMain(args, &stdout, &stderr); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		t.Fatalf("expected dir created, stat err=%v info=%v", err, info)
	}
}

func TestFilterBySeverity(t *testing.T) {
	cases := []struct {
		name   string
		filter string
		err    int
		warn   int
		info   int
		want   int
	}{
		{"empty keeps all", "", 1, 1, 1, 3},
		{"info keeps all", "info", 1, 1, 1, 3},
		{"warning drops infos", "warning", 1, 1, 1, 2},
		{"error drops warnings and infos", "error", 1, 1, 1, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			drifts := makeDrifts(tc.err, tc.warn, tc.info)
			got := filterBySeverity(drifts, tc.filter)
			if len(got) != tc.want {
				t.Errorf("filter=%q want=%d got=%d", tc.filter, tc.want, len(got))
			}
		})
	}
}

func makeDrifts(errs, warns, infos int) []validate.Drift {
	out := make([]validate.Drift, 0, errs+warns+infos)
	for i := 0; i < errs; i++ {
		out = append(out, validate.Drift{Severity: validate.SeverityError})
	}
	for i := 0; i < warns; i++ {
		out = append(out, validate.Drift{Severity: validate.SeverityWarning})
	}
	for i := 0; i < infos; i++ {
		out = append(out, validate.Drift{Severity: validate.SeverityInfo})
	}
	return out
}
