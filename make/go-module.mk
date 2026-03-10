SHELL := /bin/bash

ROOT_DIR ?= ..
GO ?= go
GOLANGCI_LINT ?= golangci-lint
BUILD_COMMAND ?= $(GO) build ./...
TEST_FLAGS ?= -short -v
TEST_ALL_FLAGS ?= -v
TEST_RACE_FLAGS ?= -short -v
BUILD_DIR ?= build

RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

include $(ROOT_DIR)/make/module-release.mk

.PHONY: help build test test-all test-race lint fmt fmt-check vet tidy deps check release-check clean

help: ## Mostrar ayuda
	@echo "$(BLUE)$(MODULE_NAME) - Comandos disponibles:$(NC)"
	@echo ""
	@grep -hE '^[a-zA-Z0-9_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-22s$(NC) %s\n", $$1, $$2}'

build: ## Compilar el módulo
	@echo "$(BLUE)Compilando $(MODULE_NAME)...$(NC)"
	@$(BUILD_COMMAND)
	@echo "$(GREEN)Compilación exitosa$(NC)"

test: ## Ejecutar tests unitarios
	@echo "$(BLUE)Ejecutando tests...$(NC)"
	@$(GO) test $(TEST_FLAGS) ./...
	@echo "$(GREEN)Tests completados$(NC)"

test-all: ## Ejecutar todos los tests
	@echo "$(BLUE)Ejecutando todos los tests...$(NC)"
	@$(GO) test $(TEST_ALL_FLAGS) ./...
	@echo "$(GREEN)Tests completados$(NC)"

test-race: ## Ejecutar tests con race detection
	@echo "$(BLUE)Ejecutando tests con race detection...$(NC)"
	@$(GO) test $(TEST_RACE_FLAGS) -race ./...
	@echo "$(GREEN)Tests con race detection completados$(NC)"

lint: ## Ejecutar golangci-lint
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "$(RED)golangci-lint no está instalado$(NC)"; exit 1; }
	@echo "$(BLUE)Ejecutando linter...$(NC)"
	@$(GOLANGCI_LINT) run --allow-parallel-runners ./...
	@echo "$(GREEN)Linter completado$(NC)"

fmt: ## Formatear código Go
	@echo "$(BLUE)Formateando código...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)Código formateado$(NC)"

fmt-check: ## Validar que el código Go ya esté formateado
	@unformatted=$$(find . -type f -name '*.go' -not -path './vendor/*' -not -path './build/*' -not -path './bin/*' -exec gofmt -l {} +); \
	if [ -z "$$unformatted" ]; then \
		echo "$(GREEN)Formato validado$(NC)"; \
		exit 0; \
	fi; \
	if [ -n "$$unformatted" ]; then \
		echo "$(RED)Archivos sin formatear:$(NC)"; \
		echo "$$unformatted"; \
		exit 1; \
	fi; \
	echo "$(GREEN)Formato validado$(NC)"

vet: ## Ejecutar go vet
	@echo "$(BLUE)Ejecutando go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)go vet completado$(NC)"

tidy: ## Ejecutar go mod tidy
	@echo "$(BLUE)Ejecutando go mod tidy...$(NC)"
	@$(GO) mod tidy
	@echo "$(GREEN)go.mod actualizado$(NC)"

deps: ## Actualizar dependencias del módulo
	@echo "$(BLUE)Actualizando dependencias...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(GREEN)Dependencias actualizadas$(NC)"

check: fmt vet lint test build ## Validación completa del módulo
	@echo "$(GREEN)Validación completa exitosa$(NC)"

release-check: fmt-check vet lint test build changelog-check ## Validación no destructiva para release
	@echo "$(GREEN)Módulo listo para release$(NC)"

clean: ## Limpiar artefactos del módulo
	@rm -rf $(BUILD_DIR)
	@rm -rf bin
	@$(GO) clean -testcache
	@echo "$(GREEN)Limpieza completada$(NC)"

.DEFAULT_GOAL := help
