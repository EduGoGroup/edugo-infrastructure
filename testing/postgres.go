package testing

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	migrationsTable = "schema_migrations"
)

// Migration representa una migración de base de datos
type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// ApplyMigrations ejecuta todas las migraciones pendientes desde un directorio
// Este método es seguro para usar en tests - crea tabla de control si no existe
func ApplyMigrations(db *sql.DB, migrationsPath string) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("error creando tabla de migraciones: %w", err)
	}

	migrations, err := loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("error cargando migraciones: %w", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("error obteniendo migraciones aplicadas: %w", err)
	}

	pendingCount := 0
	for _, m := range migrations {
		if _, exists := applied[m.Version]; exists {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error iniciando transacción: %w", err)
		}

		if _, err := tx.Exec(m.UpSQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error ejecutando migración %d (%s): %w", m.Version, m.Name, err)
		}

		insertQuery := fmt.Sprintf("INSERT INTO %s (version, name) VALUES ($1, $2)", migrationsTable)
		if _, err := tx.Exec(insertQuery, m.Version, m.Name); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error registrando migración: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error confirmando migración: %w", err)
		}

		pendingCount++
	}

	return nil
}

// ApplySeeds ejecuta todos los archivos SQL de un directorio de seeds
// Los ejecuta en orden correcto considerando dependencias de FK
func ApplySeeds(db *sql.DB, seedsPath string) error {
	files, err := os.ReadDir(seedsPath)
	if err != nil {
		return fmt.Errorf("error leyendo directorio de seeds: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// Ordenar seeds por dependencias (orden correcto para FKs)
	orderedSeeds := orderSeedsByDependencies(sqlFiles)

	// Ejecutar seeds en orden
	for _, filename := range orderedSeeds {
		fullPath := filepath.Join(seedsPath, filename)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("error leyendo seed %s: %w", filename, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("error ejecutando seed %s: %w", filename, err)
		}
	}

	return nil
}

// orderSeedsByDependencies ordena seeds respetando dependencias de FK
// Orden: users -> schools -> materials (y otros que dependan de schools)
func orderSeedsByDependencies(files []string) []string {
	// Definir orden de prioridad (menor número = se ejecuta primero)
	priority := map[string]int{
		"users.sql":     1, // Sin dependencias
		"schools.sql":   2, // Sin dependencias
		"materials.sql": 3, // Depende de schools
		// Agregar más según sea necesario
	}

	// Ordenar por prioridad, luego alfabético
	sort.Slice(files, func(i, j int) bool {
		priI, okI := priority[files[i]]
		priJ, okJ := priority[files[j]]

		// Si ambos tienen prioridad definida, usar esa
		if okI && okJ {
			return priI < priJ
		}

		// Si solo uno tiene prioridad, ese va primero
		if okI {
			return true
		}
		if okJ {
			return false
		}

		// Si ninguno tiene prioridad, orden alfabético
		return files[i] < files[j]
	})

	return files
}

// CleanDatabase trunca todas las tablas excepto schema_migrations
// Útil para limpiar datos entre tests
func CleanDatabase(db *sql.DB) error {
	// Obtener lista de tablas
	query := `
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename != $1
		ORDER BY tablename
	`

	rows, err := db.Query(query, migrationsTable)
	if err != nil {
		return fmt.Errorf("error obteniendo tablas: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}

	// Truncar tablas en orden inverso (por FKs)
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Deshabilitar triggers temporalmente
	if _, err := tx.Exec("SET session_replication_role = replica"); err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, table := range tables {
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := tx.Exec(truncateQuery); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error truncando tabla %s: %w", table, err)
		}
	}

	// Rehabilitar triggers
	if _, err := tx.Exec("SET session_replication_role = DEFAULT"); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Funciones privadas reutilizadas del CLI

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

func loadMigrations(migrationsPath string) ([]Migration, error) {
	files, err := os.ReadDir(migrationsPath)
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

		content, err := os.ReadFile(filepath.Join(migrationsPath, name))
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
	for i := 0; i < len(migrations)-1; i++ {
		for j := i + 1; j < len(migrations); j++ {
			if migrations[i].Version > migrations[j].Version {
				migrations[i], migrations[j] = migrations[j], migrations[i]
			}
		}
	}

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[int]*time.Time, error) {
	query := fmt.Sprintf("SELECT version, applied_at FROM %s ORDER BY version", migrationsTable)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
