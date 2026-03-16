package entities

import (
	"time"

	"github.com/google/uuid"
)

// Progress representa la tabla 'content.progress' en PostgreSQL.
// Reflejo exacto del schema definido en: postgres/migrations/structure/042_content_progress.sql
//
// Clave primaria compuesta (material_id, user_id) — un registro por material/usuario.
// completed_at != NULL indica que el usuario completo el material.
type Progress struct {
	MaterialID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ProgressPercentage float64    `gorm:"column:progress_percentage;type:numeric(5,2);not null;default:0"`
	LastPosition       *string    `gorm:"column:last_position;type:jsonb"`
	CompletedAt        *time.Time `gorm:"column:completed_at"`
	CreatedAt          time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL.
func (Progress) TableName() string {
	return "content.progress"
}
