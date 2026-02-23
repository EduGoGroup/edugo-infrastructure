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
	ID               uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	AttemptID        uuid.UUID `db:"attempt_id" gorm:"type:uuid;index;not null"`
	QuestionIndex    int       `db:"question_index" gorm:"not null"`
	StudentAnswer    *string   `db:"student_answer" gorm:"default:null"`
	IsCorrect        *bool     `db:"is_correct" gorm:"default:null"`
	PointsEarned     *float64  `db:"points_earned" gorm:"type:decimal(5,2)"`
	MaxPoints        *float64  `db:"max_points" gorm:"type:decimal(5,2)"`
	TimeSpentSeconds *int      `db:"time_spent_seconds" gorm:"default:null"`
	AnsweredAt       time.Time `db:"answered_at" gorm:"not null"`
	CreatedAt        time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt        time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttemptAnswer) TableName() string {
	return "assessment.assessment_attempt_answer"
}
