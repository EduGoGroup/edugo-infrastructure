package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNotifications creates the notifications collection with schema validation
// Collection: notifications (Owner: infrastructure)
// Used by: api-mobile
// Purpose: Stores user notifications for various events
func CreateNotifications(ctx context.Context, db *mongo.Database) error {
	collectionName := "notifications"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"user_id", "notification_type", "title", "is_read", "created_at"},
			"properties": bson.M{
				"user_id": bson.M{
					"bsonType":    "string",
					"description": "UUID of the user in PostgreSQL",
				},
				"notification_type": bson.M{
					"bsonType": "string",
					"enum": []string{
						"assessment.ready", "assessment.graded",
						"material.uploaded", "material.processed", "material.shared",
						"membership.added", "membership.removed",
						"deadline.approaching",
						"achievement.unlocked",
						"system.announcement", "system.maintenance",
					},
					"description": "Type of notification",
				},
				"title": bson.M{
					"bsonType":    "string",
					"description": "Notification title",
				},
				"message": bson.M{
					"bsonType":    "string",
					"description": "Notification message body",
				},
				"is_read": bson.M{
					"bsonType":    "bool",
					"description": "Whether the notification has been read",
				},
				"priority": bson.M{
					"bsonType":    "string",
					"enum":        []string{"low", "medium", "high", "urgent"},
					"description": "Priority level of the notification",
				},
				"category": bson.M{
					"bsonType":    "string",
					"enum":        []string{"academic", "administrative", "social", "system"},
					"description": "Category of the notification",
				},
				"action_url": bson.M{
					"bsonType":    "string",
					"description": "URL to navigate when clicking the notification",
				},
				"metadata": bson.M{
					"bsonType":    "object",
					"description": "Additional notification data",
				},
				"read_at": bson.M{
					"bsonType":    "date",
					"description": "Timestamp when notification was read",
				},
				"created_at": bson.M{
					"bsonType": "date",
				},
				"expires_at": bson.M{
					"bsonType":    "date",
					"description": "Expiration timestamp for the notification",
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
