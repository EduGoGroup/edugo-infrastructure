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
	ID            uuid.UUID `db:"id"`
	MaterialID    uuid.UUID `db:"material_id"`
	VersionNumber int       `db:"version_number"`
	Title         string    `db:"title"`
	ContentURL    string    `db:"content_url"`
	ChangedBy     uuid.UUID `db:"changed_by"`
	CreatedAt     time.Time `db:"created_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialVersion) TableName() string {
	return "material_versions"
}
