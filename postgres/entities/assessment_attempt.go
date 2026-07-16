package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AssessmentAttempt representa la tabla 'assessment.assessment_attempt' en
// PostgreSQL (N4 / ADR 0019). Es el HECHO: un intento de evaluacion por membresia.
//
// Cambio NUCLEAR vs viejo: student_id (→auth.users) → student_membership_id
// (→academic.memberships). La seccion se DERIVA del enrollment de la oferta; no
// hay offering_id en el hecho (premisa ADR 0018).
//
// FKs cross-schema (assessment_id, student_membership_id) y los indices parciales
// (idx_attempt_completed, idx_attempt_pending_review, y el UNIQUE parcial de un
// solo intento activo por (assessment_id, student_membership_id) WHERE
// status='in_progress') se materializan en post_gorm.sql.
type AssessmentAttempt struct {
	ID                  uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AssessmentID        uuid.UUID       `db:"assessment_id" gorm:"type:uuid;not null;index;constraint:assessment_attempt_assessment_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	StudentMembershipID uuid.UUID       `db:"student_membership_id" gorm:"type:uuid;not null;index;constraint:assessment_attempt_student_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	StartedAt           time.Time       `db:"started_at" gorm:"not null"`
	CompletedAt         *time.Time      `db:"completed_at" gorm:"default:null;check:assessment_attempt_time_check,completed_at IS NULL OR completed_at > started_at"`
	Score               *float64        `db:"score" gorm:"type:decimal(5,2)" validate:"omitempty"`
	MaxScore            *float64        `db:"max_score" gorm:"type:decimal(5,2)" validate:"omitempty"`
	Percentage          *float64        `db:"percentage" gorm:"type:decimal(5,2)" validate:"omitempty"`
	QuestionOrder       json.RawMessage `db:"question_order" gorm:"type:jsonb;default:null"`
	TimeSpentSeconds    *int            `db:"time_spent_seconds" gorm:"default:null;check:assessment_attempt_time_spent_check,time_spent_seconds IS NULL OR (time_spent_seconds > 0 AND time_spent_seconds <= 7200)" validate:"omitempty"`
	IdempotencyKey      *string         `db:"idempotency_key" gorm:"uniqueIndex;size:64" validate:"omitempty"`
	Status              string          `db:"status" gorm:"not null;type:varchar(50);index;default:'in_progress';check:assessment_attempt_status_check,status IN ('in_progress','completed','abandoned','pending_review')" validate:"required,oneof=in_progress completed abandoned pending_review"`
	// TeacherFeedback es el comentario global del profesor al finalizar la revision (plan 036 D-036.4).
	TeacherFeedback *string `db:"teacher_feedback" gorm:"type:text;default:null" validate:"omitempty"`
	// AIReviewClaimedAt marca cuando un proceso de revision por IA tomo el candado
	// «en revision por IA» sobre este intento (candado con vencimiento, T4-1).
	// NULL = sin candado activo. timestamptz aditivo.
	AIReviewClaimedAt *time.Time `db:"ai_review_claimed_at" gorm:"default:null"`
	CreatedAt         time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt         time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAttempt) TableName() string {
	return "assessment.assessment_attempt"
}
