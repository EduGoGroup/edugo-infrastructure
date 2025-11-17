# 📦 MongoDB Schema - EduGo

**Owner:** edugo-infrastructure
**Database:** edugo
**Fecha:** 16 de Noviembre, 2025

---

## 🎯 Propósito

Este documento define el esquema de MongoDB para EduGo, incluyendo todas las colecciones, sus estructuras, índices y validaciones.

**Filosofía:** MongoDB almacena datos que requieren flexibilidad de esquema, contenido no estructurado, o gran volumen de eventos/logs. PostgreSQL almacena datos relacionales y transaccionales.

---

## 📊 Colecciones

| Colección | Propósito | Relacionada con PostgreSQL | TTL |
|-----------|-----------|----------------------------|-----|
| **material_assessment** | Contenido de assessments generados por IA | assessment (metadata) | No |
| **material_content** | Contenido extraído de materiales educativos | materials (metadata) | No |
| **assessment_attempt_result** | Resultados detallados de intentos | assessment_attempt (metadata) | No |
| **audit_logs** | Logs de auditoría del sistema | - | 90 días |
| **notifications** | Notificaciones para usuarios | - | 30 días (archivadas) |
| **analytics_events** | Eventos de analítica | - | 365 días |

---

## 1️⃣ material_assessment

**Propósito:** Almacena el contenido completo de assessments/quizzes generados por IA.

**Relacionada con PostgreSQL:** La tabla `assessment` almacena metadata y referencia estos documentos via `mongo_document_id`.

### Estructura

```javascript
{
  _id: ObjectId,
  material_id: String,              // UUID del material en PostgreSQL
  questions: [
    {
      question_index: Int,           // Índice 0-based
      question_text: String,
      question_type: String,         // "multiple_choice" | "true_false" | "short_answer"
      options: [
        {
          option_index: Int,
          text: String,
          is_correct: Boolean
        }
      ],
      explanation: String            // Opcional
    }
  ],
  metadata: {
    subject: String,
    grade: String,
    difficulty: String,              // "easy" | "medium" | "hard"
    estimated_time_minutes: Int
  },
  created_at: Date,
  updated_at: Date
}
```

### Índices

- `material_id` (simple)
- `metadata.subject` (simple)
- `metadata.grade` (simple)
- `metadata.difficulty` (simple)
- `created_at` (desc)

### Validaciones

- Campos requeridos: `material_id`, `questions`, `metadata`, `created_at`, `updated_at`
- `question_type` debe ser uno de: `multiple_choice`, `true_false`, `short_answer`
- `difficulty` debe ser uno de: `easy`, `medium`, `hard`

---

## 2️⃣ material_content

**Propósito:** Almacena contenido procesado de materiales educativos (texto extraído, estructura parseada).

**Relacionada con PostgreSQL:** La tabla `materials` almacena metadata de archivos.

### Estructura

```javascript
{
  _id: ObjectId,
  material_id: String,              // UUID del material en PostgreSQL (único)
  content_type: String,             // "pdf_extracted" | "video_transcript" | "document_parsed" | "slides_extracted"
  raw_text: String,                 // Texto crudo extraído
  structured_content: {
    title: String,
    sections: [
      {
        section_index: Int,
        heading: String,
        content: String,
        page_number: Int
      }
    ],
    summary: String,                // Resumen generado por IA
    key_concepts: [String]          // Conceptos clave extraídos
  },
  processing_info: {
    processor_version: String,
    processed_at: Date,
    processing_duration_ms: Int,
    page_count: Int,
    word_count: Int
  },
  created_at: Date,
  updated_at: Date
}
```

### Índices

- `material_id` (único)
- `content_type` (simple)
- `created_at` (desc)
- `processing_info.processed_at` (desc)
- **Full-text search** en `raw_text`, `structured_content.summary`, `structured_content.key_concepts` (idioma: español)

### Validaciones

- Campos requeridos: `material_id`, `content_type`, `created_at`, `updated_at`
- `content_type` debe ser uno de: `pdf_extracted`, `video_transcript`, `document_parsed`, `slides_extracted`
- `material_id` debe ser único

---

## 3️⃣ assessment_attempt_result

**Propósito:** Almacena resultados detallados y respuestas de intentos de assessment.

**Relacionada con PostgreSQL:** `assessment_attempt` (metadata), `assessment_attempt_answer` (respuestas individuales).

