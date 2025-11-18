package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialEventIndexes creates all indexes for material_event collection
func CreateMaterialEventIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("material_event")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "material_id", Value: 1}},
			Options: options.Index().SetName("idx_material_id"),
		},
		{
			Keys:    bson.D{{Key: "event_type", Value: 1}},
			Options: options.Index().SetName("idx_event_type"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys: bson.D{
				{Key: "material_id", Value: 1},
				{Key: "event_type", Value: 1},
			},
			Options: options.Index().SetName("idx_material_event"),
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: -1}},
			Options: options.Index().SetName("idx_timestamp_desc"),
		},
		{
			Keys: bson.D{
				{Key: "material_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().SetName("idx_material_timestamp"),
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
			Options: options.Index().SetName("idx_ttl_90days").SetExpireAfterSeconds(7776000),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
