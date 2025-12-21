# 🟠 Deuda Técnica - EduGo Infrastructure

Deuda técnica identificada que debe abordarse para mantener la salud del proyecto.

---

## ~~TD-001: Módulos Go Sin Release Tags~~ ✅ RESUELTO

### Estado: ✅ **RESUELTO** (2025-12-20)

**Acción tomada:** El proyecto ya cuenta con tags de versión por módulo siguiendo el patrón `<módulo>/v<SemVer>`.

### Estado Actual ✅

```bash
# Tags existentes verificados:
postgres/v0.11.1
mongodb/v0.10.1
schemas/v0.1.2
# (y otros módulos versionados)
```

### Contexto Histórico

Inicialmente se identificó como deuda técnica la falta de release tags. Sin embargo, al investigar se descubrió que:

1. **El proyecto ya tiene versionado por módulo** desde hace tiempo
2. **Existen integraciones activas** que consumen versiones específicas
3. **El sistema de tags está funcionando correctamente**

### Documentación Creada

Se ha creado la guía completa de releases en **`documents/RELEASING.md`** que documenta:

- ✅ Visión general del sistema de versionado
- ✅ Estructura de tags: `<módulo>/v<SemVer>`
- ✅ Proceso paso a paso para crear nuevos releases
- ✅ Comandos útiles (listar, crear, eliminar tags)
- ✅ Ejemplos prácticos por módulo
- ✅ Troubleshooting

### Solución Aplicada

```bash
# El sistema ya funciona con este patrón:
git tag postgres/v0.11.2
git tag mongodb/v0.10.2
git push origin --tags

# Consumidores pueden usar versiones específicas:
go get github.com/edugo/edugo-infrastructure/postgres@v0.11.1
```

### Resuelto: Diciembre 2025
### Esfuerzo Real: 1 hora (documentación)

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

## ~~TD-005: Logs con fmt.Printf en lugar de Logger~~ ✅ RESUELTO

### Estado: ✅ **RESUELTO** (2025-12-20)

**Acción tomada:** Migración completa de `fmt.Printf` a `log/slog` en ambos CLIs de migraciones.

### Cambios Implementados

**1. PostgreSQL CLI (`postgres/cmd/migrate/migrate.go`):**
- ✅ Reemplazado `import "log"` por `import "log/slog"`
- ✅ Agregado logger global con `slog.NewTextHandler`
- ✅ ~45 llamadas migradas de `fmt.Printf` a `logger.Info/Warn/Error`
- ✅ `log.Fatalf` → `logger.Error` + `os.Exit(1)`
- ✅ Preservados outputs user-facing (`printHelp`, `showStatus`, `createMigration`)

**2. MongoDB CLI (`mongodb/cmd/migrate/migrate.go`):**
- ✅ Mismo patrón de migración que PostgreSQL
- ✅ ~20 llamadas migradas a logger estructurado
- ✅ Agregado `import "strconv"` para conversión de versiones

**3. Estructura del Logger:**
```go
var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
```

### Ejemplos de Cambios

**ANTES:**
```go
fmt.Printf("Ejecutando migración %03d: %s\n", m.Version, m.Name)
log.Fatalf("Error ejecutando migraciones: %v", err)
```

**DESPUÉS:**
```go
logger.Info("ejecutando migración", "version", m.Version, "name", m.Name)
logger.Error("error ejecutando migraciones", "error", err)
os.Exit(1)
```

### Beneficios Obtenidos

- ✅ Logs estructurados con campos parseables
- ✅ Timestamps automáticos en cada entrada
- ✅ Niveles de log (Info, Warn, Error)
- ✅ Compatible con herramientas de agregación (Splunk, DataDog, etc.)
- ✅ Sin dependencias externas (usa Go stdlib)

### Esfuerzo: ✅ COMPLETADO (3-4 horas)

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
| ~~TD-001~~ | ~~Sin release tags~~ | ✅ Resuelto | 1h | Versionado |
| TD-002 | Sin CI/CD | 🔴 Alta | 4-6h | Calidad |
| TD-003 | Error wrapping | 🟡 Media | 2-3h | Debugging |
| TD-004 | Hardcoded timeouts | 🟡 Media | 1-2h | Flexibilidad |
| ~~TD-005~~ | ~~Printf vs Logger~~ | ✅ Resuelto | 3-4h | Observabilidad |
| TD-006 | Sin métricas | 🟢 Baja | 2h | Observabilidad |

### Total Estimado: 13-18 horas

---

## 📈 Plan de Reducción

### Sprint 1 (Urgente)
- [x] TD-001: Crear release tags ✅ RESUELTO
- [ ] TD-002: Configurar CI básico

### Sprint 2 (Importante)
- [ ] TD-003: Error wrapping
- [ ] TD-004: Timeouts configurables

### Sprint 3 (Nice to Have)
- [x] TD-005: Logger estructurado ✅ RESUELTO
- [ ] TD-006: Métricas

---

**Última actualización:** Diciembre 2025
