MODULE_NAME = github.com/EduGoGroup/edugo-infrastructure/postgres
MODULE_PATH = postgres
ROOT_DIR = ..
BUILD_DIR = build
BUILD_COMMAND = mkdir -p $(BUILD_DIR) && go build -o $(BUILD_DIR)/runner ./cmd/runner && go build -o $(BUILD_DIR)/seed ./cmd/seed

include ../make/go-module.mk

.PHONY: test-integration runner-up seed seed-production seed-development seed-all

# Gate local obligatorio (plan 054, S.6): la CI no lo corre. Levanta Postgres con
# testcontainers, por eso pide Docker. GOWORK=off para medir contra los tags reales.
# `-p 1`: un paquete a la vez. Cada paquete arranca su propio Postgres con 60 s de espera;
# en paralelo y con la máquina cargada, el arranque vence y el test falla sin ser un fallo real.
TEST_INTEGRATION_FLAGS ?= -tags integration -count=1 -timeout 30m -v -p 1

test-integration: ## Ejecutar tests de integración (requiere Docker; gate local, no de CI)
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker no está instalado$(NC)"; exit 1; }
	@docker info >/dev/null 2>&1 || { echo "$(RED)Docker no está corriendo$(NC)"; exit 1; }
	@echo "$(BLUE)Ejecutando tests de integración...$(NC)"
	@ENABLE_INTEGRATION_TESTS=true GOWORK=off $(GO) test $(TEST_INTEGRATION_FLAGS) ./...
	@echo "$(GREEN)Tests de integración completados$(NC)"

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
