package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMaterialAssessmentWorkerIndexes creates all indexes for material_assessment_worker collection
func CreateMaterialAssessmentWorkerIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("material_assessment_worker")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "material_id", Value: 1}},
			Options: options.Index().SetName("idx_material_id"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "worker_id", Value: 1}},
			Options: options.Index().SetName("idx_worker_id"),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: 1},
			},
			Options: options.Index().SetName("idx_status_created"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		{
			Keys:    bson.D{{Key: "started_at", Value: -1}},
			Options: options.Index().SetName("idx_started_at_desc"),
		},
		{
			Keys:    bson.D{{Key: "completed_at", Value: -1}},
			Options: options.Index().SetName("idx_completed_at_desc"),
		},
		{
			Keys: bson.D{
				{Key: "worker_id", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("idx_worker_status"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
