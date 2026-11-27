package validate

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// Cada exención de la deuda 194 tiene aquí su caso rojo: el detector
// sigue mordiendo fuera de la exención.

func hasDrift(ds []Drift, cat, id string) bool {
	for _, d := range ds {
		if d.Category == cat && d.Identifier == id {
			return true
		}
	}
	return false
}

func TestDeadScreens_NativeWhitelistOnlyExemptsListed(t *testing.T) {
	s := seed.Snapshot{ResourceScreens: []seed.ResourceScreen{
		{ScreenKey: "messaging", ScreenType: "list", IsDefault: true},
		{ScreenKey: "ghost-native", ScreenType: "list", IsDefault: true},
	}}
	got := detectDeadScreenKeys(kmp.Snapshot{}, s)
	if hasDrift(got, CategoryScreenKeyDead, "messaging") {
		t.Fatalf("messaging está en la lista de nativas y no debía salir: %+v", got)
	}
	if !hasDrift(got, CategoryScreenKeyDead, "ghost-native") {
		t.Fatalf("una pantalla fuera de la lista debía seguir muerta: %+v", got)
	}
}

func TestPhantomScreens_ScreenInstanceCountsButNothingElse(t *testing.T) {
	k := kmp.Snapshot{ScreenKeys: map[string][]kmp.Location{
		"grades-batch": {{FilePath: "G.kt", Line: 1}},
		"ghost-list":   {{FilePath: "X.kt", Line: 1}},
	}}
	s := seed.Snapshot{ScreenInstances: []seed.ScreenInstance{{ScreenKey: "grades-batch"}}}
	got := detectPhantomScreenKeys(k, s)
	if hasDrift(got, CategoryScreenKeyPhantom, "grades-batch") {
		t.Fatalf("grades-batch tiene screen_instance y no debía salir: %+v", got)
	}
	if !hasDrift(got, CategoryScreenKeyPhantom, "ghost-list") {
		t.Fatalf("sin instancia ni mapping debía seguir fantasma: %+v", got)
	}
}

func TestRoles_DashboardSuffixNormalizedButUnknownStillPhantom(t *testing.T) {
	k := kmp.Snapshot{Roles: map[string][]kmp.Location{
		"superadmin": {{FilePath: "D.kt", Line: 1}},
		"home":       {{FilePath: "D.kt", Line: 2}},
		"homes":      {{FilePath: "D.kt", Line: 3}},
	}}
	s := seed.Snapshot{Roles: []seed.Role{
		{Code: "super_admin", Scope: "system"},
		{Code: "platform_admin", Scope: "system"},
	}}
	ph := detectPhantomRoles(k, s)
	if hasDrift(ph, CategoryRolePhantom, "superadmin") || hasDrift(ph, CategoryRolePhantom, "home") {
		t.Fatalf("superadmin (≈ super_admin) y home (genérico) no debían salir: %+v", ph)
	}
	if !hasDrift(ph, CategoryRolePhantom, "homes") {
		t.Fatalf("un sufijo desconocido debía seguir fantasma: %+v", ph)
	}
	un := detectUnusedRoles(k, s)
	if hasDrift(un, CategoryRoleUnused, "super_admin") {
		t.Fatalf("super_admin lo atiende dashboard-superadmin: %+v", un)
	}
	if !hasDrift(un, CategoryRoleUnused, "platform_admin") {
		t.Fatalf("un rol de sistema sin UI debía seguir sin usar: %+v", un)
	}
}

func TestServicePrefix_IdentityIsCanonicalAndIamIsMismatch(t *testing.T) {
	k := kmp.Snapshot{Contracts: []kmp.ContractDecl{
		{ScreenKey: "users-list", APIPrefix: "identity:", Resource: "users"},
		{ScreenKey: "audit-list", APIPrefix: "iam:", Resource: "audit"},
	}}
	got := detectServicePrefixMismatch(k)
	if hasDrift(got, CategoryServicePrefixMismatch, "users") {
		t.Fatalf("identity: es el prefijo canónico de users: %+v", got)
	}
	if !hasDrift(got, CategoryServicePrefixMismatch, "audit") {
		t.Fatalf("iam: ya no es canónico y debía salir: %+v", got)
	}
}
