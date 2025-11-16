# PHASE2 BRIDGE - Sprint-02-Validator

## 📋 Resumen

**Sprint:** Sprint-02-Validator
**Archivo principal:** `schemas/validator.go`
**Tests:** `schemas/example_test.go`
**Estado Fase 1:** ✅ COMPLETADO

---

## ✅ Completado en Fase 1

### Implementación

- [x] Struct `EventValidator` con cache de schemas
- [x] Constructor `NewEventValidator()` carga 4 schemas embebidos
- [x] Método `Validate(event)` extrae event_type y event_version automáticamente
- [x] Método `ValidateWithType(event, type, version)` validación explícita
- [x] Método `ValidateJSON(jsonBytes, type, version)` para JSON raw
- [x] Schemas embebidos con `//go:embed events/*.json`
- [x] Mensajes de error detallados con lista de violaciones
- [x] Soporte para 4 eventos: material.uploaded, assessment.generated, material.deleted, student.enrolled

### Tests Unitarios

- [x] `TestMaterialUploadedValidation` - Evento válido vs inválido
- [x] `TestValidateJSON` - Validación desde JSON bytes
- [x] Tests con UUIDs reales generados
- [x] Validación de campos requeridos
- [x] Validación de formatos (UUID, timestamp, file_url)

### Código

```go
// validator.go - 130 líneas
// Funciones principales:
- NewEventValidator()                         // Constructor, carga schemas
- loadSchema(key, filename)                   // Carga schema individual
- Validate(event)                             // Auto-detect event_type
- ValidateWithType(event, type, version)      // Validación explícita
- ValidateJSON(jsonBytes, type, version)      // Desde JSON raw

// Schemas soportados:
- material.uploaded:1.0
- assessment.generated:1.0
- material.deleted:1.0
- student.enrolled:1.0
```

**Total de líneas:** 130 (validator.go) + 78 (example_test.go)
**Tests unitarios:** 2 tests, todos passing
**Schemas embebidos:** 4 archivos JSON

---

## ⏳ Pendiente para Fase 2

### Tests Adicionales

1. **Test: Validar TODOS los schemas (4 eventos)**
   - Descripción: Cobertura completa de los 4 eventos
   - Requiere: Solo Go (no servicios externos)
   - Validar:
     - material.uploaded v1.0 (ya testeado ✅)
     - assessment.generated v1.0 (ya testeado ✅)
     - material.deleted v1.0 (pendiente)
     - student.enrolled v1.0 (pendiente)

2. **Test: Edge cases de validación**
   - Descripción: Casos límite y errores
   - Validar:
     - event_type faltante → error claro
     - event_version inválida → error claro
     - Schema no encontrado (ej: "unknown.event:1.0") → error claro
     - Payload con tipos incorrectos (string en vez de int)
     - UUIDs inválidos
     - Timestamps en formato incorrecto
     - URLs malformadas

3. **Test: Performance con grandes volúmenes**
   - Descripción: Validar 10,000 eventos en <1 segundo
   - Requiere: Benchmark tests
   - Validar:
     - Schema cache funciona (no recarga en cada validación)
     - Memoria no crece indefinidamente
     - Throughput aceptable para producción

4. **Test: Validación de todos los campos de cada schema**
   - Descripción: Exhaustivo, campo por campo
   - Schemas a validar:
     - material.uploaded: material_id, school_id, teacher_id, file_url, file_size_bytes, file_type
     - assessment.generated: material_id, mongo_document_id, questions_count
     - material.deleted: material_id, school_id, teacher_id, deleted_by_id
     - student.enrolled: student_id, school_id, academic_unit_id, enrollment_date

### Edge Cases

1. **Evento sin event_type**
   - Escenario: JSON sin campo event_type
   - Validación: Error "event_type faltante o inválido"

2. **Schema versión no soportada**
   - Escenario: material.uploaded:2.0 (no existe)
   - Validación: Error "schema no encontrado para material.uploaded:2.0"

3. **Payload con campos extra**
   - Escenario: JSON con campos no definidos en schema
   - Validación: ¿Rechazar o permitir? (verificar behavior de gojsonschema)

4. **UUIDs en formato incorrecto**
   - Escenario: "123" en vez de UUID válido
   - Validación: Error de formato

5. **file_size_bytes negativo**
   - Escenario: -100 en material.uploaded
   - Validación: Error (verificar si schema lo cubre)

---

## 🔧 Prerequisitos para Fase 2

### Servicios Requeridos

```bash
# NO requiere servicios externos
# Solo Go 1.24+

# Opcional: RabbitMQ para tests de integración end-to-end
make dev-up-messaging
```

### Variables de Entorno

```bash
# NO requiere variables de entorno
# Schemas están embebidos en binario
```

### Datos de Prueba

```bash
# Ver schemas/events/ para ejemplos de estructura
ls -la schemas/events/

# Archivos disponibles:
# - material-uploaded-v1.schema.json
# - assessment-generated-v1.schema.json
# - material-deleted-v1.schema.json
# - student-enrolled-v1.schema.json
```

---

## 🧪 Tests Adicionales a Implementar

### Archivo: `schemas/validator_comprehensive_test.go`

