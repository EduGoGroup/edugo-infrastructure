package schemas_test

import (
	"encoding/json"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/schemas"
)

func TestMaterialUploadedValidation(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	// Evento válido
	validEvent := map[string]interface{}{
		"event_id":      "01JA8XYZ-1234-5678-90AB-CDEF12345678",
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-15T10:30:00Z",
		"payload": map[string]interface{}{
			"material_id":     "66666666-6666-6666-6666-666666666666",
			"school_id":       "44444444-4444-4444-4444-444444444444",
			"teacher_id":      "22222222-2222-2222-2222-222222222222",
			"file_url":        "s3://edugo-materials/test.pdf",
			"file_size_bytes": 2048000,
			"file_type":       "application/pdf",
		},
	}

	if err := validator.Validate(validEvent); err != nil {
		t.Errorf("Evento válido rechazado: %v", err)
	}

	// Evento inválido (falta campo requerido)
	invalidEvent := map[string]interface{}{
		"event_id":      "01JA8XYZ-1234-5678-90AB-CDEF12345678",
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-15T10:30:00Z",
		"payload": map[string]interface{}{
			// Falta material_id, school_id, etc.
			"file_url": "s3://edugo-materials/test.pdf",
		},
	}

	if err := validator.Validate(invalidEvent); err == nil {
		t.Error("Evento inválido fue aceptado (esperaba error)")
	}
}

func TestValidateJSON(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	validJSON := `{
		"event_id": "01JA8XYZ-1234-5678-90AB-CDEF12345678",
		"event_type": "assessment.generated",
		"event_version": "1.0",
		"timestamp": "2025-11-15T10:35:00Z",
		"payload": {
			"material_id": "66666666-6666-6666-6666-666666666666",
			"mongo_document_id": "507f1f77bcf86cd799439011",
			"questions_count": 8
		}
	}`

	if err := validator.ValidateJSON([]byte(validJSON), "assessment.generated", "1.0"); err != nil {
		t.Errorf("JSON válido rechazado: %v", err)
	}
}

// Ejemplo de uso en publisher
func ExampleEventValidator_publisher() {
	validator, _ := schemas.NewEventValidator()

	event := map[string]interface{}{
		"event_id":      "01JA8XYZ-1234-5678-90AB-CDEF12345678",
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-15T10:30:00Z",
		"payload": map[string]interface{}{
			"material_id":     "66666666-6666-6666-6666-666666666666",
			"school_id":       "44444444-4444-4444-4444-444444444444",
			"teacher_id":      "22222222-2222-2222-2222-222222222222",
			"file_url":        "s3://edugo-materials/test.pdf",
			"file_size_bytes": 2048000,
			"file_type":       "application/pdf",
		},
	}

	// Validar antes de publicar
	if err := validator.Validate(event); err != nil {
		// Log error y NO publicar
		return
	}

	// Publicar evento (solo si es válido)
	eventJSON, _ := json.Marshal(event)
	_ = eventJSON // publisher.Publish(exchange, routingKey, eventJSON)
}

// Ejemplo de uso en consumer
func ExampleEventValidator_consumer() {
	validator, _ := schemas.NewEventValidator()

	// Mensaje recibido de RabbitMQ
	messageBytes := []byte(`{"event_type": "material.uploaded", "event_version": "1.0", ...}`)

	// Validar JSON recibido
	if err := validator.ValidateJSON(messageBytes, "material.uploaded", "1.0"); err != nil {
		// Evento inválido → enviar a DLQ
		return
	}

	// Procesar evento (solo si es válido)
	var event map[string]interface{}
	json.Unmarshal(messageBytes, &event)
	// ... procesar ...
}
