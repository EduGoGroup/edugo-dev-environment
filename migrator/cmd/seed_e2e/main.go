// Binario seed_e2e — aplica fixtures E2E del monorepo.
//
// Modo de invocación:
//
//   - SCENARIO=teacher_grades_only ./bin/seed_e2e
//     Aplica un scenario registrado por scenarios.RegisterAll().
//
//   - SCENARIO=teacher_grades_only,observer_audits ./bin/seed_e2e
//     Aplica varios scenarios secuencialmente, cada uno en su propia
//     transacción. Útil para tests que necesitan múltiples namespaces.
//
// Tras cada Apply exitoso el binario regenera el JSON de constantes
// (seeds/e2e/exports/fixtures-constants.json) que los tests Kotlin del
// KMP consumen (Fase C, C-REQ-6.3).
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/scenarios"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
)

// ExportPath es la ruta del JSON de constantes relativa al working
// directory del binario. La env var FIXTURES_EXPORT permite
// sobreescribirla en CI.
const defaultExportPath = "EduBack/edugo-infrastructure/postgres/seeds/e2e/exports/fixtures-constants.json"

func main() {
	scenarioEnv := strings.TrimSpace(os.Getenv("SCENARIO"))

	if scenarioEnv == "" {
		log.Fatal("❌ Definí SCENARIO=teacher_grades_only (o una lista separada por comas)")
	}

	cfg := config.Load()
	fmt.Println("=== EduGo Seed E2E ===")
	fmt.Printf("🎯 Scenario(s): %s\n", scenarioEnv)
	fmt.Println()

	db, err := sql.Open("postgres", cfg.Postgres.ConnStr)
	if err != nil {
		log.Fatalf("❌ Error abriendo conexión: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Error pingueando Postgres: %v", err)
	}
	fmt.Printf("✓ Conectado a PostgreSQL (user=%s)\n", cfg.Postgres.User)

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("❌ Error abriendo GORM: %v", err)
	}

	exporter := framework.NewConstantsExporter(exportPathFromEnv())

	// Modo SCENARIO=...: registramos los scenarios canónicos +
	// legacy una sola vez en un registry local y aplicamos cada
	// uno secuencialmente dentro de su propia transacción.
	reg := framework.NewRegistry()
	if err := scenarios.RegisterAll(reg); err != nil {
		log.Fatalf("❌ Error registrando scenarios: %v", err)
	}
	composer := framework.NewComposer(reg, framework.NewJSONLogger())
	for _, name := range splitAndTrim(scenarioEnv) {
		ctx, err := composer.Apply(gdb, name)
		if err != nil {
			log.Fatalf("❌ Error aplicando scenario %q: %v", name, err)
		}
		if err := exporter.WriteFromContext(ctx); err != nil {
			log.Fatalf("❌ Error exportando constantes (%s): %v", name, err)
		}
		fmt.Printf("✓ Scenario %q aplicado\n", name)
	}

	fmt.Println("\n✅ Fixtures E2E aplicadas correctamente")
}

// exportPathFromEnv permite redirigir el JSON via FIXTURES_EXPORT
// (path absoluto). Por default usa la ruta dentro del repo.
func exportPathFromEnv() string {
	if v := strings.TrimSpace(os.Getenv("FIXTURES_EXPORT")); v != "" {
		return v
	}
	if abs, err := filepath.Abs(defaultExportPath); err == nil {
		return abs
	}
	return defaultExportPath
}

func splitAndTrim(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}
