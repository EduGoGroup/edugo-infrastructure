# 🟢 Propuestas de Refactorización

Mejoras opcionales que aumentarían la calidad del código sin ser urgentes.

---

## REF-001: Unificar CLIs de Migración

### Descripción

Los CLIs de migración para PostgreSQL y MongoDB tienen estructura casi idéntica. Se podrían unificar con una interfaz común.

### Estado Actual

```
postgres/cmd/migrate/migrate.go  (453 líneas)
mongodb/cmd/migrate/migrate.go   (522 líneas)
```

Funciones duplicadas:
- `printHelp()`
- `getEnv()`
- `sanitizeName()`
- `loadMigrations()`
- `showStatus()`
- Estructura `Migration`

### Propuesta

Crear paquete `internal/migration` con interfaz común:

```go
// internal/migration/migration.go
package migration

type Migration struct {
	Version   int
	Name      string
	UpScript  string
	DownScript string
	AppliedAt *time.Time
}

type Store interface {
	EnsureMigrationsTable(ctx context.Context) error
	GetAppliedMigrations(ctx context.Context) (map[int]*time.Time, error)
	RecordMigration(ctx context.Context, version int, name string) error
	RemoveMigration(ctx context.Context, version int) error
}

type Executor interface {
	ExecuteUp(ctx context.Context, sql string) error
	ExecuteDown(ctx context.Context, sql string) error
}

type Runner struct {
	store    Store
	executor Executor
	loader   MigrationLoader
}

func (r *Runner) Up(ctx context.Context) error { ... }
func (r *Runner) Down(ctx context.Context) error { ... }
func (r *Runner) Status(ctx context.Context) error { ... }
```

Implementaciones:
```go
// internal/migration/postgres/store.go
type PostgresStore struct { db *sql.DB }

// internal/migration/mongodb/store.go  
type MongoStore struct { db *mongo.Database }
```

### Beneficios

- Reduce duplicación ~40%
- Facilita agregar nuevas DBs (ej: SQLite para tests)
- Tests más fáciles con mocks
- Comportamiento consistente

### Riesgos

- Over-engineering para 2 implementaciones
- Complejidad adicional
- Breaking change en estructura de comandos

### Esfuerzo: 8-12 horas

### Recomendación: 🟡 Considerar si se agrega tercera DB

---

## REF-002: Entities con Métodos de Validación

### Descripción

Agregar métodos de validación a entities para usar antes de INSERT.

### Estado Actual

```go
type User struct {
	ID        uuid.UUID  `db:"id"`
	Email     string     `db:"email"`
	FirstName string     `db:"first_name"`
	// ...
}

func (User) TableName() string { return "users" }
```

### Propuesta

```go
type User struct {
	ID        uuid.UUID  `db:"id"`
	Email     string     `db:"email" validate:"required,email"`
	FirstName string     `db:"first_name" validate:"required,min=1,max=100"`
	LastName  string     `db:"last_name" validate:"required,min=1,max=100"`
	Role      string     `db:"role" validate:"required,oneof=admin teacher student guardian"`
	// ...
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// O con validación manual más específica
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if u.Role != "" && !isValidRole(u.Role) {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}
```

### Beneficios

- Validación consistente en todos los proyectos
- Errores claros antes de llegar a BD
- Documentación de reglas de negocio en código

### Riesgos

- Entities dejan de ser "structs puros"
- Posible conflicto con validadores de APIs
- Dependencia adicional (validator lib)

### Esfuerzo: 4-6 horas

### Recomendación: 🟡 Evaluar necesidad real

---

## REF-003: Builder Pattern para Queries Complejas

### Descripción

Crear builders para queries comunes con múltiples filtros.

### Propuesta

```go
// internal/query/builder.go
type MaterialQueryBuilder struct {
	schoolID       *uuid.UUID
	teacherID      *uuid.UUID
	status         *string
	subject        *string
	academicUnitID *uuid.UUID
	limit          int
	offset         int
}

func NewMaterialQuery() *MaterialQueryBuilder {
	return &MaterialQueryBuilder{limit: 20}
}

func (b *MaterialQueryBuilder) BySchool(id uuid.UUID) *MaterialQueryBuilder {
	b.schoolID = &id
	return b
}

func (b *MaterialQueryBuilder) ByTeacher(id uuid.UUID) *MaterialQueryBuilder {
	b.teacherID = &id
	return b
}

func (b *MaterialQueryBuilder) WithStatus(status string) *MaterialQueryBuilder {
	b.status = &status
	return b
}

func (b *MaterialQueryBuilder) Paginate(limit, offset int) *MaterialQueryBuilder {
	b.limit = limit
	b.offset = offset
	return b
}

func (b *MaterialQueryBuilder) Build() (string, []interface{}) {
	query := "SELECT * FROM materials WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if b.schoolID != nil {
		query += fmt.Sprintf(" AND school_id = $%d", argIndex)
		args = append(args, *b.schoolID)
		argIndex++
	}
	// ... más condiciones

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, b.limit, b.offset)

	return query, args
}

// Uso
query, args := NewMaterialQuery().
	BySchool(schoolID).
	WithStatus("ready").
	Paginate(20, 0).
	Build()

rows, err := db.Query(query, args...)
```

