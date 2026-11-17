// Package ontop aplica un escenario on-top de playground_v2 (plan 053 T5.2)
// sobre una BD ya migrada, sin forzar ni recrear nada.
//
// Es el paso 4 del flujo de escenarios on-top: la BD se migró con base, los
// usuarios del escenario se dieron de alta por la API de identidad
// (cmd/provision_users) y aquí se siembra el resto del escenario.
package ontop

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Apply comprueba que name sea un escenario on-top y que la BD esté migrada a
// la SchemaVersion y al content hash actuales, y llama a
// playground_v2.ApplyOnTop sobre la misma conexión.
func Apply(db *sql.DB, name string) error {
	if !slices.Contains(playground_v2.OnTopScenarios(), name) {
		return fmt.Errorf("%q no es un escenario on-top (on-top: %v)", name, playground_v2.OnTopScenarios())
	}

	status, err := migrations.Status(db)
	if err != nil {
		return fmt.Errorf("la BD no está migrada: migra antes con base (make docker-recreate): %w", err)
	}
	if status.NeedsForce {
		return fmt.Errorf("la BD está en la versión %s y el migrator espera la %s (o el content hash no coincide): "+
			"recréala con base (make docker-recreate) antes de aplicar un escenario on-top",
			status.Version, migrations.SchemaVersion)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("abriendo gorm sobre la conexión: %w", err)
	}
	return playground_v2.ApplyOnTop(gdb, name)
}
