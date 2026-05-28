package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MaterialAssessment representa la collection 'material_assessment' (o 'material_assessment_worker') en MongoDB.
// Esta entity es el reflejo exacto del schema de BD.
//
// Seed: mongodb/seeds/material_assessment_worker.js
// Usada por: worker
//
// Nota: Contiene las preguntas completas del assessment generado por IA.
// La tabla PostgreSQL 'assessment' solo tiene metadata y referencia (mongo_document_id).
type MaterialAssessment struct {
	ID               bson.ObjectID       `bson:"_id,omitempty"`
	MaterialID       string              `bson:"material_id"`        // UUID del material en PostgreSQL
	Questions        []Question          `bson:"questions"`          // Array de preguntas
	TotalQuestions   int                 `bson:"total_questions"`    // Total de preguntas
	TotalPoints      int                 `bson:"total_points"`       // Puntos totales posibles
	Version          int                 `bson:"version"`            // Versión del assessment
	AIModel          string              `bson:"ai_model"`           // Modelo IA usado (gpt-4, gpt-4-turbo, etc)
	ProcessingTimeMs int                 `bson:"processing_time_ms"` // Tiempo de procesamiento en ms
	TokenUsage       *TokenUsage         `bson:"token_usage,omitempty"`
	Metadata         *AssessmentMetadata `bson:"metadata,omitempty"` // Metadata adicional (opcional)
	CreatedAt        time.Time           `bson:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at"`
}

// Question representa una pregunta embebida en el assessment
type Question struct {
	QuestionID    string   `bson:"question_id"`
	QuestionText  string   `bson:"question_text"`
	QuestionType  string   `bson:"question_type"` // multiple_choice, true_false, open
	Options       []Option `bson:"options,omitempty"`
	CorrectAnswer string   `bson:"correct_answer"`
	Explanation   string   `bson:"explanation"`
	Points        int      `bson:"points"`
	Difficulty    string   `bson:"difficulty"` // easy, medium, hard
	Tags          []string `bson:"tags,omitempty"`
}

// Option representa una opción de respuesta
type Option struct {
	OptionID   string `bson:"option_id"`
	OptionText string `bson:"option_text"`
}

// TokenUsage representa metadata de tokens consumidos por IA
type TokenUsage struct {
	PromptTokens     int `bson:"prompt_tokens"`
	CompletionTokens int `bson:"completion_tokens"`
	TotalTokens      int `bson:"total_tokens"`
}

// AssessmentMetadata contiene metadata adicional extensible
type AssessmentMetadata struct {
	AverageDifficulty string `bson:"average_difficulty,omitempty"`
	EstimatedTimeMin  int    `bson:"estimated_time_min,omitempty"`
	SourceLength      int    `bson:"source_length,omitempty"`
	HasImages         bool   `bson:"has_images,omitempty"`
}

// CollectionName retorna el nombre de la collection en MongoDB
func (MaterialAssessment) CollectionName() string {
	return "material_assessment_worker"
}
