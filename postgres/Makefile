MODULE_NAME = github.com/EduGoGroup/edugo-infrastructure/postgres
MODULE_PATH = postgres
ROOT_DIR = ..
BUILD_DIR = build
BUILD_COMMAND = mkdir -p $(BUILD_DIR) && go build -o $(BUILD_DIR)/migrate ./cmd/migrate && go build -o $(BUILD_DIR)/runner ./cmd/runner && go build -o $(BUILD_DIR)/seed ./cmd/seed

include ../make/go-module.mk

.PHONY: migrate-up migrate-down migrate-status migrate-create runner-up seed seed-production seed-development seed-all

migrate-up: ## Ejecutar migraciones legacy pendientes
	@go run ./cmd/migrate up

migrate-down: ## Revertir última migración legacy
	@go run ./cmd/migrate down

migrate-status: ## Ver estado de migraciones legacy
	@go run ./cmd/migrate status

migrate-create: ## Crear nueva migración legacy (uso: make migrate-create NAME=descripcion)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Debe especificar NAME=descripcion$(NC)"; \
		exit 1; \
	fi
	@go run ./cmd/migrate create "$(NAME)"

runner-up: ## Aplicar estructura + seeds de producción + seeds de desarrollo
	@go run ./cmd/runner all

seed-production: ## Aplicar solo seeds de producción
	@go run ./cmd/seed production

seed-development: ## Aplicar solo seeds de desarrollo
	@go run ./cmd/seed development

seed-all: ## Aplicar seeds de producción + desarrollo
	@go run ./cmd/seed all

seed: seed-all ## Alias de seeds completos
	@true