### Estructura

```javascript
{
  _id: ObjectId,
  attempt_id: String,               // UUID del intento en PostgreSQL (único)
  student_id: String,               // UUID del estudiante
  assessment_id: String,            // UUID del assessment
  answers: [
    {
      question_index: Int,
      question_text: String,        // Snapshot de la pregunta
      selected_option_index: Int,
      selected_option_text: String,
      correct_option_index: Int,
      is_correct: Boolean,
      time_spent_seconds: Int,
      answered_at: Date
    }
  ],
  score: {
    correct_count: Int,
    incorrect_count: Int,
    total_questions: Int,
    percentage: Double              // 0-100
  },
  time_tracking: {
    total_time_seconds: Int,
    average_time_per_question: Double
  },
  started_at: Date,
  submitted_at: Date,
  created_at: Date
}
```

### Índices

- `attempt_id` (único)
- `student_id` (simple)
- `assessment_id` (simple)
- `student_id + assessment_id` (compuesto)
- `submitted_at` (desc)
- `score.percentage` (desc)

### Validaciones

- Campos requeridos: `attempt_id`, `student_id`, `assessment_id`, `answers`, `score`, `started_at`, `submitted_at`, `created_at`
- `score.percentage` debe estar entre 0 y 100

---

## 4️⃣ audit_logs

**Propósito:** Logs de auditoría de eventos importantes del sistema.

**TTL:** Documentos se eliminan automáticamente después de 90 días.

### Estructura

```javascript
{
  _id: ObjectId,
  event_type: String,               // Ver enum abajo
  actor_id: String,                 // UUID del usuario o "system"
  actor_type: String,               // "user" | "system" | "api" | "worker"
  resource_type: String,            // "user" | "school" | "material" | "assessment" | etc.
  resource_id: String,
  action: String,                   // "create" | "read" | "update" | "delete" | etc.
  details: {
    ip_address: String,
    user_agent: String,
    changes: Object,                // Cambios antes/después
    metadata: Object,
    error: Object                   // Si la acción falló
  },
  severity: String,                 // "info" | "warning" | "error" | "critical"
  timestamp: Date,
  session_id: String,
  request_id: String
}
```

### Event Types

- User: `user.created`, `user.updated`, `user.deleted`, `user.login`, `user.logout`
- School: `school.created`, `school.updated`, `school.deleted`
- Material: `material.uploaded`, `material.updated`, `material.deleted`, `material.processed`
- Assessment: `assessment.generated`, `assessment.published`, `assessment.archived`
- Attempt: `attempt.started`, `attempt.submitted`, `attempt.graded`
- Membership: `membership.created`, `membership.updated`, `membership.deleted`
- Permission: `permission.granted`, `permission.revoked`
- System: `system.backup`, `system.restore`, `system.migration`

### Índices

- `timestamp` (desc)
- `event_type + timestamp` (compuesto)
- `actor_id + timestamp` (compuesto)
- `resource_type + resource_id` (compuesto)
- `severity + timestamp` (compuesto)
- `session_id` (simple)
- `request_id` (simple)
- **TTL:** `timestamp` expira después de 7,776,000 segundos (90 días)

### Validaciones

- Campos requeridos: `event_type`, `actor_id`, `timestamp`, `resource_type`
- Ver enums en migración para valores permitidos

---

## 5️⃣ notifications

**Propósito:** Notificaciones para usuarios (in-app, push, email).

**TTL:** Notificaciones archivadas se eliminan después de 30 días.

### Estructura

```javascript
{
  _id: ObjectId,
  user_id: String,                  // UUID del usuario
  notification_type: String,        // Ver enum abajo
  title: String,
  message: String,
  priority: String,                 // "low" | "medium" | "high" | "urgent"
  category: String,                 // "academic" | "administrative" | "social" | "system"
  data: {
    resource_type: String,
    resource_id: String,
    action_url: String,
    action_label: String,
    metadata: Object
  },
  delivery: {
    in_app: {
      enabled: Boolean,
      delivered_at: Date
    },
    push: {
      enabled: Boolean,
      sent_at: Date,
      delivered_at: Date,
      error: String
    },
    email: {
      enabled: Boolean,
      sent_at: Date,
      delivered_at: Date,
      error: String
    }
  },
  is_read: Boolean,
  read_at: Date,
  is_archived: Boolean,
  archived_at: Date,
  expires_at: Date,                 // Opcional
  created_at: Date
}
```

