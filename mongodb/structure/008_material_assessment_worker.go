package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialAssessmentWorker creates the material_assessment_worker collection with schema validation
// Collection: material_assessment_worker (Owner: infrastructure)
// Used by: worker
// Purpose: Stores AI-generated assessments processed by the worker service
func CreateMaterialAssessmentWorker(ctx context.Context, db *mongo.Database) error {
	collectionName := "material_assessment_worker"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"material_id", "questions", "total_questions", "total_points", "version", "ai_model", "processing_time_ms", "created_at", "updated_at"},
			"properties": bson.M{
				"material_id": bson.M{
					"bsonType":    "string",
					"pattern":     "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
					"description": "UUID v4 of the material in PostgreSQL",
				},
				"questions": bson.M{
					"bsonType": "array",
					"minItems": 3,
					"maxItems": 20,
					"items": bson.M{
						"bsonType": "object",
						"required": []string{"question_id", "question_text", "question_type", "correct_answer", "points", "difficulty"},
						"properties": bson.M{
							"question_id": bson.M{
								"bsonType": "string",
							},
							"question_text": bson.M{
								"bsonType": "string",
							},
							"question_type": bson.M{
								"bsonType": "string",
								"enum":     []string{"multiple_choice", "true_false", "open"},
							},
							"options": bson.M{
								"bsonType": "array",
							},
							"correct_answer": bson.M{
								"bsonType": "string",
							},
							"points": bson.M{
								"bsonType": "int",
							},
							"difficulty": bson.M{
								"bsonType": "string",
								"enum":     []string{"easy", "medium", "hard"},
							},
							"explanation": bson.M{
								"bsonType": "string",
							},
						},
					},
					"description": "Array of assessment questions",
				},
				"total_questions": bson.M{
					"bsonType":    "int",
					"minimum":     3,
					"maximum":     20,
					"description": "Total number of questions in the assessment",
				},
				"total_points": bson.M{
					"bsonType":    "int",
					"description": "Total points possible in the assessment",
				},
				"version": bson.M{
					"bsonType":    "int",
					"minimum":     1,
					"description": "Version number of the assessment",
				},
				"ai_model": bson.M{
					"bsonType":    "string",
					"enum":        []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo", "gpt-4o"},
					"description": "AI model used to generate the assessment",
				},
				"processing_time_ms": bson.M{
					"bsonType":    "int",
					"minimum":     0,
					"description": "Processing time in milliseconds",
				},
				"metadata": bson.M{
					"bsonType":    "object",
					"description": "Additional metadata about the assessment generation",
				},
				"created_at": bson.M{
					"bsonType": "date",
				},
				"updated_at": bson.M{
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
