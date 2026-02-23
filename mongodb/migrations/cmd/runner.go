package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(dbName)

	switch command {
	case "structure":
		fmt.Println("🏗️  Ejecutando Structure...")
		if err := migrations.ApplyStructure(ctx, db); err != nil {
			log.Fatalf("Error ejecutando structure: %v", err)
		}
		fmt.Println("✅ Structure completada")
	case "constraints":
		fmt.Println("🔗 Ejecutando Constraints...")
		if err := migrations.ApplyConstraints(ctx, db); err != nil {
			log.Fatalf("Error ejecutando constraints: %v", err)
		}
		fmt.Println("✅ Constraints completadas")
	case "all":
		if err := migrations.ApplyAll(ctx, db); err != nil {
			log.Fatalf("Error ejecutando migraciones: %v", err)
		}
		fmt.Println("✅ Todas las migraciones completadas")
	default:
		printHelp()
		os.Exit(1)
	}
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
