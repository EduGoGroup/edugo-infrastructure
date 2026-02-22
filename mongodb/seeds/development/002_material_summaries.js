// ============================================================
// SEED 002: Resúmenes de materiales (material_summary)
// Fecha: 2026-02-22
// Coherente con: postgres/seeds/development/
//   mat001 → "Introducción a las Fracciones"  (completed → tiene summary)
//   mat002 → "El Sistema Solar"               (completed → tiene summary)
//   mat003 → "Historia de América Latina"     (completed → tiene summary)
//   mat004 → "Álgebra Básica"                 (processing → SIN summary en Mongo)
// ============================================================

const db = db.getSiblingDB('edugo');

print("Seeding material_summary (desarrollo)...");

const MAT001 = "mat00000-0000-0000-0000-000000000001";
const MAT002 = "mat00000-0000-0000-0000-000000000002";
const MAT003 = "mat00000-0000-0000-0000-000000000003";

db.material_summary.insertMany([
  // ---- mat001: Introducción a las Fracciones ----
  {
    material_id: MAT001,
    summary: "Las fracciones son una forma de representar partes de un todo. Una fracción se compone de un numerador, que indica cuántas partes se toman, y un denominador, que indica en cuántas partes iguales se divide el todo. Las fracciones equivalentes representan la misma cantidad aunque tengan distintos numeradores y denominadores. Para simplificar fracciones se divide numerador y denominador por su máximo común divisor. Las operaciones básicas con fracciones incluyen suma, resta, multiplicación y división, cada una con sus propios procedimientos.",
    key_points: [
      "Una fracción tiene numerador (partes tomadas) y denominador (partes totales del entero)",
      "Las fracciones equivalentes representan la misma cantidad con diferentes números",
      "Para simplificar se divide por el Máximo Común Divisor (MCD)",
      "Para sumar o restar fracciones se necesita un denominador común",
      "Multiplicar fracciones: se multiplican numeradores entre sí y denominadores entre sí"
    ],
    language: "es",
    word_count: 94,
    version: 1,
    ai_model: "gpt-4-turbo-preview",
    processing_time_ms: 2340,
    token_usage: {
      prompt_tokens: 512,
      completion_tokens: 148,
      total_tokens: 660
    },
    metadata: {
      source_length: 4200,
      has_images: true
    },
    created_at: new Date("2026-02-22T10:03:00Z"),
    updated_at: new Date("2026-02-22T10:03:00Z")
  },

  // ---- mat002: El Sistema Solar ----
  {
    material_id: MAT002,
    summary: "El Sistema Solar está formado por el Sol y todos los cuerpos celestes que orbitan a su alrededor, incluyendo ocho planetas, sus lunas, planetas enanos, asteroides y cometas. Los planetas interiores son rocosos: Mercurio, Venus, Tierra y Marte. Los planetas exteriores son gaseosos o de hielo: Júpiter, Saturno, Urano y Neptuno. La Tierra es el único planeta conocido con vida, gracias a su atmósfera, agua líquida y distancia al Sol. La gravedad solar mantiene a todos los cuerpos en órbita mediante la ley de gravitación universal.",
    key_points: [
      "El Sistema Solar tiene 8 planetas divididos en interiores (rocosos) y exteriores (gaseosos/helados)",
      "El Sol contiene el 99.8% de toda la masa del Sistema Solar",
      "La Tierra es el único planeta con condiciones conocidas para la vida",
      "Júpiter es el planeta más grande y tiene más de 90 lunas conocidas",
      "Los cometas son cuerpos de hielo y roca que forman colas al acercarse al Sol"
    ],
    language: "es",
    word_count: 103,
    version: 1,
    ai_model: "gpt-4-turbo-preview",
    processing_time_ms: 2890,
    token_usage: {
      prompt_tokens: 648,
      completion_tokens: 162,
      total_tokens: 810
    },
    metadata: {
      source_length: 5800,
      has_images: true
    },
    created_at: new Date("2026-02-22T11:03:30Z"),
    updated_at: new Date("2026-02-22T11:03:30Z")
  },

  // ---- mat003: Historia de América Latina ----
  {
    material_id: MAT003,
    summary: "América Latina vivió un proceso de independencia entre 1810 y 1830, impulsado por las ideas ilustradas europeas, el ejemplo de la Revolución Francesa y las guerras napoleónicas que debilitaron a España y Portugal. Líderes como Simón Bolívar, José de San Martín, Miguel Hidalgo y José Martí encabezaron los movimientos independentistas. Tras la independencia, los nuevos países enfrentaron guerras civiles, caudillismos e inestabilidad política. En el siglo XX, la región experimentó dictaduras, revoluciones sociales y la paulatina democratización.",
    key_points: [
      "Las independencias latinoamericanas ocurrieron entre 1810 y 1830 influenciadas por la Ilustración",
      "Las guerras napoleónicas debilitaron a las metrópolis española y portuguesa facilitando la independencia",
      "Simón Bolívar liberó Venezuela, Colombia, Ecuador, Perú y Bolivia",
      "Las nuevas repúblicas sufrieron inestabilidad política y caudillismo posindependencia",
      "El siglo XX estuvo marcado por dictaduras militares y procesos de democratización"
    ],
    language: "es",
    word_count: 97,
    version: 1,
    ai_model: "gpt-4-turbo-preview",
    processing_time_ms: 3120,
    token_usage: {
      prompt_tokens: 720,
      completion_tokens: 155,
      total_tokens: 875
    },
    metadata: {
      source_length: 6400,
      has_images: false
    },
    created_at: new Date("2026-02-22T12:02:10Z"),
    updated_at: new Date("2026-02-22T12:02:10Z")
  }
]);

print("  OK: 3 summaries insertados en material_summary");
print("    - mat001 (Fracciones): language=es, word_count=94, version=1");
print("    - mat002 (Sistema Solar): language=es, word_count=103, version=1");
print("    - mat003 (Historia LATAM): language=es, word_count=97, version=1");
print("    - mat004 (Álgebra): SIN summary (processing_status=processing)");
