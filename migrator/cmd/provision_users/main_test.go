package main

import (
	"os"
	"strings"
	"testing"
)

// Las validaciones de entrada cortan antes de tocar la red o la BD.
func TestRun_RejectsBadInputBeforeAnyCall(t *testing.T) {
	cases := []struct {
		name, core, password string
		args                 []string
		want                 string
	}{
		{"sin escenario", "http://core.test", "una-clave-de-prueba-larga", nil, "Falta -scenario"},
		{"escenario no on-top", "http://core.test", "una-clave-de-prueba-larga", []string{"-scenario=base"}, "no es un escenario on-top"},
		{"sin identity-core", "", "una-clave-de-prueba-larga", []string{"-scenario=onboarding"}, "Falta IDENTITY_CORE_BASE_URL"},
		{"sin contraseña", "http://core.test", "", []string{"-scenario=onboarding"}, "Falta SCENARIO_TEST_PASSWORD"},
		{"contraseña corta", "http://core.test", "corta-11ch!", []string{"-scenario=onboarding"}, "al menos 12"},
		{"url sin esquema", "core.test", "una-clave-de-prueba-larga", []string{"-scenario=onboarding"}, "IDENTITY_CORE_BASE_URL inválida"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("IDENTITY_CORE_BASE_URL", c.core)
			t.Setenv("SCENARIO_TEST_PASSWORD", c.password)
			var out strings.Builder
			if rc := run(c.args, os.Stdin, &out); rc != 2 {
				t.Fatalf("rc=%d, esperado 2; salida:\n%s", rc, out.String())
			}
			if !strings.Contains(out.String(), c.want) {
				t.Errorf("salida sin %q:\n%s", c.want, out.String())
			}
			if c.password != "" && strings.Contains(out.String(), c.password) {
				t.Errorf("la contraseña aparece en la salida")
			}
		})
	}
}

// Sin -yes y sin terminal interactiva, se niega sin escribir nada.
func TestRun_RefusesWithoutTTYOrYes(t *testing.T) {
	t.Setenv("IDENTITY_CORE_BASE_URL", "http://core.test")
	t.Setenv("SCENARIO_TEST_PASSWORD", "una-clave-de-prueba-larga")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	defer func() { _ = w.Close() }()
	var out strings.Builder
	if rc := run([]string{"-scenario=onboarding"}, r, &out); rc != 3 {
		t.Fatalf("rc=%d, esperado 3; salida:\n%s", rc, out.String())
	}
	if !strings.Contains(out.String(), "core.test") || !strings.Contains(out.String(), "-yes") {
		t.Errorf("debía mostrar los destinos y pedir -yes:\n%s", out.String())
	}
}
