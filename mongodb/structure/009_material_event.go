package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialEvent creates the material_event collection with schema validation
// Collection: material_event (Owner: infrastructure)
// Used by: worker, api-mobile
// Purpose: Stores event queue for material-related processing tasks
func CreateMaterialEvent(ctx context.Context, db *mongo.Database) error {
	collectionName := "material_event"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"event_type", "payload", "status", "retry_count", "created_at", "updated_at"},
			"properties": bson.M{
				"event_type": bson.M{
					"bsonType": "string",
					"enum": []string{
						"material_uploaded",
						"material_reprocess",
						"material_deleted",
						"assessment_attempt",
						"student_enrolled",
						"student_unenrolled",
					},
					"description": "Type of event",
				},
				"material_id": bson.M{
					"bsonType":    "string",
					"pattern":     "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
					"description": "UUID v4 of the material in PostgreSQL",
				},
				"user_id": bson.M{
					"bsonType":    "string",
					"pattern":     "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$",
					"description": "UUID v4 of the user in PostgreSQL",
				},
				"payload": bson.M{
					"bsonType":    "object",
					"description": "Event payload with flexible structure",
				},
				"status": bson.M{
					"bsonType":    "string",
					"enum":        []string{"pending", "processing", "completed", "failed"},
					"description": "Current status of the event processing",
				},
				"error_msg": bson.M{
					"bsonType":    "string",
					"maxLength":   5000,
					"description": "Error message if the event processing failed",
				},
				"stack_trace": bson.M{
					"bsonType":    "string",
					"maxLength":   10000,
					"description": "Stack trace if the event processing failed",
				},
				"retry_count": bson.M{
					"bsonType":    "int",
					"minimum":     0,
					"description": "Number of times this event has been retried",
				},
				"next_retry_at": bson.M{
					"bsonType":    "date",
					"description": "Timestamp for the next retry attempt",
				},
				"processed_at": bson.M{
					"bsonType":    "date",
					"description": "Timestamp when the event was successfully processed",
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
