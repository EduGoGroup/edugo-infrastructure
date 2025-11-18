package constraints

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAssessmentAttemptResultIndexes creates all indexes for assessment_attempt_result collection
func CreateAssessmentAttemptResultIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("assessment_attempt_result")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "attempt_id", Value: 1}},
			Options: options.Index().SetName("idx_attempt_id").SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "assessment_id", Value: 1}},
			Options: options.Index().SetName("idx_assessment_id"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "assessment_id", Value: 1},
			},
			Options: options.Index().SetName("idx_user_assessment"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
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
				{Key: "user_id", Value: 1},
				{Key: "completed_at", Value: -1},
			},
			Options: options.Index().SetName("idx_user_completed"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
