# Makefile para postgres
MODULE_NAME = github.com/EduGoGroup/edugo-infrastructure/postgres
BUILD_DIR = build

RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

.PHONY: help
help: ## Mostrar ayuda
	@echo "$(BLUE)$(MODULE_NAME) - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# ============================================================================
# Desarrollo
# ============================================================================

.PHONY: build
build: ## Compilar el modulo
	@echo "$(BLUE)Compilando $(MODULE_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/ ./...
	@echo "$(GREEN)Compilacion exitosa$(NC)"

.PHONY: test
test: ## Ejecutar tests unitarios
	@echo "$(BLUE)Ejecutando tests...$(NC)"
	@go test -short -v ./...
	@echo "$(GREEN)Tests completados$(NC)"

.PHONY: test-all
test-all: ## Ejecutar todos los tests (incluye integracion, requiere Docker)
	@echo "$(BLUE)Ejecutando todos los tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)Tests completados$(NC)"

.PHONY: test-race
test-race: ## Tests con race detection
	@echo "$(BLUE)Ejecutando tests con race detection...$(NC)"
	@go test -race -short -v ./...
	@echo "$(GREEN)Tests con race detection completados$(NC)"

.PHONY: lint
lint: ## Ejecutar linter
	@echo "$(BLUE)Ejecutando linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)Linter completado$(NC)"

.PHONY: fmt
fmt: ## Formatear codigo
	@echo "$(BLUE)Formateando codigo...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Codigo formateado$(NC)"

.PHONY: vet
vet: ## Analisis estatico
	@echo "$(BLUE)Ejecutando go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Analisis completado$(NC)"

.PHONY: tidy
tidy: ## Limpiar go.mod
	@echo "$(BLUE)Ejecutando go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)go.mod limpio$(NC)"

.PHONY: deps
deps: ## Actualizar dependencias
	@echo "$(BLUE)Actualizando dependencias...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)Dependencias actualizadas$(NC)"

.PHONY: check
check: fmt vet lint test ## Validacion completa

.PHONY: clean
clean: ## Limpiar archivos generados
	@rm -rf $(BUILD_DIR)
	@go clean -testcache
	@echo "$(GREEN)Limpieza completada$(NC)"

# ============================================================================
# Migraciones
# ============================================================================

.PHONY: migrate-up
migrate-up: ## Ejecutar migraciones pendientes
	@echo "$(BLUE)Ejecutando migraciones PostgreSQL...$(NC)"
	@go run cmd/migrate/migrate.go up

.PHONY: migrate-down
migrate-down: ## Revertir ultima migracion
	@echo "$(BLUE)Revirtiendo ultima migracion...$(NC)"
	@go run cmd/migrate/migrate.go down

.PHONY: migrate-status
migrate-status: ## Ver estado de migraciones
	@echo "$(BLUE)Estado de migraciones:$(NC)"
	@go run cmd/migrate/migrate.go status

.PHONY: migrate-create
migrate-create: ## Crear nueva migracion (uso: make migrate-create NAME=descripcion)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: Debe especificar NAME=descripcion$(NC)"; \
		exit 1; \
	fi
	@go run cmd/migrate/migrate.go create "$(NAME)"

.PHONY: runner-up
runner-up: ## Ejecutar runner de 4 capas (estructura + constraints + seeds + testing)
	@echo "$(BLUE)Ejecutando runner de 4 capas...$(NC)"
	@go run cmd/runner/runner.go

.PHONY: seed
seed: ## Cargar datos de prueba
	@echo "$(BLUE)Cargando seeds PostgreSQL...$(NC)"
	@if [ -d "seeds" ]; then \
		for file in seeds/*.sql; do \
			echo "$(YELLOW)Ejecutando $$file...$(NC)"; \
			psql "$(shell go run -tags helper ../scripts/get-db-url.go)" -f $$file; \
		done; \
	else \
		echo "$(RED)No se encontro el directorio seeds/$(NC)"; \
	fi

.DEFAULT_GOAL := help
