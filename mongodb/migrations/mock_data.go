package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// mockDocument representa un conjunto de documentos mock para insertar en una colección
type mockDocument struct {
	collection string
	documents  []interface{}
}

// applyMockDataInternal ejecuta todos los mocks en la base de datos MongoDB
// Retorna el número de documentos insertados y cualquier error
func applyMockDataInternal(ctx context.Context, db *mongo.Database) (int, error) {
	mocks := getMockDocuments()
	totalInserted := 0

	for _, mock := range mocks {
		if len(mock.documents) == 0 {
			continue
		}

		collection := db.Collection(mock.collection)

		// Usar ordered: false para que si un documento ya existe, continúe con los demás
		opts := options.InsertMany().SetOrdered(false)

		result, err := collection.InsertMany(ctx, mock.documents, opts)
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
			return totalInserted, fmt.Errorf("error insertando mocks en %s: %w", mock.collection, err)
		}

		totalInserted += len(result.InsertedIDs)
	}

	return totalInserted, nil
}

// getMockDocuments retorna todos los mocks organizados por colección
func getMockDocuments() []mockDocument {
	return []mockDocument{
		materialAssessmentWorkerMockData(),
		materialSummaryMockData(),
		// Las siguientes colecciones fueron eliminadas por no uso:
		// - analyticsEventsMockData (analytics usa servicio externo)
		// - assessmentsMockData (duplicada por material_assessment_worker)
		// - auditLogsMockData (usará SaaS externo)
		// - notificationsMockData (push notifications no implementado)
	}
}

