package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAnalyticsEvents creates the analytics_events collection with schema validation
// Collection: analytics_events (Owner: infrastructure)
// Used by: api-mobile, worker
// Purpose: Stores analytics and tracking events for user behavior analysis
func CreateAnalyticsEvents(ctx context.Context, db *mongo.Database) error {
	collectionName := "analytics_events"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"event_name", "timestamp"},
			"properties": bson.M{
				"event_name": bson.M{
					"bsonType": "string",
					"enum": []string{
						"page.view",
						"material.view", "material.download", "material.search",
						"assessment.start", "assessment.complete", "assessment.abandon",
						"question.answer", "question.skip",
						"video.play", "video.pause", "video.complete",
						"session.start", "session.end",
						"feature.click",
						"error.occurred",
						"search.performed",
						"filter.applied",
					},
					"description": "Name of the analytics event",
				},
				"user_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the user in PostgreSQL",
				},
				"session_id": bson.M{
					"bsonType":    "string",
					"description": "Session identifier",
				},
				"timestamp": bson.M{
					"bsonType": "date",
				},
				"properties": bson.M{
					"bsonType":    "object",
					"description": "Event-specific properties",
				},
				"device": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"type": bson.M{
							"bsonType": "string",
						},
						"os": bson.M{
							"bsonType": "string",
						},
						"browser": bson.M{
							"bsonType": "string",
						},
					},
					"description": "Device information",
				},
				"location": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"country": bson.M{
							"bsonType": "string",
						},
						"city": bson.M{
							"bsonType": "string",
						},
						"timezone": bson.M{
							"bsonType": "string",
						},
					},
					"description": "Location information",
				},
				"context": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"page": bson.M{
							"bsonType": "string",
						},
						"referrer": bson.M{
							"bsonType": "string",
						},
						"url": bson.M{
							"bsonType": "string",
						},
					},
					"description": "Context information",
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
