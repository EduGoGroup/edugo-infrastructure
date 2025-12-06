# 🗄️ Database Schema - EduGo

Este documento describe el modelo de datos completo de EduGo, incluyendo PostgreSQL (relacional) y MongoDB (documentos).

---

## 📊 Diagrama Entidad-Relación (PostgreSQL)

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                    DIAGRAMA ER                                           │
└─────────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────┐          ┌─────────────────┐          ┌─────────────────┐
│     users       │          │     schools     │          │ academic_units  │
├─────────────────┤          ├─────────────────┤          ├─────────────────┤
│ id (PK)         │          │ id (PK)         │◀────────┐│ id (PK)         │
│ email           │          │ name            │         ││ parent_unit_id  │──┐
│ password_hash   │          │ code            │         ││ school_id (FK)  │──┤
│ first_name      │          │ address         │         ││ name            │  │
│ last_name       │          │ city            │         ││ code            │  │
│ role            │          │ country         │         ││ type            │  │
│ is_active       │          │ phone           │         ││ description     │  │
│ email_verified  │          │ email           │         ││ level           │  │
│ created_at      │          │ metadata (JSON) │         ││ academic_year   │  │
│ updated_at      │          │ is_active       │         ││ metadata (JSON) │  │
│ deleted_at      │          │ subscription_tier│        ││ is_active       │  │
└────────┬────────┘          │ max_teachers    │         ││ created_at      │  │
         │                   │ max_students    │         ││ updated_at      │  │
         │                   │ created_at      │         ││ deleted_at      │  │
         │                   │ updated_at      │         │└────────┬────────┘  │
         │                   │ deleted_at      │         │         │           │
         │                   └────────┬────────┘         │         │           │
         │                            │                  │         │           │
         │                            │                  │         └───────────┘
         │                            │                  │         (self-reference)
         │                            │                  │
         │     ┌─────────────────┐    │                  │
         │     │   memberships   │    │                  │
         │     ├─────────────────┤    │                  │
         └────▶│ id (PK)         │    │                  │
               │ user_id (FK)    │────┘                  │
               │ school_id (FK)  │───────────────────────┘
               │ academic_unit_id│◀──────────────────────┘
               │ role            │
               │ metadata (JSON) │
               │ is_active       │
               │ enrolled_at     │
               │ withdrawn_at    │
               │ created_at      │
               │ updated_at      │
               └─────────────────┘

┌─────────────────┐          ┌─────────────────┐
│    materials    │          │   assessment    │
├─────────────────┤          ├─────────────────┤
│ id (PK)         │◀────────┐│ id (PK)         │
│ school_id (FK)  │─────────┤│ material_id(FK) │────┐
│ uploaded_by_    │         ││ mongo_document_id│    │
│   teacher_id    │─────────┘│ questions_count │    │
│ academic_unit_id│──────────│ total_questions │    │
│ title           │          │ title           │    │
│ description     │          │ pass_threshold  │    │
│ subject         │          │ max_attempts    │    │
│ grade           │          │ time_limit_min  │    │
│ file_url        │          │ status          │    │
│ file_type       │          │ created_at      │    │
│ file_size_bytes │          │ updated_at      │    │
│ status          │          │ deleted_at      │    │
│ processing_     │          └────────┬────────┘    │
│   started_at    │                   │             │
│ processing_     │                   │             │
│   completed_at  │                   ▼             │
│ is_public       │          ┌─────────────────┐    │
│ created_at      │          │assessment_attempt│   │
│ updated_at      │          ├─────────────────┤    │
│ deleted_at      │          │ id (PK)         │    │
└─────────────────┘          │ assessment_id   │────┤
                             │ student_id (FK) │    │
                             │ started_at      │    │
                             │ submitted_at    │    │
                             │ score           │    │
                             │ status          │    │
                             │ metadata (JSON) │    │
                             │ created_at      │    │
                             │ updated_at      │    │
                             └────────┬────────┘    │
                                      │             │
                                      ▼             │
                             ┌────────────────────┐ │
                             │assessment_attempt_ │ │
                             │      answer        │ │
                             ├────────────────────┤ │
                             │ id (PK)            │ │
                             │ attempt_id (FK)    │─┘
                             │ question_index     │
                             │ student_answer     │
                             │ is_correct         │
                             │ answered_at        │
                             │ created_at         │
                             │ updated_at         │
                             └────────────────────┘
