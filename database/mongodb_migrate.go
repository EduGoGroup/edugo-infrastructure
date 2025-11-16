//go:build mongodb
// +build mongodb

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	migrationsCollection = "schema_migrations"
	migrationsDir        = "migrations/mongodb"
)

type Migration struct {
	Version   int
	Name      string
	UpScript  string
	DownScript string
	AppliedAt *time.Time
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	mongoURI := getMongoURI()
	dbName := getEnv("MONGO_DB_NAME", "edugo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error desconectando MongoDB: %v", err)
		}
	}()

	// Ping para validar conexión
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Error validando conexión: %v", err)
	}

	db := client.Database(dbName)

	if err := ensureMigrationsCollection(db); err != nil {
		log.Fatalf("Error creando colección de migraciones: %v", err)
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Error ejecutando migraciones: %v", err)
		}
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Error revirtiendo migración: %v", err)
		}
	case "status":
		if err := showStatus(db); err != nil {
			log.Fatalf("Error mostrando estado: %v", err)
		}
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Uso: go run mongodb_migrate.go create \"descripcion_migracion\"")
		}
		if err := createMigration(os.Args[2]); err != nil {
			log.Fatalf("Error creando migración: %v", err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Uso: go run mongodb_migrate.go force VERSION")
		}
		if err := forceMigration(db, os.Args[2]); err != nil {
			log.Fatalf("Error forzando versión: %v", err)
		}
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("CLI de Migraciones MongoDB - edugo-infrastructure")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run mongodb_migrate.go up                    Ejecutar migraciones pendientes")
	fmt.Println("  go run mongodb_migrate.go down                  Revertir última migración")
	fmt.Println("  go run mongodb_migrate.go status                Ver estado de migraciones")
	fmt.Println("  go run mongodb_migrate.go create \"nombre\"       Crear nueva migración")
	fmt.Println("  go run mongodb_migrate.go force VERSION         Forzar versión (¡cuidado!)")
	fmt.Println("")
	fmt.Println("Variables de entorno:")
	fmt.Println("  MONGO_HOST     (default: localhost)")
	fmt.Println("  MONGO_PORT     (default: 27017)")
	fmt.Println("  MONGO_DB_NAME  (default: edugo)")
	fmt.Println("  MONGO_USER     (opcional)")
	fmt.Println("  MONGO_PASSWORD (opcional)")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func migrateUp(db *mongo.Database) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	pendingCount := 0
	for _, m := range migrations {
		if _, exists := applied[m.Version]; exists {
			continue
		}

		fmt.Printf("Ejecutando migración %03d: %s\n", m.Version, m.Name)

		if err := executeMigrationScript(db, m.UpScript); err != nil {
			return fmt.Errorf("error en migración %d: %w", m.Version, err)
		}

		// Registrar migración aplicada
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := db.Collection(migrationsCollection)
		now := time.Now()
		_, err := collection.InsertOne(ctx, bson.M{
			"version":    m.Version,
			"name":       m.Name,
			"applied_at": now,
		})
		if err != nil {
			return err
		}

		pendingCount++
		fmt.Printf("✅ Migración %03d aplicada exitosamente\n", m.Version)
	}

	if pendingCount == 0 {
		fmt.Println("✅ No hay migraciones pendientes")
	} else {
		fmt.Printf("✅ %d migración(es) aplicada(s) exitosamente\n", pendingCount)
	}

	return nil
}

func migrateDown(db *mongo.Database) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		fmt.Println("No hay migraciones para revertir")
		return nil
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	// Encontrar la última versión aplicada
	lastVersion := 0
	for v := range applied {
		if v > lastVersion {
			lastVersion = v
		}
	}

	var targetMigration *Migration
	for i := range migrations {
		if migrations[i].Version == lastVersion {
			targetMigration = &migrations[i]
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migración %d no encontrada", lastVersion)
	}

	fmt.Printf("Revirtiendo migración %03d: %s\n", targetMigration.Version, targetMigration.Name)

	if err := executeMigrationScript(db, targetMigration.DownScript); err != nil {
		return fmt.Errorf("error revirtiendo migración: %w", err)
	}

	// Eliminar registro de migración
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(migrationsCollection)
	_, err = collection.DeleteOne(ctx, bson.M{"version": targetMigration.Version})
	if err != nil {
		return err
	}

	fmt.Printf("✅ Migración %03d revertida exitosamente\n", targetMigration.Version)
	return nil
}

