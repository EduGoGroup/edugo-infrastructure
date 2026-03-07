// =============================================================================
// EduGo Development Seeds v2 — MongoDB — 001_assessments.js
// =============================================================================
// Creates assessment documents in the `material_assessment_worker` collection.
// These documents contain the actual questions, options, and correct answers.
//
// Run with: mongosh "mongodb+srv://..." --file 001_assessments.js
// Or via: mongosh < 001_assessments.js
//
// The _id values here match the mongo_document_id in PostgreSQL assessments:
//   aaaaaa000000000000000001 → Examen Fracciones (5 preguntas)
//   aaaaaa000000000000000002 → Quiz Sistema Solar (4 preguntas)
//   aaaaaa000000000000000003 → Ejercicio Color y Forma (3 preguntas)
//   aaaaaa000000000000000004 → English Grammar Test (4 preguntas)
//   aaaaaa000000000000000005 → Evaluacion Historia Chile (3 preguntas, draft)
//   aaaaaa000000000000000006 → Proyecto Final Escultura (0 preguntas, draft)
// =============================================================================

// Switch to the correct database
const dbName = "edugo";
const db = db.getSiblingDB ? db.getSiblingDB(dbName) : use(dbName);
const collection = db.getCollection("material_assessment_worker");

// Clean existing seed documents
collection.deleteMany({
  _id: {
    $in: [
      ObjectId("aaaaaa000000000000000001"),
      ObjectId("aaaaaa000000000000000002"),
      ObjectId("aaaaaa000000000000000003"),
      ObjectId("aaaaaa000000000000000004"),
      ObjectId("aaaaaa000000000000000005"),
      ObjectId("aaaaaa000000000000000006"),
    ],
  },
});

// =========================================================================
// ass001: Examen Fracciones — 5 preguntas, 20 pts c/u = 100 pts max
// =========================================================================
collection.insertOne({
  _id: ObjectId("aaaaaa000000000000000001"),
  material_id: "mat-fracciones-001",
  total_points: NumberInt(100),
  total_questions: NumberInt(5),
  version: NumberInt(1),
  ai_model: "manual",
  questions: [
    {
      question_id: "q-frac-001",
      question_text: "Cuanto es 1/4 + 1/4?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "1/2" },
        { option_id: "B", option_text: "2/8" },
        { option_id: "C", option_text: "1/4" },
        { option_id: "D", option_text: "2/4" },
      ],
      correct_answer: "A",
      explanation: "1/4 + 1/4 = 2/4 = 1/2",
    },
    {
      question_id: "q-frac-002",
      question_text: "Cuanto es 1/4 + 2/4?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "3/8" },
        { option_id: "B", option_text: "3/4" },
        { option_id: "C", option_text: "1/2" },
        { option_id: "D", option_text: "2/4" },
      ],
      correct_answer: "B",
      explanation: "1/4 + 2/4 = 3/4",
    },
    {
      question_id: "q-frac-003",
      question_text: "Cual fraccion es equivalente a 2/6?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "medium",
      options: [
        { option_id: "A", option_text: "1/3" },
        { option_id: "B", option_text: "2/6" },
        { option_id: "C", option_text: "1/2" },
        { option_id: "D", option_text: "3/6" },
      ],
      correct_answer: "A",
      explanation: "2/6 simplificado es 1/3",
    },
    {
      question_id: "q-frac-004",
      question_text: "Cuanto es 1/5 + 1/5?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "1/5" },
        { option_id: "B", option_text: "2/5" },
        { option_id: "C", option_text: "2/10" },
        { option_id: "D", option_text: "1/10" },
      ],
      correct_answer: "B",
      explanation: "1/5 + 1/5 = 2/5",
    },
    {
      question_id: "q-frac-005",
      question_text: "Cuanto es 1/8 + 2/8?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "3/16" },
        { option_id: "B", option_text: "2/8" },
        { option_id: "C", option_text: "3/8" },
        { option_id: "D", option_text: "1/4" },
      ],
      correct_answer: "C",
      explanation: "1/8 + 2/8 = 3/8",
    },
  ],
  created_at: new Date(),
  updated_at: new Date(),
});

