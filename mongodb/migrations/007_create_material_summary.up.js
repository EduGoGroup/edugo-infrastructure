// 007_create_material_summary.up.js
// Migración: Crear collection material_summary para resúmenes generados por IA
// Proyecto: edugo-infrastructure
// Consumidor: edugo-worker
// Fecha: 2025-11-18

print("🔧 Creating collection: material_summary...");

db.createCollection("material_summary", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: [
        "material_id",
        "summary",
        "key_points",
        "language",
        "word_count",
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
        summary: {
          bsonType: "string",
          minLength: 10,
          maxLength: 5000,
          description: "Resumen generado por IA (min 10, max 5000 caracteres)"
        },
        key_points: {
          bsonType: "array",
          minItems: 1,
          maxItems: 10,
          items: {
            bsonType: "string",
            minLength: 5,
            maxLength: 500
          },
          description: "Array de puntos clave (1-10 elementos)"
        },
        language: {
          enum: ["es", "en", "pt"],
          description: "Idioma del resumen: español, inglés o portugués"
        },
        word_count: {
          bsonType: "int",
          minimum: 1,
          description: "Número de palabras del resumen (mínimo 1)"
        },
        version: {
          bsonType: "int",
          minimum: 1,
          description: "Versión del resumen (>= 1, incrementa en reprocesos)"
        },
        ai_model: {
          enum: ["gpt-4", "gpt-3.5-turbo", "gpt-4-turbo", "gpt-4o"],
          description: "Modelo de IA utilizado para generar el resumen"
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
          description: "Metadata adicional del procesamiento (opcional)"
        },
        created_at: {
          bsonType: "date",
          description: "Fecha de creación del resumen (requerido)"
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
print("📑 Creating indexes for material_summary...");

db.material_summary.createIndex(
  { material_id: 1 },
  { unique: true, name: "idx_material_id" }
);

db.material_summary.createIndex(
  { created_at: -1 },
  { name: "idx_created_at" }
);

db.material_summary.createIndex(
  { version: 1 },
  { name: "idx_version" }
);

db.material_summary.createIndex(
  { language: 1, created_at: -1 },
  { name: "idx_language_created" }
);

print("✅ Collection material_summary created successfully with 4 indexes");
