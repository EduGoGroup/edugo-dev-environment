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

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

//go:embed templates/contract_check.md.tmpl
var templateFS embed.FS

// markdownTemplate se compila una sola vez al cargar el paquete.
var markdownTemplate = template.Must(
	template.New("contract_check.md.tmpl").
		Funcs(template.FuncMap{}).
		ParseFS(templateFS, "templates/contract_check.md.tmpl"),
)

// markdownView es el modelo plano que la plantilla consume. Se mantiene
// separado de Result para no atar el template al shape interno (lo que
// permite cambiar el JSON sin romper el Markdown).
type markdownView struct {
	GeneratedAt   string
	SchemaVersion string
	Summary       Summary
	Stats         Stats
	CategoryRows  []categoryRow
	Categories    []categoryView
	Diff          *diffView
}

type categoryRow struct {
	Name  string
	Count int
}

type categoryView struct {
	Name   string
	Notes  string
	Drifts []driftView
}

type driftView struct {
	Severity   string
	Direction  string
	Identifier string
	Detail     string
	Evidence   string
}

type diffView struct {
	Regressions []driftView
	Fixes       []driftView
}

// WriteMarkdown serializa el Result como un archivo Markdown legible
// dentro de `dir`. Devuelve la ruta absoluta del archivo creado.
func WriteMarkdown(r *Result, dir string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("report.WriteMarkdown: nil result")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("report.WriteMarkdown: mkdir %q: %w", dir, err)
	}
	path := filepath.Join(dir, markdownFileName(r.Timestamp))
	payload, err := renderMarkdown(r)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return "", fmt.Errorf("report.WriteMarkdown: write %q: %w", path, err)
	}
	return path, nil
}

// renderMarkdown compila el Result en bytes Markdown listos para
// escribir o comparar contra goldens.
func renderMarkdown(r *Result) ([]byte, error) {
	view := buildMarkdownView(r)
	var buf bytes.Buffer
	if err := markdownTemplate.Execute(&buf, view); err != nil {
		return nil, fmt.Errorf("report.renderMarkdown: %w", err)
	}
	return buf.Bytes(), nil
}

func buildMarkdownView(r *Result) markdownView {
	view := markdownView{
		GeneratedAt:   r.Timestamp.UTC().Format("2006-01-02T15:04:05Z"),
		SchemaVersion: r.SchemaVersion,
		Summary:       r.Summary,
		Stats:         r.Stats,
	}

	// Conteo por categoría (orden alfabético).
	for _, name := range sortedKeys(r.Summary.ByCategory) {
		view.CategoryRows = append(view.CategoryRows, categoryRow{
			Name:  name,
			Count: r.Summary.ByCategory[name],
		})
	}

	// Tabla por categoría: agrupar drifts y emitir en el orden del
	// catálogo (alfabético por nombre de categoría).
	groups := groupByCategory(r.Drifts)
	for _, name := range sortedCategoryNames(groups) {
		notes := ""
		if meta, ok := validate.Catalog[name]; ok {
			notes = meta.Notes
		}
		drifts := groups[name]
		view.Categories = append(view.Categories, categoryView{
			Name:   name,
			Notes:  notes,
			Drifts: toDriftViews(drifts),
		})
	}

	if r.BaselineDiff != nil {
		view.Diff = &diffView{
			Regressions: toDriftViews(r.BaselineDiff.Regressions),
			Fixes:       toDriftViews(r.BaselineDiff.Fixes),
		}
	}
	return view
}

func groupByCategory(drifts []validate.Drift) map[string][]validate.Drift {
	out := map[string][]validate.Drift{}
	for _, d := range drifts {
		out[d.Category] = append(out[d.Category], d)
	}
	return out
}

func sortedCategoryNames(groups map[string][]validate.Drift) []string {
	out := make([]string, 0, len(groups))
	for name := range groups {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func toDriftViews(ds []validate.Drift) []driftView {
	out := make([]driftView, 0, len(ds))
	for _, d := range ds {
		out = append(out, driftView{
			Severity:   escapeCell(string(d.Severity)),
			Direction:  escapeCell(string(d.Direction)),
			Identifier: escapeCell(d.Identifier),
			Detail:     escapeCell(d.Detail),
			Evidence:   escapeCell(formatEvidence(d.Evidence)),
		})
	}
	return out
}

// formatEvidence convierte la lista de Locations en una sola línea
// "file.kt:line, file2.kt:line2". Mantiene el orden ya determinista
// devuelto por validate.
func formatEvidence(locs []kmp.Location) string {
	if len(locs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(locs))
	for _, l := range locs {
		parts = append(parts, fmt.Sprintf("%s:%d", l.FilePath, l.Line))
	}
	return strings.Join(parts, ", ")
}

// escapeCell sanitiza un valor para que pueda incrustarse en una celda
// de tabla Markdown. Reemplaza separadores y saltos de línea.
func escapeCell(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}
