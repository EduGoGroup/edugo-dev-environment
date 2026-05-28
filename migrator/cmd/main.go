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
	seedDemo := fs.Bool("seed-demo", true, "aplicar seed de demo (default true, sobrescribe APPLY_MOCK_DATA si explícito)")
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

	seedDemoExplicit := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed-demo" {
			seedDemoExplicit = true
		}
	})

	cfg := config.Load()
	cfg.SeedUpToLayer = *seedUpToLayer
	if seedDemoExplicit {
		cfg.SeedDemo = *seedDemo
	}

	// Modo playground: recrear desde cero y aplicar el/los playground(s)
	// pedido(s). NO forzamos un layer particular del sistema — el playground
	// puede convivir con cualquier capa, su funcionamiento depende solo de
	// L0 (siempre se siembra) y de sus propios grants para restringir el menú.
	// Si el caller necesita acotar el sistema, lo hace explícitamente con
	// -seed-up-to-layer.
	if *playgroundFlag != "" && *playgroundV2Flag != "" {
		log.Fatalf("flags -playground y -playground-v2 son mutuamente excluyentes")
	}
	if *playgroundFlag != "" {
		cfg.Playground = *playgroundFlag
		cfg.ForceMigration = true
		cfg.SeedDemo = false
	}
	if *playgroundV2Flag != "" {
		cfg.PlaygroundV2 = *playgroundV2Flag
		cfg.ForceMigration = true
		cfg.SeedDemo = false
	}

	orch := orchestrator.New(cfg)
	if err := orch.Run(); err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}
}
