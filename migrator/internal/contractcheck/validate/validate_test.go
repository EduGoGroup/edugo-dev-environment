package validate

import (
	"sort"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

// driftKey is the comparable shape used by the table-driven tests:
// (Category, Severity, Identifier). Detail and Evidence are validated
// separately when relevant.
type driftKey struct {
	Category   string
	Severity   Severity
	Identifier string
	Direction  Direction
}

func toKeys(drifts []Drift) []driftKey {
	out := make([]driftKey, 0, len(drifts))
	for _, d := range drifts {
		out = append(out, driftKey{d.Category, d.Severity, d.Identifier, d.Direction})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		if out[i].Severity != out[j].Severity {
			return out[i].Severity < out[j].Severity
		}
		return out[i].Identifier < out[j].Identifier
	})
	return out
}

func TestValidate_TableDriven(t *testing.T) {
	loc := func(file string, line int) kmp.Location {
		return kmp.Location{FilePath: file, Line: line, Snippet: "snippet"}
	}

	cases := []struct {
		name string
		k    kmp.Snapshot
		s    seed.Snapshot
		want []driftKey
	}{
		{
			name: "phantom_screen_keys",
			k: kmp.Snapshot{
				ScreenKeys: map[string][]kmp.Location{
					"ghost-form":   {loc("Ghost.kt", 10), loc("Ghost.kt", 10)}, // duplicado mismo file/line/snippet
					"schools-form": {loc("Schools.kt", 5)},
				},
			},
			s: seed.Snapshot{
				ResourceScreens: []seed.ResourceScreen{
					{ScreenKey: "schools-form", ScreenType: "form"},
				},
			},
			want: []driftKey{
				{CategoryScreenKeyPhantom, SeverityError, "ghost-form", DirectionFEOnly},
			},
		},
		{
			name: "dead_screen_keys_dashboard_escalates",
			k: kmp.Snapshot{
				ScreenKeys: map[string][]kmp.Location{
					"schools-form": {loc("Schools.kt", 5)},
				},
			},
			s: seed.Snapshot{
				ResourceScreens: []seed.ResourceScreen{
					{ScreenKey: "schools-form", ScreenType: "form"},
					{ScreenKey: "dashboard-orphan", ScreenType: "dashboard"},             // -> error
					{ScreenKey: "settings-default", ScreenType: "form", IsDefault: true}, // -> error
					{ScreenKey: "calendar-extra", ScreenType: "list"},                    // -> warning
				},
			},
			want: []driftKey{
				{CategoryScreenKeyDead, SeverityError, "dashboard-orphan", DirectionBEOnly},
				{CategoryScreenKeyDead, SeverityError, "settings-default", DirectionBEOnly},
				{CategoryScreenKeyDead, SeverityWarning, "calendar-extra", DirectionBEOnly},
			},
		},
		{
			name: "phantom_permissions_canonical_inference_and_explicit",
			k: kmp.Snapshot{
				Permissions: map[string][]kmp.Location{
					"assessments:publish": {loc("AssessmentsContract.kt", 12)},
				},
				Contracts: []kmp.ContractDecl{
					{
						ScreenKey: "ghost-form",
						APIPrefix: "academic:",
						BasePath:  "/api/v1/ghosts",
						Resource:  "ghosts", // resource no existe
						File:      loc("GhostContract.kt", 1),
					},
					{
						ScreenKey: "schools-form",
						APIPrefix: "academic:",
						BasePath:  "/api/v1/schools",
						Resource:  "schools",
						File:      loc("SchoolsContract.kt", 1),
					},
				},
			},
			s: seed.Snapshot{
				Resources: []seed.Resource{{Key: "schools"}},
				Permissions: []seed.Permission{
					{Code: "schools:read"},
					{Code: "schools:create"},
					// schools:update y schools:delete NO seedados → warning
				},
			},
			// canónicos generados:
			//   ghosts:{read,create,update,delete} → error (resource no existe)
			//   schools:{read,create,update,delete} → schools:read/create OK; schools:update/delete warning
			// explícito:
			//   assessments:publish → error (resource no existe en seed)
			want: []driftKey{
				{CategoryPermissionPhantom, SeverityError, "assessments:publish", DirectionFEOnly},
				{CategoryPermissionPhantom, SeverityError, "ghosts:create", DirectionFEOnly},
				{CategoryPermissionPhantom, SeverityError, "ghosts:delete", DirectionFEOnly},
				{CategoryPermissionPhantom, SeverityError, "ghosts:read", DirectionFEOnly},
				{CategoryPermissionPhantom, SeverityError, "ghosts:update", DirectionFEOnly},
				{CategoryPermissionPhantom, SeverityWarning, "schools:delete", DirectionFEOnly},
				{CategoryPermissionPhantom, SeverityWarning, "schools:update", DirectionFEOnly},
			},
		},
		{
			name: "zombie_permissions_info_vs_warning",
			k: kmp.Snapshot{
				// FE consume schools:read → no es zombie
				Permissions: map[string][]kmp.Location{
					"schools:read": {loc("SchoolsContract.kt", 5)},
				},
			},
			s: seed.Snapshot{
				Permissions: []seed.Permission{
					{Code: "schools:read"},
					{Code: "schools:archive"}, // sin role_perm, sin slot, sin FE → info
					{Code: "audit:read"},      // con role_perm pero sin FE/slot → warning
				},
				RolePermissions: []seed.RolePermission{
					{RoleCode: "auditor", PermissionCode: "audit:read"},
				},
			},
			want: []driftKey{
				{CategoryPermissionZombie, SeverityWarning, "audit:read", DirectionBEOnly},
				{CategoryPermissionZombie, SeverityInfo, "schools:archive", DirectionBEOnly},
			},
		},
		{
			name: "phantom_roles",
			k: kmp.Snapshot{
				Roles: map[string][]kmp.Location{
					"teacher":     {loc("DashboardSwitcher.kt", 3)},
					"phantom_one": {loc("DashboardSwitcher.kt", 4)},
				},
			},
			s: seed.Snapshot{
				Roles: []seed.Role{{Code: "teacher", Scope: "tenant"}},
			},
			want: []driftKey{
				{CategoryRolePhantom, SeverityError, "phantom_one", DirectionFEOnly},
			},
		},
		{
			name: "unused_roles_system_escalates",
			k: kmp.Snapshot{
				Roles: map[string][]kmp.Location{
					"teacher": {loc("DashboardSwitcher.kt", 3)},
				},
			},
			s: seed.Snapshot{
				Roles: []seed.Role{
					{Code: "teacher", Scope: "tenant"},
					{Code: "platform_admin", Scope: "system"}, // → error
					{Code: "tenant_lead", Scope: "tenant"},    // → warning
				},
			},
			want: []driftKey{
				{CategoryRoleUnused, SeverityError, "platform_admin", DirectionBEOnly},
				{CategoryRoleUnused, SeverityWarning, "tenant_lead", DirectionBEOnly},
			},
		},
		{
			name: "service_prefix_mismatch_and_unclassified",
			k: kmp.Snapshot{
				Contracts: []kmp.ContractDecl{
					{
						ScreenKey: "announcements-list",
						APIPrefix: "academic:", // bug F2·H3.a: debería ser platform:
						Resource:  "announcements",
						File:      loc("AnnouncementsContract.kt", 1),
					},
					{
						ScreenKey: "schools-form",
						APIPrefix: "academic:",
						Resource:  "schools",
						File:      loc("SchoolsContract.kt", 1),
					},
					{
						ScreenKey: "frobs-list",
						APIPrefix: "academic:",
						Resource:  "frobs", // no clasificado → info
						File:      loc("FrobsContract.kt", 1),
					},
				},
			},
			s: seed.Snapshot{},
			want: []driftKey{
				{CategoryServicePrefixMismatch, SeverityError, "announcements", DirectionMismatch},
				{CategoryServicePrefixMismatch, SeverityInfo, "frobs", DirectionMismatch},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Validate(tc.k, tc.s)
			// Filter: keep only the categories asserted by the case to
			// make the assertions readable. The orchestrator emits all
			// detectors, but each fixture targets one category at a time.
			expectedCats := make(map[string]struct{})
			for _, w := range tc.want {
				expectedCats[w.Category] = struct{}{}
			}
			filtered := make([]Drift, 0, len(got))
			for _, d := range got {
				if _, ok := expectedCats[d.Category]; ok {
					filtered = append(filtered, d)
				}
			}
			gotKeys := toKeys(filtered)
			wantKeys := append([]driftKey(nil), tc.want...)
			sort.Slice(wantKeys, func(i, j int) bool {
				if wantKeys[i].Category != wantKeys[j].Category {
					return wantKeys[i].Category < wantKeys[j].Category
				}
				if wantKeys[i].Severity != wantKeys[j].Severity {
					return wantKeys[i].Severity < wantKeys[j].Severity
				}
				return wantKeys[i].Identifier < wantKeys[j].Identifier
			})
			if !equalKeys(gotKeys, wantKeys) {
				t.Fatalf("\nGot drifts: %+v\nWant drifts: %+v", gotKeys, wantKeys)
			}
		})
	}
}

