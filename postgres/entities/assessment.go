package entities

import (
	"time"

	"github.com/google/uuid"
)

// Assessment representa la tabla 'assessment' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 006_create_assessments.up.sql, 009_extend_assessment_schema.up.sql
// Usada por: api-mobile, worker
//
// Nota: El contenido completo de las preguntas se almacena en MongoDB (material_assessment).
// Esta tabla solo mantiene metadata y referencia al documento MongoDB.
type Assessment struct {
	ID               uuid.UUID  `db:"id"`
	MaterialID       uuid.UUID  `db:"material_id"`
	MongoDocumentID  string     `db:"mongo_document_id"` // ObjectId de MongoDB
	QuestionsCount   int        `db:"questions_count"`   // Total de preguntas (columna canónica)
	Title            *string    `db:"title"`
	PassThreshold    *int       `db:"pass_threshold"`     // Porcentaje 0-100
	MaxAttempts      *int       `db:"max_attempts"`       // NULL = ilimitado
	TimeLimitMinutes *int       `db:"time_limit_minutes"` // NULL = sin límite
	Status           string     `db:"status"`             // draft, generated, published, archived, closed
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Assessment) TableName() string {
	return "assessment"
}
