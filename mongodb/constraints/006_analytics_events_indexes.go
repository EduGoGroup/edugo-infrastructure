package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAnalyticsEventsIndexes creates all indexes for analytics_events collection
func CreateAnalyticsEventsIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("analytics_events")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "event_type", Value: 1}},
			Options: options.Index().SetName("idx_event_type"),
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
			Keys: bson.D{
				{Key: "event_type", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().SetName("idx_event_timestamp"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("idx_ttl_365days").SetExpireAfterSeconds(31536000),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
