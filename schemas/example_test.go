package schemas_test

import (
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/schemas"
	"github.com/google/uuid"
)

func TestMaterialUploadedValidation(t *testing.T) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		t.Fatalf("Error creando validador: %v", err)
	}

	// Evento válido con UUID real
	validEvent := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-15T10:30:00Z",
		"payload": map[string]interface{}{
			"material_id":     uuid.New().String(),
			"school_id":       uuid.New().String(),
			"teacher_id":      uuid.New().String(),
			"file_url":        "s3://edugo-materials/test.pdf",
			"file_size_bytes": float64(2048000),
			"file_type":       "application/pdf",
		},
	}

	if err := validator.Validate(validEvent); err != nil {
		t.Errorf("Evento válido rechazado: %v", err)
	}

	// Evento inválido (falta campo requerido)
	invalidEvent := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     "2025-11-15T10:30:00Z",
		"payload": map[string]interface{}{
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

	materialID := uuid.New().String()
	eventID := uuid.New().String()

	validJSON := fmt.Sprintf(`{
		"event_id": "%s",
		"event_type": "assessment.generated",
		"event_version": "1.0",
		"timestamp": "2025-11-15T10:35:00Z",
		"payload": {
			"material_id": "%s",
			"mongo_document_id": "507f1f77bcf86cd799439011",
			"questions_count": 8
		}
	}`, eventID, materialID)

	if err := validator.ValidateJSON([]byte(validJSON), "assessment.generated", "1.0"); err != nil {
		t.Errorf("JSON válido rechazado: %v", err)
	}
}
