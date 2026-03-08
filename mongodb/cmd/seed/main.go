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
	command := "all"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if command == "help" || command == "--help" || command == "-h" {
		printHelp()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(buildMongoURI()))
	if err != nil {
		log.Fatalf("error conectando a MongoDB: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("error validando conexión a MongoDB: %v", err)
	}

	db := client.Database(envFirst("MONGO_DB_NAME", "MONGO_DB", "edugo"))

	switch command {
	case "canonical":
		fmt.Println("Aplicando seeds canónicos MongoDB...")
		if err := migrations.ApplySeeds(ctx, db); err != nil {
			log.Fatalf("error aplicando seeds canónicos: %v", err)
		}
	case "mock":
		fmt.Println("Aplicando mock data MongoDB...")
		if err := migrations.ApplyMockData(ctx, db); err != nil {
			log.Fatalf("error aplicando mock data: %v", err)
		}
	case "all":
		fmt.Println("Aplicando seeds canónicos MongoDB...")
		if err := migrations.ApplySeeds(ctx, db); err != nil {
			log.Fatalf("error aplicando seeds canónicos: %v", err)
		}
		fmt.Println("Aplicando mock data MongoDB...")
		if err := migrations.ApplyMockData(ctx, db); err != nil {
			log.Fatalf("error aplicando mock data: %v", err)
		}
	default:
		printHelp()
		os.Exit(1)
	}

	fmt.Println("Seeds MongoDB aplicados correctamente")
}

func printHelp() {
	fmt.Println("MongoDB Seed Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/seed canonical  Aplicar solo seeds canónicos")
	fmt.Println("  go run ./cmd/seed mock       Aplicar solo mock data")
	fmt.Println("  go run ./cmd/seed all        Aplicar ambos conjuntos")
	fmt.Println("")
	fmt.Println("Variables soportadas:")
	fmt.Println("  MONGO_URI")
	fmt.Println("  MONGO_HOST / MONGO_PORT / MONGO_USER / MONGO_PASSWORD / MONGO_DB o MONGO_DB_NAME")
}

func buildMongoURI() string {
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}

	host := envFirst("MONGO_HOST", "", "localhost")
	port := envFirst("MONGO_PORT", "", "27017")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")

	if user != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port)
}

func envFirst(primary, secondary, fallback string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if secondary != "" {
		if value := os.Getenv(secondary); value != "" {
			return value
		}
	}
	return fallback
}
