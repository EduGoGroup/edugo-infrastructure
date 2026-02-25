package migrations

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// seedDocument representa un conjunto de documentos para insertar en una colección
type seedDocument struct {
	collection string
	documents  []interface{}
}

// applySeedsInternal ejecuta todos los seeds en la base de datos MongoDB
// Retorna el número de documentos insertados y cualquier error
func applySeedsInternal(ctx context.Context, db *mongo.Database) (int, error) {
	seeds := getSeedDocuments()
	totalInserted := 0

	for _, seed := range seeds {
		if len(seed.documents) == 0 {
			continue
		}

		collection := db.Collection(seed.collection)

		// Usar ordered: false para que si un documento ya existe, continúe con los demás
		opts := options.InsertMany().SetOrdered(false)

		result, err := collection.InsertMany(ctx, seed.documents, opts)
		if err != nil {
			// Si es error de duplicados, solo reportamos cuántos se insertaron
			if mongo.IsDuplicateKeyError(err) {
				inserted := 0
				if result != nil {
					inserted = len(result.InsertedIDs)
				}
				totalInserted += inserted
				// Continuamos con la siguiente colección
				continue
			}
			return totalInserted, fmt.Errorf("error insertando seeds en %s: %w", seed.collection, err)
		}

		totalInserted += len(result.InsertedIDs)
	}

	return totalInserted, nil
}

// getSeedDocuments retorna todos los seeds organizados por colección
func getSeedDocuments() []seedDocument {
	return []seedDocument{
		materialAssessmentWorkerSeeds(),
		materialSummarySeeds(),
		// Las siguientes colecciones fueron eliminadas por no uso:
		// - analytics_events (analytics usa servicio externo)
		// - material_assessment (duplicada por material_assessment_worker)
		// - audit_logs (usará SaaS externo)
		// - notifications (push notifications no implementado)
	}
}

// materialAssessmentWorkerSeeds retorna los seeds de la colección material_assessment_worker
func materialAssessmentWorkerSeeds() seedDocument {
	return seedDocument{
		collection: "material_assessment_worker",
		documents: []interface{}{
			// Worker 1 - POO Java
			bson.D{
				{Key: "material_id", Value: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q1111111-1111-1111-1111-111111111111"},
						{Key: "question_text", Value: "¿Cuál es el principio fundamental de la Programación Orientada a Objetos que permite ocultar los detalles de implementación?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Herencia"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Polimorfismo"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Encapsulación"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Abstracción"}},
						}},
						{Key: "correct_answer", Value: "opt3"},
						{Key: "explanation", Value: "La encapsulación es el principio que permite ocultar los detalles internos de implementación y exponer solo lo necesario mediante interfaces públicas."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"POO", "conceptos"}},
					},
					bson.D{
						{Key: "question_id", Value: "q2222222-2222-2222-2222-222222222222"},
						{Key: "question_text", Value: "En Java, ¿una clase puede heredar de múltiples clases?"},
						{Key: "question_type", Value: "true_false"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "true"}, {Key: "option_text", Value: "Verdadero"}},
							bson.D{{Key: "option_id", Value: "false"}, {Key: "option_text", Value: "Falso"}},
						}},
						{Key: "correct_answer", Value: "false"},
						{Key: "explanation", Value: "Java no soporta herencia múltiple de clases para evitar el problema del diamante. Sin embargo, una clase puede implementar múltiples interfaces."},
						{Key: "points", Value: 5},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q3333333-3333-3333-3333-333333333333"},
						{Key: "question_text", Value: "Explica brevemente qué es el polimorfismo y da un ejemplo en Java."},
						{Key: "question_type", Value: "open"},
						{Key: "options", Value: bson.A{}},
						{Key: "correct_answer", Value: "El polimorfismo permite que objetos de diferentes clases sean tratados como objetos de una clase común. Ejemplo: Animal animal = new Perro(); donde Perro extiende Animal."},
						{Key: "explanation", Value: "El polimorfismo permite escribir código más flexible y reutilizable al trabajar con abstracciones en lugar de implementaciones concretas."},
						{Key: "points", Value: 15},
						{Key: "difficulty", Value: "hard"},
						{Key: "tags", Value: bson.A{"POO", "polimorfismo"}},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 30},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4"},
				{Key: "processing_time_ms", Value: 5200},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 1200},
					{Key: "completion_tokens", Value: 450},
					{Key: "total_tokens", Value: 1650},
				}},
				{Key: "created_at", Value: mustParseTime("2025-11-15T10:35:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-11-15T10:35:00Z")},
			},
			// Worker 2 - React Hooks
			bson.D{
				{Key: "material_id", Value: "f1a2b3c4-d5e6-4f5a-9b8c-7d6e5f4a3b2c"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q4444444-4444-4444-4444-444444444444"},
						{Key: "question_text", Value: "Which React Hook is used to perform side effects in functional components?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "useState"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "useEffect"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "useContext"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "useReducer"}},
						}},
						{Key: "correct_answer", Value: "opt2"},
						{Key: "explanation", Value: "useEffect is the Hook used to perform side effects such as data fetching, subscriptions, or manually changing the DOM."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
						{Key: "tags", Value: bson.A{"React", "Hooks"}},
					},
					bson.D{
						{Key: "question_id", Value: "q5555555-5555-5555-5555-555555555555"},
						{Key: "question_text", Value: "What is the correct syntax for creating a custom Hook in React?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "function myHook() {}"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "const myHook = () => {}"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "function useMyHook() {}"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "hook myHook() {}"}},
						}},
						{Key: "correct_answer", Value: "opt3"},
						{Key: "explanation", Value: "Custom Hooks must start with 'use' prefix to follow React conventions and enable linting rules."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"React", "Hooks", "custom"}},
					},
					bson.D{
						{Key: "question_id", Value: "q6666666-6666-4666-8666-666666666666"},
						{Key: "question_text", Value: "What does the useCallback Hook do in React?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Manages component state"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Memoizes a callback function"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Fetches data from an API"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Creates a ref to a DOM element"}},
						}},
						{Key: "correct_answer", Value: "opt2"},
						{Key: "explanation", Value: "useCallback returns a memoized callback that only changes if one of the dependencies has changed."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "hard"},
						{Key: "tags", Value: bson.A{"React", "Hooks", "performance"}},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 30},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 4100},
				{Key: "created_at", Value: mustParseTime("2025-11-16T14:25:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-11-16T14:25:00Z")},
			},
		},
	}
}

