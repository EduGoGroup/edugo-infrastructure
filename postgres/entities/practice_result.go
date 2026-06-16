package entities

import (
	"time"

	"github.com/google/uuid"
)

// PracticeResult representa la tabla 'academic.practice_result' en PostgreSQL
// (plan 024 F6). Es el ESPEJO de academic.grade_item para evaluaciones de tipo
// 'practice': su resultado NO va al expediente (no genera grade_item), se guarda
// aqui solo para estadisticas. El worker decide en cual tabla escribir segun el
// `kind` de la evaluacion ('final' -> grade_item, 'practice' -> practice_result).
//
// CLAVE DE CONFLICTO DEL UPSERT: el worker genera un `id` determinista derivado
// del intento de origen y hace ON CONFLICT (id). La PK(id) soporta ese upsert.
// Ademas, como defensa en profundidad (paridad con grade_item), el UNIQUE PARCIAL
// uq_practice_result_attempt (membership_id, subject_id, period_id,
// source_attempt_id) WHERE source_attempt_id IS NOT NULL previene duplicar el
// resultado auto_scored derivado del mismo intento — vive en post_gorm.sql (GORM
// no expresa indices parciales con WHERE).
//
// El grain (membership_id, subject_id, period_id) NO es unico (un alumno puede
// tener varios resultados practicos por sesion); se indexa via
// idx_practice_result_grain.
//
// FKs: las academic (membership/subject/period CASCADE, created_by RESTRICT) y
// las cross-schema a assessment.* (source_attempt_id/source_assessment_id, ambas
// SET NULL) se materializan en migrations/sql/post_gorm.sql, porque GORM no crea
// FKs desde el tag `constraint:` sin campo de relacion (mismo patron que
// academic.grade_item).
type PracticeResult struct {
	ID           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MembershipID uuid.UUID `db:"membership_id" gorm:"type:uuid;not null;constraint:practice_result_membership_fkey,OnDelete:CASCADE;index:idx_practice_result_grain,priority:1" validate:"required,uuid"`
	SubjectID    uuid.UUID `db:"subject_id" gorm:"type:uuid;not null;constraint:practice_result_subject_fkey,OnDelete:CASCADE;index:idx_practice_result_grain,priority:2" validate:"required,uuid"`
	PeriodID     uuid.UUID `db:"period_id" gorm:"type:uuid;not null;constraint:practice_result_period_fkey,OnDelete:CASCADE;index:idx_practice_result_grain,priority:3" validate:"required,uuid"`
	Title        string    `db:"title" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	// Value: % del intento (0–100). Nullable (paridad con grade_item.value).
	Value *float64 `db:"value" gorm:"type:decimal(5,2)" validate:"omitempty"`
	// Source: procedencia del resultado. CHECK inline en el tag GORM (mismo patron
	// que grade_item.source).
	Source string `db:"source" gorm:"not null;type:varchar(20);default:'auto_scored';check:practice_result_source_check,source IN ('auto_scored','manual','auto_llm')" validate:"required,oneof=auto_scored manual auto_llm"`
	// SourceAttemptID / SourceAssessmentID: trazabilidad al intento y evaluacion
	// de origen (assessment.*). Nullable; FK cross-schema SET NULL en post_gorm.sql.
	SourceAttemptID       *uuid.UUID `db:"source_attempt_id" gorm:"type:uuid;index:idx_practice_result_attempt;constraint:practice_result_source_attempt_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	SourceAssessmentID    *uuid.UUID `db:"source_assessment_id" gorm:"type:uuid;constraint:practice_result_source_assessment_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	CreatedByMembershipID uuid.UUID  `db:"created_by_membership_id" gorm:"type:uuid;not null;constraint:practice_result_created_by_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	CreatedAt             time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt             time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (PracticeResult) TableName() string {
	return "academic.practice_result"
}
