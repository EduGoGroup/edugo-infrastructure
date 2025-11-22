package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAttempt representa la tabla 'assessment_attempt' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 007_create_assessment_attempts.up.sql, 010_extend_assessment_attempt.up.sql
// Usada por: api-mobile
type AssessmentAttempt struct {
	ID               uuid.UUID  `db:"id"`
	AssessmentID     uuid.UUID  `db:"assessment_id"`
	StudentID        uuid.UUID  `db:"student_id"`
	StartedAt        time.Time  `db:"started_at"`
	CompletedAt      *time.Time `db:"completed_at"`
	Score            *float64   `db:"score"`              // DECIMAL(5,2)
	MaxScore         *float64   `db:"max_score"`          // DECIMAL(5,2)
	Percentage       *float64   `db:"percentage"`         // DECIMAL(5,2)
	TimeSpentSeconds *int       `db:"time_spent_seconds"` // Tiempo total en segundos
	IdempotencyKey   *string    `db:"idempotency_key"`    // Prevención de duplicados
	Status           string     `db:"status"`             // in_progress, completed, abandoned
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttempt) TableName() string {
	return "assessment_attempt"
}
