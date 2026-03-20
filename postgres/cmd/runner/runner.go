package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/internal/dbutil"
	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	postgresSeeds "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds"
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

	db, err := sql.Open("postgres", dbutil.BuildDBURL())
	if err != nil {
		log.Fatalf("error conectando a PostgreSQL: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("error validando conexión PostgreSQL: %v", err)
	}

	switch command {
	case "structure":
		fmt.Println("Aplicando estructura PostgreSQL...")
		if err := postgresMigrations.ApplyAll(db); err != nil {
			log.Fatalf("error aplicando estructura: %v", err)
		}
	case "production-seeds":
		fmt.Println("Aplicando seeds de producción...")
		if err := postgresSeeds.ApplyProduction(db); err != nil {
			log.Fatalf("error aplicando seeds de producción: %v", err)
		}
	case "development-seeds":
		fmt.Println("Aplicando seeds de desarrollo...")
		if err := postgresSeeds.ApplyDevelopment(db); err != nil {
			log.Fatalf("error aplicando seeds de desarrollo: %v", err)
		}
	case "all":
		fmt.Println("Aplicando estructura PostgreSQL...")
		if err := postgresMigrations.ApplyAll(db); err != nil {
			log.Fatalf("error aplicando estructura: %v", err)
		}
		fmt.Println("Aplicando seeds de producción...")
		if err := postgresSeeds.ApplyProduction(db); err != nil {
			log.Fatalf("error aplicando seeds de producción: %v", err)
		}
		fmt.Println("Aplicando seeds de desarrollo...")
		if err := postgresSeeds.ApplyDevelopment(db); err != nil {
			log.Fatalf("error aplicando seeds de desarrollo: %v", err)
		}
	default:
		printHelp()
		os.Exit(1)
	}

	fmt.Println("Runner PostgreSQL completado")
}

func printHelp() {
	fmt.Println("PostgreSQL Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/runner structure          Aplicar estructura embebida")
	fmt.Println("  go run ./cmd/runner production-seeds   Aplicar seeds de producción")
	fmt.Println("  go run ./cmd/runner development-seeds  Aplicar seeds de desarrollo")
	fmt.Println("  go run ./cmd/runner all                Aplicar estructura + producción + desarrollo")
	fmt.Println("")
	fmt.Println("Variables soportadas:")
	fmt.Println("  DATABASE_URL")
	fmt.Println("  DB_HOST / DB_PORT / DB_NAME / DB_USER / DB_PASSWORD / DB_SSL_MODE")
	fmt.Println("  POSTGRES_HOST / POSTGRES_PORT / POSTGRES_DB / POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_SSLMODE")
}
