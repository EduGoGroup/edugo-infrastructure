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
	MaterialID     uuid.UUID `db:"material_id"`
	UserID         uuid.UUID `db:"user_id"`
	Percentage     int       `db:"percentage"` // 0-100
	LastPage       int       `db:"last_page"`  // >= 0
	Status         string    `db:"status"`     // not_started, in_progress, completed
	LastAccessedAt time.Time `db:"last_accessed_at"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Progress) TableName() string {
	return "progress"
}
