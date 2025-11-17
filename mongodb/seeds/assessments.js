// Seeds de assessments en MongoDB
// Ejecutar con: mongosh --host localhost:27017/edugo < assessments.js

use edugo;

// Assessment 1 (para material de Física)
db.material_assessment.insertOne({
  _id: ObjectId("507f1f77bcf86cd799439011"),
  material_id: "66666666-6666-6666-6666-666666666666",
  questions: [
    {
      question_index: 0,
      question_text: "¿Qué es la dualidad onda-partícula?",
      question_type: "multiple_choice",
      options: [
        { option_index: 0, text: "Partículas que actúan solo como ondas", is_correct: false },
        { option_index: 1, text: "Partículas que pueden comportarse como ondas y viceversa", is_correct: true },
        { option_index: 2, text: "Ondas que no son partículas", is_correct: false },
        { option_index: 3, text: "Ninguna de las anteriores", is_correct: false }
      ]
    },
    {
      question_index: 1,
      question_text: "¿Quién propuso el principio de incertidumbre?",
      question_type: "multiple_choice",
      options: [
        { option_index: 0, text: "Einstein", is_correct: false },
        { option_index: 1, text: "Heisenberg", is_correct: true },
        { option_index: 2, text: "Bohr", is_correct: false },
        { option_index: 3, text: "Schrödinger", is_correct: false }
      ]
    }
  ],
  metadata: {
    subject: "Física",
    grade: "10th",
    difficulty: "medium"
  },
  created_at: new Date(),
  updated_at: new Date()
});

// Assessment 2 (para material de Álgebra)
db.material_assessment.insertOne({
  _id: ObjectId("507f1f77bcf86cd799439012"),
  material_id: "77777777-7777-7777-7777-777777777777",
  questions: [
    {
      question_index: 0,
      question_text: "¿Qué es una matriz identidad?",
      question_type: "multiple_choice",
      options: [
        { option_index: 0, text: "Una matriz con todos 1s", is_correct: false },
        { option_index: 1, text: "Una matriz con 1s en la diagonal y 0s en el resto", is_correct: true },
        { option_index: 2, text: "Una matriz cuadrada", is_correct: false },
        { option_index: 3, text: "Una matriz invertible", is_correct: false }
      ]
    }
  ],
  metadata: {
    subject: "Matemáticas",
    grade: "11th",
    difficulty: "easy"
  },
  created_at: new Date(),
  updated_at: new Date()
});

print("✅ 2 assessments insertados correctamente");
