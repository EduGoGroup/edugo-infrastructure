package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MaterialPipelineCandidate representa la tabla
// 'content.material_pipeline_candidate' en PostgreSQL (plan 043 F0). Es una
// pregunta candidata generada a partir de un chunk: el payload lleva la propuesta
// (JSONB), embedding su vector opcional para deduplicar, y status marca el
// destino de la candidata (candidata, descartada por duplicado/irrelevante, o
// seleccionada). dedupe_group agrupa candidatas equivalentes; score pondera la
// selección.
//
// Las FKs same-schema job_id → content.material_pipeline_job(id) y chunk_id →
// content.material_pipeline_chunk(id) (ambas ON DELETE CASCADE) se declaran en
// post_gorm.sql (GORM no materializa FKs desde el tag `constraint:` sin campo de
// relación).
type MaterialPipelineCandidate struct {
	ID          uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	JobID       uuid.UUID       `db:"job_id" gorm:"type:uuid;index;not null;constraint:material_pipeline_candidate_job_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	ChunkID     uuid.UUID       `db:"chunk_id" gorm:"type:uuid;index;not null;constraint:material_pipeline_candidate_chunk_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	Payload     json.RawMessage `db:"payload" gorm:"type:jsonb;not null"`
	Embedding   json.RawMessage `db:"embedding" gorm:"type:jsonb;default:null"`
	Status      string          `db:"status" gorm:"not null;type:varchar(30);index;default:'candidate';check:material_pipeline_candidate_status_check,status IN ('candidate','dropped_dup','dropped_irrelevant','selected')" validate:"required,oneof=candidate dropped_dup dropped_irrelevant selected"`
	DedupeGroup *uuid.UUID      `db:"dedupe_group" gorm:"type:uuid;default:null" validate:"omitempty,uuid"`
	Score       *float64        `db:"score" gorm:"type:numeric;default:null" validate:"omitempty"`
	CreatedAt   time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialPipelineCandidate) TableName() string {
	return "content.material_pipeline_candidate"
}
