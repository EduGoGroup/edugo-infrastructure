package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

const (
	migrationsCollection = "schema_migrations"

	// Timeouts configurables
	DefaultConnectTimeout   = 10 * time.Second
	DefaultOperationTimeout = 5 * time.Second
)

type Migration struct {
	Version   int
	Name      string
	AppliedAt *time.Time
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	mongoURI := getMongoURI()
	dbName := getEnv("MONGO_DB_NAME", "edugo")

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Error("error conectando a MongoDB", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			logger.Error("error desconectando MongoDB", "error", err)
		}
	}()

	// Ping para validar conexión
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("error validando conexión", "error", err)
		os.Exit(1)
	}

	db := client.Database(dbName)

	if err := ensureMigrationsCollection(db); err != nil {
		logger.Error("error creando colección de migraciones", "error", err)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "status":
		if err := showStatus(db); err != nil {
			logger.Error("error mostrando estado", "error", err)
			os.Exit(1)
		}
	case "force":
		if len(os.Args) < 3 {
			logger.Error("uso incorrecto", "mensaje", "Uso: go run migrate.go force VERSION")
			os.Exit(1)
		}
		if err := forceMigration(db, os.Args[2]); err != nil {
			logger.Error("error forzando versión", "error", err)
			os.Exit(1)
		}
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("CLI de Migraciones MongoDB - edugo-infrastructure")
	fmt.Println("")
	fmt.Println("📢 NOTA: Este CLI solo gestiona el estado de migraciones.")
	fmt.Println("   Las migraciones reales se ejecutan desde código Go usando:")
	fmt.Println("   - migrations.ApplyAll(ctx, db)")
	fmt.Println("   - migrations.ApplyStructure(ctx, db)")
	fmt.Println("   - migrations.ApplyConstraints(ctx, db)")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run migrate.go status                Ver estado de migraciones")
	fmt.Println("  go run migrate.go force VERSION         Forzar versión (¡cuidado!)")
	fmt.Println("")
	fmt.Println("Variables de entorno:")
	fmt.Println("  MONGO_HOST     (default: localhost)")
	fmt.Println("  MONGO_PORT     (default: 27017)")
	fmt.Println("  MONGO_DB_NAME  (default: edugo)")
	fmt.Println("  MONGO_USER     (opcional)")
	fmt.Println("  MONGO_PASSWORD (opcional)")
	fmt.Println("")
	fmt.Println("Para ejecutar migraciones, usar el paquete migrations desde Go:")
	fmt.Println("  import \"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations\"")
	fmt.Println("  migrations.ApplyAll(ctx, db)")
}

func getMongoURI() string {
	host := getEnv("MONGO_HOST", "localhost")
	port := getEnv("MONGO_PORT", "27017")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")

	if user != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ensureMigrationsCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultOperationTimeout)
	defer cancel()

	// Verificar si la colección existe
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": migrationsCollection})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		// Crear la colección
		if err := db.CreateCollection(ctx, migrationsCollection); err != nil {
			return err
		}

		// Crear índice único en version
		collection := db.Collection(migrationsCollection)
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "version", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
			return err
		}
	}

	return nil
}

func showStatus(db *mongo.Database) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Println("Estado de Migraciones MongoDB:")
	fmt.Println("==============================")
	fmt.Println("")

	if len(applied) == 0 {
		fmt.Println("⬜ No hay migraciones registradas")
		fmt.Println("")
		fmt.Println("Para aplicar migraciones, usar desde Go:")
		fmt.Println("  migrations.ApplyAll(ctx, db)")
		return nil
	}

	// Mostrar migraciones aplicadas ordenadas
	versions := make([]int, 0, len(applied))
	for v := range applied {
		versions = append(versions, v)
	}

	// Ordenar versiones
	for i := 0; i < len(versions); i++ {
		for j := i + 1; j < len(versions); j++ {
			if versions[i] > versions[j] {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

	for _, v := range versions {
		appliedAt := applied[v]
		fmt.Printf("✅ Versión %03d (aplicada: %s)\n",
			v, appliedAt.Format("2006-01-02 15:04"))
	}

	fmt.Println("")
	fmt.Printf("Total: %d migración(es) aplicada(s)\n", len(applied))

	return nil
}

func forceMigration(db *mongo.Database, version string) error {
	logger.Warn("forzando versión de migración", "version", version)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultOperationTimeout)
	defer cancel()

	collection := db.Collection(migrationsCollection)

	// Eliminar todos los registros
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	// Insertar versión forzada
	versionNum, err := strconv.Atoi(version)
	if err != nil {
		return fmt.Errorf("versión inválida: %s", version)
	}

	now := time.Now()
	if _, err := collection.InsertOne(ctx, bson.M{
		"version":    versionNum,
		"name":       "forced",
		"applied_at": now,
	}); err != nil {
		return err
	}

	logger.Info("versión forzada exitosamente", "version", versionNum)
	return nil
}

func getAppliedMigrations(db *mongo.Database) (map[int]*time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultOperationTimeout)
	defer cancel()

	collection := db.Collection(migrationsCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	applied := make(map[int]*time.Time)
	for cursor.Next(ctx) {
		var result struct {
			Version   int       `bson:"version"`
			AppliedAt time.Time `bson:"applied_at"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		applied[result.Version] = &result.AppliedAt
	}

	return applied, nil
}