func showStatus(db *mongo.Database) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Println("Estado de Migraciones MongoDB:")
	fmt.Println("==============================")
	fmt.Println("")

	for _, m := range migrations {
		if appliedAt, exists := applied[m.Version]; exists {
			fmt.Printf("✅ %03d: %s (aplicada: %s)\n",
				m.Version, m.Name, appliedAt.Format("2006-01-02 15:04"))
		} else {
			fmt.Printf("⬜ %03d: %s (pendiente)\n", m.Version, m.Name)
		}
	}

	fmt.Println("")
	fmt.Printf("Total: %d migraciones, %d aplicadas, %d pendientes\n",
		len(migrations), len(applied), len(migrations)-len(applied))

	return nil
}

func createMigration(description string) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	nextVersion := 1
	if len(migrations) > 0 {
		lastMigration := migrations[len(migrations)-1]
		nextVersion = lastMigration.Version + 1
	}

	filename := fmt.Sprintf("%03d_%s", nextVersion, sanitizeName(description))

	upFile := filepath.Join(migrationsDir, filename+".up.js")
	downFile := filepath.Join(migrationsDir, filename+".down.js")

	upContent := fmt.Sprintf(`// Migration: %s
// Created: %s

// TODO: Escribir código JavaScript para migración UP
// Ejemplo:
// db.createCollection("new_collection", {
//   validator: {
//     $jsonSchema: {
//       bsonType: "object",
//       required: ["field1", "field2"],
//       properties: {
//         field1: { bsonType: "string" },
//         field2: { bsonType: "int" }
//       }
//     }
//   }
// });

print("✅ Migration %s UP completed");
`,
		description, time.Now().Format("2006-01-02 15:04"), description)

	downContent := fmt.Sprintf(`// Migration DOWN: %s
// Created: %s

// TODO: Escribir código JavaScript para revertir migración
// Ejemplo:
// db.new_collection.drop();

print("✅ Migration %s DOWN completed");
`,
		description, time.Now().Format("2006-01-02 15:04"), description)

	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Migración creada:\n")
	fmt.Printf("   UP:   %s\n", upFile)
	fmt.Printf("   DOWN: %s\n", downFile)
	fmt.Println("")
	fmt.Println("Editar los archivos JavaScript y luego ejecutar: go run mongodb_migrate.go up")

	return nil
}

func forceMigration(db *mongo.Database, version string) error {
	fmt.Printf("⚠️  Forzando versión de migración a: %s\n", version)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(migrationsCollection)

	// Eliminar todos los registros
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	// Insertar versión forzada
	var versionNum int
	if _, err := fmt.Sscanf(version, "%d", &versionNum); err != nil {
		return fmt.Errorf("versión inválida: %s", version)
	}

	now := time.Now()
	_, err = collection.InsertOne(ctx, bson.M{
		"version":    versionNum,
		"name":       "forced",
		"applied_at": now,
	})
	if err != nil {
		return err
	}

	fmt.Println("✅ Versión forzada exitosamente")
	return nil
}

func loadMigrations() ([]Migration, error) {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	migrationsMap := make(map[int]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, ".js") {
			continue
		}

		parts := strings.SplitN(name, "_", 2)
		if len(parts) < 2 {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
			continue
		}

		if migrationsMap[version] == nil {
			migrationsMap[version] = &Migration{
				Version: version,
			}
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(name, ".up.js") {
			migrationsMap[version].UpScript = string(content)
			migrationsMap[version].Name = strings.TrimSuffix(strings.TrimSuffix(parts[1], ".up.js"), ".down.js")
		} else if strings.HasSuffix(name, ".down.js") {
			migrationsMap[version].DownScript = string(content)
		}
	}

	var migrations []Migration
	for _, m := range migrationsMap {
		if m.UpScript != "" && m.DownScript != "" {
			migrations = append(migrations, *m)
		}
	}

	// Ordenar por versión
	for i := 0; i < len(migrations)-1; i++ {
		for j := i + 1; j < len(migrations); j++ {
			if migrations[i].Version > migrations[j].Version {
				migrations[i], migrations[j] = migrations[j], migrations[i]
			}
		}
	}

	return migrations, nil
}

func getAppliedMigrations(db *mongo.Database) (map[int]*time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func executeMigrationScript(db *mongo.Database, script string) error {
	// Guardar script temporalmente
	tmpFile, err := os.CreateTemp("", "migration-*.js")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(script); err != nil {
		return err
	}
	tmpFile.Close()

	// Construir comando mongosh
	mongoURI := getMongoURI()
	dbName := getEnv("MONGO_DB_NAME", "edugo")

	cmd := exec.Command("mongosh", mongoURI+"/"+dbName, "--file", tmpFile.Name(), "--quiet")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
