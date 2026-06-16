# EduGo Dev Environment - Makefile raíz
# Toda la lógica vive en migrator/. Este Makefile es solo un atajo.

.PHONY: help migrator ci-local ci-docker

MIGRATOR_DIR := migrator

help: ## Mostrar esta ayuda
	@echo ""
	@echo "edugo-dev-environment"
	@echo "====================="
	@echo ""
	@echo "Toda la actividad sucede dentro de migrator/. Ejemplo:"
	@echo "  cd migrator && make help"
	@echo ""
	@echo "Atajo desde la raíz: pasa el target con make migrator T=<target>"
	@echo "  make migrator T=cloud-migrate"
	@echo "  make migrator T=cloud-status"
	@echo "  make migrator T=check"
	@echo ""

migrator: ## Delegar al Makefile de migrator. Uso: make migrator T=<target>
	@if [ -z "$(T)" ]; then \
		echo "Uso: make migrator T=<target> (ej. T=cloud-migrate)"; \
		exit 1; \
	fi
	@$(MAKE) -C $(MIGRATOR_DIR) $(T)

ci-local: ## CI local (delega a migrator)
	@$(MAKE) -C $(MIGRATOR_DIR) ci-local

ci-docker: ## CI en Docker (delega a migrator)
	@$(MAKE) -C $(MIGRATOR_DIR) ci-docker
