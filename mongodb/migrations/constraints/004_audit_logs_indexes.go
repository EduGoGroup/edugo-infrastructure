package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAuditLogsIndexes creates all indexes for audit_logs collection
func CreateAuditLogsIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("audit_logs")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "action", Value: 1}},
			Options: options.Index().SetName("idx_action"),
		},
		{
			Keys:    bson.D{{Key: "entity_type", Value: 1}},
			Options: options.Index().SetName("idx_entity_type"),
		},
		{
			Keys:    bson.D{{Key: "entity_id", Value: 1}},
			Options: options.Index().SetName("idx_entity_id"),
		},
		{
			Keys: bson.D{
				{Key: "entity_type", Value: 1},
				{Key: "entity_id", Value: 1},
			},
			Options: options.Index().SetName("idx_entity"),
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: -1}},
			Options: options.Index().SetName("idx_timestamp_desc"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().SetName("idx_user_timestamp"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("idx_ttl_90days").SetExpireAfterSeconds(7776000),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
