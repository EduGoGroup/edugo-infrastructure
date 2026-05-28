package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAttemptAnswer representa la tabla 'assessment_attempt_answer' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 008_create_assessment_answers.up.sql, 011_extend_assessment_answer.up.sql
// Usada por: api-mobile
type AssessmentAttemptAnswer struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AttemptID        uuid.UUID  `db:"attempt_id" gorm:"type:uuid;index;not null;constraint:assessment_attempt_answer_attempt_fkey,OnDelete:CASCADE;uniqueIndex:assessment_attempt_answer_unique_question" validate:"required,uuid"`
	QuestionIndex    int        `db:"question_index" gorm:"not null;uniqueIndex:assessment_attempt_answer_unique_question" validate:"required"`
	QuestionID       *uuid.UUID `db:"question_id" gorm:"type:uuid;constraint:assessment_attempt_answer_question_fkey" validate:"omitempty,uuid"`
	StudentAnswer    *string    `db:"student_answer" gorm:"default:null" validate:"omitempty"`
	ReviewStatus     *string    `db:"review_status" gorm:"type:varchar(20);default:'pending';check:assessment_attempt_answer_review_status_check,review_status IN ('pending','auto_graded','reviewed')" validate:"omitempty,oneof=pending auto_graded reviewed"`
	IsCorrect        *bool      `db:"is_correct" gorm:"default:null"`
	PointsEarned     *float64   `db:"points_earned" gorm:"type:decimal(5,2)" validate:"omitempty"`
	MaxPoints        *float64   `db:"max_points" gorm:"type:decimal(5,2)" validate:"omitempty"`
	TimeSpentSeconds *int       `db:"time_spent_seconds" gorm:"default:null;check:assessment_attempt_answer_time_spent_seconds_check,time_spent_seconds >= 0" validate:"omitempty"`
	AnsweredAt       time.Time  `db:"answered_at" gorm:"not null"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttemptAnswer) TableName() string {
	return "assessment.assessment_attempt_answer"
}
