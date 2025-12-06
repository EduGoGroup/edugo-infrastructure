# 🟠 Deuda Técnica - EduGo Infrastructure

Deuda técnica identificada que debe abordarse para mantener la salud del proyecto.

---

## TD-001: Módulos Go Sin Release Tags

### Descripción

Los módulos Go no tienen tags de versión publicados, lo que dificulta el versionado semántico.

### Estado Actual

```bash
# No existen tags como:
# postgres/v0.1.0
# mongodb/v0.1.0
# schemas/v0.1.0
# messaging/v0.1.0
```

### Problema

- Proyectos consumidores no pueden fijar versiones
- `go get` siempre trae `@latest` (HEAD de main)
- No hay changelog por versión
- Difícil hacer rollback a versión anterior

### Solución

```bash
# Crear primer release de cada módulo
cd postgres && git tag postgres/v0.1.0
cd mongodb && git tag mongodb/v0.1.0
cd schemas && git tag schemas/v0.1.0
cd messaging && git tag messaging/v0.1.0

git push origin --tags
```

### Esfuerzo: 1 hora

---

## TD-002: Sin CI/CD Configurado

### Descripción

No hay GitHub Actions configurados para:
- Ejecutar tests automáticamente
- Lint del código
- Validar migraciones
- Publicar releases

### Estado Actual

```
.github/
└── (vacío o sin workflows relevantes)
```

### Problema

- PRs se mergean sin verificación automática
- Bugs pueden introducirse sin detectarse
- No hay garantía de que el código compile

### Solución Propuesta

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4

  test-postgres:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: cd postgres && go test ./...

  test-mongodb:
    runs-on: ubuntu-latest
    services:
      mongodb:
        image: mongo:7.0
        ports:
          - 27017:27017
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: cd mongodb && go test ./...

  test-schemas:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: cd schemas && go test ./...
```

### Esfuerzo: 4-6 horas

---

## TD-003: Falta de Error Wrapping Consistente

### Descripción

Algunos errores no se envuelven correctamente con contexto.

### Ejemplo Problemático

```go
// postgres/cmd/migrate/migrate.go:164
if _, err := tx.Exec(m.UpSQL); err != nil {
	_ = tx.Rollback()
	return fmt.Errorf("error en migración %d: %w", m.Version, err)
}
```

Bien ✅ - Tiene contexto

```go
// Otros lugares
if err != nil {
	return err  // ❌ Sin contexto
}
```

### Problema

- Difícil rastrear origen de errores
- Logs no informativos
- Debugging más lento

### Solución

Revisar todos los `return err` y agregar contexto:

```go
// Antes
if err != nil {
	return err
}

// Después
if err != nil {
	return fmt.Errorf("failed to connect to database: %w", err)
}
```

### Esfuerzo: 2-3 horas

---

## TD-004: Hardcoded Timeouts

### Descripción

Timeouts están hardcodeados en lugar de ser configurables.

### Ejemplos

```go
// mongodb/cmd/migrate/migrate.go:40
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// mongodb/cmd/migrate/migrate.go:138
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

// mongodb/cmd/migrate/migrate.go:497
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
```

### Problema

- No se pueden ajustar para diferentes ambientes
- Migraciones largas pueden fallar por timeout
- Ambientes lentos (CI) pueden tener problemas

### Solución

```go
var (
	DefaultConnectTimeout   = getEnvDuration("MIGRATE_CONNECT_TIMEOUT", 10*time.Second)
	DefaultOperationTimeout = getEnvDuration("MIGRATE_OPERATION_TIMEOUT", 5*time.Second)
	DefaultMigrationTimeout = getEnvDuration("MIGRATE_MIGRATION_TIMEOUT", 2*time.Minute)
)

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
```

### Esfuerzo: 1-2 horas

---

## TD-005: Logs con fmt.Printf en lugar de Logger

### Descripción

Los CLIs usan `fmt.Printf` para output en lugar de un logger estructurado.

### Ejemplos

```go
// postgres/cmd/migrate/migrate.go
fmt.Printf("Ejecutando migración %03d: %s\n", m.Version, m.Name)
fmt.Printf("✅ Migración %03d aplicada exitosamente\n", m.Version)
```

### Problema

- No hay niveles de log (debug, info, error)
- No hay timestamps
- No es JSON parseable para sistemas de monitoreo
- Emojis pueden causar problemas en algunos terminales

### Solución

```go
import "log/slog"

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// En lugar de fmt.Printf
logger.Info("executing migration", 
	"version", m.Version, 
	"name", m.Name)

logger.Info("migration applied successfully", 
	"version", m.Version)
```

### Esfuerzo: 3-4 horas

---

## TD-006: Sin Métricas de Procesamiento

### Descripción

No hay instrumentación para medir tiempos de operaciones.

### Problema

- No se sabe cuánto tardan las migraciones
- No hay baseline de rendimiento
- Difícil detectar regresiones de performance

### Solución

```go
import "time"

func migrateUp(db *sql.DB) error {
	start := time.Now()
	defer func() {
		logger.Info("migrate up completed", 
			"duration_ms", time.Since(start).Milliseconds(),
			"migrations_applied", pendingCount)
	}()
	
	// ... código existente
}
```

### Esfuerzo: 2 horas

---

## 📊 Resumen de Deuda Técnica

| ID | Descripción | Prioridad | Esfuerzo | Impacto |
|----|-------------|-----------|----------|---------|
| TD-001 | Sin release tags | 🔴 Alta | 1h | Versionado |
| TD-002 | Sin CI/CD | 🔴 Alta | 4-6h | Calidad |
| TD-003 | Error wrapping | 🟡 Media | 2-3h | Debugging |
| TD-004 | Hardcoded timeouts | 🟡 Media | 1-2h | Flexibilidad |
| TD-005 | Printf vs Logger | 🟢 Baja | 3-4h | Observabilidad |
| TD-006 | Sin métricas | 🟢 Baja | 2h | Observabilidad |

### Total Estimado: 13-18 horas

---

## 📈 Plan de Reducción

### Sprint 1 (Urgente)
- [ ] TD-001: Crear release tags
- [ ] TD-002: Configurar CI básico

### Sprint 2 (Importante)
- [ ] TD-003: Error wrapping
- [ ] TD-004: Timeouts configurables

### Sprint 3 (Nice to Have)
- [ ] TD-005: Logger estructurado
- [ ] TD-006: Métricas

---

**Última actualización:** Diciembre 2024