func equalKeys(a, b []driftKey) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestValidate_OrderingIsDeterministic(t *testing.T) {
	loc := kmp.Location{FilePath: "F.kt", Line: 1, Snippet: "x"}
	k := kmp.Snapshot{
		ScreenKeys: map[string][]kmp.Location{
			"zzz-phantom": {loc},
			"aaa-phantom": {loc},
		},
	}
	s := seed.Snapshot{}

	first := Validate(k, s)
	second := Validate(k, s)
	if len(first) != len(second) {
		t.Fatalf("non-deterministic length: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].Category != second[i].Category ||
			first[i].Severity != second[i].Severity ||
			first[i].Identifier != second[i].Identifier ||
			first[i].Direction != second[i].Direction {
			t.Fatalf("ordering not stable at %d: %+v vs %+v", i, first[i], second[i])
		}
	}
	// Identifiers should be alphabetically sorted within the same
	// (category, severity) bucket.
	if first[0].Identifier != "aaa-phantom" || first[1].Identifier != "zzz-phantom" {
		t.Fatalf("expected alphabetical ordering, got %v", first)
	}
}

func TestValidate_DedupesEvidenceLocations(t *testing.T) {
	dup := kmp.Location{FilePath: "Ghost.kt", Line: 10, Snippet: "x"}
	other := kmp.Location{FilePath: "Ghost2.kt", Line: 1, Snippet: "y"}
	k := kmp.Snapshot{
		ScreenKeys: map[string][]kmp.Location{
			"ghost": {dup, dup, other, dup},
		},
	}
	got := detectPhantomScreenKeys(k, seed.Snapshot{})
	if len(got) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(got))
	}
	if len(got[0].Evidence) != 2 {
		t.Fatalf("expected 2 deduped locations, got %d: %+v", len(got[0].Evidence), got[0].Evidence)
	}
}

