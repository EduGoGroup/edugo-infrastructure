package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	migrationsTable = "schema_migrations"
	migrationsDir   = "migrations"
)

type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	dbURL := getDBURL()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error conectando a PostgreSQL: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error validando conexión: %v", err)
	}

	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("Error creando tabla de migraciones: %v", err)
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
			log.Fatal("Uso: go run migrate.go create \"descripcion_migracion\"")
		}
		if err := createMigration(os.Args[2]); err != nil {
			log.Fatalf("Error creando migración: %v", err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Uso: go run migrate.go force VERSION")
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
	fmt.Println("CLI de Migraciones PostgreSQL - edugo-infrastructure")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  go run migrate.go up                    Ejecutar migraciones pendientes")
	fmt.Println("  go run migrate.go down                  Revertir última migración")
	fmt.Println("  go run migrate.go status                Ver estado de migraciones")
	fmt.Println("  go run migrate.go create \"nombre\"       Crear nueva migración")
	fmt.Println("  go run migrate.go force VERSION         Forzar versión (¡cuidado!)")
	fmt.Println("")
	fmt.Println("Variables de entorno:")
	fmt.Println("  DB_HOST     (default: localhost)")
	fmt.Println("  DB_PORT     (default: 5432)")
	fmt.Println("  DB_NAME     (default: edugo_dev)")
	fmt.Println("  DB_USER     (default: edugo)")
	fmt.Println("  DB_PASSWORD (default: changeme)")
	fmt.Println("  DB_SSL_MODE (default: disable)")
}

func getDBURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "edugo_dev")
	user := getEnv("DB_USER", "edugo")
	password := getEnv("DB_PASSWORD", "changeme")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ensureMigrationsTable(db *sql.DB) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`, migrationsTable)

	_, err := db.Exec(query)
	return err
}

func migrateUp(db *sql.DB) error {
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

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(m.UpSQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error en migración %d: %w", m.Version, err)
		}

		insertQuery := fmt.Sprintf("INSERT INTO %s (version, name) VALUES ($1, $2)", migrationsTable)
		if _, err := tx.Exec(insertQuery, m.Version, m.Name); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
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

func migrateDown(db *sql.DB) error {
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

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(targetMigration.DownSQL); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error revirtiendo migración: %w", err)
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE version = $1", migrationsTable)
	if _, err := tx.Exec(deleteQuery, targetMigration.Version); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("✅ Migración %03d revertida exitosamente\n", targetMigration.Version)
	return nil
}

func showStatus(db *sql.DB) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Println("Estado de Migraciones:")
	fmt.Println("=====================")
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

	upFile := filepath.Join(migrationsDir, filename+".up.sql")
	downFile := filepath.Join(migrationsDir, filename+".down.sql")

	upContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- TODO: Escribir SQL para migración UP\n",
		description, time.Now().Format("2006-01-02 15:04"))

	downContent := fmt.Sprintf("-- Migration: %s (DOWN)\n-- Created: %s\n\n-- TODO: Escribir SQL para revertir migración\n",
		description, time.Now().Format("2006-01-02 15:04"))

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
	fmt.Println("Editar los archivos SQL y luego ejecutar: go run migrate.go up")

	return nil
}

func forceMigration(db *sql.DB, versionStr string) error {
	fmt.Printf("⚠️  Forzando versión de migración a: %s\n", versionStr)

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return fmt.Errorf("versión inválida, debe ser un número: %s", versionStr)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %s", migrationsTable)
	if _, err := tx.Exec(deleteQuery); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insertar versión forzada
	insertQuery := fmt.Sprintf("INSERT INTO %s (version, name) VALUES ($1, $2)", migrationsTable)
	if _, err := tx.Exec(insertQuery, version, "forced"); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
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
		if !strings.HasSuffix(name, ".sql") {
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

		if strings.HasSuffix(name, ".up.sql") {
			migrationsMap[version].UpSQL = string(content)
			migrationsMap[version].Name = strings.TrimSuffix(strings.TrimSuffix(parts[1], ".up.sql"), ".down.sql")
		} else if strings.HasSuffix(name, ".down.sql") {
			migrationsMap[version].DownSQL = string(content)
		}
	}

	var migrations []Migration
	for _, m := range migrationsMap {
		if m.UpSQL != "" && m.DownSQL != "" {
			migrations = append(migrations, *m)
		}
	}

	// Ordenar por versión
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[int]*time.Time, error) {
	query := fmt.Sprintf("SELECT version, applied_at FROM %s ORDER BY version", migrationsTable)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	applied := make(map[int]*time.Time)
	for rows.Next() {
		var version int
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, err
		}
		applied[version] = &appliedAt
	}

	return applied, nil
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
