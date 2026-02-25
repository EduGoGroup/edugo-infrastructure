package migrations_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	defer func() { _ = container.Terminate(ctx) }()

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
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	db := client.Database("testdb")

	t.Log("✅ Container MongoDB creado y conectado")

	// Ejecutar tests
	t.Run("ApplyAll", testApplyAll(ctx, db))
	t.Run("ApplySeeds", testApplySeeds(ctx, db))
	t.Run("ApplyMockData", testApplyMockData(ctx, db))
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
		// El driver de MongoDB v2 devuelve documentos embebidos como bson.D
		var difficultyVal interface{}
		if md, ok := retrieved["metadata"].(bson.D); ok {
			for _, elem := range md {
				if elem.Key == "difficulty" {
					difficultyVal = elem.Value
					break
				}
			}
		}
		if difficultyVal != "medium" {
			t.Errorf("Update no aplicado correctamente, difficulty=%v", difficultyVal)
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
		// Verificar índices en las collections que sí los tienen definidos
		collectionsWithIndexes := map[string]int{
			"material_summary":           2, // _id + idx_material_id (unique) + más
			"material_assessment_worker": 2, // _id + idx_material_id + más
			"material_event":             2, // _id + índices definidos
		}

		for collName, minIndexes := range collectionsWithIndexes {
			coll := db.Collection(collName)
			cursor, err := coll.Indexes().List(ctx)
			if err != nil {
				t.Fatalf("Error listando índices de %s: %v", collName, err)
			}

			var indexes []bson.M
			if err := cursor.All(ctx, &indexes); err != nil {
				_ = cursor.Close(ctx)
				t.Fatalf("Error decodificando índices de %s: %v", collName, err)
			}
			_ = cursor.Close(ctx)

			if len(indexes) < minIndexes {
				t.Errorf("Collection %s: se esperaban al menos %d índices, se encontraron %d",
					collName, minIndexes, len(indexes))
			} else {
				t.Logf("✅ Collection %s: %d índices creados correctamente", collName, len(indexes))
			}
		}
	}
}

func testApplySeeds(ctx context.Context, db *mongo.Database) func(*testing.T) {
	return func(t *testing.T) {
		// Aplicar seeds
		if err := migrations.ApplySeeds(ctx, db); err != nil {
			t.Fatalf("Error aplicando seeds: %v", err)
		}

		// Verificar que se insertaron documentos en las colecciones activas
		expectedCounts := map[string]int64{
			"material_assessment_worker": 2,
			"material_summary":           3,
		}

		for collection, expectedCount := range expectedCounts {
			coll := db.Collection(collection)
			count, err := coll.CountDocuments(ctx, bson.M{})
			if err != nil {
				t.Fatalf("Error contando documentos en %s: %v", collection, err)
			}

			if count != expectedCount {
				t.Errorf("Collection %s: se esperaban %d documentos, se encontraron %d",
					collection, expectedCount, count)
			} else {
				t.Logf("✅ Collection %s: %d documentos insertados correctamente", collection, count)
			}
		}

		// Test de idempotencia: ejecutar seeds de nuevo no debe retornar error
		if err := migrations.ApplySeeds(ctx, db); err != nil {
			t.Fatalf("Error en segunda ejecución de seeds (idempotencia): %v", err)
		}

		t.Log("✅ ApplySeeds ejecutado correctamente (re-ejecución sin error)")
	}
}

func testApplyMockData(ctx context.Context, db *mongo.Database) func(*testing.T) {
	return func(t *testing.T) {
		// Aplicar mock data
		if err := migrations.ApplyMockData(ctx, db); err != nil {
			t.Fatalf("Error aplicando mock data: %v", err)
		}

		// Verificar que las collections activas tienen documentos tras mock
		mockCollections := []string{"material_assessment_worker", "material_summary"}
		for _, collName := range mockCollections {
			coll := db.Collection(collName)
			count, err := coll.CountDocuments(ctx, bson.M{})
			if err != nil {
				t.Fatalf("Error contando documentos en %s: %v", collName, err)
			}
			if count == 0 {
				t.Errorf("Collection %s: no se encontraron documentos tras ApplyMockData", collName)
			} else {
				t.Logf("✅ Collection %s: %d documentos presentes", collName, count)
			}
		}

		// Test de idempotencia: re-ejecución no debe retornar error
		if err := migrations.ApplyMockData(ctx, db); err != nil {
			t.Fatalf("Error en segunda ejecución de mock data (idempotencia): %v", err)
		}

		t.Log("✅ ApplyMockData ejecutado correctamente (re-ejecución sin error)")
	}
}
