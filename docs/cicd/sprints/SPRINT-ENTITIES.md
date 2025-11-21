# Sprint: Centralizar Entities en Infrastructure

**Proyecto:** edugo-infrastructure  
**Fecha:** 20 de Noviembre, 2025  
**Objetivo:** Crear entities base para PostgreSQL y MongoDB que serán compartidas por todos los proyectos  
**Prioridad:** ALTA - Bloquea adaptación en otros proyectos

---

## 🎯 Contexto

**Problema actual:**
- Entities están duplicadas en cada proyecto (api-mobile, api-administracion, worker)
- Cambios en schema de BD requieren actualizar N proyectos
- Riesgo de discrepancias entre proyectos

**Solución:**
- Centralizar entities en `infrastructure` como "single source of truth"
- Entities = reflejo exacto de tablas/collections de BD
- Cada proyecto importa desde infrastructure

---

## 📊 Inventario de Entities a Crear

### PostgreSQL Entities (14 tablas)

Basado en migraciones existentes en `postgres/migrations/`:

| # | Tabla | Entity | Proyectos que la usan |
|---|-------|--------|----------------------|
| 1 | `users` | `User` | api-mobile, api-administracion |
| 2 | `schools` | `School` | api-administracion |
| 3 | `academic_units` | `AcademicUnit` | api-administracion |
| 4 | `memberships` | `Membership` | api-administracion |
| 5 | `materials` | `Material` | api-mobile, api-administracion |
| 6 | `material_versions` | `MaterialVersion` | api-mobile |
| 7 | `subjects` | `Subject` | api-administracion |
| 8 | `units` | `Unit` | api-administracion |
| 9 | `guardian_relations` | `GuardianRelation` | api-administracion |
| 10 | `assessments` | `Assessment` | api-mobile |
| 11 | `assessment_questions` | `AssessmentQuestion` | api-mobile |
| 12 | `assessment_answers` | `AssessmentAnswer` | api-mobile |
| 13 | `assessment_attempts` | `AssessmentAttempt` | api-mobile |
| 14 | `progress` | `Progress` | api-mobile |

### MongoDB Entities (3 collections)

Basado en migraciones existentes en `mongodb/migrations/`:

| # | Collection | Entity | Proyectos que la usan |
|---|------------|--------|----------------------|
| 1 | `material_assessment` | `MaterialAssessment` | worker |
| 2 | `material_summary` | `MaterialSummary` | worker |
| 3 | `material_event` | `MaterialEvent` | worker |

---

## 🏗️ Estructura Propuesta

```
edugo-infrastructure/
├── postgres/
│   ├── migrations/           # Ya existe
│   ├── testing/              # Ya existe
│   │
│   └── entities/             # ✅ NUEVA CARPETA
│       ├── user.go
│       ├── school.go
│       ├── academic_unit.go
│       ├── membership.go
│       ├── material.go
│       ├── material_version.go
│       ├── subject.go
│       ├── unit.go
│       ├── guardian_relation.go
│       ├── assessment.go
│       ├── assessment_question.go
│       ├── assessment_answer.go
│       ├── assessment_attempt.go
│       └── progress.go
│
└── mongodb/
    ├── migrations/           # Ya existe
    ├── seeds/                # Ya existe
    ├── testing/              # Ya existe
    │
    └── entities/             # ✅ NUEVA CARPETA
        ├── material_assessment.go
        ├── material_summary.go
        └── material_event.go
```

---

## 📋 Tareas del Sprint

### Fase 1: Setup de Estructura

