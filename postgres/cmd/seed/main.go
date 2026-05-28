package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/internal/dbutil"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/demo"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"
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
	case "production":
		fmt.Println("Aplicando seeds del sistema...")
		if err := system.ApplySystem(db, ""); err != nil {
			log.Fatalf("error aplicando seeds del sistema: %v", err)
		}
	case "development":
		fmt.Println("Aplicando seeds de desarrollo...")
		if err := applyDemoWithSQL(db); err != nil {
			log.Fatalf("error aplicando seeds de desarrollo: %v", err)
		}
	case "all":
		fmt.Println("Aplicando seeds del sistema...")
		if err := system.ApplySystem(db, ""); err != nil {
			log.Fatalf("error aplicando seeds del sistema: %v", err)
		}
		fmt.Println("Aplicando seeds de desarrollo...")
		if err := applyDemoWithSQL(db); err != nil {
			log.Fatalf("error aplicando seeds de desarrollo: %v", err)
		}
	default:
		printHelp()
		os.Exit(1)
	}

	fmt.Println("Seeds PostgreSQL aplicados correctamente")
}

// applyDemoWithSQL abre un *gorm.DB desde el *sql.DB dado y aplica el seed demo.
func applyDemoWithSQL(db *sql.DB) error {
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("error abriendo GORM: %w", err)
	}
	return demo.ApplyDemo(gdb)
}

func printHelp() {
	fmt.Println("PostgreSQL Seed Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/seed production   Aplicar solo seeds del sistema")
	fmt.Println("  go run ./cmd/seed development  Aplicar solo seeds de desarrollo")
	fmt.Println("  go run ./cmd/seed all          Aplicar ambos conjuntos")
	fmt.Println("")
	fmt.Println("Variables soportadas:")
	fmt.Println("  DATABASE_URL")
	fmt.Println("  DB_HOST / DB_PORT / DB_NAME / DB_USER / DB_PASSWORD / DB_SSL_MODE")
	fmt.Println("  POSTGRES_HOST / POSTGRES_PORT / POSTGRES_DB / POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_SSLMODE")
}
