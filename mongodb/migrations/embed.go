package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ApplyAll ejecuta structure + constraints (base de datos limpia lista para usar)
func ApplyAll(ctx context.Context, db *mongo.Database) error {
	if err := ApplyStructure(ctx, db); err != nil {
		return fmt.Errorf("error aplicando structure: %w", err)
	}
	if err := ApplyConstraints(ctx, db); err != nil {
		return fmt.Errorf("error aplicando constraints: %w", err)
	}
	return nil
}

// ApplyStructure ejecuta todas las funciones de creación de collections con validators
func ApplyStructure(ctx context.Context, db *mongo.Database) error {
	structureFuncs := []func(context.Context, *mongo.Database) error{
		createMaterialSummary,
		createMaterialAssessmentWorker,
		createMaterialEvent,
	}

	for _, fn := range structureFuncs {
		if err := fn(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

// ApplyConstraints ejecuta todas las funciones de creación de índices
func ApplyConstraints(ctx context.Context, db *mongo.Database) error {
	constraintsFuncs := []func(context.Context, *mongo.Database) error{
		createMaterialSummaryIndexes,
		createMaterialAssessmentWorkerIndexes,
		createMaterialEventIndexes,
	}

	for _, fn := range constraintsFuncs {
		if err := fn(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

// ApplySeeds ejecuta seeds (datos iniciales del ecosistema)
func ApplySeeds(ctx context.Context, db *mongo.Database) error {
	inserted, err := applySeedsInternal(ctx, db)
	if err != nil {
		return fmt.Errorf("error aplicando seeds: %w", err)
	}
	if inserted > 0 {
		_ = inserted
	}
	return nil
}

// ApplyMockData ejecuta datos mock para testing y desarrollo
func ApplyMockData(ctx context.Context, db *mongo.Database) error {
	inserted, err := applyMockDataInternal(ctx, db)
	if err != nil {
		return fmt.Errorf("error aplicando mock data: %w", err)
	}
	if inserted > 0 {
		_ = inserted
	}
	return nil
}

// ListFunctions lista todas las funciones disponibles por capa
func ListFunctions() map[string][]string {
	return map[string][]string{
		"structure": {
			"createMaterialSummary",
			"createMaterialAssessmentWorker",
			"createMaterialEvent",
		},
		"constraints": {
			"createMaterialSummaryIndexes",
			"createMaterialAssessmentWorkerIndexes",
			"createMaterialEventIndexes",
		},
		"seeds":   {},
		"testing": {},
	}
}

// ============================================================
// Structure functions (inlined from deprecated sub-packages)
// ============================================================

func createMaterialSummary(ctx context.Context, db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"material_id", "summary", "key_points", "language", "word_count", "version", "ai_model", "processing_time_ms", "created_at", "updated_at"},
			"properties": bson.M{
				"material_id":      bson.M{"bsonType": "string", "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"},
				"summary":          bson.M{"bsonType": "string", "minLength": 10, "maxLength": 5000},
				"key_points":       bson.M{"bsonType": "array", "minItems": 1, "maxItems": 10, "items": bson.M{"bsonType": "string"}},
				"language":         bson.M{"bsonType": "string", "enum": []string{"es", "en", "pt"}},
				"word_count":       bson.M{"bsonType": "int", "minimum": 1},
				"version":          bson.M{"bsonType": "int", "minimum": 1},
				"ai_model":         bson.M{"bsonType": "string", "enum": []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo", "gpt-4o"}},
				"processing_time_ms": bson.M{"bsonType": "int", "minimum": 0},
				"metadata":         bson.M{"bsonType": "object"},
				"created_at":       bson.M{"bsonType": "date"},
				"updated_at":       bson.M{"bsonType": "date"},
			},
		},
	}
	opts := options.CreateCollection().SetValidator(validator)
	err := db.CreateCollection(ctx, "material_summary", opts)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

func createMaterialAssessmentWorker(ctx context.Context, db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"material_id", "questions", "total_questions", "total_points", "version", "ai_model", "processing_time_ms", "created_at", "updated_at"},
			"properties": bson.M{
				"material_id": bson.M{"bsonType": "string", "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"},
				"questions": bson.M{
					"bsonType": "array", "minItems": 3, "maxItems": 20,
					"items": bson.M{
						"bsonType": "object",
						"required": []string{"question_id", "question_text", "question_type", "correct_answer", "points", "difficulty"},
						"properties": bson.M{
							"question_id":   bson.M{"bsonType": "string"},
							"question_text":  bson.M{"bsonType": "string"},
							"question_type":  bson.M{"bsonType": "string", "enum": []string{"multiple_choice", "true_false", "open"}},
							"options":        bson.M{"bsonType": "array"},
							"correct_answer": bson.M{"bsonType": "string"},
							"points":         bson.M{"bsonType": "int"},
							"difficulty":     bson.M{"bsonType": "string", "enum": []string{"easy", "medium", "hard"}},
							"explanation":    bson.M{"bsonType": "string"},
						},
					},
				},
				"total_questions":    bson.M{"bsonType": "int", "minimum": 3, "maximum": 20},
				"total_points":       bson.M{"bsonType": "int"},
				"version":            bson.M{"bsonType": "int", "minimum": 1},
				"ai_model":           bson.M{"bsonType": "string", "enum": []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo", "gpt-4o"}},
				"processing_time_ms": bson.M{"bsonType": "int", "minimum": 0},
				"metadata":           bson.M{"bsonType": "object"},
				"created_at":         bson.M{"bsonType": "date"},
				"updated_at":         bson.M{"bsonType": "date"},
			},
		},
	}
	opts := options.CreateCollection().SetValidator(validator)
	err := db.CreateCollection(ctx, "material_assessment_worker", opts)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

func createMaterialEvent(ctx context.Context, db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"event_type", "payload", "status", "retry_count", "created_at", "updated_at"},
			"properties": bson.M{
				"event_type":  bson.M{"bsonType": "string", "enum": []string{"material_uploaded", "material_reprocess", "material_deleted", "assessment_attempt", "student_enrolled", "student_unenrolled"}},
				"material_id": bson.M{"bsonType": "string", "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"},
				"user_id":     bson.M{"bsonType": "string", "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"},
				"payload":     bson.M{"bsonType": "object"},
				"status":      bson.M{"bsonType": "string", "enum": []string{"pending", "processing", "completed", "failed"}},
				"error_msg":   bson.M{"bsonType": "string", "maxLength": 5000},
				"stack_trace": bson.M{"bsonType": "string", "maxLength": 10000},
				"retry_count":  bson.M{"bsonType": "int", "minimum": 0},
				"next_retry_at": bson.M{"bsonType": "date"},
				"processed_at":  bson.M{"bsonType": "date"},
				"created_at":    bson.M{"bsonType": "date"},
				"updated_at":    bson.M{"bsonType": "date"},
			},
		},
	}
	opts := options.CreateCollection().SetValidator(validator)
	err := db.CreateCollection(ctx, "material_event", opts)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

// ============================================================
// Constraint functions (inlined from deprecated sub-packages)
// ============================================================

func createMaterialSummaryIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "material_id", Value: 1}}, Options: options.Index().SetName("idx_material_id").SetUnique(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}, Options: options.Index().SetName("idx_status")},
		{Keys: bson.D{{Key: "metadata.subject", Value: 1}}, Options: options.Index().SetName("idx_metadata_subject")},
		{Keys: bson.D{{Key: "metadata.grade", Value: 1}}, Options: options.Index().SetName("idx_metadata_grade")},
		{Keys: bson.D{{Key: "metadata.difficulty", Value: 1}}, Options: options.Index().SetName("idx_metadata_difficulty")},
		{Keys: bson.D{{Key: "created_at", Value: -1}}, Options: options.Index().SetName("idx_created_at_desc")},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}, Options: options.Index().SetName("idx_updated_at_desc")},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "updated_at", Value: -1}}, Options: options.Index().SetName("idx_status_updated")},
	}
	_, err := db.Collection("material_summary").Indexes().CreateMany(ctx, indexes)
	return err
}

func createMaterialAssessmentWorkerIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "material_id", Value: 1}}, Options: options.Index().SetName("idx_material_id")},
		{Keys: bson.D{{Key: "status", Value: 1}}, Options: options.Index().SetName("idx_status")},
		{Keys: bson.D{{Key: "worker_id", Value: 1}}, Options: options.Index().SetName("idx_worker_id")},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: 1}}, Options: options.Index().SetName("idx_status_created")},
		{Keys: bson.D{{Key: "created_at", Value: -1}}, Options: options.Index().SetName("idx_created_at_desc")},
		{Keys: bson.D{{Key: "started_at", Value: -1}}, Options: options.Index().SetName("idx_started_at_desc")},
		{Keys: bson.D{{Key: "completed_at", Value: -1}}, Options: options.Index().SetName("idx_completed_at_desc")},
		{Keys: bson.D{{Key: "worker_id", Value: 1}, {Key: "status", Value: 1}}, Options: options.Index().SetName("idx_worker_status")},
	}
	_, err := db.Collection("material_assessment_worker").Indexes().CreateMany(ctx, indexes)
	return err
}

func createMaterialEventIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "material_id", Value: 1}}, Options: options.Index().SetName("idx_material_id")},
		{Keys: bson.D{{Key: "event_type", Value: 1}}, Options: options.Index().SetName("idx_event_type")},
		{Keys: bson.D{{Key: "status", Value: 1}}, Options: options.Index().SetName("idx_status")},
		{Keys: bson.D{{Key: "material_id", Value: 1}, {Key: "event_type", Value: 1}}, Options: options.Index().SetName("idx_material_event")},
		{Keys: bson.D{{Key: "timestamp", Value: -1}}, Options: options.Index().SetName("idx_timestamp_desc")},
		{Keys: bson.D{{Key: "material_id", Value: 1}, {Key: "timestamp", Value: -1}}, Options: options.Index().SetName("idx_material_timestamp")},
		{Keys: bson.D{{Key: "event_type", Value: 1}, {Key: "timestamp", Value: -1}}, Options: options.Index().SetName("idx_event_timestamp")},
		{Keys: bson.D{{Key: "created_at", Value: 1}}, Options: options.Index().SetName("idx_ttl_90days").SetExpireAfterSeconds(7776000)},
	}
	_, err := db.Collection("material_event").Indexes().CreateMany(ctx, indexes)
	return err
}
