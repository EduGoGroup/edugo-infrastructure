package entities

import (
	"time"

	"github.com/google/uuid"
)

// RoleGrant representa un grant de permiso asociado a un rol en iam.role_grants.
// Reemplaza el mapping role_permissions+permissions por patrones glob (P1-1).
type RoleGrant struct {
	ID        uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	RoleID    uuid.UUID `db:"role_id" gorm:"type:uuid;not null;constraint:fk_role_grants_role,OnDelete:CASCADE;uniqueIndex:uq_role_grants_role_pattern_effect;index:idx_role_grants_role_effect" validate:"required,uuid"`
	Pattern   string    `db:"pattern" gorm:"type:varchar(150);not null;uniqueIndex:uq_role_grants_role_pattern_effect" validate:"required,max=150"`
	Effect    string    `db:"effect" gorm:"type:varchar(10);not null;default:allow;uniqueIndex:uq_role_grants_role_pattern_effect;index:idx_role_grants_role_effect" validate:"required,oneof=allow deny"`
	// NOTE: CHECK constraints role_grants_pattern_format y role_grants_effect_format se crean en post_gorm.sql.
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (RoleGrant) TableName() string {
	return "iam.role_grants"
}
