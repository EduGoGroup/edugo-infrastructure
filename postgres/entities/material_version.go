package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialVersion representa la tabla 'material_versions' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/012_create_material_versions.up.sql
//
// Almacena versiones históricas de materiales educativos.
// Cada vez que se realiza un cambio significativo al material, se crea una nueva versión.
type MaterialVersion struct {
	ID            uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	MaterialID    uuid.UUID `db:"material_id" gorm:"type:uuid;index;not null"`
	VersionNumber int       `db:"version_number" gorm:"not null"`
	Title         string    `db:"title" gorm:"not null"`
	ContentURL    string    `db:"content_url" gorm:"not null"`
	ChangedBy     uuid.UUID `db:"changed_by" gorm:"type:uuid;not null"`
	CreatedAt     time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialVersion) TableName() string {
	return "content.material_versions"
}
