package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Progress representa la tabla 'content.progress' en PostgreSQL (N4 / ADR 0019).
// Progreso de lectura de un material por alumno.
//
// Cambio vs viejo: PK compuesta (material_id, user_id) → (material_id,
// student_membership_id) (re-llaveado a academic.memberships).
// completed_at IS NOT NULL indica material completado. FKs (material_id →
// content.materials CASCADE, student_membership_id → academic.memberships CASCADE)
// se materializan en post_gorm.sql.
type Progress struct {
	MaterialID          uuid.UUID       `db:"material_id" gorm:"type:uuid;primaryKey;constraint:progress_material_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	StudentMembershipID uuid.UUID       `db:"student_membership_id" gorm:"type:uuid;primaryKey;constraint:progress_student_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	ProgressPercentage  int             `db:"progress_percentage" gorm:"not null;default:0"`
	LastPosition        json.RawMessage `db:"last_position" gorm:"type:jsonb;default:null"`
	CompletedAt         *time.Time      `db:"completed_at" gorm:"default:null"`
	CreatedAt           time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt           time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL.
func (Progress) TableName() string {
	return "content.progress"
}
