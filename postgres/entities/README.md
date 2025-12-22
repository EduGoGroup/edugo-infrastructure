# PostgreSQL Entities

Entities base que reflejan el schema de PostgreSQL para el ecosistema EduGo.

---

## 📋 Entities Disponibles (13 de 14 planificadas)

| # | Entity | Tabla | Migración | Status |
|---|--------|-------|-----------|--------|
| 1 | `User` | `users` | `001_create_users.up.sql` | ✅ Disponible |
| 2 | `School` | `schools` | `002_create_schools.up.sql` | ✅ Disponible |
| 3 | `AcademicUnit` | `academic_units` | `003_create_academic_units.up.sql` | ✅ Disponible |
| 4 | `Membership` | `memberships` | `004_create_memberships.up.sql` | ✅ Disponible |
| 5 | `Material` | `materials` | `005_create_materials.up.sql` | ✅ Disponible |
| 6 | `Assessment` | `assessment` | `006_create_assessments.up.sql` | ✅ Disponible |
| 7 | `AssessmentAttempt` | `assessment_attempt` | `007_create_assessment_attempts.up.sql` | ✅ Disponible |
| 8 | `AssessmentAttemptAnswer` | `assessment_attempt_answer` | `008_create_assessment_answers.up.sql` | ✅ Disponible |
| 9 | `MaterialVersion` | `material_versions` | `012_create_material_versions.up.sql` | ✅ Disponible |
| 10 | `Subject` | `subjects` | `013_create_subjects.up.sql` | ✅ Disponible |
| 11 | `Unit` | `units` | `014_create_units.up.sql` | ✅ Disponible |
| 12 | `GuardianRelation` | `guardian_relations` | `015_create_guardian_relations.up.sql` | ✅ Disponible |
| 13 | `Progress` | `progress` | `016_create_progress.up.sql` | ✅ Disponible |

---

## 🚫 Entities Pendientes (1 bloqueada)

| # | Entity | Tabla Esperada | Razón |
|---|--------|----------------|-------|
| 1 | `AssessmentQuestion` | `assessment_questions` | Migración aún no creada |

---

## 📖 Uso Básico

### Importar Entities

```go
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
```

### Ejemplo: User Entity

```go
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

// Obtener nombre de tabla
tableName := user.TableName() // "users"
```

### Ejemplo: School Entity

```go
school := &pgentities.School{
    ID:               uuid.New(),
    Name:             "Colegio Ejemplo",
    Code:             "COL001",
    Country:          "Chile",
    IsActive:         true,
    SubscriptionTier: "basic",
    MaxTeachers:      10,
    MaxStudents:      100,
    CreatedAt:        time.Now(),
    UpdatedAt:        time.Now(),
}
```

### Ejemplo: Assessment Entities

```go
// Assessment (metadata en PostgreSQL)
assessment := &pgentities.Assessment{
    ID:              uuid.New(),
    MaterialID:      materialID,
    MongoDocumentID: objectID.Hex(), // Ref a MongoDB
    QuestionsCount:  10,
    Status:          "published",
    CreatedAt:       time.Now(),
    UpdatedAt:       time.Now(),
}

// AssessmentAttempt (intento de estudiante)
attempt := &pgentities.AssessmentAttempt{
    ID:           uuid.New(),
    AssessmentID: assessment.ID,
    StudentID:    studentID,
    StartedAt:    time.Now(),
    Status:       "in_progress",
    CreatedAt:    time.Now(),
    UpdatedAt:    time.Now(),
}

// AssessmentAttemptAnswer (respuesta individual)
answer := &pgentities.AssessmentAttemptAnswer{
    ID:            uuid.New(),
    AttemptID:     attempt.ID,
    QuestionIndex: 0,
    StudentAnswer: "opt3",
    AnsweredAt:    time.Now(),
    CreatedAt:     time.Now(),
    UpdatedAt:     time.Now(),
}
```

### Nuevas Entities Disponibles

#### MaterialVersion

```go
version := &pgentities.MaterialVersion{
    ID:            uuid.New(),
    MaterialID:    material.ID,
    VersionNumber: 2,
    FileURL:       "s3://bucket/file-v2.pdf",
    FileSize:      1024000,
    UploadedBy:    userID,
    Metadata:      []byte(`{"comment": "Fixed typos"}`),
    CreatedAt:     time.Now(),
}
```

#### Subject

```go
subject := &pgentities.Subject{
    ID:          uuid.New(),
    SchoolID:    school.ID,
    Name:        "Matemáticas",
    Code:        "MAT101",
    Description: "Matemáticas nivel básico",
    IsActive:    true,
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}
```

#### Unit

```go
unit := &pgentities.Unit{
    ID:          uuid.New(),
    SubjectID:   subject.ID,
    Name:        "Álgebra",
    OrderIndex:  1,
    Description: "Fundamentos de álgebra",
    IsActive:    true,
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}
```

#### GuardianRelation

