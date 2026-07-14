package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAttemptAnswer representa la tabla 'assessment.assessment_attempt_answer'
// en PostgreSQL (N4 / ADR 0019). Respuesta del alumno a una pregunta dentro de un
// intento.
//
// FKs (attempt_id → assessment_attempt CASCADE, question_id → question SET NULL) y
// el UNIQUE (attempt_id, question_index) se materializan en post_gorm.sql.
type AssessmentAttemptAnswer struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AttemptID        uuid.UUID  `db:"attempt_id" gorm:"type:uuid;not null;index;constraint:assessment_attempt_answer_attempt_fkey,OnDelete:CASCADE;uniqueIndex:uq_attempt_answer_question,priority:1" validate:"required,uuid"`
	QuestionIndex    int        `db:"question_index" gorm:"not null;uniqueIndex:uq_attempt_answer_question,priority:2" validate:"required"`
	QuestionID       *uuid.UUID `db:"question_id" gorm:"type:uuid;constraint:assessment_attempt_answer_question_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	StudentAnswer    *string    `db:"student_answer" gorm:"default:null" validate:"omitempty"`
	ReviewStatus     string     `db:"review_status" gorm:"not null;type:varchar(20);default:'pending';check:assessment_attempt_answer_review_status_check,review_status IN ('pending','auto_graded','reviewed','ai_reviewed')" validate:"required,oneof=pending auto_graded reviewed ai_reviewed"`
	IsCorrect        *bool      `db:"is_correct" gorm:"default:null"`
	PointsEarned     *float64   `db:"points_earned" gorm:"type:decimal(5,2)" validate:"omitempty"`
	MaxPoints        *float64   `db:"max_points" gorm:"type:decimal(5,2)" validate:"omitempty"`
	TimeSpentSeconds *int       `db:"time_spent_seconds" gorm:"default:null" validate:"omitempty"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttemptAnswer) TableName() string {
	return "assessment.assessment_attempt_answer"
}
