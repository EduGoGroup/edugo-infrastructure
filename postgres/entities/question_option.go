package entities

import (
	"time"

	"github.com/google/uuid"
)

// QuestionOption representa la tabla 'assessment.question_option' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD (N4 / ADR 0019).
//
// Opcion de una pregunta. La "correcta" NO se marca aqui: se referencia desde
// question.correct_answer (decision F1, default del esquema). FK question_id en
// post_gorm.sql.
type QuestionOption struct {
	ID         uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	QuestionID uuid.UUID `db:"question_id" gorm:"type:uuid;not null;constraint:question_option_question_fkey,OnDelete:CASCADE;index:idx_question_option_order,priority:1" validate:"required,uuid"`
	OptionText string    `db:"option_text" gorm:"not null" validate:"required"`
	SortOrder  int       `db:"sort_order" gorm:"not null;default:0;index:idx_question_option_order,priority:2"`
	CreatedAt  time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt  time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (QuestionOption) TableName() string {
	return "assessment.question_option"
}