```go
relation := &pgentities.GuardianRelation{
    ID:           uuid.New(),
    GuardianID:   guardianUser.ID,
    StudentID:    studentUser.ID,
    RelationType: "parent",
    IsActive:     true,
    CreatedAt:    time.Now(),
    UpdatedAt:    time.Now(),
}
```

#### Progress

```go
progress := &pgentities.Progress{
    ID:               uuid.New(),
    StudentID:        student.ID,
    MaterialID:       material.ID,
    Status:           "in_progress",
    CompletionPercent: 45,
    LastAccessedAt:   time.Now(),
    CreatedAt:        time.Now(),
    UpdatedAt:        time.Now(),
}
```

---

## 🔧 Uso Avanzado

### Con sqlx

```go
import (
    "github.com/jmoiron/sqlx"
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func GetUser(db *sqlx.DB, id uuid.UUID) (*pgentities.User, error) {
    var user pgentities.User
    query := "SELECT * FROM users WHERE id = $1"
    err := db.Get(&user, query, id)
    return &user, err
}
```

### Con database/sql

```go
func ListSchools(db *sql.DB) ([]pgentities.School, error) {
    query := "SELECT * FROM schools WHERE is_active = true"
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schools []pgentities.School
    for rows.Next() {
        var school pgentities.School
        err := rows.Scan(
            &school.ID, &school.Name, &school.Code,
            // ... otros campos
        )
        if err != nil {
            return nil, err
        }
        schools = append(schools, school)
    }
    return schools, nil
}
```

---

## 📐 Reglas y Principios

### ✅ Las Entities SON:

- **Reflejos exactos** de tablas SQL
- **Estructuras de datos** sin lógica
- **Mapeo** de columnas con tags `db:`
- **Documentadas** con referencias a migraciones
- **Versionadas** con el módulo `postgres`

### ❌ Las Entities NO SON:

- **Lógica de negocio** (usar domain services)
- **Validaciones** (usar validators en APIs)
- **Constructores complejos** (solo structs)
- **Métodos de mutación** (solo getters simples)
- **DTOs** (usar models separados para APIs)

---

## 🎯 Campos Comunes

### JSONB Fields

Los campos JSONB se mapean como `[]byte`:

```go
type School struct {
    Metadata []byte `db:"metadata"` // JSONB
}

// Para usar:
import "encoding/json"

// Serializar
metadata := map[string]interface{}{"logo": "url"}
school.Metadata, _ = json.Marshal(metadata)

// Deserializar
var meta map[string]interface{}
json.Unmarshal(school.Metadata, &meta)
```

### Nullable Fields

Los campos nullable usan punteros:

```go
type User struct {
    DeletedAt *time.Time `db:"deleted_at"` // NULL permitido
}

// Soft delete
now := time.Now()
user.DeletedAt = &now
```

### UUID Fields

```go
import "github.com/google/uuid"

user := &User{
    ID: uuid.New(), // Generar UUID
}
```

---

## 🔗 Referencias entre Entities

### Relaciones

```go
// Material → School (FK)
material := &Material{
    SchoolID: school.ID,
}

// Assessment → Material (FK)
assessment := &Assessment{
    MaterialID: material.ID,
}

// AssessmentAttempt → Assessment + User (FKs)
attempt := &AssessmentAttempt{
    AssessmentID: assessment.ID,
    StudentID:    user.ID,
}

// Subject → School (FK)
subject := &Subject{
    SchoolID: school.ID,
}

// Unit → Subject (FK)
unit := &Unit{
    SubjectID: subject.ID,
}

// Progress → Student + Material (FKs)
progress := &Progress{
    StudentID:  student.ID,
    MaterialID: material.ID,
}

// GuardianRelation → Guardian + Student (FKs)
relation := &GuardianRelation{
    GuardianID: guardian.ID,
    StudentID:  student.ID,
}
```

**Nota:** Las entities **NO incluyen** joins automáticos. Hacer queries con joins en tu aplicación.

---

## 🧪 Testing

Ver ejemplos de tests en `*_test.go` (pendiente Fase 2).

---

## 📦 Versionado

Las entities se versionan con el módulo `postgres`:

```bash
# Release de entities
cd postgres
git tag postgres/entities/v0.1.0
git push origin postgres/entities/v0.1.0
```

---

## 🚀 Proyectos que Pueden Usar Estas Entities

| Proyecto | Entities Disponibles | Status |
|----------|---------------------|--------|
| **api-mobile** | Todas excepto AssessmentQuestion | ✅ Listo para migración |
| **api-administracion** | Todas excepto AssessmentQuestion | ✅ Listo para migración |
| **worker** | Todas excepto AssessmentQuestion | ✅ Listo para migración |

---

## 📝 Próximos Pasos

1. **Fase 2:** Crear migración para `AssessmentQuestion`
2. **Fase 2:** Validar compilación con Go 1.25
3. **Fase 2:** Ejecutar tests de integración
4. **Fase 3:** Release de `postgres/entities/v1.0.0`
5. **Proyectos:** Migrar api-mobile, api-administracion, worker

---

**Última actualización:** Diciembre 2024  
**Versión:** 2.0
