package migrations_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestIntegration tests de integración con MongoDB en testcontainer
// Solo se ejecutan si ENABLE_INTEGRATION_TESTS=true
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set ENABLE_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()

	// Crear testcontainer MongoDB
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Error creando container: %v", err)
	}
	defer container.Terminate(ctx)

	// Obtener endpoint
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo host: %v", err)
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("Error obteniendo puerto: %v", err)
	}

	// Conectar a MongoDB
	uri := "mongodb://" + host + ":" + port.Port()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("testdb")

	t.Log("✅ Container MongoDB creado y conectado")

	// Ejecutar tests
	t.Run("ApplyAll", testApplyAll(ctx, db))
	t.Run("CRUD_MaterialAssessment", testCRUDMaterialAssessment(ctx, db))
	t.Run("CRUD_Notifications", testCRUDNotifications(ctx, db))
	t.Run("Indexes_Validation", testIndexesValidation(ctx, db))
}

func testApplyAll(ctx context.Context, db *mongo.Database) func(*testing.T) {
	return func(t *testing.T) {
		// Aplicar todas las migraciones
		if err := migrations.ApplyAll(ctx, db); err != nil {
			t.Fatalf("Error aplicando migraciones: %v", err)
		}

		// Verificar que las collections existen
		collections := []string{
			"material_assessment",
			"material_content",
			"assessment_attempt_result",
			"audit_logs",
			"notifications",
			"analytics_events",
			"material_summary",
			"material_assessment_worker",
			"material_event",
		}

		existingCollections, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			t.Fatalf("Error listando collections: %v", err)
		}

		collMap := make(map[string]bool)
		for _, c := range existingCollections {
			collMap[c] = true
		}

		for _, expected := range collections {
			if !collMap[expected] {
				t.Errorf("Collection %s no fue creada", expected)
			}
		}

		t.Logf("✅ Todas las %d collections creadas correctamente", len(collections))
	}
}

func testCRUDMaterialAssessment(ctx context.Context, db *mongo.Database) func(*testing.T) {
	return func(t *testing.T) {
		coll := db.Collection("material_assessment")

		// CREATE
		doc := bson.M{
			"material_id": "mat_test_123",
			"questions": bson.A{
				bson.M{
					"question_index": 0,
					"question_text":  "Test question?",
					"question_type":  "multiple_choice",
					"options": bson.A{
						bson.M{"option_index": 0, "text": "A", "is_correct": true},
						bson.M{"option_index": 1, "text": "B", "is_correct": false},
					},
				},
			},
			"metadata": bson.M{
				"subject":    "Math",
				"difficulty": "easy",
			},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		result, err := coll.InsertOne(ctx, doc)
		if err != nil {
			t.Fatalf("Error insertando: %v", err)
		}
		insertedID := result.InsertedID

		// READ
		var retrieved bson.M
		err = coll.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&retrieved)
		if err != nil {
			t.Fatalf("Error leyendo: %v", err)
		}
		if retrieved["material_id"] != "mat_test_123" {
			t.Errorf("Material ID incorrecto: %v", retrieved["material_id"])
		}

		// UPDATE
		_, err = coll.UpdateOne(ctx,
			bson.M{"_id": insertedID},
			bson.M{"$set": bson.M{"metadata.difficulty": "medium"}},
		)
		if err != nil {
			t.Fatalf("Error actualizando: %v", err)
		}

		// Verificar update
		err = coll.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&retrieved)
		if err != nil {
			t.Fatalf("Error verificando update: %v", err)
		}
		metadata := retrieved["metadata"].(bson.M)
		if metadata["difficulty"] != "medium" {
			t.Errorf("Update no aplicado correctamente")
		}

		// DELETE
		_, err = coll.DeleteOne(ctx, bson.M{"_id": insertedID})
		if err != nil {
			t.Fatalf("Error eliminando: %v", err)
		}

		// Verificar delete
		count, _ := coll.CountDocuments(ctx, bson.M{"_id": insertedID})
		if count != 0 {
			t.Error("Documento no fue eliminado")
		}

		t.Log("✅ CRUD material_assessment OK")
	}
}

func testCRUDNotifications(ctx context.Context, db *mongo.Database) func(*testing.T) {
	return func(t *testing.T) {
		coll := db.Collection("notifications")

		// CREATE
		doc := bson.M{
			"user_id":           "usr_test_1",
			"notification_type": "system.announcement",
			"title":             "Test Title",
			"message":           "Test Message",
			"is_read":           false,
			"priority":          "medium",
			"created_at":        time.Now(),
		}

		result, err := coll.InsertOne(ctx, doc)
		if err != nil {
			t.Fatalf("Error insertando notificación: %v", err)
		}

		// READ
		var retrieved bson.M
		err = coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&retrieved)
		if err != nil {
			t.Fatalf("Error leyendo: %v", err)
		}

		// UPDATE - marcar como leída
		_, err = coll.UpdateOne(ctx,
			bson.M{"_id": result.InsertedID},
			bson.M{"$set": bson.M{"is_read": true, "read_at": time.Now()}},
		)
		if err != nil {
			t.Fatalf("Error actualizando: %v", err)
		}

		// DELETE
		_, err = coll.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		if err != nil {
			t.Fatalf("Error eliminando: %v", err)
		}

		t.Log("✅ CRUD notifications OK")
	}
}

func testIndexesValidation(ctx context.Context, db *mongo.Database) func(*testing.T) {
	return func(t *testing.T) {
		// Verificar que los índices fueron creados
		coll := db.Collection("material_assessment")

		cursor, err := coll.Indexes().List(ctx)
		if err != nil {
			t.Fatalf("Error listando índices: %v", err)
		}
		defer cursor.Close(ctx)

		var indexes []bson.M
		if err := cursor.All(ctx, &indexes); err != nil {
			t.Fatalf("Error decodificando índices: %v", err)
		}

		if len(indexes) < 2 { // Al menos _id + alguno creado
			t.Errorf("Se esperaban más índices, se encontraron: %d", len(indexes))
		}

		t.Logf("✅ Índices creados: %d", len(indexes))
	}
}
