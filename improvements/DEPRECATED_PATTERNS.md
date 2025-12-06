# 🟡 Patrones Deprecados y Malas Prácticas

Patrones de código que deben evitarse o reemplazarse.

---

## DEP-001: Ignorar Errores con Blank Identifier

### Descripción

Uso de `_ = err` para ignorar errores silenciosamente.

### Ubicaciones

```go
// postgres/cmd/migrate/migrate.go:42
defer func() { _ = db.Close() }()

// postgres/cmd/migrate/migrate.go:163
_ = tx.Rollback()

// postgres/cmd/migrate/migrate.go:424
defer func() { _ = rows.Close() }()
```

### Por Qué es Problemático

- Errores se pierden silenciosamente
- Difícil debugging cuando algo falla
- Viola principio de "handle every error"

### Cuándo es Aceptable

- En `defer` para cleanup donde el error no afecta el flujo
- Cuando ya se está manejando otro error más importante

### Patrón Recomendado

```go
// Opción 1: Log el error aunque no lo propagues
defer func() {
	if err := db.Close(); err != nil {
		logger.Warn("failed to close db connection", "error", err)
	}
}()

// Opción 2: Si realmente no importa, documentar por qué
defer func() {
	// Error ignorado intencionalmente: ya hay otro error siendo propagado
	_ = tx.Rollback()
}()
```

### Severidad: 🟡 Media

---

## DEP-002: Context con Background en Funciones

### Descripción

Crear `context.Background()` dentro de funciones en lugar de recibirlo como parámetro.

### Ubicaciones

```go
// mongodb/cmd/migrate/migrate.go:138
func ensureMigrationsCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// ...
}

// mongodb/cmd/migrate/migrate.go:191
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

### Por Qué es Problemático

- No permite cancelación desde el caller
- No propaga deadlines
- No permite pasar valores via context

### Patrón Recomendado

```go
// Antes
func ensureMigrationsCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// ...
}

// Después
func ensureMigrationsCollection(ctx context.Context, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// ...
}

// Uso
ctx := context.Background() // Solo en main()
ensureMigrationsCollection(ctx, db)
```

### Severidad: 🟡 Media

---

## DEP-003: Panic en Código de Librería

### Descripción

Uso de `log.Fatal` y `panic` que terminan el programa abruptamente.

### Ubicaciones

```go
// postgres/cmd/migrate/migrate.go:40
log.Fatalf("Error conectando a PostgreSQL: %v", err)

// postgres/cmd/migrate/migrate.go:45
log.Fatalf("Error validando conexión: %v", err)
```

### Por Qué es Problemático

- `log.Fatal` llama `os.Exit(1)` sin ejecutar defers
- No permite que el caller maneje el error
- No permite cleanup graceful

### Cuándo es Aceptable

- En `main()` de un CLI
- Errores verdaderamente irrecuperables

### Nota

En este caso, como es código de CLI en `main()`, es aceptable. Sin embargo, si este código se refactoriza a librería, debe cambiarse.

### Severidad: 🟢 Baja (en contexto de CLI)

---

## DEP-004: SQL Concatenation sin Parameterización

### Descripción

Construcción de queries SQL con `fmt.Sprintf` en lugar de parámetros.

### Ubicación

```go
// postgres/cmd/migrate/migrate.go:126
query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		version INTEGER PRIMARY KEY,
		...
	)
`, migrationsTable)
```

### Análisis

En este caso específico:
- `migrationsTable` es una constante, no input del usuario
- No hay riesgo de SQL injection
- Es patrón común para nombres de tablas dinámicas

### Cuándo es Problemático

```go
// ❌ MALO: Input de usuario en query
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userName)

// ✅ BIEN: Usar parámetros
query := "SELECT * FROM users WHERE name = $1"
rows, err := db.Query(query, userName)
```

### Severidad: 🟢 Baja (en este contexto)

---

## DEP-005: Defer en Loop

### Descripción

Uso de `defer` dentro de loops puede causar memory leaks.

### Ejemplo (Hipotético)

```go
// ❌ MALO
for _, file := range files {
	f, _ := os.Open(file)
	defer f.Close() // Se acumulan hasta que la función termina
}

// ✅ BIEN
for _, file := range files {
	func() {
		f, _ := os.Open(file)
		defer f.Close()
		// usar f
	}()
}

// ✅ MEJOR
for _, file := range files {
	f, _ := os.Open(file)
	// usar f
	f.Close() // Cerrar explícitamente
}
```

### Estado en Codebase

No se encontró este patrón problemático en el código actual. ✅

### Severidad: N/A

---

## DEP-006: Magic Numbers

### Descripción

Números sin nombre que dificultan entender el código.

### Ubicaciones

```go
// mongodb/cmd/migrate/migrate.go:40
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//                                                       ^^ magic number

// mongodb/cmd/migrate/migrate.go:497
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
//                                                       ^^ magic number
```

### Patrón Recomendado

```go
const (
	DefaultConnectTimeout   = 10 * time.Second
	DefaultOperationTimeout = 5 * time.Second
	DefaultMigrationTimeout = 2 * time.Minute
)

ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
```

### Severidad: 🟢 Baja

---

## 📊 Resumen

| ID | Patrón | Severidad | Acción |
|----|--------|-----------|--------|
| DEP-001 | Ignorar errores | 🟡 Media | Documentar o loggear |
| DEP-002 | Context.Background() | 🟡 Media | Refactorizar si se extrae a lib |
| DEP-003 | log.Fatal | 🟢 Baja | OK en CLI |
| DEP-004 | SQL concat | 🟢 Baja | OK con constantes |
| DEP-005 | Defer en loop | ✅ OK | No presente |
| DEP-006 | Magic numbers | 🟢 Baja | Extraer constantes |

---

## 📝 Guía de Estilo Recomendada

### Manejo de Errores

```go
// ✅ Siempre verificar errores
if err != nil {
	return fmt.Errorf("operación falló: %w", err)
}

// ✅ Si ignoras un error, documenta por qué
_ = cleanup() // Error ignorado: cleanup best-effort
```

### Context

```go
// ✅ Pasar context como primer parámetro
func DoSomething(ctx context.Context, args ...) error

// ✅ Solo crear context.Background() en main()
func main() {
	ctx := context.Background()
	// ...
}
```

### Constantes

```go
// ✅ Nombrar valores mágicos
const (
	MaxRetries = 3
	DefaultTimeout = 30 * time.Second
	BufferSize = 4096
)
```

---

**Última actualización:** Diciembre 2024
