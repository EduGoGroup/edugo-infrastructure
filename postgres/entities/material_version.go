package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialVersion representa la tabla 'content.material_version' en PostgreSQL
// (N4 / ADR 0019). Historial de versiones de un material.
//
// Cambio vs viejo: changed_by (global) → changed_by_membership_id (→academic.memberships).
// FKs (material_id → content.materials CASCADE, changed_by_membership_id →
// academic.memberships RESTRICT) se materializan en post_gorm.sql.
type MaterialVersion struct {
	ID                    uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MaterialID            uuid.UUID `db:"material_id" gorm:"type:uuid;index;not null;constraint:material_version_material_fkey,OnDelete:CASCADE;uniqueIndex:uq_material_version,priority:1" validate:"required,uuid"`
	VersionNumber         int       `db:"version_number" gorm:"not null;uniqueIndex:uq_material_version,priority:2" validate:"required"`
	ContentURL            string    `db:"content_url" gorm:"not null" validate:"required,url"`
	ChangedByMembershipID uuid.UUID `db:"changed_by_membership_id" gorm:"type:uuid;not null;constraint:material_version_membership_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	CreatedAt             time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialVersion) TableName() string {
	return "content.material_version"
}
