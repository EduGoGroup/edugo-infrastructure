package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialFile representa la tabla 'content.material_file' en PostgreSQL
// (rediseño F2 plan 018: maestro-detalle). Es el DETALLE: N archivos por
// tema (content.materials). Reemplaza a la vieja content.material_version
// (que versionaba el único archivo inline) y a las columnas file_* que
// bajaron del maestro.
//
// Decisiones cerradas (f2-spec):
//   - DEC-1: file_name = nombre original del PDF subido, not null (se muestra tal cual).
//   - DEC-3: el orden de presentación sale de created_at; SIN columna sort_order.
//   - DEC-4: SIN status (el status es del tema, en content.materials).
//
// La FK material_id → content.materials(id) ON DELETE CASCADE es same-schema:
// GORM la materializa desde el tag `constraint:` (no requiere post_gorm.sql).
type MaterialFile struct {
	ID            uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MaterialID    uuid.UUID `db:"material_id" gorm:"type:uuid;index;not null;constraint:material_file_material_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	FileURL       string    `db:"file_url" gorm:"type:text;not null;default:''" validate:"omitempty"`
	FileName      string    `db:"file_name" gorm:"type:text;not null" validate:"required"`
	FileType      string    `db:"file_type" gorm:"type:text;not null;default:''" validate:"omitempty"`
	FileSizeBytes int64     `db:"file_size_bytes" gorm:"not null;default:0"`
	CreatedAt     time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialFile) TableName() string {
	return "content.material_file"
}
