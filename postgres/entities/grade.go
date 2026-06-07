package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Grade representa la tabla 'grades' en PostgreSQL.
//
// Migracion: 090_academic_grades.sql
type Grade struct {
	ID               uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MembershipID     uuid.UUID       `db:"membership_id" gorm:"type:uuid;index;not null;constraint:grades_membership_fkey,OnDelete:CASCADE;uniqueIndex:grades_unique" validate:"required,uuid"`
	SubjectID        uuid.UUID       `db:"subject_id" gorm:"type:uuid;index;not null;constraint:grades_subject_fkey,OnDelete:CASCADE;uniqueIndex:grades_unique" validate:"required,uuid"`
	PeriodID         uuid.UUID       `db:"period_id" gorm:"type:uuid;index;not null;constraint:grades_period_fkey,OnDelete:CASCADE;uniqueIndex:grades_unique" validate:"required,uuid"`
	GradeValue       *float64        `db:"grade_value" gorm:"type:decimal(5,2)" validate:"omitempty"`
	GradeLetter      *string         `db:"grade_letter" gorm:"type:varchar(5)" validate:"omitempty"`
	AssessmentScores json.RawMessage `db:"assessment_scores" gorm:"type:jsonb;default:'[]'"`
	// Source es la procedencia de la nota (N4 / ADR 0020): 'manual' (capturada por
	// el docente), 'auto_scored' (derivada de un intento auto-calificado) o
	// 'auto_llm' (generada por LLM). El CHECK es inline-expresable, asi que vive en
	// el tag GORM (mismo patron que schools.subscription_tier).
	Source string `db:"source" gorm:"not null;type:varchar(20);default:'manual';check:grades_source_check,source IN ('auto_scored','manual','auto_llm')" validate:"required,oneof=auto_scored manual auto_llm"`
	TeacherID        *uuid.UUID      `db:"teacher_id" gorm:"type:uuid;constraint:grades_teacher_fkey" validate:"omitempty,uuid"`
	Notes            *string         `db:"notes" validate:"omitempty"`
	FinalizedAt      *time.Time      `db:"finalized_at"`
	CreatedAt        time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

func (Grade) TableName() string {
	return "academic.grades"
}
