package entities

import (
	"time"

	"github.com/google/uuid"
)

// AttemptReview representa la tabla 'assessment.attempt_review' en PostgreSQL
// (N4 / ADR 0019). Calificacion manual de respuestas abiertas, por membresia.
//
// Cambio vs viejo: reviewer_id (global) → reviewer_membership_id (→academic.memberships).
// FKs (attempt_answer_id → assessment_attempt_answer CASCADE [UNIQUE],
// reviewer_membership_id → academic.memberships RESTRICT) se materializan en
// post_gorm.sql.
type AttemptReview struct {
	ID                   uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AttemptAnswerID      uuid.UUID `db:"attempt_answer_id" gorm:"type:uuid;not null;uniqueIndex:uq_attempt_review_answer;constraint:attempt_review_answer_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	ReviewerMembershipID uuid.UUID `db:"reviewer_membership_id" gorm:"type:uuid;not null;constraint:attempt_review_reviewer_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	PointsAwarded        float64   `db:"points_awarded" gorm:"type:decimal(5,2);not null" validate:"required"`
	Feedback             *string   `db:"feedback" gorm:"default:null" validate:"omitempty"`
	CreatedAt            time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt            time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AttemptReview) TableName() string {
	return "assessment.attempt_review"
}
