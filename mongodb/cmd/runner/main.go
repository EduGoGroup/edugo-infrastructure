package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/internal/mongodbutil"
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

	client, err := mongo.Connect(options.Client().ApplyURI(mongodbutil.BuildMongoURI()))
	if err != nil {
		log.Fatalf("error conectando a MongoDB: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("error validando conexión a MongoDB: %v", err)
	}

	db := client.Database(mongodbutil.EnvFirst("MONGO_DB_NAME", "MONGO_DB", "edugo"))

	switch command {
	case "structure":
		fmt.Println("Aplicando structure...")
		if err := migrations.ApplyStructure(ctx, db); err != nil {
			log.Fatalf("error aplicando structure: %v", err)
		}
	case "constraints":
		fmt.Println("Aplicando constraints...")
		if err := migrations.ApplyConstraints(ctx, db); err != nil {
			log.Fatalf("error aplicando constraints: %v", err)
		}
	case "all":
		fmt.Println("Aplicando migraciones MongoDB...")
		if err := migrations.ApplyAll(ctx, db); err != nil {
			log.Fatalf("error aplicando migraciones: %v", err)
		}
	default:
		printHelp()
		os.Exit(1)
	}

	fmt.Println("Runner MongoDB completado")
}

func printHelp() {
	fmt.Println("MongoDB Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/runner structure    Aplicar solo structure")
	fmt.Println("  go run ./cmd/runner constraints  Aplicar solo constraints")
	fmt.Println("  go run ./cmd/runner all          Aplicar structure + constraints")
	fmt.Println("")
	fmt.Println("Variables soportadas:")
	fmt.Println("  MONGO_URI")
	fmt.Println("  MONGO_HOST / MONGO_PORT / MONGO_USER / MONGO_PASSWORD / MONGO_DB o MONGO_DB_NAME")
}
