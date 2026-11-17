package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/orchestrator"
	pgRunner "github.com/EduGoGroup/edugo-dev-environment/migrator/internal/postgres"
	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
)

func main() {
	fs := flag.NewFlagSet("migrator", flag.ExitOnError)
	seedUpToLayer := fs.String("seed-up-to-layer", "", "layer del seed system hasta la cual aplicar (vacío = todas)")
	playgroundV2Flag := fs.String(
		"playground-v2",
		"",
		fmt.Sprintf(
			"nombre del playground v2 a aplicar tras el sistema, \"all\" para aplicar todos, o uno de los registrados (%s). "+
				"Implica force migration. "+
				"Por defecto el sistema se siembra completo; los playgrounds v2 asumen L4 completo.",
			strings.Join(playground_v2.Available(), "|"),
		),
	)
	onTopFlag := fs.String(
		"on-top",
		"",
		fmt.Sprintf(
			"escenario on-top a aplicar SIN forzar sobre una BD ya migrada con base (%s). "+
				"Sus usuarios deben existir antes en auth.users (cmd/provision_users). "+
				"Excluyente con -playground-v2 y -seed-up-to-layer.",
			strings.Join(playground_v2.OnTopScenarios(), "|"),
		),
	)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("flag parse: %v", err)
	}

	cfg := config.Load()

	// Modo on-top (plan 053 T5.2): ramifica antes del switch de abajo, cuyo
	// default fuerza la migración. Nunca recrea ni migra: solo siembra el
	// escenario encima de lo que hay.
	if *onTopFlag != "" {
		if *playgroundV2Flag != "" || *seedUpToLayer != "" {
			log.Fatalf("❌ -on-top es excluyente con -playground-v2 y -seed-up-to-layer")
		}
		if cfg.ForceMigration {
			log.Fatalf("❌ -on-top no fuerza: quita FORCE_MIGRATION=true del entorno")
		}
		fmt.Println("=== EduGo Migrator (on-top) ===")
		fmt.Printf("📋 Schema version esperada: %s | destino: %s\n", postgresMigrations.SchemaVersion, cfg.Postgres.Target)
		if err := pgRunner.New(cfg.Postgres).ApplyOnTop(*onTopFlag); err != nil {
			log.Fatalf("❌ Error: %v\n", err)
		}
		return
	}
	cfg.SeedUpToLayer = *seedUpToLayer

	// Default (sin flags de seed): se siembra playground_v2/base, ya fijado por
	// config.Load(). Los flags explícitos de abajo lo sobrescriben. La selección
	// de fixture implica force migration (recreación desde cero), ya que los
	// playgrounds v2 asumen una BD limpia.
	//
	// Modo playground: NO forzamos un layer particular del sistema — el
	// playground puede convivir con cualquier capa, su funcionamiento depende
	// solo de L0 (siempre se siembra) y de sus propios grants para restringir el
	// menú. Si el caller necesita acotar el sistema, lo hace explícitamente con
	// -seed-up-to-layer.
	switch {
	case *playgroundV2Flag != "":
		cfg.PlaygroundV2 = *playgroundV2Flag
		cfg.ForceMigration = true
	case *seedUpToLayer != "":
		// Modo "seed hasta una capa": operación de scope del system seed, sin
		// fixture de datos encima. No inyecta base ni fuerza recreación; respeta
		// FORCE_MIGRATION del entorno si el caller lo pide.
		cfg.PlaygroundV2 = ""
	default:
		// Default base: heredado de config.Load(). Implica force migration
		// porque los playgrounds v2 asumen una BD limpia.
		cfg.ForceMigration = true
	}

	orch := orchestrator.New(cfg)
	if err := orch.Run(); err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}
}
