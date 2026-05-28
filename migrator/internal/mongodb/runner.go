package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	mongoMigrations "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
)

// Runner ejecuta las migraciones de MongoDB.
type Runner struct {
	cfg config.MongoConfig
}

// New crea un nuevo Runner con la configuración dada.
func New(cfg config.MongoConfig) *Runner {
	return &Runner{cfg: cfg}
}

// Run ejecuta las migraciones de MongoDB.
func (r *Runner) Run(force bool, applyMock bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(r.cfg.URI))
	if err != nil {
		return fmt.Errorf("no se pudo conectar a MongoDB: %w", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("⚠️  Error desconectando cliente MongoDB: %v", err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("no se pudo hacer ping a MongoDB: %w", err)
	}

	fmt.Printf("✓ Conectado a MongoDB (db=%s)\n", r.cfg.DBName)

	db := client.Database(r.cfg.DBName)

	result, err := mongoMigrations.Migrate(ctx, db, mongoMigrations.MigrateOptions{
		Force:     force,
		ApplyMock: applyMock,
	})
	if err != nil {
		return fmt.Errorf("error en migraciones MongoDB: %w", err)
	}

	if result.Skipped {
		fmt.Println("✅ MongoDB actualizado - migraciones omitidas")
	} else {
		fmt.Println("✅ Migraciones de MongoDB completadas")
	}
	return nil
}
