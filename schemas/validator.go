package schemas

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed events/*.json
var schemasFS embed.FS

// EventValidator validates events against JSON Schemas.
type EventValidator struct {
	schemas map[string]*gojsonschema.Schema
}

// NewEventValidator creates a new validator by loading all schemas.
func NewEventValidator() (*EventValidator, error) {
	v := &EventValidator{
		schemas: make(map[string]*gojsonschema.Schema),
	}

	files, err := schemasFS.ReadDir("events")
	if err != nil {
		return nil, fmt.Errorf("error reading schemas directory: %w", err)
	}

	for _, file := range files {
		filename := file.Name()
		if !strings.HasSuffix(filename, ".schema.json") {
			continue
		}

		// Derive key from filename, e.g., "material-uploaded-v1.schema.json" -> "material.uploaded:1.0"
		base := strings.TrimSuffix(filename, ".schema.json")
		parts := strings.Split(base, "-")
		if len(parts) < 2 {
			continue
		}
		eventType := strings.Join(parts[:len(parts)-1], ".")
		version := strings.Replace(parts[len(parts)-1], "v", "", 1) + ".0"
		key := fmt.Sprintf("%s:%s", eventType, version)

		filepath := path.Join("events", filename)
		if err := v.loadSchema(key, filepath); err != nil {
			return nil, fmt.Errorf("error loading schema %s: %w", key, err)
		}
	}

	return v, nil
}

func (v *EventValidator) loadSchema(key, filename string) error {
	schemaBytes, err := schemasFS.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", filename, err)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("failed to compile schema %s: %w", filename, err)
	}

	v.schemas[key] = schema
	return nil
}

// Validate validates an event against its corresponding schema.
func (v *EventValidator) Validate(event interface{}) error {
	// Safely extract event_type and event_version by marshaling and unmarshaling.
	var tempEvent struct {
		EventType    string `json:"event_type"`
		EventVersion string `json:"event_version"`
	}

	// Convert the event to JSON bytes to safely inspect its contents.
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}
	if err := json.Unmarshal(eventBytes, &tempEvent); err != nil {
		return fmt.Errorf("failed to unmarshal event metadata: %w", err)
	}

	if tempEvent.EventType == "" {
		return errors.New("event_type missing or invalid")
	}
	if tempEvent.EventVersion == "" {
		return errors.New("event_version missing or invalid")
	}

	return v.ValidateWithType(gojsonschema.NewBytesLoader(eventBytes), tempEvent.EventType, tempEvent.EventVersion)
}

// ValidateWithType validates by explicitly specifying the type and version.
func (v *EventValidator) ValidateWithType(documentLoader gojsonschema.JSONLoader, eventType, eventVersion string) error {
	key := fmt.Sprintf("%s:%s", eventType, eventVersion)

	schema, exists := v.schemas[key]
	if !exists {
		return fmt.Errorf("schema not found for %s", key)
	}

	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("error during validation for %s: %w", key, err)
	}

	if !result.Valid() {
		return v.formatValidationErrors(key, result.Errors())
	}

	return nil
}

// ValidateJSON validates an event in JSON byte format.
func (v *EventValidator) ValidateJSON(jsonBytes []byte, eventType, eventVersion string) error {
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	return v.ValidateWithType(documentLoader, eventType, eventVersion)
}

// formatValidationErrors formats validation errors into a single error message.
func (v *EventValidator) formatValidationErrors(schemaKey string, validationErrors []gojsonschema.ResultError) error {
	var errorMessages strings.Builder
	errorMessages.WriteString(fmt.Sprintf("validation failed for %s:", schemaKey))

	for _, desc := range validationErrors {
		errorMessages.WriteString(fmt.Sprintf("\n  - %s", desc))
	}

	return errors.New(errorMessages.String())
}
