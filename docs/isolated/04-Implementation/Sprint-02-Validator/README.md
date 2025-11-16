# Sprint-02: Validator

## 🎯 Objetivo

Crear validador de eventos con JSON Schemas para garantizar contratos consistentes entre microservicios.

---

## ✅ Estado: FASE 1 COMPLETADA

**Archivo principal:** `schemas/validator.go` (130 líneas)
**Tests:** `schemas/example_test.go` (78 líneas, 2 tests)
**Fecha de completitud:** 2025-11-16

---

## 📦 Implementación

### API del Validador

```go
// Crear validador (carga todos los schemas)
validator, err := schemas.NewEventValidator()
if err != nil {
    log.Fatal(err)
}

// Opción 1: Auto-detect event_type y event_version
err = validator.Validate(event)

// Opción 2: Especificar explícitamente
err = validator.ValidateWithType(event, "material.uploaded", "1.0")

// Opción 3: Validar JSON bytes directamente
err = validator.ValidateJSON(jsonBytes, "material.uploaded", "1.0")
```

### Eventos Soportados

- ✅ `material.uploaded:1.0` - Material subido por profesor
- ✅ `assessment.generated:1.0` - Quiz generado por IA
- ✅ `material.deleted:1.0` - Material eliminado
- ✅ `student.enrolled:1.0` - Estudiante inscrito en curso

### Características implementadas

- ✅ EventValidator con cache de schemas en memoria
- ✅ Constructor carga 4 schemas automáticamente
- ✅ Schemas embebidos en binario con `//go:embed`
- ✅ 3 métodos de validación (Validate, ValidateWithType, ValidateJSON)
- ✅ Mensajes de error detallados con lista de violaciones
- ✅ No requiere servicios externos (schemas embebidos)

---

## 🧪 Tests

### Tests de Validación (Fase 1)

```bash
cd schemas
go test -v
```

**Tests implementados:**
- `TestMaterialUploadedValidation` - Evento válido vs inválido
- `TestValidateJSON` - Validación desde JSON bytes

**Resultado:** 2/2 tests passing

### Tests Exhaustivos (Fase 2)

Ver: `PHASE2_BRIDGE.md` para detalles completos

Pendiente:
- Tests para material.deleted v1.0
- Tests para student.enrolled v1.0
- Edge cases (event_type faltante, UUIDs inválidos)
- Benchmarks de performance

---

## 📁 Estructura de Archivos

```
schemas/
├── validator.go           # Validador (130 líneas)
├── example_test.go        # Tests (78 líneas)
├── go.mod                 # Dependencias
├── go.sum
├── README.md
└── events/
    ├── material-uploaded-v1.schema.json
    ├── assessment-generated-v1.schema.json
    ├── material-deleted-v1.schema.json
    └── student-enrolled-v1.schema.json
```

---

## 🚀 Uso

### Ejemplo 1: Validar evento de material subido

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

validator, _ := schemas.NewEventValidator()

event := map[string]interface{}{
    "event_id":      "550e8400-e29b-41d4-a716-446655440000",
    "event_type":    "material.uploaded",
    "event_version": "1.0",
    "timestamp":     "2025-11-16T10:30:00Z",
    "payload": map[string]interface{}{
        "material_id":     "123e4567-e89b-12d3-a456-426614174000",
        "school_id":       "234e5678-e89b-12d3-a456-426614174000",
        "teacher_id":      "345e6789-e89b-12d3-a456-426614174000",
        "file_url":        "s3://edugo-materials/test.pdf",
        "file_size_bytes": 2048000,
        "file_type":       "application/pdf",
    },
}

if err := validator.Validate(event); err != nil {
    log.Printf("Evento inválido: %v", err)
    return
}

// ✅ Evento válido, safe to publish
publisher.Publish(event)
```

### Ejemplo 2: Validar JSON raw (ej: desde RabbitMQ)

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

validator, _ := schemas.NewEventValidator()

// JSON recibido de RabbitMQ
jsonBytes := []byte(`{
    "event_id": "550e8400-e29b-41d4-a716-446655440000",
    "event_type": "assessment.generated",
    "event_version": "1.0",
    "timestamp": "2025-11-16T10:35:00Z",
    "payload": {
        "material_id": "123e4567-e89b-12d3-a456-426614174000",
        "mongo_document_id": "507f1f77bcf86cd799439011",
        "questions_count": 8
    }
}`)

err := validator.ValidateJSON(jsonBytes, "assessment.generated", "1.0")
if err != nil {
    log.Printf("Evento inválido: %v", err)
    // Rechazar mensaje o enviarlo a DLQ
    return
}

// ✅ Evento válido, procesar
processAssessmentGenerated(jsonBytes)
```

