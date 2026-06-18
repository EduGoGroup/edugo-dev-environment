// Binario contract-check — detecta drift entre el contrato KMP del frontend
// (screenKey, apiPrefix, requiredPermission, roles) y el seed de producción
// del backend (resource_screens, iam.permissions, iam.roles, slot_data).
//
// Ejecuta el flujo:
//
//  1. Extrae artefactos del frontend KMP (regex sobre archivos *.kt).
//  2. Carga la snapshot del seed real vía Fase A (seed.ProductionLoader).
//  3. Cruza ambos universos con los 7 detectores de validate.RunAll.
//  4. Persiste reporte JSON + Markdown bajo --output-dir.
//  5. Computa diff contra baseline (si existe) o lo refresca con
//     --update-baseline.
//
// Exit codes (B-REQ-8.3):
//
//	0 -> sin drift de severidad error (puede haber warning/info).
//	1 -> al menos un drift error después del filtro --severity.
//	2 -> seed loader o KMP extractor fallaron (entorno roto).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/report"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

const (
	exitOK            = 0
	exitDriftError    = 1
	exitEnvironmental = 2
)

// errDriftErrors signals that --severity rules produced at least one
// error-severity drift. main() translates it into exit code 1.
var errDriftErrors = errors.New("drift errors present")

// defaultKMPRoots replica B-REQ-9.2. Paths relativas al monorepo
// (EduGo/) — el binario tolera que algunas no existan (B-REQ-9.3).
//
// Los módulos individuales bajo `modules/` se descubren recursivamente
// por el walker; basta listar el padre.
var defaultKMPRoots = []string{
	"EduUI/edugo-ui-kmp/kmp-screens/src/commonMain",
	"EduUI/edugo-ui-kmp/kmp-design/src/commonMain",
	"EduUI/edugo-ui-kmp/kmp-resources/src/commonMain",
	"EduUI/edugo-ui-kmp/modules",
}

const baselineFilename = "contract-check-baseline.json"

type cliFlags struct {
	kmpRoots       []string
	severity       string
	updateBaseline bool
	outputDir      string
	seedSource     string
}

func main() {
	if err := realMain(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, errDriftErrors) {
			os.Exit(exitDriftError)
		}
		fmt.Fprintf(os.Stderr, "contract-check: %v\n", err)
		os.Exit(exitEnvironmental)
	}
}