func TestValidate_PermissionPhantomNonCanonicalFormat(t *testing.T) {
	k := kmp.Snapshot{
		Permissions: map[string][]kmp.Location{
			"weirdformat": {{FilePath: "F.kt", Line: 1}}, // no contiene ":"
		},
	}
	got := detectPhantomPermissions(k, seed.Snapshot{})
	if len(got) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(got))
	}
	if got[0].Severity != SeverityError {
		t.Fatalf("expected error severity, got %s", got[0].Severity)
	}
}

func TestValidate_ZombieIgnoresSlotReferenced(t *testing.T) {
	s := seed.Snapshot{
		Permissions: []seed.Permission{
			{Code: "audit:read"},
		},
		ScreenInstances: []seed.ScreenInstance{
			{ScreenKey: "x", SlotData: []byte(`{"requiredPermission": "audit:read"}`)},
		},
	}
	got := detectZombiePermissions(kmp.Snapshot{}, s)
	if len(got) != 0 {
		t.Fatalf("permiso referenciado en slot no debería ser zombie, got %+v", got)
	}
}

func TestSeverityFor_FallbackForUnknownCategory(t *testing.T) {
	if got := SeverityFor("unknown_category"); got != SeverityWarning {
		t.Fatalf("expected SeverityWarning fallback, got %s", got)
	}
	if got := SeverityFor(CategoryScreenKeyPhantom); got != SeverityError {
		t.Fatalf("expected SeverityError for phantom, got %s", got)
	}
}

func TestServicePrefixTable_NormalizesMissingColon(t *testing.T) {
	k := kmp.Snapshot{
		Contracts: []kmp.ContractDecl{
			{
				ScreenKey: "schools-form",
				APIPrefix: "academic", // sin trailing ":" — debería normalizarse
				Resource:  "schools",
				File:      kmp.Location{FilePath: "F.kt", Line: 1},
			},
		},
	}
	got := detectServicePrefixMismatch(k)
	if len(got) != 0 {
		t.Fatalf("apiPrefix sin ':' debería normalizarse a 'academic:' y no emitir drift, got %+v", got)
	}
}
