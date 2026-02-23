package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CreateMaterialSummaryIndexes creates all indexes for material_summary collection
func CreateMaterialSummaryIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("material_summary")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "material_id", Value: 1}},
			Options: options.Index().SetName("idx_material_id").SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "metadata.subject", Value: 1}},
			Options: options.Index().SetName("idx_metadata_subject"),
		},
		{
			Keys:    bson.D{{Key: "metadata.grade", Value: 1}},
			Options: options.Index().SetName("idx_metadata_grade"),
		},
		{
			Keys:    bson.D{{Key: "metadata.difficulty", Value: 1}},
			Options: options.Index().SetName("idx_metadata_difficulty"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		{
			Keys:    bson.D{{Key: "updated_at", Value: -1}},
			Options: options.Index().SetName("idx_updated_at_desc"),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "updated_at", Value: -1},
			},
			Options: options.Index().SetName("idx_status_updated"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
