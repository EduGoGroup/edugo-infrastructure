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
	ID               uuid.UUID `db:"id"`
	AttemptID        uuid.UUID `db:"attempt_id"`
	QuestionIndex    int       `db:"question_index"` // 0-based index
	StudentAnswer    *string   `db:"student_answer"` // TEXT flexible: JSON, string, etc
	IsCorrect        *bool     `db:"is_correct"`
	PointsEarned     *float64  `db:"points_earned"`      // DECIMAL(5,2)
	MaxPoints        *float64  `db:"max_points"`         // DECIMAL(5,2)
	TimeSpentSeconds *int      `db:"time_spent_seconds"` // Tiempo en esta pregunta
	AnsweredAt       time.Time `db:"answered_at"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttemptAnswer) TableName() string {
	return "assessment_attempt_answer"
}
