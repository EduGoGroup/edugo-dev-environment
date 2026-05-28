package validate

import (
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/kmp"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/contractcheck/seed"
)

func TestIsPreAuthScreen(t *testing.T) {
	cases := []struct {
		name      string
		screenKey string
		want      bool
	}{
		{name: "app-login_is_pre_auth", screenKey: "app-login", want: true},
		{name: "app-settings_is_not_pre_auth", screenKey: "app-settings", want: false},
		{name: "anything_is_not_pre_auth", screenKey: "anything", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isPreAuthScreen(tc.screenKey); got != tc.want {
				t.Fatalf("isPreAuthScreen(%q) = %v, want %v", tc.screenKey, got, tc.want)
			}
		})
	}
}

func TestIsStaticCompliantScreen(t *testing.T) {
	cases := []struct {
		name      string
		screenKey string
		want      bool
	}{
		{name: "app-settings_is_static_compliant", screenKey: "app-settings", want: true},
		{name: "system-settings_is_static_compliant", screenKey: "system-settings", want: true},
		{name: "app-login_is_not_static_compliant", screenKey: "app-login", want: false},
		{name: "anything_is_not_static_compliant", screenKey: "anything", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isStaticCompliantScreen(tc.screenKey); got != tc.want {
				t.Fatalf("isStaticCompliantScreen(%q) = %v, want %v", tc.screenKey, got, tc.want)
			}
		})
	}
}

func TestDetectPhantomScreenKeys_ExcludesPreAuth(t *testing.T) {
	loc := kmp.Location{FilePath: "F.kt", Line: 1, Snippet: "x"}
	k := kmp.Snapshot{
		ScreenKeys: map[string][]kmp.Location{
			"app-login": {loc},
			"app-foo":   {loc},
		},
	}
	s := seed.Snapshot{}

	got := detectPhantomScreenKeys(k, s)
	if len(got) != 1 {
		t.Fatalf("expected 1 phantom drift (app-foo only), got %d: %+v", len(got), got)
	}
	if got[0].Identifier != "app-foo" {
		t.Fatalf("expected drift for app-foo, got identifier %q", got[0].Identifier)
	}
	for _, d := range got {
		if d.Identifier == "app-login" {
			t.Fatalf("app-login should be excluded by pre-auth whitelist, but got drift: %+v", d)
		}
	}
}

func TestDetectPhantomScreenKeys_ExcludesStaticCompliant(t *testing.T) {
	loc := kmp.Location{FilePath: "F.kt", Line: 1, Snippet: "x"}
	k := kmp.Snapshot{
		ScreenKeys: map[string][]kmp.Location{
			"app-settings":    {loc},
			"system-settings": {loc},
			"app-foo":         {loc},
		},
	}
	s := seed.Snapshot{}

	got := detectPhantomScreenKeys(k, s)
	if len(got) != 1 {
		t.Fatalf("expected 1 phantom drift (app-foo only), got %d: %+v", len(got), got)
	}
	if got[0].Identifier != "app-foo" {
		t.Fatalf("expected drift for app-foo, got identifier %q", got[0].Identifier)
	}
	for _, d := range got {
		if d.Identifier == "app-settings" || d.Identifier == "system-settings" {
			t.Fatalf("%s should be excluded by static-compliant whitelist, but got drift: %+v", d.Identifier, d)
		}
	}
}
