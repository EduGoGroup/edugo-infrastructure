.PHONY: help migrate-up migrate-down migrate-status migrate-create seed test

help:
	@echo "Makefile - MongoDB Module"
	@echo ""
	@echo "Comandos disponibles:"
	@echo "  make migrate-up       Ejecutar migraciones pendientes"
	@echo "  make migrate-down     Revertir última migración"
	@echo "  make migrate-status   Ver estado de migraciones"
	@echo "  make migrate-create   Crear nueva migración"
	@echo "  make seed             Cargar datos de prueba"
	@echo "  make test             Ejecutar tests"

migrate-up:
	@echo "Ejecutando migraciones MongoDB..."
	go run migrate.go up

migrate-down:
	@echo "Revirtiendo última migración..."
	go run migrate.go down

migrate-status:
	@echo "Estado de migraciones:"
	go run migrate.go status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Debe especificar name=descripcion"; \
		exit 1; \
	fi
	go run migrate.go create "$(name)"

seed:
	@echo "Cargando seeds MongoDB..."
	@if [ -d "seeds" ]; then \
		for file in seeds/*.js; do \
			echo "Ejecutando $$file..."; \
			mongosh "$(MONGO_URI)/$(MONGO_DB_NAME)" --file $$file --quiet; \
		done; \
	else \
		echo "No se encontró el directorio seeds/"; \
	fi

test:
	@echo "Ejecutando tests..."
	go test -v ./...
