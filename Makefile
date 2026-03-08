SHELL := /bin/bash

GO_MODULES = postgres mongodb schemas tools/mock-generator
RELEASE_MODULES = postgres mongodb schemas tools/mock-generator docker

RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= edugo_dev
DB_USER ?= edugo
DB_PASSWORD ?= changeme
DB_SSL_MODE ?= disable
MONGO_HOST ?= localhost
MONGO_PORT ?= 27017
MONGO_DB ?= edugo

.PHONY: help require-module dev-setup dev-up-core dev-up-messaging dev-up-full dev-logs dev-ps dev-down dev-teardown dev-reset db-bootstrap migrate-up migrate-down migrate-status migrate-create runner-up seed seed-production seed-development seed-minimal validate-env validate-schemas build-all test-all lint-all fmt-all fmt-check-all vet-all tidy-all deps-all check-all release-check-all release-check release-prepare release-notes release-tag release-push-tag release-github clean status

help: ## Mostrar ayuda
	@echo "$(BLUE)EduGo Infrastructure - Comandos disponibles:$(NC)"
	@echo ""
	@grep -hE '^[a-zA-Z0-9_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Módulos Go: $(GO_MODULES)$(NC)"
	@echo "$(YELLOW)Módulos release: $(RELEASE_MODULES)$(NC)"

require-module:
	@if [ -z "$(MODULE)" ]; then \
		echo "$(RED)Debe especificar MODULE=<ruta-del-modulo>$(NC)"; \
		exit 1; \
	fi
	@if [ ! -d "$(MODULE)" ]; then \
		echo "$(RED)No existe el módulo $(MODULE)$(NC)"; \
		exit 1; \
	fi
	@if [ ! -f "$(MODULE)/Makefile" ]; then \
		echo "$(RED)El módulo $(MODULE) no tiene Makefile$(NC)"; \
		exit 1; \
	fi

# ===================
# DESARROLLO LOCAL
# ===================

dev-setup: ## Setup local del repo (Docker + bootstrap de datos)
	@./scripts/dev-setup.sh

dev-up-core: ## Levantar PostgreSQL + MongoDB
	@$(MAKE) -C docker up-core

dev-up-messaging: ## Levantar core + RabbitMQ
	@$(MAKE) -C docker up-messaging

dev-up-full: ## Levantar core + RabbitMQ + Redis + tools
	@$(MAKE) -C docker up-full

dev-logs: ## Ver logs del stack local
	@$(MAKE) -C docker logs

dev-ps: ## Ver estado del stack local
	@$(MAKE) -C docker ps

dev-down: ## Detener servicios locales
	@$(MAKE) -C docker down

dev-teardown: ## Detener y eliminar servicios/volúmenes locales
	@$(MAKE) -C docker teardown

dev-reset: dev-teardown dev-setup ## Reiniciar completamente el stack local
	@echo "$(GREEN)Ambiente local reiniciado$(NC)"

db-bootstrap: ## Aplicar estructura y datos base usando los runners del repo
	@echo "$(BLUE)Aplicando bootstrap PostgreSQL...$(NC)"
	@$(MAKE) -C postgres runner-up
	@echo "$(BLUE)Aplicando bootstrap MongoDB...$(NC)"
	@$(MAKE) -C mongodb runner-up
	@$(MAKE) -C mongodb seed-all
	@echo "$(GREEN)Bootstrap completado$(NC)"

# ===================
# MIGRACIONES Y SEEDS
# ===================

migrate-up: ## Ejecutar migraciones legacy pendientes de PostgreSQL
	@$(MAKE) -C postgres migrate-up

migrate-down: ## Revertir última migración legacy de PostgreSQL
	@$(MAKE) -C postgres migrate-down

migrate-status: ## Ver estado de migraciones legacy de PostgreSQL
	@$(MAKE) -C postgres migrate-status

migrate-create: ## Crear nueva migración legacy de PostgreSQL (uso: make migrate-create NAME=nombre)
	@$(MAKE) -C postgres migrate-create NAME="$(NAME)"

runner-up: ## Ejecutar runners embebidos de PostgreSQL y MongoDB
	@$(MAKE) -C postgres runner-up
	@$(MAKE) -C mongodb runner-up

seed: ## Aplicar seeds de PostgreSQL y MongoDB
	@./scripts/seed-data.sh

