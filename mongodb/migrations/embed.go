package migrations

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations/constraints"
	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations/structure"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
// Inserta documentos base en las colecciones necesarias para el funcionamiento del sistema
//
// Características:
//   - Idempotente: Se puede ejecutar múltiples veces sin duplicar datos
//   - Usa ordered: false para continuar aunque algunos documentos ya existan
//   - Retorna error solo si falla la inserción por razones NO de duplicados
//
// Collections pobladas:
//   - analytics_events (6 eventos de ejemplo)
//   - material_assessment (2 assessments de Física y Matemáticas)
//   - audit_logs (5 registros de auditoría)
//   - material_assessment_worker (2 workers con preguntas generadas por IA)
//   - material_summary (3 resúmenes en español, inglés y portugués)
//   - notifications (4 notificaciones de ejemplo)
//
// Uso típico: Inicializar datos mínimos en ambiente de desarrollo/staging
//
// Ejemplo:
//
//	migrations.ApplyAll(ctx, db)
//	migrations.ApplySeeds(ctx, db)  // Datos iniciales
func ApplySeeds(ctx context.Context, db *mongo.Database) error {
	inserted, err := applySeedsInternal(ctx, db)
	if err != nil {
		return fmt.Errorf("error aplicando seeds: %w", err)
	}
	if inserted > 0 {
		// Solo logueamos si se insertó algo (opcional, puede removerse)
		_ = inserted // Evitar warning de variable no usada
	}
	return nil
}

// ApplyMockData ejecuta datos mock para testing y desarrollo
// Inserta documentos de prueba más variados y numerosos que los seeds
//
// Características:
//   - Idempotente: Se puede ejecutar múltiples veces sin duplicar datos
//   - Usa ordered: false para continuar aunque algunos documentos ya existan
//   - Retorna error solo si falla la inserción por razones NO de duplicados
//
// Collections pobladas (35 documentos totales):
//   - analytics_events (10 eventos con diferentes plataformas y roles)
//   - material_assessment (3 assessments de diferentes dificultades)
//   - audit_logs (8 registros de auditoría con más casos de uso)
//   - material_assessment_worker (3 workers con diferentes subjects/idiomas)
//   - material_summary (5 resúmenes en español, inglés, portugués, francés, alemán)
//   - notifications (6 notificaciones de diferentes tipos y prioridades)
//
// Diferencias vs ApplySeeds:
//   - Más documentos por colección (35 vs 22)
//   - Mayor variedad de datos (más países, plataformas, tipos)
//   - Cubre más casos de uso y escenarios
//
// Uso típico: Tests de integración, ambiente de desarrollo, demos
//
// Ejemplo:
//
//	migrations.ApplyAll(ctx, db)
//	migrations.ApplyMockData(ctx, db)  // Datos de prueba
func ApplyMockData(ctx context.Context, db *mongo.Database) error {
	inserted, err := applyMockDataInternal(ctx, db)
	if err != nil {
		return fmt.Errorf("error aplicando mock data: %w", err)
	}
	if inserted > 0 {
		// Solo logueamos si se insertó algo (opcional, puede removerse)
		_ = inserted // Evitar warning de variable no usada
	}
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
			"CreateMaterialSummary",
			"CreateMaterialAssessmentWorker",
			"CreateMaterialEvent",
		},
		"constraints": {
			"CreateMaterialSummaryIndexes",
			"CreateMaterialAssessmentWorkerIndexes",
			"CreateMaterialEventIndexes",
		},
		"seeds":   {},
		"testing": {},
	}
}
