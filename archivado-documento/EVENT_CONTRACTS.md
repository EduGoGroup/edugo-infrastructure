# 📋 Contratos de Eventos - RabbitMQ

**Fecha:** 15 de Noviembre, 2025  
**Validación:** JSON Schema automática

---

## 🎯 Propósito

Define el formato exacto de TODOS los eventos RabbitMQ en el ecosistema EduGo.

**Regla:** Todos los eventos DEBEN validarse contra su schema antes de publicar/consumir.

---

## 📊 Configuración de RabbitMQ

### Exchange Principal

```yaml
Exchange:
  Name: edugo.topic
  Type: topic
  Durable: true
  Auto-delete: false
```

### Queues y Bindings

| Queue | Routing Key | Consumidor | DLQ |
|-------|-------------|------------|-----|
| material.processing | material.uploaded | worker | dlq.material.processing |
| assessment.notifications | assessment.generated | api-mobile | dlq.assessment.notifications |
| material.cleanup | material.deleted | worker | dlq.material.cleanup |
| student.sync | student.enrolled | api-mobile | dlq.student.sync |

---

## 📋 Eventos Definidos

### 1. material.uploaded (v1.0)

**Publicado por:** api-mobile  
**Consumido por:** worker  
**Routing key:** `material.uploaded`  
**Schema:** `schemas/events/material-uploaded-v1.schema.json`

**Propósito:** Notificar que un docente subió un nuevo material educativo que necesita ser procesado (generar resumen + quiz).

**Ejemplo:**
```json
{
  "event_id": "01JA8XYZ-1234-5678-90AB-CDEF12345678",
  "event_type": "material.uploaded",
  "event_version": "1.0",
  "timestamp": "2025-11-15T10:30:00Z",
  "payload": {
    "material_id": "66666666-6666-6666-6666-666666666666",
    "school_id": "44444444-4444-4444-4444-444444444444",
    "teacher_id": "22222222-2222-2222-2222-222222222222",
    "file_url": "s3://edugo-materials-dev/fisica-cuantica.pdf",
    "file_size_bytes": 2048000,
    "file_type": "application/pdf",
    "metadata": {
      "title": "Introducción a la Física Cuántica",
      "description": "Material educativo sobre conceptos básicos",
      "subject": "Física",
      "grade": "10th"
    }
  }
}
```

**Campos requeridos:**
- `event_id`: UUID único del evento
- `event_type`: Siempre "material.uploaded"
- `event_version`: Siempre "1.0" (para esta versión del schema)
- `timestamp`: ISO 8601
- `payload.material_id`: UUID del material en PostgreSQL
- `payload.school_id`: UUID de la escuela
- `payload.teacher_id`: UUID del docente
- `payload.file_url`: URL completa del archivo
- `payload.file_size_bytes`: Tamaño en bytes
- `payload.file_type`: MIME type

**Campos opcionales:**
- `payload.metadata.*`: Todos opcionales

---

### 2. assessment.generated (v1.0)

**Publicado por:** worker  
**Consumido por:** api-mobile  
**Routing key:** `assessment.generated`  
**Schema:** `schemas/events/assessment-generated-v1.schema.json`

**Propósito:** Notificar que el worker terminó de generar un assessment/quiz con IA.

**Ejemplo:**
```json
{
  "event_id": "01JA8XYZ-ABCD-EFGH-IJKL-MNOPQRSTUVWX",
  "event_type": "assessment.generated",
  "event_version": "1.0",
  "timestamp": "2025-11-15T10:35:00Z",
  "payload": {
    "material_id": "66666666-6666-6666-6666-666666666666",
    "mongo_document_id": "507f1f77bcf86cd799439011",
    "questions_count": 8,
    "processing_time_ms": 45000
  }
}
```

**Campos requeridos:**
- `payload.material_id`: UUID del material original
- `payload.mongo_document_id`: ObjectId (24 chars hex) del documento en MongoDB
- `payload.questions_count`: Número de preguntas generadas

**Campos opcionales:**
- `payload.processing_time_ms`: Tiempo de procesamiento

---

### 3. material.deleted (v1.0)

**Publicado por:** api-mobile  
**Consumido por:** worker  
**Routing key:** `material.deleted`  
**Schema:** `schemas/events/material-deleted-v1.schema.json`

**Propósito:** Notificar que un material fue eliminado (soft delete) para cleanup de archivos y datos.

**Ejemplo:**
```json
{
  "event_id": "01JA8XYZ-1111-2222-3333-444444444444",
  "event_type": "material.deleted",
  "event_version": "1.0",
  "timestamp": "2025-11-15T11:00:00Z",
  "payload": {
    "material_id": "88888888-8888-8888-8888-888888888888",
    "school_id": "44444444-4444-4444-4444-444444444444",
    "deleted_by_user_id": "11111111-1111-1111-1111-111111111111",
    "reason": "Contenido desactualizado"
  }
}
```

