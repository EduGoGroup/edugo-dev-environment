package config

import (
	"strings"
	"testing"
)

func TestTargetNeverCarriesThePassword(t *testing.T) {
	cases := map[string]string{
		"postgres://u:secreto-uri@db.test:5432/edugo?sslmode=require":     "db.test:5432/edugo",
		"host=db.test port=5432 user=u password=secreto-dsn dbname=edugo": "(POSTGRES_URI ilegible)",
	}
	for uri, want := range cases {
		t.Setenv("POSTGRES_URI", uri)
		got := Load().Postgres.Target
		if got != want {
			t.Errorf("Target(%q) = %q, esperado %q", uri, got, want)
		}
		if strings.Contains(got, "secreto") {
			t.Errorf("Target filtra la contraseña: %q", got)
		}
	}

	t.Setenv("POSTGRES_URI", "")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_DB", "edugo")
	t.Setenv("POSTGRES_PASSWORD", "secreto-env")
	if got := Load().Postgres.Target; got != "localhost:5433/edugo" {
		t.Errorf("Target = %q", got)
	}
}
