package entities

import (
	"time"

	"github.com/google/uuid"
)

// QuestionOption representa la tabla 'question_options' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 055_assessment_question_options.sql
// Usada por: api-mobile, worker
type QuestionOption struct {
	ID         uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	QuestionID uuid.UUID `db:"question_id" gorm:"type:uuid;index;not null;constraint:question_options_question_fkey,OnDelete:CASCADE;index:idx_options_question" validate:"required,uuid"`
	OptionText string    `db:"option_text" gorm:"not null" validate:"required"`
	SortOrder  int       `db:"sort_order" gorm:"not null;default:0;index:idx_options_question" validate:"required"`
	CreatedAt  time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (QuestionOption) TableName() string {
	return "assessment.question_options"
}
