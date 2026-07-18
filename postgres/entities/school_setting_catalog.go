package entities

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Catálogo de claves de configuración por escuela (plan 039, D-039.2). Vive en
// código —no en BD— y es la ÚNICA puerta de escritura de academic.school_settings:
// lo importan academic (validación del endpoint M2M), admin-go (render/validación
// del form) y los seeds. La tabla clave/valor no lleva CHECK por clave; toda
// validez se resuelve aquí. El worker NO lo importa: recibe valores ya resueltos.

// SettingType es el tipo lógico de una clave de configuración de escuela.
type SettingType string

const (
	// SettingTypeEnum: valor de un conjunto cerrado (AllowedValues).
	SettingTypeEnum SettingType = "enum"
	// SettingTypeInt: entero positivo (> 0).
	SettingTypeInt SettingType = "int"
)

// Claves del catálogo. Son el contrato entre admin-go, academic (M2M) y los
// seeds. No agregar una clave sin registrar su SettingSpec en el catálogo.
const (
	SettingLLMGenerationMode  = "llm.generation.mode"
	SettingLLMReviewMode      = "llm.review.mode"
	SettingLLMReviewFlow      = "llm.review.flow"
	SettingLLMPipelineMode    = "llm.pipeline.mode"
	SettingImportMaxQuestions = "import.max_questions"
	SettingImportMaxJSONBytes = "import.max_json_bytes"
)

// SettingSpec describe una clave del catálogo: su tipo, los valores permitidos
// (solo enum), la env var que aporta el default de plataforma, el default duro
// (fallback si la env no está o es inválida) y una descripción en español para
// la UI admin.
type SettingSpec struct {
	Key           string
	Type          SettingType
	AllowedValues []string // solo para Type == SettingTypeEnum
	EnvDefault    string   // nombre de la env var con el default de plataforma
	HardDefault   string   // default duro si la env no está o es inválida
	Description   string
}

// schoolSettingCatalog es el catálogo inmutable de claves válidas. El orden del
// slice es el de presentación en la UI admin.
//
// Los nombres de env de import (EDUGO_IMPORT_MAX_*) espejan los que ya lee el
// validador del import 038 en edugo-api-learning
// (internal/app/api/dto/assessment_import_validator.go), y sus defaults duros
// replican DefaultMaxImportQuestions (100) y DefaultMaxImportJSONBytes (1 MiB).
var schoolSettingCatalog = []SettingSpec{
	{
		Key:           SettingLLMGenerationMode,
		Type:          SettingTypeEnum,
		AllowedValues: []string{"local", "api", "off"},
		EnvDefault:    "LLM_GENERATION_MODE_DEFAULT",
		HardDefault:   "off",
		Description:   "Modo de generación de evaluaciones por IA (plan 041): local, api u off.",
	},
	{
		Key:           SettingLLMReviewMode,
		Type:          SettingTypeEnum,
		AllowedValues: []string{"local", "api", "off"},
		EnvDefault:    "LLM_REVIEW_MODE_DEFAULT",
		HardDefault:   "off",
		Description:   "Modo de corrección de respuestas por IA (plan 040): local, api u off.",
	},
	{
		Key:           SettingLLMReviewFlow,
		Type:          SettingTypeEnum,
		AllowedValues: []string{"direct", "teacher"},
		EnvDefault:    "LLM_REVIEW_FLOW_DEFAULT",
		HardDefault:   "teacher",
		Description:   "Flujo de la corrección IA: direct (publica directo) o teacher (la revisa el docente).",
	},
	{
		Key:           SettingLLMPipelineMode,
		Type:          SettingTypeEnum,
		AllowedValues: []string{"off", "on"},
		EnvDefault:    "LLM_PIPELINE_MODE_DEFAULT",
		HardDefault:   "off",
		// Habilita el riel material→evaluación por escuela (plan 043). NO elige provider:
		// la fase 1 fuerza el LLM local por código (ADR 0036 §4); esta llave solo prende/apaga.
		Description: "Habilita la generación de evaluaciones desde material por IA (plan 043): off u on.",
	},
	{
		Key:         SettingImportMaxQuestions,
		Type:        SettingTypeInt,
		EnvDefault:  "EDUGO_IMPORT_MAX_QUESTIONS",
		HardDefault: "100",
		Description: "Máximo de preguntas por import de evaluación (deuda 019).",
	},
	{
		Key:         SettingImportMaxJSONBytes,
		Type:        SettingTypeInt,
		EnvDefault:  "EDUGO_IMPORT_MAX_JSON_BYTES",
		HardDefault: "1048576",
		Description: "Tamaño máximo del JSON de import en bytes (deuda 019).",
	},
}

// catalogIndex indexa el catálogo por clave para lookups O(1).
var catalogIndex = func() map[string]SettingSpec {
	m := make(map[string]SettingSpec, len(schoolSettingCatalog))
	for _, spec := range schoolSettingCatalog {
		m[spec.Key] = spec
	}
	return m
}()

// Catalog devuelve una copia del catálogo de claves, en orden de presentación
// (para iterar en la UI admin sin exponer el slice interno).
func Catalog() []SettingSpec {
	out := make([]SettingSpec, len(schoolSettingCatalog))
	copy(out, schoolSettingCatalog)
	return out
}

// LookupSetting devuelve la spec de una clave y si existe en el catálogo.
func LookupSetting(key string) (SettingSpec, bool) {
	spec, ok := catalogIndex[key]
	return spec, ok
}

// ValidateSetting valida que key exista en el catálogo y que value cumpla su
// tipo (enum: valor permitido; int: entero > 0). Devuelve un error descriptivo
// en caso contrario.
func ValidateSetting(key, value string) error {
	spec, ok := catalogIndex[key]
	if !ok {
		return fmt.Errorf("clave de configuración desconocida: %q", key)
	}
	switch spec.Type {
	case SettingTypeEnum:
		for _, allowed := range spec.AllowedValues {
			if value == allowed {
				return nil
			}
		}
		return fmt.Errorf("valor inválido %q para %q: permitidos %v", value, key, spec.AllowedValues)
	case SettingTypeInt:
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || n <= 0 {
			return fmt.Errorf("valor inválido %q para %q: se espera un entero > 0", value, key)
		}
		return nil
	default:
		return fmt.Errorf("tipo de configuración no soportado para %q: %s", key, spec.Type)
	}
}

// ResolveDefault devuelve el default de PLATAFORMA de una clave siguiendo la
// regla D-039.2: env var (EnvDefault) si está presente Y es válida contra el
// catálogo; si no, el default duro. Es el fallback cuando la escuela no tiene
// fila propia para la clave. Panic si la clave no existe (bug de programación).
func ResolveDefault(key string) string {
	spec, ok := catalogIndex[key]
	if !ok {
		panic(fmt.Sprintf("ResolveDefault: clave desconocida %q", key))
	}
	if raw, ok := os.LookupEnv(spec.EnvDefault); ok {
		trimmed := strings.TrimSpace(raw)
		if ValidateSetting(key, trimmed) == nil {
			return trimmed
		}
	}
	return spec.HardDefault
}
