package migrations

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// seedDocument representa un conjunto de documentos para insertar en una coleccion
type seedDocument struct {
	collection string
	documents  []interface{}
}

// applySeedsInternal ejecuta todos los seeds en la base de datos MongoDB
// Retorna el numero de documentos insertados y cualquier error
func applySeedsInternal(ctx context.Context, db *mongo.Database) (int, error) {
	seeds := getSeedDocuments()
	totalInserted := 0

	for _, seed := range seeds {
		if len(seed.documents) == 0 {
			continue
		}

		collection := db.Collection(seed.collection)

		// Usar ordered: false para que si un documento ya existe, continue con los demas
		opts := options.InsertMany().SetOrdered(false)

		result, err := collection.InsertMany(ctx, seed.documents, opts)
		if err != nil {
			// Si es error de duplicados, solo reportamos cuantos se insertaron
			if mongo.IsDuplicateKeyError(err) {
				inserted := 0
				if result != nil {
					inserted = len(result.InsertedIDs)
				}
				totalInserted += inserted
				continue
			}
			return totalInserted, fmt.Errorf("error insertando seeds en %s: %w", seed.collection, err)
		}

		totalInserted += len(result.InsertedIDs)
	}

	return totalInserted, nil
}

// getSeedDocuments retorna todos los seeds organizados por coleccion
// IDs alineados con PostgreSQL v2 seeds (007_materials.sql):
//   - mat001: aa100000-0000-0000-0000-000000000001 (Fracciones)
//   - mat002: aa100000-0000-0000-0000-000000000002 (Sistema Solar)
//   - mat003: aa100000-0000-0000-0000-000000000003 (Historia Chile)
//   - mat004: aa100000-0000-0000-0000-000000000004 (Teoria del Color)
//   - mat005: aa100000-0000-0000-0000-000000000005 (English Grammar)
func getSeedDocuments() []seedDocument {
	return []seedDocument{
		materialAssessmentWorkerSeeds(),
		materialSummarySeeds(),
	}
}

// materialAssessmentWorkerSeeds retorna los seeds de la coleccion material_assessment_worker
func materialAssessmentWorkerSeeds() seedDocument {
	return seedDocument{
		collection: "material_assessment_worker",
		documents: []interface{}{
			// Worker 1 - Fracciones (mat001)
			bson.D{
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000001"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q1111111-1111-1111-1111-111111111111"},
						{Key: "question_text", Value: "Cual es la fraccion equivalente a 2/4?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "3/6"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "1/2"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "4/8"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Todas las anteriores"}},
						}},
						{Key: "correct_answer", Value: "opt4"},
						{Key: "explanation", Value: "Todas son equivalentes a 2/4 ya que al simplificar dan 1/2."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
						{Key: "tags", Value: bson.A{"fracciones", "equivalentes"}},
					},
					bson.D{
						{Key: "question_id", Value: "q2222222-2222-2222-2222-222222222222"},
						{Key: "question_text", Value: "3/5 + 1/5 es igual a 4/10?"},
						{Key: "question_type", Value: "true_false"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "true"}, {Key: "option_text", Value: "Verdadero"}},
							bson.D{{Key: "option_id", Value: "false"}, {Key: "option_text", Value: "Falso"}},
						}},
						{Key: "correct_answer", Value: "false"},
						{Key: "explanation", Value: "3/5 + 1/5 = 4/5, no 4/10. Cuando los denominadores son iguales, se suman los numeradores."},
						{Key: "points", Value: 5},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q3333333-3333-3333-3333-333333333333"},
						{Key: "question_text", Value: "Explica como se multiplican dos fracciones y da un ejemplo."},
						{Key: "question_type", Value: "open"},
						{Key: "options", Value: bson.A{}},
						{Key: "correct_answer", Value: "Se multiplican los numeradores entre si y los denominadores entre si. Ejemplo: 2/3 x 4/5 = 8/15."},
						{Key: "explanation", Value: "La multiplicacion de fracciones es directa: numerador por numerador y denominador por denominador."},
						{Key: "points", Value: 15},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"fracciones", "operaciones"}},
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
				{Key: "created_at", Value: mustParseTime("2026-02-10T10:05:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-10T10:05:00Z")},
			},
			// Worker 2 - Sistema Solar (mat002)
			bson.D{
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000002"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q4444444-4444-4444-4444-444444444444"},
						{Key: "question_text", Value: "Cual es el planeta mas grande del Sistema Solar?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Saturno"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Jupiter"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Urano"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Neptuno"}},
						}},
						{Key: "correct_answer", Value: "opt2"},
						{Key: "explanation", Value: "Jupiter es el planeta mas grande del Sistema Solar, con un diametro de aproximadamente 139.820 km."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
						{Key: "tags", Value: bson.A{"planetas", "sistema solar"}},
					},
					bson.D{
						{Key: "question_id", Value: "q5555555-5555-5555-5555-555555555555"},
						{Key: "question_text", Value: "La Tierra es el tercer planeta mas cercano al Sol?"},
						{Key: "question_type", Value: "true_false"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "true"}, {Key: "option_text", Value: "Verdadero"}},
							bson.D{{Key: "option_id", Value: "false"}, {Key: "option_text", Value: "Falso"}},
						}},
						{Key: "correct_answer", Value: "true"},
						{Key: "explanation", Value: "La Tierra es el tercer planeta en distancia al Sol, despues de Mercurio y Venus."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q6666666-6666-4666-8666-666666666666"},
						{Key: "question_text", Value: "Nombra los 4 planetas rocosos del Sistema Solar y explica por que se llaman asi."},
						{Key: "question_type", Value: "open"},
						{Key: "options", Value: bson.A{}},
						{Key: "correct_answer", Value: "Mercurio, Venus, Tierra y Marte. Se llaman rocosos porque tienen una superficie solida compuesta principalmente de silicatos y metales."},
						{Key: "explanation", Value: "Los planetas rocosos o terrestres se distinguen de los gaseosos por tener una corteza solida."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"planetas", "clasificacion"}},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 30},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 4100},
				{Key: "created_at", Value: mustParseTime("2026-02-12T11:04:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-12T11:04:00Z")},
			},
		},
	}
}

