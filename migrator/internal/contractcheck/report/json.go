package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// timestampLayout es el layout RFC3339 compactado usado en el nombre de
// archivo del reporte JSON/Markdown (sin separadores `-` ni `:`).
const timestampLayout = "20060102T150405Z"

// jsonFileName devuelve el nombre canónico del archivo JSON del reporte
// para un timestamp dado.
func jsonFileName(ts time.Time) string {
	return "contract-check-" + ts.UTC().Format(timestampLayout) + ".json"
}

// markdownFileName devuelve el nombre canónico del archivo Markdown del
// reporte para un timestamp dado.
func markdownFileName(ts time.Time) string {
	return "contract-check-" + ts.UTC().Format(timestampLayout) + ".md"
}

// WriteJSON serializa el Result completo a `dir/contract-check-<ts>.json`
// y devuelve la ruta absoluta del archivo escrito. El directorio se
// crea si no existe (B-REQ-12.3).
//
// Determinismo: dos llamadas con el mismo Result (mismo Timestamp) y los
// mismos drifts producen archivos byte-idénticos, porque las llaves de
// mapas se serializan con encoding/json (ordena keys de map[string]X
// alfabéticamente por defecto) y los slices ya vienen ordenados desde
// validate.Validate.
func WriteJSON(r *Result, dir string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("report.WriteJSON: nil result")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("report.WriteJSON: mkdir %q: %w", dir, err)
	}
	path := filepath.Join(dir, jsonFileName(r.Timestamp))
	payload, err := marshalResult(r)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return "", fmt.Errorf("report.WriteJSON: write %q: %w", path, err)
	}
	return path, nil
}

// UpdateBaseline serializa el Result en `path` con la sección
// BaselineDiff omitida (porque este archivo PASA a ser el nuevo
// baseline; el diff se computará en runs futuros contra él). El
// directorio padre se crea si no existe.
func UpdateBaseline(r *Result, path string) error {
	if r == nil {
		return fmt.Errorf("report.UpdateBaseline: nil result")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("report.UpdateBaseline: mkdir: %w", err)
	}
	clone := *r
	clone.BaselineDiff = nil
	payload, err := marshalResult(&clone)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return fmt.Errorf("report.UpdateBaseline: write %q: %w", path, err)
	}
	return nil
}

// marshalResult emite el JSON canónico (indent 2 espacios, llaves de
// mapas ordenadas alfabéticamente por encoding/json, salto de línea
// final). Los snapshots completos NO se incluyen por defecto: si el
// caller asignó r.KMPSnapshot/r.SeedSnapshot manualmente se serializan
// igual.
func marshalResult(r *Result) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return nil, fmt.Errorf("report.marshalResult: %w", err)
	}
	// json.Encoder ya añade un \n final; nada más que hacer.
	return buf.Bytes(), nil
}
