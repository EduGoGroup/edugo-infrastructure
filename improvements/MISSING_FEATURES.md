# 🟡 Funcionalidades Incompletas - TODOs Pendientes

Funcionalidades marcadas como TODO que requieren implementación.

---

## TODO-001: ApplySeeds() No Implementada

### Ubicación

```
mongodb/migrations/embed.go:100-103
```

### Código Actual

```go
// ApplySeeds ejecuta seeds (datos iniciales del ecosistema)
// Por ahora no implementado - agregar cuando se definan seeds necesarios
//
// Uso típico: Inicializar datos mínimos en ambiente de producción/staging
func ApplySeeds(ctx context.Context, db *mongo.Database) error {
	// TODO: Implementar cuando se definan seeds
	return nil
}
```

### Problema

- Función pública que no hace nada
- Puede confundir a usuarios del módulo
- Seeds existen en `mongodb/seeds/` pero no se cargan

### Solución Propuesta

```go
func ApplySeeds(ctx context.Context, db *mongo.Database) error {
	seedFiles := []struct {
		collection string
		filename   string
	}{
		{"material_assessment_worker", "material_assessment_worker.js"},
		{"material_summary", "material_summary.js"},
		{"material_event", "material_event.js"},
	}

	for _, sf := range seedFiles {
		content, err := seedsFS.ReadFile("seeds/" + sf.filename)
		if err != nil {
			return fmt.Errorf("error reading seed %s: %w", sf.filename, err)
		}
		
		if err := executeSeedScript(ctx, db, sf.collection, string(content)); err != nil {
			return fmt.Errorf("error applying seed %s: %w", sf.filename, err)
		}
	}
	
	return nil
}
```

### Esfuerzo Estimado

- **Complejidad:** Media
- **Tiempo:** 2-4 horas
- **Dependencias:** Definir formato de seeds

---

## TODO-002: ApplyMockData() No Implementada

### Ubicación

```
mongodb/migrations/embed.go:105-112
```

### Código Actual

```go
// ApplyMockData ejecuta datos mock para testing
// Por ahora no implementado - agregar cuando se definan datos de prueba
//
// Uso típico: Tests de integración, ambiente de desarrollo
func ApplyMockData(ctx context.Context, db *mongo.Database) error {
	// TODO: Implementar cuando se definan datos mock
	return nil
}
```

### Problema

- Similar a TODO-001
- Tests de integración no tienen datos mock centralizados

### Solución Propuesta

Implementar carga de datos mock desde archivos JSON/JS en directorio `testing/`.

### Esfuerzo Estimado

- **Complejidad:** Media
- **Tiempo:** 2-4 horas

---

## TODO-003: Entities Sin Migraciones SQL

### Descripción

Existen 6 entities Go definidas pero cuyas migraciones SQL no están activas o son incompletas.

### Entities Afectadas

| Entity | Archivo | Tabla Esperada | Estado |
|--------|---------|----------------|--------|
| `MaterialVersion` | `postgres/entities/material_version.go` | `material_versions` | Migración existe (012) |
| `Subject` | `postgres/entities/subject.go` | `subjects` | Migración existe (013) |
| `Unit` | `postgres/entities/unit.go` | `units` | Migración existe (014) |
| `GuardianRelation` | `postgres/entities/guardian_relation.go` | `guardian_relations` | Migración existe (015) |
| `Progress` | `postgres/entities/progress.go` | `progress` | Migración existe (016) |

### Código Actual

```go
// postgres/entities/material_version.go
type MaterialVersion struct {
	ID           uuid.UUID  `db:"id"`
	MaterialID   uuid.UUID  `db:"material_id"`
	VersionNumber int       `db:"version_number"`
	FileURL      string     `db:"file_url"`
	// ...
}
```

### Problema

- README indica que están "bloqueadas" pero migraciones existen
- Posible desincronización entre documentación y código
- Usuarios no saben qué entities pueden usar

### Solución Propuesta

1. Verificar que migraciones 012-016 funcionan correctamente
2. Actualizar README de entities indicando que están disponibles
3. Agregar tests de integración

### Esfuerzo Estimado

- **Complejidad:** Baja-Media
- **Tiempo:** 2-3 horas (verificación y documentación)

---

## TODO-004: Tests de Integración MongoDB

### Ubicación

```
mongodb/testing/  (directorio vacío)
```

### Problema

- Directorio `testing/` existe pero está vacío
- No hay tests de integración para entities MongoDB
- No hay tests para migraciones MongoDB

### Solución Propuesta

Crear tests similares a `postgres/migrations/migrations_integration_test.go`:

```go
// mongodb/migrations/migrations_integration_test.go
package migrations_test

import (
	"context"
	"testing"
	
	"go.mongodb.org/mongo-driver/mongo"
)

func TestApplyAllMigrations(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	
	// Test migrations
	if err := migrations.ApplyAll(ctx, db); err != nil {
		t.Fatalf("ApplyAll failed: %v", err)
	}
	
	// Verify collections exist
	collections, _ := db.ListCollectionNames(ctx, bson.M{})
	expected := []string{
		"material_assessment_worker",
		"material_summary",
		"material_event",
	}
	// Assert collections exist
}
```

### Esfuerzo Estimado

- **Complejidad:** Media
- **Tiempo:** 4-6 horas

---

## TODO-005: Validación de Schemas en Runtime

### Descripción

Los JSON Schemas se cargan al inicializar el validador, pero no hay validación de que todos los schemas esperados existan.

### Código Actual

```go
// schemas/validator.go
func NewEventValidator() (*EventValidator, error) {
	// ... carga schemas dinámicamente
	// No valida que todos los esperados existan
}
```

### Problema

- Si falta un schema, el error ocurre al validar (runtime)
- No hay lista definida de schemas requeridos
- Difícil detectar schemas faltantes en CI

### Solución Propuesta

```go
var RequiredSchemas = []string{
	"material.uploaded:1.0",
	"assessment.generated:1.0",
	"material.deleted:1.0",
	"student.enrolled:1.0",
}

func NewEventValidator() (*EventValidator, error) {
	v := &EventValidator{schemas: make(map[string]*gojsonschema.Schema)}
	
	// ... cargar schemas
	
	// Validar que todos los requeridos estén cargados
	for _, required := range RequiredSchemas {
		if _, exists := v.schemas[required]; !exists {
			return nil, fmt.Errorf("required schema missing: %s", required)
		}
	}
	
	return v, nil
}
```

### Esfuerzo Estimado

- **Complejidad:** Baja
- **Tiempo:** 1 hora

---

## 📊 Resumen de TODOs

| ID | Descripción | Prioridad | Esfuerzo |
|----|-------------|-----------|----------|
| TODO-001 | ApplySeeds() vacía | 🟡 Media | 2-4h |
| TODO-002 | ApplyMockData() vacía | 🟡 Media | 2-4h |
| TODO-003 | Entities sin doc actualizada | 🟡 Media | 2-3h |
| TODO-004 | Tests MongoDB faltantes | 🟠 Media-Alta | 4-6h |
| TODO-005 | Validación schemas | 🟢 Baja | 1h |

### Total Estimado: 11-18 horas

---

## ✅ Completados

| Fecha | ID | Descripción | PR |
|-------|-----|-------------|-----|
| - | - | - | - |

---

**Última actualización:** Diciembre 2024
