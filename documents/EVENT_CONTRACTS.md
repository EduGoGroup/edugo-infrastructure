# 📬 Event Contracts - EduGo

Este documento describe los contratos de eventos para la comunicación entre microservicios vía RabbitMQ.

---

## 🏗️ Arquitectura de Mensajería

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              FLUJO DE EVENTOS                                            │
└─────────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────┐                ┌─────────────┐                ┌─────────────┐
│ api-mobile  │───publish────▶│  RabbitMQ   │───consume─────▶│   worker    │
│             │                │             │                │             │
│             │◀──consume─────│             │◀───publish────│             │
└─────────────┘                └─────────────┘                └─────────────┘
       │                              │                              │
       │                              │                              │
       ▼                              ▼                              ▼
┌─────────────┐                ┌─────────────┐                ┌─────────────┐
│api-admin    │───publish────▶│  Exchanges  │                │   OpenAI    │
│             │                │ • materials │                │   (AI)      │
│             │◀──consume─────│ • assessments│               │             │
└─────────────┘                │ • students  │                └─────────────┘
                               └─────────────┘
```

---

## 📋 Eventos Disponibles

| Evento | Versión | Publisher | Consumer | Descripción |
|--------|---------|-----------|----------|-------------|
| `material.uploaded` | v1.0 | api-mobile | worker | Material subido por docente |
| `material.deleted` | v1.0 | api-mobile | worker | Material eliminado |
| `assessment.generated` | v1.0 | worker | api-mobile | Assessment generado por IA |
| `student.enrolled` | v1.0 | api-admin | api-mobile | Estudiante matriculado |

---

## 📄 Estructura Base de Eventos

Todos los eventos siguen esta estructura base:

```json
{
  "event_id": "uuid-v7",           // ID único del evento
  "event_type": "material.uploaded", // Tipo de evento
  "event_version": "1.0",          // Versión del schema
  "timestamp": "2024-01-01T00:00:00Z", // ISO 8601
  "payload": { ... }               // Datos específicos del evento
}
```

---

## 📤 Event: `material.uploaded` (v1.0)

**Publicado por:** `api-mobile`  
**Consumido por:** `worker`  
**Trigger:** Docente sube un material educativo

### Schema JSON

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Material Uploaded Event",
  "description": "Evento publicado cuando un docente sube un material educativo",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "payload"],
  "properties": {
    "event_id": {
      "type": "string",
      "format": "uuid",
      "description": "ID único del evento (UUID v7)"
    },
    "event_type": {
      "type": "string",
      "const": "material.uploaded"
    },
    "event_version": {
      "type": "string",
      "const": "1.0"
    },
    "timestamp": {
      "type": "string",
      "format": "date-time"
    },
    "payload": {
      "type": "object",
      "required": ["material_id", "school_id", "teacher_id", "file_url", "file_size_bytes", "file_type"],
      "properties": {
        "material_id": {
          "type": "string",
          "format": "uuid",
          "description": "ID del material en PostgreSQL"
        },
        "school_id": {
          "type": "string",
          "format": "uuid",
          "description": "ID de la escuela"
        },
        "teacher_id": {
          "type": "string",
          "format": "uuid",
          "description": "ID del docente que subió el material"
        },
        "file_url": {
          "type": "string",
          "format": "uri",
          "description": "URL del archivo en S3"
        },
        "file_size_bytes": {
          "type": "integer",
          "minimum": 0,
          "description": "Tamaño del archivo en bytes"
        },
        "file_type": {
          "type": "string",
          "description": "MIME type (ej: application/pdf)"
        },
        "metadata": {
          "type": "object",
          "properties": {
            "title": { "type": "string" },
            "description": { "type": "string" },
            "subject": { "type": "string" },
            "grade": { "type": "string" }
          }
        }
      }
    }
  }
}
```

### Ejemplo

```json
{
  "event_id": "01916a3c-4d2e-7000-8000-000000000001",
  "event_type": "material.uploaded",
  "event_version": "1.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "material_id": "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    "school_id": "s1111111-1111-1111-1111-111111111111",
    "teacher_id": "t2222222-2222-2222-2222-222222222222",
    "file_url": "https://edugo-materials.s3.amazonaws.com/materials/abc123.pdf",
    "file_size_bytes": 1048576,
    "file_type": "application/pdf",
    "metadata": {
      "title": "Introducción a Java",
      "description": "Material sobre programación orientada a objetos",
      "subject": "Informática",
      "grade": "3° Medio"
    }
  }
}
```

