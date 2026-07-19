package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MaterialPipelineJob representa la tabla 'content.material_pipeline_job' en
// PostgreSQL (plan 043 F0: pipeline material → evaluación). Es la CABECERA del
// trabajo de generación: un material que el docente manda a convertir en una
// evaluación pasa por un job con fases (phase) y estados (status). El resultado,
// cuando termina, apunta a la evaluación creada (assessment_id).
//
// FKs cross-schema (school_id → academic.schools, requested_by_membership_id →
// academic.memberships) NO se materializan como FK dura (mismo criterio que
// content.user_material_tags): school_id/requested_by_membership_id quedan como
// UUID indexado. La FK same-schema material_id → content.materials(id) se
// declara en post_gorm.sql (GORM no materializa FKs desde el tag `constraint:`
// sin campo de relación, mismo caso que content.material_assignment).
type MaterialPipelineJob struct {
	ID                      uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MaterialID              uuid.UUID       `db:"material_id" gorm:"type:uuid;index;not null;constraint:material_pipeline_job_material_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	SchoolID                uuid.UUID       `db:"school_id" gorm:"type:uuid;index;not null" validate:"required,uuid"`
	RequestedByMembershipID uuid.UUID       `db:"requested_by_membership_id" gorm:"type:uuid;index;not null" validate:"required,uuid"`
	Status                  string          `db:"status" gorm:"not null;type:varchar(20);index;default:'pending';check:material_pipeline_job_status_check,status IN ('pending','processing','done','failed')" validate:"required,oneof=pending processing done failed"`
	Phase                   int16           `db:"phase" gorm:"not null;type:smallint;default:0"`
	Params                  json.RawMessage `db:"params" gorm:"type:jsonb;default:null"`
	AssessmentID            *uuid.UUID      `db:"assessment_id" gorm:"type:uuid;default:null" validate:"omitempty,uuid"`
	LastError               *string         `db:"last_error" gorm:"type:text;default:null" validate:"omitempty"`
	CreatedAt               time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt               time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	CompletedAt             *time.Time      `db:"completed_at" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialPipelineJob) TableName() string {
	return "content.material_pipeline_job"
}
