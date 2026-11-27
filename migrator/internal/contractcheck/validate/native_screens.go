package validate

import "strings"

// nativeScreenWhitelist contiene screenKeys que el seed mapea en
// resource_screens (el menú los resuelve) pero que KMP pinta con una
// pantalla nativa: `MainScreen.kt` los intercepta por
// `currentScreenKey == "..."` o `RouteRegistry.kt` los registra, sin un
// `override val screenKey = "..."` que el extractor pueda leer.
//
// Estas screens están exentas del drift screen_key_dead.
//
// Criterio de admisión (deuda 194, decidido por el dueño el 2026-11-26):
//   - El seed las declara en resource_screens.
//   - En KMP hay una rama nativa que las pinta: un
//     `currentScreenKey == "<key>"` en MainScreen.kt o un
//     `register("<key>")` en RouteRegistry.kt (cítalo al añadirla).
//   - No existe un ScreenContract con ese screenKey (si existe, el
//     extractor ya la ve y no hace falta la exención).
//
// Es una lista a mano: si una pantalla nativa se retira de KMP, hay que
// quitarla de aquí o el detector deja de verla.
var nativeScreenWhitelist = map[string]struct{}{
	"batch-enroll":        {}, // MainScreen.kt, NATIVE_ENROLLMENT_SCREEN_KEYS
	"join-requests-inbox": {}, // MainScreen.kt
	"materials-list":      {}, // MainScreen.kt + RouteRegistry.kt
	"messaging":           {}, // MainScreen.kt + RouteRegistry.kt
	"system-settings":     {}, // MainScreen.kt (junto a app-settings)
}

// isNativeScreen retorna true si el screenKey está en la lista blanca de
// pantallas nativas y debe excluirse de los drifts screen_key_dead.
func isNativeScreen(screenKey string) bool {
	_, ok := nativeScreenWhitelist[screenKey]
	return ok
}

// normalizeRoleCode compara códigos de rol sin guiones bajos ni
// mayúsculas: los dashboards se llaman `dashboard-schooladmin` y
// `dashboard-superadmin` y los roles del seed `school_admin` y
// `super_admin`.
func normalizeRoleCode(code string) string {
	return strings.ToLower(strings.ReplaceAll(code, "_", ""))
}

// genericDashboardSuffixes contiene sufijos de `dashboard-<X>` que no
// nombran un rol sino una pantalla común: `dashboard-home` es el landing
// por defecto de la cascada del backend (ADR 0024 F0). Criterio de
// admisión: el sufijo no es un rol en ningún seed y la pantalla existe
// como screen_instance.
var genericDashboardSuffixes = map[string]struct{}{
	"home": {},
}

func isGenericDashboardSuffix(code string) bool {
	_, ok := genericDashboardSuffixes[code]
	return ok
}
