# 📦 Módulos - EduGo Infrastructure

Descripción detallada de cada módulo del repositorio y cómo consumirlos.

---

## 🗂️ Estructura de Módulos

```
edugo-infrastructure/
├── postgres/          # Módulo Go: entities + migraciones PostgreSQL
├── mongodb/           # Módulo Go: entities + seeds MongoDB
├── schemas/           # Módulo Go: JSON Schemas de validación
├── messaging/         # Módulo Go: validador de eventos RabbitMQ
├── docker/            # Docker Compose (no es módulo Go)
├── seeds/             # Seeds de datos (no es módulo Go)
├── scripts/           # Scripts de utilidad (no es módulo Go)
└── tools/             # Herramientas internas (no es módulo Go)
```

---

## 🐘 Módulo `postgres`

### Propósito
Entities Go que reflejan las tablas de PostgreSQL y herramientas de migración.

### Estructura

```
postgres/
├── entities/              # Structs Go
│   ├── user.go
│   ├── school.go
│   ├── academic_unit.go
│   ├── membership.go
│   ├── material.go
│   ├── assessment.go
│   ├── assessment_attempt.go
│   └── assessment_attempt_answer.go
├── migrations/            # Archivos SQL
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_schools.up.sql
│   └── ...
├── cmd/
│   ├── migrate/           # CLI de migraciones
│   └── runner/            # Runner de 4 capas
├── go.mod
└── go.sum
```

### Instalación

```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres
```

### Uso de Entities

```go
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"

// Crear usuario
user := &pgentities.User{
    ID:        uuid.New(),
    Email:     "teacher@school.com",
    FirstName: "John",
    LastName:  "Doe",
    Role:      "teacher",
    IsActive:  true,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// Obtener nombre de tabla
tableName := user.TableName() // "users"
```

### Entities Disponibles

| Entity | Tabla | Descripción |
|--------|-------|-------------|
| `User` | `users` | Usuarios del sistema |
| `School` | `schools` | Instituciones educativas |
| `AcademicUnit` | `academic_units` | Unidades académicas jerárquicas |
| `Membership` | `memberships` | Relación usuario-escuela |
| `Material` | `materials` | Materiales educativos |
| `Assessment` | `assessment` | Metadata de quizzes |
| `AssessmentAttempt` | `assessment_attempt` | Intentos de estudiantes |
| `AssessmentAttemptAnswer` | `assessment_attempt_answer` | Respuestas individuales |

### Proyectos que lo Usan

- **api-mobile:** Todas las entities
- **api-administracion:** User, School, AcademicUnit, Membership
- **worker:** Todas las entities

---

## 🍃 Módulo `mongodb`

### Propósito
Entities Go que reflejan las collections de MongoDB.

### Estructura

```
mongodb/
├── entities/
│   ├── material_assessment.go   # Assessment con preguntas
│   ├── material_summary.go      # Resúmenes generados
│   └── material_event.go        # Log de eventos
├── migrations/                  # Scripts de índices
├── seeds/                       # Datos de prueba
├── go.mod
└── go.sum
```

### Instalación

```bash
go get github.com/EduGoGroup/edugo-infrastructure/mongodb
```

### Uso de Entities

```go
import mongoentities "github.com/EduGoGroup/edugo-infrastructure/mongodb/entities"

// Crear assessment
assessment := &mongoentities.MaterialAssessment{
    MaterialID: "uuid-string",
    Questions: []mongoentities.Question{
        {
            QuestionID:   "q1",
            QuestionText: "¿Qué es POO?",
            QuestionType: "multiple_choice",
            Options: []mongoentities.Option{
                {OptionID: "a", OptionText: "Programación orientada a objetos"},
                {OptionID: "b", OptionText: "Otro concepto"},
            },
            CorrectAnswer: "a",
            Points:        10,
            Difficulty:    "easy",
        },
    },
    TotalQuestions: 1,
    TotalPoints:    10,
    AIModel:        "gpt-4",
    CreatedAt:      time.Now(),
    UpdatedAt:      time.Now(),
}

// Obtener nombre de collection
collectionName := assessment.CollectionName() // "material_assessment_worker"
```

### Entities Disponibles

