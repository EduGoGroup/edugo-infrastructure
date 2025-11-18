// Seeds para material_assessment_worker
// Datos de prueba para evaluaciones/quizzes generados por IA

db = db.getSiblingDB("edugo");

print("🌱 Seeding material_assessment_worker...");

db.material_assessment_worker.insertMany([
  {
    material_id: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    questions: [
      {
        question_id: "q1111111-1111-1111-1111-111111111111",
        question_text: "¿Cuál es el principio fundamental de la Programación Orientada a Objetos que permite ocultar los detalles de implementación?",
        question_type: "multiple_choice",
        options: [
          { option_id: "opt1", option_text: "Herencia" },
          { option_id: "opt2", option_text: "Polimorfismo" },
          { option_id: "opt3", option_text: "Encapsulación" },
          { option_id: "opt4", option_text: "Abstracción" }
        ],
        correct_answer: "opt3",
        explanation: "La encapsulación es el principio que permite ocultar los detalles internos de implementación y exponer solo lo necesario mediante interfaces públicas.",
        points: 10,
        difficulty: "medium",
        tags: ["POO", "conceptos"]
      },
      {
        question_id: "q2222222-2222-2222-2222-222222222222",
        question_text: "En Java, ¿una clase puede heredar de múltiples clases?",
        question_type: "true_false",
        options: [
          { option_id: "true", option_text: "Verdadero" },
          { option_id: "false", option_text: "Falso" }
        ],
        correct_answer: "false",
        explanation: "Java no soporta herencia múltiple de clases para evitar el problema del diamante. Sin embargo, una clase puede implementar múltiples interfaces.",
        points: 5,
        difficulty: "easy"
      },
      {
        question_id: "q3333333-3333-3333-3333-333333333333",
        question_text: "Explica brevemente qué es el polimorfismo y da un ejemplo en Java.",
        question_type: "open",
        options: [],
        correct_answer: "El polimorfismo permite que objetos de diferentes clases sean tratados como objetos de una clase común. Ejemplo: Animal animal = new Perro(); donde Perro extiende Animal.",
        explanation: "El polimorfismo permite escribir código más flexible y reutilizable al trabajar con abstracciones en lugar de implementaciones concretas.",
        points: 15,
        difficulty: "hard",
        tags: ["POO", "polimorfismo"]
      }
    ],
    total_questions: 3,
    total_points: 30,
    version: 1,
    ai_model: "gpt-4",
    processing_time_ms: 5200,
    token_usage: {
      prompt_tokens: 1200,
      completion_tokens: 450,
      total_tokens: 1650
    },
    created_at: new Date("2025-11-15T10:35:00Z"),
    updated_at: new Date("2025-11-15T10:35:00Z")
  },
  {
    material_id: "f1a2b3c4-d5e6-4f5a-9b8c-7d6e5f4a3b2c",
    questions: [
      {
        question_id: "q4444444-4444-4444-4444-444444444444",
        question_text: "Which React Hook is used to perform side effects in functional components?",
        question_type: "multiple_choice",
        options: [
          { option_id: "opt1", option_text: "useState" },
          { option_id: "opt2", option_text: "useEffect" },
          { option_id: "opt3", option_text: "useContext" },
          { option_id: "opt4", option_text: "useReducer" }
        ],
        correct_answer: "opt2",
        explanation: "useEffect is the Hook used to perform side effects such as data fetching, subscriptions, or manually changing the DOM.",
        points: 10,
        difficulty: "easy",
        tags: ["React", "Hooks"]
      },
      {
        question_id: "q5555555-5555-5555-5555-555555555555",
        question_text: "What is the correct syntax for creating a custom Hook in React?",
        question_type: "multiple_choice",
        options: [
          { option_id: "opt1", option_text: "function myHook() {}" },
          { option_id: "opt2", option_text: "const myHook = () => {}" },
          { option_id: "opt3", option_text: "function useMyHook() {}" },
          { option_id: "opt4", option_text: "hook myHook() {}" }
        ],
        correct_answer: "opt3",
        explanation: "Custom Hooks must start with 'use' prefix to follow React conventions and enable linting rules.",
        points: 10,
        difficulty: "medium",
        tags: ["React", "Hooks", "custom"]
      }
    ],
    total_questions: 2,
    total_points: 20,
    version: 1,
    ai_model: "gpt-4-turbo",
    processing_time_ms: 4100,
    created_at: new Date("2025-11-16T14:25:00Z"),
    updated_at: new Date("2025-11-16T14:25:00Z")
  }
]);

print("✅ 2 material_assessment_worker documents inserted");
