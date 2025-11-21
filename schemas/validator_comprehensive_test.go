package schemas_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/schemas"
	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
)

// TestMaterialDeletedValidation valida el evento material.deleted
func TestMaterialDeletedValidation(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	t.Run("happy_path_valid_event", func(t *testing.T) {
		validEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(validEvent); err != nil {
			t.Errorf("Evento válido rechazado: %v", err)
		}
	})

	t.Run("with_optional_reason", func(t *testing.T) {
		validEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
				"reason":             "Material obsoleto",
			},
		}

		if err := validator.Validate(validEvent); err != nil {
			t.Errorf("Evento válido con reason rechazado: %v", err)
		}
	})

	t.Run("missing_required_field", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id": uuid.New().String(),
				"school_id":   uuid.New().String(),
				// Falta deleted_by_user_id
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por campo faltante")
		}
	})

	t.Run("invalid_uuid_format", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        "not-a-uuid",
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por UUID inválido")
		}
	})

	t.Run("invalid_timestamp_format", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16 10:00:00", // Formato incorrecto
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por timestamp inválido")
		}
	})

	t.Run("wrong_event_type", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.uploaded", // Tipo incorrecto
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		// El validador usará material.uploaded:1.0 que requiere campos diferentes
		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por payload incompatible")
		}
	})
}

// TestStudentEnrolledValidation valida el evento student.enrolled
func TestStudentEnrolledValidation(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	t.Run("happy_path_valid_event", func(t *testing.T) {
		validEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "student.enrolled",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"student_id":    uuid.New().String(),
				"school_id":     uuid.New().String(),
				"membership_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(validEvent); err != nil {
			t.Errorf("Evento válido rechazado: %v", err)
		}
	})

	t.Run("with_optional_fields", func(t *testing.T) {
		validEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "student.enrolled",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"student_id":          uuid.New().String(),
				"school_id":           uuid.New().String(),
				"membership_id":       uuid.New().String(),
				"academic_unit_id":    uuid.New().String(),
				"enrolled_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(validEvent); err != nil {
			t.Errorf("Evento válido con campos opcionales rechazado: %v", err)
		}
	})

	t.Run("missing_required_field", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "student.enrolled",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"student_id": uuid.New().String(),
				"school_id":  uuid.New().String(),
				// Falta membership_id
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por campo faltante")
		}
	})

	t.Run("invalid_uuid_in_payload", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "student.enrolled",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"student_id":    "123",
				"school_id":     uuid.New().String(),
				"membership_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por UUID inválido")
		}
	})
}

// TestEventTypeValidation valida edge cases relacionados con event_type y event_version
func TestEventTypeValidation(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	tests := []struct {
		name          string
		event         map[string]interface{}
		shouldError   bool
		errorContains string
	}{
		{
			name: "missing_event_type",
			event: map[string]interface{}{
				"event_id":      uuid.New().String(),
				"event_version": "1.0",
				"timestamp":     "2025-11-16T10:00:00Z",
				"payload":       map[string]interface{}{},
			},
			shouldError:   true,
			errorContains: "event_type missing",
		},
		{
			name: "missing_event_version",
			event: map[string]interface{}{
				"event_id":   uuid.New().String(),
				"event_type": "material.uploaded",
				"timestamp":  "2025-11-16T10:00:00Z",
				"payload":    map[string]interface{}{},
			},
			shouldError:   true,
			errorContains: "event_version missing",
		},
		{
			name: "unknown_event_type",
			event: map[string]interface{}{
				"event_id":      uuid.New().String(),
				"event_type":    "unknown.event",
				"event_version": "1.0",
				"timestamp":     "2025-11-16T10:00:00Z",
				"payload":       map[string]interface{}{},
			},
			shouldError:   true,
			errorContains: "schema not found",
		},
		{
			name: "unknown_event_version",
			event: map[string]interface{}{
				"event_id":      uuid.New().String(),
				"event_type":    "material.uploaded",
				"event_version": "99.0",
				"timestamp":     "2025-11-16T10:00:00Z",
				"payload":       map[string]interface{}{},
			},
			shouldError:   true,
			errorContains: "schema not found",
		},
		{
			name: "event_type_wrong_type",
			event: map[string]interface{}{
				"event_id":      uuid.New().String(),
				"event_type":    123, // Número en vez de string
				"event_version": "1.0",
				"timestamp":     "2025-11-16T10:00:00Z",
				"payload":       map[string]interface{}{},
			},
			shouldError:   true,
			errorContains: "failed to unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.event)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Se esperaba error pero no ocurrió")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Error esperaba contener %q, pero obtuvo: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba error pero ocurrió: %v", err)
				}
			}
		})
	}
}

