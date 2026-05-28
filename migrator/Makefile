# EduGo Migrator
# Makefile independiente — ejecutar desde migrator/

.PHONY: help build clean fmt vet lint tidy test test-postgres test-mongo run status docker-up docker-down docker-migrate docker-recreate cloud-migrate cloud-recreate cloud-status cloud-seed-layer docker-seed-layer cloud-seed-scenario docker-playground docker-playground-list docker-playground-v2 docker-playground-v2-list cloud-playground cloud-playground-v2 seed-audit seed-audit-strict seed-audit-test contract-check contract-check-strict contract-check-update-baseline contract-check-test check

# Configuracion
BINARY  := migrator
BIN_DIR := bin
CMD_PKG := ./cmd
TEST_PKG := ./tests
GO      := go
LDFLAGS := -ldflags="-s -w"

# Colores
BLUE   := \033[0;34m
GREEN  := \033[0;32m
YELLOW := \033[1;33m
NC     := \033[0m

# ── Ayuda ───────────────────────────────────────────────────────────

help: ## Mostrar esta ayuda
	@printf "\n$(BLUE)EduGo Migrator$(NC)\n"
	@printf "==============\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-14s$(NC) %s\n", $$1, $$2}'
	@printf "\n"

# ── Build ────────────────────────────────────────────────────────────

build: ## Compilar binario en bin/
	@printf "$(BLUE)Compilando...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY) $(CMD_PKG)
	@printf "$(GREEN)OK $(BIN_DIR)/$(BINARY)$(NC)\n"

clean: ## Eliminar binarios y artefactos compilados
	@printf "$(BLUE)Limpiando...$(NC)\n"
	@rm -rf $(BIN_DIR)
	@rm -f *.test *.out
	@printf "$(GREEN)OK Limpio$(NC)\n"

# ── Calidad de codigo ────────────────────────────────────────────────

fmt: ## Formatear codigo con go fmt
	@printf "$(BLUE)Formateando...$(NC)\n"
	@$(GO) fmt ./...
	@printf "$(GREEN)OK Formateado$(NC)\n"

vet: ## Analizar codigo con go vet
	@printf "$(BLUE)Analizando (vet)...$(NC)\n"
	@$(GO) vet ./...
	@printf "$(GREEN)OK Vet$(NC)\n"

lint: ## Ejecutar golangci-lint (fallback: go vet)
	@printf "$(BLUE)Lint...$(NC)\n"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		printf "$(YELLOW)golangci-lint no instalado, usando go vet como fallback$(NC)\n"; \
		$(GO) vet ./...; \
	fi
	@printf "$(GREEN)OK Lint$(NC)\n"

tidy: ## Sincronizar go.mod y go.sum
	@printf "$(BLUE)go mod tidy...$(NC)\n"
	@$(GO) mod tidy
	@printf "$(GREEN)OK go.mod actualizado$(NC)\n"

# ── Tests ────────────────────────────────────────────────────────────

test: ## Ejecutar todos los tests de integracion (requiere Docker)
	@printf "$(BLUE)Tests...$(NC)\n"
	@$(GO) test -v -count=1 -timeout 120s $(TEST_PKG)

test-postgres: ## Ejecutar solo tests de PostgreSQL
	@printf "$(BLUE)Tests PostgreSQL...$(NC)\n"
	@$(GO) test -v -count=1 -timeout 120s $(TEST_PKG) -run TestPostgresIntegration

test-mongo: ## Ejecutar solo tests de MongoDB
	@printf "$(BLUE)Tests MongoDB...$(NC)\n"
	@$(GO) test -v -count=1 -timeout 120s $(TEST_PKG) -run TestMongoDBIntegration

# ── Ejecucion ────────────────────────────────────────────────────────

run: build ## Compilar y ejecutar localmente
	@printf "$(BLUE)Ejecutando migrator...$(NC)\n"
	@./$(BIN_DIR)/$(BINARY)

status: build ## Ver estado de la BD sin aplicar cambios
	@STATUS_ONLY=true ./$(BIN_DIR)/$(BINARY)

# ── Docker Local ─────────────────────────────────────────────────────

docker-up: ## Levantar PostgreSQL + MongoDB locales para testing
	@printf "$(BLUE)Levantando bases de datos locales...$(NC)\n"
	@docker compose -f docker-compose.local.yml up -d
	@printf "$(GREEN)OK Bases de datos listas$(NC)\n"

docker-down: ## Detener y eliminar contenedores locales
	@printf "$(BLUE)Deteniendo bases de datos locales...$(NC)\n"
	@docker compose -f docker-compose.local.yml down -v
	@printf "$(GREEN)OK Contenedores eliminados$(NC)\n"

docker-migrate: build ## Levantar DBs y aplicar migraciones (flujo completo local)
	@printf "$(BLUE)Levantando bases de datos y esperando healthchecks...$(NC)\n"
	@docker compose -f docker-compose.local.yml up -d --wait
	@printf "$(BLUE)Aplicando migraciones en BD local...$(NC)\n"
	@POSTGRES_PORT=5433 MONGO_PORT=27018 ./$(BIN_DIR)/$(BINARY)
	@printf "$(GREEN)OK Migraciones aplicadas$(NC)\n"

docker-recreate: build ## Recrear BDs desde cero con datos reales (BORRA DATOS)
	@printf "$(BLUE)Levantando bases de datos y esperando healthchecks...$(NC)\n"
	@docker compose -f docker-compose.local.yml up -d --wait
	@printf "$(YELLOW)Recreando bases de datos...$(NC)\n"
	@POSTGRES_PORT=5433 MONGO_PORT=27018 FORCE_MIGRATION=true ./$(BIN_DIR)/$(BINARY)
	@printf "$(GREEN)OK Bases de datos recreadas$(NC)\n"

# ── Cloud (Neon + Atlas) ─────────────────────────────────────────────

cloud-migrate: build ## Aplicar migraciones en Neon + Atlas (usa .env.cloud)
	@printf "$(BLUE)Aplicando migraciones en cloud (Neon + Atlas)...$(NC)\n"
	@set -a && . ./.env.cloud && set +a && ./$(BIN_DIR)/$(BINARY)
	@printf "$(GREEN)OK Migraciones aplicadas$(NC)\n"

cloud-recreate: build ## Recrear BDs desde cero en Neon + Atlas (BORRA DATOS)
	@printf "$(YELLOW)Recreando bases de datos en cloud...$(NC)\n"
	@set -a && . ./.env.cloud && set +a && FORCE_MIGRATION=true ./$(BIN_DIR)/$(BINARY)
	@printf "$(GREEN)OK Bases de datos recreadas$(NC)\n"

cloud-status: build ## Ver estado de la BD cloud sin aplicar cambios
	@printf "$(BLUE)Estado de BD cloud...$(NC)\n"
	@set -a && . ./.env.cloud && set +a && STATUS_ONLY=true ./$(BIN_DIR)/$(BINARY)

cloud-seed-layer: build ## Aplicar seeds hasta una capa específica en cloud (LAYER=legacy|l0|l1|...)
	@if [ -z "$(LAYER)" ]; then echo "Uso: make cloud-seed-layer LAYER=<nombre>"; exit 1; fi
	@set -a && . ./.env.cloud && set +a && ./$(BIN_DIR)/$(BINARY) --seed-up-to-layer=$(LAYER)

docker-seed-layer: build ## Aplicar seeds hasta una capa específica en docker local (LAYER=legacy|l0|l1|...)
	@if [ -z "$(LAYER)" ]; then echo "Uso: make docker-seed-layer LAYER=<nombre>"; exit 1; fi
	@POSTGRES_PORT=5433 MONGO_PORT=27018 ./$(BIN_DIR)/$(BINARY) --seed-up-to-layer=$(LAYER)

# ── Playgrounds (datasets pequeños y focalizados) ────────────────────

docker-playground: build ## Recrear BD docker con playground(s). P=all (default) | P=<name>
	@PLAY="$${P:-all}"; \
	printf "$(BLUE)Levantando bases de datos y esperando healthchecks...$(NC)\n"; \
	docker compose -f docker-compose.local.yml up -d --wait; \
	printf "$(YELLOW)Recreando BD con playground=$$PLAY (sistema completo + playground)...$(NC)\n"; \
	POSTGRES_PORT=5433 MONGO_PORT=27018 ./$(BIN_DIR)/$(BINARY) --playground=$$PLAY; \
	printf "$(GREEN)OK Playground(s) $$PLAY aplicado(s)$(NC)\n"

docker-playground-list: build ## Listar playgrounds disponibles (dinámico, descubre del registry)
	@printf "$(BLUE)Playgrounds disponibles:$(NC)\n"
	@./$(BIN_DIR)/$(BINARY) --playground=__list__ 2>&1 | grep -oE 'disponibles: \[[^]]+\]' | head -1 | tr -d '[]' | sed 's/disponibles: //' | tr ' ' '\n' | sed 's/^/  /'
	@printf "$(BLUE)Uso:$(NC)\n"
	@printf "  make docker-playground P=all          (todos)\n"
	@printf "  make docker-playground P=<name>       (uno)\n"
	@printf "  make docker-playground                (default: all)\n"

cloud-playground: build ## Aplicar playground(s) en cloud. P=all (default) | P=<name>
	@PLAY="$${P:-all}"; \
	set -a && . ./.env.cloud && set +a && ./$(BIN_DIR)/$(BINARY) --playground=$$PLAY

# ── Playgrounds v2 (sobre sistema completo L0..L4) ───────────────────

docker-playground-v2: build ## Recrear BD docker con playground v2. P=all (default) | P=<name>
	@PLAY="$${P:-all}"; \
	printf "$(BLUE)Levantando bases de datos y esperando healthchecks...$(NC)\n"; \
	docker compose -f docker-compose.local.yml up -d --wait; \
	printf "$(YELLOW)Recreando BD con playground-v2=$$PLAY (sistema completo + playground v2)...$(NC)\n"; \
	POSTGRES_PORT=5433 MONGO_PORT=27018 ./$(BIN_DIR)/$(BINARY) --playground-v2=$$PLAY; \
	printf "$(GREEN)OK Playground(s) v2 $$PLAY aplicado(s)$(NC)\n"

docker-playground-v2-list: build ## Listar playgrounds v2 disponibles (dinámico, descubre del registry)
	@printf "$(BLUE)Playgrounds v2 disponibles:$(NC)\n"
	@./$(BIN_DIR)/$(BINARY) --playground-v2=__list__ 2>&1 | grep -oE 'disponibles: \[[^]]+\]' | head -1 | tr -d '[]' | sed 's/disponibles: //' | tr ' ' '\n' | sed 's/^/  /'
	@printf "$(BLUE)Uso:$(NC)\n"
	@printf "  make docker-playground-v2 P=all       (todos)\n"
	@printf "  make docker-playground-v2 P=<name>    (uno)\n"
	@printf "  make docker-playground-v2             (default: all)\n"

cloud-playground-v2: build ## Aplicar playground(s) v2 en cloud. P=all (default) | P=<name>
	@PLAY="$${P:-all}"; \
	set -a && . ./.env.cloud && set +a && ./$(BIN_DIR)/$(BINARY) --playground-v2=$$PLAY

# ── E2E seeds (cloud) ────────────────────────────────────────────────

cloud-seed-scenario: ## Aplicar scenario(s) compositivo(s) Fase C (uso: make cloud-seed-scenario SCENARIO=teacher_grades_only)
	@if [ -z "$(SCENARIO)" ]; then \
		printf "$(YELLOW)⚠️  SCENARIO vacío. Uso: make cloud-seed-scenario SCENARIO=teacher_grades_only$(NC)\n"; \
		printf "$(YELLOW)    Disponibles: legacy_e2e, observer_audits, teacher_grades_only, guardian_views_child$(NC)\n"; \
		exit 1; \
	fi
	@printf "$(BLUE)Compilando seed_e2e...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/seed_e2e ./cmd/seed_e2e
	@printf "$(BLUE)Aplicando scenario(s) E2E (SCENARIO=$(SCENARIO)) en cloud...$(NC)\n"
	@set -a && . ./.env.cloud && set +a && SCENARIO=$(SCENARIO) ./$(BIN_DIR)/seed_e2e
	@printf "$(GREEN)OK Scenario(s) aplicado(s) — constantes en seeds/e2e/exports/fixtures-constants.json$(NC)\n"

# ── Auditor estático del seed ────────────────────────────────────────

seed-audit: ## Auditar el production seed (modo report-only, exit 0 siempre)
	@printf "$(BLUE)Auditando production seed (report-only)...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/seed-audit ./cmd/seed-audit
	@./$(BIN_DIR)/seed-audit --report-only --output-dir ./audit-reports

seed-audit-strict: ## Auditar el seed en modo strict (exit 1 si hay errores) — uso CI/precommit
	@printf "$(BLUE)Auditando production seed (strict)...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/seed-audit ./cmd/seed-audit
	@./$(BIN_DIR)/seed-audit --strict --output-dir ./audit-reports

seed-audit-test: ## Tests del auditor (falla si la cobertura es < 80%)
	@printf "$(BLUE)Tests seed-audit...$(NC)\n"
	@$(GO) test -count=1 -cover -coverprofile=/tmp/seed-audit-cover.out \
		./internal/seedaudit/... ./cmd/seed-audit/...
	@COV=$$($(GO) tool cover -func=/tmp/seed-audit-cover.out | grep '^total:' | awk '{print $$3}' | tr -d '%'); \
		printf "$(BLUE)Cobertura total: %s%%$(NC)\n" "$$COV"; \
		awk -v c="$$COV" 'BEGIN { if (c+0 < 80) { exit 1 } }' || \
		(printf "$(YELLOW)⚠️  cobertura $$COV%% < 80%%$(NC)\n" && exit 1)
	@printf "$(GREEN)OK seed-audit tests$(NC)\n"

# ── Cross-checker FE↔BE ──────────────────────────────────────────────
#
# Las paths default del binario son relativas al monorepo (EduGo/), por
# lo que estos targets cambian a la raíz antes de ejecutar. Si invocas
# el binario manualmente, hazlo desde EduGo/ o pasa --kmp-roots con
# paths absolutos.
MONOREPO_ROOT := $(abspath $(CURDIR)/../../..)
MIGRATOR_ABS  := $(abspath $(CURDIR))

contract-check: ## Auditar drift FE↔BE (KMP vs seed). Reporte informativo
	@printf "$(BLUE)Auditando drift FE↔BE...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/contract-check ./cmd/contract-check
	@cd $(MONOREPO_ROOT) && $(MIGRATOR_ABS)/$(BIN_DIR)/contract-check --output-dir $(MIGRATOR_ABS)/audit-reports

contract-check-strict: ## Auditar drift FE↔BE; exit 1 si hay errores. CI/precommit
	@printf "$(BLUE)Auditando drift FE↔BE (strict)...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/contract-check ./cmd/contract-check
	@cd $(MONOREPO_ROOT) && $(MIGRATOR_ABS)/$(BIN_DIR)/contract-check --severity error --output-dir $(MIGRATOR_ABS)/audit-reports

contract-check-update-baseline: ## Refrescar contract-check-baseline.json
	@printf "$(YELLOW)Refrescando baseline contract-check...$(NC)\n"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(LDFLAGS) -o $(BIN_DIR)/contract-check ./cmd/contract-check
	@cd $(MONOREPO_ROOT) && $(MIGRATOR_ABS)/$(BIN_DIR)/contract-check --update-baseline --output-dir $(MIGRATOR_ABS)/audit-reports/baseline
	@printf "$(GREEN)OK baseline actualizado$(NC)\n"

contract-check-test: ## Tests del cross-checker (falla si cobertura < 80%)
	@printf "$(BLUE)Tests contract-check...$(NC)\n"
	@$(GO) test -count=1 -cover -coverprofile=/tmp/contract-check-cover.out \
		./internal/contractcheck/... ./cmd/contract-check/...
	@COV=$$($(GO) tool cover -func=/tmp/contract-check-cover.out | grep '^total:' | awk '{print $$3}' | tr -d '%'); \
		printf "$(BLUE)Cobertura total: %s%%$(NC)\n" "$$COV"; \
		awk -v c="$$COV" 'BEGIN { if (c+0 < 80) { exit 1 } }' || \
		(printf "$(YELLOW)⚠️  cobertura $$COV%% < 80%%$(NC)\n" && exit 1)
	@printf "$(GREEN)OK contract-check tests$(NC)\n"

# ── Pipeline CI ──────────────────────────────────────────────────────

check: fmt vet lint test build ## Pipeline completo: fmt + vet + lint + test + build
	@printf "\n$(GREEN)Todos los checks pasaron$(NC)\n"
