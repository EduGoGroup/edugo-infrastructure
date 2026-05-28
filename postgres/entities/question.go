package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Question representa la tabla 'questions' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 054_assessment_questions.sql
// Usada por: api-mobile, worker
type Question struct {
	ID            uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AssessmentID  uuid.UUID      `db:"assessment_id" gorm:"type:uuid;index;not null;constraint:questions_assessment_fkey,OnDelete:CASCADE;index:idx_questions_assessment" validate:"required,uuid"`
	SortOrder     int            `db:"sort_order" gorm:"not null;default:0;index:idx_questions_assessment" validate:"required"`
	QuestionText  string         `db:"question_text" gorm:"not null" validate:"required"`
	QuestionType  string         `db:"question_type" gorm:"not null;type:varchar(50);check:questions_question_type_check,question_type IN ('multiple_choice','true_false','short_answer','open_ended')" validate:"required,oneof=multiple_choice true_false short_answer open_ended"`
	CorrectAnswer *string        `db:"correct_answer" gorm:"default:null" validate:"omitempty"`
	Explanation   *string        `db:"explanation" gorm:"default:null" validate:"omitempty"`
	Points        float64        `db:"points" gorm:"type:numeric(5,2);not null;default:1" validate:"required"`
	Difficulty    *string        `db:"difficulty" gorm:"type:varchar(20);check:questions_difficulty_check,difficulty IN ('easy','medium','hard')" validate:"omitempty,oneof=easy medium hard"`
	Tags          pq.StringArray `db:"tags" gorm:"type:text[]" validate:"-"`
	CreatedAt     time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt     time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`

	Options []QuestionOption `gorm:"foreignKey:QuestionID" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Question) TableName() string {
	return "assessment.questions"
}