// realMain wraps flag parsing + run so tests drive the CLI without
// invoking os.Exit. Returns errDriftErrors for the strict-fail path.
func realMain(args []string, stdout, stderr io.Writer) error {
	cfg, err := parseFlags(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	return run(cfg, stdout, stderr)
}

// run orchestrates extractor → loader → validate → report. Returns
// errDriftErrors if any error-severity drift survives the --severity
// filter; any other error is environmental.
func run(cfg cliFlags, stdout, stderr io.Writer) error {
	ctx := context.Background()

	// KMP extractor and seed loader are I/O-bound and independent;
	// run them in parallel.
	var (
		kmpSnap  kmp.Snapshot
		kmpErrs  []kmp.ExtractError
		kmpErr   error
		seedSnap seed.Snapshot
		seedErr  error
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		kmpSnap, kmpErrs, kmpErr = kmp.Extract(cfg.kmpRoots)
	}()
	go func() {
		defer wg.Done()
		loader := seed.NewProductionLoader(cfg.seedSource)
		seedSnap, seedErr = loader.Load(ctx)
	}()
	wg.Wait()

	if seedErr != nil {
		return fmt.Errorf("seed loader: %w", seedErr)
	}
	if kmpErr != nil {
		return fmt.Errorf("kmp extractor: %w", kmpErr)
	}
	for _, ke := range kmpErrs {
		fmt.Fprintf(stderr, "contract-check: kmp extract warning — %s\n", ke.Error())
	}

	drifts := validate.Validate(kmpSnap, seedSnap)
	drifts = filterBySeverity(drifts, cfg.severity)

	res := report.NewResult(time.Now().UTC(), kmpSnap, seedSnap, drifts)

	baselinePath := filepath.Join(cfg.outputDir, baselineFilename)
	if !cfg.updateBaseline {
		prev, err := report.LoadBaseline(baselinePath)
		if err != nil {
			return fmt.Errorf("load baseline: %w", err)
		}
		if prev != nil {
			diff := report.ComputeDiff(prev, res)
			res.BaselineDiff = &diff
		}
	}

	if err := os.MkdirAll(cfg.outputDir, 0o755); err != nil {
		return fmt.Errorf("crear output-dir: %w", err)
	}

	jsonPath, err := report.WriteJSON(res, cfg.outputDir)
	if err != nil {
		return fmt.Errorf("escribir JSON: %w", err)
	}
	fmt.Fprintf(stdout, "contract-check: reporte JSON     → %s\n", jsonPath)

	mdPath, err := report.WriteMarkdown(res, cfg.outputDir)
	if err != nil {
		return fmt.Errorf("escribir Markdown: %w", err)
	}
	fmt.Fprintf(stdout, "contract-check: reporte Markdown → %s\n", mdPath)

	if cfg.updateBaseline {
		if err := report.UpdateBaseline(res, baselinePath); err != nil {
			return fmt.Errorf("actualizar baseline: %w", err)
		}
		fmt.Fprintf(stdout, "contract-check: baseline actualizado → %s\n", baselinePath)
	}

	fmt.Fprintf(stdout,
		"contract-check: drifts errors=%d warnings=%d infos=%d\n",
		res.Summary.Errors, res.Summary.Warnings, res.Summary.Infos,
	)

	if res.Summary.Errors > 0 && cfg.severity != "" && cfg.severity != string(validate.SeverityInfo) && cfg.severity != string(validate.SeverityWarning) {
		return errDriftErrors
	}
	if res.Summary.Errors > 0 && cfg.severity == string(validate.SeverityError) {
		return errDriftErrors
	}
	if res.Summary.Errors > 0 && cfg.severity == "" {
		// Sin filtro: el binario es informativo por default. Solo
		// falla cuando se invoca explícitamente con --severity=error
		// (ver Makefile target contract-check-strict).
		return nil
	}
	return nil
}

func parseFlags(args []string, stderr io.Writer) (cliFlags, error) {
	fs := flag.NewFlagSet("contract-check", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		kmpRootsCSV    string
		severity       string
		updateBaseline bool
		outputDir      string
		seedSource     string
	)

	fs.StringVar(&kmpRootsCSV, "kmp-roots", "",
		"Comma-separated list of paths to scan for *.kt files. Empty -> default list (see B-REQ-9.2).")
	fs.StringVar(&severity, "severity", "",
		"Severity filter: error -> include only error drifts; warning -> error+warning; info|empty -> include all.")
	fs.BoolVar(&updateBaseline, "update-baseline", false,
		"If set, overwrite contract-check-baseline.json with the current run (omits Regressions/Fixes).")
	fs.StringVar(&outputDir, "output-dir", "audit-reports",
		"Directory where contract-check-<ts>.{json,md} are written. Created if missing.")
	fs.StringVar(&seedSource, "seed-source", "",
		"Optional override for the seed loader source (delegated to Fase A).")

	if err := fs.Parse(args); err != nil {
		return cliFlags{}, err
	}

	switch severity {
	case "", "error", "warning", "info":
	default:
		return cliFlags{}, fmt.Errorf("--severity inválido %q (válidos: error|warning|info o vacío)", severity)
	}

	roots := defaultKMPRoots
	if strings.TrimSpace(kmpRootsCSV) != "" {
		roots = splitAndTrim(kmpRootsCSV)
	}

	return cliFlags{
		kmpRoots:       roots,
		severity:       severity,
		updateBaseline: updateBaseline,
		outputDir:      outputDir,
		seedSource:     seedSource,
	}, nil
}

func splitAndTrim(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// filterBySeverity drops drifts strictly below the requested floor.
// Empty filter keeps all drifts.
//
//	"error"   → keep only error
//	"warning" → keep error + warning
//	"info"    → keep all (same as empty)
func filterBySeverity(in []validate.Drift, severity string) []validate.Drift {
	if severity == "" || severity == string(validate.SeverityInfo) {
		return in
	}
	out := make([]validate.Drift, 0, len(in))
	for _, d := range in {
		if severity == string(validate.SeverityError) && d.Severity == validate.SeverityError {
			out = append(out, d)
			continue
		}
		if severity == string(validate.SeverityWarning) && d.Severity != validate.SeverityInfo {
			out = append(out, d)
		}
	}
	return out
}
