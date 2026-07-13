package entities

import (
	"time"

	"github.com/google/uuid"
)

// PracticeSession representa la tabla 'assessment.practice_session' en PostgreSQL
// (plan 035 D-035.4 / §Modelo de datos). Es la CABECERA del log de practica: una
// fila por sesion que un alumno abre sobre una evaluacion con purpose in
// (practice, both). El log es APPEND-ONLY y PRESCINDIBLE por diseño: nada
// funcional lo lee (la logica adaptativa F2 se alimenta del acumulador
// user_question_stat), por lo que es podable sin culpa (≥6 meses / N sesiones).
//
// assessment_id es NULLABLE con ON DELETE SET NULL (D-035.4): el historial de
// practica NO cuelga del assessment; si el profesor lo borra, la sesion sobrevive
// huerfana. La escritura es SINCRONA en learning, sin worker (D-035.3).
//
// FKs cross-schema (school_id, membership_id, subject_id) y la de assessment_id
// (SET NULL) se materializan en migrations/sql/post_gorm.sql (GORM no crea FK
// desde el tag `constraint:` sin campo de relacion).
type PracticeSession struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID         uuid.UUID  `db:"school_id" gorm:"type:uuid;not null;index" validate:"required,uuid"`
	MembershipID     uuid.UUID  `db:"membership_id" gorm:"type:uuid;not null;index" validate:"required,uuid"`
	AssessmentID     *uuid.UUID `db:"assessment_id" gorm:"type:uuid;default:null;index" validate:"omitempty"`
	SubjectID        uuid.UUID  `db:"subject_id" gorm:"type:uuid;not null;index" validate:"required,uuid"`
	StartedAt        time.Time  `db:"started_at" gorm:"not null;default:now()"`
	FinishedAt       *time.Time `db:"finished_at" gorm:"default:null"`
	QuestionsTotal   int        `db:"questions_total" gorm:"not null;default:0"`
	QuestionsCorrect int        `db:"questions_correct" gorm:"not null;default:0"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (PracticeSession) TableName() string {
	return "assessment.practice_session"
}
