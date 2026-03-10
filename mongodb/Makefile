MODULE_NAME = github.com/EduGoGroup/edugo-infrastructure/mongodb
MODULE_PATH = mongodb
ROOT_DIR = ..
BUILD_DIR = build
BUILD_COMMAND = mkdir -p $(BUILD_DIR) && go build -o $(BUILD_DIR)/migrate ./cmd/migrate && go build -o $(BUILD_DIR)/runner ./cmd/runner && go build -o $(BUILD_DIR)/seed ./cmd/seed

include ../make/go-module.mk

.PHONY: migrate-status migrate-force runner-up runner-structure runner-constraints seed seed-canonical seed-mock seed-all

migrate-status: ## Ver estado de migraciones MongoDB
	@go run ./cmd/migrate status

migrate-force: ## Forzar versión de migración (uso: make migrate-force VERSION_NUMBER=3)
	@if [ -z "$(VERSION_NUMBER)" ]; then \
		echo "$(RED)Debe especificar VERSION_NUMBER=<numero>$(NC)"; \
		exit 1; \
	fi
	@go run ./cmd/migrate force "$(VERSION_NUMBER)"

runner-up: ## Aplicar structure + constraints
	@go run ./cmd/runner all

runner-structure: ## Aplicar solo structure
	@go run ./cmd/runner structure

runner-constraints: ## Aplicar solo constraints
	@go run ./cmd/runner constraints

seed-canonical: ## Aplicar seeds canónicos
	@go run ./cmd/seed canonical

seed-mock: ## Aplicar mock data
	@go run ./cmd/seed mock

seed-all: ## Aplicar seeds canónicos + mock data
	@go run ./cmd/seed all

seed: seed-all ## Alias de seeds completos
	@true
