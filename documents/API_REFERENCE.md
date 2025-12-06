# 📡 API Reference - EduGo

Referencia completa de los endpoints esperados para cada API del ecosistema EduGo.

> **Nota:** Este documento describe los endpoints **esperados** que las APIs deben implementar consumiendo la infraestructura compartida. No son endpoints de este repositorio de infraestructura.

---

## 📋 Índice

1. [API Mobile](#-api-mobile)
2. [API Administración](#-api-administración)
3. [Worker](#-worker-internal)
4. [Códigos de Error Comunes](#-códigos-de-error-comunes)
5. [Autenticación](#-autenticación)

---

## 📱 API Mobile

**Base URL:** `https://api.edugo.com/mobile/v1`  
**Autenticación:** Bearer JWT  
**Usuarios:** Docentes, Estudiantes, Apoderados

### Autenticación

#### POST /auth/login

Iniciar sesión con credenciales.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "token_type": "Bearer",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "teacher"
  }
}
```

**Response 401:**
```json
{
  "error": "invalid_credentials",
  "message": "Email o contraseña incorrectos"
}
```

---

#### POST /auth/refresh

Renovar access token.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

---

#### POST /auth/logout

Cerrar sesión (invalidar tokens).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response 204:** No Content

---

### Usuarios

#### GET /users/me

Obtener perfil del usuario autenticado.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response 200:**
```json
{
  "id": "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
  "email": "teacher@school.com",
  "first_name": "María",
  "last_name": "González",
  "role": "teacher",
  "is_active": true,
  "email_verified": true,
  "memberships": [
    {
      "id": "m1111111-1111-1111-1111-111111111111",
      "school": {
        "id": "s2222222-2222-2222-2222-222222222222",
        "name": "Colegio San Martín",
        "code": "CSM001"
      },
      "academic_unit": {
        "id": "au333333-3333-3333-3333-333333333333",
        "name": "3° Medio A",
        "type": "class"
      },
      "role": "teacher",
      "enrolled_at": "2024-03-01T00:00:00Z"
    }
  ],
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

#### PATCH /users/me

Actualizar perfil.

**Request:**
```json
{
  "first_name": "María José",
  "phone": "+56912345678"
}
```

**Response 200:** Usuario actualizado

---

### Materiales

#### GET /materials

Listar materiales del usuario.

**Query Parameters:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `school_id` | uuid | Filtrar por escuela |
| `academic_unit_id` | uuid | Filtrar por unidad académica |
| `status` | string | Filtrar por estado (uploaded, processing, ready, failed) |
| `subject` | string | Filtrar por materia |
| `page` | int | Página (default: 1) |
| `limit` | int | Items por página (default: 20, max: 100) |

**Response 200:**
```json
{
  "data": [
    {
      "id": "mat11111-1111-1111-1111-111111111111",
      "title": "Introducción a Java",
      "description": "Material sobre programación orientada a objetos",
      "subject": "Informática",
      "grade": "3° Medio",
      "file_url": "https://edugo-materials.s3.amazonaws.com/...",
      "file_type": "application/pdf",
      "file_size_bytes": 1048576,
      "status": "ready",
      "has_assessment": true,
      "assessment_id": "ass22222-2222-2222-2222-222222222222",
      "uploaded_by": {
        "id": "teacher-uuid",
        "name": "María González"
      },
      "created_at": "2024-06-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

---

#### POST /materials

Subir nuevo material (solo docentes).

**Headers:**
```
Content-Type: multipart/form-data
```

**Request:**
```
file: <binary>
title: "Introducción a Java"
description: "Material sobre POO"
subject: "Informática"
grade: "3° Medio"
school_id: "s2222222-2222-2222-2222-222222222222"
academic_unit_id: "au333333-3333-3333-3333-333333333333" (opcional)
is_public: false
```

**Response 201:**
```json
{
  "id": "mat11111-1111-1111-1111-111111111111",
  "title": "Introducción a Java",
  "status": "uploaded",
  "message": "Material subido. El procesamiento comenzará en breve."
}
```

---

#### GET /materials/:id

Obtener detalle de material.

**Response 200:**
```json
{
  "id": "mat11111-1111-1111-1111-111111111111",
  "title": "Introducción a Java",
  "description": "Material sobre programación orientada a objetos",
  "subject": "Informática",
  "grade": "3° Medio",
  "file_url": "https://edugo-materials.s3.amazonaws.com/...",
  "file_type": "application/pdf",
  "file_size_bytes": 1048576,
  "status": "ready",
  "processing_started_at": "2024-06-15T10:31:00Z",
  "processing_completed_at": "2024-06-15T10:32:30Z",
  "summary": {
    "text": "Este material cubre los fundamentos de la programación orientada a objetos...",
    "key_points": [
      "Conceptos básicos de POO",
      "Clases y objetos",
      "Herencia y polimorfismo",
      "Encapsulación"
    ],
    "word_count": 150
  },
  "assessment": {
    "id": "ass22222-2222-2222-2222-222222222222",
    "questions_count": 10,
    "status": "published",
    "total_attempts": 45
  },
  "uploaded_by": {
    "id": "teacher-uuid",
    "name": "María González"
  },
  "school": {
    "id": "school-uuid",
    "name": "Colegio San Martín"
  },
  "created_at": "2024-06-15T10:30:00Z",
  "updated_at": "2024-06-15T10:32:30Z"
}
```

---

#### DELETE /materials/:id

Eliminar material (solo dueño).

**Response 204:** No Content

---

### Assessments

#### GET /assessments/:id

Obtener assessment para rendirlo.

**Response 200:**
```json
{
  "id": "ass22222-2222-2222-2222-222222222222",
  "material_id": "mat11111-1111-1111-1111-111111111111",
  "title": "Quiz: Introducción a Java",
  "questions_count": 10,
  "total_points": 100,
  "time_limit_minutes": 30,
  "max_attempts": 3,
  "pass_threshold": 60,
  "status": "published",
  "questions": [
    {
      "index": 0,
      "question_id": "q1111111-1111-1111-1111-111111111111",
      "question_text": "¿Qué es la encapsulación en POO?",
      "question_type": "multiple_choice",
      "options": [
        {"id": "opt1", "text": "Ocultar detalles de implementación"},
        {"id": "opt2", "text": "Heredar de otra clase"},
        {"id": "opt3", "text": "Crear múltiples instancias"},
        {"id": "opt4", "text": "Ejecutar código en paralelo"}
      ],
      "points": 10,
      "difficulty": "medium"
    }
  ],
  "my_attempts": [
    {
      "id": "attempt-uuid",
      "score": 80,
      "status": "graded",
      "submitted_at": "2024-06-16T14:30:00Z"
    }
  ],
  "remaining_attempts": 2
}
```

---

#### POST /assessments/:id/start

Iniciar intento de assessment.

**Response 201:**
```json
{
  "attempt_id": "att33333-3333-3333-3333-333333333333",
  "started_at": "2024-06-16T14:00:00Z",
  "expires_at": "2024-06-16T14:30:00Z"
}
```

---

#### POST /attempts/:id/answer

Enviar respuesta a pregunta.

**Request:**
```json
{
  "question_index": 0,
  "answer": "opt1"
}
```

**Response 200:**
```json
{
  "saved": true,
  "question_index": 0
}
```

---

#### POST /attempts/:id/submit

Finalizar y enviar intento.

**Response 200:**
```json
{
  "attempt_id": "att33333-3333-3333-3333-333333333333",
  "score": 80.00,
  "correct_answers": 8,
  "total_questions": 10,
  "passed": true,
  "results": [
    {
      "question_index": 0,
      "question_text": "¿Qué es la encapsulación en POO?",
      "student_answer": "opt1",
      "correct_answer": "opt1",
      "is_correct": true,
      "explanation": "La encapsulación es el mecanismo de ocultar los detalles de implementación...",
      "points_earned": 10
    }
  ],
  "submitted_at": "2024-06-16T14:25:00Z"
}
```

---

### Progreso

#### GET /progress

Obtener progreso del estudiante.

**Query Parameters:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `school_id` | uuid | Filtrar por escuela |
| `subject` | string | Filtrar por materia |
| `from` | date | Fecha inicio |
| `to` | date | Fecha fin |

**Response 200:**
```json
{
  "overall": {
    "total_assessments_taken": 25,
    "average_score": 78.5,
    "total_time_minutes": 450,
    "improvement_trend": "+5.2%"
  },
  "by_subject": [
    {
      "subject": "Informática",
      "assessments_taken": 10,
      "average_score": 82.3,
      "best_score": 95,
      "worst_score": 65
    },
    {
      "subject": "Matemáticas",
      "assessments_taken": 8,
      "average_score": 75.0,
      "best_score": 90,
      "worst_score": 55
    }
  ],
  "recent_activity": [
    {
      "date": "2024-06-16",
      "assessment_title": "Quiz: Introducción a Java",
      "score": 80,
      "subject": "Informática"
    }
  ]
}
```

---

## 🖥️ API Administración

**Base URL:** `https://api.edugo.com/admin/v1`  
**Autenticación:** Bearer JWT  
**Usuarios:** Administradores, Coordinadores

### Escuelas

#### GET /schools

Listar escuelas (superadmin).

#### GET /schools/:id

Obtener detalle de escuela.

**Response 200:**
```json
{
  "id": "s2222222-2222-2222-2222-222222222222",
  "name": "Colegio San Martín",
  "code": "CSM001",
  "address": "Av. Principal 123",
  "city": "Santiago",
  "country": "Chile",
  "phone": "+56212345678",
  "email": "contacto@csm.cl",
  "subscription_tier": "premium",
  "max_teachers": 50,
  "max_students": 1000,
  "current_teachers": 35,
  "current_students": 850,
  "is_active": true,
  "metadata": {
    "logo_url": "https://...",
    "website": "https://csm.cl"
  },
  "created_at": "2023-01-01T00:00:00Z"
}
```

#### PATCH /schools/:id

Actualizar escuela.

---

### Unidades Académicas

#### GET /schools/:school_id/academic-units

Listar unidades académicas (árbol jerárquico).

**Response 200:**
```json
{
  "data": [
    {
      "id": "au111111-1111-1111-1111-111111111111",
      "name": "Enseñanza Media",
      "code": "EM",
      "type": "grade",
      "level": null,
      "children": [
        {
          "id": "au222222-2222-2222-2222-222222222222",
          "name": "3° Medio",
          "code": "3M",
          "type": "grade",
          "level": "3",
          "children": [
            {
              "id": "au333333-3333-3333-3333-333333333333",
              "name": "3° Medio A",
              "code": "3MA",
              "type": "class",
              "students_count": 35,
              "teachers_count": 8
            },
            {
              "id": "au444444-4444-4444-4444-444444444444",
              "name": "3° Medio B",
              "code": "3MB",
              "type": "class",
              "students_count": 32,
              "teachers_count": 8
            }
          ]
        }
      ]
    }
  ]
}
```

#### POST /schools/:school_id/academic-units

Crear unidad académica.

**Request:**
```json
{
  "name": "3° Medio C",
  "code": "3MC",
  "type": "class",
  "parent_unit_id": "au222222-2222-2222-2222-222222222222",
  "academic_year": 2024
}
```

---

### Membresías

#### GET /schools/:school_id/memberships

Listar membresías de una escuela.

**Query Parameters:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `role` | string | Filtrar por rol |
| `academic_unit_id` | uuid | Filtrar por unidad |
| `is_active` | bool | Filtrar activos/inactivos |

#### POST /schools/:school_id/memberships

Matricular usuario.

**Request:**
```json
{
  "user_id": "u4444444-4444-4444-4444-444444444444",
  "academic_unit_id": "au333333-3333-3333-3333-333333333333",
  "role": "student"
}
```

**Response 201:**
```json
{
  "id": "m5555555-5555-5555-5555-555555555555",
  "user_id": "u4444444-4444-4444-4444-444444444444",
  "school_id": "s2222222-2222-2222-2222-222222222222",
  "academic_unit_id": "au333333-3333-3333-3333-333333333333",
  "role": "student",
  "enrolled_at": "2024-06-01T00:00:00Z",
  "message": "Estudiante matriculado exitosamente"
}
```

#### DELETE /schools/:school_id/memberships/:id

Dar de baja membresía.

---

### Usuarios

#### GET /users

Listar usuarios del sistema.

#### POST /users

Crear usuario.

#### GET /users/:id

Obtener usuario.

#### PATCH /users/:id

Actualizar usuario.

---

### Reportes

#### GET /reports/assessments

Reporte de assessments.

**Query Parameters:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `school_id` | uuid | Escuela |
| `academic_unit_id` | uuid | Unidad académica |
| `from` | date | Fecha inicio |
| `to` | date | Fecha fin |
| `group_by` | string | subject, teacher, class |

**Response 200:**
```json
{
  "period": {
    "from": "2024-01-01",
    "to": "2024-06-30"
  },
  "summary": {
    "total_assessments": 150,
    "total_attempts": 4500,
    "average_score": 72.5,
    "pass_rate": 68.2
  },
  "by_subject": [
    {
      "subject": "Matemáticas",
      "assessments": 45,
      "attempts": 1350,
      "average_score": 68.3,
      "pass_rate": 62.1
    }
  ]
}
```

---

## ⚙️ Worker (Internal)

El worker no expone API HTTP. Se comunica vía RabbitMQ.

### Eventos Consumidos

| Evento | Acción |
|--------|--------|
| `material.uploaded` | Procesar material, generar assessment y resumen |
| `material.deleted` | Limpiar archivos y documentos relacionados |

### Eventos Publicados

| Evento | Trigger |
|--------|---------|
| `assessment.generated` | Assessment generado exitosamente |

---

## ❌ Códigos de Error Comunes

| Código | Error | Descripción |
|--------|-------|-------------|
| 400 | `bad_request` | Request malformado |
| 401 | `unauthorized` | Token inválido o expirado |
| 403 | `forbidden` | Sin permisos para la acción |
| 404 | `not_found` | Recurso no encontrado |
| 409 | `conflict` | Conflicto (ej: duplicado) |
| 422 | `validation_error` | Error de validación |
| 429 | `rate_limit` | Demasiadas requests |
| 500 | `internal_error` | Error interno |

### Formato de Error

```json
{
  "error": "validation_error",
  "message": "Error de validación",
  "details": [
    {
      "field": "email",
      "message": "Email inválido"
    }
  ],
  "request_id": "req-uuid-for-tracing"
}
```

---

## 🔐 Autenticación

### JWT Payload

```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "role": "teacher",
  "school_ids": ["school-uuid-1", "school-uuid-2"],
  "exp": 1704067200,
  "iat": 1704066300,
  "jti": "token-unique-id"
}
```

### Headers Requeridos

```
Authorization: Bearer <access_token>
Content-Type: application/json
Accept: application/json
X-Request-ID: <uuid> (opcional, para tracing)
```

---

## 📊 Rate Limiting

| Endpoint | Límite |
|----------|--------|
| `/auth/*` | 10 req/min |
| `POST /materials` | 5 req/min |
| `GET /*` | 100 req/min |
| `POST /*` | 30 req/min |

Headers de respuesta:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067260
```

---

**Última actualización:** Diciembre 2024  
**Versión API:** v1
