package schemas

import (
	"embed"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed events/*.json
var schemasFS embed.FS

// EventValidator valida eventos contra JSON Schemas
type EventValidator struct {
	schemas map[string]*gojsonschema.Schema
}

// NewEventValidator crea un nuevo validador cargando todos los schemas
func NewEventValidator() (*EventValidator, error) {
	v := &EventValidator{
		schemas: make(map[string]*gojsonschema.Schema),
	}

	// Cargar schemas desde filesystem embebido
	schemaFiles := map[string]string{
		"material.uploaded:1.0":    "events/material-uploaded-v1.schema.json",
		"assessment.generated:1.0": "events/assessment-generated-v1.schema.json",
		"material.deleted:1.0":     "events/material-deleted-v1.schema.json",
		"student.enrolled:1.0":     "events/student-enrolled-v1.schema.json",
	}

	for key, filename := range schemaFiles {
		if err := v.loadSchema(key, filename); err != nil {
			return nil, fmt.Errorf("error cargando schema %s: %w", key, err)
		}
	}

	return v, nil
}

func (v *EventValidator) loadSchema(key, filename string) error {
	schemaBytes, err := schemasFS.ReadFile(filename)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return err
	}

	v.schemas[key] = schema
	return nil
}

// Validate valida un evento contra su schema correspondiente
func (v *EventValidator) Validate(event interface{}) error {
	// Extraer event_type y event_version del evento
	eventMap, ok := event.(map[string]interface{})
	if !ok {
		return errors.New("evento debe ser un objeto JSON")
	}

	eventType, ok := eventMap["event_type"].(string)
	if !ok {
		return errors.New("event_type faltante o inválido")
	}

	eventVersion, ok := eventMap["event_version"].(string)
	if !ok {
		return errors.New("event_version faltante o inválido")
	}

	return v.ValidateWithType(event, eventType, eventVersion)
}

// ValidateWithType valida especificando el tipo y versión explícitamente
func (v *EventValidator) ValidateWithType(event interface{}, eventType, eventVersion string) error {
	key := fmt.Sprintf("%s:%s", eventType, eventVersion)

	schema, exists := v.schemas[key]
	if !exists {
		return fmt.Errorf("schema no encontrado para %s", key)
	}

	documentLoader := gojsonschema.NewGoLoader(event)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		errMsg := "validación falló para " + key + ":"
		for _, desc := range result.Errors() {
			errMsg += fmt.Sprintf("\n  - %s", desc)
		}
		return errors.New(errMsg)
	}

	return nil
}

// ValidateJSON valida un evento en formato JSON bytes
func (v *EventValidator) ValidateJSON(jsonBytes []byte, eventType, eventVersion string) error {
	key := fmt.Sprintf("%s:%s", eventType, eventVersion)

	schema, exists := v.schemas[key]
	if !exists {
		return fmt.Errorf("schema no encontrado para %s", key)
	}

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		errMsg := "validación falló para " + key + ":"
		for _, desc := range result.Errors() {
			errMsg += fmt.Sprintf("\n  - %s", desc)
		}
		return errors.New(errMsg)
	}

	return nil
}
