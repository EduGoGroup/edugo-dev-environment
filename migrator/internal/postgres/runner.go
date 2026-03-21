package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
)

// Runner ejecuta las migraciones de PostgreSQL.
type Runner struct {
	cfg config.PostgresConfig
}

// New crea un nuevo Runner con la configuración dada.
func New(cfg config.PostgresConfig) *Runner {
	return &Runner{cfg: cfg}
}

// Status imprime el estado actual del schema sin modificar nada.
func (r *Runner) Status() error {
	db, err := r.connect()
	if err != nil {
		return err
	}
	defer r.closeDB(db)

	result, err := postgresMigrations.Status(db)
	if err != nil {
		fmt.Printf("⚠️  No se pudo leer version de BD: %v\n", err)
		return nil
	}

	fmt.Println("\n📊 Version actual en BD:")
	fmt.Printf("   Version:       %s\n", result.Version)
	fmt.Printf("   Content hash:  %s\n", result.ContentHash)
	fmt.Printf("   Execution ID:  %s\n", result.ExecutionID)

	if !result.NeedsForce {
		fmt.Println("\n✅ BD está actualizada - version y content hash coinciden")
	} else {
		fmt.Printf("\n⚠️  BD DESACTUALIZADA — esperado: v%s\n", postgresMigrations.SchemaVersion)
		fmt.Println("   Ejecuta FORCE_MIGRATION=true para actualizar")
	}
	return nil
}

// Run ejecuta las migraciones de PostgreSQL.
func (r *Runner) Run(force bool, applyMock bool) error {
	db, err := r.connect()
	if err != nil {
		return err
	}
	defer r.closeDB(db)

	fmt.Printf("✓ Conectado a PostgreSQL (user=%s)\n", r.cfg.User)

	result, err := postgresMigrations.Migrate(db, postgresMigrations.MigrateOptions{
		Force:     force,
		ApplyMock: applyMock,
		DBUser:    r.cfg.User,
	})
	if err != nil {
		return fmt.Errorf("error en migraciones PostgreSQL: %w", err)
	}

	if result.Skipped {
		if result.NeedsForce {
			fmt.Println("⚠️  PostgreSQL DESACTUALIZADO — usa FORCE_MIGRATION=true para recrear")
		} else {
			fmt.Println("✅ PostgreSQL actualizado - migraciones omitidas")
		}
	} else {
		fmt.Printf("✅ PostgreSQL migrado → Version: %s | Hash: %s | Execution: %s\n",
			result.Version, result.ContentHash, result.ExecutionID)
	}
	return nil
}

func (r *Runner) connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", r.cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a PostgreSQL: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("no se pudo hacer ping a PostgreSQL: %w", err)
	}
	return db, nil
}

func (r *Runner) closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("⚠️  Error cerrando conexión PostgreSQL: %v", err)
	}
}
