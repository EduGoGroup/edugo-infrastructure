// Migration: Create material_assessment collection
// Collection: material_assessment (Owner: infrastructure)
// Created by: edugo-infrastructure
// Used by: api-mobile, worker
//
// Purpose: Stores AI-generated assessments/quizzes for educational materials
// Related PostgreSQL table: assessment (stores metadata, references this via mongo_document_id)

db.createCollection("material_assessment", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["material_id", "questions", "metadata", "created_at", "updated_at"],
      properties: {
        material_id: {
          bsonType: "string",
          description: "UUID of the material in PostgreSQL (materials table)"
        },
        questions: {
          bsonType: "array",
          description: "Array of questions for the assessment",
          items: {
            bsonType: "object",
            required: ["question_index", "question_text", "question_type", "options"],
            properties: {
              question_index: {
                bsonType: "int",
                minimum: 0,
                description: "Index/position of the question (0-based)"
              },
              question_text: {
                bsonType: "string",
                description: "The question text"
              },
              question_type: {
                bsonType: "string",
                enum: ["multiple_choice", "true_false", "short_answer"],
                description: "Type of question"
              },
              options: {
                bsonType: "array",
                description: "Array of answer options",
                items: {
                  bsonType: "object",
                  required: ["option_index", "text", "is_correct"],
                  properties: {
                    option_index: {
                      bsonType: "int",
                      minimum: 0,
                      description: "Index of the option (0-based)"
                    },
                    text: {
                      bsonType: "string",
                      description: "Option text"
                    },
                    is_correct: {
                      bsonType: "bool",
                      description: "Whether this option is correct"
                    }
                  }
                }
              },
              explanation: {
                bsonType: "string",
                description: "Optional explanation of the correct answer"
              }
            }
          }
        },
        metadata: {
          bsonType: "object",
          description: "Additional metadata about the assessment",
          properties: {
            subject: {
              bsonType: "string",
              description: "Subject/topic of the assessment"
            },
            grade: {
              bsonType: "string",
              description: "Grade level"
            },
            difficulty: {
              bsonType: "string",
              enum: ["easy", "medium", "hard"],
              description: "Difficulty level"
            },
            estimated_time_minutes: {
              bsonType: "int",
              minimum: 1,
              description: "Estimated time to complete in minutes"
            }
          }
        },
        created_at: {
          bsonType: "date",
          description: "Timestamp when the assessment was created"
        },
        updated_at: {
          bsonType: "date",
          description: "Timestamp when the assessment was last updated"
        }
      }
    }
  }
});

// Create indexes for efficient queries
db.material_assessment.createIndex({ "material_id": 1 }, { name: "idx_material_id" });
db.material_assessment.createIndex({ "metadata.subject": 1 }, { name: "idx_metadata_subject" });
db.material_assessment.createIndex({ "metadata.grade": 1 }, { name: "idx_metadata_grade" });
db.material_assessment.createIndex({ "metadata.difficulty": 1 }, { name: "idx_metadata_difficulty" });
db.material_assessment.createIndex({ "created_at": -1 }, { name: "idx_created_at_desc" });

print("✅ Collection 'material_assessment' created successfully");
