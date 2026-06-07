package entities

import (
	"time"

	"github.com/google/uuid"
)

// GradeItem representa la tabla 'academic.grade_item' en PostgreSQL (N4 / ADR 0020).
//
// Componente de nota: una pieza individual que contribuye a la nota de un alumno
// en una sesion (materia + periodo). El perfil de escuela 'detailed' usa varios
// grade_item por (membership, subject, period); el perfil 'basic' usa solo
// academic.grades. La procedencia (Source) distingue notas capturadas a mano,
// derivadas de un intento auto-calificado o generadas por LLM.
//
// El grain (membership_id, subject_id, period_id) NO es unico (un alumno tiene
// varios componentes por sesion); se indexa via idx_grade_item_grain. El UNIQUE
// PARCIAL uq_grade_item_attempt (membership_id, subject_id, period_id,
// source_attempt_id) WHERE source_attempt_id IS NOT NULL previene duplicar el
// componente auto_scored derivado del mismo intento — vive en post_gorm.sql
// (GORM no expresa indices parciales con WHERE).
//
// FKs: las academic (membership/subject/period CASCADE, created_by RESTRICT) y
// las cross-schema a assessment.* (source_attempt_id/source_assessment_id, ambas
// SET NULL) se materializan en migrations/sql/post_gorm.sql, porque GORM no crea
// FKs desde el tag `constraint:` sin campo de relacion (mismo patron que
// academic.grades.teacher_id y assessment.*).
type GradeItem struct {
	ID           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MembershipID uuid.UUID `db:"membership_id" gorm:"type:uuid;not null;constraint:grade_item_membership_fkey,OnDelete:CASCADE;index:idx_grade_item_grain,priority:1" validate:"required,uuid"`
	SubjectID    uuid.UUID `db:"subject_id" gorm:"type:uuid;not null;constraint:grade_item_subject_fkey,OnDelete:CASCADE;index:idx_grade_item_grain,priority:2" validate:"required,uuid"`
	PeriodID     uuid.UUID `db:"period_id" gorm:"type:uuid;not null;constraint:grade_item_period_fkey,OnDelete:CASCADE;index:idx_grade_item_grain,priority:3" validate:"required,uuid"`
	Title        string    `db:"title" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	Value        *float64  `db:"value" gorm:"type:decimal(5,2)" validate:"omitempty"`
	// Weight es informativo en la generacion 1 (no se usa para ponderar todavia).
	Weight *float64 `db:"weight" gorm:"type:decimal(5,2)" validate:"omitempty"`
	// Source: procedencia del componente. CHECK inline en el tag GORM (mismo patron
	// que grades.source / schools.subscription_tier).
	Source string `db:"source" gorm:"not null;type:varchar(20);default:'manual';check:grade_item_source_check,source IN ('auto_scored','manual','auto_llm')" validate:"required,oneof=auto_scored manual auto_llm"`
	// SourceAttemptID / SourceAssessmentID: trazabilidad al origen auto_scored/
	// auto_llm (intento y evaluacion de assessment.*). Nullable (manual no tiene
	// origen); FK cross-schema SET NULL en post_gorm.sql.
	SourceAttemptID       *uuid.UUID `db:"source_attempt_id" gorm:"type:uuid;index:idx_grade_item_attempt;constraint:grade_item_source_attempt_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	SourceAssessmentID    *uuid.UUID `db:"source_assessment_id" gorm:"type:uuid;constraint:grade_item_source_assessment_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	CreatedByMembershipID uuid.UUID  `db:"created_by_membership_id" gorm:"type:uuid;not null;constraint:grade_item_created_by_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	CreatedAt             time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt             time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (GradeItem) TableName() string {
	return "academic.grade_item"
}
