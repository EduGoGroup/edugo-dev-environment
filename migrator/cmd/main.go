package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/orchestrator"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
)

func main() {
	fs := flag.NewFlagSet("migrator", flag.ExitOnError)
	seedUpToLayer := fs.String("seed-up-to-layer", "", "layer del seed system hasta la cual aplicar (vacío = todas)")
	seedDemo := fs.Bool("seed-demo", false, "aplicar el seed de demo legacy en vez del default (playground_v2/base). Implica force migration y omite base.")
	playgroundFlag := fs.String(
		"playground",
		"",
		fmt.Sprintf(
			"nombre del playground a aplicar tras el sistema, \"all\" para aplicar todos, o uno de los registrados (%s). "+
				"Implica force migration y skip demo. "+
				"Por defecto el sistema se siembra completo (todas las capas); si querés acotar el sistema usá -seed-up-to-layer.",
			strings.Join(playground.Available(), "|"),
		),
	)
	playgroundV2Flag := fs.String(
		"playground-v2",
		"",
		fmt.Sprintf(
			"nombre del playground v2 a aplicar tras el sistema, \"all\" para aplicar todos, o uno de los registrados (%s). "+
				"Mutuamente excluyente con -playground. Implica force migration y skip demo. "+
				"Por defecto el sistema se siembra completo; los playgrounds v2 asumen L4 completo.",
			strings.Join(playground_v2.Available(), "|"),
		),
	)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("flag parse: %v", err)
	}

	cfg := config.Load()
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
	if *playgroundFlag != "" && *playgroundV2Flag != "" {
		log.Fatalf("flags -playground y -playground-v2 son mutuamente excluyentes")
	}
	switch {
	case *seedDemo:
		// Seed de demo legacy explícito: reemplaza el default base.
		cfg.SeedDemo = true
		cfg.PlaygroundV2 = ""
		cfg.ForceMigration = true
	case *playgroundFlag != "":
		cfg.Playground = *playgroundFlag
		cfg.PlaygroundV2 = ""
		cfg.ForceMigration = true
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