// materialAssessmentWorkerMockData retorna los mocks de la colección material_assessment_worker
func materialAssessmentWorkerMockData() mockDocument {
	return mockDocument{
		collection: "material_assessment_worker",
		documents: []interface{}{
			// Worker 1 - Química
			bson.D{
				{Key: "material_id", Value: "c1e2d3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "qchem001-0001-0001-0001-000000000001"},
						{Key: "question_text", Value: "¿Cuál es la configuración electrónica del oxígeno (Z=8)?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "1s² 2s² 2p⁴"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "1s² 2s² 2p⁶"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "1s² 2s² 2p³"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "1s² 2s² 2p⁵"}},
						}},
						{Key: "correct_answer", Value: "opt1"},
						{Key: "explanation", Value: "El oxígeno tiene 8 electrones distribuidos: 2 en 1s, 2 en 2s y 4 en 2p."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"química", "configuración electrónica"}},
					},
					bson.D{
						{Key: "question_id", Value: "qchem002-0002-0002-0002-000000000002"},
						{Key: "question_text", Value: "¿Qué tipo de reacción es: 2H₂ + O₂ → 2H₂O?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Descomposición"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Síntesis"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Sustitución"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Doble sustitución"}},
						}},
						{Key: "correct_answer", Value: "opt2"},
						{Key: "explanation", Value: "Es una reacción de síntesis o combinación donde dos sustancias simples forman una compuesta."},
						{Key: "points", Value: 8},
						{Key: "difficulty", Value: "easy"},
					},
					bson.D{
						{Key: "question_id", Value: "qchem003-0003-4003-8003-000000000003"},
						{Key: "question_text", Value: "¿Cuántos electrones puede alojar el tercer nivel de energía?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "8"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "18"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "32"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "2"}},
						}},
						{Key: "correct_answer", Value: "opt2"},
						{Key: "explanation", Value: "El tercer nivel (n=3) puede alojar hasta 18 electrones: 2 en 3s, 6 en 3p y 10 en 3d."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "hard"},
						{Key: "tags", Value: bson.A{"química", "electrones"}},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 28},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4o"},
				{Key: "processing_time_ms", Value: 4800},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 1100},
					{Key: "completion_tokens", Value: 380},
					{Key: "total_tokens", Value: 1480},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T10:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T10:00:00Z")},
			},
			// Worker 2 - Historia (inglés)
			bson.D{
				{Key: "material_id", Value: "e1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "qhist001-0001-0001-0001-000000000001"},
						{Key: "question_text", Value: "Which country was NOT part of the Allied Powers in WWII?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "United States"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Soviet Union"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Italy"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "United Kingdom"}},
						}},
						{Key: "correct_answer", Value: "opt3"},
						{Key: "explanation", Value: "Italy was part of the Axis Powers alongside Germany and Japan."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
						{Key: "tags", Value: bson.A{"history", "WWII"}},
					},
					bson.D{
						{Key: "question_id", Value: "qhist002-0002-4002-8002-000000000002"},
						{Key: "question_text", Value: "In which year did WWII end?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "1943"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "1944"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "1945"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "1946"}},
						}},
						{Key: "correct_answer", Value: "opt3"},
						{Key: "explanation", Value: "WWII ended in 1945 with Germany surrendering in May and Japan in September."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
						{Key: "tags", Value: bson.A{"history", "WWII"}},
					},
					bson.D{
						{Key: "question_id", Value: "qhist003-0003-4003-8003-000000000003"},
						{Key: "question_text", Value: "What was the codename for the Allied invasion of Normandy?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Operation Overlord"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Operation Barbarossa"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Operation Market Garden"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Operation Torch"}},
						}},
						{Key: "correct_answer", Value: "opt1"},
						{Key: "explanation", Value: "Operation Overlord was the codename for the D-Day invasion of Normandy on June 6, 1944."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"history", "WWII", "D-Day"}},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 30},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 3200},
				{Key: "created_at", Value: mustParseTime("2025-01-20T11:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T11:00:00Z")},
			},
			// Worker 3 - Biología (portugués)
			bson.D{
				{Key: "material_id", Value: "b1c2d3e4-f5a6-4b7c-8d9e-0f1a2b3c4d5e"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_id", Value: "qbio001-0001-0001-0001-000000000001"},
						{Key: "question_text", Value: "Qual é a função principal da mitocôndria?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Síntese de proteínas"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Produção de energia (ATP)"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Digestão celular"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Fotossíntese"}},
						}},
						{Key: "correct_answer", Value: "opt2"},
						{Key: "explanation", Value: "A mitocôndria é responsável pela respiração celular e produção de ATP."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "medium"},
						{Key: "tags", Value: bson.A{"biologia", "célula"}},
					},
					bson.D{
						{Key: "question_id", Value: "qbio002-0002-0002-0002-000000000002"},
						{Key: "question_text", Value: "O que é a fotossíntese?"},
						{Key: "question_type", Value: "open"},
						{Key: "options", Value: bson.A{}},
						{Key: "correct_answer", Value: "Processo pelo qual plantas convertem luz solar, água e CO₂ em glicose e oxigênio."},
						{Key: "explanation", Value: "A fotossíntese ocorre nos cloroplastos e é fundamental para a vida na Terra."},
						{Key: "points", Value: 15},
						{Key: "difficulty", Value: "hard"},
					},
					bson.D{
						{Key: "question_id", Value: "qbio003-0003-4003-8003-000000000003"},
						{Key: "question_text", Value: "Qual é a função do núcleo celular?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_id", Value: "opt1"}, {Key: "option_text", Value: "Produção de energia"}},
							bson.D{{Key: "option_id", Value: "opt2"}, {Key: "option_text", Value: "Digestão de substâncias"}},
							bson.D{{Key: "option_id", Value: "opt3"}, {Key: "option_text", Value: "Controle das atividades celulares e armazenamento do DNA"}},
							bson.D{{Key: "option_id", Value: "opt4"}, {Key: "option_text", Value: "Síntese de lipídios"}},
						}},
						{Key: "correct_answer", Value: "opt3"},
						{Key: "explanation", Value: "O núcleo controla todas as atividades celulares e armazena o material genético (DNA)."},
						{Key: "points", Value: 10},
						{Key: "difficulty", Value: "easy"},
						{Key: "tags", Value: bson.A{"biologia", "núcleo"}},
					},
				}},
				{Key: "total_questions", Value: 3},
				{Key: "total_points", Value: 35},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4"},
				{Key: "processing_time_ms", Value: 5500},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 950},
					{Key: "completion_tokens", Value: 420},
					{Key: "total_tokens", Value: 1370},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T12:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T12:00:00Z")},
			},
		},
	}
}

