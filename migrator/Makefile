# EduGo Migrator
# Makefile independiente — ejecutar desde migrator/

.PHONY: help build clean fmt vet lint tidy test test-postgres test-mongo run status docker-up docker-down docker-migrate docker-recreate check

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

docker-migrate: docker-up build ## Levantar DBs y aplicar migraciones (flujo completo local)
	@printf "$(BLUE)Aplicando migraciones en BD local...$(NC)\n"
	@sleep 3
	@POSTGRES_PORT=5433 MONGO_PORT=27018 ./$(BIN_DIR)/$(BINARY)
	@printf "$(GREEN)OK Migraciones aplicadas$(NC)\n"

docker-recreate: docker-up build ## Recrear BDs desde cero con datos reales (BORRA DATOS)
	@printf "$(YELLOW)Recreando bases de datos...$(NC)\n"
	@sleep 3
	@POSTGRES_PORT=5433 MONGO_PORT=27018 FORCE_MIGRATION=true ./$(BIN_DIR)/$(BINARY)
	@printf "$(GREEN)OK Bases de datos recreadas$(NC)\n"

# ── Pipeline CI ──────────────────────────────────────────────────────

check: fmt vet lint test build ## Pipeline completo: fmt + vet + lint + test + build
	@printf "\n$(GREEN)Todos los checks pasaron$(NC)\n"