### Flujo de Procesamiento

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ api-mobile  │────▶│  RabbitMQ   │────▶│   worker    │────▶│  OpenAI     │
│             │     │             │     │             │     │             │
│ 1. Sube PDF │     │ 2. Queue    │     │ 3. Consume  │     │ 4. Genera   │
│    a S3     │     │    evento   │     │    evento   │     │    quiz     │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │ 5. Guarda   │
                                        │    MongoDB  │
                                        │    + PG     │
                                        └─────────────┘
```

---

## 📤 Event: `assessment.generated` (v1.0)

**Publicado por:** `worker`  
**Consumido por:** `api-mobile`  
**Trigger:** Worker completa generación de assessment con IA

### Schema JSON

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Assessment Generated Event",
  "description": "Evento publicado cuando el worker genera un assessment con IA",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "payload"],
  "properties": {
    "event_id": {
      "type": "string",
      "format": "uuid"
    },
    "event_type": {
      "type": "string",
      "const": "assessment.generated"
    },
    "event_version": {
      "type": "string",
      "const": "1.0"
    },
    "timestamp": {
      "type": "string",
      "format": "date-time"
    },
    "payload": {
      "type": "object",
      "required": ["material_id", "mongo_document_id", "questions_count"],
      "properties": {
        "material_id": {
          "type": "string",
          "format": "uuid",
          "description": "ID del material para el cual se generó"
        },
        "mongo_document_id": {
          "type": "string",
          "pattern": "^[0-9a-fA-F]{24}$",
          "description": "ObjectId del documento en MongoDB"
        },
        "questions_count": {
          "type": "integer",
          "minimum": 1,
          "description": "Número de preguntas generadas"
        },
        "processing_time_ms": {
          "type": "integer",
          "minimum": 0,
          "description": "Tiempo de procesamiento en ms"
        }
      }
    }
  }
}
```

### Ejemplo

```json
{
  "event_id": "01916a3c-5e4f-7000-8000-000000000002",
  "event_type": "assessment.generated",
  "event_version": "1.0",
  "timestamp": "2024-01-15T10:35:00Z",
  "payload": {
    "material_id": "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    "mongo_document_id": "507f1f77bcf86cd799439011",
    "questions_count": 10,
    "processing_time_ms": 5200
  }
}
```

---

## 📤 Event: `material.deleted` (v1.0)

**Publicado por:** `api-mobile`  
**Consumido por:** `worker`  
**Trigger:** Docente elimina un material

### Schema JSON

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Material Deleted Event",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "payload"],
  "properties": {
    "event_id": { "type": "string", "format": "uuid" },
    "event_type": { "type": "string", "const": "material.deleted" },
    "event_version": { "type": "string", "const": "1.0" },
    "timestamp": { "type": "string", "format": "date-time" },
    "payload": {
      "type": "object",
      "required": ["material_id", "school_id", "deleted_by"],
      "properties": {
        "material_id": { "type": "string", "format": "uuid" },
        "school_id": { "type": "string", "format": "uuid" },
        "deleted_by": { "type": "string", "format": "uuid" },
        "delete_files": { "type": "boolean", "default": true }
      }
    }
  }
}
```

### Ejemplo

```json
{
  "event_id": "01916a3c-6f70-7000-8000-000000000003",
  "event_type": "material.deleted",
  "event_version": "1.0",
  "timestamp": "2024-01-15T11:00:00Z",
  "payload": {
    "material_id": "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    "school_id": "s1111111-1111-1111-1111-111111111111",
    "deleted_by": "t2222222-2222-2222-2222-222222222222",
    "delete_files": true
  }
}
```

---

## 📤 Event: `student.enrolled` (v1.0)

**Publicado por:** `api-administracion`  
**Consumido por:** `api-mobile`  
**Trigger:** Admin matricula estudiante en unidad académica

### Schema JSON

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Student Enrolled Event",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "payload"],
  "properties": {
    "event_id": { "type": "string", "format": "uuid" },
    "event_type": { "type": "string", "const": "student.enrolled" },
    "event_version": { "type": "string", "const": "1.0" },
    "timestamp": { "type": "string", "format": "date-time" },
    "payload": {
      "type": "object",
      "required": ["membership_id", "student_id", "school_id", "academic_unit_id"],
      "properties": {
        "membership_id": { "type": "string", "format": "uuid" },
        "student_id": { "type": "string", "format": "uuid" },
        "school_id": { "type": "string", "format": "uuid" },
        "academic_unit_id": { "type": "string", "format": "uuid" },
        "enrolled_at": { "type": "string", "format": "date-time" },
        "enrolled_by": { "type": "string", "format": "uuid" }
      }
    }
  }
}
```

