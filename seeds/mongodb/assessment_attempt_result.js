// Seeds for assessment_attempt_result collection
// Execute with: mongosh edugo < assessment_attempt_result.js
// Or: mongosh --eval "$(cat assessment_attempt_result.js)"

db = db.getSiblingDB('edugo');

// Attempt result 1 (Student completed Physics assessment)
db.assessment_attempt_result.insertOne({
  attempt_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
  student_id: "33333333-3333-3333-3333-333333333333",
  assessment_id: "99999999-9999-9999-9999-999999999999",
  answers: [
    {
      question_index: 0,
      question_text: "¿Qué es la dualidad onda-partícula?",
      selected_option_index: 1,
      selected_option_text: "Partículas que pueden comportarse como ondas y viceversa",
      correct_option_index: 1,
      is_correct: true,
      time_spent_seconds: 45,
      answered_at: new Date("2025-01-15T10:15:00Z")
    },
    {
      question_index: 1,
      question_text: "¿Quién propuso el principio de incertidumbre?",
      selected_option_index: 1,
      selected_option_text: "Heisenberg",
      correct_option_index: 1,
      is_correct: true,
      time_spent_seconds: 30,
      answered_at: new Date("2025-01-15T10:16:00Z")
    }
  ],
  score: {
    correct_count: 2,
    incorrect_count: 0,
    total_questions: 2,
    percentage: 100.0
  },
  time_tracking: {
    total_time_seconds: 75,
    average_time_per_question: 37.5
  },
  started_at: new Date("2025-01-15T10:14:00Z"),
  submitted_at: new Date("2025-01-15T10:16:15Z"),
  created_at: new Date("2025-01-15T10:16:15Z")
});

// Attempt result 2 (Student partially completed Algebra assessment)
db.assessment_attempt_result.insertOne({
  attempt_id: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
  student_id: "44444444-4444-4444-4444-444444444444",
  assessment_id: "88888888-8888-8888-8888-888888888888",
  answers: [
    {
      question_index: 0,
      question_text: "¿Qué es una matriz identidad?",
      selected_option_index: 0,
      selected_option_text: "Una matriz con todos 1s",
      correct_option_index: 1,
      is_correct: false,
      time_spent_seconds: 60,
      answered_at: new Date("2025-01-15T11:20:00Z")
    }
  ],
  score: {
    correct_count: 0,
    incorrect_count: 1,
    total_questions: 1,
    percentage: 0.0
  },
  time_tracking: {
    total_time_seconds: 60,
    average_time_per_question: 60.0
  },
  started_at: new Date("2025-01-15T11:19:00Z"),
  submitted_at: new Date("2025-01-15T11:20:00Z"),
  created_at: new Date("2025-01-15T11:20:00Z")
});

print("✅ 2 assessment_attempt_result documents inserted successfully");
