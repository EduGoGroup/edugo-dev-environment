package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRealMain_Version(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := realMain([]string{"--version"}, &stdout, &stderr); err != nil {
		t.Fatalf("--version returned error: %v", err)
	}
	if !strings.Contains(stdout.String(), "seed-audit "+Version) {
		t.Fatalf("stdout %q does not contain version banner", stdout.String())
	}
}

func TestRealMain_InvalidFormat(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := realMain([]string{"--format", "yaml"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
	if !strings.Contains(err.Error(), "formato inválido") {
		t.Fatalf("error %q lacks expected prefix", err.Error())
	}
}

func TestRealMain_UnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := realMain([]string{"--bogus"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected flag.ErrHelp-style error for unknown flag")
	}
}

func TestRealMain_ReportOnlySuccess(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	args := []string{"--output-dir", dir, "--format", "json"}
	if err := realMain(args, &stdout, &stderr); err != nil {
		t.Fatalf("default report-only run failed: %v (stderr=%s)", err, stderr.String())
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "seed-audit-*.json"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 JSON report, got %d (%v)", len(matches), matches)
	}
	mdMatches, _ := filepath.Glob(filepath.Join(dir, "seed-audit-*.md"))
	if len(mdMatches) != 0 {
		t.Fatalf("--format=json should not emit Markdown, got %v", mdMatches)
	}
}

func TestRealMain_BothFormats(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	if err := realMain([]string{"--output-dir", dir}, &stdout, &stderr); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	json, _ := filepath.Glob(filepath.Join(dir, "seed-audit-*.json"))
	md, _ := filepath.Glob(filepath.Join(dir, "seed-audit-*.md"))
	if len(json) != 1 || len(md) != 1 {
		t.Fatalf("expected 1 JSON + 1 MD, got %d/%d", len(json), len(md))
	}
}

func TestRealMain_StrictAgainstProductionSeed(t *testing.T) {
	// Smoke: against the real production seed --strict should not
	// fail today. If it does, the seed has regressed below the
	// expected baseline (zero error-severity violations as of v1).
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	err := realMain([]string{"--output-dir", dir, "--strict", "--format", "json"}, &stdout, &stderr)
	if err != nil && !errors.Is(err, errStrictViolations) {
		t.Fatalf("strict run hit a fatal error: %v", err)
	}
	if errors.Is(err, errStrictViolations) {
		// Surface the report path so a developer can investigate.
		t.Fatalf("--strict found error-severity violations against the seed; check %s", dir)
	}
}

func TestRealMain_StrictReportOnlyConflictWarning(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	args := []string{"--output-dir", dir, "--strict", "--report-only=true", "--format", "json"}
	_ = realMain(args, &stdout, &stderr)
	if !strings.Contains(stderr.String(), "--report-only se ignora") {
		t.Fatalf("expected D-7 warning on stderr, got %q", stderr.String())
	}
}

func TestRealMain_OutputDirCreated(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "nested", "audit-reports")
	var stdout, stderr bytes.Buffer
	if err := realMain([]string{"--output-dir", target, "--format", "json"}, &stdout, &stderr); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("expected --output-dir created, stat err: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected directory, got %v", info.Mode())
	}
}

func TestRealMain_ReproducibleJSON(t *testing.T) {
	// A-REQ-10.2/10.3: two consecutive runs must yield byte-identical
	// JSON modulo `generated_at`. We strip the line that contains the
	// timestamp before comparing so we don't tightly couple the test
	// to the exact serialisation format.
	dir1, dir2 := t.TempDir(), t.TempDir()
	var stdout1, stderr1, stdout2, stderr2 bytes.Buffer
	if err := realMain([]string{"--output-dir", dir1, "--format", "json"}, &stdout1, &stderr1); err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	if err := realMain([]string{"--output-dir", dir2, "--format", "json"}, &stdout2, &stderr2); err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	a := mustReadOnlyJSON(t, dir1)
	b := mustReadOnlyJSON(t, dir2)
	if stripTimestamp(a) != stripTimestamp(b) {
		t.Fatalf("JSON outputs diverge across runs (excluding timestamps)")
	}
}

func mustReadOnlyJSON(t *testing.T, dir string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, "seed-audit-*.json"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("expected 1 JSON report in %s, got %v (err=%v)", dir, matches, err)
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read %s: %v", matches[0], err)
	}
	return string(data)
}

// TestRealMain_Help_ListsAllFlags satisfies A-REQ-11.1/11.4 (test 7.4):
// seed-audit --help must print every documented flag so CI runners can
// discover the surface area without reading source.
func TestRealMain_Help_ListsAllFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := realMain([]string{"--help"}, &stdout, &stderr)
	// flag.ContinueOnError prints usage to stderr and returns flag.ErrHelp,
	// which realMain wraps as a generic error. Either way help text lands
	// on stderr and we must see all flag names there.
	help := stderr.String()
	if err == nil && help == "" {
		t.Fatal("expected --help to emit usage on stderr")
	}
	for _, want := range []string{
		"-seed-source",
		"-output-dir",
		"-format",
		"-strict",
		"-report-only",
		"-version",
	} {
		if !strings.Contains(help, want) {
			t.Errorf("--help output missing flag %q\nstderr=%s", want, help)
		}
	}
}

// TestRealMain_Performance_LessThan30s satisfies A-REQ-10.1 (test 7.3):
// a full run against the production seed must complete well under 30s.
// We measure with a generous wall-clock budget (10s) so the test does
// not flake on slow CI runners — anything close to 30s would already be
// a real regression.
func TestRealMain_Performance_LessThan30s(t *testing.T) {
	if testing.Short() {
		t.Skip("performance check skipped under -short")
	}
	dir := t.TempDir()
	start := time.Now()
	var stdout, stderr bytes.Buffer
	if err := realMain([]string{"--output-dir", dir, "--format", "json"}, &stdout, &stderr); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if d := time.Since(start); d > 10*time.Second {
		t.Fatalf("seed-audit run took %s, expected <10s (budget against 30s spec)", d)
	}
}

func stripTimestamp(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, `"generated_at"`) {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}
