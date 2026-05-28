package validate

import (
	"fmt"
	"sort"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// detectDeadScreenKeys implementa B-REQ-2: resource_screens.screen_key
// que el seed declara pero que ningún archivo Kotlin atiende.
//
// Severidad por defecto: `warning`. Se eleva a `error` cuando el seed
// declara `screen_type=dashboard` o `is_default=true` (B-REQ-2.2): un
// dashboard inalcanzable es un bug bloqueante.
func detectDeadScreenKeys(k kmp.Snapshot, s seed.Snapshot) []Drift {
	feKeys := make(map[string]struct{}, len(k.ScreenKeys))
	for key := range k.ScreenKeys {
		feKeys[key] = struct{}{}
	}

	// Agrupar por screenKey conservando los flags más severos.
	type meta struct {
		screenType string
		isDefault  bool
		critical   bool
	}
	grouped := make(map[string]meta)
	order := make([]string, 0)
	for _, rs := range s.ResourceScreens {
		if _, ok := feKeys[rs.ScreenKey]; ok {
			continue
		}
		m, exists := grouped[rs.ScreenKey]
		if !exists {
			order = append(order, rs.ScreenKey)
			m = meta{screenType: rs.ScreenType, isDefault: rs.IsDefault}
		}
		if rs.ScreenType == "dashboard" {
			m.screenType = "dashboard"
		}
		if rs.IsDefault {
			m.isDefault = true
		}
		if m.screenType == "dashboard" || m.isDefault {
			m.critical = true
		}
		grouped[rs.ScreenKey] = m
	}
	sort.Strings(order)

	out := make([]Drift, 0, len(order))
	for _, key := range order {
		m := grouped[key]
		severity := SeverityFor(CategoryScreenKeyDead)
		detail := fmt.Sprintf(
			"El seed declara resource_screens.screen_key %q pero ningún composable KMP lo implementa.",
			key,
		)
		if m.critical {
			severity = SeverityError
			detail = fmt.Sprintf(
				"%s screen_type=%q is_default=%t — pantalla crítica inalcanzable.",
				detail, m.screenType, m.isDefault,
			)
		}
		out = append(out, Drift{
			Direction:  DirectionBEOnly,
			Category:   CategoryScreenKeyDead,
			Severity:   severity,
			Identifier: key,
			Detail:     detail,
		})
	}
	return out
}
