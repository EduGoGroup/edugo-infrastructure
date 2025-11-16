package schemas

import (
	"embed"
	"fmt"
	"path/filepath"

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
	schemaFiles := []string{
		"material.uploaded:1.0",
		"assessment.generated:1.0",
		"material.deleted:1.0",
		"student.enrolled:1.0",
	}

	for _, key := range schemaFiles {
		if err := v.loadSchema(key); err != nil {
			return nil, fmt.Errorf("error cargando schema %s: %w", key, err)
		}
	}

	return v, nil
}

func (v *EventValidator) loadSchema(key string) error {
	filename := getSchemaFilename(key)
	
	schemaBytes, err := schemasFS.ReadFile(filepath.Join("events", filename))
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
		return fmt.Errorf("evento debe ser un objeto JSON")
	}

	eventType, ok := eventMap["event_type"].(string)
	if !ok {
		return fmt.Errorf("event_type faltante o inválido")
	}

	eventVersion, ok := eventMap["event_version"].(string)
	if !ok {
		return fmt.Errorf("event_version faltante o inválido")
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
		errMsg := fmt.Sprintf("validación falló para %s:", key)
		for _, desc := range result.Errors() {
			errMsg += fmt.Sprintf("\n  - %s", desc)
		}
		return fmt.Errorf(errMsg)
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
		errMsg := fmt.Sprintf("validación falló para %s:", key)
		for _, desc := range result.Errors() {
			errMsg += fmt.Sprintf("\n  - %s", desc)
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}

func getSchemaFilename(key string) string {
	// material.uploaded:1.0 → material-uploaded-v1.schema.json
	parts := splitKey(key)
	eventType := parts[0]
	version := parts[1]
	
	eventType = replaceAll(eventType, ".", "-")
	version = replaceAll(version, ".", "")
	
	return fmt.Sprintf("%s-v%s.schema.json", eventType, version)
}

func splitKey(key string) [2]string {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == ':' {
			return [2]string{key[:i], key[i+1:]}
		}
	}
	return [2]string{key, ""}
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