### Beneficios

- Queries type-safe
- Fácil de testear
- Evita SQL injection
- Código más legible

### Riesgos

- Abstracción que puede limitar queries complejas
- Overhead para queries simples

### Esfuerzo: 6-8 horas

### Recomendación: 🟢 Nice to have

---

## REF-004: Event Types con Generics

### Descripción

Usar generics de Go 1.18+ para tipos de eventos.

### Estado Actual

```go
type MaterialUploadedEvent struct {
	EventID      string
	EventType    string
	EventVersion string
	Timestamp    time.Time
	Payload      MaterialUploadedPayload
}

type AssessmentGeneratedEvent struct {
	EventID      string
	EventType    string
	EventVersion string
	Timestamp    time.Time
	Payload      AssessmentGeneratedPayload
}
```

### Propuesta

```go
// Base event con generic para payload
type Event[T any] struct {
	EventID      string    `json:"event_id"`
	EventType    string    `json:"event_type"`
	EventVersion string    `json:"event_version"`
	Timestamp    time.Time `json:"timestamp"`
	Payload      T         `json:"payload"`
}

// Payloads específicos
type MaterialUploadedPayload struct {
	MaterialID    string `json:"material_id"`
	SchoolID      string `json:"school_id"`
	TeacherID     string `json:"teacher_id"`
	FileURL       string `json:"file_url"`
	FileSizeBytes int64  `json:"file_size_bytes"`
	FileType      string `json:"file_type"`
}

// Type aliases para compatibilidad
type MaterialUploadedEvent = Event[MaterialUploadedPayload]
type AssessmentGeneratedEvent = Event[AssessmentGeneratedPayload]

// Constructor genérico
func NewEvent[T any](eventType, version string, payload T) Event[T] {
	return Event[T]{
		EventID:      uuid.New().String(),
		EventType:    eventType,
		EventVersion: version,
		Timestamp:    time.Now().UTC(),
		Payload:      payload,
	}
}

// Uso
event := NewEvent("material.uploaded", "1.0", MaterialUploadedPayload{
	MaterialID: materialID,
	// ...
})
```

### Beneficios

- Menos código duplicado
- Type safety en compile time
- Constructor unificado

### Riesgos

- Requiere Go 1.18+
- Cambio en API pública

### Esfuerzo: 3-4 horas

### Recomendación: 🟢 Considerar para v2.0

---

## REF-005: Separar Concerns en Migraciones

### Descripción

El CLI de migraciones mezcla:
- Parsing de argumentos
- Conexión a BD
- Lógica de migraciones
- Output formatting

### Propuesta

```
cmd/migrate/
├── main.go           # Solo parsing y wiring
├── commands/
│   ├── up.go
│   ├── down.go
│   ├── status.go
│   └── create.go
└── output/
    └── formatter.go  # Output formatting
```

### Beneficios

- Código más testeable
- Fácil agregar nuevos comandos
- Separación de concerns

### Esfuerzo: 4-6 horas

### Recomendación: 🟢 Si CLI crece

---

## 📊 Resumen de Propuestas

| ID | Propuesta | Beneficio | Esfuerzo | Prioridad |
|----|-----------|-----------|----------|-----------|
| REF-001 | Unificar CLIs | Alto | 8-12h | 🟡 Media |
| REF-002 | Validación en entities | Medio | 4-6h | 🟡 Media |
| REF-003 | Query builders | Medio | 6-8h | 🟢 Baja |
| REF-004 | Events con generics | Medio | 3-4h | 🟢 Baja |
| REF-005 | Separar CLI | Bajo | 4-6h | 🟢 Baja |

### Total si se implementan todas: 25-36 horas

---

## 🎯 Criterios de Decisión

Implementar refactorización cuando:

1. **El código se toca frecuentemente** - ROI de mejora es alto
2. **Hay bugs recurrentes** - Señal de código problemático
3. **Onboarding es difícil** - Código confuso
4. **Performance es problema** - Optimización necesaria

NO implementar cuando:

1. **Código estable** - "If it ain't broke, don't fix it"
2. **Único uso** - Over-engineering
3. **Deadline cercano** - Priorizar features
4. **Sin tests** - Refactorizar sin tests es peligroso

---

**Última actualización:** Diciembre 2024
