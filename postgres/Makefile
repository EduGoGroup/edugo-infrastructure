MODULE_NAME = github.com/EduGoGroup/edugo-infrastructure/postgres
MODULE_PATH = postgres
ROOT_DIR = ..
BUILD_DIR = build
BUILD_COMMAND = mkdir -p $(BUILD_DIR) && go build -o $(BUILD_DIR)/runner ./cmd/runner && go build -o $(BUILD_DIR)/seed ./cmd/seed

include ../make/go-module.mk

.PHONY: runner-up seed seed-production seed-development seed-all

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
