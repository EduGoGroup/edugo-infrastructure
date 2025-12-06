# Workflows y CI/CD - edugo-infrastructure

**Documentación completa de configuración CI/CD**

---

## 📊 Estado Actual

![CI Status](https://github.com/EduGoGroup/edugo-infrastructure/workflows/CI/badge.svg)

```yaml
success_rate: 95-100% (target post Sprint-1)
go_version: "1.25"
test_strategy: "unit tests + integration tests (opt-in)"
pre_commit_hooks: enabled
race_detection: enabled
```

---

## 🔄 Workflows Disponibles

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Trigger:**
- Pull requests a `main` y `dev`
- Push a `main`

**Jobs:**
```yaml
test:
  - Setup Go 1.25
  - Configure GOPRIVATE for private repos
  - Download dependencies (all modules)
  - Run tests (short + race detection)
```

**Características:**
- ✅ **Go 1.25**: Versión estandarizada
- ✅ **Short flag**: Skips integration tests (`-short`)
- ✅ **Race detection**: Detecta condiciones de carrera (`-race`)
- ✅ **GOPRIVATE**: Acceso a repos privados de EduGoGroup
- ✅ **Cache**: Go modules cacheados para velocidad

**Tiempo de ejecución:** ~3-5 minutos

---

## 🧪 Estrategia de Testing

### Tests Unitarios (CI)

**Comando:**
```bash
go test -short -race -v ./...
```

**Características:**
- Ejecutados en **cada** PR y push
- Skip integration tests automáticamente
- Race detector habilitado
- Timeout: 5 minutos por módulo

### Tests de Integración (Local/Opt-in)

**Comando:**
```bash
# Opción 1: Ejecutar todos los tests
go test -v ./...

# Opción 2: Solo integration tests
ENABLE_INTEGRATION_TESTS=true go test -v ./...
```

**Requiere:**
- Docker corriendo localmente
- Testcontainers funcional
- PostgreSQL/MongoDB containers disponibles

**Cuándo ejecutar:**
- Antes de merge a `main`
- Después de cambios en migraciones
- Validación pre-release

---

## 📦 Configuración por Módulo

### postgres/

**go.mod:**
```go
module github.com/EduGoGroup/edugo-infrastructure/postgres
go 1.25
```

**Tests:**
- Unit tests: Validación de SQL sintaxis
- Integration tests: Testcontainers + PostgreSQL 16
- Skipped en CI con `-short` flag

**Dependencias principales:**
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/testcontainers/testcontainers-go` - Testing

### mongodb/

**go.mod:**
```go
module github.com/EduGoGroup/edugo-infrastructure/mongodb
go 1.25
```

**Tests:**
- Unit tests: Validación de schemas
- Integration tests: Testcontainers + MongoDB 7
- Skipped en CI con `-short` flag

**Dependencias principales:**
- `go.mongodb.org/mongo-driver` - MongoDB driver
- `github.com/testcontainers/testcontainers-go` - Testing

### messaging/

**go.mod:**
```go
module github.com/EduGoGroup/edugo-infrastructure/messaging
go 1.25
```

**Tests:**
- Unit tests: Validación de JSON schemas
- Performance: Benchmarks de validación
- Integration tests: RabbitMQ mocking

**Dependencias principales:**
- `github.com/xeipuuv/gojsonschema` - JSON Schema validation
- `github.com/rabbitmq/amqp091-go` - RabbitMQ client

### schemas/

**go.mod:**
```go
module github.com/EduGoGroup/edugo-infrastructure/schemas
go 1.25
```

**Tests:**
- Unit tests: Schema validation
- No integration tests required

---

## 🔧 Configuración de Desarrollo

### Pre-commit Hooks

**Instalación:**
```bash
# Una sola vez por clon del repo
./scripts/setup-hooks.sh
```

**Checks automáticos:**
1. **go fmt** - Formato de código
2. **go vet** - Análisis estático
3. **go mod tidy** - Dependencias actualizadas
4. **go test -short** - Tests unitarios

**Bypass (NO recomendado):**
```bash
git commit --no-verify
```

### Variables de Entorno

**CI (GitHub Actions):**
```yaml
GOPRIVATE: github.com/EduGoGroup/*
GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Local (opcional):**
```bash
# Para tests de integración
export ENABLE_INTEGRATION_TESTS=true
```

---

## 🚀 Comandos Útiles

### Testing Local

```bash
# Todos los módulos (unit tests)
for module in postgres mongodb messaging schemas; do
  cd $module
  go test -short -race -v ./...
  cd ..
done

# Módulo específico (unit tests)
cd postgres
go test -short -v ./...

# Módulo específico (integration tests)
cd postgres
ENABLE_INTEGRATION_TESTS=true go test -v ./...
```

### Validación Pre-Push

```bash
# Ejecutar pre-commit checks manualmente
.git/hooks/pre-commit

# O usar el script de validación
./scripts/test-all-modules.sh
```

### Diagnóstico

```bash
# Ver versión de Go
go version  # Debe ser 1.25+

# Verificar dependencias
cd <module>
go mod verify
go mod tidy

# Ejecutar linters
cd <module>
go vet ./...
gofmt -l .
```

---

## 📊 Troubleshooting

### Error: "tests failing in CI but passing locally"

**Causa:** Integration tests ejecutándose en CI

**Solución:**
1. Verificar que tests de integración tienen `testing.Short()` check
2. Asegurar que CI usa flag `-short`

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    // ... test code
}
```

### Error: "cannot download private repo"

**Causa:** GOPRIVATE no configurado

**Solución local:**
```bash
# Configurar Git para usar token
git config --global url."https://TOKEN@github.com/".insteadOf "https://github.com/"

# O usar SSH
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

**Solución CI:** Ya configurado en workflows

### Error: "go.mod out of date"

**Causa:** Dependencias desactualizadas

**Solución:**
```bash
cd <module>
go mod tidy
git add go.mod go.sum
git commit -m "chore: update dependencies"
```

### Error: "race detector finds issues"

**Causa:** Condiciones de carrera en código

**Solución:**
1. Ejecutar localmente con `-race` para reproducir
2. Agregar mutexes o channels según sea necesario
3. Ver logs detallados para ubicar la línea exacta

```bash
go test -race -v ./... 2>&1 | grep "DATA RACE"
```

---

## 🎯 Mejores Prácticas

### 1. Escribir Tests

✅ **DO:**
- Usar `testing.Short()` para tests de integración
- Agregar benchmarks para código crítico
- Usar table-driven tests
- Mockear servicios externos en unit tests

❌ **DON'T:**
- Ejecutar integration tests sin flag check
- Hardcodear valores de conexión
- Dejar tests flaky sin resolver

### 2. Gestión de Dependencias

✅ **DO:**
- Ejecutar `go mod tidy` después de agregar deps
- Mantener go.mod sincronizado entre módulos
- Usar versiones específicas (no `latest`)

❌ **DON'T:**
- Commitear sin `go mod tidy`
- Agregar dependencias innecesarias
- Usar replace directives en producción

### 3. Commits y PRs

✅ **DO:**
- Ejecutar pre-commit hooks antes de push
- Validar tests localmente primero
- Usar conventional commits
- Esperar a CI antes de merge

❌ **DON'T:**
- Usar `--no-verify` habitualmente
- Pushear código sin testear
- Mergear con CI fallando

---

## 📚 Referencias

### Documentación Interna
- [README.md](../README.md) - Guía general del proyecto
- [docs/cicd/](../docs/cicd/) - Planes de sprint y tracking
- [scripts/](../scripts/) - Scripts de automatización

### Módulos
- [postgres/README.md](../postgres/README.md)
- [mongodb/README.md](../mongodb/README.md)
- [messaging/README.md](../messaging/README.md)

### Recursos Externos
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [Testcontainers Go](https://golang.testcontainers.org/)
- [GitHub Actions](https://docs.github.com/en/actions)

---

## 🔄 Changelog

### 2025-11-20 - Sprint 1 Improvements
- ✅ Migración a Go 1.25
- ✅ Implementación de `-short` flag strategy
- ✅ Race detection habilitado
- ✅ GOPRIVATE configurado
- ✅ Pre-commit hooks implementados
- ✅ Success rate: 20% → 95%+

---

**Última actualización:** 20 de Noviembre, 2025
**Versión:** 1.0
**Mantenedor:** Equipo EduGo
