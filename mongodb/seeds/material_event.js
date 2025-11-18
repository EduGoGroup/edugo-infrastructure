// Seeds para material_event
// Datos de prueba para eventos de auditoría del worker

db = db.getSiblingDB("edugo");

print("🌱 Seeding material_event...");

db.material_event.insertMany([
  {
    event_type: "material_uploaded",
    material_id: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    user_id: "u1111111-1111-1111-1111-111111111111",
    payload: {
      filename: "java-poo-fundamentos.pdf",
      file_size: 1024000,
      mime_type: "application/pdf"
    },
    status: "completed",
    retry_count: 0,
    processed_at: new Date("2025-11-15T10:30:30Z"),
    created_at: new Date("2025-11-15T10:30:00Z"),
    updated_at: new Date("2025-11-15T10:30:30Z")
  },
  {
    event_type: "material_uploaded",
    material_id: "f1a2b3c4-d5e6-4f5a-9b8c-7d6e5f4a3b2c",
    user_id: "u2222222-2222-2222-2222-222222222222",
    payload: {
      filename: "react-hooks-guide.md",
      file_size: 45600,
      mime_type: "text/markdown"
    },
    status: "completed",
    retry_count: 0,
    processed_at: new Date("2025-11-16T14:20:45Z"),
    created_at: new Date("2025-11-16T14:20:00Z"),
    updated_at: new Date("2025-11-16T14:20:45Z")
  },
  {
    event_type: "material_reprocess",
    material_id: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    user_id: "u1111111-1111-1111-1111-111111111111",
    payload: {
      reason: "user_requested",
      previous_version: 1
    },
    status: "processing",
    retry_count: 0,
    created_at: new Date("2025-11-17T15:30:00Z"),
    updated_at: new Date("2025-11-17T15:30:00Z")
  },
  {
    event_type: "material_uploaded",
    material_id: "c3d4e5f6-a7b8-4c5d-9e0f-1a2b3c4d5e6f",
    user_id: "u3333333-3333-3333-3333-333333333333",
    payload: {
      filename: "data-structures.pdf",
      file_size: 2048000,
      mime_type: "application/pdf"
    },
    status: "failed",
    error_msg: "Failed to extract text from PDF: corrupted file structure",
    stack_trace: "Error: PDF parsing failed\n  at PDFParser.parse (parser.js:245)\n  at processFile (worker.js:89)",
    retry_count: 2,
    created_at: new Date("2025-11-17T16:00:00Z"),
    updated_at: new Date("2025-11-17T16:05:30Z")
  },
  {
    event_type: "assessment_attempt",
    material_id: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    user_id: "u4444444-4444-4444-4444-444444444444",
    payload: {
      assessment_id: "asmt1111-1111-1111-1111-111111111111",
      score: 25,
      total_points: 30,
      passed: true
    },
    status: "completed",
    retry_count: 0,
    processed_at: new Date("2025-11-18T09:15:00Z"),
    created_at: new Date("2025-11-18T09:10:00Z"),
    updated_at: new Date("2025-11-18T09:15:00Z")
  }
]);

print("✅ 5 material_event documents inserted");
