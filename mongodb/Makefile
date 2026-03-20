MODULE_NAME = github.com/EduGoGroup/edugo-infrastructure/mongodb
MODULE_PATH = mongodb
ROOT_DIR = ..
BUILD_DIR = build
BUILD_COMMAND = mkdir -p $(BUILD_DIR) && go build -o $(BUILD_DIR)/runner ./cmd/runner && go build -o $(BUILD_DIR)/seed ./cmd/seed

include ../make/go-module.mk

.PHONY: runner-up runner-structure runner-constraints seed seed-canonical seed-mock seed-all

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
