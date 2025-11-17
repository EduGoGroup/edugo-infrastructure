.PHONY: help test coverage benchmark

help:
	@echo "Makefile - Messaging Module"
	@echo ""
	@echo "Comandos disponibles:"
	@echo "  make test        Ejecutar tests"
	@echo "  make coverage    Ver cobertura de tests"
	@echo "  make benchmark   Ejecutar benchmarks"

test:
	@echo "Ejecutando tests..."
	go test -v ./...

coverage:
	@echo "Generando reporte de cobertura..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

benchmark:
	@echo "Ejecutando benchmarks..."
	go test -bench=. -benchmem ./...
