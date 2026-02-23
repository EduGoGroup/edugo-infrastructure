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
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	AssessmentID     uuid.UUID  `db:"assessment_id" gorm:"type:uuid;index;not null"`
	StudentID        uuid.UUID  `db:"student_id" gorm:"type:uuid;index;not null"`
	StartedAt        time.Time  `db:"started_at" gorm:"not null"`
	CompletedAt      *time.Time `db:"completed_at" gorm:"default:null"`
	Score            *float64   `db:"score" gorm:"type:decimal(5,2)"`
	MaxScore         *float64   `db:"max_score" gorm:"type:decimal(5,2)"`
	Percentage       *float64   `db:"percentage" gorm:"type:decimal(5,2)"`
	TimeSpentSeconds *int       `db:"time_spent_seconds" gorm:"default:null"`
	IdempotencyKey   *string    `db:"idempotency_key" gorm:"uniqueIndex"`
	Status           string     `db:"status" gorm:"not null;type:varchar(50)"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttempt) TableName() string {
	return "assessment.assessment_attempt"
}
