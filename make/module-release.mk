SHELL := /bin/bash

ROOT_DIR ?= ..
MODULE_PATH ?= $(notdir $(CURDIR))
CHANGELOG_FILE ?= CHANGELOG.md
RELEASE_SCRIPT := $(abspath $(ROOT_DIR)/scripts/module-release.sh)

RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m

.PHONY: changelog-check guard-version release-prepare release-notes release-tag release-push-tag release-github

changelog-check: ## Validar que exista CHANGELOG y sección Unreleased
	@if [ ! -f "$(CHANGELOG_FILE)" ]; then \
		echo "$(RED)Falta $(CHANGELOG_FILE)$(NC)"; \
		exit 1; \
	fi
	@if ! grep -q '^## \[Unreleased\]' "$(CHANGELOG_FILE)"; then \
		echo "$(RED)$(CHANGELOG_FILE) debe contener '## [Unreleased]'$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)CHANGELOG validado$(NC)"

guard-version:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Debe especificar VERSION=vX.Y.Z$(NC)"; \
		exit 1; \
	fi

release-prepare: guard-version release-check ## Congelar Unreleased en CHANGELOG para VERSION
	@$(RELEASE_SCRIPT) prepare $(MODULE_PATH) $(VERSION)

release-notes: guard-version ## Mostrar notas de release para VERSION
	@$(RELEASE_SCRIPT) notes $(MODULE_PATH) $(VERSION)

release-tag: guard-version ## Crear tag local del módulo
	@$(RELEASE_SCRIPT) notes $(MODULE_PATH) $(VERSION) >/dev/null
	@if git rev-parse "$(MODULE_PATH)/$(VERSION)" >/dev/null 2>&1; then \
		echo "$(RED)El tag $(MODULE_PATH)/$(VERSION) ya existe$(NC)"; \
		exit 1; \
	fi
	@git tag -a "$(MODULE_PATH)/$(VERSION)" -m "$(MODULE_PATH) $(VERSION)"
	@echo "$(GREEN)Tag creado: $(MODULE_PATH)/$(VERSION)$(NC)"

release-push-tag: guard-version ## Publicar tag del módulo en origin
	@git push origin "$(MODULE_PATH)/$(VERSION)"

release-github: guard-version ## Crear GitHub Release del módulo a partir del tag existente
	@$(RELEASE_SCRIPT) github $(MODULE_PATH) $(VERSION)
