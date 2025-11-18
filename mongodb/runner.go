package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/constraints"
	"github.com/EduGoGroup/edugo-infrastructure/mongodb/structure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	mongoURI := getMongoURI()
	dbName := getEnv("MONGO_DB_NAME", "edugo")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(dbName)

	switch command {
	case "structure":
		if err := runStructure(ctx, db); err != nil {
			log.Fatalf("Error ejecutando structure: %v", err)
		}
	case "constraints":
		if err := runConstraints(ctx, db); err != nil {
			log.Fatalf("Error ejecutando constraints: %v", err)
		}
	case "all":
		if err := runStructure(ctx, db); err != nil {
			log.Fatalf("Error ejecutando structure: %v", err)
		}
		if err := runConstraints(ctx, db); err != nil {
			log.Fatalf("Error ejecutando constraints: %v", err)
		}
		fmt.Println("✅ Todas las migraciones completadas")
	default:
		printHelp()
		os.Exit(1)
	}
}

func runStructure(ctx context.Context, db *mongo.Database) error {
	fmt.Println("🏗️  Ejecutando Structure...")

	if err := structure.CreateMaterialAssessment(ctx, db); err != nil {
		return fmt.Errorf("001_material_assessment: %w", err)
	}
	fmt.Println("✅ 001_material_assessment")

	if err := structure.CreateMaterialContent(ctx, db); err != nil {
		return fmt.Errorf("002_material_content: %w", err)
	}
	fmt.Println("✅ 002_material_content")

	if err := structure.CreateAssessmentAttemptResult(ctx, db); err != nil {
		return fmt.Errorf("003_assessment_attempt_result: %w", err)
	}
	fmt.Println("✅ 003_assessment_attempt_result")

	if err := structure.CreateAuditLogs(ctx, db); err != nil {
		return fmt.Errorf("004_audit_logs: %w", err)
	}
	fmt.Println("✅ 004_audit_logs")

	if err := structure.CreateNotifications(ctx, db); err != nil {
		return fmt.Errorf("005_notifications: %w", err)
	}
	fmt.Println("✅ 005_notifications")

	if err := structure.CreateAnalyticsEvents(ctx, db); err != nil {
		return fmt.Errorf("006_analytics_events: %w", err)
	}
	fmt.Println("✅ 006_analytics_events")

	if err := structure.CreateMaterialSummary(ctx, db); err != nil {
		return fmt.Errorf("007_material_summary: %w", err)
	}
	fmt.Println("✅ 007_material_summary")

	if err := structure.CreateMaterialAssessmentWorker(ctx, db); err != nil {
		return fmt.Errorf("008_material_assessment_worker: %w", err)
	}
	fmt.Println("✅ 008_material_assessment_worker")

	if err := structure.CreateMaterialEvent(ctx, db); err != nil {
		return fmt.Errorf("009_material_event: %w", err)
	}
	fmt.Println("✅ 009_material_event")

	return nil
}

func runConstraints(ctx context.Context, db *mongo.Database) error {
	fmt.Println("🔗 Ejecutando Constraints...")

	if err := constraints.CreateMaterialAssessmentIndexes(ctx, db); err != nil {
		return fmt.Errorf("001_material_assessment_indexes: %w", err)
	}
	fmt.Println("✅ 001_material_assessment_indexes")

	if err := constraints.CreateMaterialContentIndexes(ctx, db); err != nil {
		return fmt.Errorf("002_material_content_indexes: %w", err)
	}
	fmt.Println("✅ 002_material_content_indexes")

	if err := constraints.CreateAssessmentAttemptResultIndexes(ctx, db); err != nil {
		return fmt.Errorf("003_assessment_attempt_result_indexes: %w", err)
	}
	fmt.Println("✅ 003_assessment_attempt_result_indexes")

	if err := constraints.CreateAuditLogsIndexes(ctx, db); err != nil {
		return fmt.Errorf("004_audit_logs_indexes: %w", err)
	}
	fmt.Println("✅ 004_audit_logs_indexes")

	if err := constraints.CreateNotificationsIndexes(ctx, db); err != nil {
		return fmt.Errorf("005_notifications_indexes: %w", err)
	}
	fmt.Println("✅ 005_notifications_indexes")

	if err := constraints.CreateAnalyticsEventsIndexes(ctx, db); err != nil {
		return fmt.Errorf("006_analytics_events_indexes: %w", err)
	}
	fmt.Println("✅ 006_analytics_events_indexes")

	if err := constraints.CreateMaterialSummaryIndexes(ctx, db); err != nil {
		return fmt.Errorf("007_material_summary_indexes: %w", err)
	}
	fmt.Println("✅ 007_material_summary_indexes")

	if err := constraints.CreateMaterialAssessmentWorkerIndexes(ctx, db); err != nil {
		return fmt.Errorf("008_material_assessment_worker_indexes: %w", err)
	}
	fmt.Println("✅ 008_material_assessment_worker_indexes")

	if err := constraints.CreateMaterialEventIndexes(ctx, db); err != nil {
		return fmt.Errorf("009_material_event_indexes: %w", err)
	}
	fmt.Println("✅ 009_material_event_indexes")

	return nil
}

func printHelp() {
	fmt.Println("MongoDB Migration Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run runner.go structure    Ejecutar solo structure")
	fmt.Println("  go run runner.go constraints  Ejecutar solo constraints")
	fmt.Println("  go run runner.go all          Ejecutar todo")
}

func getMongoURI() string {
	host := getEnv("MONGO_HOST", "localhost")
	port := getEnv("MONGO_PORT", "27017")
	user := getEnv("MONGO_USER", "edugo")
	password := getEnv("MONGO_PASSWORD", "edugo123")
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
