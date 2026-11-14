# EduGo Dev Environment - Makefile raíz
# Toda la lógica vive en migrator/. Este Makefile es solo un atajo.

.PHONY: help migrator vuln ci-local ci-docker cross-api-check cross-api-test

MIGRATOR_DIR := migrator
TEST_DIR := test
# Todos los paquetes de test/ llevan el build tag `integration`: sin el tag,
# `./...` no casa ningún paquete y los gates pasan en vacío.
CROSS_API_TAGS := integration,notest_swagger

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

vuln: ## govulncheck sobre migrator/ (delega a migrator)
	@$(MAKE) -C $(MIGRATOR_DIR) vuln

ci-local: ## CI local (delega a migrator)
	@$(MAKE) -C $(MIGRATOR_DIR) ci-local

ci-docker: ## CI en Docker (delega a migrator)
	@$(MAKE) -C $(MIGRATOR_DIR) ci-docker

cross-api-check: ## Gate local de test/: gofmt + vet + build + govulncheck (sin Docker)
	@cd $(TEST_DIR) && out="$$(gofmt -l .)" && if [ -n "$$out" ]; then echo "gofmt: archivos sin formatear en test/:"; echo "$$out"; exit 1; fi
	@cd $(TEST_DIR) && GOWORK=off go vet -tags=$(CROSS_API_TAGS) ./...
	@cd $(TEST_DIR) && GOWORK=off go build -tags=$(CROSS_API_TAGS) ./...
	@cd $(TEST_DIR) && GOWORK=off govulncheck -test -tags=$(CROSS_API_TAGS) ./...

cross-api-test: ## Tests cross-API de test/ (requiere Docker y los repos edugo-api-* al lado)
	@cd $(TEST_DIR) && GOWORK=off ENABLE_INTEGRATION_TESTS=true \
		go test -tags=$(CROSS_API_TAGS) -count=1 -timeout=1200s -p 1 ./integration/...
