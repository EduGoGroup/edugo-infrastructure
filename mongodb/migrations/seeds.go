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
// IDs alineados con PostgreSQL seeds (008_assessments.sql) via mongo_document_id
func materialAssessmentWorkerSeeds() seedDocument {
	return seedDocument{
		collection: "material_assessment_worker",
		documents: []interface{}{
			// ass001 - Examen Fracciones (mat001)
			bson.D{
				{Key: "_id", Value: mustObjectID("aaaaaa000000000000000001")},
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000001"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q-frac-001"},
						{Key: "question_text", Value: "Cuanto es 1/4 + 1/4?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "1/2"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "2/8"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "1/4"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "2/4"}},
						}},
						{Key: "correct_answer", Value: "A"},
						{Key: "explanation", Value: "1/4 + 1/4 = 2/4 = 1/2"},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-frac-002"},
						{Key: "question_text", Value: "Cuanto es 1/4 + 2/4?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "3/8"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "3/4"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "1/2"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "2/4"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "1/4 + 2/4 = 3/4"},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-frac-003"},
						{Key: "question_text", Value: "Cual fraccion es equivalente a 2/6?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "1/3"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "2/6"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "1/2"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "3/6"}},
						}},
						{Key: "correct_answer", Value: "A"},
						{Key: "explanation", Value: "2/6 simplificado es 1/3"},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "medium"},
					},
					bson.D{
						{Key: "question_id", Value: "q-frac-004"},
						{Key: "question_text", Value: "Cuanto es 1/5 + 1/5?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "1/5"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "2/5"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "2/10"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "1/10"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "1/5 + 1/5 = 2/5"},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-frac-005"},
						{Key: "question_text", Value: "Cuanto es 1/8 + 2/8?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "3/16"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "2/8"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "3/8"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "1/4"}},
						}},
						{Key: "correct_answer", Value: "C"},
						{Key: "explanation", Value: "1/8 + 2/8 = 3/8"},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
				}},
				{Key: "total_questions", Value: 5},
				{Key: "total_points", Value: 100},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "manual"},
				{Key: "created_at", Value: mustParseTime("2026-02-10T10:05:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-10T10:05:00Z")},
			},
			// ass002 - Quiz Sistema Solar (mat002)
			bson.D{
				{Key: "_id", Value: mustObjectID("aaaaaa000000000000000002")},
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000002"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q-solar-001"},
						{Key: "question_text", Value: "Cual es el planeta mas grande del sistema solar?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Saturno"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Jupiter"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Urano"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Neptuno"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "Jupiter es el planeta mas grande del sistema solar."},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-solar-002"},
						{Key: "question_text", Value: "Cual es el planeta mas cercano al Sol?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Venus"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Tierra"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Mercurio"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Marte"}},
						}},
						{Key: "correct_answer", Value: "C"},
						{Key: "explanation", Value: "Mercurio es el planeta mas cercano al Sol."},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-solar-003"},
						{Key: "question_text", Value: "Cuantos planetas hay en el sistema solar?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "7"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "8"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "9"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "10"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "El sistema solar tiene 8 planetas desde que Pluton fue reclasificado."},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-solar-004"},
						{Key: "question_text", Value: "Que planeta es conocido como el planeta rojo?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Venus"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Jupiter"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Marte"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Saturno"}},
						}},
						{Key: "correct_answer", Value: "C"},
						{Key: "explanation", Value: "Marte es conocido como el planeta rojo por el oxido de hierro en su superficie."},
						{Key: "points", Value: 20},
						{Key: "difficulty", Value: "easy"},
					},
				}},
				{Key: "total_questions", Value: 4},
				{Key: "total_points", Value: 80},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "manual"},
				{Key: "created_at", Value: mustParseTime("2026-02-12T11:04:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-12T11:04:00Z")},
			},
			// ass003 - Ejercicio Color y Forma (mat004)
			bson.D{
				{Key: "_id", Value: mustObjectID("aaaaaa000000000000000003")},
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000004"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q-color-001"},
						{Key: "question_text", Value: "Cuales son los colores primarios en pintura?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Rojo, Azul, Amarillo"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Rojo, Verde, Azul"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Amarillo, Verde, Naranja"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Blanco, Negro, Gris"}},
						}},
						{Key: "correct_answer", Value: "A"},
						{Key: "explanation", Value: "Los colores primarios en pintura son rojo, azul y amarillo."},
						{Key: "points", Value: 34},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-color-002"},
						{Key: "question_text", Value: "Que color se obtiene al mezclar azul y amarillo?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Naranja"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Verde"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Morado"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Marron"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "La mezcla de azul y amarillo produce verde."},
						{Key: "points", Value: 33},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-color-003"},
						{Key: "question_text", Value: "Que son los colores complementarios?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Colores que estan juntos en el circulo cromatico"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Colores opuestos en el circulo cromatico"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Colores primarios mezclados"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Colores claros y oscuros"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "Los colores complementarios son los que se encuentran opuestos en el circulo cromatico."},
						{Key: "points", Value: 33},
						{Key: "difficulty", Value: "medium"},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 100},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "manual"},
				{Key: "created_at", Value: mustParseTime("2026-02-14T09:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-14T09:00:00Z")},
			},
			// ass004 - English Grammar Test (mat005)
			bson.D{
				{Key: "_id", Value: mustObjectID("aaaaaa000000000000000004")},
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000005"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q-eng-001"},
						{Key: "question_text", Value: "Choose the correct article: ___ apple a day keeps the doctor away."},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "A"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "An"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "The"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "No article"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "Use 'an' before words that start with a vowel sound."},
						{Key: "points", Value: 25},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-eng-002"},
						{Key: "question_text", Value: "Which is the correct form? She ___ to school every day."},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "go"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "goes"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "going"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "gone"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "Third person singular in present simple requires 'goes'."},
						{Key: "points", Value: 25},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-eng-003"},
						{Key: "question_text", Value: "What is the plural of child?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "childs"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "childes"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "children"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "childrens"}},
						}},
						{Key: "correct_answer", Value: "C"},
						{Key: "explanation", Value: "Child has an irregular plural form: children."},
						{Key: "points", Value: 25},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "q-eng-004"},
						{Key: "question_text", Value: "Choose the correct pronoun: John gave the book to ___."},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "I"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "me"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "my"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "mine"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "After a preposition, use the object pronoun 'me'."},
						{Key: "points", Value: 25},
						{Key: "difficulty", Value: "easy"},
					},
				}},
				{Key: "total_questions", Value: 4},
				{Key: "total_points", Value: 100},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "manual"},
				{Key: "created_at", Value: mustParseTime("2026-02-16T08:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-16T08:00:00Z")},
			},
			// ass005 - Evaluacion Historia Chile (mat003)
			bson.D{
				{Key: "_id", Value: mustObjectID("aaaaaa000000000000000005")},
				{Key: "material_id", Value: "aa100000-0000-0000-0000-000000000003"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "q-hist-001"},
						{Key: "question_text", Value: "En que ano se firmo el Acta de Independencia de Chile?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "1810"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "1818"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "1820"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "1825"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "El Acta de Independencia de Chile se firmo el 12 de febrero de 1818."},
						{Key: "points", Value: 34},
						{Key: "difficulty", Value: "medium"},
					},
					bson.D{
						{Key: "question_id", Value: "q-hist-002"},
						{Key: "question_text", Value: "Quien fue el Director Supremo que firmo la independencia?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Jose Miguel Carrera"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Bernardo O'Higgins"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Manuel Blanco Encalada"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Diego Portales"}},
						}},
						{Key: "correct_answer", Value: "B"},
						{Key: "explanation", Value: "Bernardo O'Higgins fue el Director Supremo que firmo el Acta de Independencia."},
						{Key: "points", Value: 33},
						{Key: "difficulty", Value: "medium"},
					},
					bson.D{
						{Key: "question_id", Value: "q-hist-003"},
						{Key: "question_text", Value: "Que batalla fue decisiva para la independencia de Chile?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "A"}, {Key: "option_text", Value: "Batalla de Rancagua"}},
							bson.D{{Key: "option_id", Value: "B"}, {Key: "option_text", Value: "Batalla de Chacabuco"}},
							bson.D{{Key: "option_id", Value: "C"}, {Key: "option_text", Value: "Batalla de Maipu"}},
							bson.D{{Key: "option_id", Value: "D"}, {Key: "option_text", Value: "Batalla de Ayacucho"}},
						}},
						{Key: "correct_answer", Value: "C"},
						{Key: "explanation", Value: "La Batalla de Maipu (1818) fue decisiva para consolidar la independencia de Chile."},
						{Key: "points", Value: 33},
						{Key: "difficulty", Value: "hard"},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 100},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "manual"},
				{Key: "created_at", Value: mustParseTime("2026-02-15T14:06:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-15T14:06:00Z")},
			},
			// ass006 - Proyecto Final Escultura (sin material, borrador)
			bson.D{
				{Key: "_id", Value: mustObjectID("aaaaaa000000000000000006")},
				{Key: "material_id", Value: nil},
				{Key: "questions", Value: bson.A{}},
				{Key: "total_questions", Value: 0},
				{Key: "total_points", Value: 0},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "manual"},
				{Key: "created_at", Value: mustParseTime("2026-02-18T10:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2026-02-18T10:00:00Z")},
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

// mustObjectID convierte un hex string de 24 caracteres a bson.ObjectID o entra en panico (solo para seeds)
func mustObjectID(hex string) bson.ObjectID {
	oid, err := bson.ObjectIDFromHex(hex)
	if err != nil {
		panic(fmt.Sprintf("invalid ObjectID hex: %s", hex))
	}
	return oid
}

// mustParseTime parsea una fecha RFC3339 o entra en panico (solo para seeds)
func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(fmt.Sprintf("invalid time format: %s", s))
	}
	return t
}
