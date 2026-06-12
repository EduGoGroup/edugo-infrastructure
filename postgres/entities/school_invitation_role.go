package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolInvitationRole representa la tabla 'school_invitation_roles' en el
// schema academic.
//
// Equivalencia por escuela (MP-08): (school, tipo de invitacion) -> rol IAM.
// Modela en datos que rol concreto otorga cada tipo de invitacion dentro de una
// escuela. El par (school_id, invitation_type_id) es unico para no duplicar la
// equivalencia.
//
// GORM no materializa las FKs desde el tag `constraint:` sin campo de relacion
// (mismo caso que school_invitations/subject_offerings), por eso TODAS las FKs
// viven en sql/post_gorm.sql: school_id -> academic.schools (same-schema) e
// invitation_type_id -> academic.invitation_types (same-schema) van con los
// otros constraints academic; iam_role_id -> iam.roles es CROSS-SCHEMA
// (academic -> iam) y por eso obligatoriamente fuera del tag. El uniqueIndex
// compuesto si lo materializa GORM desde el tag.
type SchoolInvitationRole struct {
	ID               uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID         uuid.UUID `db:"school_id" gorm:"type:uuid;index;not null;constraint:school_invitation_roles_school_fkey,OnDelete:CASCADE;uniqueIndex:uq_school_invitation_roles_school_type" validate:"required,uuid"`
	InvitationTypeID uuid.UUID `db:"invitation_type_id" gorm:"type:uuid;index;not null;constraint:school_invitation_roles_type_fkey,OnDelete:CASCADE;uniqueIndex:uq_school_invitation_roles_school_type" validate:"required,uuid"`
	IAMRoleID        uuid.UUID `db:"iam_role_id" gorm:"type:uuid;index;not null" validate:"required,uuid"`
	CreatedAt        time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SchoolInvitationRole) TableName() string {
	return "academic.school_invitation_roles"
}