```go
package schemas_test

import (
    "testing"
    "github.com/google/uuid"
    "github.com/EduGoGroup/edugo-infrastructure/schemas"
)

func TestMaterialDeletedValidation(t *testing.T) {
    validator, _ := schemas.NewEventValidator()

    validEvent := map[string]interface{}{
        "event_id":      uuid.New().String(),
        "event_type":    "material.deleted",
        "event_version": "1.0",
        "timestamp":     "2025-11-16T10:00:00Z",
        "payload": map[string]interface{}{
            "material_id":   uuid.New().String(),
            "school_id":     uuid.New().String(),
            "teacher_id":    uuid.New().String(),
            "deleted_by_id": uuid.New().String(),
        },
    }

    if err := validator.Validate(validEvent); err != nil {
        t.Errorf("Evento válido rechazado: %v", err)
    }
}

func TestStudentEnrolledValidation(t *testing.T) {
    validator, _ := schemas.NewEventValidator()

    validEvent := map[string]interface{}{
        "event_id":      uuid.New().String(),
        "event_type":    "student.enrolled",
        "event_version": "1.0",
        "timestamp":     "2025-11-16T10:00:00Z",
        "payload": map[string]interface{}{
            "student_id":       uuid.New().String(),
            "school_id":        uuid.New().String(),
            "academic_unit_id": uuid.New().String(),
            "enrollment_date":  "2025-11-16",
        },
    }

    if err := validator.Validate(validEvent); err != nil {
        t.Errorf("Evento válido rechazado: %v", err)
    }
}

func TestEventTypeValidation(t *testing.T) {
    validator, _ := schemas.NewEventValidator()

    tests := []struct {
        name        string
        event       map[string]interface{}
        shouldError bool
        errorMsg    string
    }{
        {
            name: "missing event_type",
            event: map[string]interface{}{
                "event_id":      uuid.New().String(),
                "event_version": "1.0",
                "timestamp":     "2025-11-16T10:00:00Z",
                "payload":       map[string]interface{}{},
            },
            shouldError: true,
            errorMsg:    "event_type faltante o inválido",
        },
        {
            name: "missing event_version",
            event: map[string]interface{}{
                "event_id":   uuid.New().String(),
                "event_type": "material.uploaded",
                "timestamp":  "2025-11-16T10:00:00Z",
                "payload":    map[string]interface{}{},
            },
            shouldError: true,
            errorMsg:    "event_version faltante o inválido",
        },
        {
            name: "unknown event type",
            event: map[string]interface{}{
                "event_id":      uuid.New().String(),
                "event_type":    "unknown.event",
                "event_version": "1.0",
                "timestamp":     "2025-11-16T10:00:00Z",
                "payload":       map[string]interface{}{},
            },
            shouldError: true,
            errorMsg:    "schema no encontrado",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.Validate(tt.event)

            if tt.shouldError {
                if err == nil {
                    t.Errorf("Expected error but got nil")
                }
                // Opcionalmente verificar mensaje de error
            } else {
                if err != nil {
                    t.Errorf("Expected no error but got: %v", err)
                }
            }
        })
    }
}

func BenchmarkValidation(b *testing.B) {
    validator, _ := schemas.NewEventValidator()

    event := map[string]interface{}{
        "event_id":      uuid.New().String(),
        "event_type":    "material.uploaded",
        "event_version": "1.0",
        "timestamp":     "2025-11-16T10:00:00Z",
        "payload": map[string]interface{}{
            "material_id":     uuid.New().String(),
            "school_id":       uuid.New().String(),
            "teacher_id":      uuid.New().String(),
            "file_url":        "s3://edugo-materials/test.pdf",
            "file_size_bytes": float64(2048000),
            "file_type":       "application/pdf",
        },
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        validator.Validate(event)
    }
}
```

### Casos de Prueba

1. **Happy path:** Validar los 4 eventos con datos correctos
2. **Missing fields:** Eventos con campos requeridos faltantes
3. **Invalid types:** Strings donde deberían ser integers, etc.
4. **Invalid formats:** UUIDs malformados, timestamps incorrectos
5. **Unknown schemas:** event_type que no existe
6. **Performance:** 10,000 validaciones en <1 segundo

---

## 📝 Notas para Fase 2

- Los schemas están embebidos con `//go:embed` - binario autónomo ✅
- Actualmente solo hay 4 eventos - fácil extender agregando más schemas
- Considerar agregar versionado de schemas (v1.0, v2.0) en el futuro
- gojsonschema es robusto pero considerar alternativas si performance es crítica
- Cache de schemas en memoria - OK para <100 schemas
- ValidateJSON es útil para RabbitMQ consumers que reciben JSON raw

---

## ✅ Checklist Fase 2

- [ ] Test: material.deleted v1.0 (happy path + invalid)
- [ ] Test: student.enrolled v1.0 (happy path + invalid)
- [ ] Test: edge cases (event_type faltante, schema no encontrado, etc.)
- [ ] Test: validación exhaustiva de todos los campos
- [ ] Test: UUIDs inválidos
- [ ] Test: timestamps en formato incorrecto
- [ ] Test: file_size_bytes negativo o cero
- [ ] Benchmark: validar 10,000 eventos
- [ ] Medir cobertura de tests (objetivo: >90%)
- [ ] Documentar comportamiento con campos extra en payload
- [ ] Agregar ejemplos de integración con RabbitMQ (opcional)
- [ ] Commit y push

---

**Fase 1 completada:** 2025-11-16
**Próximo paso:** Implementar tests exhaustivos y benchmarks
**Estimado Fase 2:** 1-2 horas
**Ventaja:** NO requiere PostgreSQL ni servicios externos ✅