// =========================================================================
// ass002: Quiz Sistema Solar — 4 preguntas, 20 pts c/u = 80 pts max
// =========================================================================
collection.insertOne({
  _id: ObjectId("aaaaaa000000000000000002"),
  material_id: "mat-sistema-solar-001",
  total_points: NumberInt(80),
  total_questions: NumberInt(4),
  version: NumberInt(1),
  ai_model: "manual",
  questions: [
    {
      question_id: "q-solar-001",
      question_text: "Cual es el planeta mas grande del sistema solar?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "Saturno" },
        { option_id: "B", option_text: "Jupiter" },
        { option_id: "C", option_text: "Neptuno" },
        { option_id: "D", option_text: "Urano" },
      ],
      correct_answer: "B",
      explanation: "Jupiter es el planeta mas grande del sistema solar",
    },
    {
      question_id: "q-solar-002",
      question_text: "Cual es el planeta mas cercano al Sol?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "Venus" },
        { option_id: "B", option_text: "Marte" },
        { option_id: "C", option_text: "Mercurio" },
        { option_id: "D", option_text: "Tierra" },
      ],
      correct_answer: "C",
      explanation: "Mercurio es el planeta mas cercano al Sol",
    },
    {
      question_id: "q-solar-003",
      question_text: "Cuantos planetas hay en el sistema solar?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "7" },
        { option_id: "B", option_text: "8" },
        { option_id: "C", option_text: "9" },
        { option_id: "D", option_text: "10" },
      ],
      correct_answer: "B",
      explanation: "Hay 8 planetas en el sistema solar",
    },
    {
      question_id: "q-solar-004",
      question_text: "Que planeta es conocido como el planeta rojo?",
      question_type: "multiple_choice",
      points: NumberInt(20),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "Jupiter" },
        { option_id: "B", option_text: "Venus" },
        { option_id: "C", option_text: "Marte" },
        { option_id: "D", option_text: "Saturno" },
      ],
      correct_answer: "C",
      explanation: "Marte es conocido como el planeta rojo",
    },
  ],
  created_at: new Date(),
  updated_at: new Date(),
});

// =========================================================================
// ass003: Ejercicio Color y Forma — 3 preguntas, ~33 pts c/u = 100 pts max
// =========================================================================
collection.insertOne({
  _id: ObjectId("aaaaaa000000000000000003"),
  material_id: "mat-color-forma-001",
  total_points: NumberInt(100),
  total_questions: NumberInt(3),
  version: NumberInt(1),
  ai_model: "manual",
  questions: [
    {
      question_id: "q-color-001",
      question_text: "Cuales son los colores primarios en pintura?",
      question_type: "multiple_choice",
      points: NumberInt(34),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "Rojo, Azul, Amarillo" },
        { option_id: "B", option_text: "Rojo, Verde, Azul" },
        { option_id: "C", option_text: "Naranja, Violeta, Verde" },
        { option_id: "D", option_text: "Blanco, Negro, Gris" },
      ],
      correct_answer: "A",
      explanation: "Los colores primarios en pintura son Rojo, Azul y Amarillo",
    },
    {
      question_id: "q-color-002",
      question_text: "Que color se obtiene al mezclar azul y amarillo?",
      question_type: "multiple_choice",
      points: NumberInt(33),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "Naranja" },
        { option_id: "B", option_text: "Verde" },
        { option_id: "C", option_text: "Violeta" },
        { option_id: "D", option_text: "Marron" },
      ],
      correct_answer: "B",
      explanation: "Azul + Amarillo = Verde",
    },
    {
      question_id: "q-color-003",
      question_text: "Que son los colores complementarios?",
      question_type: "multiple_choice",
      points: NumberInt(33),
      difficulty: "medium",
      options: [
        { option_id: "A", option_text: "Colores que estan uno al lado del otro en el circulo cromatico" },
        { option_id: "B", option_text: "Colores que estan opuestos en el circulo cromatico" },
        { option_id: "C", option_text: "Colores que tienen el mismo tono" },
        { option_id: "D", option_text: "Colores que son variaciones de un mismo color" },
      ],
      correct_answer: "B",
      explanation: "Los colores complementarios son los que estan opuestos en el circulo cromatico",
    },
  ],
  created_at: new Date(),
  updated_at: new Date(),
});

