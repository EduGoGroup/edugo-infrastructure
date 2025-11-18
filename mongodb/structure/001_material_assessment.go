package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialAssessment creates the material_assessment collection with schema validation
// Collection: material_assessment (Owner: infrastructure)
// Used by: api-mobile, worker
// Purpose: Stores AI-generated assessments/quizzes for educational materials
func CreateMaterialAssessment(ctx context.Context, db *mongo.Database) error {
	collectionName := "material_assessment"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"material_id", "questions", "metadata", "created_at", "updated_at"},
			"properties": bson.M{
				"material_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the material in PostgreSQL",
				},
				"questions": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType": "object",
						"required": []string{"question_index", "question_text", "question_type", "options"},
					},
				},
				"metadata": bson.M{
					"bsonType": "object",
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
