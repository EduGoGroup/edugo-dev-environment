package validate

import (
	"fmt"
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// detectPhantomScreenKeys implementa B-REQ-1: screenKeys que el FE
// declara pero que no existen en `resource_screens.screen_key` del seed.
//
// La severidad por defecto es `error` (Catalog) — un screenKey fantasma
// significa que el composable nunca puede cargar su contenido.
func detectPhantomScreenKeys(k kmp.Snapshot, s seed.Snapshot) []Drift {
	seedKeys := make(map[string]struct{}, len(s.ResourceScreens))
	for _, rs := range s.ResourceScreens {
		seedKeys[rs.ScreenKey] = struct{}{}
	}

	keys := make([]string, 0, len(k.ScreenKeys))
	for key := range k.ScreenKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]Drift, 0)
	for _, key := range keys {
		if _, ok := seedKeys[key]; ok {
			continue
		}
		if isPreAuthScreen(key) {
			continue
		}
		if isStaticCompliantScreen(key) {
			continue
		}
		out = append(out, Drift{
			Direction:  DirectionFEOnly,
			Category:   CategoryScreenKeyPhantom,
			Severity:   SeverityFor(CategoryScreenKeyPhantom),
			Identifier: key,
			Detail: fmt.Sprintf(
				"El frontend declara screenKey %q pero el seed de producción no tiene ningún resource_screens.screen_key con ese valor.",
				key,
			),
			Evidence: dedupLocations(k.ScreenKeys[key]),
		})
	}
	return out
}

// dedupLocations colapsa duplicados (mismo file/line/snippet) y ordena
// por (FilePath, Line) para que el reporte sea determinista.
func dedupLocations(locs []kmp.Location) []kmp.Location {
	if len(locs) == 0 {
		return nil
	}
	seen := make(map[kmp.Location]struct{}, len(locs))
	out := make([]kmp.Location, 0, len(locs))
	for _, l := range locs {
		if _, ok := seen[l]; ok {
			continue
		}
		seen[l] = struct{}{}
		out = append(out, l)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].FilePath != out[j].FilePath {
			return out[i].FilePath < out[j].FilePath
		}
		return out[i].Line < out[j].Line
	})
	return out
}
