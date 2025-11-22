# MongoDB Entities

Entities base que reflejan los schemas de MongoDB para el ecosistema EduGo.

---

## 📋 Entities Disponibles (3 de 3 planificadas)

| # | Entity | Collection | Seed | Status |
|---|--------|------------|------|--------|
| 1 | `MaterialAssessment` | `material_assessment_worker` | `material_assessment_worker.js` | ✅ Disponible |
| 2 | `MaterialSummary` | `material_summary` | `material_summary.js` | ✅ Disponible |
| 3 | `MaterialEvent` | `material_event` | `material_event.js` | ✅ Disponible |

---

## 📖 Uso Básico

### Importar Entities

```go
import mongoentities "github.com/EduGoGroup/edugo-infrastructure/mongodb/entities"
```

### Ejemplo: MaterialAssessment Entity

```go
import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    mongoentities "github.com/EduGoGroup/edugo-infrastructure/mongodb/entities"
)

// Crear assessment con preguntas
assessment := &mongoentities.MaterialAssessment{
    MaterialID:     "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    Questions: []mongoentities.Question{
        {
            QuestionID:   "q1111111-1111-1111-1111-111111111111",
            QuestionText: "¿Qué es encapsulación?",
            QuestionType: "multiple_choice",
            Options: []mongoentities.Option{
                {OptionID: "opt1", OptionText: "Herencia"},
                {OptionID: "opt2", OptionText: "Polimorfismo"},
                {OptionID: "opt3", OptionText: "Encapsulación"},
            },
            CorrectAnswer: "opt3",
            Explanation:   "La encapsulación...",
            Points:        10,
            Difficulty:    "medium",
            Tags:          []string{"POO", "conceptos"},
        },
    },
    TotalQuestions:   1,
    TotalPoints:      10,
    Version:          1,
    AIModel:          "gpt-4",
    ProcessingTimeMs: 5200,
    TokenUsage: &mongoentities.TokenUsage{
        PromptTokens:     1200,
        CompletionTokens: 450,
        TotalTokens:      1650,
    },
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// Obtener nombre de collection
collectionName := assessment.CollectionName() // "material_assessment_worker"
```

### Ejemplo: MaterialSummary Entity

```go
summary := &mongoentities.MaterialSummary{
    MaterialID: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    Summary:    "Este material cubre los fundamentos de POO...",
    KeyPoints: []string{
        "Introducción a POO",
        "Clases y objetos",
        "Herencia y polimorfismo",
    },
    Language:         "es",
    WordCount:        42,
    Version:          1,
    AIModel:          "gpt-4",
    ProcessingTimeMs: 3500,
    TokenUsage: &mongoentities.TokenUsage{
        PromptTokens:     850,
        CompletionTokens: 180,
        TotalTokens:      1030,
    },
    Metadata: &mongoentities.SummaryMetadata{
        SourceLength: 5420,
        HasImages:    false,
    },
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
```

### Ejemplo: MaterialEvent Entity

```go
event := &mongoentities.MaterialEvent{
    EventType:  "material_uploaded",
    MaterialID: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    UserID:     "u1111111-1111-1111-1111-111111111111",
    Payload: primitive.M{
        "filename":  "java-poo.pdf",
        "file_size": 1024000,
        "mime_type": "application/pdf",
    },
    Status:     "completed",
    RetryCount: 0,
    CreatedAt:  time.Now(),
    UpdatedAt:  time.Now(),
}
```

---

## 🔧 Uso con MongoDB Driver

### Insertar Documento

```go
import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    mongoentities "github.com/EduGoGroup/edugo-infrastructure/mongodb/entities"
)

func InsertAssessment(ctx context.Context, db *mongo.Database, assessment *mongoentities.MaterialAssessment) error {
    collection := db.Collection(assessment.CollectionName())
    result, err := collection.InsertOne(ctx, assessment)
    if err != nil {
        return err
    }
    // Asignar ID generado
    assessment.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}
```

### Buscar Documento

```go
func FindAssessmentByMaterialID(ctx context.Context, db *mongo.Database, materialID string) (*mongoentities.MaterialAssessment, error) {
    var assessment mongoentities.MaterialAssessment
    collection := db.Collection(assessment.CollectionName())

    filter := primitive.M{"material_id": materialID}
    err := collection.FindOne(ctx, filter).Decode(&assessment)
    if err != nil {
        return nil, err
    }
    return &assessment, nil
}
```

### Actualizar Documento

```go
func UpdateSummary(ctx context.Context, db *mongo.Database, summary *mongoentities.MaterialSummary) error {
    collection := db.Collection(summary.CollectionName())

    filter := primitive.M{"_id": summary.ID}
    update := primitive.M{
        "$set": primitive.M{
            "summary":     summary.Summary,
            "key_points":  summary.KeyPoints,
            "updated_at":  time.Now(),
        },
    }

    _, err := collection.UpdateOne(ctx, filter, update)
    return err
}
```

### Listar Documentos

