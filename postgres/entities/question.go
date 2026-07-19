package entities

import (
	"encoding/json"
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
	// LLMPrep es el artefacto de preparacion para el LLM (contrato versionado v1,
	// plan 042 D-042.2). JSONB aditivo NULL; solo lo consume el carril LLM, no SQL.
	LLMPrep json.RawMessage `db:"llm_prep" gorm:"type:jsonb;default:null"`
	// LLMPrepStatus es el estado del artefacto de preparacion (plan 042 D-042.1).
	LLMPrepStatus string `db:"llm_prep_status" gorm:"not null;type:varchar(20);default:'none';check:question_llm_prep_status_check,llm_prep_status IN ('none','pending','processing','ready','failed','stale')" validate:"omitempty,oneof=none pending processing ready failed stale"`
	// LLMPrepSourceHash es el sha256 de question_type+question_text+correct_answer+explanation
	// que ancla la concurrencia optimista de la preparacion (plan 042 D-042.5). NULL = sin prep aun.
	LLMPrepSourceHash *string `db:"llm_prep_source_hash" gorm:"type:varchar(64);default:null" validate:"omitempty"`
	// LLMPrepFeedback es el comentario del profesor para re-preparar; se consume 1 vez
	// y se limpia al procesarlo (plan 042 D-042.7). NULL = sin comentario pendiente.
	LLMPrepFeedback *string `db:"llm_prep_feedback" gorm:"type:text;default:null" validate:"omitempty"`
	// LLMPrepUpdatedAt marca cuando se escribio el prep por ultima vez. timestamptz aditivo NULL.
	LLMPrepUpdatedAt *time.Time `db:"llm_prep_updated_at" gorm:"default:null"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Question) TableName() string {
	return "assessment.question"
}