### Ejemplo

```json
{
  "event_id": "01916a3c-7081-7000-8000-000000000004",
  "event_type": "student.enrolled",
  "event_version": "1.0",
  "timestamp": "2024-01-15T09:00:00Z",
  "payload": {
    "membership_id": "m3333333-3333-3333-3333-333333333333",
    "student_id": "u4444444-4444-4444-4444-444444444444",
    "school_id": "s1111111-1111-1111-1111-111111111111",
    "academic_unit_id": "au555555-5555-5555-5555-555555555555",
    "enrolled_at": "2024-01-15T09:00:00Z",
    "enrolled_by": "admin123-1111-1111-1111-111111111111"
  }
}
```

---

## 🔄 RabbitMQ Configuration

### Exchanges

| Exchange | Type | Durability |
|----------|------|------------|
| `edugo.materials` | topic | durable |
| `edugo.assessments` | topic | durable |
| `edugo.students` | topic | durable |

### Queues

| Queue | Binding Key | Consumer |
|-------|-------------|----------|
| `worker.materials.process` | `material.uploaded`, `material.deleted` | worker |
| `api-mobile.assessments.ready` | `assessment.generated` | api-mobile |
| `api-mobile.students.enrolled` | `student.enrolled` | api-mobile |

### Routing Keys

```
material.uploaded     → worker.materials.process
material.deleted      → worker.materials.process
assessment.generated  → api-mobile.assessments.ready
student.enrolled      → api-mobile.students.enrolled
```

---

## 🛡️ Validación de Eventos

### Uso del Validador

```go
import "github.com/EduGoGroup/edugo-infrastructure/messaging"

// Publisher
validator := messaging.NewEventValidator()
if err := validator.Validate(event); err != nil {
    return fmt.Errorf("invalid event: %w", err)
}
publisher.Publish(event)

// Consumer
if err := validator.ValidateJSON(msgBody, "material.uploaded", "1.0"); err != nil {
    logger.Error("invalid event received", err)
    return sendToDLQ(msg, err)
}
```

### Dead Letter Queue (DLQ)

Eventos que fallan validación o procesamiento van a DLQ:
- `worker.materials.dlq`
- `api-mobile.assessments.dlq`

---

## 📈 Versionamiento de Eventos

### Reglas

| Tipo Cambio | Versión | Ejemplo |
|-------------|---------|---------|
| Campo opcional nuevo | 1.0 → 1.1 | Agregar `metadata.tags` |
| Campo requerido nuevo | 1.0 → 2.0 | Agregar `payload.priority` |
| Cambio de tipo | 1.0 → 2.0 | `file_size` string → integer |
| Eliminar campo | 1.0 → 2.0 | Remover `payload.legacy_id` |

### Manejo Multi-versión

```go
switch event.EventVersion {
case "1.0", "1.1":
    return handleV1(event)
case "2.0":
    return handleV2(event)
default:
    return fmt.Errorf("unsupported version: %s", event.EventVersion)
}
```

---

## 🔍 Trazabilidad

Cada evento incluye `event_id` (UUID v7) que permite:
- Rastrear evento en logs
- Correlacionar con operaciones en BD
- Debugging end-to-end

```
[api-mobile] Published event_id=01916a3c-4d2e-7000-8000-000000000001
[worker] Consumed event_id=01916a3c-4d2e-7000-8000-000000000001
[worker] Processing completed event_id=01916a3c-4d2e-7000-8000-000000000001
```

---

**Última actualización:** Diciembre 2024
