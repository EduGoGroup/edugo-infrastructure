package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialAssessmentIndexes creates all indexes for material_assessment collection
// Indexes optimize queries by material_id, metadata fields, and created_at
func CreateMaterialAssessmentIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("material_assessment")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "material_id", Value: 1}},
			Options: options.Index().SetName("idx_material_id"),
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
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
