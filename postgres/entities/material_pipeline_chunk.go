package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MaterialPipelineChunk representa la tabla 'content.material_pipeline_chunk' en
// PostgreSQL (plan 043 F0). Es el DETALLE del job: el material se parte en trozos
// (chunk_text) que se procesan de forma independiente (summary + artifacts). El
// par (job_id, seq) es único: cada trozo ocupa una posición fija dentro del job.
//
// La FK same-schema job_id → content.material_pipeline_job(id) ON DELETE CASCADE
// se declara en post_gorm.sql (GORM no materializa FKs desde el tag `constraint:`
// sin campo de relación). El UNIQUE (job_id, seq) sí lo materializa GORM desde el
// tag `uniqueIndex`.
type MaterialPipelineChunk struct {
	ID        uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	JobID     uuid.UUID       `db:"job_id" gorm:"type:uuid;index;not null;constraint:material_pipeline_chunk_job_fkey,OnDelete:CASCADE;uniqueIndex:uq_material_pipeline_chunk_job_seq,priority:1" validate:"required,uuid"`
	Seq       int             `db:"seq" gorm:"not null;uniqueIndex:uq_material_pipeline_chunk_job_seq,priority:2" validate:"required"`
	ChunkText string          `db:"chunk_text" gorm:"type:text;not null" validate:"required"`
	Summary   *string         `db:"summary" gorm:"type:text;default:null" validate:"omitempty"`
	Artifacts json.RawMessage `db:"artifacts" gorm:"type:jsonb;default:null"`
	Status    string          `db:"status" gorm:"not null;type:varchar(20);index;default:'pending';check:material_pipeline_chunk_status_check,status IN ('pending','processing','done','failed')" validate:"required,oneof=pending processing done failed"`
	Attempts  int             `db:"attempts" gorm:"not null;default:0"`
	UpdatedAt time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialPipelineChunk) TableName() string {
	return "content.material_pipeline_chunk"
}
