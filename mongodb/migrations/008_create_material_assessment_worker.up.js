// 008_create_material_assessment_worker.up.js  
// Migración: Crear collection material_assessment_worker para quizzes generados por IA
// Proyecto: edugo-infrastructure
// Consumidor: edugo-worker
// Fecha: 2025-11-18
//
// NOTA: Se llama "material_assessment_worker" (no "material_assessment") para evitar
// conflicto con la collection existente de api-admin y permitir evolución independiente.

print("🔧 Creating collection: material_assessment_worker...");

db.createCollection("material_assessment_worker", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: [
        "material_id",
        "questions",
        "total_questions",
        "total_points",
        "version",
        "ai_model",
        "processing_time_ms",
        "created_at",
        "updated_at"
      ],
      properties: {
        material_id: {
          bsonType: "string",
          pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
          description: "UUID v4 del material en PostgreSQL (requerido)"
        },
        questions: {
          bsonType: "array",
          minItems: 3,
          maxItems: 20,
          items: {
            bsonType: "object",
            required: ["question_id", "question_text", "question_type", "correct_answer", "points", "difficulty"],
            properties: {
              question_id: {
                bsonType: "string",
                pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
                description: "UUID v4 de la pregunta"
              },
              question_text: {
                bsonType: "string",
                minLength: 10,
                maxLength: 1000,
                description: "Texto de la pregunta (10-1000 caracteres)"
              },
              question_type: {
                enum: ["multiple_choice", "true_false", "open"],
                description: "Tipo de pregunta"
              },
              options: {
                bsonType: "array",
                minItems: 2,
                maxItems: 5,
                items: {
                  bsonType: "object",
                  required: ["option_id", "option_text"],
                  properties: {
                    option_id: { bsonType: "string" },
                    option_text: { bsonType: "string", minLength: 1, maxLength: 500 }
                  }
                },
                description: "Opciones de respuesta (para multiple_choice y true_false)"
              },
              correct_answer: {
                bsonType: "string",
                minLength: 1,
                description: "Respuesta correcta (option_id para multiple_choice, texto para open)"
              },
              explanation: {
                bsonType: "string",
                maxLength: 2000,
                description: "Explicación de la respuesta correcta (opcional)"
              },
              points: {
                bsonType: "int",
                minimum: 1,
                maximum: 100,
                description: "Puntos de la pregunta (1-100)"
              },
              difficulty: {
                enum: ["easy", "medium", "hard"],
                description: "Nivel de dificultad"
              },
              tags: {
                bsonType: "array",
                items: { bsonType: "string" },
                description: "Tags/categorías de la pregunta (opcional)"
              }
            }
          },
          description: "Array de preguntas (3-20 elementos)"
        },
        total_questions: {
          bsonType: "int",
          minimum: 3,
          maximum: 20,
          description: "Número total de preguntas"
        },
        total_points: {
          bsonType: "int",
          minimum: 1,
          description: "Suma total de puntos de todas las preguntas"
        },
        version: {
          bsonType: "int",
          minimum: 1,
          description: "Versión del assessment (>= 1, incrementa en reprocesos)"
        },
        ai_model: {
          enum: ["gpt-4", "gpt-3.5-turbo", "gpt-4-turbo", "gpt-4o"],
          description: "Modelo de IA utilizado"
        },
        processing_time_ms: {
          bsonType: "int",
          minimum: 0,
          description: "Tiempo de procesamiento en milisegundos"
        },
        token_usage: {
          bsonType: "object",
          properties: {
            prompt_tokens: { bsonType: "int", minimum: 0 },
            completion_tokens: { bsonType: "int", minimum: 0 },
            total_tokens: { bsonType: "int", minimum: 0 }
          },
          description: "Metadata de tokens consumidos (opcional)"
        },
        metadata: {
          bsonType: "object",
          description: "Metadata adicional (opcional)"
        },
        created_at: {
          bsonType: "date",
          description: "Fecha de creación (requerido)"
        },
        updated_at: {
          bsonType: "date",
          description: "Fecha de última actualización (requerido)"
        }
      },
      additionalProperties: true
    }
  },
  validationLevel: "strict",
  validationAction: "error"
});

// Crear índices
print("📑 Creating indexes for material_assessment_worker...");

db.material_assessment_worker.createIndex(
  { material_id: 1 },
  { unique: true, name: "idx_material_id" }
);

db.material_assessment_worker.createIndex(
  { created_at: -1 },
  { name: "idx_created_at" }
);

db.material_assessment_worker.createIndex(
  { version: 1 },
  { name: "idx_version" }
);

db.material_assessment_worker.createIndex(
  { "questions.difficulty": 1 },
  { name: "idx_questions_difficulty" }
);

db.material_assessment_worker.createIndex(
  { total_questions: 1, created_at: -1 },
  { name: "idx_total_questions_created" }
);

print("✅ Collection material_assessment_worker created successfully with 5 indexes");