### Notification Types

- Assessment: `assessment.ready`, `assessment.graded`
- Material: `material.uploaded`, `material.processed`, `material.shared`
- Membership: `membership.added`, `membership.removed`
- General: `deadline.approaching`, `achievement.unlocked`
- System: `system.announcement`, `system.maintenance`

### Índices

- `user_id + created_at` (compuesto, desc)
- `user_id + is_read` (compuesto)
- `notification_type` (simple)
- `priority + created_at` (compuesto)
- `created_at` (desc)
- `data.resource_type + data.resource_id` (compuesto)
- **TTL:** `expires_at` expira inmediatamente
- **TTL:** `archived_at` expira después de 2,592,000 segundos (30 días)

### Validaciones

- Campos requeridos: `user_id`, `notification_type`, `title`, `is_read`, `created_at`
- Ver enums en migración para valores permitidos

---

## 6️⃣ analytics_events

**Propósito:** Eventos de analítica y comportamiento de usuarios.

**TTL:** Eventos se eliminan automáticamente después de 365 días.

### Estructura

```javascript
{
  _id: ObjectId,
  event_name: String,               // Ver enum abajo
  user_id: String,                  // UUID o null para eventos anónimos
  session_id: String,
  timestamp: Date,
  properties: {
    page_path: String,
    page_title: String,
    resource_id: String,
    resource_type: String,
    duration_seconds: Int,
    search_query: String,
    search_results_count: Int,
    button_label: String,
    error_message: String,
    custom_data: Object
  },
  device: {
    platform: String,               // "web" | "ios" | "android"
    os: String,
    os_version: String,
    browser: String,
    browser_version: String,
    device_type: String,            // "mobile" | "tablet" | "desktop"
    screen_resolution: String
  },
  location: {
    ip_address: String,             // Anonimizado
    country: String,                // ISO 3166-1 alpha-2
    city: String,
    timezone: String
  },
  context: {
    school_id: String,
    academic_unit_id: String,
    user_role: String,              // "admin" | "teacher" | "student" | "guardian"
    ab_test_variant: String
  }
}
```

### Event Names

- Navigation: `page.view`
- Material: `material.view`, `material.download`, `material.search`
- Assessment: `assessment.start`, `assessment.complete`, `assessment.abandon`
- Question: `question.answer`, `question.skip`
- Video: `video.play`, `video.pause`, `video.complete`
- Session: `session.start`, `session.end`
- Interaction: `feature.click`, `error.occurred`
- Search: `search.performed`, `filter.applied`

### Índices

- `timestamp` (desc)
- `event_name + timestamp` (compuesto)
- `user_id + timestamp` (compuesto)
- `session_id + timestamp` (compuesto)
- `properties.resource_type + properties.resource_id` (compuesto)
- `context.school_id + timestamp` (compuesto)
- `device.platform + timestamp` (compuesto)
- `event_name + context.school_id + timestamp` (compuesto)
- **TTL:** `timestamp` expira después de 31,536,000 segundos (365 días)

### Validaciones

- Campos requeridos: `event_name`, `timestamp`
- Ver enums en migración para valores permitidos

---

## 🔄 Relación PostgreSQL ↔️ MongoDB

### Patrón Híbrido

EduGo usa un **patrón híbrido** donde PostgreSQL y MongoDB trabajan juntos:

```
PostgreSQL (Metadata)          MongoDB (Content)
┌─────────────────┐           ┌──────────────────────┐
│ materials       │           │ material_content     │
│ ├─ id (UUID)    │──────────▶│ ├─ material_id       │
│ ├─ title        │           │ ├─ raw_text          │
│ ├─ file_url     │           │ └─ structured_content│
│ └─ status       │           └──────────────────────┘
└─────────────────┘

┌─────────────────┐           ┌──────────────────────┐
│ assessment      │           │ material_assessment  │
│ ├─ id (UUID)    │           │ ├─ _id (ObjectId)    │
│ ├─ material_id  │           │ ├─ material_id       │
│ └─ mongo_doc_id │──────────▶│ ├─ questions[]       │
└─────────────────┘           │ └─ metadata          │
                              └──────────────────────┘

┌─────────────────┐           ┌────────────────────────┐
│ assessment_     │           │ assessment_attempt_    │
│   attempt       │           │   result               │
│ ├─ id (UUID)    │──────────▶│ ├─ attempt_id          │
│ ├─ student_id   │           │ ├─ answers[]           │
│ └─ score        │           │ └─ detailed_results    │
└─────────────────┘           └────────────────────────┘
```

