# EduGo Infrastructure - Makefile raiz
# Orquesta desarrollo local, migraciones y operaciones multi-modulo

# Modulos Go
MODULES = mongodb postgres schemas tools/mock-generator

# Colores para output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

# Variables por defecto (pueden sobreescribirse con .env)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= edugo_dev
DB_USER ?= edugo
DB_PASSWORD ?= changeme

.PHONY: help
help: ## Mostrar ayuda
	@echo "$(BLUE)EduGo Infrastructure - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Modulos: $(MODULES)$(NC)"

# ===================
# DESARROLLO LOCAL
# ===================

.PHONY: dev-setup
dev-setup: ## Setup completo (primera vez)
	@echo "Iniciando setup completo de EduGo..."
	@./scripts/dev-setup.sh

.PHONY: dev-up-core
dev-up-core: ## Levantar solo PostgreSQL + MongoDB
	@echo "Levantando servicios core (PostgreSQL + MongoDB)..."
	@cd docker && docker-compose up -d postgres mongodb
	@echo "Servicios core corriendo"
	@echo "   PostgreSQL: localhost:5432"
	@echo "   MongoDB: localhost:27017"

.PHONY: dev-up-messaging
dev-up-messaging: ## Levantar core + RabbitMQ
	@echo "Levantando servicios core + messaging..."
	@cd docker && docker-compose --profile messaging up -d
	@echo "Servicios corriendo"
	@echo "   PostgreSQL: localhost:5432"
	@echo "   MongoDB: localhost:27017"
	@echo "   RabbitMQ: localhost:5672"
	@echo "   RabbitMQ UI: http://localhost:15672"

.PHONY: dev-up-full
dev-up-full: ## Levantar todos los servicios + tools
	@echo "Levantando todos los servicios..."
	@cd docker && docker-compose --profile messaging --profile cache --profile tools up -d
	@echo "Todos los servicios corriendo"
	@echo "   PostgreSQL: localhost:5432"
	@echo "   MongoDB: localhost:27017"
	@echo "   RabbitMQ: localhost:5672 (UI: http://localhost:15672)"
	@echo "   Redis: localhost:6379"
	@echo "   PgAdmin: http://localhost:5050"
	@echo "   Mongo Express: http://localhost:8082"

.PHONY: dev-logs
dev-logs: ## Ver logs de servicios
	@cd docker && docker-compose logs -f

.PHONY: dev-ps
dev-ps: ## Ver estado de servicios
	@cd docker && docker-compose ps

.PHONY: dev-down
dev-down: ## Detener servicios (mantener datos)
	@cd docker && docker-compose down

.PHONY: dev-teardown
dev-teardown: ## Detener y eliminar todo (incluye volumenes)
	@echo "Eliminando servicios y datos..."
	@cd docker && docker-compose down -v
	@echo "Ambiente limpio"

.PHONY: dev-reset
dev-reset: dev-teardown dev-setup ## Reset completo (teardown + setup)
	@echo "Ambiente reseteado completamente"

# ===================
# MIGRACIONES
# ===================

.PHONY: migrate-up
migrate-up: ## Ejecutar migraciones pendientes
	@echo "$(BLUE)Ejecutando migraciones PostgreSQL...$(NC)"
	@cd postgres && go run cmd/migrate/migrate.go up

.PHONY: migrate-down
migrate-down: ## Revertir ultima migracion
	@echo "$(BLUE)Revirtiendo migracion PostgreSQL...$(NC)"
	@cd postgres && go run cmd/migrate/migrate.go down

.PHONY: migrate-status
migrate-status: ## Ver estado de migraciones
	@cd postgres && go run cmd/migrate/migrate.go status

.PHONY: migrate-create
migrate-create: ## Crear nueva migracion (uso: make migrate-create NAME="add_column")
	@cd postgres && go run cmd/migrate/migrate.go create "$(NAME)"

.PHONY: runner-up
runner-up: ## Ejecutar runner de 4 capas
	@echo "$(BLUE)Ejecutando runner de 4 capas...$(NC)"
	@cd postgres && go run cmd/runner/runner.go

# ===================
# SEEDS
# ===================

.PHONY: seed
seed: ## Cargar datos de prueba
	@echo "$(BLUE)Cargando seeds...$(NC)"
	@./scripts/seed-data.sh

