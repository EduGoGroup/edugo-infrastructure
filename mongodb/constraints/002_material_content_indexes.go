package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialContentIndexes creates all indexes for material_content collection
func CreateMaterialContentIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("material_content")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "material_id", Value: 1}},
			Options: options.Index().SetName("idx_material_id").SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "content_type", Value: 1}},
			Options: options.Index().SetName("idx_content_type"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		{
			Keys:    bson.D{{Key: "processing_info.processed_at", Value: -1}},
			Options: options.Index().SetName("idx_processed_at_desc"),
		},
		{
			Keys: bson.D{
				{Key: "raw_text", Value: "text"},
				{Key: "structured_content.summary", Value: "text"},
				{Key: "structured_content.key_concepts", Value: "text"},
			},
			Options: options.Index().SetName("idx_fulltext_search").SetDefaultLanguage("spanish"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