// materialSummarySeeds retorna los seeds de la colección material_summary
func materialSummarySeeds() seedDocument {
	return seedDocument{
		collection: "material_summary",
		documents: []interface{}{
			// Summary 1 - POO Java (español)
			bson.D{
				{Key: "material_id", Value: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d"},
				{Key: "summary", Value: "Este material cubre los fundamentos de la programación orientada a objetos en Java. Se explican conceptos clave como clases, objetos, herencia, polimorfismo y encapsulación con ejemplos prácticos."},
				{Key: "key_points", Value: bson.A{
					"Introducción a POO y sus principios fundamentales",
					"Clases y objetos: definición y uso",
					"Herencia y polimorfismo en Java",
					"Encapsulación y modificadores de acceso",
					"Ejemplos prácticos con código",
				}},
				{Key: "language", Value: "es"},
				{Key: "word_count", Value: 42},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4"},
				{Key: "processing_time_ms", Value: 3500},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 850},
					{Key: "completion_tokens", Value: 180},
					{Key: "total_tokens", Value: 1030},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "source_length", Value: 5420},
					{Key: "has_images", Value: false},
				}},
				{Key: "created_at", Value: mustParseTime("2025-11-15T10:30:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-11-15T10:30:00Z")},
			},
			// Summary 2 - React Hooks (inglés)
			bson.D{
				{Key: "material_id", Value: "f1a2b3c4-d5e6-4f5a-9b8c-7d6e5f4a3b2c"},
				{Key: "summary", Value: "A comprehensive guide to React Hooks covering useState, useEffect, useContext, and custom hooks. Learn how to manage state and side effects in functional components effectively."},
				{Key: "key_points", Value: bson.A{
					"Introduction to React Hooks and their benefits",
					"useState for state management",
					"useEffect for side effects and lifecycle",
					"useContext for global state sharing",
					"Creating custom hooks for reusable logic",
				}},
				{Key: "language", Value: "en"},
				{Key: "word_count", Value: 38},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 2800},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 920},
					{Key: "completion_tokens", Value: 165},
					{Key: "total_tokens", Value: 1085},
				}},
				{Key: "created_at", Value: mustParseTime("2025-11-16T14:20:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-11-16T14:20:00Z")},
			},
			// Summary 3 - Estruturas de dados (portugués)
			bson.D{
				{Key: "material_id", Value: "b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e"},
				{Key: "summary", Value: "Material sobre estruturas de dados fundamentais: arrays, listas encadeadas, pilhas e filas. Inclui análise de complexidade e implementações práticas em Python."},
				{Key: "key_points", Value: bson.A{
					"Arrays e suas operações básicas",
					"Listas encadeadas: simples e duplas",
					"Pilhas (LIFO) e suas aplicações",
					"Filas (FIFO) e variantes",
					"Análise de complexidade temporal e espacial",
				}},
				{Key: "language", Value: "pt"},
				{Key: "word_count", Value: 35},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4o"},
				{Key: "processing_time_ms", Value: 3100},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 780},
					{Key: "completion_tokens", Value: 155},
					{Key: "total_tokens", Value: 935},
				}},
				{Key: "created_at", Value: mustParseTime("2025-11-17T09:45:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-11-17T09:45:00Z")},
			},
		},
	}
}

// mustParseTime parsea una fecha RFC3339 o entra en pánico (solo para seeds)
func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(fmt.Sprintf("invalid time format: %s", s))
	}
	return t
}