// materialSummaryMockData retorna los mocks de la colección material_summary
func materialSummaryMockData() mockDocument {
	return mockDocument{
		collection: "material_summary",
		documents: []interface{}{
			// Summary 1 - Química (español)
			bson.D{
				{Key: "material_id", Value: "c1e2d3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f"},
				{Key: "summary", Value: "Este material explora conceptos fundamentales de química general incluyendo estructura atómica, tabla periódica, enlaces químicos y reacciones. Se enfatiza la configuración electrónica y los diferentes tipos de reacciones químicas con ejemplos aplicados."},
				{Key: "key_points", Value: bson.A{
					"Estructura atómica y configuración electrónica",
					"Organización de la tabla periódica",
					"Tipos de enlaces: iónico, covalente, metálico",
					"Clasificación de reacciones químicas",
					"Balanceo de ecuaciones químicas",
				}},
				{Key: "language", Value: "es"},
				{Key: "word_count", Value: 48},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4o"},
				{Key: "processing_time_ms", Value: 3800},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 920},
					{Key: "completion_tokens", Value: 195},
					{Key: "total_tokens", Value: 1115},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "source_length", Value: 6200},
					{Key: "has_images", Value: true},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T10:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T10:00:00Z")},
			},
			// Summary 2 - Historia (inglés)
			bson.D{
				{Key: "material_id", Value: "e1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c"},
				{Key: "summary", Value: "Comprehensive overview of World War II covering major events, key battles, political alliances, and the war's global impact. Includes analysis of the causes, major turning points, and the Holocaust, ending with post-war consequences."},
				{Key: "key_points", Value: bson.A{
					"Causes of World War II and the rise of totalitarianism",
					"Major battles: Pearl Harbor, D-Day, Stalingrad",
					"Allied and Axis Powers composition",
					"The Holocaust and war crimes",
					"Post-war reconstruction and the United Nations",
				}},
				{Key: "language", Value: "en"},
				{Key: "word_count", Value: 52},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 3100},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 1050},
					{Key: "completion_tokens", Value: 210},
					{Key: "total_tokens", Value: 1260},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T11:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T11:00:00Z")},
			},
			// Summary 3 - Biología (portugués)
			bson.D{
				{Key: "material_id", Value: "b1c2d3e4-f5a6-4b7c-8d9e-0f1a2b3c4d5e"},
				{Key: "summary", Value: "Material sobre biologia celular abordando estrutura e função das organelas, processos de respiração celular e fotossíntese. Apresenta detalhes sobre mitocôndrias, cloroplastos, núcleo e outras estruturas celulares essenciais."},
				{Key: "key_points", Value: bson.A{
					"Estrutura básica da célula eucariota",
					"Função das organelas celulares",
					"Respiração celular e produção de ATP",
					"Fotossíntese em plantas",
					"Membrana plasmática e transporte celular",
				}},
				{Key: "language", Value: "pt"},
				{Key: "word_count", Value: 44},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4"},
				{Key: "processing_time_ms", Value: 3500},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 880},
					{Key: "completion_tokens", Value: 185},
					{Key: "total_tokens", Value: 1065},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T12:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T12:00:00Z")},
			},
			// Summary 4 - Cálculo (francés)
			bson.D{
				{Key: "material_id", Value: "ca1b2c3d-4e5f-4a6b-8c7d-8e9f0a1b2c3d"},
				{Key: "summary", Value: "Introduction au calcul différentiel et intégral couvrant les concepts de limites, dérivées et intégrales. Présentation de techniques de résolution et applications pratiques dans divers domaines scientifiques."},
				{Key: "key_points", Value: bson.A{
					"Notion de limite et continuité",
					"Dérivées et règles de dérivation",
					"Applications des dérivées",
					"Intégrales définies et indéfinies",
					"Théorème fondamental du calcul",
				}},
				{Key: "language", Value: "pt"},
				{Key: "word_count", Value: 40},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 3300},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 900},
					{Key: "completion_tokens", Value: 175},
					{Key: "total_tokens", Value: 1075},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T13:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T13:00:00Z")},
			},
			// Summary 5 - Física (alemán)
			bson.D{
				{Key: "material_id", Value: "f1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5e"},
				{Key: "summary", Value: "Grundlagen der klassischen Mechanik mit Schwerpunkt auf Newtonsche Gesetze, Bewegung, Kraft und Energie. Enthält theoretische Konzepte sowie praktische Beispiele und Problemlösungen."},
				{Key: "key_points", Value: bson.A{
					"Newtonsche Bewegungsgesetze",
					"Kinematik und Dynamik",
					"Arbeit, Energie und Leistung",
					"Impuls und Impulserhaltung",
					"Anwendungen in der Praxis",
				}},
				{Key: "language", Value: "en"},
				{Key: "word_count", Value: 38},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4"},
				{Key: "processing_time_ms", Value: 3200},
				{Key: "token_usage", Value: bson.D{
					{Key: "prompt_tokens", Value: 870},
					{Key: "completion_tokens", Value: 170},
					{Key: "total_tokens", Value: 1040},
				}},
				{Key: "created_at", Value: mustParseTime("2025-01-20T14:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T14:00:00Z")},
			},
		},
	}
}
