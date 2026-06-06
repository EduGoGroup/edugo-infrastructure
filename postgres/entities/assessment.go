package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Assessment representa la tabla 'assessment.assessment' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD (N4 / ADR 0019).
//
// Evaluacion como artefacto autoral, reutilizable. Reemplaza el esquema viejo
// llaveado a auth.users + subject/grade texto-libre:
//   - created_by_user_id (→auth.users) → created_by_membership_id (→academic.memberships)
//   - subject/grade texto → subject_id FK (→academic.subjects, catalogo de escuela; ADR 0016)
//   - school_id pasa a NOT NULL (tenant del JWT)
//
// FKs cross-schema (school_id, created_by_membership_id, subject_id) se
// materializan en migrations/sql/post_gorm.sql (GORM no crea FK desde el tag
// `constraint:` sin campo de relacion). mongo_document_id queda reservado para
// V2 (source_type='ai_generated').
type Assessment struct {
	ID                    uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID              uuid.UUID      `db:"school_id" gorm:"type:uuid;not null;index;constraint:assessment_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	CreatedByMembershipID uuid.UUID      `db:"created_by_membership_id" gorm:"type:uuid;not null;index;constraint:assessment_created_by_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	SubjectID             uuid.UUID      `db:"subject_id" gorm:"type:uuid;not null;index;constraint:assessment_subject_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	Title                 string         `db:"title" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	Description           *string        `db:"description" gorm:"default:null" validate:"omitempty"`
	SourceType            string         `db:"source_type" gorm:"not null;type:varchar(20);default:'manual';check:assessment_source_type_check,source_type IN ('manual','ai_generated')" validate:"required,oneof=manual ai_generated"`
	Status                string         `db:"status" gorm:"not null;type:varchar(20);default:'draft';index;check:assessment_status_check,status IN ('draft','published','archived')" validate:"required,oneof=draft published archived"`
	QuestionsCount        int            `db:"questions_count" gorm:"not null;default:0"`
	PassThreshold         int            `db:"pass_threshold" gorm:"not null;default:70"`
	MaxAttempts           *int           `db:"max_attempts" gorm:"default:null" validate:"omitempty"`
	TimeLimitMinutes      *int           `db:"time_limit_minutes" gorm:"default:null" validate:"omitempty"`
	IsTimed               bool           `db:"is_timed" gorm:"not null;default:false"`
	ShuffleQuestions      bool           `db:"shuffle_questions" gorm:"not null;default:false"`
	ShowCorrectAnswers    bool           `db:"show_correct_answers" gorm:"not null;default:true"`
	AvailableFrom         *time.Time     `db:"available_from" gorm:"default:null"`
	AvailableUntil        *time.Time     `db:"available_until" gorm:"default:null"`
	MongoDocumentID       *string        `db:"mongo_document_id" gorm:"default:null;size:24;uniqueIndex:assessment_mongo_unique" validate:"omitempty"`
	CreatedAt             time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt             time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	// NOTE: indice parcial idx_assessment_active (WHERE deleted_at IS NULL) en post_gorm.sql
	DeletedAt gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Assessment) TableName() string {
	return "assessment.assessment"
}
