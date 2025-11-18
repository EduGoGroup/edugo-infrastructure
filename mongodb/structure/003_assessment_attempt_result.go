package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAssessmentAttemptResult creates the assessment_attempt_result collection with schema validation
// Collection: assessment_attempt_result (Owner: infrastructure)
// Used by: api-mobile
// Purpose: Stores detailed results and answers from assessment attempts
func CreateAssessmentAttemptResult(ctx context.Context, db *mongo.Database) error {
	collectionName := "assessment_attempt_result"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"attempt_id", "student_id", "assessment_id", "answers", "score", "started_at", "submitted_at", "created_at"},
			"properties": bson.M{
				"attempt_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the attempt in PostgreSQL",
				},
				"student_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the student in PostgreSQL",
				},
				"assessment_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the assessment in PostgreSQL",
				},
				"answers": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType": "object",
						"required": []string{"question_index", "selected_option_index", "is_correct", "time_spent_seconds"},
						"properties": bson.M{
							"question_index": bson.M{
								"bsonType": "int",
							},
							"selected_option_index": bson.M{
								"bsonType": "int",
							},
							"is_correct": bson.M{
								"bsonType": "bool",
							},
							"time_spent_seconds": bson.M{
								"bsonType": "int",
							},
						},
					},
					"description": "Array of individual question answers",
				},
				"score": bson.M{
					"bsonType": "object",
					"required": []string{"correct_count", "total_questions", "percentage"},
					"properties": bson.M{
						"correct_count": bson.M{
							"bsonType": "int",
						},
						"total_questions": bson.M{
							"bsonType": "int",
						},
						"percentage": bson.M{
							"bsonType": "double",
						},
					},
					"description": "Score summary",
				},
				"started_at": bson.M{
					"bsonType": "date",
				},
				"submitted_at": bson.M{
					"bsonType": "date",
				},
				"created_at": bson.M{
					"bsonType": "date",
				},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	err := db.CreateCollection(ctx, collectionName, opts)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	return nil
}