**Worker debe:**
- Eliminar archivo de S3
- Marcar assessment en MongoDB como deleted
- Cleanup de datos temporales

---

### 4. student.enrolled (v1.0)

**Publicado por:** api-admin  
**Consumido por:** api-mobile  
**Routing key:** `student.enrolled`  
**Schema:** `schemas/events/student-enrolled-v1.schema.json`

**Propósito:** Notificar que un estudiante fue inscrito en una escuela/curso.

**Ejemplo:**
```json
{
  "event_id": "01JA8XYZ-AAAA-BBBB-CCCC-DDDDDDDDDDDD",
  "event_type": "student.enrolled",
  "event_version": "1.0",
  "timestamp": "2025-11-15T09:00:00Z",
  "payload": {
    "student_id": "33333333-3333-3333-3333-333333333333",
    "school_id": "44444444-4444-4444-4444-444444444444",
    "academic_unit_id": "99999999-9999-9999-9999-999999999999",
    "membership_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "enrolled_by_user_id": "11111111-1111-1111-1111-111111111111"
  }
}
```

**api-mobile debe:**
- Sincronizar datos del estudiante si es necesario
- Actualizar caché de memberships
- Notificar al estudiante

---

## 🔄 Estrategia de Versionamiento

### Versión en el Evento

Todos los eventos incluyen `event_version` en formato "MAJOR.MINOR":

```json
{
  "event_version": "1.0"  // Major: 1, Minor: 0
}
```

### Reglas de Versionamiento

| Cambio | Version Bump | Backward Compatible |
|--------|--------------|---------------------|
| Agregar campo opcional | 1.0 → 1.1 | ✅ Sí |
| Modificar campo existente | 1.0 → 2.0 | ❌ No |
| Eliminar campo | 1.0 → 2.0 | ❌ No |
| Cambiar tipo de campo | 1.0 → 2.0 | ❌ No |
| Renombrar campo | 1.0 → 2.0 | ❌ No |

### Manejo en Consumer

```go
func (c *MaterialConsumer) Handle(msg Message) error {
    switch msg.EventVersion {
    case "1.0", "1.1":
        return c.handleV1(msg)  // v1.1 backward compatible con v1.0
    case "2.0":
        return c.handleV2(msg)
    default:
        return fmt.Errorf("unsupported event version: %s", msg.EventVersion)
    }
}
```

### Deprecación de Versiones

```
1. Publicar nueva versión (ej: v2.0)
2. Publisher publica AMBAS versiones por 2 sprints
3. Consumers actualizan para soportar v2.0
4. Después de 2 sprints, publisher solo envía v2.0
5. Eliminar soporte de v1.0 en consumers
```

---

## ✅ Validación Automática

### En Publisher (api-mobile, worker)

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

validator := schemas.NewEventValidator()

event := MaterialUploadedEvent{...}

// Validar antes de publicar
if err := validator.Validate(event); err != nil {
    logger.Error("event validation failed", err)
    return fmt.Errorf("invalid event: %w", err)
}

// Publicar solo si es válido
publisher.Publish(exchange, routingKey, event)
```

### En Consumer (worker, api-mobile)

```go
validator := schemas.NewEventValidator()

func (c *Consumer) HandleMessage(msg []byte) error {
    // Validar JSON contra schema
    if err := validator.ValidateJSON(msg, "material.uploaded", "1.0"); err != nil {
        logger.Error("invalid event received", 
            "error", err,
            "message", string(msg))
        
        // Enviar a DLQ (evento mal formado)
        return c.sendToDLQ(msg, err)
    }
    
    // Deserializar y procesar
    var event MaterialUploadedEvent
    json.Unmarshal(msg, &event)
    
    return c.processEvent(event)
}
```

---

## 🐛 Debugging de Eventos

### RabbitMQ Management UI

http://localhost:15672

- Ver mensajes en queues
- Inspeccionar payloads
- Verificar bindings

### Logs Estructurados

```go
logger.Info("publishing event",
    "event_id", event.EventID,
    "event_type", event.EventType,
    "event_version", event.EventVersion,
    "payload", event.Payload)
```

### DLQ (Dead Letter Queue)

Eventos inválidos van a DLQ:
- `dlq.material.processing`
- `dlq.assessment.notifications`
- `dlq.material.cleanup`
- `dlq.student.sync`

Revisar DLQ para detectar problemas:
```bash
# En RabbitMQ UI
Queues → dlq.* → Get Messages
```

---

## 📚 Referencias

- **JSON Schemas:** `schemas/events/*.schema.json`
- **Validador Go:** `schemas/validator.go` (pendiente crear)
- **Ejemplos de uso:** `schemas/README.md`

---

**Última actualización:** 15 de Noviembre, 2025  
**Mantenedor:** Equipo EduGo
