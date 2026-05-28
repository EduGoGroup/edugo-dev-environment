package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// timestampLayout is the compact RFC 3339 layout used in report
// filenames (no separators, UTC, second resolution).
const timestampLayout = "20060102T150405Z"

// WriteJSON serialises r as an indented JSON document under dir. The
// filename embeds r.GeneratedAt formatted as RFC 3339 compact
// ("seed-audit-20060102T150405Z.json"). The output directory is
// created with mode 0o755 if missing (A-REQ-8.4).
//
// Encoding rules (A-REQ-10.2):
//   - 2-space indentation, UTF-8.
//   - Trailing newline so the file is git-friendly.
//   - Map keys are emitted sorted (encoding/json default) so
//     Summary.ByCode is deterministic.
//
// On success the absolute path of the written file is returned.
func WriteJSON(r *AuditReport, dir string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("report.WriteJSON: nil report")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("report.WriteJSON: create dir %q: %w", dir, err)
	}

	name := "seed-audit-" + r.GeneratedAt.UTC().Format(timestampLayout) + ".json"
	path := filepath.Join(dir, name)

	buf, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("report.WriteJSON: marshal: %w", err)
	}
	buf = append(buf, '\n')

	if err := os.WriteFile(path, buf, 0o644); err != nil {
		return "", fmt.Errorf("report.WriteJSON: write %q: %w", path, err)
	}
	return path, nil
}
