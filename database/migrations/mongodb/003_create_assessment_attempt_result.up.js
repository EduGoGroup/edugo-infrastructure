// Migration: Create assessment_attempt_result collection
// Collection: assessment_attempt_result (Owner: infrastructure)
// Created by: edugo-infrastructure
// Used by: api-mobile
//
// Purpose: Stores detailed results and answers from student assessment attempts
// Related PostgreSQL tables: assessment_attempt (stores metadata), assessment_attempt_answer (stores individual answers)

db.createCollection("assessment_attempt_result", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["attempt_id", "student_id", "assessment_id", "answers", "score", "started_at", "submitted_at", "created_at"],
      properties: {
        attempt_id: {
          bsonType: "string",
          description: "UUID of the attempt in PostgreSQL (assessment_attempt table)"
        },
        student_id: {
          bsonType: "string",
          description: "UUID of the student (users table)"
        },
        assessment_id: {
          bsonType: "string",
          description: "UUID of the assessment (assessment table in PostgreSQL)"
        },
        answers: {
          bsonType: "array",
          description: "Detailed answers provided by the student",
          items: {
            bsonType: "object",
            required: ["question_index", "selected_option_index", "is_correct", "time_spent_seconds"],
            properties: {
              question_index: {
                bsonType: "int",
                minimum: 0,
                description: "Index of the question answered"
              },
              question_text: {
                bsonType: "string",
                description: "Text of the question (snapshot)"
              },
              selected_option_index: {
                bsonType: "int",
                minimum: 0,
                description: "Index of the option selected by student"
              },
              selected_option_text: {
                bsonType: "string",
                description: "Text of the selected option (snapshot)"
              },
              correct_option_index: {
                bsonType: "int",
                minimum: 0,
                description: "Index of the correct option"
              },
              is_correct: {
                bsonType: "bool",
                description: "Whether the answer was correct"
              },
              time_spent_seconds: {
                bsonType: "int",
                minimum: 0,
                description: "Time spent on this question in seconds"
              },
              answered_at: {
                bsonType: "date",
                description: "Timestamp when this question was answered"
              }
            }
          }
        },
        score: {
          bsonType: "object",
          required: ["correct_count", "total_questions", "percentage"],
          properties: {
            correct_count: {
              bsonType: "int",
              minimum: 0,
              description: "Number of correct answers"
            },
            incorrect_count: {
              bsonType: "int",
              minimum: 0,
              description: "Number of incorrect answers"
            },
            total_questions: {
              bsonType: "int",
              minimum: 1,
              description: "Total number of questions"
            },
            percentage: {
              bsonType: "double",
              minimum: 0,
              maximum: 100,
              description: "Score as percentage (0-100)"
            }
          }
        },
        time_tracking: {
          bsonType: "object",
          properties: {
            total_time_seconds: {
              bsonType: "int",
              minimum: 0,
              description: "Total time spent on the attempt"
            },
            average_time_per_question: {
              bsonType: "double",
              minimum: 0,
              description: "Average time per question in seconds"
            }
          }
        },
        started_at: {
          bsonType: "date",
          description: "When the attempt was started"
        },
        submitted_at: {
          bsonType: "date",
          description: "When the attempt was submitted"
        },
        created_at: {
          bsonType: "date",
          description: "Timestamp when this result was created"
        }
      }
    }
  }
});

// Create indexes for efficient queries
db.assessment_attempt_result.createIndex({ "attempt_id": 1 }, { name: "idx_attempt_id", unique: true });
db.assessment_attempt_result.createIndex({ "student_id": 1 }, { name: "idx_student_id" });
db.assessment_attempt_result.createIndex({ "assessment_id": 1 }, { name: "idx_assessment_id" });
db.assessment_attempt_result.createIndex({ "student_id": 1, "assessment_id": 1 }, { name: "idx_student_assessment" });
db.assessment_attempt_result.createIndex({ "submitted_at": -1 }, { name: "idx_submitted_at_desc" });
db.assessment_attempt_result.createIndex({ "score.percentage": -1 }, { name: "idx_score_percentage_desc" });

print("✅ Collection 'assessment_attempt_result' created successfully");
