package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

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

	db, err := sql.Open("postgres", buildDBURL())
	if err != nil {
		log.Fatalf("error conectando a PostgreSQL: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("error validando conexión PostgreSQL: %v", err)
	}

	switch command {
	case "production":
		fmt.Println("Aplicando seeds de producción...")
		if err := postgresSeeds.ApplyProduction(db); err != nil {
			log.Fatalf("error aplicando seeds de producción: %v", err)
		}
	case "development":
		fmt.Println("Aplicando seeds de desarrollo...")
		if err := postgresSeeds.ApplyDevelopment(db); err != nil {
			log.Fatalf("error aplicando seeds de desarrollo: %v", err)
		}
	case "all":
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

	fmt.Println("Seeds PostgreSQL aplicados correctamente")
}

func printHelp() {
	fmt.Println("PostgreSQL Seed Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/seed production   Aplicar solo seeds de producción")
	fmt.Println("  go run ./cmd/seed development  Aplicar solo seeds de desarrollo")
	fmt.Println("  go run ./cmd/seed all          Aplicar ambos conjuntos")
	fmt.Println("")
	fmt.Println("Variables soportadas:")
	fmt.Println("  DATABASE_URL")
	fmt.Println("  DB_HOST / DB_PORT / DB_NAME / DB_USER / DB_PASSWORD / DB_SSL_MODE")
	fmt.Println("  POSTGRES_HOST / POSTGRES_PORT / POSTGRES_DB / POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_SSLMODE")
}

func buildDBURL() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	host := envFirst("DB_HOST", "POSTGRES_HOST", "localhost")
	port := envFirst("DB_PORT", "POSTGRES_PORT", "5432")
	name := envFirst("DB_NAME", "POSTGRES_DB", "edugo_dev")
	user := envFirst("DB_USER", "POSTGRES_USER", "edugo")
	password := envFirst("DB_PASSWORD", "POSTGRES_PASSWORD", "changeme")
	sslmode := envFirst("DB_SSL_MODE", "POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
}

func envFirst(primary, secondary, fallback string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if value := os.Getenv(secondary); value != "" {
		return value
	}
	return fallback
}
