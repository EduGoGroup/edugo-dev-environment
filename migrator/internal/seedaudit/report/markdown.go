package report

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/audit_report.md.tmpl
var templatesFS embed.FS

// markdownTemplate is parsed once on package init. The template lives
// in templates/audit_report.md.tmpl per Decision D-2.
var markdownTemplate = template.Must(
	template.New("audit_report.md.tmpl").
		ParseFS(templatesFS, "templates/audit_report.md.tmpl"),
)

// byCodeRow is a flattened, sorted entry for the ByCode table.
type byCodeRow struct {
	Code  string
	Count int
}

// violationRow is the per-violation projection shown in the markdown
// tables. References is pre-rendered into a deterministic string so
// the template never iterates a map directly.
type violationRow struct {
	Code       string
	Entity     string
	EntityID   string
	Message    string
	References string
	Path       string
}

// violationGroup holds the rows of one severity bucket. Heading is
// the human-readable section title (Spanish); Count is len(Items).
type violationGroup struct {
	Severity Severity
	Heading  string
	Count    int
	Items    []violationRow
}

// markdownView is the data passed to the template. We keep all
// presentation-level transformations here (sort, escape, format) so
// the template stays declarative and the renderer remains testable.
type markdownView struct {
	SchemaVersion  string
	SeedSource     string
	GeneratedAtRFC string
	Stats          Stats
	Summary        Summary
	ByCodeRows     []byCodeRow
	Groups         []violationGroup
	HasViolations  bool
}

var severityHeadings = map[Severity]string{
	SeverityError:   "Errores",
	SeverityWarning: "Advertencias",
	SeverityInfo:    "Informativos",
}

// severityOrder defines the bucket order in the rendered report.
// Mirrors Build's sorting so visual scan and JSON ordering agree.
var severityOrder = []Severity{SeverityError, SeverityWarning, SeverityInfo}

// renderReferences serialises a References map into a deterministic
// "k1=v1, k2=v2" string. Markdown pipe characters are escaped so the
// table cell stays well-formed even when a reference contains "|".
func renderReferences(refs map[string]string) string {
	if len(refs) == 0 {
		return ""
	}
	keys := make([]string, 0, len(refs))
	for k := range refs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", mdEscape(k), mdEscape(refs[k])))
	}
	return strings.Join(parts, ", ")
}

// mdEscape neutralises the few characters that would corrupt the
// table layout. Newlines collapse to spaces; pipes are backslash-
// escaped so they stay inside the same cell.
func mdEscape(s string) string {
	if s == "" {
		return ""
	}
	r := strings.NewReplacer(
		"|", `\|`,
		"\r\n", " ",
		"\n", " ",
		"\r", " ",
	)
	return r.Replace(s)
}

func buildView(r *AuditReport) markdownView {
	v := markdownView{
		SchemaVersion:  r.SchemaVersion,
		SeedSource:     r.SeedSource,
		GeneratedAtRFC: r.GeneratedAt.UTC().Format(time.RFC3339),
		Stats:          r.Stats,
		Summary:        r.Summary,
		HasViolations:  len(r.Violations) > 0,
	}

	// ByCode rows: sorted alphabetically (deterministic, matches
	// JSON map-key ordering).
	if len(r.Summary.ByCode) > 0 {
		codes := make([]string, 0, len(r.Summary.ByCode))
		for c := range r.Summary.ByCode {
			codes = append(codes, c)
		}
		sort.Strings(codes)
		v.ByCodeRows = make([]byCodeRow, 0, len(codes))
		for _, c := range codes {
			v.ByCodeRows = append(v.ByCodeRows, byCodeRow{Code: c, Count: r.Summary.ByCode[c]})
		}
	}

	// Group violations by severity. r.Violations is already sorted
	// by Build, so we preserve order within each bucket.
	buckets := make(map[Severity][]violationRow)
	for _, vio := range r.Violations {
		buckets[vio.Severity] = append(buckets[vio.Severity], violationRow{
			Code:       vio.Code,
			Entity:     mdEscape(vio.Entity),
			EntityID:   mdEscape(vio.EntityID),
			Message:    mdEscape(vio.Message),
			References: renderReferences(vio.References),
			Path:       mdEscape(vio.Path),
		})
	}
	for _, sev := range severityOrder {
		rows := buckets[sev]
		if len(rows) == 0 {
			continue
		}
		v.Groups = append(v.Groups, violationGroup{
			Severity: sev,
			Heading:  severityHeadings[sev],
			Count:    len(rows),
			Items:    rows,
		})
	}
	return v
}

// WriteMarkdown renders r as Markdown and writes it under dir using
// the same timestamped filename convention as WriteJSON
// ("seed-audit-<RFC3339-compact>.md"). The template lives in
// templates/audit_report.md.tmpl (Decision D-2) and is embedded at
// build time.
//
// The renderer pre-flattens every map (Summary.ByCode,
// Violation.References) into sorted slices so identical inputs
// always produce identical bytes (A-REQ-10.2).
func WriteMarkdown(r *AuditReport, dir string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("report.WriteMarkdown: nil report")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("report.WriteMarkdown: create dir %q: %w", dir, err)
	}

	view := buildView(r)
	var buf bytes.Buffer
	if err := markdownTemplate.Execute(&buf, view); err != nil {
		return "", fmt.Errorf("report.WriteMarkdown: render: %w", err)
	}
	out := buf.Bytes()
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}

	name := "seed-audit-" + r.GeneratedAt.UTC().Format(timestampLayout) + ".md"
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return "", fmt.Errorf("report.WriteMarkdown: write %q: %w", path, err)
	}
	return path, nil
}
