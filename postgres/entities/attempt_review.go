package entities

import (
	"time"

	"github.com/google/uuid"
)

// AttemptReview representa la tabla 'attempt_reviews' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 057_assessment_attempt_reviews.sql
// Usada por: api-mobile
type AttemptReview struct {
	ID              uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AttemptAnswerID uuid.UUID `db:"attempt_answer_id" gorm:"type:uuid;uniqueIndex:unique_review;not null;constraint:attempt_reviews_answer_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	ReviewerID      uuid.UUID `db:"reviewer_id" gorm:"type:uuid;not null;constraint:attempt_reviews_reviewer_fkey" validate:"required,uuid"`
	PointsAwarded   float64   `db:"points_awarded" gorm:"type:numeric(5,2);not null" validate:"required"`
	Feedback        *string   `db:"feedback" gorm:"default:null" validate:"omitempty"`
	ReviewedAt      time.Time `db:"reviewed_at" gorm:"not null;autoCreateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AttemptReview) TableName() string {
	return "assessment.attempt_reviews"
}
