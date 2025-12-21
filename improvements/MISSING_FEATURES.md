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

## ~~TODO-002: ApplyMockData() No Implementada~~ ✅ RESUELTO

### Estado: ✅ **RESUELTO** (2025-12-20)

### Ubicación Original

```
mongodb/migrations/embed.go:105-112
```

### Problema Original

- Función pública que no hacía nada (retornaba nil)
- Tests de integración no tenían datos mock centralizados
- Similar a TODO-001 pero con más variedad de datos

### Solución Implementada

**Archivos creados:**
- `mongodb/migrations/mock_data.go` (1,089 líneas) - Contiene todos los datos mock en Go

**Archivos modificados:**
- `mongodb/migrations/embed.go` - Función `ApplyMockData()` ahora invoca `applyMockDataInternal()`
- `mongodb/migrations/migrations_integration_test.go` - Agregado test `testApplyMockData()`

### Implementación

**1. Estructura similar a seeds.go:**
```go
func analyticsEventsMockData() mockDocument {
    return mockDocument{
        collection: "analytics_events",
        documents: []interface{}{
            // 10 eventos con diferentes plataformas, países, roles
            bson.D{
                {Key: "event_name", Value: "user.login"},
                {Key: "device", Value: bson.D{
                    {Key: "platform", Value: "mobile"},
                    {Key: "os", Value: "iOS"},
                    // ... más variedad
                }},
            },
        },
    }
}
```

**2. Función principal:**
```go
func ApplyMockData(ctx context.Context, db *mongo.Database) error {
    inserted, err := applyMockDataInternal(ctx, db)
    if err != nil {
        return fmt.Errorf("error aplicando mock data: %w", err)
    }
    return nil
}
```

**3. Idempotencia:**
- Igual que ApplySeeds(), usa `InsertMany` con `ordered: false`
- Ignora errores de clave duplicada
- Permite ejecutar múltiples veces

### Collections Pobladas

| Collection | Documentos | Descripción |
|------------|-----------|-------------|
| `analytics_events` | 10 | Eventos variados (mobile/tablet/web, diferentes países y plataformas) |
| `material_assessment` | 3 | Assessments de Química (hard), Historia (easy), Cálculo (medium) |
| `audit_logs` | 8 | Registros extendidos (material deleted, user created, password changed, brute force, etc.) |
| `material_assessment_worker` | 3 | Workers en español, inglés y portugués con diferentes subjects |
| `material_summary` | 5 | Resúmenes en español, inglés, portugués, francés y alemán |
| `notifications` | 6 | Notificaciones variadas (material ready, system update, deadline, comment, achievement, security alert) |

**Total:** 35 documentos mock

### Diferencias vs ApplySeeds()

| Aspecto | Seeds (22 docs) | MockData (35 docs) |
|---------|----------------|-------------------|
| **Propósito** | Datos mínimos funcionales | Datos variados para testing |
| **Variedad** | Casos básicos | Múltiples escenarios |
| **Plataformas** | Principalmente web | Web, mobile, tablet |
| **Países** | Chile | 10+ países latinoamericanos |
| **Idiomas** | 3 (es, en, pt) | 5 (es, en, pt, fr, de) |
| **Dificultades** | easy, medium | easy, medium, hard |
| **Tipos evento** | 6 tipos | 10 tipos |

### Tests Agregados

```go
func testApplyMockData(ctx context.Context, db *mongo.Database) func(*testing.T) {
    // 1. Aplica mock data
    // 2. Verifica conteo: 10 + 3 + 8 + 3 + 5 + 6 = 35 documentos
    // 3. Test de idempotencia (ejecuta 2 veces)
    // 4. Verifica que NO se duplican documentos con _id explícito
}
```

### Beneficios

- ✅ **Type-safe**: Go verifica tipos en tiempo de compilación
- ✅ **Mayor cobertura**: 35 documentos vs 22 en seeds
- ✅ **Más variedad**: Diferentes plataformas, países, idiomas
- ✅ **Idempotente**: Se puede ejecutar múltiples veces
- ✅ **Testeable**: Tests de integración incluidos
- ✅ **Consistente**: Sigue mismo patrón que ApplySeeds()
- ✅ **Documentado**: GoDoc completo con comparación vs seeds

### Uso

```go
import "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"

// Ambiente de desarrollo con datos de prueba
migrations.ApplyAll(ctx, db)
migrations.ApplySeeds(ctx, db)      // Datos mínimos
migrations.ApplyMockData(ctx, db)   // Datos variados para testing
```

### Casos de Uso

**ApplySeeds()**: Ideal para inicialización mínima
- Datos esenciales del ecosistema
- Ambientes productivos
- CI/CD básico

**ApplyMockData()**: Ideal para desarrollo y demos
- Tests de integración complejos
- Demostración de features
- Desarrollo local
- QA/Staging con datos variados

### Impacto en Proyectos

- **edugo-worker**: Ahora puede usar `ApplyMockData()` para tests con más variedad
- **edugo-api-mobile**: Datos mock con eventos mobile/tablet para testing realista
- **Todos**: Consistencia con patrón PostgreSQL (que tiene `testing/*.sql`)

### Esfuerzo Real

- **Complejidad:** Media
- **Tiempo:** ~2.5 horas (creación de 35 documentos variados)
- **Líneas agregadas:** +1,089 líneas en mock_data.go
- **Líneas modificadas:** ~35 líneas en embed.go + test

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
| ~~TODO-001~~ | ~~ApplySeeds() vacía~~ | 🟡 Media | ✅ Resuelto | 2h |
| ~~TODO-002~~ | ~~ApplyMockData() vacía~~ | 🟡 Media | ✅ Resuelto | 2.5h |
| TODO-004 | Tests MongoDB | 🟠 Media-Alta | Parcial | 2-3h |
| TODO-005 | Validación schemas | 🟢 Baja | Pendiente | 1h |

### Total Completado: 4.5h
### Total Pendiente: 3-4h

---

**Última actualización:** Diciembre 2024
