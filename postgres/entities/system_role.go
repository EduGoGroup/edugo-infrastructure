package entities

import (
	"time"

	"github.com/google/uuid"
)

// SystemRole representa la tabla 'system_roles' en el schema iam.
//
// Tabla puente sistema<->rol (MP-08): que roles IAM entran a cada sistema/app
// (iam.systems). Reemplaza la enumeracion hardcodeada de acceso por app. El par
// (system_id, iam_role_id) es unico para no duplicar la relacion.
//
// GORM no materializa las FKs desde el tag `constraint:` sin campo de relacion
// (mismo caso que role_grants/school_invitations), por eso ambas FKs
// (system_id -> iam.systems, iam_role_id -> iam.roles, ambas same-schema iam)
// se declaran en sql/post_gorm.sql espejando los nombres de constraint. El
// uniqueIndex compuesto si lo materializa GORM desde el tag.
type SystemRole struct {
	ID        uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SystemID  uuid.UUID `db:"system_id" gorm:"type:uuid;index;not null;constraint:system_roles_system_fkey,OnDelete:CASCADE;uniqueIndex:uq_system_roles_system_role" validate:"required,uuid"`
	IAMRoleID uuid.UUID `db:"iam_role_id" gorm:"type:uuid;index;not null;constraint:system_roles_role_fkey,OnDelete:CASCADE;uniqueIndex:uq_system_roles_system_role" validate:"required,uuid"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SystemRole) TableName() string {
	return "iam.system_roles"
}
