SHELL := /bin/bash

MODULE_NAME = docker
MODULE_PATH = docker
ROOT_DIR = ..
COMPOSE ?= docker compose

include ../make/module-release.mk

.PHONY: help validate up-core up-messaging up-full logs ps down teardown release-check

help: ## Mostrar ayuda
	@echo "$(BLUE)docker - Comandos disponibles:$(NC)"
	@echo ""
	@grep -hE '^[a-zA-Z0-9_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-22s$(NC) %s\n", $$1, $$2}'

validate: ## Validar docker-compose.yml
	@$(COMPOSE) -f docker-compose.yml config -q
	@echo "$(GREEN)docker-compose.yml validado$(NC)"

up-core: ## Levantar PostgreSQL + MongoDB
	@$(COMPOSE) up -d postgres mongodb

up-messaging: ## Levantar core + RabbitMQ
	@$(COMPOSE) --profile messaging up -d

up-full: ## Levantar core + RabbitMQ + Redis + tools
	@$(COMPOSE) --profile messaging --profile cache --profile tools up -d

logs: ## Ver logs del stack local
	@$(COMPOSE) logs -f

ps: ## Ver estado del stack local
	@$(COMPOSE) ps

down: ## Detener servicios conservando volúmenes
	@$(COMPOSE) down

teardown: ## Detener servicios y eliminar volúmenes
	@$(COMPOSE) down -v

release-check: validate changelog-check ## Validación para release de docker/
	@echo "$(GREEN)Módulo docker listo para release$(NC)"

.DEFAULT_GOAL := help