| Entity | Collection | Descripción |
|--------|------------|-------------|
| `MaterialAssessment` | `material_assessment_worker` | Preguntas de quizzes |
| `MaterialSummary` | `material_summary` | Resúmenes de materiales |
| `MaterialEvent` | `material_event` | Log de procesamiento |

### Tipos Embebidos

```go
// Question - Pregunta del assessment
type Question struct {
    QuestionID   string
    QuestionText string
    QuestionType string   // multiple_choice, true_false, open
    Options      []Option
    CorrectAnswer string
    Explanation  string
    Points       int
    Difficulty   string   // easy, medium, hard
    Tags         []string
}

// Option - Opción de respuesta
type Option struct {
    OptionID   string
    OptionText string
}

// TokenUsage - Consumo de tokens de IA
type TokenUsage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}
```

### Proyectos que lo Usan

- **worker:** MaterialAssessment, MaterialSummary, MaterialEvent
- **api-mobile:** MaterialAssessment (read-only)

---

## 📋 Módulo `schemas`

### Propósito
JSON Schemas para validación de eventos y datos.

### Estructura

```
schemas/
├── events/
│   ├── material-uploaded-v1.schema.json
│   ├── assessment-generated-v1.schema.json
│   ├── material-deleted-v1.schema.json
│   └── student-enrolled-v1.schema.json
├── validator.go
├── go.mod
└── go.sum
```

### Instalación

```bash
go get github.com/EduGoGroup/edugo-infrastructure/schemas
```

### Uso

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

validator := schemas.NewValidator()

// Validar JSON contra schema
jsonData := []byte(`{"event_type": "material.uploaded", ...}`)
err := validator.ValidateEvent(jsonData, "material-uploaded-v1")
if err != nil {
    log.Error("Invalid event", err)
}
```

### Schemas Disponibles

| Schema | Versión | Descripción |
|--------|---------|-------------|
| `material-uploaded-v1` | 1.0 | Material subido |
| `assessment-generated-v1` | 1.0 | Assessment generado |
| `material-deleted-v1` | 1.0 | Material eliminado |
| `student-enrolled-v1` | 1.0 | Estudiante matriculado |

---

## 📬 Módulo `messaging`

### Propósito
Validador de eventos RabbitMQ con JSON Schema integrado.

### Estructura

```
messaging/
├── events/
│   ├── material_uploaded.go
│   ├── assessment_generated.go
│   ├── material_deleted.go
│   └── student_enrolled.go
├── validator.go
├── go.mod
└── go.sum
```

### Instalación

```bash
go get github.com/EduGoGroup/edugo-infrastructure/messaging
```

### Uso - Publisher

```go
import "github.com/EduGoGroup/edugo-infrastructure/messaging"

// Crear evento
event := messaging.MaterialUploadedEvent{
    EventID:      uuid.New().String(),
    EventType:    "material.uploaded",
    EventVersion: "1.0",
    Timestamp:    time.Now(),
    Payload: messaging.MaterialUploadedPayload{
        MaterialID:    materialID.String(),
        SchoolID:      schoolID.String(),
        TeacherID:     teacherID.String(),
        FileURL:       s3URL,
        FileSizeBytes: fileSize,
        FileType:      "application/pdf",
    },
}

// Validar antes de publicar
validator := messaging.NewEventValidator()
if err := validator.Validate(event); err != nil {
    return fmt.Errorf("invalid event: %w", err)
}

// Publicar
publisher.Publish("edugo.materials", "material.uploaded", event)
```

### Uso - Consumer

```go
import "github.com/EduGoGroup/edugo-infrastructure/messaging"

func handleMessage(msg amqp.Delivery) error {
    validator := messaging.NewEventValidator()
    
    // Validar mensaje recibido
    if err := validator.ValidateJSON(msg.Body, "material.uploaded", "1.0"); err != nil {
        logger.Error("Invalid event", "error", err)
        return sendToDLQ(msg, err)
    }
    
    // Deserializar
    var event messaging.MaterialUploadedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }
    
    // Procesar
    return processEvent(event)
}
```

### Tipos de Eventos

```go
// MaterialUploadedEvent
type MaterialUploadedEvent struct {
    EventID      string
    EventType    string // "material.uploaded"
    EventVersion string // "1.0"
    Timestamp    time.Time
    Payload      MaterialUploadedPayload
}