```

---

## 🐘 PostgreSQL Tables

### 1. `users`

Usuarios del sistema (docentes, estudiantes, apoderados, administradores).

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `email` | VARCHAR(255) | NO | Email único |
| `password_hash` | VARCHAR(255) | NO | Hash bcrypt |
| `first_name` | VARCHAR(100) | NO | Nombre |
| `last_name` | VARCHAR(100) | NO | Apellido |
| `role` | VARCHAR(20) | NO | admin, teacher, student, guardian |
| `is_active` | BOOLEAN | NO | Estado activo |
| `email_verified` | BOOLEAN | NO | Email verificado |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |
| `deleted_at` | TIMESTAMP | YES | Soft delete |

**Índices:**
- `idx_users_email` (UNIQUE)
- `idx_users_role`

---

### 2. `schools`

Instituciones educativas registradas.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `name` | VARCHAR(200) | NO | Nombre escuela |
| `code` | VARCHAR(50) | NO | Código único |
| `address` | VARCHAR(500) | YES | Dirección |
| `city` | VARCHAR(100) | YES | Ciudad |
| `country` | VARCHAR(100) | NO | País |
| `phone` | VARCHAR(20) | YES | Teléfono |
| `email` | VARCHAR(255) | YES | Email contacto |
| `metadata` | JSONB | YES | Metadata extensible |
| `is_active` | BOOLEAN | NO | Estado activo |
| `subscription_tier` | VARCHAR(20) | NO | free, basic, premium, enterprise |
| `max_teachers` | INTEGER | NO | Límite docentes |
| `max_students` | INTEGER | NO | Límite estudiantes |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |
| `deleted_at` | TIMESTAMP | YES | Soft delete |

**Índices:**
- `idx_schools_code` (UNIQUE)
- `idx_schools_country`

---

### 3. `academic_units`

Unidades académicas jerárquicas (grados, cursos, secciones).

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `parent_unit_id` | UUID | YES | FK a parent (jerarquía) |
| `school_id` | UUID | NO | FK a schools |
| `name` | VARCHAR(200) | NO | Nombre unidad |
| `code` | VARCHAR(50) | NO | Código |
| `type` | VARCHAR(50) | NO | school, grade, class, section, club, department |
| `description` | TEXT | YES | Descripción |
| `level` | VARCHAR(50) | YES | Nivel educativo |
| `academic_year` | INTEGER | NO | Año académico (0 = sin año) |
| `metadata` | JSONB | YES | Metadata extensible |
| `is_active` | BOOLEAN | NO | Estado activo |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |
| `deleted_at` | TIMESTAMP | YES | Soft delete |

**Índices:**
- `idx_academic_units_school_id`
- `idx_academic_units_parent_unit_id`
- `idx_academic_units_type`

**Ejemplo de jerarquía:**
```
School (Colegio ABC)
├── Grade (1° Básico)
│   ├── Class (1°A)
│   └── Class (1°B)
└── Grade (2° Básico)
    └── Class (2°A)
```

---

### 4. `memberships`

Relación usuario-escuela-unidad académica.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `user_id` | UUID | NO | FK a users |
| `school_id` | UUID | NO | FK a schools |
| `academic_unit_id` | UUID | YES | FK a academic_units |
| `role` | VARCHAR(20) | NO | teacher, student, guardian, coordinator, admin, assistant |
| `metadata` | JSONB | YES | Metadata extensible |
| `is_active` | BOOLEAN | NO | Estado activo |
| `enrolled_at` | TIMESTAMP | NO | Fecha matrícula |
| `withdrawn_at` | TIMESTAMP | YES | Fecha retiro |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |

**Índices:**
- `idx_memberships_user_school` (user_id, school_id) UNIQUE
- `idx_memberships_academic_unit_id`

---

### 5. `materials`

Materiales educativos subidos por docentes.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `school_id` | UUID | NO | FK a schools |
| `uploaded_by_teacher_id` | UUID | NO | FK a users (docente) |
| `academic_unit_id` | UUID | YES | FK a academic_units |
| `title` | VARCHAR(300) | NO | Título |
| `description` | TEXT | YES | Descripción |
| `subject` | VARCHAR(100) | YES | Materia |
| `grade` | VARCHAR(50) | YES | Grado |
| `file_url` | VARCHAR(1000) | NO | URL S3 |
| `file_type` | VARCHAR(100) | NO | MIME type |
| `file_size_bytes` | BIGINT | NO | Tamaño archivo |
| `status` | VARCHAR(20) | NO | uploaded, processing, ready, failed |
| `processing_started_at` | TIMESTAMP | YES | Inicio procesamiento |
| `processing_completed_at` | TIMESTAMP | YES | Fin procesamiento |
| `is_public` | BOOLEAN | NO | Público/Privado |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |
| `deleted_at` | TIMESTAMP | YES | Soft delete |

**Índices:**
- `idx_materials_school_id`
- `idx_materials_uploaded_by`
- `idx_materials_status`

**Estados del material:**
```
uploaded → processing → ready
                     ↘ failed
