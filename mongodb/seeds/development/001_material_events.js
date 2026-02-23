// ============================================================
// SEED 001: Eventos de materiales (material_event)
// Fecha: 2026-02-22
// Coherente con: postgres/seeds/development/
//   mat001 → "Introducción a las Fracciones"  (processing_status=completed)
//   mat002 → "El Sistema Solar"               (processing_status=completed)
//   mat003 → "Historia de América Latina"     (processing_status=completed)
//   mat004 → "Álgebra Básica"                 (processing_status=processing → evento en progreso)
// ============================================================

const db = db.getSiblingDB('edugo');

print("Seeding material_event (desarrollo)...");

// IDs de materiales y usuarios coherentes con seeds SQL
const MAT001 = "mat00000-0000-0000-0000-000000000001";
const MAT002 = "mat00000-0000-0000-0000-000000000002";
const MAT003 = "mat00000-0000-0000-0000-000000000003";
const MAT004 = "mat00000-0000-0000-0000-000000000004";

// Usuario docente que subió los materiales (coherente con seeds SQL)
const USER_TEACHER = "usr00000-0000-0000-0000-000000000002";

const now = new Date("2026-02-22T10:00:00Z");

db.material_event.insertMany([
  // ---- mat001: Introducción a las Fracciones → completed ----
  {
    event_type: "material_uploaded",
    material_id: MAT001,
    user_id: USER_TEACHER,
    payload: {
      material_title: "Introducción a las Fracciones",
      file_url: "https://storage.edugo.dev/materials/mat001_fracciones.pdf",
      file_size_bytes: 524288,
      mime_type: "application/pdf",
      requested_tasks: ["summary", "assessment"]
    },
    status: "completed",
    error_msg: null,
    stack_trace: null,
    retry_count: 0,
    processed_at: new Date("2026-02-22T10:05:30Z"),
    created_at: new Date("2026-02-22T10:00:00Z"),
    updated_at: new Date("2026-02-22T10:05:30Z")
  },

  // ---- mat002: El Sistema Solar → completed ----
  {
    event_type: "material_uploaded",
    material_id: MAT002,
    user_id: USER_TEACHER,
    payload: {
      material_title: "El Sistema Solar",
      file_url: "https://storage.edugo.dev/materials/mat002_sistema_solar.pdf",
      file_size_bytes: 786432,
      mime_type: "application/pdf",
      requested_tasks: ["summary", "assessment"]
    },
    status: "completed",
    error_msg: null,
    stack_trace: null,
    retry_count: 0,
    processed_at: new Date("2026-02-22T11:08:15Z"),
    created_at: new Date("2026-02-22T11:00:00Z"),
    updated_at: new Date("2026-02-22T11:08:15Z")
  },

  // ---- mat003: Historia de América Latina → completed (solo summary, sin assessment publicado) ----
  {
    event_type: "material_uploaded",
    material_id: MAT003,
    user_id: USER_TEACHER,
    payload: {
      material_title: "Historia de América Latina",
      file_url: "https://storage.edugo.dev/materials/mat003_historia_latam.pdf",
      file_size_bytes: 655360,
      mime_type: "application/pdf",
      requested_tasks: ["summary"]
    },
    status: "completed",
    error_msg: null,
    stack_trace: null,
    retry_count: 0,
    processed_at: new Date("2026-02-22T12:04:45Z"),
    created_at: new Date("2026-02-22T12:00:00Z"),
    updated_at: new Date("2026-02-22T12:04:45Z")
  },

  // ---- mat004: Álgebra Básica → processing (en progreso, sin docs MongoDB) ----
  {
    event_type: "material_uploaded",
    material_id: MAT004,
    user_id: USER_TEACHER,
    payload: {
      material_title: "Álgebra Básica",
      file_url: "https://storage.edugo.dev/materials/mat004_algebra.pdf",
      file_size_bytes: 471040,
      mime_type: "application/pdf",
      requested_tasks: ["summary", "assessment"]
    },
    status: "processing",
    error_msg: null,
    stack_trace: null,
    retry_count: 0,
    processed_at: null,
    created_at: new Date("2026-02-22T13:00:00Z"),
    updated_at: new Date("2026-02-22T13:00:00Z")
  }
]);

print("  OK: 4 eventos insertados en material_event");
print("    - mat001 (Fracciones): completed");
print("    - mat002 (Sistema Solar): completed");
print("    - mat003 (Historia LATAM): completed");
print("    - mat004 (Álgebra): processing (sin docs Mongo)");
