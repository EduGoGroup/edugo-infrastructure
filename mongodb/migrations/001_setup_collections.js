// ============================================================
// MIGRACIÓN 001: Setup de colecciones MongoDB
// Fecha: 2026-02-22
// Reemplaza: Todos los scripts anteriores en migrations/structure/
// Colecciones activas: material_summary, material_assessment_worker, material_event
// Colecciones eliminadas: analytics_events, assessment_attempt_result, audit_logs,
//                         material_content, notifications (huérfanas)
// ============================================================

const db = db.getSiblingDB('edugo');

// ============================================================
// ---- material_summary ----
// Usada por: worker
// Propósito: Resúmenes de materiales generados por IA
// Campos reales: material_id, summary, key_points, language, word_count,
//                version, ai_model, processing_time_ms, token_usage,
//                metadata(source_length, has_images), created_at, updated_at
// ============================================================

print("Creando colección: material_summary");

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
          pattern: "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
          description: "UUID v4 del material en PostgreSQL (requerido)"
        },
        summary: {
          bsonType: "string",
          minLength: 10,
          maxLength: 5000,
          description: "Texto del resumen generado por IA (requerido)"
        },
        key_points: {
          bsonType: "array",
          minItems: 1,
          maxItems: 10,
          items: {
            bsonType: "string"
          },
          description: "Puntos clave extraídos del material (requerido)"
        },
        language: {
          bsonType: "string",
          enum: ["es", "en", "pt"],
          description: "Idioma del resumen: es, en, pt (requerido)"
        },
        word_count: {
          bsonType: "int",
          minimum: 1,
          description: "Conteo de palabras del resumen (requerido)"
        },
        version: {
          bsonType: "int",
          minimum: 1,
          description: "Versión del resumen (requerido)"
        },
        ai_model: {
          bsonType: "string",
          description: "Modelo IA usado para generar el resumen (requerido)"
        },
        processing_time_ms: {
          bsonType: "int",
          minimum: 0,
          description: "Tiempo de procesamiento en milisegundos (requerido)"
        },
        token_usage: {
          bsonType: "object",
          description: "Metadata de tokens consumidos por IA (opcional)",
          properties: {
            prompt_tokens: {
              bsonType: "int"
            },
            completion_tokens: {
              bsonType: "int"
            },
            total_tokens: {
              bsonType: "int"
            }
          }
        },
        metadata: {
          bsonType: "object",
          description: "Metadata adicional del resumen (opcional)",
          properties: {
            source_length: {
              bsonType: "int",
              description: "Longitud del material fuente"
            },
            has_images: {
              bsonType: "bool",
              description: "Si el material tiene imágenes"
            }
          }
        },
        created_at: {
          bsonType: "date",
          description: "Fecha de creación (requerido)"
        },
        updated_at: {
          bsonType: "date",
          description: "Fecha de última actualización (requerido)"
        }
      }
    }
  },
  validationLevel: "moderate",
  validationAction: "error"
});

db.material_summary.createIndex({ material_id: 1 }, { unique: true, name: "idx_material_summary_material_id_unique" });
db.material_summary.createIndex({ language: 1 }, { name: "idx_material_summary_language" });
db.material_summary.createIndex({ version: 1 }, { name: "idx_material_summary_version" });
db.material_summary.createIndex({ created_at: -1 }, { name: "idx_material_summary_created_at" });

print("  OK: material_summary creada con 4 índices");

// ============================================================
// ---- material_assessment_worker ----
// Usada por: worker, api-mobile
// Propósito: Preguntas completas de assessments generados por IA
// Campos reales: material_id, questions[], total_questions, total_points,
//                version, ai_model, processing_time_ms, token_usage,
//                metadata(average_difficulty, estimated_time_min, source_length, has_images),
//                created_at, updated_at
// ============================================================

