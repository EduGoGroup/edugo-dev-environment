package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestBinaryExitCodes satisfies B-REQ-8.3 / task 6.6: build the binary
// and invoke it via exec.Command to verify the three documented exit
// codes (0, 1, 2). We use the in-tree fixture roots so the test does
// not depend on the production seed having any specific drift profile.
func TestBinaryExitCodes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("exit-code expectations match POSIX semantics")
	}
	binPath := buildBinary(t)
	fixtureHappy := absPath(t, "../../internal/contractcheck/testdata/kmp/happy")
	fixturePhantom := absPath(t, "../../internal/contractcheck/testdata/kmp/phantom-screen")

	cases := []struct {
		name     string
		args     []string
		wantExit int
	}{
		{
			name:     "info severity → exit 0 even with drifts",
			args:     []string{"--kmp-roots", fixtureHappy, "--severity", "info"},
			wantExit: 0,
		},
		{
			name:     "error severity on phantom-screen fixture → exit 1",
			args:     []string{"--kmp-roots", fixturePhantom, "--severity", "error"},
			wantExit: 1,
		},
		{
			name:     "invalid flag → exit 2",
			args:     []string{"--bogus-flag"},
			wantExit: 2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outDir := t.TempDir()
			args := append([]string{}, tc.args...)
			args = append(args, "--output-dir", outDir)
			cmd := exec.Command(binPath, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			gotExit := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					gotExit = exitErr.ExitCode()
				} else {
					t.Fatalf("unexpected error type: %v", err)
				}
			}
			if gotExit != tc.wantExit {
				t.Fatalf("exit code = %d, want %d", gotExit, tc.wantExit)
			}
		})
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "contract-check-test")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
	return bin
}

func absPath(t *testing.T, rel string) string {
	t.Helper()
	abs, err := filepath.Abs(rel)
	if err != nil {
		t.Fatalf("abs(%s): %v", rel, err)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("fixture missing: %s — %v", abs, err)
	}
	if !strings.HasSuffix(abs, rel[strings.Index(rel, "testdata"):]) {
		t.Logf("absolute path: %s", abs)
	}
	return abs
}