.PHONY: seed-minimal
seed-minimal: ## Cargar solo datos minimos
	@echo "$(BLUE)Cargando seeds minimos...$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h localhost -U $(DB_USER) -d $(DB_NAME) -f seeds/postgres/users.sql
	@PGPASSWORD=$(DB_PASSWORD) psql -h localhost -U $(DB_USER) -d $(DB_NAME) -f seeds/postgres/schools.sql

# ===================
# VALIDACION
# ===================

.PHONY: validate-env
validate-env: ## Validar variables de entorno
	@./scripts/validate-env.sh

.PHONY: validate-schemas
validate-schemas: ## Validar JSON schemas
	@cd schemas && go test -v ./...

# ===================
# MULTI-MODULO
# ===================

.PHONY: build-all
build-all: ## Compilar todos los modulos
	@echo "$(BLUE)Compilando todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Building $$module...$(NC)"; \
		(cd $$module && go build ./...) || exit 1; \
		echo "$(GREEN)  $$module compilado$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos compilados$(NC)"

.PHONY: test-all
test-all: ## Ejecutar tests unitarios en todos los modulos
	@echo "$(BLUE)Ejecutando tests en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Testing $$module...$(NC)"; \
		(cd $$module && go test -short -v ./...) || exit 1; \
		echo "$(GREEN)  $$module tests passed$(NC)"; \
		echo ""; \
	done
	@echo "$(GREEN)Todos los modulos pasaron los tests$(NC)"

.PHONY: lint-all
lint-all: ## Ejecutar linter en todos los modulos
	@echo "$(BLUE)Ejecutando linter en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Linting $$module...$(NC)"; \
		(cd $$module && golangci-lint run ./...) || exit 1; \
		echo "$(GREEN)  $$module linted$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos pasaron el linter$(NC)"

.PHONY: fmt-all
fmt-all: ## Formatear codigo en todos los modulos
	@echo "$(BLUE)Formateando codigo en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Formatting $$module...$(NC)"; \
		(cd $$module && go fmt ./...); \
		echo "$(GREEN)  $$module formatted$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos formateados$(NC)"

.PHONY: vet-all
vet-all: ## Ejecutar go vet en todos los modulos
	@echo "$(BLUE)Ejecutando go vet en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Vetting $$module...$(NC)"; \
		(cd $$module && go vet ./...) || exit 1; \
		echo "$(GREEN)  $$module vetted$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos pasaron go vet$(NC)"

.PHONY: tidy-all
tidy-all: ## Ejecutar go mod tidy en todos los modulos
	@echo "$(BLUE)Ejecutando go mod tidy en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Tidying $$module...$(NC)"; \
		(cd $$module && go mod tidy); \
		echo "$(GREEN)  $$module tidied$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos tidied$(NC)"

.PHONY: deps-all
deps-all: ## Actualizar dependencias en todos los modulos
	@echo "$(BLUE)Actualizando dependencias en todos los modulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)Updating $$module...$(NC)"; \
		(cd $$module && go get -u ./... && go mod tidy); \
		echo "$(GREEN)  $$module updated$(NC)"; \
	done
	@echo "$(GREEN)Todos los modulos actualizados$(NC)"

.PHONY: check-all
check-all: fmt-all vet-all lint-all test-all ## Validacion completa de todos los modulos
	@echo "$(GREEN)Validacion completa exitosa$(NC)"

# ===================
# UTILIDADES
# ===================

.PHONY: clean
clean: ## Limpiar archivos temporales
	@echo "$(BLUE)Limpiando...$(NC)"
	@rm -rf tmp/ temp/ *.log
	@for module in $(MODULES); do \
		(cd $$module && rm -rf build && go clean -testcache); \
	done
	@echo "$(GREEN)Limpieza completada$(NC)"

.PHONY: status
status: ## Ver estado general del ambiente
	@echo "$(BLUE)Estado del Ambiente EduGo$(NC)"
	@echo ""
	@echo "Docker:"
	@cd docker && docker-compose ps || echo "  Servicios no corriendo"
	@echo ""
	@echo "Migraciones PostgreSQL:"
	@cd postgres && go run cmd/migrate/migrate.go status 2>/dev/null || echo "  No ejecutadas aun"

.DEFAULT_GOAL := help
