package migrations

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		analyticsEventsSeeds(),
		assessmentsSeeds(),
		auditLogsSeeds(),
		materialAssessmentWorkerSeeds(),
		materialSummarySeeds(),
		notificationsSeeds(),
		// assessment_attempt_result, material_content, material_event se agregarán después
	}
}

// analyticsEventsSeeds retorna los seeds de la colección analytics_events
func analyticsEventsSeeds() seedDocument {
	return seedDocument{
		collection: "analytics_events",
		documents: []interface{}{
			// Event 1 - Page view
			bson.D{
				{Key: "event_name", Value: "page.view"},
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "session_id", Value: "sess_student_abc123"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:00:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "page_path", Value: "/materials"},
					{Key: "page_title", Value: "Mis Materiales"},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "macOS"},
					{Key: "os_version", Value: "14.0"},
					{Key: "browser", Value: "Chrome"},
					{Key: "browser_version", Value: "120.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "1920x1080"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Santiago"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 2 - Material view
			bson.D{
				{Key: "event_name", Value: "material.view"},
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "session_id", Value: "sess_student_abc123"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:01:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "66666666-6666-6666-6666-666666666666"},
					{Key: "resource_type", Value: "material"},
					{Key: "custom_data", Value: bson.D{
						{Key: "material_title", Value: "Física Cuántica - Introducción"},
						{Key: "subject", Value: "Física"},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "macOS"},
					{Key: "os_version", Value: "14.0"},
					{Key: "browser", Value: "Chrome"},
					{Key: "browser_version", Value: "120.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "1920x1080"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Santiago"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 3 - Assessment start
			bson.D{
				{Key: "event_name", Value: "assessment.start"},
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "session_id", Value: "sess_student_abc123"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:14:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "99999999-9999-9999-9999-999999999999"},
					{Key: "resource_type", Value: "assessment"},
					{Key: "custom_data", Value: bson.D{
						{Key: "questions_count", Value: 2},
						{Key: "subject", Value: "Física"},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "macOS"},
					{Key: "os_version", Value: "14.0"},
					{Key: "browser", Value: "Chrome"},
					{Key: "browser_version", Value: "120.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "1920x1080"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Santiago"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 4 - Assessment complete
			bson.D{
				{Key: "event_name", Value: "assessment.complete"},
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "session_id", Value: "sess_student_abc123"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:16:15Z")},
				{Key: "properties", Value: bson.D{
					{Key: "resource_id", Value: "99999999-9999-9999-9999-999999999999"},
					{Key: "resource_type", Value: "assessment"},
					{Key: "duration_seconds", Value: 135},
					{Key: "custom_data", Value: bson.D{
						{Key: "score", Value: 100},
						{Key: "questions_count", Value: 2},
						{Key: "correct_answers", Value: 2},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "macOS"},
					{Key: "os_version", Value: "14.0"},
					{Key: "browser", Value: "Chrome"},
					{Key: "browser_version", Value: "120.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "1920x1080"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Santiago"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 5 - Search performed
			bson.D{
				{Key: "event_name", Value: "search.performed"},
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "session_id", Value: "sess_student_abc123"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:30:00Z")},
				{Key: "properties", Value: bson.D{
					{Key: "search_query", Value: "álgebra matrices"},
					{Key: "search_results_count", Value: 3},
					{Key: "custom_data", Value: bson.D{
						{Key: "filters_applied", Value: bson.D{
							{Key: "subject", Value: "Matemáticas"},
						}},
					}},
				}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "web"},
					{Key: "os", Value: "macOS"},
					{Key: "os_version", Value: "14.0"},
					{Key: "browser", Value: "Chrome"},
					{Key: "browser_version", Value: "120.0"},
					{Key: "device_type", Value: "desktop"},
					{Key: "screen_resolution", Value: "1920x1080"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Santiago"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					{Key: "user_role", Value: "student"},
				}},
			},
			// Event 6 - Mobile app session
			bson.D{
				{Key: "event_name", Value: "session.start"},
				{Key: "user_id", Value: "44444444-4444-4444-4444-444444444444"},
				{Key: "session_id", Value: "sess_mobile_xyz789"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T11:00:00Z")},
				{Key: "properties", Value: bson.D{}},
				{Key: "device", Value: bson.D{
					{Key: "platform", Value: "android"},
					{Key: "os", Value: "Android"},
					{Key: "os_version", Value: "13"},
					{Key: "device_type", Value: "mobile"},
					{Key: "screen_resolution", Value: "1080x2400"},
				}},
				{Key: "location", Value: bson.D{
					{Key: "country", Value: "CL"},
					{Key: "city", Value: "Valparaíso"},
					{Key: "timezone", Value: "America/Santiago"},
				}},
				{Key: "context", Value: bson.D{
					{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					{Key: "user_role", Value: "student"},
				}},
			},
		},
	}
}

// assessmentsSeeds retorna los seeds de la colección material_assessment
func assessmentsSeeds() seedDocument {
	return seedDocument{
		collection: "material_assessment",
		documents: []interface{}{
			// Assessment 1 (Física)
			bson.D{
				{Key: "_id", Value: mustParseObjectID("507f1f77bcf86cd799439011")},
				{Key: "material_id", Value: "66666666-6666-6666-6666-666666666666"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_index", Value: 0},
						{Key: "question_text", Value: "¿Qué es la dualidad onda-partícula?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{
								{Key: "option_index", Value: 0},
								{Key: "text", Value: "Partículas que actúan solo como ondas"},
								{Key: "is_correct", Value: false},
							},
							bson.D{
								{Key: "option_index", Value: 1},
								{Key: "text", Value: "Partículas que pueden comportarse como ondas y viceversa"},
								{Key: "is_correct", Value: true},
							},
							bson.D{
								{Key: "option_index", Value: 2},
								{Key: "text", Value: "Ondas que no son partículas"},
								{Key: "is_correct", Value: false},
							},
							bson.D{
								{Key: "option_index", Value: 3},
								{Key: "text", Value: "Ninguna de las anteriores"},
								{Key: "is_correct", Value: false},
							},
						}},
					},
					bson.D{
						{Key: "question_index", Value: 1},
						{Key: "question_text", Value: "¿Quién propuso el principio de incertidumbre?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{
								{Key: "option_index", Value: 0},
								{Key: "text", Value: "Einstein"},
								{Key: "is_correct", Value: false},
							},
							bson.D{
								{Key: "option_index", Value: 1},
								{Key: "text", Value: "Heisenberg"},
								{Key: "is_correct", Value: true},
							},
							bson.D{
								{Key: "option_index", Value: 2},
								{Key: "text", Value: "Bohr"},
								{Key: "is_correct", Value: false},
							},
							bson.D{
								{Key: "option_index", Value: 3},
								{Key: "text", Value: "Schrödinger"},
								{Key: "is_correct", Value: false},
							},
						}},
					},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "subject", Value: "Física"},
					{Key: "grade", Value: "10th"},
					{Key: "difficulty", Value: "medium"},
				}},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
			// Assessment 2 (Álgebra)
			bson.D{
				{Key: "_id", Value: mustParseObjectID("507f1f77bcf86cd799439012")},
				{Key: "material_id", Value: "77777777-7777-7777-7777-777777777777"},
				{Key: "questions", Value: bson.A{
					bson.D{
						{Key: "question_index", Value: 0},
						{Key: "question_text", Value: "¿Qué es una matriz identidad?"},
						{Key: "question_type", Value: "multiple_choice"},
						{Key: "options", Value: bson.A{
							bson.D{
								{Key: "option_index", Value: 0},
								{Key: "text", Value: "Una matriz con todos 1s"},
								{Key: "is_correct", Value: false},
							},
							bson.D{
								{Key: "option_index", Value: 1},
								{Key: "text", Value: "Una matriz con 1s en la diagonal y 0s en el resto"},
								{Key: "is_correct", Value: true},
							},
							bson.D{
								{Key: "option_index", Value: 2},
								{Key: "text", Value: "Una matriz cuadrada"},
								{Key: "is_correct", Value: false},
							},
							bson.D{
								{Key: "option_index", Value: 3},
								{Key: "text", Value: "Una matriz invertible"},
								{Key: "is_correct", Value: false},
							},
						}},
					},
				}},
				{Key: "metadata", Value: bson.D{
					{Key: "subject", Value: "Matemáticas"},
					{Key: "grade", Value: "11th"},
					{Key: "difficulty", Value: "easy"},
				}},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			},
		},
	}
}

// auditLogsSeeds retorna los seeds de la colección audit_logs
func auditLogsSeeds() seedDocument {
	return seedDocument{
		collection: "audit_logs",
		documents: []interface{}{
			// Audit log 1 - User login
			bson.D{
				{Key: "event_type", Value: "user.login"},
				{Key: "actor_id", Value: "11111111-1111-1111-1111-111111111111"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "user"},
				{Key: "resource_id", Value: "11111111-1111-1111-1111-111111111111"},
				{Key: "action", Value: "login"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.1.100"},
					{Key: "user_agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"},
					{Key: "metadata", Value: bson.D{
						{Key: "login_method", Value: "email_password"},
						{Key: "remember_me", Value: true},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T09:00:00Z")},
				{Key: "session_id", Value: "sess_abc123xyz"},
				{Key: "request_id", Value: "req_001"},
			},
			// Audit log 2 - Material uploaded
			bson.D{
				{Key: "event_type", Value: "material.uploaded"},
				{Key: "actor_id", Value: "22222222-2222-2222-2222-222222222222"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "material"},
				{Key: "resource_id", Value: "66666666-6666-6666-6666-666666666666"},
				{Key: "action", Value: "upload"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.1.101"},
					{Key: "user_agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
					{Key: "changes", Value: bson.D{
						{Key: "file_name", Value: "fisica_cuantica.pdf"},
						{Key: "file_size", Value: 2048576},
						{Key: "school_id", Value: "55555555-5555-5555-5555-555555555555"},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:00:00Z")},
				{Key: "session_id", Value: "sess_def456uvw"},
				{Key: "request_id", Value: "req_002"},
			},
			// Audit log 3 - Assessment published
			bson.D{
				{Key: "event_type", Value: "assessment.published"},
				{Key: "actor_id", Value: "22222222-2222-2222-2222-222222222222"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "assessment"},
				{Key: "resource_id", Value: "99999999-9999-9999-9999-999999999999"},
				{Key: "action", Value: "update"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.1.101"},
					{Key: "user_agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
					{Key: "changes", Value: bson.D{
						{Key: "status", Value: bson.D{
							{Key: "from", Value: "generated"},
							{Key: "to", Value: "published"},
						}},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T10:30:00Z")},
				{Key: "session_id", Value: "sess_def456uvw"},
				{Key: "request_id", Value: "req_003"},
			},
			// Audit log 4 - Failed login attempt
			bson.D{
				{Key: "event_type", Value: "user.login"},
				{Key: "actor_id", Value: "unknown"},
				{Key: "actor_type", Value: "user"},
				{Key: "resource_type", Value: "user"},
				{Key: "resource_id", Value: "unknown"},
				{Key: "action", Value: "login"},
				{Key: "details", Value: bson.D{
					{Key: "ip_address", Value: "192.168.1.200"},
					{Key: "user_agent", Value: "Mozilla/5.0 (X11; Linux x86_64)"},
					{Key: "error", Value: bson.D{
						{Key: "code", Value: "INVALID_CREDENTIALS"},
						{Key: "message", Value: "Invalid email or password"},
					}},
					{Key: "metadata", Value: bson.D{
						{Key: "attempted_email", Value: "test@example.com"},
					}},
				}},
				{Key: "severity", Value: "warning"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T11:00:00Z")},
				{Key: "request_id", Value: "req_004"},
			},
			// Audit log 5 - System backup
			bson.D{
				{Key: "event_type", Value: "system.backup"},
				{Key: "actor_id", Value: "system"},
				{Key: "actor_type", Value: "system"},
				{Key: "resource_type", Value: "system"},
				{Key: "action", Value: "create"},
				{Key: "details", Value: bson.D{
					{Key: "metadata", Value: bson.D{
						{Key: "backup_type", Value: "automated_daily"},
						{Key: "backup_size_mb", Value: 1024},
						{Key: "backup_location", Value: "s3://edugo-backups/2025-01-15/"},
					}},
				}},
				{Key: "severity", Value: "info"},
				{Key: "timestamp", Value: mustParseTime("2025-01-15T02:00:00Z")},
				{Key: "request_id", Value: "req_005"},
			},
		},
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
				}},
				{Key: "total_questions", Value: 2},
				{Key: "total_points", Value: 20},
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

// notificationsSeeds retorna los seeds de la colección notifications
func notificationsSeeds() seedDocument {
	return seedDocument{
		collection: "notifications",
		documents: []interface{}{
			// Notification 1 - Assessment ready
			bson.D{
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "notification_type", Value: "assessment.ready"},
				{Key: "title", Value: "Nuevo Assessment Disponible"},
				{Key: "message", Value: "Tu profesor ha publicado un nuevo assessment de Física Cuántica. ¡Es hora de demostrar lo que has aprendido!"},
				{Key: "priority", Value: "medium"},
				{Key: "category", Value: "academic"},
				{Key: "data", Value: bson.D{
					{Key: "resource_type", Value: "assessment"},
					{Key: "resource_id", Value: "99999999-9999-9999-9999-999999999999"},
					{Key: "action_url", Value: "/assessments/99999999-9999-9999-9999-999999999999"},
					{Key: "action_label", Value: "Comenzar Assessment"},
					{Key: "metadata", Value: bson.D{
						{Key: "subject", Value: "Física"},
						{Key: "teacher_name", Value: "Prof. García"},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T10:30:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-15T10:30:01Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T10:30:02Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: false},
					}},
				}},
				{Key: "is_read", Value: false},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-15T10:30:00Z")},
			},
			// Notification 2 - Assessment graded
			bson.D{
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "notification_type", Value: "assessment.graded"},
				{Key: "title", Value: "Assessment Calificado"},
				{Key: "message", Value: "¡Felicitaciones! Has obtenido 100% en el assessment de Física Cuántica."},
				{Key: "priority", Value: "high"},
				{Key: "category", Value: "academic"},
				{Key: "data", Value: bson.D{
					{Key: "resource_type", Value: "attempt"},
					{Key: "resource_id", Value: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
					{Key: "action_url", Value: "/results/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
					{Key: "action_label", Value: "Ver Resultados"},
					{Key: "metadata", Value: bson.D{
						{Key: "score", Value: 100},
						{Key: "total_questions", Value: 2},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T10:16:15Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-15T10:16:16Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T10:16:17Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-15T10:16:18Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T10:16:25Z")},
					}},
				}},
				{Key: "is_read", Value: true},
				{Key: "read_at", Value: mustParseTime("2025-01-15T10:20:00Z")},
				{Key: "is_archived", Value: false},
				{Key: "created_at", Value: mustParseTime("2025-01-15T10:16:15Z")},
			},
			// Notification 3 - Material uploaded
			bson.D{
				{Key: "user_id", Value: "33333333-3333-3333-3333-333333333333"},
				{Key: "notification_type", Value: "material.uploaded"},
				{Key: "title", Value: "Nuevo Material Disponible"},
				{Key: "message", Value: "El Prof. García ha subido un nuevo material: Álgebra Lineal - Matrices"},
				{Key: "priority", Value: "low"},
				{Key: "category", Value: "academic"},
				{Key: "data", Value: bson.D{
					{Key: "resource_type", Value: "material"},
					{Key: "resource_id", Value: "77777777-7777-7777-7777-777777777777"},
					{Key: "action_url", Value: "/materials/77777777-7777-7777-7777-777777777777"},
					{Key: "action_label", Value: "Ver Material"},
					{Key: "metadata", Value: bson.D{
						{Key: "subject", Value: "Matemáticas"},
						{Key: "teacher_name", Value: "Prof. García"},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-11T14:20:00Z")},
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
				{Key: "created_at", Value: mustParseTime("2025-01-11T14:20:00Z")},
			},
			// Notification 4 - System announcement
			bson.D{
				{Key: "user_id", Value: "11111111-1111-1111-1111-111111111111"},
				{Key: "notification_type", Value: "system.announcement"},
				{Key: "title", Value: "Mantenimiento Programado"},
				{Key: "message", Value: "El sistema estará en mantenimiento el sábado 18 de enero de 02:00 a 04:00 hrs."},
				{Key: "priority", Value: "urgent"},
				{Key: "category", Value: "system"},
				{Key: "data", Value: bson.D{
					{Key: "metadata", Value: bson.D{
						{Key: "maintenance_start", Value: "2025-01-18T02:00:00Z"},
						{Key: "maintenance_end", Value: "2025-01-18T04:00:00Z"},
						{Key: "services_affected", Value: bson.A{"assessments", "materials"}},
					}},
				}},
				{Key: "delivery", Value: bson.D{
					{Key: "in_app", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T08:00:00Z")},
					}},
					{Key: "push", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-15T08:00:01Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T08:00:02Z")},
					}},
					{Key: "email", Value: bson.D{
						{Key: "enabled", Value: true},
						{Key: "sent_at", Value: mustParseTime("2025-01-15T08:00:03Z")},
						{Key: "delivered_at", Value: mustParseTime("2025-01-15T08:00:10Z")},
					}},
				}},
				{Key: "is_read", Value: true},
				{Key: "read_at", Value: mustParseTime("2025-01-15T09:00:00Z")},
				{Key: "is_archived", Value: false},
				{Key: "expires_at", Value: mustParseTime("2025-01-18T05:00:00Z")},
				{Key: "created_at", Value: mustParseTime("2025-01-15T08:00:00Z")},
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

// mustParseObjectID parsea un ObjectID hex o entra en pánico (solo para seeds)
func mustParseObjectID(s string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(fmt.Sprintf("invalid ObjectID: %s", s))
	}
	return oid
}
