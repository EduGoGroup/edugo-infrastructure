package migrations

import (
	"context"
	"fmt"
	"time"

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

// analyticsEventsMockData retorna los mocks de la colección analytics_events
func analyticsEventsMockData() mockDocument {
	return mockDocument{
		collection: "analytics_events",
		documents: []interface{}{
			// Event 1 - Mobile login
			bson.D{
				{Key: "event_name", Value: "user.login"},
				{Key: "user_id", Value: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
				{Key: "session_id", Value: "sess_mobile_login_001"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T08:00:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "login_method", Value: "google_oauth"},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "mobile"},
					{Key: "os", Value: "iOS"},
					{Key: "os_version", Value: "17.2"},
					{Key: "device_type", Value: "mobile"},
					{Key: "screen_resolution", Value: "1170x2532"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "AR"},
					{Key: "city", Value: "Buenos Aires"},
					{Key: "timezone", Value: "America/Argentina/Buenos_Aires"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"},
					{Key: "user_role", Value: "teacher"},
				}},
			},
			// Event 2 - Tablet material view
			bson.D{
				{Key: "event_name", Value: "material.view"},
				{Key: "user_id", Value: "cccccccc-cccc-cccc-cccc-cccccccccccc"},
				{Key: "session_id", Value: "sess_tablet_view_002"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T09:15:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "dddddddd-dddd-dddd-dddd-dddddddddddd"},
					{Key: "resource_type", Value: "material"},
					{Key: "custom_data", Value: bson.D{
						{Key: "material_title", Value: "Historia Universal - Segunda Guerra Mundial"},
						{Key: "subject", Value: "Historia"},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "tablet"},
					{Key: "os", Value: "iPadOS"},
					{Key: "os_version", Value: "16.5"},
					{Key: "device_type", Value: "tablet"},
					{Key: "screen_resolution", Value: "2048x2732"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "MX"},
					{Key: "city", Value: "Ciudad de México"},
					{Key: "timezone", Value: "America/Mexico_City"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 3 - Web assessment start
			bson.D{
				{Key: "event_name", Value: "assessment.start"},
				{Key: "user_id", Value: "ffffffff-ffff-ffff-ffff-ffffffffffff"},
				{Key: "session_id", Value: "sess_web_assess_003"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T10:30:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "gggggggg-gggg-gggg-gggg-gggggggggggg"},
					{Key: "resource_type", Value: "assessment"},
					{Key: "custom_data", Value: bson.D{
						{Key: "questions_count", Value: 3},
						{Key: "subject", Value: "Química"},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "Windows"},
					{Key: "os_version", Value: "11"},
					{Key: "browser", Value: "Edge"},
					{Key: "browser_version", Value: "119.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "1920x1080"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CO"},
					{Key: "city", Value: "Bogotá"},
					{Key: "timezone", Value: "America/Bogota"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 4 - Mobile assessment complete
			bson.D{
				{Key: "event_name", Value: "assessment.complete"},
				{Key: "user_id", Value: "iiiiiiii-iiii-iiii-iiii-iiiiiiiiiiii"},
				{Key: "session_id", Value: "sess_mobile_complete_004"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T11:45:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "jjjjjjjj-jjjj-jjjj-jjjj-jjjjjjjjjjjj"},
					{Key: "resource_type", Value: "assessment"},
					{Key: "duration_seconds", Value: 240},
					{Key: "custom_data", Value: bson.D{
						{Key: "score", Value: 85},
						{Key: "questions_count", Value: 3},
						{Key: "correct_answers", Value: 2},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "android"},
					{Key: "os", Value: "Android"},
					{Key: "os_version", Value: "14"},
					{Key: "device_type", Value: "mobile"},
					{Key: "screen_resolution", Value: "1440x3200"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "PE"},
					{Key: "city", Value: "Lima"},
					{Key: "timezone", Value: "America/Lima"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "kkkkkkkk-kkkk-kkkk-kkkk-kkkkkkkkkkkk"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 5 - Tablet search
			bson.D{
				{Key: "event_name", Value: "search.performed"},
				{Key: "user_id", Value: "llllllll-llll-llll-llll-llllllllllll"},
				{Key: "session_id", Value: "sess_tablet_search_005"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T13:00:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "search_query", Value: "biología celular"},
					{Key: "search_results_count", Value: 5},
					{Key: "custom_data", Value: bson.D{
						{Key: "filters_applied", Value: bson.D{
							{Key: "subject", Value: "Biología"},
						}},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "tablet"},
					{Key: "os", Value: "Android"},
					{Key: "os_version", Value: "13"},
					{Key: "device_type", Value: "tablet"},
					{Key: "screen_resolution", Value: "2560x1600"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "UY"},
					{Key: "city", Value: "Montevideo"},
					{Key: "timezone", Value: "America/Montevideo"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "mmmmmmmm-mmmm-mmmm-mmmm-mmmmmmmmmmmm"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 6 - Web page view
			bson.D{
				{Key: "event_name", Value: "page.view"},
				{Key: "user_id", Value: "nnnnnnnn-nnnn-nnnn-nnnn-nnnnnnnnnnnn"},
				{Key: "session_id", Value: "sess_web_page_006"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T14:20:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "page_path", Value: "/dashboard"},
					{Key: "page_title", Value: "Dashboard Principal"},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "Linux"},
					{Key: "os_version", Value: "Ubuntu 22.04"},
					{Key: "browser", Value: "Firefox"},
					{Key: "browser_version", Value: "121.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "2560x1440"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "EC"},
					{Key: "city", Value: "Quito"},
					{Key: "timezone", Value: "America/Guayaquil"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "oooooooo-oooo-oooo-oooo-oooooooooooo"},
					{Key: "user_role", Value: "admin"},
				}},
			},
			// Event 7 - Mobile material download
			bson.D{
				{Key: "event_name", Value: "material.download"},
				{Key: "user_id", Value: "pppppppp-pppp-pppp-pppp-pppppppppppp"},
				{Key: "session_id", Value: "sess_mobile_download_007"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T15:30:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "qqqqqqqq-qqqq-qqqq-qqqq-qqqqqqqqqqqq"},
					{Key: "resource_type", Value: "material"},
					{Key: "custom_data", Value: bson.D{
						{Key: "material_title", Value: "Geografía - Continentes"},
						{Key: "subject", Value: "Geografía"},
						{Key: "file_size_mb", Value: 8.5},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "mobile"},
					{Key: "os", Value: "iOS"},
					{Key: "os_version", Value: "17.3"},
					{Key: "device_type", Value: "mobile"},
					{Key: "screen_resolution", Value: "1284x2778"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "VE"},
					{Key: "city", Value: "Caracas"},
					{Key: "timezone", Value: "America/Caracas"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "rrrrrrrr-rrrr-rrrr-rrrr-rrrrrrrrrrrr"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 8 - Web video playback
			bson.D{
				{Key: "event_name", Value: "video.play"},
				{Key: "user_id", Value: "ssssssss-ssss-ssss-ssss-ssssssssssss"},
				{Key: "session_id", Value: "sess_web_video_008"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T16:00:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "tttttttt-tttt-tttt-tttt-tttttttttttt"},
					{Key: "resource_type", Value: "video"},
					{Key: "custom_data", Value: bson.D{
						{Key: "video_title", Value: "Tutorial: Integrales Definidas"},
						{Key: "duration_seconds", Value: 320},
						{Key: "subject", Value: "Cálculo"},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "macOS"},
					{Key: "os_version", Value: "14.2"},
					{Key: "browser", Value: "Safari"},
					{Key: "browser_version", Value: "17.2"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "2880x1800"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Concepción"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "uuuuuuuu-uuuu-uuuu-uuuu-uuuuuuuuuuuu"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 9 - Tablet session end
			bson.D{
				{Key: "event_name", Value: "session.end"},
				{Key: "user_id", Value: "vvvvvvvv-vvvv-vvvv-vvvv-vvvvvvvvvvvv"},
				{Key: "session_id", Value: "sess_tablet_end_009"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T17:00:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "session_duration_seconds", Value: 1800},
					{Key: "custom_data", Value: bson.D{
						{Key: "pages_viewed", Value: 12},
						{Key: "materials_accessed", Value: 3},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "tablet"},
					{Key: "os", Value: "iPadOS"},
					{Key: "os_version", Value: "17.0"},
					{Key: "device_type", Value: "tablet"},
					{Key: "screen_resolution", Value: "2388x1668"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "BO"},
					{Key: "city", Value: "La Paz"},
					{Key: "timezone", Value: "America/La_Paz"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "wwwwwwww-wwww-wwww-wwww-wwwwwwwwwwww"},
					{Key: "user_role", Value: "teacher"},
				}},
			},
			// Event 10 - Mobile notification click
			bson.D{
				{Key: "event_name", Value: "notification.click"},
				{Key: "user_id", Value: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"},
				{Key: "session_id", Value: "sess_mobile_notif_010"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T18:00:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "notification_type", Value: "assessment.graded"},
					{Key: "custom_data", Value: bson.D{
						{Key: "notification_id", Value: "notif_mock_123"},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "mobile"},
					{Key: "os", Value: "Android"},
					{Key: "os_version", Value: "14"},
					{Key: "device_type", Value: "mobile"},
					{Key: "screen_resolution", Value: "1080x2400"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "PY"},
					{Key: "city", Value: "Asunción"},
					{Key: "timezone", Value: "America/Asuncion"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"},
					{Key: "user_role", Value: "student"},
				}},
			},
		},
	}
}

// assessmentsMockData retorna los mocks de la colección material_assessment
func assessmentsMockData() mockDocument {
	return mockDocument{
		collection: "material_assessment",
		documents: []interface{}{
			// Assessment 1 (Química - Difícil)
			bson.D{
				{Key: "_id", Value: mustParseObjectID("607f1f77bcf86cd799439021")},
				{Key: "material_id", Value: "gggggggg-gggg-gggg-gggg-gggggggggggg"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_index", Value: 0},
						{Key: "question_text", Value: "¿Cuál es el número de oxidación del Cr en K₂Cr₂O₇?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_index", Value: 0}, {Key: "text", Value: "+3"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 1}, {Key: "text", Value: "+6"}, {Key: "is_correct", Value: true}},
							bson.D{{Key: "option_index", Value: 2}, {Key: "text", Value: "+7"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 3}, {Key: "text", Value: "+2"}, {Key: "is_correct", Value: false}},
						}},
					},
					bson.D{
						{Key: "question_index", Value: 1},
						{Key: "question_text", Value: "¿Qué tipo de enlace forma el NaCl?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_index", Value: 0}, {Key: "text", Value: "Covalente"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 1}, {Key: "text", Value: "Iónico"}, {Key: "is_correct", Value: true}},
							bson.D{{Key: "option_index", Value: 2}, {Key: "text", Value: "Metálico"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 3}, {Key: "text", Value: "Van der Waals"}, {Key: "is_correct", Value: false}},
						}},
					},
					bson.D{
						{Key: "question_index", Value: 2},
						{Key: "question_text", Value: "¿Cuál es el pH de una solución con [H⁺] = 1x10⁻⁵ M?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_index", Value: 0}, {Key: "text", Value: "3"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 1}, {Key: "text", Value: "5"}, {Key: "is_correct", Value: true}},
							bson.D{{Key: "option_index", Value: 2}, {Key: "text", Value: "7"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 3}, {Key: "text", Value: "9"}, {Key: "is_correct", Value: false}},
						}},
					},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "subject", Value: "Química"},
					{Key: "grade", Value: "11th"},
					{Key: "difficulty", Value: "hard"},
				}},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
			// Assessment 2 (Historia - Fácil)
			bson.D{
				{Key: "_id", Value: mustParseObjectID("607f1f77bcf86cd799439022")},
				{Key: "material_id", Value: "dddddddd-dddd-dddd-dddd-dddddddddddd"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_index", Value: 0},
						{Key: "question_text", Value: "¿En qué año comenzó la Segunda Guerra Mundial?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_index", Value: 0}, {Key: "text", Value: "1914"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 1}, {Key: "text", Value: "1939"}, {Key: "is_correct", Value: true}},
							bson.D{{Key: "option_index", Value: 2}, {Key: "text", Value: "1941"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 3}, {Key: "text", Value: "1945"}, {Key: "is_correct", Value: false}},
						}},
					},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "subject", Value: "Historia"},
					{Key: "grade", Value: "9th"},
					{Key: "difficulty", Value: "easy"},
				}},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
			// Assessment 3 (Cálculo - Medio)
			bson.D{
				{Key: "_id", Value: mustParseObjectID("607f1f77bcf86cd799439023")},
				{Key: "material_id", Value: "tttttttt-tttt-tttt-tttt-tttttttttttt"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_index", Value: 0},
						{Key: "question_text", Value: "¿Cuál es la derivada de x²?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_index", Value: 0}, {Key: "text", Value: "x"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 1}, {Key: "text", Value: "2x"}, {Key: "is_correct", Value: true}},
							bson.D{{Key: "option_index", Value: 2}, {Key: "text", Value: "x²"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 3}, {Key: "text", Value: "2x²"}, {Key: "is_correct", Value: false}},
						}},
					},
					bson.D{
						{Key: "question_index", Value: 1},
						{Key: "question_text", Value: "¿Cuál es la integral de 2x?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{{Key: "option_index", Value: 0}, {Key: "text", Value: "x"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 1}, {Key: "text", Value: "x² + C"}, {Key: "is_correct", Value: true}},
							bson.D{{Key: "option_index", Value: 2}, {Key: "text", Value: "2x² + C"}, {Key: "is_correct", Value: false}},
							bson.D{{Key: "option_index", Value: 3}, {Key: "text", Value: "x²"}, {Key: "is_correct", Value: false}},
						}},
					},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "subject", Value: "Cálculo"},
					{Key: "grade", Value: "12th"},
					{Key: "difficulty", Value: "medium"},
				}},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
		},
	}
}

