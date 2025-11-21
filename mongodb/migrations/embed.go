package migrations

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations/constraints"
	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations/structure"
	"go.mongodb.org/mongo-driver/mongo"
)

// Nota: structure y constraints son subpaquetes del módulo migrations

// ApplyAll ejecuta structure + constraints (base de datos limpia lista para usar)
// Equivalente a: ApplyStructure() + ApplyConstraints()
//
// Uso típico: Inicializar base de datos en ambiente de desarrollo o testing
//
// Ejemplo:
//
//	import "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
//	if err := migrations.ApplyAll(ctx, db); err != nil {
//	    log.Fatal(err)
//	}
func ApplyAll(ctx context.Context, db *mongo.Database) error {
	if err := ApplyStructure(ctx, db); err != nil {
		return fmt.Errorf("error aplicando structure: %w", err)
	}
	if err := ApplyConstraints(ctx, db); err != nil {
		return fmt.Errorf("error aplicando constraints: %w", err)
	}
	return nil
}

// ApplyStructure ejecuta todas las funciones de structure/ (createCollection con validators)
// Crea las collections base con validación de schema
//
// Uso típico: Cuando necesitas crear collections en orden específico
//
// Ejemplo:
//
//	migrations.ApplyStructure(ctx, db)
func ApplyStructure(ctx context.Context, db *mongo.Database) error {
	structureFuncs := []func(context.Context, *mongo.Database) error{
		structure.CreateMaterialAssessment,
		structure.CreateMaterialContent,
		structure.CreateAssessmentAttemptResult,
		structure.CreateAuditLogs,
		structure.CreateNotifications,
		structure.CreateAnalyticsEvents,
		structure.CreateMaterialSummary,
		structure.CreateMaterialAssessmentWorker,
		structure.CreateMaterialEvent,
	}

	for _, fn := range structureFuncs {
		if err := fn(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

// ApplyConstraints ejecuta todas las funciones de constraints/ (createIndex)
// DEBE ejecutarse DESPUÉS de ApplyStructure()
//
// Uso típico: Agregar índices después de haber creado las collections
//
// Ejemplo:
//
//	migrations.ApplyStructure(ctx, db)
//	migrations.ApplyConstraints(ctx, db)
func ApplyConstraints(ctx context.Context, db *mongo.Database) error {
	constraintsFuncs := []func(context.Context, *mongo.Database) error{
		constraints.CreateMaterialAssessmentIndexes,
		constraints.CreateMaterialContentIndexes,
		constraints.CreateAssessmentAttemptResultIndexes,
		constraints.CreateAuditLogsIndexes,
		constraints.CreateNotificationsIndexes,
		constraints.CreateAnalyticsEventsIndexes,
		constraints.CreateMaterialSummaryIndexes,
		constraints.CreateMaterialAssessmentWorkerIndexes,
		constraints.CreateMaterialEventIndexes,
	}

	for _, fn := range constraintsFuncs {
		if err := fn(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

// ApplySeeds ejecuta seeds (datos iniciales del ecosistema)
// Por ahora no implementado - agregar cuando se definan seeds necesarios
//
// Uso típico: Inicializar datos mínimos en ambiente de producción/staging
func ApplySeeds(ctx context.Context, db *mongo.Database) error {
	// TODO: Implementar cuando se definan seeds
	return nil
}

// ApplyMockData ejecuta datos mock para testing
// Por ahora no implementado - agregar cuando se definan datos de prueba
//
// Uso típico: Tests de integración, ambiente de desarrollo
func ApplyMockData(ctx context.Context, db *mongo.Database) error {
	// TODO: Implementar cuando se definan datos mock
	return nil
}

// ListFunctions lista todas las funciones disponibles por capa
// Útil para debugging y documentación
//
// Retorna map con estructura:
//
//	{
//	  "structure": ["CreateMaterialAssessment", "CreateMaterialContent", ...],
//	  "constraints": ["CreateMaterialAssessmentIndexes", ...],
//	  "seeds": [],
//	  "testing": []
//	}
func ListFunctions() map[string][]string {
	return map[string][]string{
		"structure": {
			"CreateMaterialAssessment",
			"CreateMaterialContent",
			"CreateAssessmentAttemptResult",
			"CreateAuditLogs",
			"CreateNotifications",
			"CreateAnalyticsEvents",
			"CreateMaterialSummary",
			"CreateMaterialAssessmentWorker",
			"CreateMaterialEvent",
		},
		"constraints": {
			"CreateMaterialAssessmentIndexes",
			"CreateMaterialContentIndexes",
			"CreateAssessmentAttemptResultIndexes",
			"CreateAuditLogsIndexes",
			"CreateNotificationsIndexes",
			"CreateAnalyticsEventsIndexes",
			"CreateMaterialSummaryIndexes",
			"CreateMaterialAssessmentWorkerIndexes",
			"CreateMaterialEventIndexes",
		},
		"seeds":   {},
		"testing": {},
	}
}
