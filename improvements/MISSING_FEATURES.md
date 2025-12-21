# 🟡 Funcionalidades Incompletas - TODOs Pendientes

Funcionalidades marcadas como TODO que requieren implementación.

---

## ~~TODO-001: ApplySeeds() No Implementada~~ ✅ RESUELTO

### Estado: ✅ **RESUELTO** (2025-12-20)

### Ubicación Original

```
mongodb/migrations/embed.go:95-103
```

### Problema Original

- Función pública que no hacía nada (retornaba nil)
- Confundía a usuarios del módulo
- Seeds existían en `mongodb/seeds/` (9 archivos JavaScript) pero no se cargaban

### Solución Implementada

**Archivos creados:**
- `mongodb/migrations/seeds.go` (1,053 líneas) - Contiene todos los datos de seeds en Go

**Archivos modificados:**
- `mongodb/migrations/embed.go` - Función `ApplySeeds()` ahora invoca `applySeedsInternal()`
- `mongodb/migrations/migrations_integration_test.go` - Agregado test `testApplySeeds()`

### Implementación

**1. Conversión JavaScript → Go:**
Los 9 archivos JavaScript fueron convertidos a estructuras Go usando `bson.D` y `bson.A`:

```go
func analyticsEventsSeeds() seedDocument {
    return seedDocument{
        collection: "analytics_events",
        documents: []interface{}{
            bson.D{
                {Key: "event_name", Value: "page.view"},
                {Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
                // ... 6 eventos completos
            },
        },
    }
}
```

**2. Función principal:**
```go
func ApplySeeds(ctx context.Context, db *mongo.Database) error {
    inserted, err := applySeedsInternal(ctx, db)
    if err != nil {
        return fmt.Errorf("error aplicando seeds: %w", err)
    }
    return nil
}
```

**3. Idempotencia:**
- Usa `InsertMany` con `ordered: false`
- Ignora errores de clave duplicada
- Permite ejecutar múltiples veces sin duplicar datos (para colecciones con `_id` explícito)

### Collections Pobladas

| Collection | Documentos | Descripción |
|------------|-----------|-------------|
| `analytics_events` | 6 | Eventos de ejemplo (page.view, material.view, assessment.start/complete, search.performed) |
| `material_assessment` | 2 | Assessments de Física y Matemáticas con ObjectID explícito |
| `audit_logs` | 5 | Registros de auditoría (login, material uploaded, failed login, system backup) |
| `material_assessment_worker` | 2 | Workers con preguntas generadas por IA (POO Java, React Hooks) |
| `material_summary` | 3 | Resúmenes en español, inglés y portugués |
| `notifications` | 4 | Notificaciones de ejemplo (assessment ready/graded, material uploaded, system announcement) |

**Total:** 22 documentos de ejemplo

### Tests Agregados

```go
func testApplySeeds(ctx context.Context, db *mongo.Database) func(*testing.T) {
    // 1. Aplica seeds
    // 2. Verifica conteo de documentos por colección
    // 3. Test de idempotencia (ejecuta seeds 2 veces)
    // 4. Verifica que NO se duplican documentos con _id explícito
}
```

### Beneficios

- ✅ **Type-safe**: Go verifica tipos en tiempo de compilación
- ✅ **Sin dependencias externas**: No necesita intérprete JavaScript
- ✅ **Idempotente**: Se puede ejecutar múltiples veces
- ✅ **Testeable**: Tests de integración incluidos
- ✅ **Consistente**: Sigue el patrón de PostgreSQL
- ✅ **Documentado**: GoDoc completo con ejemplos

### Uso

```go
import "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"

// Inicializar base de datos completa
migrations.ApplyAll(ctx, db)
migrations.ApplySeeds(ctx, db)  // ← Ahora funcional
```

### Impacto en Proyectos

- **edugo-worker**: Ahora puede usar `ApplySeeds()` en tests de integración
- **edugo-api-mobile**: Consistencia con el patrón ya usado en PostgreSQL

### Archivos JavaScript Originales

Los archivos en `mongodb/seeds/*.js` se mantienen como **documentación de referencia** pero ya no son necesarios para la ejecución. La implementación en Go es la fuente de verdad.

### Esfuerzo Real

- **Complejidad:** Media
- **Tiempo:** ~2 horas (conversión manual de JavaScript a Go)
- **Líneas agregadas:** +1,053 líneas en seeds.go
- **Líneas modificadas:** ~30 líneas en embed.go + test

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
- No existe directorio `testing/` con archivos de prueba

### Comparación con PostgreSQL

PostgreSQL SÍ tiene implementación funcional con 10 archivos en `testing/`:
- `001_demo_users.sql`
- `002_demo_schools.sql`
- `003_demo_academic_units.sql`
- etc.

### Solución Propuesta

1. Crear directorio `mongodb/testing/`
2. Agregar archivos `.js` con datos mock
3. Implementar carga similar a PostgreSQL

### Esfuerzo Estimado

- **Complejidad:** Media
- **Tiempo:** 2-4 horas

---

## TODO-004: Tests de Integración MongoDB

### Ubicación

```
mongodb/migrations/migrations_integration_test.go
```

### Estado Actual

🟡 **Parcialmente implementado**

✅ **Existe y funciona:**
- Archivo `migrations_integration_test.go` creado
- 5 tests implementados:
  - `TestIntegration` - Suite principal
  - `testApplyAll` - Verifica aplicación de migraciones
  - `testCRUDMaterialAssessment` - Prueba CRUD completo
  - `testCRUDNotifications` - Prueba CRUD de notificaciones
  - `testIndexesValidation` - Verifica creación de índices

❌ **Faltante:**
- Tests para `ApplySeeds()` (depende de TODO-001)
- Tests para `ApplyMockData()` (depende de TODO-002)
- Directorio `testing/` con archivos de prueba

### Conclusión

El framework de tests existe y funciona, pero está incompleto porque depende de funcionalidades no implementadas.

### Esfuerzo Estimado

- **Complejidad:** Baja (ya existe base)
- **Tiempo:** 2-3 horas (cuando TODO-001/002 estén listos)

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

### Schemas Actuales

- `assessment.generated:1.0`
- `material.deleted:1.0`
- `material.uploaded:1.0`
- `student.enrolled:1.0`

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

| ID | Descripción | Prioridad | Estado | Esfuerzo |
|----|-------------|-----------|--------|----------|
| TODO-001 | ApplySeeds() vacía | 🟡 Media | Pendiente | 2-4h |
| TODO-002 | ApplyMockData() vacía | 🟡 Media | Pendiente | 2-4h |
| TODO-004 | Tests MongoDB | 🟠 Media-Alta | Parcial | 2-3h |
| TODO-005 | Validación schemas | 🟢 Baja | Pendiente | 1h |

### Total Estimado: 7-14 horas

---

**Última actualización:** Diciembre 2024
