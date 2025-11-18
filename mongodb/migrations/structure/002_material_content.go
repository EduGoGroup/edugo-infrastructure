package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialContent creates the material_content collection with schema validation
// Collection: material_content (Owner: infrastructure)
// Used by: worker
// Purpose: Stores extracted/processed content from educational materials
func CreateMaterialContent(ctx context.Context, db *mongo.Database) error {
	collectionName := "material_content"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"material_id", "content_type", "created_at", "updated_at"},
			"properties": bson.M{
				"material_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the material in PostgreSQL",
				},
				"content_type": bson.M{
					"bsonType":    "string",
					"enum":        []string{"pdf_extracted", "video_transcript", "document_parsed", "slides_extracted"},
					"description": "Type of content extracted",
				},
				"raw_text": bson.M{
					"bsonType":    "string",
					"description": "Raw extracted text content",
				},
				"structured_content": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"title": bson.M{
							"bsonType": "string",
						},
						"sections": bson.M{
							"bsonType": "array",
						},
						"summary": bson.M{
							"bsonType": "string",
						},
						"key_concepts": bson.M{
							"bsonType": "array",
						},
					},
					"description": "Structured parsed content",
				},
				"processing_info": bson.M{
					"bsonType":    "object",
					"description": "Information about the processing",
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
