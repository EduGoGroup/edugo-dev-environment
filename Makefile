# EduGo Dev Environment - Makefile
# Comandos simplificados para gestión del ambiente

.PHONY: help up down stop restart logs status clean setup validate diagnose seed update \
	migrator-build migrator-test migrator-lint migrator-check \
	db-migrate db-migrate-cloud db-recreate db-recreate-cloud

# Colores
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m
MIGRATOR_DIR := migrator
MIGRATOR_BIN := bin/migrator

help: ## Mostrar esta ayuda
	@echo ""
	@echo "$(BLUE)EduGo Dev Environment - Comandos disponibles$(NC)"
	@echo "=============================================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

# ==========================================
# COMANDOS PRINCIPALES
# ==========================================

up: ## Iniciar todos los servicios
	@echo "$(BLUE)🚀 Iniciando servicios...$(NC)"
	@cd docker && docker-compose up -d
	@echo "$(GREEN)✅ Servicios iniciados$(NC)"
	@$(MAKE) status

down: ## Detener y eliminar contenedores (mantiene datos)
	@echo "$(BLUE)🛑 Deteniendo servicios...$(NC)"
	@cd docker && docker-compose down
	@echo "$(GREEN)✅ Servicios detenidos$(NC)"

stop: ## Pausar servicios (mantiene contenedores)
	@echo "$(BLUE)⏸️  Pausando servicios...$(NC)"
	@cd docker && docker-compose stop
	@echo "$(GREEN)✅ Servicios pausados$(NC)"

restart: ## Reiniciar todos los servicios
	@echo "$(BLUE)🔄 Reiniciando servicios...$(NC)"
	@cd docker && docker-compose restart
	@echo "$(GREEN)✅ Servicios reiniciados$(NC)"

# ==========================================
# LOGS Y ESTADO
# ==========================================

logs: ## Ver logs de todos los servicios (Ctrl+C para salir)
	@cd docker && docker-compose logs -f

logs-api: ## Ver logs de API Mobile
	@cd docker && docker-compose logs -f api-mobile

logs-admin: ## Ver logs de API Admin
	@cd docker && docker-compose logs -f api-administracion

logs-worker: ## Ver logs de Worker
	@cd docker && docker-compose logs -f worker

status: ## Ver estado de los servicios
	@echo ""
	@echo "$(BLUE)📊 Estado de servicios$(NC)"
	@echo "======================"
	@cd docker && docker-compose ps
	@echo ""

# ==========================================
# SETUP Y CONFIGURACIÓN
# ==========================================

setup: ## Setup inicial completo
	@./scripts/setup.sh

setup-with-seeds: ## Setup completo con datos de prueba
	@./scripts/setup.sh --seed

validate: ## Validar configuración docker-compose
	@./scripts/validate.sh

diagnose: ## Ejecutar diagnóstico del ambiente
	@./scripts/diagnose.sh

seed: ## Recrear base de datos con migrator (edugo-infrastructure)
	@cd migrator && FORCE_MIGRATION=true go run cmd/main.go

update: ## Actualizar imágenes Docker
	@./scripts/update-images.sh

# ==========================================
# LIMPIEZA
# ==========================================

clean: ## Limpiar ambiente (interactivo)
	@./scripts/cleanup.sh

reset: ## Reset completo (⚠️ BORRA DATOS)
	@echo "$(YELLOW)⚠️  Esto eliminará todos los datos. ¿Continuar? [y/N]$(NC)"
	@read -r response; \
	if [ "$$response" = "y" ] || [ "$$response" = "Y" ]; then \
		cd docker && docker-compose down -v; \
		echo "$(GREEN)✅ Datos eliminados. Ejecuta 'make setup' para reiniciar.$(NC)"; \
	else \
		echo "Cancelado."; \
	fi

# ==========================================
# ACCESO A BASES DE DATOS
# ==========================================

psql: ## Conectar a PostgreSQL
	@docker exec -it edugo-postgres psql -U edugo -d edugo

mongo: ## Conectar a MongoDB
	@docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 edugo --authSource admin

# ==========================================
# HEALTH CHECKS
# ==========================================

