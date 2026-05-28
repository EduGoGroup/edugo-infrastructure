package migrations

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/demo"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	postgresSeeds "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"
)

// MigrateOptions configura el comportamiento de Migrate.
type MigrateOptions struct {
	Force          bool   // Eliminar y recrear todos los schemas
	SeedDemo       bool   // Incluir datos de desarrollo (demo seed)
	SeedUpToLayer  string // Aplicar system seed hasta esta capa (vacío = todas)
	Playground     string // Si se setea, aplica el playground tras ApplySystem (omite demo)
	PlaygroundV2   string // Si se setea, aplica el playground v2 tras ApplySystem (omite demo). Mutuamente excluyente con Playground.
	DBUser         string // Usuario PostgreSQL (para GRANT tras drop)
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
// Flujo: ApplyPreGORM → autoMigrateAll (GORM) → ApplyPostGORM → seeds.
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

	// Step 1: schemas, extensions, ENUM types, shared trigger functions
	if err := ApplyPreGORM(db); err != nil {
		return MigrateResult{}, fmt.Errorf("error en pre-GORM: %w", err)
	}

	// Step 2: GORM AutoMigrate (idempotent — adds missing columns/indexes)
	gdb, err := openGORM(db)
	if err != nil {
		return MigrateResult{}, fmt.Errorf("error abriendo GORM: %w", err)
	}
	if err := autoMigrateAll(gdb); err != nil {
		return MigrateResult{}, fmt.Errorf("error en AutoMigrate: %w", err)
	}

	// Step 3: triggers, views, IAM functions, partial indexes, analytics tables
	if err := ApplyPostGORM(db); err != nil {
		return MigrateResult{}, fmt.Errorf("error en post-GORM: %w", err)
	}

	if err := system.ApplySystem(db, opts.SeedUpToLayer); err != nil {
		return MigrateResult{}, fmt.Errorf("error aplicando system seeds: %w", err)
	}

	if opts.Playground != "" {
		// playground.Apply maneja "all" expandiendo a todos los registrados
		// (ver seeds/playground/playground.go). Acá tratamos al valor como
		// opaco: pasamos lo que venga desde el CLI sin special-casing.
		if err := playground.Apply(gdb, opts.Playground); err != nil {
			return MigrateResult{}, fmt.Errorf("error aplicando playground %q: %w", opts.Playground, err)
		}
	} else if opts.PlaygroundV2 != "" {
		// playground_v2.Apply usa registry propio. Sigue el mismo contrato
		// que playground.Apply: "all" expande a todos.
		if err := playground_v2.Apply(gdb, opts.PlaygroundV2); err != nil {
			return MigrateResult{}, fmt.Errorf("error aplicando playground_v2 %q: %w", opts.PlaygroundV2, err)
		}
	} else if opts.SeedDemo {
		if err := demo.ApplyDemo(gdb); err != nil {
			return MigrateResult{}, fmt.Errorf("error aplicando demo seeds: %w", err)
		}
	}

	desc := "Migracion completa (GORM AutoMigrate)"
	if opts.Force {
		desc = "Migracion forzada (recreacion completa, GORM AutoMigrate)"
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
