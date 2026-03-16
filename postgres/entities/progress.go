package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Progress representa la tabla 'content.progress' en PostgreSQL.
// Reflejo exacto del schema definido en: postgres/migrations/structure/042_content_progress.sql
//
// Clave primaria compuesta (material_id, user_id) — un registro por material/usuario.
// completed_at IS NOT NULL indica que el usuario completó el material.
type Progress struct {
	MaterialID         uuid.UUID       `db:"material_id" gorm:"type:uuid;primaryKey"`
	UserID             uuid.UUID       `db:"user_id" gorm:"type:uuid;primaryKey"`
	ProgressPercentage float64         `db:"progress_percentage" gorm:"column:progress_percentage;type:numeric(5,2);not null;default:0"`
	LastPosition       json.RawMessage `db:"last_position" gorm:"column:last_position;type:jsonb;default:'{}'"`
	CompletedAt        *time.Time      `db:"completed_at" gorm:"column:completed_at"`
	CreatedAt          time.Time       `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL.
func (Progress) TableName() string {
	return "content.progress"
}
