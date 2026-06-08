package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Question representa la tabla 'assessment.question' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD (N4 / ADR 0019).
//
// Pregunta de una evaluacion. Sin cambio estructural respecto al viejo (el
// modelo era valido): la opcion correcta se referencia en correct_answer
// (NO hay columna is_correct en la opcion). FK assessment_id en post_gorm.sql.
type Question struct {
	ID            uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AssessmentID  uuid.UUID      `db:"assessment_id" gorm:"type:uuid;not null;constraint:question_assessment_fkey,OnDelete:CASCADE;index:idx_question_assessment_order,priority:1" validate:"required,uuid"`
	SortOrder     int            `db:"sort_order" gorm:"not null;default:0;index:idx_question_assessment_order,priority:2"`
	QuestionText  string         `db:"question_text" gorm:"not null" validate:"required"`
	QuestionType  string         `db:"question_type" gorm:"not null;type:varchar(20);check:question_type_check,question_type IN ('multiple_choice','multiple_select','true_false','short_answer','open_ended')" validate:"required,oneof=multiple_choice multiple_select true_false short_answer open_ended"`
	CorrectAnswer *string        `db:"correct_answer" gorm:"default:null" validate:"omitempty"`
	Explanation   *string        `db:"explanation" gorm:"default:null" validate:"omitempty"`
	Points        int            `db:"points" gorm:"not null;default:1"`
	Difficulty    *string        `db:"difficulty" gorm:"type:varchar(20);default:null" validate:"omitempty"`
	Tags          pq.StringArray `db:"tags" gorm:"type:text[]" validate:"-"`
	CreatedAt     time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt     time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Question) TableName() string {
	return "assessment.question"
}