```

---

### 6. `assessment`

Metadata de assessments (quizzes) generados.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `material_id` | UUID | NO | FK a materials |
| `mongo_document_id` | VARCHAR(24) | NO | ObjectId de MongoDB |
| `questions_count` | INTEGER | NO | Total preguntas |
| `total_questions` | INTEGER | YES | Sincronizado |
| `title` | VARCHAR(300) | YES | Título |
| `pass_threshold` | INTEGER | YES | % para aprobar (0-100) |
| `max_attempts` | INTEGER | YES | Intentos máximos (NULL = ilimitado) |
| `time_limit_minutes` | INTEGER | YES | Tiempo límite (NULL = sin límite) |
| `status` | VARCHAR(20) | NO | draft, generated, published, archived, closed |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |
| `deleted_at` | TIMESTAMP | YES | Soft delete |

**Índices:**
- `idx_assessment_material_id`
- `idx_assessment_mongo_document_id`
- `idx_assessment_status`

---

### 7. `assessment_attempt`

Intentos de estudiantes en assessments.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `assessment_id` | UUID | NO | FK a assessment |
| `student_id` | UUID | NO | FK a users |
| `started_at` | TIMESTAMP | NO | Inicio intento |
| `submitted_at` | TIMESTAMP | YES | Fin intento |
| `score` | DECIMAL(5,2) | YES | Puntaje (0-100) |
| `status` | VARCHAR(20) | NO | in_progress, submitted, graded, abandoned |
| `metadata` | JSONB | YES | Metadata extensible |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |

**Índices:**
- `idx_attempt_assessment_id`
- `idx_attempt_student_id`
- `idx_attempt_status`

---

### 8. `assessment_attempt_answer`

Respuestas individuales por intento.

| Campo | Tipo | Nullable | Descripción |
|-------|------|----------|-------------|
| `id` | UUID | NO | Primary key |
| `attempt_id` | UUID | NO | FK a assessment_attempt |
| `question_index` | INTEGER | NO | Índice pregunta (0-based) |
| `student_answer` | TEXT | NO | Respuesta del estudiante |
| `is_correct` | BOOLEAN | YES | Correcto/Incorrecto |
| `answered_at` | TIMESTAMP | NO | Timestamp respuesta |
| `created_at` | TIMESTAMP | NO | Fecha creación |
| `updated_at` | TIMESTAMP | NO | Última actualización |

**Índices:**
- `idx_answer_attempt_id`
- `idx_answer_attempt_question` (attempt_id, question_index) UNIQUE

---

## 🍃 MongoDB Collections

### 1. `material_assessment_worker`

Contenido completo de assessments generados por IA.

```javascript
{
  "_id": ObjectId("..."),                    // ID MongoDB
  "material_id": "uuid-string",              // Ref a PostgreSQL
  "questions": [
    {
      "question_id": "q-uuid",               // ID único pregunta
      "question_text": "¿Qué es POO?",       // Texto pregunta
      "question_type": "multiple_choice",    // multiple_choice | true_false | open
      "options": [
        { "option_id": "opt1", "option_text": "Opción 1" },
        { "option_id": "opt2", "option_text": "Opción 2" },
        { "option_id": "opt3", "option_text": "Opción 3" }
      ],
      "correct_answer": "opt3",              // ID respuesta correcta
      "explanation": "Porque...",            // Explicación
      "points": 10,                          // Puntos
      "difficulty": "medium",                // easy | medium | hard
      "tags": ["POO", "conceptos"]           // Tags
    }
  ],
  "total_questions": 10,                     // Total preguntas
  "total_points": 100,                       // Puntos totales
  "version": 1,                              // Versión
  "ai_model": "gpt-4",                       // Modelo usado
  "processing_time_ms": 5200,                // Tiempo procesamiento
  "token_usage": {
    "prompt_tokens": 1200,
    "completion_tokens": 450,
    "total_tokens": 1650
  },
  "metadata": {
    "average_difficulty": "medium",
    "estimated_time_min": 15,
    "source_length": 5420,
    "has_images": false
  },
  "created_at": ISODate("2024-01-01T00:00:00Z"),
  "updated_at": ISODate("2024-01-01T00:00:00Z")
}
```

**Índices:**
- `material_id` (único)
- `created_at`

---

### 2. `material_summary`

Resúmenes generados por IA.

```javascript
{
  "_id": ObjectId("..."),
  "material_id": "uuid-string",
  "summary": "Este material cubre...",       // Resumen texto
  "key_points": [                            // Puntos clave
    "Introducción a POO",
    "Clases y objetos",
    "Herencia y polimorfismo"
  ],
  "language": "es",                          // Idioma detectado
  "word_count": 150,                         // Palabras
  "version": 1,
  "ai_model": "gpt-4",
  "processing_time_ms": 3500,
  "token_usage": {
    "prompt_tokens": 850,
    "completion_tokens": 180,
    "total_tokens": 1030
  },
  "metadata": {
    "source_length": 5420,
    "has_images": false
  },
  "created_at": ISODate("2024-01-01T00:00:00Z"),
  "updated_at": ISODate("2024-01-01T00:00:00Z")
}
```

---

### 3. `material_event`

Log de eventos de procesamiento de materiales.

```javascript
{
  "_id": ObjectId("..."),
  "event_type": "material_uploaded",         // Tipo evento
  "material_id": "uuid-string",
  "user_id": "uuid-string",
  "payload": {                               // Payload flexible
    "filename": "java-poo.pdf",
    "file_size": 1024000,
    "mime_type": "application/pdf"
  },
  "status": "completed",                     // pending | processing | completed | failed
  "retry_count": 0,
  "error_message": null,
  "created_at": ISODate("2024-01-01T00:00:00Z"),
  "updated_at": ISODate("2024-01-01T00:00:00Z")
}
```

---

## 🔗 Relación PostgreSQL ↔ MongoDB

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    RELACIÓN ENTRE BASES DE DATOS                         │
└─────────────────────────────────────────────────────────────────────────┘

PostgreSQL                              MongoDB
┌─────────────────┐                     ┌─────────────────────────────┐
│    assessment   │                     │  material_assessment_worker  │
├─────────────────┤                     ├─────────────────────────────┤
│ mongo_document_ │────────────────────▶│ _id (ObjectId)              │
│       id        │                     │                             │
│                 │                     │ questions: [...]            │
│ questions_count │◀───sincronizado────▶│ total_questions             │
└─────────────────┘                     └─────────────────────────────┘

PostgreSQL                              MongoDB
┌─────────────────┐                     ┌─────────────────────────────┐
│    materials    │                     │     material_summary        │
├─────────────────┤                     ├─────────────────────────────┤
│      id         │────────────────────▶│ material_id (string)        │
│                 │                     │                             │
│    status       │◀───actualizado─────▶│ (genera cuando ready)       │
└─────────────────┘                     └─────────────────────────────┘
```