### Ejemplo 3: Manejo de errores

```go
event := map[string]interface{}{
    "event_id":   "invalid-uuid",  // UUID inválido
    "event_type": "material.uploaded",
    // event_version faltante
}

err := validator.Validate(event)
if err != nil {
    // Error: "event_version faltante o inválido"
    log.Printf("Error: %v", err)
}
```

---

## 🔍 Detalles de Implementación

### Constructor: NewEventValidator()

```go
func NewEventValidator() (*EventValidator, error) {
    v := &EventValidator{
        schemas: make(map[string]*gojsonschema.Schema),
    }

    // Cargar schemas desde filesystem embebido
    schemaFiles := map[string]string{
        "material.uploaded:1.0":    "events/material-uploaded-v1.schema.json",
        "assessment.generated:1.0": "events/assessment-generated-v1.schema.json",
        "material.deleted:1.0":     "events/material-deleted-v1.schema.json",
        "student.enrolled:1.0":     "events/student-enrolled-v1.schema.json",
    }

    for key, filename := range schemaFiles {
        if err := v.loadSchema(key, filename); err != nil {
            return nil, fmt.Errorf("error cargando schema %s: %w", key, err)
        }
    }

    return v, nil
}
```

### Método principal: Validate()

```go
func (v *EventValidator) Validate(event interface{}) error {
    // Extraer event_type y event_version del evento
    eventMap, ok := event.(map[string]interface{})
    if !ok {
        return errors.New("evento debe ser un objeto JSON")
    }

    eventType, ok := eventMap["event_type"].(string)
    if !ok {
        return errors.New("event_type faltante o inválido")
    }

    eventVersion, ok := eventMap["event_version"].(string)
    if !ok {
        return errors.New("event_version faltante o inválido")
    }

    return v.ValidateWithType(event, eventType, eventVersion)
}
```

---

## 📊 Schemas Embebidos

### material.uploaded v1.0

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "payload"],
  "properties": {
    "event_id": { "type": "string", "format": "uuid" },
    "event_type": { "const": "material.uploaded" },
    "event_version": { "const": "1.0" },
    "timestamp": { "type": "string", "format": "date-time" },
    "payload": {
      "type": "object",
      "required": ["material_id", "school_id", "teacher_id", "file_url", "file_size_bytes", "file_type"],
      "properties": {
        "material_id": { "type": "string", "format": "uuid" },
        "school_id": { "type": "string", "format": "uuid" },
        "teacher_id": { "type": "string", "format": "uuid" },
        "file_url": { "type": "string", "pattern": "^s3://" },
        "file_size_bytes": { "type": "integer", "minimum": 1 },
        "file_type": { "type": "string" }
      }
    }
  }
}
```

Ver carpeta `events/` para los demás schemas.

---

## 📝 Próximos Pasos (Fase 2)

1. Tests para material.deleted v1.0 y student.enrolled v1.0
2. Tests de edge cases (event_type faltante, UUIDs inválidos, etc.)
3. Benchmarks de performance (objetivo: <1s para 10,000 eventos)
4. Validación exhaustiva de todos los campos
5. Documentar comportamiento con campos extra en payload

Ver: `PHASE2_BRIDGE.md` para instrucciones detalladas

---

## 📚 Referencias

- Documentación principal: `README.md` (raíz del proyecto)
- Contratos de eventos: `EVENT_CONTRACTS.md`
- JSON Schemas: `schemas/events/`
- Phase 2 Bridge: `PHASE2_BRIDGE.md`
- gojsonschema docs: https://github.com/xeipuuv/gojsonschema

---

**Versión:** 0.1.1
**Estado:** Fase 1 COMPLETADA
**Próximo paso:** Fase 2 - Tests exhaustivos
