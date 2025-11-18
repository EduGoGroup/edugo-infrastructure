// 009_create_material_event.up.js
// Migración: Crear collection material_event para auditoría de eventos (TTL 90 días)
// Proyecto: edugo-infrastructure  
// Consumidor: edugo-worker
// Fecha: 2025-11-18

print("🔧 Creating collection: material_event...");

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
          enum: [
            "material_uploaded",
            "material_reprocess",
            "material_deleted",
            "assessment_attempt",
            "student_enrolled",
            "student_unenrolled"
          ],
          description: "Tipo de evento procesado"
        },
        material_id: {
          bsonType: "string",
          pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
          description: "UUID v4 del material (opcional)"
        },
        user_id: {
          bsonType: "string",
          pattern: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
          description: "UUID v4 del usuario (opcional)"
        },
        payload: {
          bsonType: "object",
          description: "Payload flexible del evento (requerido)"
        },
        status: {
          enum: ["pending", "processing", "completed", "failed"],
          description: "Estado del procesamiento del evento"
        },
        error_msg: {
          bsonType: "string",
          maxLength: 5000,
          description: "Mensaje de error si status=failed (opcional)"
        },
        stack_trace: {
          bsonType: "string",
          maxLength: 10000,
          description: "Stack trace del error (opcional)"
        },
        retry_count: {
          bsonType: "int",
          minimum: 0,
          description: "Número de reintentos (>= 0)"
        },
        processed_at: {
          bsonType: "date",
          description: "Fecha de procesamiento exitoso (opcional)"
        },
        created_at: {
          bsonType: "date",
          description: "Fecha de creación del evento (requerido)"
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
print("📑 Creating indexes for material_event...");

db.material_event.createIndex(
  { event_type: 1 },
  { name: "idx_event_type" }
);

db.material_event.createIndex(
  { material_id: 1 },
  { name: "idx_material_id" }
);

db.material_event.createIndex(
  { status: 1 },
  { name: "idx_status" }
);

db.material_event.createIndex(
  { created_at: -1 },
  { name: "idx_created_at" }
);

db.material_event.createIndex(
  { processed_at: -1 },
  { name: "idx_processed_at" }
);

db.material_event.createIndex(
  { status: 1, created_at: -1 },
  { name: "idx_status_created" }
);

// TTL Index: Auto-eliminar documentos después de 90 días
db.material_event.createIndex(
  { created_at: 1 },
  {
    expireAfterSeconds: 7776000,  // 90 días = 7776000 segundos
    name: "idx_ttl_created_at"
  }
);

print("✅ Collection material_event created successfully with 7 indexes (including TTL)");
print("⏰ TTL configured: Documents will be auto-deleted after 90 days");