health: ## Verificar health de las APIs
	@echo ""
	@echo "$(BLUE)🏥 Health Check$(NC)"
	@echo "==============="
	@echo ""
	@echo -n "API Mobile:  "
	@curl -s http://localhost:8081/health | head -c 100 || echo "❌ No responde"
	@echo ""
	@echo -n "API Admin:   "
	@curl -s http://localhost:8082/health | head -c 100 || echo "❌ No responde"
	@echo ""

# ==========================================
# MIGRATOR
# ==========================================

migrator-build: ## Compilar proyecto migrator
	@echo "$(BLUE)🔨 Compilando migrator...$(NC)"
	@mkdir -p $(MIGRATOR_DIR)/bin
	@cd $(MIGRATOR_DIR) && go build -o $(MIGRATOR_BIN) ./cmd
	@echo "$(GREEN)✅ Migrator compilado: $(MIGRATOR_DIR)/$(MIGRATOR_BIN)$(NC)"

migrator-test: ## Ejecutar tests del migrator
	@echo "$(BLUE)🧪 Ejecutando tests de migrator...$(NC)"
	@cd $(MIGRATOR_DIR) && go test -v ./tests
	@echo "$(GREEN)✅ Tests de migrator completados$(NC)"

migrator-lint: ## Ejecutar lint del migrator
	@echo "$(BLUE)🧹 Ejecutando lint de migrator...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(MIGRATOR_DIR) && golangci-lint run ./...; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint no está instalado. Ejecutando go vet como fallback...$(NC)"; \
		cd $(MIGRATOR_DIR) && go vet ./...; \
	fi
	@echo "$(GREEN)✅ Lint de migrator completado$(NC)"

migrator-check: migrator-lint migrator-test migrator-build ## Ejecutar lint, tests y compilación del migrator

# ==========================================
# DATABASE MIGRATIONS
# ==========================================

db-migrate: migrator-build ## Ejecutar migraciones (idempotente, Docker local)
	@echo "$(BLUE)📦 Ejecutando migraciones (Docker local)...$(NC)"
	@cd $(MIGRATOR_DIR) && ./$(MIGRATOR_BIN)
	@echo "$(GREEN)✅ Migraciones completadas$(NC)"

db-migrate-cloud: migrator-build ## Ejecutar migraciones en cloud (idempotente). Requiere .env.cloud
	@echo "$(BLUE)☁️  Ejecutando migraciones (cloud)...$(NC)"
	@if [ ! -f docker/.env.cloud ]; then \
		echo "$(YELLOW)⚠️  Archivo docker/.env.cloud no encontrado. Copia docker/.env.cloud.example y configura tus valores.$(NC)"; \
		exit 1; \
	fi
	@set -a && . docker/.env.cloud && set +a && cd $(MIGRATOR_DIR) && ./$(MIGRATOR_BIN)
	@echo "$(GREEN)✅ Migraciones cloud completadas$(NC)"

db-recreate: migrator-build ## Recrear BD en Docker local. DESTRUYE DATOS
	@echo "$(YELLOW)⚠️  Esto eliminará y recreará TODAS las bases de datos locales. ¿Continuar? [y/N]$(NC)"
	@read -r response; \
	if [ "$$response" = "y" ] || [ "$$response" = "Y" ]; then \
		cd $(MIGRATOR_DIR) && FORCE_MIGRATION=true ./$(MIGRATOR_BIN); \
		echo "$(GREEN)✅ Bases de datos recreadas$(NC)"; \
	else \
		echo "Cancelado."; \
	fi

db-recreate-cloud: migrator-build ## Recrear BD en cloud. Requiere .env.cloud. DESTRUYE DATOS
	@echo "$(YELLOW)⚠️  Esto eliminará y recreará TODAS las bases de datos en CLOUD. ¿Continuar? [y/N]$(NC)"
	@read -r response; \
	if [ "$$response" = "y" ] || [ "$$response" = "Y" ]; then \
		if [ ! -f docker/.env.cloud ]; then \
			echo "$(YELLOW)⚠️  Archivo docker/.env.cloud no encontrado. Copia docker/.env.cloud.example y configura tus valores.$(NC)"; \
			exit 1; \
		fi; \
		set -a && . docker/.env.cloud && set +a && cd $(MIGRATOR_DIR) && FORCE_MIGRATION=true ./$(MIGRATOR_BIN); \
		echo "$(GREEN)✅ Bases de datos cloud recreadas$(NC)"; \
	else \
		echo "Cancelado."; \
	fi
