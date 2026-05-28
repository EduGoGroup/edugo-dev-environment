package validate

// preAuthScreenWhitelist contiene screenKeys que el KMP declara como
// pantallas pre-autenticación: existen visualmente antes de que el
// usuario tenga JWT válido, por lo tanto NO tienen contrato de
// permisos y NO deben sembrarse en resource_screens.
//
// Estas screens están exentas del drift screen_key_phantom.
//
// Criterio de admisión: solo pantallas previas a JWT válido. Cualquier
// otra pantalla huérfana es bug.
//
// Ver: e2e-integration-plan/seed-rebuild-spec/phase-7-static-screens-audit/
var preAuthScreenWhitelist = map[string]struct{}{
	"app-login": {},
}

// staticCompliantScreenWhitelist contiene screenKeys de pantallas
// estáticas post-auth que cumplen el contrato de permisos completo
// pero cuya UI no se resuelve dinámicamente y por lo tanto NO se
// siembran en resource_screens (decisión B5 del seed-rebuild-spec:
// el FE las resuelve por screen_key directo).
//
// Criterio de admisión:
//   - Tiene ScreenContract declarado.
//   - resource no vacío.
//   - permissionFor() no nulo para todos los ScreenEvent relevantes.
//   - Cada CustomEventHandler que muta estado declara
//     requiredPermission no nulo.
//   - Los permisos consumidos existen sembrados en iam.permissions.
//
// Estas screens están exentas del drift screen_key_phantom por la
// limitación del cross-checker (mira solo resource_screens), no por
// ser huérfanas en sentido de contrato. La validación de contrato
// se hace por revisión de código + presencia de permission en seed.
//
// Ver: docs/static-screens-contract.md en EduUI/edugo-ui-kmp
var staticCompliantScreenWhitelist = map[string]struct{}{
	"app-settings":    {},
	"system-settings": {},
}

// isPreAuthScreen retorna true si el screenKey está en la lista blanca
// pre-auth y debe excluirse de los drifts de tipo screen_key_phantom.
func isPreAuthScreen(screenKey string) bool {
	_, ok := preAuthScreenWhitelist[screenKey]
	return ok
}

// isStaticCompliantScreen retorna true si el screenKey está en la lista
// blanca de pantallas estáticas post-auth con contrato completo, y
// debe excluirse de los drifts de tipo screen_key_phantom.
func isStaticCompliantScreen(screenKey string) bool {
	_, ok := staticCompliantScreenWhitelist[screenKey]
	return ok
}
