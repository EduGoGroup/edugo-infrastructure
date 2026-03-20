package migrations

import (
	"database/sql"
	"fmt"
)

// schemaRecord representa un registro interno de la tabla public.schema_version.
type schemaRecord struct {
	Version     string
	ContentHash string
	ExecutionID string
}

// readSchemaVersion lee el último registro de schema_version en la BD.
func readSchemaVersion(db *sql.DB) (schemaRecord, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'schema_version'
	)`).Scan(&exists)
	if err != nil || !exists {
		return schemaRecord{}, fmt.Errorf("tabla schema_version no existe")
	}

	var r schemaRecord
	err = db.QueryRow(`
		SELECT version, content_hash, execution_id::text
		FROM public.schema_version
		ORDER BY id DESC LIMIT 1
	`).Scan(&r.Version, &r.ContentHash, &r.ExecutionID)
	if err != nil {
		return schemaRecord{}, fmt.Errorf("sin registros en schema_version")
	}
	return r, nil
}

// writeSchemaVersion inserta un registro en schema_version.
// execution_id lo genera PostgreSQL (gen_random_uuid).
func writeSchemaVersion(db *sql.DB, version, contentHash string, forced bool, description string) (schemaRecord, error) {
	var r schemaRecord
	err := db.QueryRow(`
		INSERT INTO public.schema_version (version, content_hash, forced, description)
		VALUES ($1, $2, $3, $4)
		RETURNING version, content_hash, execution_id::text
	`, version, contentHash, forced, description).Scan(&r.Version, &r.ContentHash, &r.ExecutionID)
	return r, err
}

// hasTables verifica si ya existen las tablas del dominio (idempotencia).
func hasTables(db *sql.DB) bool {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'auth' AND table_name = 'users'
	)`).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// dropSchemas elimina y recrea todos los schemas para una migración forzada.
func dropSchemas(db *sql.DB, user string) error {
	schemas := []string{"audit", "ui_config", "assessment", "content", "academic", "iam", "auth", "public"}
	for _, schema := range schemas {
		if _, err := db.Exec("DROP SCHEMA IF EXISTS " + schema + " CASCADE"); err != nil {
			return fmt.Errorf("error eliminando schema %s: %w", schema, err)
		}
	}
	if _, err := db.Exec("CREATE SCHEMA public"); err != nil {
		return fmt.Errorf("error creando schema public: %w", err)
	}
	if _, err := db.Exec("GRANT ALL ON SCHEMA public TO " + user); err != nil {
		return fmt.Errorf("error otorgando permisos al usuario: %w", err)
	}
	if _, err := db.Exec("GRANT ALL ON SCHEMA public TO public"); err != nil {
		return fmt.Errorf("error otorgando permisos públicos: %w", err)
	}
	return nil
}
