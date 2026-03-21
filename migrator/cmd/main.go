package main

import (
	"log"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/config"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/orchestrator"
)

func main() {
	cfg := config.Load()
	orch := orchestrator.New(cfg)
	if err := orch.Run(); err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}
}