// =========================================================================
// ass004: English Grammar Test — 4 preguntas, 25 pts c/u = 100 pts max
// =========================================================================
collection.insertOne({
  _id: ObjectId("aaaaaa000000000000000004"),
  material_id: "mat-english-grammar-001",
  total_points: NumberInt(100),
  total_questions: NumberInt(4),
  version: NumberInt(1),
  ai_model: "manual",
  questions: [
    {
      question_id: "q-eng-001",
      question_text: "Choose the correct article: ___ apple a day keeps the doctor away.",
      question_type: "multiple_choice",
      points: NumberInt(25),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "A" },
        { option_id: "B", option_text: "An" },
        { option_id: "C", option_text: "The" },
        { option_id: "D", option_text: "No article" },
      ],
      correct_answer: "B",
      explanation: "An is used before words starting with a vowel sound",
    },
    {
      question_id: "q-eng-002",
      question_text: "Which is the correct form? She ___ to school every day.",
      question_type: "multiple_choice",
      points: NumberInt(25),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "go" },
        { option_id: "B", option_text: "goes" },
        { option_id: "C", option_text: "going" },
        { option_id: "D", option_text: "gone" },
      ],
      correct_answer: "B",
      explanation: "Third person singular uses goes in simple present",
    },
    {
      question_id: "q-eng-003",
      question_text: "What is the plural of child?",
      question_type: "multiple_choice",
      points: NumberInt(25),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "childs" },
        { option_id: "B", option_text: "childrens" },
        { option_id: "C", option_text: "children" },
        { option_id: "D", option_text: "childes" },
      ],
      correct_answer: "C",
      explanation: "Children is the irregular plural of child",
    },
    {
      question_id: "q-eng-004",
      question_text: "Choose the correct pronoun: John gave the book to ___.",
      question_type: "multiple_choice",
      points: NumberInt(25),
      difficulty: "easy",
      options: [
        { option_id: "A", option_text: "I" },
        { option_id: "B", option_text: "me" },
        { option_id: "C", option_text: "my" },
        { option_id: "D", option_text: "mine" },
      ],
      correct_answer: "B",
      explanation: "Me is the object pronoun used after prepositions",
    },
  ],
  created_at: new Date(),
  updated_at: new Date(),
});

// =========================================================================
// ass005: Evaluacion Historia Chile — 3 preguntas (draft, pero con contenido)
// =========================================================================
collection.insertOne({
  _id: ObjectId("aaaaaa000000000000000005"),
  material_id: "mat-historia-chile-001",
  total_points: NumberInt(100),
  total_questions: NumberInt(3),
  version: NumberInt(1),
  ai_model: "manual",
  questions: [
    {
      question_id: "q-hist-001",
      question_text: "En que ano se firmo el Acta de Independencia de Chile?",
      question_type: "multiple_choice",
      points: NumberInt(34),
      difficulty: "medium",
      options: [
        { option_id: "A", option_text: "1810" },
        { option_id: "B", option_text: "1818" },
        { option_id: "C", option_text: "1821" },
        { option_id: "D", option_text: "1826" },
      ],
      correct_answer: "B",
      explanation: "El Acta de Independencia de Chile se firmo en 1818",
    },
    {
      question_id: "q-hist-002",
      question_text: "Quien fue el Director Supremo que firmo la independencia?",
      question_type: "multiple_choice",
      points: NumberInt(33),
      difficulty: "medium",
      options: [
        { option_id: "A", option_text: "Jose Miguel Carrera" },
        { option_id: "B", option_text: "Bernardo O'Higgins" },
        { option_id: "C", option_text: "Manuel Rodriguez" },
        { option_id: "D", option_text: "Simon Bolivar" },
      ],
      correct_answer: "B",
      explanation: "Bernardo O'Higgins fue el Director Supremo que firmo la independencia",
    },
    {
      question_id: "q-hist-003",
      question_text: "Que batalla fue decisiva para la independencia de Chile?",
      question_type: "multiple_choice",
      points: NumberInt(33),
      difficulty: "hard",
      options: [
        { option_id: "A", option_text: "Batalla de Rancagua" },
        { option_id: "B", option_text: "Batalla de Chacabuco" },
        { option_id: "C", option_text: "Batalla de Maipu" },
        { option_id: "D", option_text: "Batalla de Ayacucho" },
      ],
      correct_answer: "C",
      explanation: "La Batalla de Maipu en 1818 fue decisiva para la independencia de Chile",
    },
  ],
  created_at: new Date(),
  updated_at: new Date(),
});

// =========================================================================
// ass006: Proyecto Final Escultura — 0 preguntas (draft vacio)
// =========================================================================
collection.insertOne({
  _id: ObjectId("aaaaaa000000000000000006"),
  material_id: "mat-escultura-001",
  total_points: NumberInt(0),
  total_questions: NumberInt(0),
  version: NumberInt(1),
  ai_model: "manual",
  questions: [],
  created_at: new Date(),
  updated_at: new Date(),
});

print("6 assessment documents inserted successfully");