#### Tarea 1.1: Crear carpetas de entities
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
mkdir -p postgres/entities
mkdir -p mongodb/entities
```

**Criterio de éxito:** Carpetas creadas

---

### Fase 2: PostgreSQL Entities

#### Tarea 2.1: Analizar schema de cada tabla
Para cada tabla en `postgres/migrations/*.up.sql`:
1. Leer el schema SQL
2. Identificar campos, tipos, constraints
3. Mapear a tipos Go

**Criterio de éxito:** Documento de mapeo SQL → Go creado

#### Tarea 2.2: Crear entities base de PostgreSQL

Para cada entity, seguir esta estructura:

```go
// postgres/entities/user.go
package entities

import (
    "time"
    "github.com/google/uuid"
)

// User representa la tabla 'users' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD
type User struct {
    ID        uuid.UUID  `db:"id"`
    Email     string     `db:"email"`
    FirstName string     `db:"first_name"`
    LastName  string     `db:"last_name"`
    Role      string     `db:"role"`
    IsActive  bool       `db:"is_active"`
    CreatedAt time.Time  `db:"created_at"`
    UpdatedAt time.Time  `db:"updated_at"`
}

// TableName retorna el nombre de la tabla
func (User) TableName() string {
    return "users"
}
```

**Reglas importantes:**
- ✅ Struct tags `db:` con nombre exacto de columna
- ✅ Tipos Go que mapean correctamente a PostgreSQL
- ✅ Sin lógica de negocio (solo estructura)
- ✅ Sin validaciones (eso va en los proyectos)
- ✅ Comentarios que referencian la tabla SQL

**Lista de entities a crear:**
- [ ] `postgres/entities/user.go`
- [ ] `postgres/entities/school.go`
- [ ] `postgres/entities/academic_unit.go`
- [ ] `postgres/entities/membership.go`
- [ ] `postgres/entities/material.go`
- [ ] `postgres/entities/material_version.go`
- [ ] `postgres/entities/subject.go`
- [ ] `postgres/entities/unit.go`
- [ ] `postgres/entities/guardian_relation.go`
- [ ] `postgres/entities/assessment.go`
- [ ] `postgres/entities/assessment_question.go`
- [ ] `postgres/entities/assessment_answer.go`
- [ ] `postgres/entities/assessment_attempt.go`
- [ ] `postgres/entities/progress.go`

**Criterio de éxito:** 14 entities creados, compilación exitosa

---

### Fase 3: MongoDB Entities

#### Tarea 3.1: Crear entities de MongoDB

Para cada collection, seguir esta estructura:

```go
// mongodb/entities/material_assessment.go
package entities

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// MaterialAssessment representa la collection 'material_assessment' en MongoDB
// Esta entity es el reflejo exacto del schema de BD
type MaterialAssessment struct {
    ID               primitive.ObjectID `bson:"_id,omitempty"`
    MaterialID       string             `bson:"material_id"`
    Questions        []Question         `bson:"questions"`
    TotalQuestions   int                `bson:"total_questions"`
    TotalPoints      int                `bson:"total_points"`
    Version          int                `bson:"version"`
    AIModel          string             `bson:"ai_model"`
    ProcessingTimeMs int                `bson:"processing_time_ms"`
    TokenUsage       *TokenUsage        `bson:"token_usage,omitempty"`
    Metadata         *AssessmentMetadata `bson:"metadata,omitempty"`
    CreatedAt        time.Time          `bson:"created_at"`
    UpdatedAt        time.Time          `bson:"updated_at"`
}

// Question representa una pregunta embebida
type Question struct {
    QuestionID    string   `bson:"question_id"`
    QuestionText  string   `bson:"question_text"`
    QuestionType  string   `bson:"question_type"`
    Options       []Option `bson:"options,omitempty"`
    CorrectAnswer string   `bson:"correct_answer"`
    Explanation   string   `bson:"explanation"`
    Points        int      `bson:"points"`
    Difficulty    string   `bson:"difficulty"`
    Tags          []string `bson:"tags,omitempty"`
}

// Option representa una opción de respuesta
type Option struct {
    OptionID   string `bson:"option_id"`
    OptionText string `bson:"option_text"`
}

// TokenUsage representa metadata de tokens
type TokenUsage struct {
    PromptTokens     int `bson:"prompt_tokens"`
    CompletionTokens int `bson:"completion_tokens"`
    TotalTokens      int `bson:"total_tokens"`
}

// AssessmentMetadata contiene metadata adicional
type AssessmentMetadata struct {
    AverageDifficulty string `bson:"average_difficulty"`
    EstimatedTimeMin  int    `bson:"estimated_time_min"`
}

// CollectionName retorna el nombre de la collection
func (MaterialAssessment) CollectionName() string {
    return "material_assessment"
}
```

**Lista de entities a crear:**
- [ ] `mongodb/entities/material_assessment.go` (con structs embebidos)
- [ ] `mongodb/entities/material_summary.go`
- [ ] `mongodb/entities/material_event.go`

**Criterio de éxito:** 3 entities creados, compilación exitosa

---

### Fase 4: Actualizar go.mod

#### Tarea 4.1: Verificar dependencias

Asegurar que `infrastructure/postgres/go.mod` y `infrastructure/mongodb/go.mod` tengan:

```go
// postgres/go.mod
require (
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9  // Para tipos específicos si es necesario
)

// mongodb/go.mod
require (
    go.mongodb.org/mongo-driver v1.13.1
)
```

**Criterio de éxito:** Dependencies correctas, `go mod tidy` sin errores

---

### Fase 5: Crear README de Entities

#### Tarea 5.1: Documentar uso de entities

Crear `postgres/entities/README.md`:

```markdown
# PostgreSQL Entities

Entities base que reflejan el schema de PostgreSQL.

## Uso

```go
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"

user := &pgentities.User{
    ID:        uuid.New(),
    Email:     "test@example.com",
    FirstName: "John",
    LastName:  "Doe",
    Role:      "student",
    IsActive:  true,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
\```

## Reglas

- **NO agregar lógica de negocio** aquí
- **NO agregar validaciones** aquí
- Solo reflejar estructura de BD
- Para lógica de negocio, usar domain services en tu proyecto
\```
```

Crear `mongodb/entities/README.md` similar.

**Criterio de éxito:** READMEs creados

---

### Fase 6: Testing

#### Tarea 6.1: Tests básicos de entities

Crear tests que validen:
- Struct tags correctos
- TableName()/CollectionName() funciona
- JSON/BSON marshaling funciona

```go
// postgres/entities/user_test.go
package entities_test

import (
    "testing"
    "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestUser_TableName(t *testing.T) {
    u := entities.User{}
    if got := u.TableName(); got != "users" {
        t.Errorf("expected 'users', got %s", got)
    }
}
```

**Criterio de éxito:** Tests pasan

---

### Fase 7: Release

#### Tarea 7.1: Crear release de entities

```bash
# Crear tag para postgres
cd postgres
git tag postgres/entities/v0.1.0
git push origin postgres/entities/v0.1.0

# Crear tag para mongodb
cd ../mongodb
git tag mongodb/entities/v0.1.0
git push origin mongodb/entities/v0.1.0
```

**Criterio de éxito:** Tags creados, disponibles en GitHub

---

## 📊 Estimación de Esfuerzo

| Fase | Tareas | Tiempo Estimado |
|------|--------|-----------------|
| Fase 1: Setup | 1 tarea | 5 min |
| Fase 2: PostgreSQL | 14 entities | 3-4 horas |
| Fase 3: MongoDB | 3 entities | 1-2 horas |
| Fase 4: go.mod | 1 tarea | 10 min |
| Fase 5: README | 2 archivos | 30 min |
| Fase 6: Testing | Tests básicos | 1 hora |
| Fase 7: Release | Tags | 10 min |
| **TOTAL** | | **6-8 horas** |

---

## 🔗 Dependencias

**Antes de este sprint:**
- ✅ Migraciones PostgreSQL existen
- ✅ Migraciones MongoDB existen

**Después de este sprint:**
- ➡️ api-mobile puede adaptar sus entities
- ➡️ api-administracion puede adaptar sus entities
- ➡️ worker puede adaptar sus entities

---

## ⚠️ Notas Importantes

1. **Entities son SOLO estructura**, sin lógica:
   - ✅ Struct con campos
   - ✅ Tags de mapeo (db, bson)
   - ✅ TableName/CollectionName
   - ❌ NO validaciones de negocio
   - ❌ NO constructores complejos
   - ❌ NO métodos de lógica

2. **Sincronización con migraciones:**
   - Cada entity debe reflejar EXACTAMENTE la migración
   - Si cambias migración → cambias entity
   - Mantener versiones sincronizadas

3. **Imports en proyectos:**
   ```go
   // En api-mobile, api-administracion, worker:
   import (
       pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
       mongoentities "github.com/EduGoGroup/edugo-infrastructure/mongodb/entities"
   )
   ```

---

## 📈 Criterios de Éxito del Sprint

- [ ] 14 entities PostgreSQL creados
- [ ] 3 entities MongoDB creados
- [ ] Todos compilan sin errores
- [ ] Tests básicos pasan
- [ ] READMEs documentan uso
- [ ] Tags de release creados
- [ ] Disponibles en GitHub para `go get`

---

**Siguiente paso:** Una vez completado este sprint, los otros proyectos pueden ejecutar sus sprints de adaptación en paralelo.

---

**Generado por:** Claude Code  
**Fecha:** 20 de Noviembre, 2025
