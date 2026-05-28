package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AssessmentAttempt representa la tabla 'assessment_attempt' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 051_assessment_assessment_attempt.sql
// Usada por: api-mobile
type AssessmentAttempt struct {
	ID               uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AssessmentID     uuid.UUID       `db:"assessment_id" gorm:"type:uuid;index;not null;constraint:assessment_attempt_assessment_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	StudentID        uuid.UUID       `db:"student_id" gorm:"type:uuid;index;not null;constraint:assessment_attempt_student_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	StartedAt        time.Time       `db:"started_at" gorm:"not null"`
	CompletedAt      *time.Time      `db:"completed_at" gorm:"default:null;check:check_attempt_time_logical,completed_at IS NULL OR completed_at > started_at"`
	Score            *float64        `db:"score" gorm:"type:decimal(5,2)" validate:"omitempty"`
	MaxScore         *float64        `db:"max_score" gorm:"type:decimal(5,2)" validate:"omitempty"`
	Percentage       *float64        `db:"percentage" gorm:"type:decimal(5,2)" validate:"omitempty"`
	QuestionOrder    json.RawMessage `db:"question_order" gorm:"type:jsonb;default:null"`
	TimeSpentSeconds *int            `db:"time_spent_seconds" gorm:"default:null;check:assessment_attempt_time_spent_seconds_check,time_spent_seconds IS NULL OR (time_spent_seconds > 0 AND time_spent_seconds <= 7200)" validate:"omitempty"`
	IdempotencyKey   *string         `db:"idempotency_key" gorm:"uniqueIndex;size:64" validate:"omitempty"`
	// NOTE: partial index idx_attempt_completed (WHERE status='completed') must be created in post_gorm.sql
	// NOTE: partial index idx_attempt_pending_review (WHERE status='pending_review') must be created in post_gorm.sql
	Status           string          `db:"status" gorm:"not null;type:varchar(50);check:assessment_attempt_status_check,status IN ('in_progress','completed','abandoned','pending_review');index:idx_assessment_attempt_status;default:'in_progress'" validate:"required,oneof=in_progress completed abandoned pending_review"`
	CreatedAt        time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttempt) TableName() string {
	return "assessment.assessment_attempt"
}