### ¿Cuándo usar MongoDB?

✅ **Usar MongoDB para:**
- Contenido largo y no estructurado (texto extraído, JSONs grandes)
- Datos con esquema flexible (preguntas de assessment pueden variar)
- Alto volumen de escritura (logs, eventos, analítica)
- Datos que no requieren JOINs complejos
- TTL automático (auto-delete de logs antiguos)

❌ **NO usar MongoDB para:**
- Relaciones complejas (foreign keys)
- Transacciones ACID críticas
- Datos que requieren integridad referencial estricta
- Queries que requieren múltiples JOINs

---

## 🛠️ Uso de Migraciones

### Ejecutar Migraciones

```bash
# Ir al directorio database
cd database

# Ejecutar migraciones MongoDB (Sprint-04)
go run mongodb_migrate.go up

# Ver estado de migraciones
go run mongodb_migrate.go status

# Revertir última migración
go run mongodb_migrate.go down
```

### Crear Nueva Migración

```bash
# Crear nueva migración
go run mongodb_migrate.go create "add_field_to_collection"

# Genera:
# - migrations/mongodb/007_add_field_to_collection.up.js
# - migrations/mongodb/007_add_field_to_collection.down.js
```

---

## 📋 Cargar Seeds

Después de ejecutar migraciones, cargar datos de prueba:

```bash
cd seeds/mongodb

# Cargar todos los seeds
mongosh --host localhost:27017/edugo < assessments.js
mongosh --host localhost:27017/edugo < material_content.js
mongosh --host localhost:27017/edugo < assessment_attempt_result.js
mongosh --host localhost:27017/edugo < audit_logs.js
mongosh --host localhost:27017/edugo < notifications.js
mongosh --host localhost:27017/edugo < analytics_events.js
```

O usar el script helper:

```bash
# Crear script para cargar todos los seeds
for file in seeds/mongodb/*.js; do
  mongosh --host localhost:27017/edugo < "$file"
done
```

---

## ✅ Checklist para Nuevas Colecciones

- [ ] Crear migración UP en `database/migrations/mongodb/00X_create_*.up.js`
- [ ] Crear migración DOWN en `database/migrations/mongodb/00X_create_*.down.js`
- [ ] Definir validación de esquema JSON
- [ ] Crear índices necesarios (incluyendo TTL si aplica)
- [ ] Agregar colección a este documento (MONGODB_SCHEMA.md)
- [ ] Documentar relación con PostgreSQL si existe
- [ ] Crear seeds en `seeds/mongodb/*.js`
- [ ] Testear migración UP y DOWN localmente
- [ ] Commit en rama `dev` de infrastructure
- [ ] PR y merge

---

## 🔍 Consultas Comunes

### Buscar assessment por material

```javascript
db.material_assessment.findOne({ material_id: "66666666-6666-6666-6666-666666666666" })
```

### Ver resultados de un estudiante

```javascript
db.assessment_attempt_result.find({ student_id: "33333333-3333-3333-3333-333333333333" }).sort({ submitted_at: -1 })
```

### Buscar en contenido de materiales (full-text)

```javascript
db.material_content.find({ $text: { $search: "física cuántica" } })
```

### Ver logs de auditoría de un usuario

```javascript
db.audit_logs.find({ actor_id: "11111111-1111-1111-1111-111111111111" }).sort({ timestamp: -1 }).limit(50)
```

### Notificaciones no leídas de un usuario

```javascript
db.notifications.find({ user_id: "33333333-3333-3333-3333-333333333333", is_read: false }).sort({ created_at: -1 })
```

### Analítica: eventos por tipo en las últimas 24 horas

```javascript
db.analytics_events.aggregate([
  { $match: { timestamp: { $gte: new Date(Date.now() - 24*60*60*1000) } } },
  { $group: { _id: "$event_name", count: { $sum: 1 } } },
  { $sort: { count: -1 } }
])
```

---

**Última actualización:** 16 de Noviembre, 2025
**Mantenedor:** Equipo EduGo
**Versión de Schema:** 1.0.0