// TestInvalidFormats valida casos de formatos inválidos
func TestInvalidFormats(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	t.Run("invalid_event_id_uuid", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      "not-a-uuid",
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por event_id inválido")
		}
	})

	t.Run("empty_string_uuid", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        "",
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por UUID vacío")
		}
	})

	t.Run("timestamp_without_timezone", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por timestamp sin timezone")
		}
	})

	t.Run("timestamp_invalid_format", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "16/11/2025 10:00:00",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por timestamp en formato incorrecto")
		}
	})

	t.Run("number_as_uuid", func(t *testing.T) {
		invalidEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        12345,
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.Validate(invalidEvent); err == nil {
			t.Error("Se esperaba error por tipo incorrecto")
		}
	})
}

// TestValidateJSONMethod valida el método ValidateJSON con bytes (casos exhaustivos)
func TestValidateJSONMethod(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	t.Run("valid_json_bytes", func(t *testing.T) {
		validEvent := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		jsonBytes, err := json.Marshal(validEvent)
		if err != nil {
			t.Fatalf("Error marshaling JSON: %v", err)
		}

		if err := validator.ValidateJSON(jsonBytes, "material.deleted", "1.0"); err != nil {
			t.Errorf("Validación falló: %v", err)
		}
	})

	t.Run("invalid_json_bytes", func(t *testing.T) {
		invalidJSON := []byte(`{"event_id": "not-a-uuid"}`)

		if err := validator.ValidateJSON(invalidJSON, "material.deleted", "1.0"); err == nil {
			t.Error("Se esperaba error por JSON inválido")
		}
	})

	t.Run("malformed_json", func(t *testing.T) {
		malformedJSON := []byte(`{invalid json}`)

		if err := validator.ValidateJSON(malformedJSON, "material.deleted", "1.0"); err == nil {
			t.Error("Se esperaba error por JSON malformado")
		}
	})

	t.Run("unknown_schema", func(t *testing.T) {
		validJSON := []byte(`{"event_id": "` + uuid.New().String() + `"}`)

		if err := validator.ValidateJSON(validJSON, "unknown.event", "1.0"); err == nil {
			t.Error("Se esperaba error por schema desconocido")
		}
	})
}

// TestValidateWithType valida el método ValidateWithType
func TestValidateWithType(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	t.Run("explicit_type_version", func(t *testing.T) {
		event := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		if err := validator.ValidateWithType(gojsonschema.NewGoLoader(event), "material.deleted", "1.0"); err != nil {
			t.Errorf("Validación falló: %v", err)
		}
	})

	t.Run("type_mismatch", func(t *testing.T) {
		event := map[string]interface{}{
			"event_id":      uuid.New().String(),
			"event_type":    "material.deleted",
			"event_version": "1.0",
			"timestamp":     "2025-11-16T10:00:00Z",
			"payload": map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		}

		// Validar como student.enrolled aunque el payload es de material.deleted
		if err := validator.ValidateWithType(gojsonschema.NewGoLoader(event), "student.enrolled", "1.0"); err == nil {
			t.Error("Se esperaba error por payload incompatible")
		}
	})
}

