# Tareas: Sprint-02-Validator

**Sprint:** Sprint-02-Validator  
**Duración:** 2-3 horas

---

## TASK-001: Crear schemas/validator.go

**Descripción:** Implementar validador de eventos con JSON Schema

**Pasos:**

1. Crear archivo `schemas/validator.go`

```go
package schemas

import (
    "embed"
    "fmt"
    "sync"

    "github.com/xeipuuv/gojsonschema"
)

//go:embed events/*.schema.json
var schemaFiles embed.FS

// Validator valida eventos contra JSON Schemas
type Validator struct {
    schemas map[string]*gojsonschema.Schema
    mu      sync.RWMutex
}

// NewValidator crea un nuevo validador cargando todos los schemas
func NewValidator() (*Validator, error) {
    v := &Validator{
        schemas: make(map[string]*gojsonschema.Schema),
    }

    // Cargar todos los schemas
    entries, err := schemaFiles.ReadDir("events")
    if err != nil {
        return nil, fmt.Errorf("error reading schemas directory: %w", err)
    }

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        // Leer schema file
        schemaData, err := schemaFiles.ReadFile("events/" + entry.Name())
        if err != nil {
            return nil, fmt.Errorf("error reading schema %s: %w", entry.Name(), err)
        }

        // Parsear schema
        schemaLoader := gojsonschema.NewBytesLoader(schemaData)
        schema, err := gojsonschema.NewSchema(schemaLoader)
        if err != nil {
            return nil, fmt.Errorf("error parsing schema %s: %w", entry.Name(), err)
        }

        // Extraer nombre del schema (sin extensión)
        // material-uploaded-v1.schema.json → material-uploaded-v1
        schemaName := entry.Name()
        schemaName = schemaName[:len(schemaName)-len(".schema.json")]

        v.schemas[schemaName] = schema
    }

    return v, nil
}

// Validate valida un evento contra un schema específico
func (v *Validator) Validate(event interface{}, schemaName string) error {
    v.mu.RLock()
    schema, exists := v.schemas[schemaName]
    v.mu.RUnlock()

    if !exists {
        return fmt.Errorf("schema not found: %s", schemaName)
    }

    // Validar
    documentLoader := gojsonschema.NewGoLoader(event)
    result, err := schema.Validate(documentLoader)
    if err != nil {
        return fmt.Errorf("validation error: %w", err)
    }

    if !result.Valid() {
        // Construir error con todos los problemas
        errMsg := "validation failed:"
        for _, err := range result.Errors() {
            errMsg += fmt.Sprintf("\n  - %s", err.String())
        }
        return fmt.Errorf("%s", errMsg)
    }

    return nil
}

// ListSchemas retorna lista de schemas disponibles
func (v *Validator) ListSchemas() []string {
    v.mu.RLock()
    defer v.mu.RUnlock()

    schemas := make([]string, 0, len(v.schemas))
    for name := range v.schemas {
        schemas = append(schemas, name)
    }
    return schemas
}
```

2. Actualizar `schemas/go.mod`

```go
module github.com/EduGoGroup/edugo-infrastructure/schemas

go 1.24

require (
    github.com/xeipuuv/gojsonschema v1.2.0
)
```

3. Ejecutar `go mod tidy`

**Estimación:** 60 minutos

---

## TASK-002: Crear tests del validador

**Descripción:** Tests con eventos válidos e inválidos

**Pasos:**

1. Crear `schemas/validator_test.go`

```go
package schemas

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
    validator, err := NewValidator()
    require.NoError(t, err)
    require.NotNil(t, validator)

    // Verificar que cargó los 4 schemas
    schemas := validator.ListSchemas()
    assert.Len(t, schemas, 4)
    assert.Contains(t, schemas, "material-uploaded-v1")
    assert.Contains(t, schemas, "assessment-generated-v1")
    assert.Contains(t, schemas, "material-deleted-v1")
    assert.Contains(t, schemas, "student-enrolled-v1")
}

func TestValidator_Validate_ValidEvent(t *testing.T) {
    validator, err := NewValidator()
    require.NoError(t, err)

    // Evento válido
    event := map[string]interface{}{
        "event_id":      "01JCXYZ123ABC456DEF789GHI",
        "event_type":    "material.uploaded",
        "event_version": "1.0",
        "timestamp":     "2025-11-16T10:00:00Z",
        "payload": map[string]interface{}{
            "material_id": "550e8400-e29b-41d4-a716-446655440000",
            "school_id":   "660e8400-e29b-41d4-a716-446655440001",
            "teacher_id":  "770e8400-e29b-41d4-a716-446655440002",
            "unit_id":     "880e8400-e29b-41d4-a716-446655440003",
            "file_url":    "s3://edugo-materials/test.pdf",
            "file_size_bytes": 2048000,
            "file_type":   "application/pdf",
            "title":       "Test Material",
        },
    }

    err = validator.Validate(event, "material-uploaded-v1")
    assert.NoError(t, err)
}

func TestValidator_Validate_InvalidEvent(t *testing.T) {
    validator, err := NewValidator()
    require.NoError(t, err)

    tests := []struct {
        name          string
        event         map[string]interface{}
        schemaName    string
        expectedError string
    }{
        {
            name: "missing event_id",
            event: map[string]interface{}{
                "event_type": "material.uploaded",
                "timestamp":  "2025-11-16T10:00:00Z",
                "payload":    map[string]interface{}{},
            },
            schemaName:    "material-uploaded-v1",
            expectedError: "event_id",
        },
        {
            name: "invalid event_type",
            event: map[string]interface{}{
                "event_id":   "01JCXYZ123",
                "event_type": "wrong.type",
                "timestamp":  "2025-11-16T10:00:00Z",
                "payload":    map[string]interface{}{},
            },
            schemaName:    "material-uploaded-v1",
            expectedError: "event_type",
        },
        {
            name: "missing payload fields",
            event: map[string]interface{}{
                "event_id":   "01JCXYZ123",
                "event_type": "material.uploaded",
                "event_version": "1.0",
                "timestamp":  "2025-11-16T10:00:00Z",
                "payload":    map[string]interface{}{
                    // Falta material_id, file_url, etc.
                },
            },
            schemaName:    "material-uploaded-v1",
            expectedError: "material_id",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.Validate(tt.event, tt.schemaName)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedError)
        })
    }
}

func TestValidator_Validate_SchemaNotFound(t *testing.T) {
    validator, err := NewValidator()
    require.NoError(t, err)

    event := map[string]interface{}{}
    err = validator.Validate(event, "non-existent-schema")
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "schema not found")
}
```

