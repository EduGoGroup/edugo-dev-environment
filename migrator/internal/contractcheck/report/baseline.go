package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/validate"
)

// LoadBaseline lee el archivo baseline JSON previo. Si el archivo no
// existe (primer run, o antes de `--update-baseline`), devuelve `nil,
// nil`: el caller tratará la ausencia como "sin baseline" y emitirá el
// reporte sin las secciones Regressions/Fixes (B-REQ-11.1, design §5).
//
// Cualquier otro error de I/O o de parsing se propaga al caller para
// que decida (típicamente: log warning + continuar sin diff).
func LoadBaseline(path string) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("report.LoadBaseline: read %q: %w", path, err)
	}
	var r Result
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("report.LoadBaseline: parse %q: %w", path, err)
	}
	return &r, nil
}

// ComputeDiff cruza los drifts del baseline (`prev`) con los del run
// actual (`curr`) y devuelve regresiones (drifts nuevos en curr) y
// fixes (drifts del baseline que ya no aparecen). El match es por la
// tupla `(Category, Identifier)` para que un cambio de severidad — de
// warning a error — cuente como regresión nueva del lado humano sin
// que el algoritmo lo trate igual: lo registramos como regresión sólo
// si el (Category, Identifier) no estaba en prev.
//
// Si `prev` es nil, devolvemos un BaselineDiff vacío (sin regresiones
// ni fixes). El reporter omite la sección si Regressions y Fixes están
// ambos vacíos y el caller no asignó BaselineDiff al Result.
func ComputeDiff(prev, curr *Result) BaselineDiff {
	diff := BaselineDiff{}
	if prev == nil || curr == nil {
		return diff
	}
	prevSet := indexByKey(prev.Drifts)
	currSet := indexByKey(curr.Drifts)

	for k, d := range currSet {
		if _, ok := prevSet[k]; !ok {
			diff.Regressions = append(diff.Regressions, d)
		}
	}
	for k, d := range prevSet {
		if _, ok := currSet[k]; !ok {
			diff.Fixes = append(diff.Fixes, d)
		}
	}
	sortDrifts(diff.Regressions)
	sortDrifts(diff.Fixes)
	return diff
}

// driftKey compone la clave canónica usada por ComputeDiff para
// emparejar drifts entre runs. Usa categoría e identificador (ignora
// severidad e idioma del Detail) para que la heurística sea estable a
// reformulaciones del mensaje.
func driftKey(d validate.Drift) string {
	return d.Category + "\x1f" + d.Identifier
}

func indexByKey(drifts []validate.Drift) map[string]validate.Drift {
	out := make(map[string]validate.Drift, len(drifts))
	for _, d := range drifts {
		out[driftKey(d)] = d
	}
	return out
}

// sortDrifts replica el orden canónico que validate.Validate aplica
// (B-REQ-10.1: Category asc, Severity desc, Identifier asc) para que
// los Regressions/Fixes salgan determinísticos en JSON y Markdown.
func sortDrifts(drifts []validate.Drift) {
	severityRank := map[validate.Severity]int{
		validate.SeverityError:   0,
		validate.SeverityWarning: 1,
		validate.SeverityInfo:    2,
	}
	sort.SliceStable(drifts, func(i, j int) bool {
		a, b := drifts[i], drifts[j]
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		ra, rb := severityRank[a.Severity], severityRank[b.Severity]
		if ra != rb {
			return ra < rb
		}
		return a.Identifier < b.Identifier
	})
}