```go
func ListEventsByMaterial(ctx context.Context, db *mongo.Database, materialID string) ([]mongoentities.MaterialEvent, error) {
    var event mongoentities.MaterialEvent
    collection := db.Collection(event.CollectionName())

    filter := primitive.M{"material_id": materialID}
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var events []mongoentities.MaterialEvent
    if err := cursor.All(ctx, &events); err != nil {
        return nil, err
    }
    return events, nil
}
```

---

## 📐 Reglas y Principios

### ✅ Las Entities SON:

- **Reflejos exactos** de documents MongoDB
- **Estructuras de datos** sin lógica
- **Mapeo** de fields con tags `bson:`
- **Documentadas** con referencias a seeds
- **Versionadas** con el módulo `mongodb`

### ❌ Las Entities NO SON:

- **Lógica de negocio** (usar domain services)
- **Validaciones** (usar validators en worker)
- **Constructores complejos** (solo structs)
- **Métodos de mutación** (solo getters simples)
- **DTOs** (usar models separados para worker)

---

## 🎯 Campos Especiales

### ObjectID (\_id)

```go
import "go.mongodb.org/mongo-driver/bson/primitive"

// Crear con ID autogenerado
assessment := &MaterialAssessment{
    // ID se omite, MongoDB lo genera
    MaterialID: "...",
}

// Después de InsertOne:
// assessment.ID contiene el ObjectID generado
```

### primitive.M (Payload flexible)

```go
// MaterialEvent usa primitive.M para payload flexible
event := &MaterialEvent{
    Payload: primitive.M{
        "key1": "value1",
        "key2": 123,
        "nested": primitive.M{
            "subkey": "subvalue",
        },
    },
}
```

### Embedded Documents

```go
// MaterialAssessment tiene arrays de structs embebidos
assessment := &MaterialAssessment{
    Questions: []Question{
        {
            QuestionID:   "q1",
            QuestionText: "...",
            Options: []Option{
                {OptionID: "opt1", OptionText: "..."},
                {OptionID: "opt2", OptionText: "..."},
            },
        },
    },
}
```

---

## 🔗 Relación con PostgreSQL

### MaterialAssessment → PostgreSQL Assessment

```go
// 1. Crear assessment en MongoDB (contenido completo)
mongoAssessment := &mongoentities.MaterialAssessment{
    MaterialID:     materialID,
    Questions:      questions,
    TotalQuestions: len(questions),
    // ...
}
insertedID, _ := insertIntoMongo(mongoAssessment)

// 2. Crear metadata en PostgreSQL (referencia)
pgAssessment := &pgentities.Assessment{
    MaterialID:      materialID,
    MongoDocumentID: insertedID.Hex(), // Referencia a MongoDB
    QuestionsCount:  len(questions),
    Status:          "published",
}
insertIntoPostgres(pgAssessment)
```

### MaterialSummary → PostgreSQL Material

```go
// Material en PostgreSQL ya existe
material := &pgentities.Material{
    ID:     materialID,
    Status: "processing",
}

// Summary generado en MongoDB
summary := &mongoentities.MaterialSummary{
    MaterialID: material.ID.String(),
    Summary:    "...",
    KeyPoints:  []string{"..."},
}

// Actualizar status en PostgreSQL
material.Status = "ready"
updateMaterial(material)
```

---

## 🧪 Testing

Ver ejemplos de tests en `*_test.go` (pendiente Fase 2).

---

## 📦 Versionado

Las entities se versionan con el módulo `mongodb`:

```bash
# Release de entities
cd mongodb
git tag mongodb/entities/v0.1.0
git push origin mongodb/entities/v0.1.0
```

---

## 🚀 Proyectos que Pueden Usar Estas Entities

| Proyecto | Entities Disponibles | Status |
|----------|---------------------|--------|
| **worker** | MaterialAssessment, MaterialSummary, MaterialEvent | ✅ Listo para migración |
| **api-mobile** | MaterialAssessment (read-only para mostrar assessments) | ✅ Listo para migración |

---

## 📋 Tipos de Preguntas Soportadas

### MaterialAssessment.Questions

**Tipos de pregunta (`question_type`):**

| Tipo | Descripción | Options | Ejemplo CorrectAnswer |
|------|-------------|---------|----------------------|
| `multiple_choice` | Opción múltiple | Array de Options | `"opt3"` |
| `true_false` | Verdadero/Falso | Array con 2 Options | `"true"` o `"false"` |
| `open` | Respuesta abierta | Array vacío | Texto esperado |

**Niveles de dificultad (`difficulty`):**
- `easy`
- `medium`
- `hard`

---

## 📝 Próximos Pasos

1. **Fase 2:** Validar compilación con Go 1.25
2. **Fase 2:** Ejecutar tests de integración con MongoDB
3. **Fase 3:** Release de `mongodb/entities/v0.1.0`
4. **Proyectos:** Migrar worker a usar estas entities

---

**Generado por:** Claude Code - Sprint Entities Fase 1
**Fecha:** 2025-11-22
**Versión:** 1.0
