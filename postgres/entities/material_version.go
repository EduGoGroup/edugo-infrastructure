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
	ID            uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MaterialID    uuid.UUID `db:"material_id" gorm:"type:uuid;index;not null;constraint:material_versions_material_fkey,OnDelete:CASCADE;uniqueIndex:material_versions_unique" validate:"required,uuid"`
	VersionNumber int       `db:"version_number" gorm:"not null;uniqueIndex:material_versions_unique" validate:"required"`
	Title         string    `db:"title" gorm:"not null" validate:"required,min=2,max=255"`
	ContentURL    string    `db:"content_url" gorm:"not null" validate:"required,url"`
	ChangedBy     uuid.UUID `db:"changed_by" gorm:"type:uuid;not null;constraint:material_versions_created_by_fkey,OnDelete:SET NULL" validate:"required,uuid"`
	CreatedAt     time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialVersion) TableName() string {
	return "content.material_versions"
}