// AssessmentGeneratedEvent
type AssessmentGeneratedEvent struct {
    EventID      string
    EventType    string // "assessment.generated"
    EventVersion string // "1.0"
    Timestamp    time.Time
    Payload      AssessmentGeneratedPayload
}
```

---

## 🐳 Docker

### Propósito
Configuración de Docker Compose para desarrollo local.

### Estructura

```
docker/
├── docker-compose.yml
└── README.md
```

### Servicios Definidos

| Servicio | Imagen | Profile |
|----------|--------|---------|
| `postgres` | postgres:15-alpine | core |
| `mongodb` | mongo:7.0 | core |
| `rabbitmq` | rabbitmq:3.12-management-alpine | messaging |
| `redis` | redis:7-alpine | cache |
| `pgadmin` | dpage/pgadmin4:latest | tools |
| `mongo-express` | mongo-express:latest | tools |

### Profiles

```bash
# Core (default)
docker-compose up -d postgres mongodb

# Con messaging
docker-compose --profile messaging up -d

# Con cache
docker-compose --profile cache up -d

# Con tools
docker-compose --profile tools up -d

# Todo
docker-compose --profile messaging --profile cache --profile tools up -d
```

---

## 🌱 Seeds

### Propósito
Datos de prueba para desarrollo y testing.

### Estructura

```
seeds/
├── postgres/
│   ├── users.sql
│   ├── schools.sql
│   └── memberships.sql
└── mongodb/
    ├── material_assessment_worker.js
    ├── material_summary.js
    └── material_event.js
```

### Cargar Seeds

```bash
# Todos los seeds
make seed

# Solo PostgreSQL mínimo
make seed-minimal
```

---

## 🔧 Scripts

### Propósito
Scripts de utilidad para desarrollo.

### Estructura

```
scripts/
├── dev-setup.sh         # Setup inicial completo
├── seed-data.sh         # Cargar seeds
├── validate-env.sh      # Validar .env
└── ...
```

### Uso

```bash
# Setup completo primera vez
./scripts/dev-setup.sh

# Validar variables de entorno
./scripts/validate-env.sh
```

---

## 🛠️ Tools

### Propósito
Herramientas internas del proyecto.

### Estructura

```
tools/
├── generate-entities/   # Generador de entities
├── schema-validator/    # Validador de schemas
└── ...
```

---

## 📊 Dependencias entre Módulos

```
┌─────────────┐
│  messaging  │──────────────────┐
└──────┬──────┘                  │
       │ imports                 │ imports
       ▼                         ▼
┌─────────────┐           ┌─────────────┐
│   schemas   │           │  postgres/  │
│             │           │  entities   │
└─────────────┘           └─────────────┘
                                 │
                                 │ reference (ID)
                                 ▼
                          ┌─────────────┐
                          │  mongodb/   │
                          │  entities   │
                          └─────────────┘
```

---

## 🔄 Versionado de Módulos

Cada módulo Go tiene su propio tag de versión:

```bash
# PostgreSQL entities
git tag postgres/v0.1.0
git push origin postgres/v0.1.0

# MongoDB entities
git tag mongodb/v0.1.0
git push origin mongodb/v0.1.0

# Messaging
git tag messaging/v0.1.0
git push origin messaging/v0.1.0

# Schemas
git tag schemas/v0.1.0
git push origin schemas/v0.1.0
```

### Consumir Versión Específica

```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.1.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.1.0
```

---

## 📝 Checklist de Integración

### Para integrar `postgres/entities`:

- [ ] Agregar dependencia: `go get .../postgres`
- [ ] Importar: `import pgentities ".../postgres/entities"`
- [ ] Configurar conexión DB con mismas credenciales
- [ ] Usar `entity.TableName()` para queries

### Para integrar `mongodb/entities`:

- [ ] Agregar dependencia: `go get .../mongodb`
- [ ] Importar: `import mongoentities ".../mongodb/entities"`
- [ ] Configurar conexión MongoDB
- [ ] Usar `entity.CollectionName()` para operaciones

### Para integrar `messaging`:

- [ ] Agregar dependencia: `go get .../messaging`
- [ ] Configurar RabbitMQ
- [ ] Usar `NewEventValidator()` para validar
- [ ] Implementar publisher/consumer según necesidad

---

**Última actualización:** Diciembre 2024
