// ============================================================
// SEED 003: Assessments de materiales (material_assessment_worker)
// Fecha: 2026-02-22
// Coherente con: postgres/seeds/development/
//   Solo materiales con assessment published en PostgreSQL:
//   mat001 → "Introducción a las Fracciones" → mongo_document_id = "mongo_assessment_mat001"
//   mat002 → "El Sistema Solar"              → mongo_document_id = "mongo_assessment_mat002"
//   mat003 → "Historia de América Latina"    → solo summary (SIN assessment en Mongo)
//   mat004 → "Álgebra Básica"               → processing (SIN assessment en Mongo)
//
// question_ids coherentes con assessment_attempt_answer del SQL:
//   mat001: q001_mat001, q002_mat001, q003_mat001, q004_mat001, q005_mat001
//   mat002: q001_mat002, q002_mat002, q003_mat002, q004_mat002, q005_mat002
// ============================================================

const db = db.getSiblingDB('edugo');

print("Seeding material_assessment_worker (desarrollo)...");

const MAT001 = "mat00000-0000-0000-0000-000000000001";
const MAT002 = "mat00000-0000-0000-0000-000000000002";

db.material_assessment_worker.insertMany([
  // ============================================================
  // ---- mat001: Introducción a las Fracciones ----
  // _id string "mongo_assessment_mat001" referenciado desde assessment.mongo_document_id en PostgreSQL
  // ============================================================
  {
    _id: "mongo_assessment_mat001",
    material_id: MAT001,
    questions: [
      {
        question_id: "q001_mat001",
        question_text: "¿Qué representa el denominador en una fracción?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "El número de partes que se toman del entero" },
          { option_id: "b", option_text: "El número total de partes iguales en que se divide el entero" },
          { option_id: "c", option_text: "El resultado de dividir dos números" },
          { option_id: "d", option_text: "La cantidad de fracciones equivalentes que existen" }
        ],
        correct_answer: "b",
        explanation: "El denominador indica en cuántas partes iguales se divide el todo. Por ejemplo, en 3/4 el denominador 4 indica que el entero está dividido en 4 partes iguales.",
        points: 10,
        difficulty: "easy",
        tags: ["fracciones", "conceptos_básicos"]
      },
      {
        question_id: "q002_mat001",
        question_text: "¿Cuál de las siguientes fracciones es equivalente a 1/2?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "2/6" },
          { option_id: "b", option_text: "3/4" },
          { option_id: "c", option_text: "4/8" },
          { option_id: "d", option_text: "5/12" }
        ],
        correct_answer: "c",
        explanation: "4/8 es equivalente a 1/2 porque al dividir numerador y denominador entre 4 obtenemos 1/2. Las fracciones equivalentes representan la misma cantidad del entero.",
        points: 10,
        difficulty: "easy",
        tags: ["fracciones_equivalentes"]
      },
      {
        question_id: "q003_mat001",
        question_text: "¿Cuál es el resultado de sumar 1/4 + 2/4?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "3/8" },
          { option_id: "b", option_text: "3/4" },
          { option_id: "c", option_text: "2/8" },
          { option_id: "d", option_text: "1/2" }
        ],
        correct_answer: "b",
        explanation: "Al sumar fracciones con el mismo denominador se suman solo los numeradores: 1/4 + 2/4 = 3/4. El denominador permanece igual.",
        points: 10,
        difficulty: "medium",
        tags: ["suma_fracciones", "mismo_denominador"]
      },
      {
        question_id: "q004_mat001",
        question_text: "Para simplificar la fracción 6/9, ¿por qué número se debe dividir numerador y denominador?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "2" },
          { option_id: "b", option_text: "3" },
          { option_id: "c", option_text: "6" },
          { option_id: "d", option_text: "9" }
        ],
        correct_answer: "b",
        explanation: "El Máximo Común Divisor (MCD) de 6 y 9 es 3. Al dividir 6/3 = 2 y 9/3 = 3, la fracción simplificada es 2/3.",
        points: 10,
        difficulty: "medium",
        tags: ["simplificación", "mcd"]
      },
      {
        question_id: "q005_mat001",
        question_text: "¿Cuál es el resultado de multiplicar 2/3 × 3/5?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "5/8" },
          { option_id: "b", option_text: "6/8" },
          { option_id: "c", option_text: "6/15" },
          { option_id: "d", option_text: "2/5" }
        ],
        correct_answer: "c",
        explanation: "Para multiplicar fracciones se multiplican los numeradores entre sí y los denominadores entre sí: (2×3)/(3×5) = 6/15. Esta fracción se puede simplificar a 2/5 dividiendo entre 3.",
        points: 10,
        difficulty: "hard",
        tags: ["multiplicación_fracciones"]
      }
    ],
    total_questions: 5,
    total_points: 50,
    version: 1,
    ai_model: "gpt-4-turbo-preview",
    processing_time_ms: 5240,
    token_usage: {
      prompt_tokens: 820,
      completion_tokens: 642,
      total_tokens: 1462
    },
    metadata: {
      average_difficulty: "medium",
      estimated_time_min: 10,
      source_length: 4200,
      has_images: true
    },
    created_at: new Date("2026-02-22T10:05:00Z"),
    updated_at: new Date("2026-02-22T10:05:00Z")
  },

  // ============================================================
  // ---- mat002: El Sistema Solar ----
  // _id string "mongo_assessment_mat002" referenciado desde assessment.mongo_document_id en PostgreSQL
  // ============================================================
  {
    _id: "mongo_assessment_mat002",
    material_id: MAT002,
    questions: [
      {
        question_id: "q001_mat002",
        question_text: "¿Cuántos planetas tiene el Sistema Solar?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "7" },
          { option_id: "b", option_text: "8" },
          { option_id: "c", option_text: "9" },
          { option_id: "d", option_text: "10" }
        ],
        correct_answer: "b",
        explanation: "El Sistema Solar tiene 8 planetas desde que Plutón fue reclasificado como planeta enano en 2006 por la Unión Astronómica Internacional.",
        points: 10,
        difficulty: "easy",
        tags: ["sistema_solar", "planetas"]
      },
      {
        question_id: "q002_mat002",
        question_text: "¿Cuál es el planeta más grande del Sistema Solar?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "Saturno" },
          { option_id: "b", option_text: "Neptuno" },
          { option_id: "c", option_text: "Júpiter" },
          { option_id: "d", option_text: "Urano" }
        ],
        correct_answer: "c",
        explanation: "Júpiter es el planeta más grande del Sistema Solar, con una masa más de 300 veces mayor que la de la Tierra. Es un gigante gaseoso compuesto principalmente de hidrógeno y helio.",
        points: 10,
        difficulty: "easy",
        tags: ["planetas", "gigantes_gaseosos"]
      },
      {
        question_id: "q003_mat002",
        question_text: "¿Qué planeta es conocido por sus anillos prominentes visibles desde la Tierra?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "Urano" },
          { option_id: "b", option_text: "Neptuno" },
          { option_id: "c", option_text: "Júpiter" },
          { option_id: "d", option_text: "Saturno" }
        ],
        correct_answer: "d",
        explanation: "Saturno es famoso por su espectacular sistema de anillos compuesto por hielo y roca. Aunque otros planetas gigantes también tienen anillos, los de Saturno son los más visibles y extensos.",
        points: 10,
        difficulty: "easy",
        tags: ["saturno", "anillos"]
      },
      {
        question_id: "q004_mat002",
        question_text: "¿Cuáles son los planetas interiores (rocosos) del Sistema Solar?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "Júpiter, Saturno, Urano, Neptuno" },
          { option_id: "b", option_text: "Mercurio, Venus, Tierra, Marte" },
          { option_id: "c", option_text: "Mercurio, Venus, Tierra, Júpiter" },
          { option_id: "d", option_text: "Tierra, Marte, Júpiter, Saturno" }
        ],
        correct_answer: "b",
        explanation: "Los planetas interiores o terrestres son Mercurio, Venus, Tierra y Marte. Son planetas rocosos, más pequeños y más cercanos al Sol. Los planetas exteriores (Júpiter, Saturno, Urano, Neptuno) son gigantes gaseosos o de hielo.",
        points: 10,
        difficulty: "medium",
        tags: ["planetas_interiores", "planetas_rocosos"]
      },
      {
        question_id: "q005_mat002",
        question_text: "¿Qué porcentaje aproximado de la masa total del Sistema Solar contiene el Sol?",
        question_type: "multiple_choice",
        options: [
          { option_id: "a", option_text: "50%" },
          { option_id: "b", option_text: "75%" },
          { option_id: "c", option_text: "90%" },
          { option_id: "d", option_text: "99.8%" }
        ],
        correct_answer: "d",
        explanation: "El Sol contiene aproximadamente el 99.8% de toda la masa del Sistema Solar. Su enorme masa es la que genera la gravedad suficiente para mantener en órbita a todos los planetas, lunas y demás cuerpos celestes.",
        points: 10,
        difficulty: "hard",
        tags: ["sol", "masa", "gravedad"]
      }
    ],
    total_questions: 5,
    total_points: 50,
    version: 1,
    ai_model: "gpt-4-turbo-preview",
    processing_time_ms: 6110,
    token_usage: {
      prompt_tokens: 940,
      completion_tokens: 718,
      total_tokens: 1658
    },
    metadata: {
      average_difficulty: "medium",
      estimated_time_min: 10,
      source_length: 5800,
      has_images: true
    },
    created_at: new Date("2026-02-22T11:07:00Z"),
    updated_at: new Date("2026-02-22T11:07:00Z")
  }
]);

print("  OK: 2 assessments insertados en material_assessment_worker");
print("    - mongo_assessment_mat001 (Fracciones): 5 preguntas, 50 pts");
print("    - mongo_assessment_mat002 (Sistema Solar): 5 preguntas, 50 pts");
print("    - mat003 (Historia LATAM): SIN assessment (solo summary)");
print("    - mat004 (Álgebra): SIN assessment (processing_status=processing)");
