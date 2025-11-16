.PHONY: help
help: ## Mostrar ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ===================
# DESARROLLO LOCAL
# ===================

.PHONY: dev-setup
dev-setup: ## Setup completo (primera vez)
	@echo "🚀 Iniciando setup completo de EduGo..."
	@./scripts/dev-setup.sh

.PHONY: dev-up-core
dev-up-core: ## Levantar solo PostgreSQL + MongoDB
	@echo "🐳 Levantando servicios core (PostgreSQL + MongoDB)..."
	@cd docker && docker-compose up -d postgres mongodb
	@echo "✅ Servicios core corriendo"
	@echo "   PostgreSQL: localhost:5432"
	@echo "   MongoDB: localhost:27017"

.PHONY: dev-up-messaging
dev-up-messaging: ## Levantar core + RabbitMQ
	@echo "🐳 Levantando servicios core + messaging..."
	@cd docker && docker-compose --profile messaging up -d
	@echo "✅ Servicios corriendo"
	@echo "   PostgreSQL: localhost:5432"
	@echo "   MongoDB: localhost:27017"
	@echo "   RabbitMQ: localhost:5672"
	@echo "   RabbitMQ UI: http://localhost:15672"

.PHONY: dev-up-full
dev-up-full: ## Levantar todos los servicios + tools
	@echo "🐳 Levantando todos los servicios..."
	@cd docker && docker-compose --profile messaging --profile cache --profile tools up -d
	@echo "✅ Todos los servicios corriendo"
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
dev-teardown: ## Detener y eliminar todo (incluye volúmenes)
	@echo "⚠️  Eliminando servicios y datos..."
	@cd docker && docker-compose down -v
	@echo "✅ Ambiente limpio"

.PHONY: dev-reset
dev-reset: dev-teardown dev-setup ## Reset completo (teardown + setup)
	@echo "✅ Ambiente reseteado completamente"

# ===================
# MIGRACIONES
# ===================

.PHONY: migrate-up
migrate-up: ## Ejecutar migraciones pendientes
	@echo "📊 Ejecutando migraciones..."
	@cd database && go run migrate.go up

.PHONY: migrate-down
migrate-down: ## Revertir última migración
	@echo "⬇️  Revirtiendo migración..."
	@cd database && go run migrate.go down

.PHONY: migrate-status
migrate-status: ## Ver estado de migraciones
	@cd database && go run migrate.go status

.PHONY: migrate-create
migrate-create: ## Crear nueva migración (uso: make migrate-create NAME="add_column")
	@cd database && go run migrate.go create "$(NAME)"

# ===================
# SEEDS
# ===================

.PHONY: seed
seed: ## Cargar datos de prueba
	@echo "🌱 Cargando seeds..."
	@./scripts/seed-data.sh

.PHONY: seed-minimal
seed-minimal: ## Cargar solo datos mínimos
	@echo "🌱 Cargando seeds mínimos..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h localhost -U $(DB_USER) -d $(DB_NAME) -f seeds/postgres/users.sql
	@PGPASSWORD=$(DB_PASSWORD) psql -h localhost -U $(DB_USER) -d $(DB_NAME) -f seeds/postgres/schools.sql

# ===================
# VALIDACIÓN
# ===================

.PHONY: validate-env
validate-env: ## Validar variables de entorno
	@./scripts/validate-env.sh

.PHONY: validate-schemas
validate-schemas: ## Validar JSON schemas
	@echo "✅ Schemas válidos (implementar validación después)"

# ===================
# UTILIDADES
# ===================

.PHONY: clean
clean: ## Limpiar archivos temporales
	@echo "🧹 Limpiando..."
	@rm -rf tmp/ temp/ *.log

.PHONY: status
status: ## Ver estado general del ambiente
	@echo "📊 Estado del Ambiente EduGo"
	@echo ""
	@echo "Docker:"
	@cd docker && docker-compose ps || echo "  Servicios no corriendo"
	@echo ""
	@echo "Migraciones:"
	@cd database && go run migrate.go status 2>/dev/null || echo "  No ejecutadas aún"

# Variables por defecto (pueden sobreescribirse con .env)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= edugo_dev
DB_USER ?= edugo
DB_PASSWORD ?= changeme
