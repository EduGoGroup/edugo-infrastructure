package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserGrant representa overrides puntuales por usuario en iam.user_grants
// (P1-1 permissions-redesign). Permite allow/deny ad-hoc sobre un
// permission_pattern (ADR 0024 DEC-4: el scope_pattern decorativo se eliminó
// porque el motor de auth nunca lo evaluaba; el scope efectivo vive en el JWT).
type UserGrant struct {
	ID                uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID            uuid.UUID `db:"user_id" gorm:"type:uuid;not null;constraint:fk_user_grants_user,OnDelete:CASCADE;uniqueIndex:uq_user_grants_user_perm_effect" validate:"required,uuid"`
	PermissionPattern string    `db:"permission_pattern" gorm:"type:varchar(150);not null;uniqueIndex:uq_user_grants_user_perm_effect" validate:"required,max=150"`
	Effect            string    `db:"effect" gorm:"type:varchar(10);not null;default:deny;uniqueIndex:uq_user_grants_user_perm_effect" validate:"required,oneof=allow deny"`
	// NOTE: índice parcial idx_user_grants_user_active (WHERE expires_at IS NULL OR expires_at > NOW()) en post_gorm.sql.
	ExpiresAt *time.Time `db:"expires_at" gorm:"default:null"`
	GrantedBy *uuid.UUID `db:"granted_by" gorm:"type:uuid;constraint:fk_user_grants_granted_by,OnDelete:SET NULL" validate:"omitempty,uuid"`
	CreatedAt time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (UserGrant) TableName() string {
	return "iam.user_grants"
}