print("Creando colección: material_assessment_worker");

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
          pattern: "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
          description: "UUID v4 del material en PostgreSQL (requerido)"
        },
        questions: {
          bsonType: "array",
          minItems: 1,
          maxItems: 20,
          items: {
            bsonType: "object",
            required: ["question_id", "question_text", "question_type", "correct_answer", "points", "difficulty"],
            properties: {
              question_id: {
                bsonType: "string"
              },
              question_text: {
                bsonType: "string"
              },
              question_type: {
                bsonType: "string",
                enum: ["multiple_choice", "true_false", "open"]
              },
              options: {
                bsonType: "array",
                items: {
                  bsonType: "object",
                  required: ["option_id", "option_text"],
                  properties: {
                    option_id: { bsonType: "string" },
                    option_text: { bsonType: "string" }
                  }
                }
              },
              correct_answer: {
                bsonType: "string"
              },
              explanation: {
                bsonType: "string"
              },
              points: {
                bsonType: "int",
                minimum: 1
              },
              difficulty: {
                bsonType: "string",
                enum: ["easy", "medium", "hard"]
              },
              tags: {
                bsonType: "array",
                items: { bsonType: "string" }
              }
            }
          },
          description: "Array de preguntas del assessment (requerido)"
        },
        total_questions: {
          bsonType: "int",
          minimum: 1,
          maximum: 20,
          description: "Total de preguntas en el assessment (requerido)"
        },
        total_points: {
          bsonType: "int",
          minimum: 1,
          description: "Puntos totales posibles (requerido)"
        },
        version: {
          bsonType: "int",
          minimum: 1,
          description: "Versión del assessment (requerido)"
        },
        ai_model: {
          bsonType: "string",
          description: "Modelo IA usado para generar el assessment (requerido)"
        },
        processing_time_ms: {
          bsonType: "int",
          minimum: 0,
          description: "Tiempo de procesamiento en milisegundos (requerido)"
        },
        token_usage: {
          bsonType: "object",
          description: "Metadata de tokens consumidos por IA (opcional)",
          properties: {
            prompt_tokens: {
              bsonType: "int"
            },
            completion_tokens: {
              bsonType: "int"
            },
            total_tokens: {
              bsonType: "int"
            }
          }
        },
        metadata: {
          bsonType: "object",
          description: "Metadata adicional del assessment (opcional)",
          properties: {
            average_difficulty: {
              bsonType: "string"
            },
            estimated_time_min: {
              bsonType: "int"
            },
            source_length: {
              bsonType: "int"
            },
            has_images: {
              bsonType: "bool"
            }
          }
        },
        created_at: {
          bsonType: "date",
          description: "Fecha de creación (requerido)"
        },
        updated_at: {
          bsonType: "date",
          description: "Fecha de última actualización (requerido)"
        }
      }
    }
  },
  validationLevel: "moderate",
  validationAction: "error"
});

db.material_assessment_worker.createIndex({ material_id: 1 }, { unique: true, name: "idx_assessment_worker_material_id_unique" });
db.material_assessment_worker.createIndex({ total_questions: 1 }, { name: "idx_assessment_worker_total_questions" });
db.material_assessment_worker.createIndex({ version: 1 }, { name: "idx_assessment_worker_version" });
db.material_assessment_worker.createIndex({ created_at: -1 }, { name: "idx_assessment_worker_created_at" });

print("  OK: material_assessment_worker creada con 4 índices");

// ============================================================
// ---- material_event ----
// Usada por: worker
// Propósito: Cola de eventos de procesamiento de materiales (auditoría)
// Campos reales: event_type, material_id, user_id, payload, status,
//                error_msg, stack_trace, retry_count, processed_at,
//                created_at, updated_at
// ============================================================

print("Creando colección: material_event");

db.createCollection("material_event", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: [
        "event_type",
        "payload",
        "status",
        "retry_count",
        "created_at",
        "updated_at"
      ],
      properties: {
        event_type: {
          bsonType: "string",
          enum: [
            "material_uploaded",
            "material_reprocess",
            "assessment_attempt"
          ],
          description: "Tipo de evento (requerido)"
        },
        material_id: {
          bsonType: "string",
          pattern: "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
          description: "UUID v4 del material en PostgreSQL (opcional)"
        },
        user_id: {
          bsonType: "string",
          pattern: "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
          description: "UUID v4 del usuario en PostgreSQL (opcional)"
        },
        payload: {
          bsonType: "object",
          description: "Datos del evento con estructura flexible (requerido)"
        },
        status: {
          bsonType: "string",
          enum: ["processing", "completed", "failed"],
          description: "Estado actual del procesamiento del evento (requerido)"
        },
        error_msg: {
          bsonType: "string",
          maxLength: 5000,
          description: "Mensaje de error si el procesamiento falló (opcional)"
        },
        stack_trace: {
          bsonType: "string",
          maxLength: 10000,
          description: "Stack trace si el procesamiento falló (opcional)"
        },
        retry_count: {
          bsonType: "int",
          minimum: 0,
          description: "Número de reintentos del evento (requerido)"
        },
        processed_at: {
          bsonType: "date",
          description: "Timestamp de procesamiento exitoso (opcional)"
        },
        created_at: {
          bsonType: "date",
          description: "Fecha de creación (requerido)"
        },
        updated_at: {
          bsonType: "date",
          description: "Fecha de última actualización (requerido)"
        }
      }
    }
  },
  validationLevel: "moderate",
  validationAction: "error"
});

db.material_event.createIndex({ material_id: 1 }, { name: "idx_material_event_material_id" });
db.material_event.createIndex({ user_id: 1 }, { name: "idx_material_event_user_id" });
db.material_event.createIndex({ event_type: 1 }, { name: "idx_material_event_event_type" });
db.material_event.createIndex({ status: 1 }, { name: "idx_material_event_status" });
db.material_event.createIndex({ created_at: -1 }, { name: "idx_material_event_created_at" });

print("  OK: material_event creada con 5 índices");

print("Migración 001 completada. Colecciones activas: material_summary, material_assessment_worker, material_event");
