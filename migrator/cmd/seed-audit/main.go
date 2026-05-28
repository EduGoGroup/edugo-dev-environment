// Binario seed-audit — auditor estático del production seed.
//
// Carga el seed en memoria, ejecuta los siete validadores estáticos
// declarados en `phase-a-static-auditor` y persiste un AuditReport en
// formato JSON y/o Markdown bajo `--output-dir`. El exit code se
// determina según `--strict` / `--report-only` (Decisión D-7).
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/loader"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/report"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/seedaudit/validators"
)

// Version sigue el `AuditReport.SchemaVersion`. Bump conjunto cuando
// el contrato del reporte cambie.
const Version = "0.1.0-dev"

const (
	formatJSON = "json"
	formatMD   = "md"
	formatBoth = "both"
)

// errStrictViolations is returned by run() when --strict is active and
// the report contains at least one error-severity violation. main()
// translates it into exit code 1.
var errStrictViolations = errors.New("strict mode: error-severity violations present")

// options captures the parsed CLI flags. Keep this struct serialisable
// and free of side effects so run() stays unit-testable.
type options struct {
	seedSource string
	outputDir  string
	format     string
	strict     bool
	reportOnly bool
}

func main() {
	if err := realMain(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, errStrictViolations) {
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "seed-audit: %v\n", err)
		os.Exit(2)
	}
}

// realMain wraps flag parsing + run so tests can drive the CLI without
// invoking os.Exit. It returns errStrictViolations for the strict-fail
// path and other errors for fatal paths.
func realMain(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("seed-audit", flag.ContinueOnError)
	fs.SetOutput(stderr)

	showVersion := fs.Bool("version", false, "imprime la versión del binario y termina")
	seedSource := fs.String("seed-source", loader.SeedSourceProduction, "fuente del seed (v1: solo `production`)")
	outputDir := fs.String("output-dir", "./audit-reports", "directorio donde se persisten los reportes")
	format := fs.String("format", formatBoth, "formato del reporte: json | md | both")
	strict := fs.Bool("strict", false, "exit 1 si el reporte contiene violaciones severity=error")
	reportOnly := fs.Bool("report-only", true, "modo no-bloqueante (exit 0 siempre que el binario complete)")

	reportOnlyExplicit := false
	if err := fs.Parse(args); err != nil {
		return err
	}
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "report-only" {
			reportOnlyExplicit = true
		}
	})

	if *showVersion {
		fmt.Fprintf(stdout, "seed-audit %s\n", Version)
		return nil
	}

	if err := validateFormat(*format); err != nil {
		return err
	}

	if *strict && reportOnlyExplicit && *reportOnly {
		// Decisión D-7: gana --strict, warning a stderr.
		fmt.Fprintln(stderr, "seed-audit: aviso — --report-only se ignora porque --strict está activo")
	}

	opts := options{
		seedSource: *seedSource,
		outputDir:  *outputDir,
		format:     *format,
		strict:     *strict,
		reportOnly: *reportOnly,
	}
	return run(opts, stdout)
}

// run orchestrates loader → validators → reporter and writes the
// requested formats. It returns errStrictViolations when --strict is
// active and the report contains error-severity violations; any other
// error is fatal.
func run(opts options, stdout io.Writer) error {
	snap, err := loader.Load(loader.RunOptions{SeedSource: opts.seedSource})
	if err != nil {
		return fmt.Errorf("cargar seed: %w", err)
	}

	violations := validators.RunAll(snap)
	rep := report.Build(snap, violations, opts.seedSource)

	if err := os.MkdirAll(opts.outputDir, 0o755); err != nil {
		return fmt.Errorf("crear output-dir: %w", err)
	}

	if opts.format == formatJSON || opts.format == formatBoth {
		path, err := report.WriteJSON(rep, opts.outputDir)
		if err != nil {
			return fmt.Errorf("escribir JSON: %w", err)
		}
		fmt.Fprintf(stdout, "seed-audit: reporte JSON     → %s\n", path)
	}
	if opts.format == formatMD || opts.format == formatBoth {
		path, err := report.WriteMarkdown(rep, opts.outputDir)
		if err != nil {
			return fmt.Errorf("escribir Markdown: %w", err)
		}
		fmt.Fprintf(stdout, "seed-audit: reporte Markdown → %s\n", path)
	}

	fmt.Fprintf(stdout,
		"seed-audit: violaciones errors=%d warnings=%d infos=%d\n",
		rep.Summary.Errors, rep.Summary.Warnings, rep.Summary.Infos,
	)

	if opts.strict && rep.Summary.Errors > 0 {
		return errStrictViolations
	}
	return nil
}

func validateFormat(f string) error {
	switch f {
	case formatJSON, formatMD, formatBoth:
		return nil
	default:
		return fmt.Errorf("formato inválido %q (válidos: json, md, both)", f)
	}
}
