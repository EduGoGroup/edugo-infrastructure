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
	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/base"
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
	case "structure":
		fmt.Println("Aplicando estructura PostgreSQL...")
		if err := postgresMigrations.ApplyAll(db); err != nil {
			log.Fatalf("error aplicando estructura: %v", err)
		}
	case "production-seeds":
		fmt.Println("Aplicando seeds del sistema...")
		if err := system.ApplySystem(db, ""); err != nil {
			log.Fatalf("error aplicando seeds del sistema: %v", err)
		}
	case "development-seeds":
		fmt.Println("Aplicando seeds de desarrollo...")
		if err := applyDemoWithSQL(db); err != nil {
			log.Fatalf("error aplicando seeds de desarrollo: %v", err)
		}
	case "all":
		fmt.Println("Aplicando estructura PostgreSQL...")
		if err := postgresMigrations.ApplyAll(db); err != nil {
			log.Fatalf("error aplicando estructura: %v", err)
		}
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

	fmt.Println("Runner PostgreSQL completado")
}

// applyDemoWithSQL abre un *gorm.DB desde el *sql.DB dado y aplica el seed de
// desarrollo (playground_v2/base, mundo de datos por defecto de EduGo).
func applyDemoWithSQL(db *sql.DB) error {
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("error abriendo GORM: %w", err)
	}
	return base.Apply(gdb)
}

func printHelp() {
	fmt.Println("PostgreSQL Runner")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/runner structure          Aplicar estructura embebida")
	fmt.Println("  go run ./cmd/runner production-seeds   Aplicar seeds del sistema")
	fmt.Println("  go run ./cmd/runner development-seeds  Aplicar seeds de desarrollo")
	fmt.Println("  go run ./cmd/runner all                Aplicar estructura + sistema + desarrollo")
	fmt.Println("")
	fmt.Println("Variables soportadas:")
	fmt.Println("  DATABASE_URL")
	fmt.Println("  DB_HOST / DB_PORT / DB_NAME / DB_USER / DB_PASSWORD / DB_SSL_MODE")
	fmt.Println("  POSTGRES_HOST / POSTGRES_PORT / POSTGRES_DB / POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_SSLMODE")
}
