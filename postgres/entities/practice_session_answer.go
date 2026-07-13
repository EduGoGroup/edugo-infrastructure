package entities

import (
	"time"

	"github.com/google/uuid"
)

// PracticeSessionAnswer representa la tabla 'assessment.practice_session_answer'
// en PostgreSQL (plan 035 D-035.4 / §Modelo de datos). Es el DETALLE del log de
// practica: una fila por respuesta dentro de una sesion. APPEND-ONLY y
// PRESCINDIBLE por diseño (nada funcional lo lee; el acumulador
// user_question_stat ya alimenta F2); podable ≥6 meses sin culpa.
//
// session_id lleva FK CASCADE a practice_session; question_id es NULLABLE con ON
// DELETE SET NULL (el historial sobrevive al borrado de la pregunta). Ambas se
// materializan en migrations/sql/post_gorm.sql (GORM no crea FK desde el tag
// `constraint:` sin campo de relacion). No lleva updated_at: es append-only, solo
// answered_at.
type PracticeSessionAnswer struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SessionID        uuid.UUID  `db:"session_id" gorm:"type:uuid;not null;index" validate:"required,uuid"`
	QuestionID       *uuid.UUID `db:"question_id" gorm:"type:uuid;default:null;index" validate:"omitempty"`
	QuestionIndex    int        `db:"question_index" gorm:"not null;default:0"`
	IsCorrect        bool       `db:"is_correct" gorm:"not null"`
	TimeSpentSeconds *int       `db:"time_spent_seconds" gorm:"default:null" validate:"omitempty"`
	AnsweredAt       time.Time  `db:"answered_at" gorm:"not null;default:now()"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (PracticeSessionAnswer) TableName() string {
	return "assessment.practice_session_answer"
}
