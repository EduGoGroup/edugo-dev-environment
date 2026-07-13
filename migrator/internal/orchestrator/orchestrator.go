package orchestrator

import (
	"fmt"

	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
	pgRunner "github.com/EduGoGroup/edugo-dev-environment/migrator/internal/postgres"
)

// Orchestrator coordina la ejecución de las migraciones de PostgreSQL.
// MongoDB fue retirado del ecosistema (plan 037, D-037.11).
type Orchestrator struct {
	cfg config.Config
	pg  *pgRunner.Runner
}

// New crea un nuevo Orchestrator con la configuración dada.
func New(cfg config.Config) *Orchestrator {
	return &Orchestrator{
		cfg: cfg,
		pg:  pgRunner.New(cfg.Postgres),
	}
}

// Run ejecuta el flujo completo de migraciones según la configuración.
func (o *Orchestrator) Run() error {
	fmt.Println("=== EduGo Migrator ===")
	fmt.Printf("📋 Schema version: %s\n", postgresMigrations.SchemaVersion)

	if o.cfg.ForceMigration {
		fmt.Println("⚠️  MODO FORZADO ACTIVADO")
	}
	fmt.Println()

	if o.cfg.StatusOnly {
		if err := o.pg.Status(); err != nil {
			return fmt.Errorf("error leyendo estado PostgreSQL: %w", err)
		}
		return nil
	}

	fmt.Println("\n--- PostgreSQL ---")
	if err := o.pg.Run(o.cfg.ForceMigration, o.cfg.SeedUpToLayer, o.cfg.PlaygroundV2); err != nil {
		return fmt.Errorf("error en PostgreSQL: %w", err)
	}

	fmt.Println("\n✅ Migraciones completadas")
	return nil
}
