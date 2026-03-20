package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MigrateOptions configura el comportamiento de Migrate.
type MigrateOptions struct {
	Force     bool // Eliminar y recrear la base de datos
	ApplyMock bool // Incluir datos de desarrollo
}

// MigrateResult contiene el resultado de una migración.
type MigrateResult struct {
	Skipped bool // true si la BD ya tenía colecciones
	Forced  bool // true si se ejecutó con Force=true
}

// Migrate ejecuta el flujo completo de migración de MongoDB.
// Si la BD ya tiene colecciones y no se fuerza, retorna Skipped=true.
func Migrate(ctx context.Context, db *mongo.Database, opts MigrateOptions) (MigrateResult, error) {
	if opts.Force {
		if err := db.Drop(ctx); err != nil {
			return MigrateResult{}, fmt.Errorf("error eliminando database: %w", err)
		}
	} else if hasCollections(ctx, db) {
		return MigrateResult{Skipped: true}, nil
	}

	if err := ApplyAll(ctx, db); err != nil {
		return MigrateResult{}, fmt.Errorf("error aplicando migraciones: %w", err)
	}

	if err := ApplySeeds(ctx, db); err != nil {
		return MigrateResult{}, fmt.Errorf("error aplicando seeds: %w", err)
	}

	if opts.ApplyMock {
		if err := ApplyMockData(ctx, db); err != nil {
			return MigrateResult{}, fmt.Errorf("error aplicando mock data: %w", err)
		}
	}

	return MigrateResult{Forced: opts.Force}, nil
}

// hasCollections verifica si la base de datos ya tiene colecciones.
func hasCollections(ctx context.Context, db *mongo.Database) bool {
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return false
	}
	return len(collections) > 0
}