2. Ejecutar tests

```bash
cd schemas
go test -v
```

**Estimación:** 45 minutos

---

## TASK-003: Crear ejemplo de uso

**Descripción:** Documentar cómo integrar en api-mobile y worker

**Pasos:**

1. Crear `schemas/examples/validate_event.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/EduGoGroup/edugo-infrastructure/schemas"
)

func main() {
    // Inicializar validador
    validator, err := schemas.NewValidator()
    if err != nil {
        log.Fatalf("Error creating validator: %v", err)
    }

    // Ejemplo 1: Validar evento válido
    validEvent := map[string]interface{}{
        "event_id":      "01JCXYZ123ABC456DEF789GHI",
        "event_type":    "material.uploaded",
        "event_version": "1.0",
        "timestamp":     "2025-11-16T10:00:00Z",
        "payload": map[string]interface{}{
            "material_id":     "550e8400-e29b-41d4-a716-446655440000",
            "file_url":        "s3://edugo-materials/test.pdf",
            "file_size_bytes": 2048000,
            "file_type":       "application/pdf",
        },
    }

    if err := validator.Validate(validEvent, "material-uploaded-v1"); err != nil {
        fmt.Printf("❌ Validation failed: %v\n", err)
    } else {
        fmt.Println("✅ Event is valid")
    }

    // Ejemplo 2: Validar evento inválido
    invalidEvent := map[string]interface{}{
        "event_type": "material.uploaded",
        // Falta event_id, timestamp, payload
    }

    if err := validator.Validate(invalidEvent, "material-uploaded-v1"); err != nil {
        fmt.Printf("✅ Correctly detected invalid event: %v\n", err)
    }

    // Ejemplo 3: Uso en api-mobile (publisher)
    fmt.Println("\nEjemplo de uso en api-mobile:")
    fmt.Println("------------------------------")
    examplePublisher()

    // Ejemplo 4: Uso en worker (consumer)
    fmt.Println("\nEjemplo de uso en worker:")
    fmt.Println("-------------------------")
    exampleConsumer()
}

func examplePublisher() {
    code := `
// En api-mobile al publicar evento
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

type MaterialPublisher struct {
    validator *schemas.Validator
    rabbit    *rabbit.Publisher
}

func (p *MaterialPublisher) PublishMaterialUploaded(material *Material) error {
    event := Event{
        EventID:      generateUUIDv7(),
        EventType:    "material.uploaded",
        EventVersion: "1.0",
        Timestamp:    time.Now(),
        Payload: map[string]interface{}{
            "material_id": material.ID,
            "file_url":    material.FileURL,
            // ... otros campos
        },
    }

    // Validar ANTES de publicar
    if err := p.validator.Validate(event, "material-uploaded-v1"); err != nil {
        return fmt.Errorf("invalid event: %w", err)
    }

    // Publicar evento válido
    return p.rabbit.Publish("material.uploaded", event)
}
`
    fmt.Println(code)
}

func exampleConsumer() {
    code := `
// En worker al consumir evento
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

type MaterialConsumer struct {
    validator *schemas.Validator
}

func (c *MaterialConsumer) HandleMessage(msg []byte) error {
    var event Event
    if err := json.Unmarshal(msg, &event); err != nil {
        return err
    }

    // Validar evento recibido
    if err := c.validator.Validate(event, "material-uploaded-v1"); err != nil {
        logger.Error("invalid event received", err)
        // NO reintentar (enviar a DLQ directamente)
        return nil
    }

    // Procesar evento válido
    return c.processMaterial(event.Payload)
}
`
    fmt.Println(code)
}
```

2. Documentar en README principal

**Estimación:** 30 minutos

---

## ✅ Checklist de Completitud

- [ ] validator.go creado
- [ ] go.mod actualizado con gojsonschema
- [ ] Tests con eventos válidos pasan
- [ ] Tests con eventos inválidos detectan errores
- [ ] Ejemplo de uso creado
- [ ] README actualizado