// TestAllFourSchemas valida que los 4 schemas funcionan correctamente
func TestAllFourSchemas(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	schemas := []struct {
		eventType    string
		eventVersion string
		payload      map[string]interface{}
	}{
		{
			eventType:    "material.uploaded",
			eventVersion: "1.0",
			payload: map[string]interface{}{
				"material_id":     uuid.New().String(),
				"school_id":       uuid.New().String(),
				"teacher_id":      uuid.New().String(),
				"file_url":        "s3://bucket/file.pdf",
				"file_size_bytes": float64(1024000),
				"file_type":       "application/pdf",
			},
		},
		{
			eventType:    "assessment.generated",
			eventVersion: "1.0",
			payload: map[string]interface{}{
				"material_id":       uuid.New().String(),
				"mongo_document_id": "507f1f77bcf86cd799439011", // MongoDB ObjectID válido (24 hex chars)
				"questions_count":   10,                         // integer, no float
			},
		},
		{
			eventType:    "material.deleted",
			eventVersion: "1.0",
			payload: map[string]interface{}{
				"material_id":        uuid.New().String(),
				"school_id":          uuid.New().String(),
				"deleted_by_user_id": uuid.New().String(),
			},
		},
		{
			eventType:    "student.enrolled",
			eventVersion: "1.0",
			payload: map[string]interface{}{
				"student_id":    uuid.New().String(),
				"school_id":     uuid.New().String(),
				"membership_id": uuid.New().String(),
			},
		},
	}

	for _, schema := range schemas {
		t.Run(schema.eventType, func(t *testing.T) {
			event := map[string]interface{}{
				"event_id":      uuid.New().String(),
				"event_type":    schema.eventType,
				"event_version": schema.eventVersion,
				"timestamp":     "2025-11-16T10:00:00Z",
				"payload":       schema.payload,
			}

			if err := validator.Validate(event); err != nil {
				t.Errorf("Schema %s falló: %v", schema.eventType, err)
			}
		})
	}
}

// TestNotObjectEvent valida el comportamiento cuando el evento no es un objeto
func TestNotObjectEvent(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	invalidEvents := []interface{}{
		"string",
		123,
		[]string{"array"},
		nil,
	}

	for _, invalidEvent := range invalidEvents {
		if err := validator.Validate(invalidEvent); err == nil {
			t.Errorf("Se esperaba error para evento no-objeto: %v", invalidEvent)
		}
	}
}

// BenchmarkValidation mide la performance de validación con 10,000 eventos
func BenchmarkValidation(b *testing.B) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		b.Fatalf("Error creando validador: %v", err)
	}

	event := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-16T10:00:00Z",
		"payload": map[string]interface{}{
			"material_id":     uuid.New().String(),
			"school_id":       uuid.New().String(),
			"teacher_id":      uuid.New().String(),
			"file_url":        "s3://edugo-materials/test.pdf",
			"file_size_bytes": float64(2048000),
			"file_type":       "application/pdf",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(event)
	}
}

// BenchmarkValidationJSON mide la performance de ValidateJSON
func BenchmarkValidationJSON(b *testing.B) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		b.Fatalf("Error creando validador: %v", err)
	}

	event := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-16T10:00:00Z",
		"payload": map[string]interface{}{
			"material_id":     uuid.New().String(),
			"school_id":       uuid.New().String(),
			"teacher_id":      uuid.New().String(),
			"file_url":        "s3://edugo-materials/test.pdf",
			"file_size_bytes": float64(2048000),
			"file_type":       "application/pdf",
		},
	}

	jsonBytes, _ := json.Marshal(event)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateJSON(jsonBytes, "material.uploaded", "1.0")
	}
}

// BenchmarkValidation10000 valida específicamente 10,000 eventos
func BenchmarkValidation10000(b *testing.B) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		b.Fatalf("Error creando validador: %v", err)
	}

	event := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-16T10:00:00Z",
		"payload": map[string]interface{}{
			"material_id":     uuid.New().String(),
			"school_id":       uuid.New().String(),
			"teacher_id":      uuid.New().String(),
			"file_url":        "s3://edugo-materials/test.pdf",
			"file_size_bytes": float64(2048000),
			"file_type":       "application/pdf",
		},
	}

	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		validator.Validate(event)
	}
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