seed-production: ## Aplicar solo datos canónicos de PostgreSQL y MongoDB
	@$(MAKE) -C postgres seed-production
	@$(MAKE) -C mongodb seed-canonical

seed-development: ## Aplicar datos de desarrollo completos
	@$(MAKE) -C postgres seed-all
	@$(MAKE) -C mongodb seed-all

seed-minimal: seed-production ## Alias de datos mínimos/canónicos
	@true

# ===================
# VALIDACIÓN
# ===================

validate-env: ## Validar variables de entorno locales
	@./scripts/validate-env.sh

validate-schemas: ## Ejecutar tests del módulo schemas
	@$(MAKE) -C schemas test

# ===================
# MULTI-MÓDULO
# ===================

build-all: ## Compilar todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Building $$module...$(NC)"; \
		$(MAKE) -C $$module build || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Build completado para todos los módulos$(NC)"

test-all: ## Ejecutar tests en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Testing $$module...$(NC)"; \
		$(MAKE) -C $$module test || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Tests completados para todos los módulos$(NC)"

lint-all: ## Ejecutar linter en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Linting $$module...$(NC)"; \
		$(MAKE) -C $$module lint || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Lint completado para todos los módulos$(NC)"

fmt-all: ## Formatear código en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Formatting $$module...$(NC)"; \
		$(MAKE) -C $$module fmt || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Formato aplicado a todos los módulos$(NC)"

fmt-check-all: ## Validar formato en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Checking fmt $$module...$(NC)"; \
		$(MAKE) -C $$module fmt-check || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Formato validado en todos los módulos$(NC)"

vet-all: ## Ejecutar go vet en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Vetting $$module...$(NC)"; \
		$(MAKE) -C $$module vet || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)go vet completado en todos los módulos$(NC)"

tidy-all: ## Ejecutar go mod tidy en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Tidying $$module...$(NC)"; \
		$(MAKE) -C $$module tidy || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)go mod tidy completado en todos los módulos$(NC)"

deps-all: ## Actualizar dependencias en todos los módulos Go
	@for module in $(GO_MODULES); do \
		echo "$(YELLOW)Updating deps $$module...$(NC)"; \
		$(MAKE) -C $$module deps || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Dependencias actualizadas en todos los módulos$(NC)"

check-all: fmt-all vet-all lint-all test-all build-all ## Validación completa con formateo
	@echo "$(GREEN)Validación completa exitosa$(NC)"

release-check-all: ## Validar todos los módulos listos para release
	@for module in $(RELEASE_MODULES); do \
		echo "$(YELLOW)Release check $$module...$(NC)"; \
		$(MAKE) -C $$module release-check || exit 1; \
		echo ""; \
	done
	@echo "$(GREEN)Todos los módulos están validados para release$(NC)"

release-check: require-module ## Ejecutar release-check en un módulo (MODULE=postgres)
	@$(MAKE) -C "$(MODULE)" release-check

release-prepare: require-module ## Actualizar CHANGELOG para VERSION en un módulo
	@$(MAKE) -C "$(MODULE)" release-prepare VERSION="$(VERSION)"

release-notes: require-module ## Mostrar notas de release de un módulo
	@$(MAKE) -C "$(MODULE)" release-notes VERSION="$(VERSION)"

release-tag: require-module ## Crear tag local de un módulo
	@$(MAKE) -C "$(MODULE)" release-tag VERSION="$(VERSION)"

release-push-tag: require-module ## Publicar tag de un módulo
	@$(MAKE) -C "$(MODULE)" release-push-tag VERSION="$(VERSION)"

release-github: require-module ## Crear GitHub Release de un módulo
	@$(MAKE) -C "$(MODULE)" release-github VERSION="$(VERSION)"

# ===================
# UTILIDADES
# ===================

clean: ## Limpiar artefactos temporales del repo y módulos
	@rm -rf tmp temp *.log logs/*.log
	@for module in $(GO_MODULES); do \
		$(MAKE) -C $$module clean || exit 1; \
	done
	@echo "$(GREEN)Limpieza completada$(NC)"

status: ## Ver estado general del repo local
	@echo "$(BLUE)Estado del stack local$(NC)"
	@$(MAKE) -C docker ps || true
	@echo ""
	@echo "$(BLUE)Estado migraciones PostgreSQL$(NC)"
	@$(MAKE) -C postgres migrate-status || true

.DEFAULT_GOAL := help
