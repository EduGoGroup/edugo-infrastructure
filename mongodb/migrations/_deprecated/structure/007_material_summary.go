package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CreateMaterialSummary creates the material_summary collection with schema validation
// Collection: material_summary (Owner: infrastructure)
// Used by: worker, api-mobile
// Purpose: Stores AI-generated summaries of educational materials
func CreateMaterialSummary(ctx context.Context, db *mongo.Database) error {
	collectionName := "material_summary"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"material_id", "summary", "key_points", "language", "word_count", "version", "ai_model", "processing_time_ms", "created_at", "updated_at"},
			"properties": bson.M{
				"material_id": bson.M{
					"bsonType":    "string",
					"pattern":     "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
					"description": "UUID v4 of the material in PostgreSQL",
				},
				"summary": bson.M{
					"bsonType":    "string",
					"minLength":   10,
					"maxLength":   5000,
					"description": "AI-generated summary text",
				},
				"key_points": bson.M{
					"bsonType": "array",
					"minItems": 1,
					"maxItems": 10,
					"items": bson.M{
						"bsonType": "string",
					},
					"description": "Array of key points extracted from the material",
				},
				"language": bson.M{
					"bsonType":    "string",
					"enum":        []string{"es", "en", "pt"},
					"description": "Language of the summary",
				},
				"word_count": bson.M{
					"bsonType":    "int",
					"minimum":     1,
					"description": "Word count of the summary",
				},
				"version": bson.M{
					"bsonType":    "int",
					"minimum":     1,
					"description": "Version number of the summary",
				},
				"ai_model": bson.M{
					"bsonType":    "string",
					"enum":        []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo", "gpt-4o"},
					"description": "AI model used to generate the summary",
				},
				"processing_time_ms": bson.M{
					"bsonType":    "int",
					"minimum":     0,
					"description": "Processing time in milliseconds",
				},
				"metadata": bson.M{
					"bsonType":    "object",
					"description": "Additional metadata about the summary generation",
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
