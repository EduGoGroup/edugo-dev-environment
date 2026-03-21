package orchestrator

import (
	"fmt"

	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
	mongoRunner "github.com/EduGoGroup/edugo-dev-environment/migrator/internal/mongodb"
	pgRunner "github.com/EduGoGroup/edugo-dev-environment/migrator/internal/postgres"
)

// Orchestrator coordina la ejecución de las migraciones de PostgreSQL y MongoDB.
type Orchestrator struct {
	cfg   config.Config
	pg    *pgRunner.Runner
	mongo *mongoRunner.Runner
}

// New crea un nuevo Orchestrator con la configuración dada.
func New(cfg config.Config) *Orchestrator {
	return &Orchestrator{
		cfg:   cfg,
		pg:    pgRunner.New(cfg.Postgres),
		mongo: mongoRunner.New(cfg.Mongo),
	}
}

// Run ejecuta el flujo completo de migraciones según la configuración.
func (o *Orchestrator) Run() error {
	fmt.Println("=== EduGo Migrator ===")
	fmt.Printf("📋 Schema version: %s\n", postgresMigrations.SchemaVersion)

	if o.cfg.ForceMigration {
		fmt.Println("⚠️  MODO FORZADO ACTIVADO")
	}
	if !o.cfg.ApplyMockData {
		fmt.Println("ℹ️  Mock data deshabilitado")
	}
	fmt.Println()

	if o.cfg.StatusOnly {
		if o.cfg.MongoOnly {
			return fmt.Errorf("STATUS_ONLY para MongoDB no esta soportado: no hay implementacion de estado para Mongo")
		}
		if !o.cfg.MongoOnly {
			if err := o.pg.Status(); err != nil {
				return fmt.Errorf("error leyendo estado PostgreSQL: %w", err)
			}
		}
		return nil
	}

	if !o.cfg.MongoOnly {
		fmt.Println("\n--- PostgreSQL ---")
		if err := o.pg.Run(o.cfg.ForceMigration, o.cfg.ApplyMockData); err != nil {
			return fmt.Errorf("error en PostgreSQL: %w", err)
		}
	}

	if !o.cfg.PostgresOnly {
		fmt.Println("\n--- MongoDB ---")
		if err := o.mongo.Run(o.cfg.ForceMigration, o.cfg.ApplyMockData); err != nil {
			return fmt.Errorf("error en MongoDB: %w", err)
		}
	}

	fmt.Println("\n✅ Migraciones completadas")
	return nil
}
