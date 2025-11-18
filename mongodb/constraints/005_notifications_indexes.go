package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNotificationsIndexes creates all indexes for notifications collection
func CreateNotificationsIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("notifications")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "type", Value: 1}},
			Options: options.Index().SetName("idx_type"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("idx_user_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("idx_user_created"),
		},
		{
			Keys:    bson.D{{Key: "read_at", Value: -1}},
			Options: options.Index().SetName("idx_read_at_desc"),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetName("idx_ttl_expires").SetExpireAfterSeconds(0),
		},
		{
			Keys:    bson.D{{Key: "archived_at", Value: 1}},
			Options: options.Index().SetName("idx_ttl_archived_30days").SetExpireAfterSeconds(2592000),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
