package entities

import (
	"time"

	"github.com/google/uuid"
)

// Progress representa la tabla 'progress' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/016_create_progress.up.sql
//
// Representa el progreso de lectura de un material por parte de un usuario.
// La clave primaria es compuesta (material_id, user_id), permitiendo un registro por material/usuario.
type Progress struct {
	MaterialID     uuid.UUID `db:"material_id" gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `db:"user_id" gorm:"type:uuid;primaryKey"`
	Percentage     int       `db:"percentage" gorm:"not null;default:0"`
	LastPage       int       `db:"last_page" gorm:"not null;default:0"`
	Status         string    `db:"status" gorm:"not null;type:varchar(50)"`
	LastAccessedAt time.Time `db:"last_accessed_at" gorm:"not null"`
	CreatedAt      time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Progress) TableName() string {
	return "content.progress"
}