// materialSummarySeeds retorna los seeds de la coleccion material_summary
func materialSummarySeeds() seedDocument {
	return seedDocument{
		collection: "material_summary",
		documents: []interface{}{
			// Summary 1 - Fracciones (mat001)
			bson.D{
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000001"},
				{Key: "summary", Value: "Material introductorio sobre fracciones simples, equivalentes y operaciones basicas. Cubre suma, resta, multiplicacion y division de fracciones con ejemplos practicos."},
				{Key: "key_points", Value: bson.A{
					"Concepto de fraccion: numerador y denominador",
					"Fracciones equivalentes y simplificacion",
					"Suma y resta de fracciones con igual y distinto denominador",
					"Multiplicacion y division de fracciones",
					"Problemas de aplicacion con fracciones",
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
				{Key: "created_at", Value: mustParseTime("2026-02-10T10:05:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-10T10:05:00Z")},
			},
			// Summary 2 - Sistema Solar (mat002)
			bson.D{
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000002"},
				{Key: "summary", Value: "Descripcion completa de los planetas del Sistema Solar, el Sol y sus caracteristicas principales. Incluye datos como tamano, distancia al Sol y composicion."},
				{Key: "key_points", Value: bson.A{
					"El Sol como estrella central del sistema",
					"Planetas rocosos: Mercurio, Venus, Tierra, Marte",
					"Planetas gaseosos: Jupiter, Saturno, Urano, Neptuno",
					"Cinturon de asteroides y otros cuerpos",
					"Comparacion de tamanos y distancias",
				}},
				{Key: "language", Value: "es"},
				{Key: "word_count", Value: 38},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 2800},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 920},
					{Key: "completion_tokens", Value: 165},
					{Key: "total_tokens", Value: 1085},
				}},
				{Key: "created_at", Value: mustParseTime("2026-02-12T11:04:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-12T11:04:00Z")},
			},
			// Summary 3 - Historia de Chile (mat003)
			bson.D{
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000003"},
				{Key: "summary", Value: "Resumen de los principales procesos de la independencia de Chile, desde la Primera Junta de Gobierno hasta la consolidacion de la republica."},
				{Key: "key_points", Value: bson.A{
					"Contexto historico: invasion napoleonica a Espana",
					"Primera Junta de Gobierno (1810)",
					"Patria Vieja y Reconquista",
					"Batalla de Chacabuco y Maipu",
					"Proclamacion de la independencia",
				}},
				{Key: "language", Value: "es"},
				{Key: "word_count", Value: 35},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4o"},
				{Key: "processing_time_ms", Value: 3100},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 780},
					{Key: "completion_tokens", Value: 155},
					{Key: "total_tokens", Value: 935},
				}},
				{Key: "created_at", Value: mustParseTime("2026-02-15T14:06:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-15T14:06:00Z")},
			},
		},
	}
}

// mustParseTime parsea una fecha RFC3339 o entra en panico (solo para seeds)
func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(fmt.Sprintf("invalid time format: %s", s))
	}
	return t
}