**Flujo:**
1. Material se crea en PostgreSQL (status: `uploaded`)
2. Worker procesa y crea documento en MongoDB
3. Worker actualiza PostgreSQL con `mongo_document_id` y status `ready`

---

## 📋 Queries Comunes

### PostgreSQL

```sql
-- Obtener materiales de un docente
SELECT * FROM materials 
WHERE uploaded_by_teacher_id = $1 
AND deleted_at IS NULL
ORDER BY created_at DESC;

-- Obtener estudiantes de una unidad académica
SELECT u.* FROM users u
JOIN memberships m ON u.id = m.user_id
WHERE m.academic_unit_id = $1
AND m.role = 'student'
AND m.is_active = true;

-- Obtener intentos de un estudiante en un assessment
SELECT * FROM assessment_attempt
WHERE student_id = $1
AND assessment_id = $2
ORDER BY started_at DESC;
```

### MongoDB

```javascript
// Obtener assessment con preguntas
db.material_assessment_worker.findOne({ material_id: "uuid-string" })

// Obtener resumen de un material
db.material_summary.findOne({ material_id: "uuid-string" })

// Listar eventos de un material
db.material_event.find({ material_id: "uuid-string" }).sort({ created_at: -1 })
```

---

**Última actualización:** Diciembre 2024