// auditLogsMockData retorna los mocks de la colección audit_logs
func auditLogsMockData() mockDocument {
	return mockDocument{
		collection: "audit_logs",
		documents: []interface{}{
			// Audit log 1 - Material deleted
			bson.D{
				{Key: "event_type", Value: "material.deleted"},
				{Key: "actor_id", Value: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "material"},
				{Key: "resource_id", Value: "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"},
				{Key: "action", Value: "delete"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "10.0.1.50"},
					{Key: "user_agent", Value: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X)"},
					{Key: "metadata", Value: bson.D{
						{Key: "file_name", Value: "material_obsoleto.pdf"},
						{Key: "reason", Value: "contenido desactualizado"},
					}},
				}},
				{Key: "severity", Value: "warning"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T08:30:00Z")},
				{Key: "session_id", Value: "sess_mobile_delete_001"},
				{Key: "request_id", Value: "req_mock_001"},
			},
			// Audit log 2 - User created
			bson.D{
				{Key: "event_type", Value: "user.created"},
				{Key: "actor_id", Value: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"},
				{Key: "actor_type", Value: "admin"},
				{Key: "resource_type", Value: "user"},
				{Key: "resource_id", Value: "newuser1-1111-1111-1111-111111111111"},
				{Key: "action", Value: "create"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.10.5"},
					{Key: "user_agent", Value: "Mozilla/5.0 (Windows NT 11.0; Win64; x64)"},
					{Key: "changes", Value: bson.D{
						{Key: "email", Value: "newuser@example.com"},
						{Key: "role", Value: "student"},
						{Key: "school_id", Value: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T09:00:00Z")},
				{Key: "session_id", Value: "sess_admin_create_002"},
				{Key: "request_id", Value: "req_mock_002"},
			},
			// Audit log 3 - Password changed
			bson.D{
				{Key: "event_type", Value: "user.password_changed"},
				{Key: "actor_id", Value: "cccccccc-cccc-cccc-cccc-cccccccccccc"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "user"},
				{Key: "resource_id", Value: "cccccccc-cccc-cccc-cccc-cccccccccccc"},
				{Key: "action", Value: "update"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "172.16.0.100"},
					{Key: "user_agent", Value: "Mozilla/5.0 (iPad; CPU OS 16_5 like Mac OS X)"},
					{Key: "metadata", Value: bson.D{
						{Key: "password_strength", Value: "strong"},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T10:00:00Z")},
				{Key: "session_id", Value: "sess_tablet_password_003"},
				{Key: "request_id", Value: "req_mock_003"},
			},
			// Audit log 4 - Multiple failed login attempts
			bson.D{
				{Key: "event_type", Value: "security.brute_force"},
				{Key: "actor_id", Value: "unknown"},
				{Key: "actor_type", Value: "unknown"},
				{Key: "resource_type", Value: "auth"},
				{Key: "action", Value: "login"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "203.0.113.42"},
					{Key: "user_agent", Value: "curl/7.68.0"},
					{Key: "error", Value: bson.D{
						{Key: "code", Value: "RATE_LIMIT_EXCEEDED"},
						{Key: "message", Value: "Too many failed login attempts"},
					}},
					{Key: "metadata", Value: bson.D{
						{Key: "attempts_count", Value: 15},
						{Key: "attempted_emails", Value: bson.A{"admin@test.com", "root@test.com"}},
					}},
				}},
				{Key: "severity", Value: "critical"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T11:00:00Z")},
				{Key: "request_id", Value: "req_mock_004"},
			},
			// Audit log 5 - Assessment published by teacher
			bson.D{
				{Key: "event_type", Value: "assessment.published"},
				{Key: "actor_id", Value: "ffffffff-ffff-ffff-ffff-ffffffffffff"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "assessment"},
				{Key: "resource_id", Value: "gggggggg-gggg-gggg-gggg-gggggggggggg"},
				{Key: "action", Value: "update"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.5.25"},
					{Key: "user_agent", Value: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64)"},
					{Key: "changes", Value: bson.D{
						{Key: "status", Value: bson.D{
							{Key: "from", Value: "draft"},
							{Key: "to", Value: "published"},
						}},
						{Key: "questions_count", Value: 3},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T12:00:00Z")},
				{Key: "session_id", Value: "sess_teacher_publish_005"},
				{Key: "request_id", Value: "req_mock_005"},
			},
			// Audit log 6 - School settings updated
			bson.D{
				{Key: "event_type", Value: "school.settings_updated"},
				{Key: "actor_id", Value: "adminsch-0000-0000-0000-000000000001"},
				{Key: "actor_type", Value: "admin"},
				{Key: "resource_type", Value: "school"},
				{Key: "resource_id", Value: "hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh"},
				{Key: "action", Value: "update"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.1.10"},
					{Key: "user_agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2)"},
					{Key: "changes", Value: bson.D{
						{Key: "max_students_per_class", Value: bson.D{
							{Key: "from", Value: 30},
							{Key: "to", Value: 35},
						}},
						{Key: "timezone", Value: bson.D{
							{Key: "from", Value: "America/Santiago"},
							{Key: "to", Value: "America/Bogota"},
						}},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T13:00:00Z")},
				{Key: "session_id", Value: "sess_admin_settings_006"},
				{Key: "request_id", Value: "req_mock_006"},
			},
			// Audit log 7 - Database backup completed
			bson.D{
				{Key: "event_type", Value: "system.backup_completed"},
				{Key: "actor_id", Value: "system"},
				{Key: "actor_type", Value: "system"},
				{Key: "resource_type", Value: "system"},
				{Key: "action", Value: "backup"},
				{Key: "details", Value: bson.D{
					{Key: "metadata", Value: bson.D{
						{Key: "backup_type", Value: "automated_weekly"},
						{Key: "backup_size_mb", Value: 2048},
						{Key: "backup_location", Value: "s3://edugo-backups/2025-01-20/"},
						{Key: "collections_backed_up", Value: 9},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T03:00:00Z")},
				{Key: "request_id", Value: "req_mock_007"},
			},
			// Audit log 8 - API key rotated
			bson.D{
				{Key: "event_type", Value: "security.api_key_rotated"},
				{Key: "actor_id", Value: "adminsch-0000-0000-0000-000000000002"},
				{Key: "actor_type", Value: "admin"},
				{Key: "resource_type", Value: "api_key"},
				{Key: "resource_id", Value: "apikey_school_001"},
				{Key: "action", Value: "rotate"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "10.0.2.15"},
					{Key: "user_agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edge/119.0"},
					{Key: "metadata", Value: bson.D{
						{Key: "reason", Value: "scheduled_rotation"},
						{Key: "old_key_last_used", Value: "2025-01-15T10:00:00Z"},
					}},
				}},
				{Key: "severity", Value: "warning"},
				{Key: "timestamp", Value: mustParseTime("2025-01-20T14:00:00Z")},
				{Key: "session_id", Value: "sess_admin_apikey_008"},
				{Key: "request_id", Value: "req_mock_008"},
			},
		},
	}
}

// materialAssessmentWorkerMockData retorna los mocks de la colección material_assessment_worker
func materialAssessmentWorkerMockData() mockDocument {
	return mockDocument{
		collection: "material_assessment_worker",
		documents: []interface{}{
			// Worker 1 - Química
			bson.D{
				{Key: "material_id", Value: "chem_mat_001-0001-0001-0001-000000000001"},
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
				}},
				{Key: "total_questions", Value: 2},
				{Key: "total_points", Value: 18},
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
				{Key: "material_id", Value: "hist_mat_001-0001-0001-0001-000000000001"},
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
				}},
				{Key: "total_questions", Value: 1},
				{Key: "total_points", Value: 10},
				{Key: "version", Value: 1},
				{Key: "ai_model", Value: "gpt-4-turbo"},
				{Key: "processing_time_ms", Value: 3200},
				{Key: "created_at", Value: mustParseTime("2025-01-20T11:00:00Z")},
				{Key: "updated_at", Value: mustParseTime("2025-01-20T11:00:00Z")},
			},
			// Worker 3 - Biología (portugués)
			bson.D{
				{Key: "material_id", Value: "bio_mat_001-0001-0001-0001-000000000001"},
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
				}},
				{Key: "total_questions", Value: 2},
				{Key: "total_points", Value: 25},
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
				{Key: "material_id", Value: "chem_mat_001-0001-0001-0001-000000000001"},
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
				{Key: "material_id", Value: "hist_mat_001-0001-0001-0001-000000000001"},
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
				{Key: "material_id", Value: "bio_mat_001-0001-0001-0001-000000000001"},
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
				{Key: "material_id", Value: "calc_mat_001-0001-0001-0001-000000000001"},
				{Key: "summary", Value: "Introduction au calcul différentiel et intégral couvrant les concepts de limites, dérivées et intégrales. Présentation de techniques de résolution et applications pratiques dans divers domaines scientifiques."},
				{Key: "key_points", Value: bson.A{
					"Notion de limite et continuité",
					"Dérivées et règles de dérivation",
					"Applications des dérivées",
					"Intégrales définies et indéfinies",
					"Théorème fondamental du calcul",
				}},
				{Key: "language", Value: "fr"},
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
				{Key: "material_id", Value: "phys_mat_001-0001-0001-0001-000000000001"},
				{Key: "summary", Value: "Grundlagen der klassischen Mechanik mit Schwerpunkt auf Newtonsche Gesetze, Bewegung, Kraft und Energie. Enthält theoretische Konzepte sowie praktische Beispiele und Problemlösungen."},
				{Key: "key_points", Value: bson.A{
					"Newtonsche Bewegungsgesetze",
					"Kinematik und Dynamik",
					"Arbeit, Energie und Leistung",
					"Impuls und Impulserhaltung",
					"Anwendungen in der Praxis",
				}},
				{Key: "language", Value: "de"},
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

// notificationsMockData retorna los mocks de la colección notifications
func notificationsMockData() mockDocument {
	return mockDocument{
		collection: "notifications",
		documents: []interface{}{
			// Notification 1 - Material ready
			bson.D{
				{Key: "user_id", Value: "llllllll-llll-llll-llll-llllllllllll"},
				{Key: "notification_type", Value: "material.ready"},
				{Key: "title", Value: "Material Procesado"},
				{Key: "message", Value: "Tu material de Biología Celular ha sido procesado exitosamente y ya está disponible."},
				{Key: "priority", Value: "low"},
				{Key: "category", Value: "academic"},
				{Key: "data", Value: bson.D{
					{Key: "resource_type", Value: "material"},
					{Key: "resource_id", Value: "bio_mat_001-0001-0001-0001-000000000001"},
					{Key: "action_url", Value: "/materials/bio_mat_001-0001-0001-0001-000000000001"},
					{Key: "action_label", Value: "Ver Material"},
					{Key: "metadata", Value: bson.D{
						{Key: "subject", Value: "Biología"},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T12:05:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
				}},
				{Key: "is_read", Value: false},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-20T12:05:00Z")},
			},
			// Notification 2 - System update
			bson.D{
				{Key: "user_id", Value: "nnnnnnnn-nnnn-nnnn-nnnn-nnnnnnnnnnnn"},
				{Key: "notification_type", Value: "system.update"},
				{Key: "title", Value: "Actualización del Sistema"},
				{Key: "message", Value: "Nueva versión 2.5.0 disponible con mejoras en rendimiento y nuevas funcionalidades."},
				{Key: "priority", Value: "medium"},
				{Key: "category", Value: "system"},
				{Key: "data", Value: bson.D{
					{Key: "metadata", Value: bson.D{
						{Key: "version", Value: "2.5.0"},
						{Key: "release_date", Value: "2025-01-20"},
						{Key: "features", Value: bson.A{"Nuevo dashboard", "Búsqueda mejorada", "Correcciones de bugs"}},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T07:00:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-20T07:00:01Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T07:00:03Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
				}},
				{Key: "is_read", Value: true},
				{Key: "read_at", Value: mustParseTime("2025-01-20T08:00:00Z")},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-20T07:00:00Z")},
			},
			// Notification 3 - Deadline reminder
			bson.D{
				{Key: "user_id", Value: "iiiiiiii-iiii-iiii-iiii-iiiiiiiiiiii"},
				{Key: "notification_type", Value: "assessment.deadline"},
				{Key: "title", Value: "Recordatorio de Evaluación"},
				{Key: "message", Value: "La evaluación de Química vence en 24 horas. ¡No olvides completarla!"},
				{Key: "priority", Value: "high"},
				{Key: "category", Value: "academic"},
				{Key: "data", Value: bson.D{
					{Key: "resource_type", Value: "assessment"},
					{Key: "resource_id", Value: "gggggggg-gggg-gggg-gggg-gggggggggggg"},
					{Key: "action_url", Value: "/assessments/gggggggg-gggg-gggg-gggg-gggggggggggg"},
					{Key: "action_label", Value: "Realizar Evaluación"},
					{Key: "metadata", Value: bson.D{
						{Key: "deadline", Value: "2025-01-21T23:59:59Z"},
						{Key: "subject", Value: "Química"},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T23:00:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-20T23:00:01Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T23:00:02Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-20T23:00:03Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T23:00:10Z")},
					}},
				}},
				{Key: "is_read", Value: false},
				{Key: "is_archived", Value: false},
				{Key: "expires_at", Value: mustParseTime("2025-01-22T00:00:00Z")},
				{Key: "created_at", Value: mustParseTime("2025-01-20T23:00:00Z")},
			},
			// Notification 4 - New comment
			bson.D{
				{Key: "user_id", Value: "ffffffff-ffff-ffff-ffff-ffffffffffff"},
				{Key: "notification_type", Value: "comment.new"},
				{Key: "title", Value: "Nuevo Comentario"},
				{Key: "message", Value: "El Prof. López respondió tu consulta sobre integrales."},
				{Key: "priority", Value: "medium"},
				{Key: "category", Value: "social"},
				{Key: "data", Value: bson.D{
					{Key: "resource_type", Value: "comment"},
					{Key: "resource_id", Value: "comment_mock_001"},
					{Key: "action_url", Value: "/materials/tttttttt-tttt-tttt-tttt-tttttttttttt/comments#comment_mock_001"},
					{Key: "action_label", Value: "Ver Comentario"},
					{Key: "metadata", Value: bson.D{
						{Key: "commenter_name", Value: "Prof. López"},
						{Key: "material_title", Value: "Tutorial: Integrales Definidas"},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T16:30:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
				}},
				{Key: "is_read", Value: false},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-20T16:30:00Z")},
			},
			// Notification 5 - Achievement unlocked
			bson.D{
				{Key: "user_id", Value: "ssssssss-ssss-ssss-ssss-ssssssssssss"},
				{Key: "notification_type", Value: "achievement.unlocked"},
				{Key: "title", Value: "¡Logro Desbloqueado!"},
				{Key: "message", Value: "Has completado 10 evaluaciones consecutivas con nota superior a 90%. ¡Excelente trabajo!"},
				{Key: "priority", Value: "low"},
				{Key: "category", Value: "gamification"},
				{Key: "data", Value: bson.D{
					{Key: "metadata", Value: bson.D{
						{Key: "achievement_id", Value: "perfect_streak_10"},
						{Key: "achievement_name", Value: "Racha Perfecta"},
						{Key: "badge_icon", Value: "🏆"},
						{Key: "points_earned", Value: 500},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T16:05:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-20T16:05:01Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T16:05:02Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
				}},
				{Key: "is_read", Value: true},
				{Key: "read_at", Value: mustParseTime("2025-01-20T17:00:00Z")},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-20T16:05:00Z")},
			},
			// Notification 6 - Security alert
			bson.D{
				{Key: "user_id", Value: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
				{Key: "notification_type", Value: "security.alert"},
				{Key: "title", Value: "Alerta de Seguridad"},
				{Key: "message", Value: "Detectamos un inicio de sesión desde un nuevo dispositivo en Buenos Aires, Argentina."},
				{Key: "priority", Value: "urgent"},
				{Key: "category", Value: "security"},
				{Key: "data", Value: bson.D{
					{Key: "metadata", Value: bson.D{
						{Key: "device_type", Value: "iPhone"},
						{Key: "location", Value: "Buenos Aires, Argentina"},
						{Key: "ip_address", Value: "181.x.x.x"},
						{Key: "timestamp", Value: "2025-01-20T08:00:00Z"},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T08:00:05Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-20T08:00:06Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T08:00:07Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-20T08:00:08Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-20T08:00:15Z")},
					}},
				}},
				{Key: "is_read", Value: true},
				{Key: "read_at", Value: mustParseTime("2025-01-20T08:05:00Z")},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-20T08:00:05Z")},
			},
		},
	}
}
