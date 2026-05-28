package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// FixtureLoader carga un Snapshot desde un archivo JSON local.
// Su único uso es testing y desarrollo mientras Fase A no expone su API.
//
// Implementa la interfaz Loader.
type FixtureLoader struct {
	Path string
}

// NewFixtureLoader devuelve un Loader respaldado por el archivo dado.
// El path debe apuntar a un JSON que serialice la struct Snapshot.
func NewFixtureLoader(path string) *FixtureLoader {
	return &FixtureLoader{Path: path}
}

// Load lee el fixture, lo deserializa y devuelve el Snapshot. Errores
// de I/O o de JSON se devuelven crudos para que el caller decida.
func (l *FixtureLoader) Load(ctx context.Context) (Snapshot, error) {
	if l == nil || l.Path == "" {
		return Snapshot{}, fmt.Errorf("seed.FixtureLoader: empty path")
	}
	if err := ctx.Err(); err != nil {
		return Snapshot{}, err
	}
	raw, err := os.ReadFile(l.Path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("seed.FixtureLoader: read %s: %w", l.Path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(raw, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("seed.FixtureLoader: unmarshal %s: %w", l.Path, err)
	}
	return snap, nil
}
