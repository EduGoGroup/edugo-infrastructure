package migrations

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	postgresSeeds "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds"
)

// MigrateOptions configura el comportamiento de Migrate.
type MigrateOptions struct {
	Force     bool   // Eliminar y recrear todos los schemas
	ApplyMock bool   // Incluir datos de desarrollo
	DBUser    string // Usuario PostgreSQL (para GRANT tras drop)
}

// MigrateResult contiene el resultado de una migración.
type MigrateResult struct {
	Version     string // Versión del schema aplicado
	ContentHash string // Hash combinado de archivos SQL
	ExecutionID string // UUID generado por PostgreSQL
	Skipped     bool   // true si la BD ya estaba actualizada
	NeedsForce  bool   // true si hay tablas pero version/hash no coincide
	Forced      bool   // true si se ejecutó con Force=true
}

// ExpectedContentHash calcula el hash SHA256 combinado de migraciones + seeds.
func ExpectedContentHash() string {
	mHash := ComputeFilesHash()
	sHash := postgresSeeds.ComputeFilesHash()
	combined := sha256.New()
	combined.Write([]byte(mHash))
	combined.Write([]byte(sHash))
	return fmt.Sprintf("%x", combined.Sum(nil))[:16]
}

// Migrate ejecuta el flujo completo de migración de PostgreSQL.
// Si la BD ya tiene tablas y no se fuerza, retorna Skipped=true.
func Migrate(db *sql.DB, opts MigrateOptions) (MigrateResult, error) {
	expectedHash := ExpectedContentHash()

	if opts.Force {
		if err := dropSchemas(db, opts.DBUser); err != nil {
			return MigrateResult{}, fmt.Errorf("error eliminando schemas: %w", err)
		}
	} else if hasTables(db) {
		current, err := readSchemaVersion(db)
		if err == nil && current.Version == SchemaVersion && current.ContentHash == expectedHash {
			return MigrateResult{
				Version:     current.Version,
				ContentHash: current.ContentHash,
				ExecutionID: current.ExecutionID,
				Skipped:     true,
			}, nil
		}
		return MigrateResult{
			Skipped:    true,
			NeedsForce: true,
		}, nil
	}

	if err := ApplyAll(db); err != nil {
		return MigrateResult{}, fmt.Errorf("error aplicando migraciones: %w", err)
	}

	if err := postgresSeeds.ApplyProduction(db); err != nil {
		return MigrateResult{}, fmt.Errorf("error aplicando seeds de producción: %w", err)
	}

	if opts.ApplyMock {
		if err := postgresSeeds.ApplyDevelopment(db); err != nil {
			return MigrateResult{}, fmt.Errorf("error aplicando datos de desarrollo: %w", err)
		}
	}

	desc := "Migracion completa"
	if opts.Force {
		desc = "Migracion forzada (recreacion completa)"
	}

	record, err := writeSchemaVersion(db, SchemaVersion, expectedHash, opts.Force, desc)
	if err != nil {
		return MigrateResult{
			Version:     SchemaVersion,
			ContentHash: expectedHash,
			Forced:      opts.Force,
		}, fmt.Errorf("migracion exitosa pero error registrando version: %w", err)
	}

	return MigrateResult{
		Version:     record.Version,
		ContentHash: record.ContentHash,
		ExecutionID: record.ExecutionID,
		Forced:      opts.Force,
	}, nil
}

// Status retorna el estado actual del schema en la BD sin modificar nada.
func Status(db *sql.DB) (*MigrateResult, error) {
	current, err := readSchemaVersion(db)
	if err != nil {
		return nil, err
	}
	expectedHash := ExpectedContentHash()
	return &MigrateResult{
		Version:     current.Version,
		ContentHash: current.ContentHash,
		ExecutionID: current.ExecutionID,
		NeedsForce:  current.Version != SchemaVersion || current.ContentHash != expectedHash,
	}, nil
}
